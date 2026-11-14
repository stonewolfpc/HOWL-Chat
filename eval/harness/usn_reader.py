"""Bounded USN-journal reader — implementation of the frozen v6.9 specification.

This module answers a narrower question than it might appear to:

    Which documented-format change-journal records matching a complete
    documented-reason mask were returned over a cursor-completed traversal,
    and were they well formed?

It does **not** claim that every journal record in the range was returned. The
API's ReasonMask is a selection filter, so a record whose only reason bits lie
outside the mask is legitimately omitted and invisible from inside the read
stream. "Extent-checked" and "identity-matched" describe validation of the
records that were returned; they establish neither that every relevant record
was returned nor that the journal retained every historical event.

Consequences that the code enforces rather than merely documents:

* Claim R' requires a completed cursor traversal **and** valid records **and**
  unchanged journal retention across the scan. Any one failing yields a named
  INCONCLUSIVE outcome.
* Claim T' (temporal reach) is never produced here. V2/V3 records carry
  timestamps; V4 records do not, and no timestamp is invented for them.
* An absent matching record is reported as "no matching record returned",
  never as evidence that no modification occurred.
* No volume is opened at import time, and nothing in this module can open one
  unless `volume_io()` is called with a raw volume path.

Authority: MEMORY-09 WO-ARCH-06 §1 (implementation phase), controlling
specification USN Reader v6.9. The three live-volume operations remain separate
authorisation gates: G: disposable-file validation, event-log enumeration, and
the bounded S1-B journal query.
"""

from __future__ import annotations

import ctypes
import math
import struct
from dataclasses import dataclass, field
from typing import Callable, Iterable, Sequence

# --------------------------------------------------------------------------
# Constants (v6.9 §1)
# --------------------------------------------------------------------------

#: Every documented USN reason flag, OR-ed together. Reserved bits are left
#: clear: 0x00000008, 0x00000080 and 0x01000000..0x40000000 are not defined by
#: the API, and guessing at them in a parameter that decides which evidence
#: exists would be unacceptable.
COMPLETE_REASON_MASK = 0x80FFFF77

REASON_FLAGS: dict[str, int] = {
    "DATA_OVERWRITE": 0x00000001,
    "DATA_EXTEND": 0x00000002,
    "DATA_TRUNCATION": 0x00000004,
    "NAMED_DATA_OVERWRITE": 0x00000010,
    "NAMED_DATA_EXTEND": 0x00000020,
    "NAMED_DATA_TRUNCATION": 0x00000040,
    "FILE_CREATE": 0x00000100,
    "FILE_DELETE": 0x00000200,
    "EA_CHANGE": 0x00000400,
    "SECURITY_CHANGE": 0x00000800,
    "RENAME_OLD_NAME": 0x00001000,
    "RENAME_NEW_NAME": 0x00002000,
    "INDEXABLE_CHANGE": 0x00004000,
    "BASIC_INFO_CHANGE": 0x00008000,
    "HARD_LINK_CHANGE": 0x00010000,
    "COMPRESSION_CHANGE": 0x00020000,
    "ENCRYPTION_CHANGE": 0x00040000,
    "OBJECT_ID_CHANGE": 0x00080000,
    "REPARSE_POINT_CHANGE": 0x00100000,
    "STREAM_CHANGE": 0x00200000,
    "TRANSACTED_CHANGE": 0x00400000,
    "INTEGRITY_CHANGE": 0x00800000,
    "CLOSE": 0x80000000,
}

#: Bits the API does not define. A record carrying one invalidates the
#: completeness premise for the whole scan (v6.9 §3, assertion 2).
RESERVED_REASON_BITS = 0x7F000088

#: Documented SourceInfo flags (v6.9 §2). DATA_MANAGEMENT and AUXILIARY_DATA
#: are both cases where a reason being logged does not imply the
#: application-visible bytes changed.
SOURCE_INFO_FLAGS: dict[str, int] = {
    "DATA_MANAGEMENT": 0x00000001,
    "AUXILIARY_DATA": 0x00000002,
    "REPLICATION_MANAGEMENT": 0x00000004,
    "CLIENT_REPLICATION_MANAGEMENT": 0x00000008,
}

#: Reasons that constitute a class-A data change for this investigation.
CLASS_A_DATA_REASONS = (
    REASON_FLAGS["DATA_OVERWRITE"] | REASON_FLAGS["DATA_EXTEND"] | REASON_FLAGS["DATA_TRUNCATION"]
)

MIN_MAJOR_VERSION = 2
MAX_MAJOR_VERSION = 4
SUPPORTED_VERSIONS = (2, 3, 4)

DEFAULT_BUFFER_BYTES = 1024 * 1024

#: The 8-byte next-USN cursor that prefixes every read buffer.
CURSOR_BYTES = 8

#: Header sizes, excluding the variable-length name / extent payload.
HEADER_BYTES = {2: 60, 3: 76, 4: 64}

#: sizeof(USN_RECORD_EXTENT): two signed 64-bit fields.
EXTENT_BYTES = 16

# --------------------------------------------------------------------------
# Win32 control codes and structures
# --------------------------------------------------------------------------
#
# Hoisted to module level so their layout and arithmetic can be asserted by the
# offline suite. Defining them here does not open anything; only `volume_io()`
# can do that, and only under its own gate.

FILE_DEVICE_FILE_SYSTEM = 0x00000009
METHOD_BUFFERED = 0
METHOD_NEITHER = 3

#: CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 61, METHOD_BUFFERED, FILE_ANY_ACCESS)
FSCTL_QUERY_USN_JOURNAL = 0x000900F4
#: CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 46, METHOD_NEITHER, FILE_ANY_ACCESS)
FSCTL_READ_USN_JOURNAL = 0x000900BB


def ctl_code(function: int, method: int, access: int = 0) -> int:
    """Recompute a control code, so the constants above can be checked."""
    return (FILE_DEVICE_FILE_SYSTEM << 16) | (access << 14) | (function << 2) | method


class USN_JOURNAL_DATA_V0(ctypes.Structure):
    """Output structure of FSCTL_QUERY_USN_JOURNAL."""

    _fields_ = [
        ("UsnJournalID", ctypes.c_ulonglong),
        ("FirstUsn", ctypes.c_longlong),
        ("NextUsn", ctypes.c_longlong),
        ("LowestValidUsn", ctypes.c_longlong),
        ("MaxUsn", ctypes.c_longlong),
        ("MaximumSize", ctypes.c_ulonglong),
        ("AllocationDelta", ctypes.c_ulonglong),
    ]


class READ_USN_JOURNAL_DATA_V1(ctypes.Structure):
    """Input structure of FSCTL_READ_USN_JOURNAL (the V1 layout)."""

    _fields_ = [
        ("StartUsn", ctypes.c_longlong),
        ("ReasonMask", ctypes.c_uint32),
        ("ReturnOnlyOnClose", ctypes.c_uint32),
        ("Timeout", ctypes.c_ulonglong),
        ("BytesToWaitFor", ctypes.c_ulonglong),
        ("UsnJournalID", ctypes.c_ulonglong),
        ("MinMajorVersion", ctypes.c_uint16),
        ("MaxMajorVersion", ctypes.c_uint16),
    ]

# --------------------------------------------------------------------------
# Outcomes (v6.9 §1.4, §4)
# --------------------------------------------------------------------------
#
# Every failure path is named. No path converts a failure into a negative
# result about the file, and no path retries.

#: The scan completed its bound and every returned record was valid.
OUTCOME_OK = "OK"

#: The read returned no matching records and left the cursor where it was.
#: This is documented behaviour with BytesToWaitFor = 0, but the cursor never
#: reached the fixed boundary, so it is neither a completed traversal nor a
#: statement that the range holds nothing.
OUTCOME_NO_MATCH_AT_CURSOR = "NO_MATCH_AT_CURSOR"

OUTCOME_NOT_ATTEMPTED = "NOT_ATTEMPTED"

I_JOURNAL_ID_CHANGED = "INCONCLUSIVE_JOURNAL_ID_CHANGED"
I_RETENTION_ADVANCED = "INCONCLUSIVE_RETENTION_ADVANCED"
I_NEXT_USN_REGRESSED = "INCONCLUSIVE_NEXT_USN_REGRESSED"
#: The three journal-status failures are distinct conditions with distinct
#: meanings, and are never collapsed into one another:
#:
#:   1178 (0x49A) ERROR_JOURNAL_DELETE_IN_PROGRESS — the journal is being deleted;
#:                                                  the current operation cannot proceed
#:   1179 (0x49B) ERROR_JOURNAL_NOT_ACTIVE         — the journal is not active;
#:                                                  the current operation cannot proceed
#:   1181 (0x49D) ERROR_JOURNAL_ENTRY_DELETED       — the requested entry has been
#:                                                  deleted from the journal
#:
#: Source: "System Error Codes (1000-1299) (WinError.h)", Microsoft Learn,
#: https://learn.microsoft.com/en-us/windows/win32/debug/system-error-codes--1000-1299-
I_JOURNAL_DELETE_IN_PROGRESS = "INCONCLUSIVE_JOURNAL_DELETE_IN_PROGRESS"
I_JOURNAL_ENTRY_DELETED = "INCONCLUSIVE_JOURNAL_ENTRY_DELETED"
I_JOURNAL_INACTIVE = "INCONCLUSIVE_JOURNAL_INACTIVE"
I_NO_PROGRESS_WITH_RECORDS = "INCONCLUSIVE_NO_PROGRESS_WITH_RECORDS"
I_ADVANCED_WITHOUT_RECORDS = "INCONCLUSIVE_ADVANCED_WITHOUT_RECORDS"
I_CURSOR_REGRESSED = "INCONCLUSIVE_CURSOR_REGRESSED"
I_UNSUPPORTED_VERSION = "INCONCLUSIVE_UNSUPPORTED_VERSION"
I_UNDOCUMENTED_REASON = "INCONCLUSIVE_UNDOCUMENTED_REASON"
I_MALFORMED_RECORD = "INCONCLUSIVE_MALFORMED_RECORD"
I_MASK_MISMATCH = "INCONCLUSIVE_MASK_MISMATCH"
I_BUDGET_CALLS = "INCONCLUSIVE_BUDGET_CALLS"
I_BUDGET_BYTES = "INCONCLUSIVE_BUDGET_BYTES"
I_BUDGET_PREFLIGHT = "INCONCLUSIVE_BUDGET_FEASIBILITY_NOT_PROVEN"
I_BAD_INPUT = "INCONCLUSIVE_INVALID_QUERY_RESULT"

INCONCLUSIVE_OUTCOMES = frozenset(
    {
        I_JOURNAL_ID_CHANGED,
        I_RETENTION_ADVANCED,
        I_NEXT_USN_REGRESSED,
        I_JOURNAL_DELETE_IN_PROGRESS,
        I_JOURNAL_ENTRY_DELETED,
        I_JOURNAL_INACTIVE,
        I_NO_PROGRESS_WITH_RECORDS,
        I_ADVANCED_WITHOUT_RECORDS,
        I_CURSOR_REGRESSED,
        I_UNSUPPORTED_VERSION,
        I_UNDOCUMENTED_REASON,
        I_MALFORMED_RECORD,
        I_MASK_MISMATCH,
        I_BUDGET_CALLS,
        I_BUDGET_BYTES,
        I_BUDGET_PREFLIGHT,
        I_BAD_INPUT,
    }
)

#: Win32 statuses as documented in WinError.h (System Error Codes 1000-1299).
WINERROR_JOURNAL_DELETE_IN_PROGRESS = 1178
WINERROR_JOURNAL_NOT_ACTIVE = 1179
WINERROR_JOURNAL_ENTRY_DELETED = 1181

#: Canonical diagnostic text per capability outcome. A report that says only
#: "INCONCLUSIVE" hides which condition occurred; these strings are what the
#: log carries, and they quote the documented status meaning.
OUTCOME_DESCRIPTIONS: dict[str, str] = {
    I_JOURNAL_DELETE_IN_PROGRESS: (
        "1178 ERROR_JOURNAL_DELETE_IN_PROGRESS: the volume change journal is being "
        "deleted; the current operation cannot proceed"
    ),
    I_JOURNAL_INACTIVE: (
        "1179 ERROR_JOURNAL_NOT_ACTIVE: the volume change journal is not active; "
        "the current operation cannot proceed"
    ),
    I_JOURNAL_ENTRY_DELETED: (
        "1181 ERROR_JOURNAL_ENTRY_DELETED: the requested journal entry has been "
        "deleted from the journal"
    ),
}

#: Range-tracking (V4) sequence outcomes (v6.9 §1).
SEQ_CONSISTENT = "CONSISTENT"
SEQ_BOUNDARY_TRUNCATED = "BOUNDARY_TRUNCATED"
SEQ_INVALID = "INVALID_INCOMPLETE"

V4_SEQUENCE_OUTCOMES = frozenset({SEQ_CONSISTENT, SEQ_BOUNDARY_TRUNCATED, SEQ_INVALID})


# --------------------------------------------------------------------------
# Bounds (v6.9 §4)
# --------------------------------------------------------------------------


@dataclass(frozen=True)
class Caps:
    """Hard limits, enforced by the reader rather than by its caller."""

    max_calls: int = 128
    max_bytes: int = 128 * 1024 * 1024
    max_output_bytes: int = 4 * 1024 * 1024
    buffer_bytes: int = DEFAULT_BUFFER_BYTES

    def __post_init__(self) -> None:
        if self.max_calls <= 0 or self.max_bytes <= 0 or self.max_output_bytes <= 0:
            raise ValueError("caps must be positive")
        if self.buffer_bytes <= CURSOR_BYTES + HEADER_BYTES[4]:
            raise ValueError("buffer must be able to hold the cursor and one record header")


# --------------------------------------------------------------------------
# Journal state
# --------------------------------------------------------------------------


@dataclass(frozen=True)
class JournalState:
    """The four values captured by FSCTL_QUERY_USN_JOURNAL.

    `FirstUsn` is the first record that can currently be read; `LowestValidUsn`
    is the first record written for the current journal instance. Their
    inequality documents past truncation and is never used as a pass/fail test.
    """

    journal_id: int
    first_usn: int
    next_usn: int
    lowest_valid_usn: int
    maximum_size: int = 0

    @property
    def retained_bytes(self) -> int:
        return max(0, self.next_usn - self.first_usn)

    @property
    def is_empty_range(self) -> bool:
        return self.next_usn <= self.first_usn


def parse_journal_state(payload: Sequence[int] | dict) -> JournalState:
    """Build a JournalState from the ioctl's five USN_JOURNAL_DATA fields.

    Accepts either a 5-tuple in structure order or a mapping with those names.
    Unknown input is rejected rather than guessed at.
    """
    if isinstance(payload, dict):
        try:
            return JournalState(
                journal_id=int(payload["usn_journal_id"]),
                first_usn=int(payload["first_usn"]),
                next_usn=int(payload["next_usn"]),
                lowest_valid_usn=int(payload["lowest_valid_usn"]),
                maximum_size=int(payload.get("maximum_size", 0)),
            )
        except (KeyError, TypeError, ValueError) as exc:
            raise ValueError(f"unrecognised journal state: {exc}") from exc
    try:
        values = list(payload)
    except TypeError as exc:
        raise ValueError("journal state must be a mapping or a sequence") from exc
    if len(values) < 4:
        raise ValueError("journal state needs at least four fields")
    return JournalState(
        journal_id=int(values[0]),
        first_usn=int(values[1]),
        next_usn=int(values[2]),
        lowest_valid_usn=int(values[3]),
        maximum_size=int(values[4]) if len(values) > 4 else 0,
    )


@dataclass(frozen=True)
class Preflight:
    """What can be said *before* a single byte is read (v6.9 §4)."""

    initial_span_bytes: int
    payload_per_call_bytes: int
    estimated_calls: int
    budget_appears_sufficient: bool

    @property
    def completion_guaranteed(self) -> bool:
        # Deliberately always False. The estimate cannot see journal growth
        # during the scan, nor truncation mid-scan.
        return False

    def to_dict(self) -> dict:
        return {
            "initial_span_bytes": self.initial_span_bytes,
            "payload_per_call_bytes": self.payload_per_call_bytes,
            "estimated_calls": self.estimated_calls,
            "budget_appears_sufficient": self.budget_appears_sufficient,
            "completion_guaranteed": self.completion_guaranteed,
            "statement": (
                "budget appears sufficient (estimate {n} calls); completion not guaranteed".format(
                    n=self.estimated_calls
                )
                if self.budget_appears_sufficient
                else "budget feasibility not proven"
            ),
        }


def preflight(state: JournalState, caps: Caps | None = None) -> Preflight:
    caps = caps or Caps()
    payload = max(1, caps.buffer_bytes - CURSOR_BYTES)
    span = state.retained_bytes
    estimated = math.ceil(span / payload) if span else 0
    # Two calls of margin for the final over-read and unavoidable boundary slack.
    feasible = span <= caps.max_bytes and estimated <= max(0, caps.max_calls - 2)
    return Preflight(
        initial_span_bytes=span,
        payload_per_call_bytes=payload,
        estimated_calls=estimated,
        budget_appears_sufficient=feasible,
    )


# --------------------------------------------------------------------------
# Records
# --------------------------------------------------------------------------


@dataclass(frozen=True)
class UsnRecord:
    """One parsed journal record, with raw bytes preserved."""

    usn: int
    version: int
    minor_version: int
    file_reference: bytes
    parent_reference: bytes
    reason: int
    source_info: int
    record_length: int
    timestamp: int | None = None
    remaining_extents: int | None = None
    number_of_extents: int | None = None
    extent_size: int | None = None
    extents: tuple[tuple[int, int], ...] = ()

    @property
    def dated(self) -> bool:
        """V4 records carry no timestamp (v6.9 §1, §2)."""
        return self.timestamp is not None

    @property
    def reason_names(self) -> tuple[str, ...]:
        return tuple(name for name, bit in REASON_FLAGS.items() if self.reason & bit)

    @property
    def has_class_a_data_reason(self) -> bool:
        return bool(self.reason & CLASS_A_DATA_REASONS)

    @property
    def is_close(self) -> bool:
        return bool(self.reason & REASON_FLAGS["CLOSE"])

    def frn64(self) -> int:
        """The NTFS file reference number, extracted version-accurately.

        V2 carries an 8-byte FRN. V3/V4 carry a 16-byte FILE_ID_128 whose low
        64 bits are the FRN on NTFS; a non-zero high half is reported by the
        validator rather than silently truncated here.
        """
        if self.version == 2:
            return int.from_bytes(self.file_reference, "little")
        return int.from_bytes(self.file_reference[:8], "little")

    def high_half(self) -> int:
        if self.version == 2:
            return 0
        return int.from_bytes(self.file_reference[8:], "little")


@dataclass(frozen=True)
class RecordProblem:
    code: str
    usn: int | None
    detail: str


def _u16(buf: bytes, off: int) -> int:
    return struct.unpack_from("<H", buf, off)[0]


def _u32(buf: bytes, off: int) -> int:
    return struct.unpack_from("<I", buf, off)[0]


def _u64(buf: bytes, off: int) -> int:
    return struct.unpack_from("<Q", buf, off)[0]


def _i64(buf: bytes, off: int) -> int:
    return struct.unpack_from("<q", buf, off)[0]


def parse_records(buffer: bytes, start_usn: int) -> tuple[int, list[UsnRecord], list[RecordProblem]]:
    """Parse one read buffer.

    Returns (next_usn_cursor, records, problems). The leading 8 bytes are the
    resume cursor and are consumed, never parsed as a record. A record whose
    declared length runs past the buffer, whose length is zero, or whose header
    cannot hold the fields its version requires is reported as malformed — the
    buffer is sized to hold any single record, so a short tail is a defect
    rather than something to carry forward.
    """
    if len(buffer) < CURSOR_BYTES:
        return start_usn, [], [RecordProblem(I_MALFORMED_RECORD, None, "buffer shorter than cursor")]
    cursor = _u64(buffer, 0)
    records: list[UsnRecord] = []
    problems: list[RecordProblem] = []
    offset = CURSOR_BYTES
    while offset < len(buffer):
        remaining = len(buffer) - offset
        if remaining < 8:
            problems.append(
                RecordProblem(I_MALFORMED_RECORD, None, f"{remaining} trailing bytes cannot hold a header")
            )
            break
        length = _u32(buffer, offset)
        version = _u16(buffer, offset + 4)
        minor = _u16(buffer, offset + 6)
        if version not in SUPPORTED_VERSIONS:
            problems.append(
                RecordProblem(I_UNSUPPORTED_VERSION, None, f"record version {version} is not documented")
            )
            break
        header = HEADER_BYTES[version]
        if length < header:
            problems.append(
                RecordProblem(
                    I_MALFORMED_RECORD, None, f"record length {length} shorter than the {header}-byte header"
                )
            )
            break
        if length > remaining:
            problems.append(
                RecordProblem(
                    I_MALFORMED_RECORD, None, f"record length {length} exceeds the {remaining} bytes available"
                )
            )
            break
        chunk = buffer[offset : offset + length]
        try:
            record = _parse_one(chunk, version, minor, length)
        except (struct.error, ValueError) as exc:
            problems.append(RecordProblem(I_MALFORMED_RECORD, None, f"unparsable record: {exc}"))
            break
        if record.usn < start_usn:
            problems.append(
                RecordProblem(
                    I_MALFORMED_RECORD,
                    record.usn,
                    f"record USN {record.usn} precedes the call's start USN {start_usn}",
                )
            )
            break
        records.append(record)
        offset += length
    return cursor, records, problems


def _parse_one(chunk: bytes, version: int, minor: int, length: int) -> UsnRecord:
    if version == 2:
        return UsnRecord(
            usn=_u64(chunk, 24),
            version=2,
            minor_version=minor,
            file_reference=bytes(chunk[8:16]),
            parent_reference=bytes(chunk[16:24]),
            reason=_u32(chunk, 40),
            source_info=_u32(chunk, 44),
            record_length=length,
            timestamp=_u64(chunk, 32),
        )
    if version == 3:
        return UsnRecord(
            usn=_u64(chunk, 40),
            version=3,
            minor_version=minor,
            file_reference=bytes(chunk[8:24]),
            parent_reference=bytes(chunk[24:40]),
            reason=_u32(chunk, 56),
            source_info=_u32(chunk, 60),
            record_length=length,
            timestamp=_u64(chunk, 48),
        )
    # Version 4: range tracking, and no timestamp field at all.
    remaining = _u32(chunk, 56)
    count = _u16(chunk, 60)
    extent_size = _u16(chunk, 62)
    extents: list[tuple[int, int]] = []
    if count and extent_size == EXTENT_BYTES:
        need = count * EXTENT_BYTES
        if length < HEADER_BYTES[4] + need:
            raise ValueError(f"{count} extents need {need} bytes but the record holds {length - HEADER_BYTES[4]}")
        for index in range(count):
            base = HEADER_BYTES[4] + index * EXTENT_BYTES
            extents.append((_i64(chunk, base), _i64(chunk, base + 8)))
    elif count:
        raise ValueError(f"extent size {extent_size} is not sizeof(USN_RECORD_EXTENT)")
    return UsnRecord(
        usn=_u64(chunk, 40),
        version=4,
        minor_version=minor,
        file_reference=bytes(chunk[8:24]),
        parent_reference=bytes(chunk[24:40]),
        reason=_u32(chunk, 48),
        source_info=_u32(chunk, 52),
        record_length=length,
        timestamp=None,
        remaining_extents=remaining,
        number_of_extents=count,
        extent_size=extent_size,
        extents=tuple(extents),
    )


def validate_record(record: UsnRecord, mask: int = COMPLETE_REASON_MASK) -> RecordProblem | None:
    """Structural and mask validation of a single returned record (v6.9 §3)."""
    if record.version not in SUPPORTED_VERSIONS:
        return RecordProblem(I_UNSUPPORTED_VERSION, record.usn, f"version {record.version}")
    if not record.reason & mask:
        # The API should not have returned it; either the mask is incomplete or
        # the record is malformed. Both invalidate the completeness premise.
        return RecordProblem(I_MASK_MISMATCH, record.usn, "returned record has no reason bit inside the mask")
    undocumented = record.reason & ~COMPLETE_REASON_MASK
    if undocumented:
        return RecordProblem(
            I_UNDOCUMENTED_REASON, record.usn, f"reason bits outside the documented set: 0x{undocumented:08X}"
        )
    if record.version in (3, 4) and record.high_half() != 0:
        return RecordProblem(I_MALFORMED_RECORD, record.usn, "FILE_ID_128 high half is not zero on NTFS")
    # Reserved bits are a subset of the undocumented set, so they are already
    # rejected above; RESERVED_REASON_BITS exists to name them and to let the
    # mask constant be checked arithmetically. A separate branch here would be
    # unreachable, and an unreachable check is worse than none: it reads as
    # coverage that the code does not provide.
    return None


# --------------------------------------------------------------------------
# Identity and SourceInfo
# --------------------------------------------------------------------------


@dataclass(frozen=True)
class FileIdentity:
    """The evidence used to attribute a returned record to one file.

    The documented identity is the volume serial together with the file ID; the
    serial alone carries no general uniqueness guarantee.
    """

    frn: int
    volume_serial_64: int
    file_id_128: bytes | None = None

    @classmethod
    def from_file_id(cls, file_id_128: bytes, volume_serial_64: int) -> "FileIdentity":
        if len(file_id_128) != 16:
            raise ValueError("FILE_ID_INFO file id must be 16 bytes")
        return cls(
            frn=int.from_bytes(file_id_128[:8], "little"),
            volume_serial_64=int(volume_serial_64),
            file_id_128=bytes(file_id_128),
        )


def verify_volume_identity(
    serial_64_file: int,
    serial_64_dir: int,
    serial_32_file: int | None = None,
    serial_32_volume: int | None = None,
) -> dict:
    """64-bit handle-to-handle comparison, with the 32-bit values supplementary.

    The 32-bit volume serial is an opaque value that is not required to be
    unique, so agreement between two 32-bit serials is reported as a
    consistency check and never as volume-identity proof.
    """
    identity_established = serial_64_file == serial_64_dir
    consistency: bool | None = None
    if serial_32_file is not None and serial_32_volume is not None:
        consistency = serial_32_file == serial_32_volume
    return {
        "identity_established": identity_established,
        "volume_serial_64_file": serial_64_file,
        "volume_serial_64_dir": serial_64_dir,
        "serial_32_consistency_check": consistency,
        "statement": (
            "64-bit volume serial matches between the file handle and the volume root; "
            "the serial-plus-file-ID pair is the documented identity"
            if identity_established
            else "64-bit volume serial differs between handles: identity not established"
        ),
        "serial_32_is_not_identity_proof": True,
    }


def decode_source_info(mask: int) -> dict:
    """Decode SourceInfo without discarding anything (v6.9 §2)."""
    names = [name for name, bit in SOURCE_INFO_FLAGS.items() if mask & bit]
    unknown = mask & ~sum(SOURCE_INFO_FLAGS.values())
    weight = "ordinary"
    if mask & SOURCE_INFO_FLAGS["DATA_MANAGEMENT"]:
        weight = "data-management: a logged overwrite need not be an application-visible change"
    elif mask & SOURCE_INFO_FLAGS["AUXILIARY_DATA"]:
        weight = "auxiliary-data: documented as not changing application data"
    elif mask & SOURCE_INFO_FLAGS["REPLICATION_MANAGEMENT"]:
        weight = "replication: identify the source; not an ordinary application write"
    elif mask & SOURCE_INFO_FLAGS["CLIENT_REPLICATION_MANAGEMENT"]:
        weight = "client replication: identify the source; not an ordinary application write"
    return {
        "raw": mask,
        "flags": names,
        "unknown_bits": unknown,
        "weight": weight,
        "changes_application_data": not bool(
            mask
            & (
                SOURCE_INFO_FLAGS["DATA_MANAGEMENT"]
                | SOURCE_INFO_FLAGS["AUXILIARY_DATA"]
            )
        ),
    }


# --------------------------------------------------------------------------
# V4 range-tracking sequences (v6.9 §1)
# --------------------------------------------------------------------------


@dataclass
class V4Sequence:
    records: list[UsnRecord] = field(default_factory=list)
    outcome: str = SEQ_INVALID
    detail: str = ""
    companion: UsnRecord | None = None

    @property
    def total_extents(self) -> int:
        return sum(r.number_of_extents or 0 for r in self.records)

    def to_dict(self) -> dict:
        return {
            "outcome": self.outcome,
            "detail": self.detail,
            "records": len(self.records),
            "total_extents": self.total_extents,
            "remaining_extents": [r.remaining_extents for r in self.records],
            "number_of_extents": [r.number_of_extents for r in self.records],
            "undated": True,
            "companion_usn": self.companion.usn if self.companion else None,
            "companion_timestamp": self.companion.timestamp if self.companion else None,
            "companion_timestamp_label": "close-record timestamp (not the modification time)",
        }


def collect_v4_sequence(records: Sequence[UsnRecord], start_index: int) -> tuple[V4Sequence, int]:
    """Group consecutive same-FRN V4 records beginning at `start_index`.

    Returns the sequence and the index of the first record after it. Validation
    applies the exact recurrence: `R[i+1] == R[i] - N[i+1]`, with no underflow,
    a terminating zero, and per-record extent-array agreement. The
    reconciliation sum is computed independently as a parsing cross-check; it is
    not an independent completeness detector, because the recurrence plus the
    terminal zero already imply it.
    """
    seq = V4Sequence()
    index = start_index
    first = records[start_index]
    while (
        index < len(records)
        and records[index].version == 4
        and records[index].frn64() == first.frn64()
    ):
        seq.records.append(records[index])
        index += 1

    detail = "sequence validated"
    seq.outcome = SEQ_CONSISTENT
    for position, record in enumerate(seq.records):
        count = record.number_of_extents or 0
        remaining = record.remaining_extents
        if remaining is None:
            seq.outcome, detail = SEQ_INVALID, "V4 record lacks RemainingExtents"
            break
        if record.extent_size != (EXTENT_BYTES if count else record.extent_size):
            seq.outcome, detail = SEQ_INVALID, "ExtentSize does not match sizeof(USN_RECORD_EXTENT)"
            break
        if len(record.extents) != count:
            seq.outcome, detail = SEQ_INVALID, "extent array length disagrees with NumberOfExtents"
            break
        if position == 0:
            continue
        previous = seq.records[position - 1]
        prev_remaining = previous.remaining_extents or 0
        prev_count = previous.number_of_extents or 0
        if prev_count > prev_remaining:
            seq.outcome, detail = SEQ_INVALID, f"underflow at record {position}"
            break
        if remaining != prev_remaining - count:
            seq.outcome, detail = (
                SEQ_INVALID,
                f"recurrence failure at record {position}: expected {prev_remaining - count}, saw {remaining}",
            )
            break

    if seq.outcome == SEQ_CONSISTENT and (seq.records[-1].remaining_extents or 0) != 0:
        # Not terminated inside the scanned range: the boundary cut the sequence,
        # or records were lost. The reader cannot tell those apart, so it reports
        # the weaker of the two and claims no companion.
        seq.outcome = SEQ_BOUNDARY_TRUNCATED
        detail = "sequence does not terminate inside the scanned range"

    if seq.outcome == SEQ_CONSISTENT:
        # Independent reconciliation, computed from the record counts rather than
        # from the recurrence. The recurrence plus a terminal zero already imply
        # the identity, so disagreement signals an implementation or parsing
        # error; agreement is a consistency check, not proof of completeness.
        first_remaining = seq.records[0].remaining_extents or 0
        later_total = sum((r.number_of_extents or 0) for r in seq.records[1:])
        if first_remaining != later_total:
            seq.outcome, detail = (
                SEQ_INVALID,
                f"reconciliation disagrees: R1={first_remaining}, sum of later extents={later_total}",
            )

    if seq.outcome == SEQ_CONSISTENT:
        companion_index = index
        while companion_index < len(records) and records[companion_index].frn64() != first.frn64():
            companion_index += 1
        if companion_index < len(records):
            candidate = records[companion_index]
            if candidate.version == 3 and candidate.is_close:
                seq.companion = candidate
                detail = "sequence validated; V3 close companion observed"
            else:
                detail = "sequence validated; next same-FRN record is not a V3 close record, so no companion"
    seq.detail = detail
    return seq, index


# --------------------------------------------------------------------------
# Scan
# --------------------------------------------------------------------------

ReadFn = Callable[[int, int], bytes]
QueryFn = Callable[[], JournalState]


@dataclass
class ScanResult:
    outcome: str = OUTCOME_NOT_ATTEMPTED
    detail: str = ""
    claim_r_prime: bool = False
    traversal: dict = field(default_factory=dict)
    validity: dict = field(default_factory=dict)
    gaps: dict = field(default_factory=dict)
    preflight: dict = field(default_factory=dict)
    journal_state_before: JournalState | None = None
    journal_state_after: JournalState | None = None
    records: list[UsnRecord] = field(default_factory=list)
    matches: list[dict] = field(default_factory=list)
    v4_sequences: list[dict] = field(default_factory=list)
    over_read_records: int = 0
    output_summarised: bool = False

    @property
    def ok(self) -> bool:
        return self.outcome == OUTCOME_OK

    def to_dict(self) -> dict:
        return {
            "outcome": self.outcome,
            "detail": self.detail,
            "claim_r_prime": self.claim_r_prime,
            "claim_t_prime": "not established by this reader",
            "traversal": self.traversal,
            "validity": self.validity,
            "gaps": self.gaps,
            "preflight": self.preflight,
            "matches": self.matches,
            "v4_sequences": self.v4_sequences,
            "over_read_records": self.over_read_records,
            "output_summarised": self.output_summarised,
            "completeness": (
                "Claim R' covers the documented V2-V4 sequence under the API's selection "
                "semantics. It does not establish that every journal record in the range "
                "was returned."
            ),
            "absence_means": (
                "No matching record was returned. This is not evidence that the file was "
                "not modified, and says nothing about raw-volume or device-level changes."
            ),
        }


def scan(
    identity: FileIdentity | None,
    query_fn: QueryFn,
    read_fn: ReadFn,
    caps: Caps | None = None,
) -> ScanResult:
    """Run one bounded traversal.

    `query_fn` returns the journal state; it is called once before and once
    after the read loop. `read_fn(start_usn, buffer_bytes)` returns one raw
    buffer. Neither is implemented here, so tests cannot touch a volume and this
    function cannot open one.
    """
    caps = caps or Caps()
    result = ScanResult()

    before = _safe_query(query_fn, result)
    if before is None:
        return result
    result.journal_state_before = before
    plan = preflight(before, caps)
    result.preflight = plan.to_dict()
    if not plan.budget_appears_sufficient:
        result.outcome = I_BUDGET_PREFLIGHT
        result.detail = plan.to_dict()["statement"]
        # No records are read. A failed preflight is a statement about budget,
        # not about the journal.
        _post_query(query_fn, result)
        return result

    boundary = before.next_usn
    cursor = before.first_usn
    calls = 0
    # The cap bounds *requested* bytes: the device read is up to the buffer
    # size whatever comes back, so counting only returned bytes would understate
    # the I/O and let a slow journal evade the budget.
    requested_total = 0
    returned_total = 0
    all_records: list[UsnRecord] = []

    while cursor < boundary:
        if calls >= caps.max_calls:
            result.outcome = I_BUDGET_CALLS
            result.detail = f"call cap {caps.max_calls} reached at cursor {cursor}"
            break
        if requested_total + caps.buffer_bytes > caps.max_bytes:
            result.outcome = I_BUDGET_BYTES
            result.detail = f"byte cap {caps.max_bytes} reached at cursor {cursor}"
            break
        calls += 1
        requested_total += caps.buffer_bytes
        try:
            buffer = read_fn(cursor, caps.buffer_bytes)
        except Exception as exc:  # noqa: BLE001 - the reader must not raise past its contract
            result.outcome = _exception_outcome(exc)
            result.detail = _detail_for(result.outcome, f"read failed at cursor {cursor}", exc)
            break
        if not isinstance(buffer, (bytes, bytearray)):
            result.outcome = I_BAD_INPUT
            result.detail = "read function returned a non-bytes object"
            break
        returned_total += len(buffer)
        next_cursor, records, problems = parse_records(bytes(buffer), cursor)
        if problems:
            problem = problems[0]
            result.outcome = (
                problem.code
                if problem.code in INCONCLUSIVE_OUTCOMES
                else I_MALFORMED_RECORD
            )
            result.detail = problem.detail
            break
        if next_cursor < cursor:
            result.outcome = I_CURSOR_REGRESSED
            result.detail = f"cursor moved backwards: {cursor} -> {next_cursor}"
            break
        if next_cursor > boundary:
            # The API has no end-boundary parameter: the final call may return
            # records appended after the captured NextUsn. They are excluded.
            kept = [r for r in records if r.usn < boundary]
            result.over_read_records += len(records) - len(kept)
            records = kept
        if next_cursor == cursor:
            if records:
                result.outcome = I_NO_PROGRESS_WITH_RECORDS
                result.detail = f"records returned without traversal progress at cursor {cursor}"
            else:
                # Documented behaviour with BytesToWaitFor = 0: no matching
                # record at or after this position. The cursor did not advance,
                # so Claim R' is withheld: this says nothing about the range.
                result.outcome = OUTCOME_NO_MATCH_AT_CURSOR
                result.detail = f"no records matching the mask at or after USN {cursor}"
            break
        if not records:
            result.outcome = I_ADVANCED_WITHOUT_RECORDS
            result.detail = f"cursor advanced {cursor} -> {next_cursor} with no records"
            break
        for record in records:
            if record.usn >= next_cursor:
                result.outcome = I_MALFORMED_RECORD
                result.detail = f"record USN {record.usn} is not below the returned cursor {next_cursor}"
                break
        if result.outcome == I_MALFORMED_RECORD:
            break
        all_records.extend(records)
        cursor = next_cursor

    for record in all_records:
        problem = validate_record(record)
        if problem is not None:
            result.outcome = problem.code if problem.code in INCONCLUSIVE_OUTCOMES else I_MALFORMED_RECORD
            result.detail = problem.detail
            break

    after = _post_query(query_fn, result)
    if after is not None:
        retention = _check_retention(before, after)
        if retention is not None:
            result.outcome, result.detail = retention

    traversal_complete = cursor >= boundary
    valid = result.outcome in (OUTCOME_OK, OUTCOME_NOT_ATTEMPTED)
    result.claim_r_prime = bool(traversal_complete and valid and after is not None)
    if result.outcome == OUTCOME_NOT_ATTEMPTED:
        result.outcome = OUTCOME_OK

    result.records = all_records
    result.traversal = {
        "start_cursor": before.first_usn,
        "boundary": boundary,
        "end_cursor": cursor,
        "reached_boundary": traversal_complete,
        "calls": calls,
        "bytes_requested": requested_total,
        "bytes_returned": returned_total,
    }
    result.validity = {
        "records_returned": len(all_records),
        "by_version": {v: sum(1 for r in all_records if r.version == v) for v in SUPPORTED_VERSIONS},
        "reason_bits_seen": sorted({bit for r in all_records for bit in _bits(r.reason)}),
        "mask": f"0x{COMPLETE_REASON_MASK:08X}",
        "validation": "all returned records well formed and mask-consistent"
        if result.claim_r_prime
        else "not established",
    }
    result.gaps = _gap_report(all_records)

    if identity is not None:
        result.matches = _collect_matches(all_records, identity, caps, result)
        result.v4_sequences = _collect_sequences(all_records)
    return result


def _bits(value: int) -> Iterable[int]:
    for shift in range(32):
        bit = 1 << shift
        if value & bit:
            yield bit


def _safe_query(query_fn: QueryFn, result: ScanResult) -> JournalState | None:
    try:
        state = query_fn()
    except Exception as exc:  # noqa: BLE001
        result.outcome = _exception_outcome(exc)
        result.detail = _detail_for(result.outcome, "journal query failed", exc)
        return None
    if not isinstance(state, JournalState):
        result.outcome = I_BAD_INPUT
        result.detail = "journal query did not return a JournalState"
        return None
    return state


def _post_query(query_fn: QueryFn, result: ScanResult) -> JournalState | None:
    try:
        state = query_fn()
    except Exception as exc:  # noqa: BLE001
        # A failed closing query cannot establish retention stability.
        result.outcome = _exception_outcome(exc)
        result.detail = _detail_for(result.outcome, "post-scan journal query failed", exc)
        return None
    if not isinstance(state, JournalState):
        result.outcome = I_BAD_INPUT
        result.detail = "post-scan journal query did not return a JournalState"
        return None
    result.journal_state_after = state
    return state


def _winerror(exc: Exception) -> int | None:
    """The Win32 status carried by an exception, if any.

    `ctypes.WinError`/`OSError` set `winerror`; a plain `OSError(code, msg)` sets
    `errno`. Both are consulted, and neither is guessed from prose alone when a
    numeric status is present.
    """
    for attribute in ("winerror", "errno"):
        value = getattr(exc, attribute, None)
        if isinstance(value, int) and value:
            return value
    return None


def _exception_outcome(exc: Exception) -> str:
    """Map a failure to one named capability outcome, never to a negative result.

    The three journal statuses stay distinct because they say different things:

    * 1178 says the journal **is being deleted**. It does not say what journal,
      if any, will exist afterwards, and nothing here infers one.
    * 1179 says the journal **is not active**. It does not establish that no
      journal was ever created, nor anything about its history.
    * 1181 says the **requested entry has been deleted** from the journal. It is
      a statement about the requested range, not about whether a journal exists
      or is active.

    None of the three licenses a conclusion about the file under investigation.

    Numeric statuses are preferred; text matching exists only for providers that
    surface a message without a code, and it is deliberately narrow.
    """
    code = _winerror(exc)
    if code is not None:
        if code == WINERROR_JOURNAL_ENTRY_DELETED:
            return I_JOURNAL_ENTRY_DELETED
        if code == WINERROR_JOURNAL_DELETE_IN_PROGRESS:
            return I_JOURNAL_DELETE_IN_PROGRESS
        if code == WINERROR_JOURNAL_NOT_ACTIVE:
            return I_JOURNAL_INACTIVE

    text = str(exc).lower()
    if "entry deleted" in text or "journal entry has been deleted" in text:
        return I_JOURNAL_ENTRY_DELETED
    if "delete in progress" in text or "being deleted" in text:
        return I_JOURNAL_DELETE_IN_PROGRESS
    if "not active" in text or "journal is not active" in text:
        return I_JOURNAL_INACTIVE
    return I_BAD_INPUT


def _detail_for(outcome: str, context: str, exc: Exception) -> str:
    """Diagnostic text that names the condition rather than only the class."""
    description = OUTCOME_DESCRIPTIONS.get(outcome)
    if description is None:
        return f"{context}: {exc}"
    return f"{description} [{context}: {exc}]"


def _check_retention(before: JournalState, after: JournalState) -> tuple[str, str] | None:
    """Retention is checked across the scan, not only at its start (v6.9 §4)."""
    if after.journal_id != before.journal_id:
        return I_JOURNAL_ID_CHANGED, "journal id changed during the scan"
    if after.first_usn > before.first_usn:
        return (
            I_RETENTION_ADVANCED,
            f"retention advanced past the original start: {before.first_usn} -> {after.first_usn}",
        )
    if after.next_usn < before.next_usn:
        return I_NEXT_USN_REGRESSED, f"next usn moved backwards: {before.next_usn} -> {after.next_usn}"
    return None


def _gap_report(records: Sequence[UsnRecord]) -> dict:
    """Non-adjacent returned USNs — informational only.

    A gap is not truncation and not omission: a record whose only reason bits
    lie outside the mask is filtered out by design, and both explanations stay
    open.
    """
    gaps = []
    for previous, current in zip(records, records[1:]):
        end = previous.usn + previous.record_length
        if current.usn != end:
            gaps.append({"after_usn": previous.usn, "next_usn": current.usn})
    return {
        "count": len(gaps),
        "locations": gaps[:64],
        "interpretation": (
            "reported, not interpreted: a non-adjacent pair may reflect records filtered "
            "out by the reason mask, or discarded history; neither is concluded here"
        ),
    }


def _collect_matches(
    records: Sequence[UsnRecord], identity: FileIdentity, caps: Caps, result: ScanResult
) -> list[dict]:
    matches: list[dict] = []
    summarised = 0
    for record in records:
        if record.frn64() != identity.frn:
            continue
        entry = {
            "usn": record.usn,
            "version": record.version,
            "frn": f"0x{record.frn64():016X}",
            "reason": f"0x{record.reason:08X}",
            "reason_names": list(record.reason_names),
            "class_a_data_reason": record.has_class_a_data_reason,
            "timestamp": record.timestamp,
            "dated": record.dated,
            "undated_label": None if record.dated else "V4 (range-tracking, undated)",
            "source_info": decode_source_info(record.source_info),
            "volume_serial_64": identity.volume_serial_64,
            "journal_id": result.journal_state_before.journal_id if result.journal_state_before else None,
        }
        if summarised or _matches_bytes(matches, entry) > caps.max_output_bytes:
            summarised += 1
            continue
        matches.append(entry)
    if summarised:
        result.output_summarised = True
        matches.append(
            {
                "summarised": True,
                "omitted_records": summarised,
                "note": "output cap reached; counts retained, per-record detail omitted",
            }
        )
    return matches


def _matches_bytes(matches: Sequence[dict], entry: dict) -> int:
    # Cheap deterministic size proxy: avoids json import at module scope for a
    # quantity the cap only needs to bound loosely.
    def size(value) -> int:
        if isinstance(value, dict):
            return sum(size(k) + size(v) for k, v in value.items())
        if isinstance(value, (list, tuple)):
            return sum(size(v) for v in value)
        return len(str(value))

    return sum(size(m) for m in matches) + size(entry)


def _collect_sequences(records: Sequence[UsnRecord]) -> list[dict]:
    sequences: list[dict] = []
    index = 0
    while index < len(records):
        if records[index].version == 4:
            sequence, index = collect_v4_sequence(records, index)
            sequences.append(sequence.to_dict())
            continue
        index += 1
    return sequences


# --------------------------------------------------------------------------
# Real volume access — never reached by tests, never run at import
# --------------------------------------------------------------------------


def volume_io(volume_path: str, *, allow_open: bool = False):
    """Return (query_fn, read_fn) bound to a raw volume handle.

    This is the only code path that can open a volume. It requires a raw-volume
    path (`\\\\.\\X:`), Windows, and an explicit `allow_open=True`; the three
    live-volume operations are separate authorisation gates and are not
    performed by this module.
    """
    import os
    import sys

    if os.name != "nt" or sys.platform != "win32":
        raise RuntimeError("USN journal access requires Windows")
    if not allow_open:
        raise PermissionError(
            "volume access is not authorised by this module; the live operations are "
            "separate gates (G: validation, event-log enumeration, S1-B query)"
        )
    if not volume_path.startswith("\\\\.\\") or not volume_path.endswith(":"):
        raise ValueError("a raw volume path such as \\\\.\\G: is required")

    from ctypes import wintypes

    GENERIC_READ = 0x80000000
    FILE_SHARE_READ = 0x00000001
    FILE_SHARE_WRITE = 0x00000002
    OPEN_EXISTING = 3

    kernel32 = ctypes.WinDLL("kernel32", use_last_error=True)

    # Signatures must be declared. Without them ctypes guesses `int` for the
    # return, which truncates a 64-bit HANDLE, and passes the handle back as a
    # 32-bit integer — so a valid handle would be silently corrupted and
    # INVALID_HANDLE_VALUE (-1) would not compare equal to c_void_p(-1).
    kernel32.CreateFileW.restype = wintypes.HANDLE
    kernel32.CreateFileW.argtypes = [
        wintypes.LPCWSTR,
        wintypes.DWORD,
        wintypes.DWORD,
        wintypes.LPVOID,
        wintypes.DWORD,
        wintypes.DWORD,
        wintypes.HANDLE,
    ]
    kernel32.DeviceIoControl.restype = wintypes.BOOL
    kernel32.DeviceIoControl.argtypes = [
        wintypes.HANDLE,
        wintypes.DWORD,
        wintypes.LPVOID,
        wintypes.DWORD,
        wintypes.LPVOID,
        wintypes.DWORD,
        ctypes.POINTER(wintypes.DWORD),
        wintypes.LPVOID,
    ]

    invalid_handle = ctypes.c_void_p(-1).value
    handle = kernel32.CreateFileW(
        volume_path,
        GENERIC_READ,
        FILE_SHARE_READ | FILE_SHARE_WRITE,
        None,
        OPEN_EXISTING,
        0,
        None,
    )
    if handle is None or handle in (0, invalid_handle):
        raise OSError(ctypes.get_last_error(), f"CreateFileW failed for {volume_path}")

    def query() -> JournalState:
        data = USN_JOURNAL_DATA_V0()
        returned = wintypes.DWORD(0)
        ok = kernel32.DeviceIoControl(
            handle,
            FSCTL_QUERY_USN_JOURNAL,
            None,
            0,
            ctypes.byref(data),
            ctypes.sizeof(data),
            ctypes.byref(returned),
            None,
        )
        if not ok:
            raise OSError(ctypes.get_last_error(), "FSCTL_QUERY_USN_JOURNAL failed")
        return JournalState(
            journal_id=data.UsnJournalID,
            first_usn=data.FirstUsn,
            next_usn=data.NextUsn,
            lowest_valid_usn=data.LowestValidUsn,
            maximum_size=data.MaximumSize,
        )

    # The journal id is captured in a mutable holder rather than a rebinding
    # local, so `read` cannot observe a stale value after a failed query.
    journal = {"id": 0}

    def read(start_usn: int, buffer_bytes: int) -> bytes:
        request = READ_USN_JOURNAL_DATA_V1(
            StartUsn=start_usn,
            ReasonMask=COMPLETE_REASON_MASK,
            ReturnOnlyOnClose=0,
            Timeout=0,
            BytesToWaitFor=0,
            UsnJournalID=journal["id"],
            MinMajorVersion=MIN_MAJOR_VERSION,
            MaxMajorVersion=MAX_MAJOR_VERSION,
        )
        buffer = ctypes.create_string_buffer(buffer_bytes)
        returned = wintypes.DWORD(0)
        ok = kernel32.DeviceIoControl(
            handle,
            FSCTL_READ_USN_JOURNAL,
            ctypes.byref(request),
            ctypes.sizeof(request),
            ctypes.cast(buffer, ctypes.c_void_p),
            buffer_bytes,
            ctypes.byref(returned),
            None,
        )
        if not ok:
            raise OSError(ctypes.get_last_error(), "FSCTL_READ_USN_JOURNAL failed")
        return buffer.raw[: returned.value]

    try:
        journal["id"] = query().journal_id
    except OSError:
        journal["id"] = 0
    return query, read
