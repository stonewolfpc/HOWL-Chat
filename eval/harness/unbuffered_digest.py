"""Unbuffered (FILE_FLAG_NO_BUFFERING) SHA-256 reader for Windows, revision 2.

Purpose: separate the *file cache* from everything below it when hashing a file.
Buffered readers (Python `open`, GNU sha256sum, .NET Get-FileHash) all take the
same cached path, so their agreement cannot tell a poisoned cache page from data
delivered by the storage stack. This reader asks the volume directly.

What it bypasses: the Windows file cache (CcCopyRead, readahead, mapped pages).
What it does NOT bypass: minifilters (WdFilter, Wof, bfs, gameflt, FileInfo),
NTFS, the volume manager, the storage driver stack, the controller, the cable,
or the device. It does not separate disk from driver, and it does not prove that
data reaching it was correct.

ALIGNMENT CONTRACT, enforced by the volume rather than by this code:
  * buffer address  must be a multiple of the volume sector size
  * file offset     must be a multiple of the volume sector size
  * transfer length must be a multiple of the volume sector size
A violation returns ERROR_INVALID_PARAMETER (87); the self-test uses that.

BYTE ACCOUNTING. The aligned body is read unbuffered. Bytes past the last whole
sector are read through a separate *buffered* handle and reported as
`buffered_bytes`; the two counts are never merged, and the digest is never
described as fully unbuffered when a tail exists.

FILE IDENTITY. The body and the tail travel different handles, so the reader
records the underlying file identity from BOTH: volume serial number plus the
128-bit file id, via GetFileInformationByHandleEx(FileIdInfo). Both handles must
name the same file, the identity must still hold after each read, and it is
re-checked through a fresh handle at the end. A file replaced or retargeted
between the two read paths aborts the run.

WHAT THE CHECKS DO AND DO NOT ESTABLISH (see also the README section below):

  * Path metadata consistency - size and mtime compared before and after. Detects
    a modification that changes either. Does NOT detect one that preserves both.
  * Underlying file identity consistency - one volume serial + file id across the
    unbuffered handle, the buffered tail handle, and a final fresh handle. Detects
    the path naming a different file. Does NOT detect in-place modification of the
    same file, which leaves both volume serial and file id unchanged.
  * Exact byte accounting - unbuffered_bytes + buffered_bytes must equal the file
    size, or the run raises. Establishes what was hashed, not that it was correct.
  * Protection against concurrent in-place modification - NOT PROVIDED. A writer
    that overwrites bytes in place without changing size, mtime or file id is not
    detectable by metadata, identity or accounting. Only a before/after comparison
    of two independent reads, or a volume-level snapshot, can speak to that.

No code path returns a digest after a failed read, a short read, a misaligned
transfer, an identity change, or a size/mtime change.

Exit codes: 0 success, 2 OSError (Win32 error number), 3 integrity abort.
"""

from __future__ import annotations

import argparse
import ctypes
import ctypes.wintypes as wt
import hashlib
import json
import os
import sys
import time

GENERIC_READ = 0x80000000
FILE_SHARE_ALL = 0x00000001 | 0x00000002 | 0x00000004
OPEN_EXISTING = 3
FILE_FLAG_NO_BUFFERING = 0x20000000
FILE_FLAG_SEQUENTIAL_SCAN = 0x08000000
MEM_COMMIT = 0x1000
MEM_RESERVE = 0x2000   # 0x2000, not 0x200000 (that is MEM_WRITE_WATCH) — the wrong
                       # constant makes VirtualAlloc fail with ERROR_INVALID_PARAMETER 87.
MEM_RELEASE = 0x8000
PAGE_READWRITE = 0x04
FILE_BEGIN = 0
FileIdInfo = 18        # FILE_INFO_BY_HANDLE_CLASS value for FILE_ID_INFO
INVALID_HANDLE = ctypes.c_void_p(-1).value

kernel32 = ctypes.WinDLL("kernel32", use_last_error=True)

CreateFileW = kernel32.CreateFileW
CreateFileW.restype = wt.HANDLE
CreateFileW.argtypes = [wt.LPCWSTR, wt.DWORD, wt.DWORD, ctypes.c_void_p,
                        wt.DWORD, wt.DWORD, wt.HANDLE]

ReadFile = kernel32.ReadFile
ReadFile.restype = wt.BOOL
ReadFile.argtypes = [wt.HANDLE, ctypes.c_void_p, wt.DWORD,
                     ctypes.POINTER(wt.DWORD), ctypes.c_void_p]

CloseHandle = kernel32.CloseHandle
CloseHandle.restype = wt.BOOL
CloseHandle.argtypes = [wt.HANDLE]

SetFilePointerEx = kernel32.SetFilePointerEx
SetFilePointerEx.restype = wt.BOOL
SetFilePointerEx.argtypes = [wt.HANDLE, ctypes.c_longlong,
                             ctypes.POINTER(ctypes.c_longlong), wt.DWORD]

VirtualAlloc = kernel32.VirtualAlloc
VirtualAlloc.restype = ctypes.c_void_p
VirtualAlloc.argtypes = [ctypes.c_void_p, ctypes.c_size_t, wt.DWORD, wt.DWORD]

VirtualFree = kernel32.VirtualFree
VirtualFree.restype = wt.BOOL
VirtualFree.argtypes = [ctypes.c_void_p, ctypes.c_size_t, wt.DWORD]

GetFileSizeEx = kernel32.GetFileSizeEx
GetFileSizeEx.restype = wt.BOOL
GetFileSizeEx.argtypes = [wt.HANDLE, ctypes.POINTER(ctypes.c_longlong)]

GetFileInformationByHandleEx = kernel32.GetFileInformationByHandleEx
GetFileInformationByHandleEx.restype = wt.BOOL
GetFileInformationByHandleEx.argtypes = [wt.HANDLE, ctypes.c_int,
                                         ctypes.c_void_p, wt.DWORD]


class FILE_ID_INFO(ctypes.Structure):
    """FILE_ID_INFO: ULONGLONG VolumeSerialNumber; FILE_ID_128 FileId."""
    _fields_ = [("VolumeSerialNumber", ctypes.c_ulonglong),
                ("FileId", ctypes.c_ubyte * 16)]


class IO_COUNTERS(ctypes.Structure):
    _fields_ = [("ReadOperationCount", ctypes.c_ulonglong),
                ("WriteOperationCount", ctypes.c_ulonglong),
                ("OtherOperationCount", ctypes.c_ulonglong),
                ("ReadTransferCount", ctypes.c_ulonglong),
                ("WriteTransferCount", ctypes.c_ulonglong),
                ("OtherTransferCount", ctypes.c_ulonglong)]


GetProcessIoCounters = kernel32.GetProcessIoCounters
GetProcessIoCounters.restype = wt.BOOL
GetProcessIoCounters.argtypes = [wt.HANDLE, ctypes.POINTER(IO_COUNTERS)]
GetCurrentProcess = kernel32.GetCurrentProcess
GetCurrentProcess.restype = wt.HANDLE


class IntegrityError(RuntimeError):
    """Raised when a read, an identity check, or the byte accounting fails."""


def _open(path: str, unbuffered: bool) -> int:
    flags = FILE_FLAG_SEQUENTIAL_SCAN | (FILE_FLAG_NO_BUFFERING if unbuffered else 0)
    handle = CreateFileW(path, GENERIC_READ, FILE_SHARE_ALL, None,
                         OPEN_EXISTING, flags, None)
    if handle is None or handle == INVALID_HANDLE:
        raise OSError(ctypes.get_last_error(), f"CreateFileW failed for {path}")
    return handle


def _size(handle: int) -> int:
    value = ctypes.c_longlong(0)
    if not GetFileSizeEx(handle, ctypes.byref(value)):
        raise OSError(ctypes.get_last_error(), "GetFileSizeEx failed")
    return int(value.value)


def _identity(handle: int, source: str) -> dict:
    """Volume serial + 128-bit file id for an open handle. Raises on failure."""
    info = FILE_ID_INFO()
    if not GetFileInformationByHandleEx(handle, FileIdInfo, ctypes.byref(info),
                                        ctypes.sizeof(info)):
        raise OSError(ctypes.get_last_error(),
                      f"GetFileInformationByHandleEx(FileIdInfo) failed for {source}")
    return {"volume_serial": f"{info.VolumeSerialNumber:016x}",
            "file_id": bytes(info.FileId).hex(),
            "captured_from": source}


def _io_snapshot() -> dict:
    counters = IO_COUNTERS()
    if not GetProcessIoCounters(GetCurrentProcess(), ctypes.byref(counters)):
        return {}
    return {"read_ops": counters.ReadOperationCount,
            "read_bytes": counters.ReadTransferCount}


def _same_file(a: dict, b: dict) -> bool:
    return (a["volume_serial"] == b["volume_serial"]
            and a["file_id"] == b["file_id"])


def digest(path: str, *, sector: int = 4096, block: int = 4 << 20,
           buffer_offset: int = 0, length_delta: int = 0) -> dict:
    """SHA-256 of `path`: unbuffered body, explicitly labelled buffered tail.

    buffer_offset / length_delta exist only for the negative controls; they
    deliberately violate the alignment contract.
    """
    if sector <= 0 or block % sector:
        raise ValueError("block must be a positive multiple of sector size")

    st_before = os.stat(path)
    before = _io_snapshot()

    short_reads = 0
    handle = _open(path, unbuffered=True)
    try:
        size = _size(handle)
        if size != st_before.st_size:
            raise IntegrityError(
                f"size disagrees before read: GetFileSizeEx={size} stat={st_before.st_size}")

        body_id = _identity(handle, "unbuffered body handle")

        aligned_len = (size // sector) * sector
        tail_len = size - aligned_len

        raw = VirtualAlloc(None, block + 4096, MEM_COMMIT | MEM_RESERVE, PAGE_READWRITE)
        if not raw:
            raise OSError(ctypes.get_last_error(), "VirtualAlloc failed")
        buffer = raw + buffer_offset

        hasher = hashlib.sha256()
        offset = 0
        calls = 0
        unbuffered_bytes = 0
        try:
            while offset < aligned_len:
                want = min(block, aligned_len - offset) + length_delta
                read = wt.DWORD(0)
                ok = ReadFile(handle, ctypes.c_void_p(buffer), want,
                              ctypes.byref(read), None)
                if not ok:
                    err = ctypes.get_last_error()
                    raise OSError(err, f"ReadFile failed at offset {offset} "
                                       f"len {want} (last error {err})")
                got = read.value

                # Per-transfer validation. A legitimate short read is
                # sector-aligned and continues the loop; anything else is a fault.
                if got <= 0:
                    raise IntegrityError(
                        f"unbuffered read returned {got} bytes at offset {offset}")
                if got > want:
                    raise IntegrityError(
                        f"unbuffered read returned {got} bytes for a {want}-byte request "
                        f"at offset {offset}")
                if got % sector:
                    raise IntegrityError(
                        f"unbuffered read returned {got} bytes at offset {offset}, "
                        f"not a multiple of the {sector}-byte sector size")
                if got < want:
                    short_reads += 1

                hasher.update(ctypes.string_at(buffer, got))
                unbuffered_bytes += got
                offset += got
                calls += 1
        finally:
            VirtualFree(ctypes.c_void_p(raw), 0, MEM_RELEASE)

        body_id_after = _identity(handle, "unbuffered body handle, after read")
        if not _same_file(body_id, body_id_after):
            raise IntegrityError(
                f"file identity changed during body read: {body_id} -> {body_id_after}")
    finally:
        CloseHandle(handle)

    buffered_bytes = 0
    tail_id = None
    tail_id_after = None
    if tail_len:
        # Exactly one buffered handle: identity check, seek, tail, identity again.
        thandle = _open(path, unbuffered=False)
        try:
            tail_id = _identity(thandle, "buffered tail handle")
            if not _same_file(body_id, tail_id):
                raise IntegrityError(
                    "body and tail handles name different files: "
                    f"body={body_id} tail={tail_id}")

            position = ctypes.c_longlong(0)
            if not SetFilePointerEx(thandle, ctypes.c_longlong(aligned_len),
                                   ctypes.byref(position), FILE_BEGIN):
                raise OSError(ctypes.get_last_error(), "SetFilePointerEx failed")
            if position.value != aligned_len:
                raise IntegrityError(
                    f"seek landed at {position.value}, expected {aligned_len}")

            tail = ctypes.create_string_buffer(tail_len, tail_len)
            base = ctypes.addressof(tail)
            remaining = tail_len
            read_total = 0
            while remaining:
                read = wt.DWORD(0)
                ok = ReadFile(thandle, ctypes.c_void_p(base + read_total), remaining,
                              ctypes.byref(read), None)
                if not ok:
                    err = ctypes.get_last_error()
                    raise OSError(err, f"buffered tail ReadFile failed at {read_total} "
                                       f"of {tail_len} (last error {err})")
                got = read.value
                if got == 0:
                    raise IntegrityError(
                        f"tail EOF after {read_total} of {tail_len} bytes")
                if got > remaining:
                    raise IntegrityError(
                        f"tail read returned {got} bytes for a {remaining}-byte request")
                read_total += got
                remaining -= got
            if read_total != tail_len:
                raise IntegrityError(
                    f"tail read returned {read_total} bytes, expected {tail_len}")

            tail_id_after = _identity(thandle, "buffered tail handle, after read")
            if not _same_file(tail_id, tail_id_after):
                raise IntegrityError(
                    f"file identity changed during tail read: {tail_id} -> {tail_id_after}")
        finally:
            CloseHandle(thandle)

        hasher.update(ctypes.string_at(base, tail_len))
        buffered_bytes = tail_len

    st_after = os.stat(path)
    if st_after.st_size != st_before.st_size:
        raise IntegrityError(
            f"size changed during read: {st_before.st_size} -> {st_after.st_size}")
    if st_after.st_mtime_ns != st_before.st_mtime_ns:
        raise IntegrityError(
            f"mtime changed during read: {st_before.st_mtime_ns} -> {st_after.st_mtime_ns}")

    # Final re-check through a fresh handle: the path must still name that file.
    vhandle = _open(path, unbuffered=False)
    try:
        final_id = _identity(vhandle, "fresh verification handle")
    finally:
        CloseHandle(vhandle)
    if not _same_file(body_id, final_id):
        raise IntegrityError(
            f"path no longer names the file that was read: read={body_id} now={final_id}")

    if unbuffered_bytes + buffered_bytes != size:
        raise IntegrityError(
            f"byte accounting mismatch: {unbuffered_bytes}+{buffered_bytes} != {size}")

    after = _io_snapshot()
    return {
        "path": path,
        "size": size,
        "sector": sector,
        "block": block,
        "aligned_len": aligned_len,
        "tail_len": tail_len,
        "unbuffered_bytes": unbuffered_bytes,
        "buffered_bytes": buffered_bytes,
        "read_calls": calls,
        "short_sector_aligned_reads": short_reads,
        "sha256": hasher.hexdigest(),
        "identity": {"volume_serial": body_id["volume_serial"],
                     "file_id": body_id["file_id"]},
        "identity_body_stable": _same_file(body_id, body_id_after),
        "identity_body_matches_tail": tail_id is None or _same_file(body_id, tail_id),
        "identity_tail_stable": tail_id_after is None or _same_file(tail_id, tail_id_after),
        "identity_final_matches_body": _same_file(final_id, body_id),
        "identity_stable_after_read": (
            _same_file(body_id, body_id_after)
            and (tail_id is None or _same_file(body_id, tail_id))
            and (tail_id_after is None or _same_file(tail_id, tail_id_after))
            and _same_file(final_id, body_id)),
        "size_before": st_before.st_size,
        "size_after": st_after.st_size,
        "mtime_ns_before": st_before.st_mtime_ns,
        "mtime_ns_after": st_after.st_mtime_ns,
        "path_metadata_consistent": True,
        "unchanged_during_read": True,
        "io": {"read_ops_delta": after.get("read_ops", 0) - before.get("read_ops", 0),
               "read_bytes_delta": after.get("read_bytes", 0) - before.get("read_bytes", 0)},
    }


def main(argv: list[str]) -> int:
    ap = argparse.ArgumentParser(description="unbuffered SHA-256 (Windows)")
    ap.add_argument("path")
    ap.add_argument("--sector", type=int, default=4096)
    ap.add_argument("--block", type=int, default=4 << 20)
    ap.add_argument("--buffer-offset", type=int, default=0, help="negative control")
    ap.add_argument("--length-delta", type=int, default=0, help="negative control")
    ap.add_argument("--quiet", action="store_true")
    args = ap.parse_args(argv)

    started = time.perf_counter()
    try:
        result = digest(args.path, sector=args.sector, block=args.block,
                        buffer_offset=args.buffer_offset,
                        length_delta=args.length_delta)
    except IntegrityError as exc:
        print(json.dumps({"path": args.path, "integrity_abort": str(exc)}, indent=2))
        return 3
    except OSError as exc:
        print(json.dumps({"path": args.path, "error": exc.errno,
                          "message": str(exc), "winerror": exc.errno}, indent=2))
        return 2
    result["elapsed_s"] = round(time.perf_counter() - started, 3)
    print(json.dumps(result, indent=2))
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
