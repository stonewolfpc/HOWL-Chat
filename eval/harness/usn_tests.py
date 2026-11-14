"""Offline synthetic tests for `eval.harness.usn_reader` (v6.9 specification).

Everything here is a synthetic buffer, a fake response or a controlled fixture.
No test opens a volume, queries a journal, enumerates an event log, or touches a
model file. The only volume-facing function, `volume_io`, is asserted to refuse
unless explicitly authorised, and is never called.

Usage:  python -m eval.harness.usn_tests
Exit:   0 = every check passed
        1 = at least one check failed
"""

from __future__ import annotations

import struct
import sys

from . import usn_reader as R

CHECKS: list[tuple[str, callable]] = []
RESULTS: list[tuple[str, bool, str]] = []


def check(name: str):
    def register(fn):
        CHECKS.append((name, fn))
        return fn

    return register


# --------------------------------------------------------------------------
# Fixtures
# --------------------------------------------------------------------------

S1B_FRN = 0x0003000000014F91
OTHER_FRN = 0x0000000000000123
VOLUME_SERIAL_64 = 0x880ECBB30ECB9918


def _file_id(frn: int, high: int = 0) -> bytes:
    return frn.to_bytes(8, "little") + high.to_bytes(8, "little")


def v2(frn=S1B_FRN, usn=1000, reason=R.REASON_FLAGS["DATA_OVERWRITE"], ts=133_000_000_000_000_000,
       source=0, name=b"file.gguf", minor=0, length=None, version=2):
    header = 60
    total = header + len(name) if length is None else length
    return struct.pack(
        "<IHHQQqQIIIIHH",
        total,
        version,
        minor,
        frn,
        frn,
        usn,
        ts,
        reason,
        source,
        0,
        0x20,
        len(name),
        header,
    ) + name


def v3(frn=S1B_FRN, usn=2000, reason=R.REASON_FLAGS["CLOSE"], ts=133_000_000_000_000_001,
       source=0, name=b"file.gguf", high=0, length=None):
    header = 76
    total = header + len(name) if length is None else length
    return struct.pack(
        "<IHHQQQQQQIIIIHH",
        total,
        3,
        0,
        frn,
        high,
        frn,
        0,
        usn,
        ts,
        reason,
        source,
        0,
        0x20,
        len(name),
        header,
    ) + name


def v4(frn=S1B_FRN, usn=3000, reason=R.REASON_FLAGS["DATA_OVERWRITE"], source=0,
       remaining=0, extents=((0, 4096),)):
    count = len(extents)
    total = 64 + count * 16
    blob = struct.pack(
        "<IHHQQQQQIIIHH",
        total,
        4,
        0,
        frn,
        frn,
        frn,
        frn,
        usn,
        reason,
        source,
        remaining,
        count,
        16,
    )
    for offset, length in extents:
        blob += struct.pack("<qq", offset, length)
    return blob


def buf(cursor: int, *records: bytes) -> bytes:
    return struct.pack("<Q", cursor) + b"".join(records)


def parsed(*blobs: bytes):
    """Parse record blobs into UsnRecord objects for sequence-level checks."""
    cursor = 10**9
    _, records, problems = R.parse_records(buf(cursor, *blobs), 0)
    assert not problems, problems
    return records


def state(journal_id=7, first=1000, next=5000, lowest=900, maximum=32 * 1024 * 1024) -> R.JournalState:
    return R.JournalState(journal_id=journal_id, first_usn=first, next_usn=next,
                          lowest_valid_usn=lowest, maximum_size=maximum)


class Vol:
    """A scripted volume: query states in order, read buffers in order."""

    def __init__(self, states, buffers):
        self.states = list(states)
        self.buffers = list(buffers)
        self.reads: list[tuple[int, int]] = []
        self.queries = 0

    def query(self):
        self.queries += 1
        if not self.states:
            return state()
        return self.states.pop(0) if len(self.states) > 1 else self.states[0]

    def read(self, start_usn, buffer_bytes):
        self.reads.append((start_usn, buffer_bytes))
        if not self.buffers:
            raise AssertionError("test volume ran out of scripted buffers")
        return self.buffers.pop(0)


IDENTITY = R.FileIdentity(frn=S1B_FRN, volume_serial_64=VOLUME_SERIAL_64,
                          file_id_128=_file_id(S1B_FRN))


def run(states, buffers, header_records, identity=IDENTITY, caps=None):
    vol = Vol(states, buffers)
    return R.scan(identity, vol.query, vol.read, caps or R.Caps(max_calls=8, max_bytes=8 * 1024 * 1024)), vol


# --------------------------------------------------------------------------
# 1. Parsing: headers, lengths, versions
# --------------------------------------------------------------------------

@check("parse: a valid V2 record is read with FRN, reason and timestamp")
def _p1():
    cursor, records, problems = R.parse_records(buf(2000, v2(usn=1000)), 1000)
    assert not problems, problems
    assert cursor == 2000
    rec = records[0]
    assert rec.version == 2 and rec.frn64() == S1B_FRN
    assert rec.reason == R.REASON_FLAGS["DATA_OVERWRITE"]
    assert rec.dated and rec.timestamp == 133_000_000_000_000_000


@check("parse: a truncated header is malformed, not carried forward")
def _p2():
    _, records, problems = R.parse_records(buf(1500) + b"\x01\x02\x03", 1000)
    assert not records
    assert problems and problems[0].code == R.I_MALFORMED_RECORD


@check("parse: a record claiming more bytes than the buffer holds is malformed")
def _p3():
    blob = v2(usn=1000, length=4000)
    _, _, problems = R.parse_records(buf(1400) + blob, 1000)
    assert problems and "exceeds" in problems[0].detail


@check("parse: a zero or under-header record length is malformed")
def _p4():
    _, _, problems = R.parse_records(buf(1000) + struct.pack("<IHH", 0, 2, 0), 1000)
    assert problems and problems[0].code == R.I_MALFORMED_RECORD


@check("parse: an unsupported record version is rejected, not guessed")
def _p5():
    _, _, problems = R.parse_records(buf(1000) + v2(usn=1000, version=9), 1000)
    assert problems and problems[0].code == R.I_UNSUPPORTED_VERSION


@check("journal state: the ioctl fields map onto the captured four values")
def _p8():
    st = R.parse_journal_state(
        {"usn_journal_id": 7, "first_usn": 100, "next_usn": 500, "lowest_valid_usn": 90,
         "maximum_size": 1024}
    )
    assert (st.journal_id, st.first_usn, st.next_usn, st.lowest_valid_usn) == (7, 100, 500, 90)
    assert st.retained_bytes == 400 and st.is_empty_range is False


@check("journal state: a positional sequence is accepted")
def _p9():
    st = R.parse_journal_state((7, 100, 500, 90, 1024))
    assert st.next_usn == 500 and st.maximum_size == 1024


@check("journal state: unrecognised input is rejected, not guessed")
def _p10():
    for bad in ({}, {"usn_journal_id": 1}, (1, 2), "nonsense"):
        try:
            R.parse_journal_state(bad)
        except ValueError:
            continue
        raise AssertionError(f"malformed journal state accepted: {bad!r}")


@check("parse: a buffer shorter than the cursor is malformed")
def _p6():
    _, _, problems = R.parse_records(b"\x00\x01", 1000)
    assert problems and problems[0].code == R.I_MALFORMED_RECORD


@check("parse: a record before the call's start USN is malformed")
def _p7():
    _, _, problems = R.parse_records(buf(1200, v2(usn=999)), 1000)
    assert problems and "precedes" in problems[0].detail


# --------------------------------------------------------------------------
# 2. Validation: mask and reason bits
# --------------------------------------------------------------------------

@check("validate: a documented reason inside the mask is accepted")
def _v1():
    _, records, _ = R.parse_records(buf(1100, v2(usn=1000)), 1000)
    assert R.validate_record(records[0]) is None


@check("api: the documented structures have the documented layout")
def _api1():
    import ctypes

    assert ctypes.sizeof(R.USN_JOURNAL_DATA_V0) == 56
    journal_offsets = {
        "UsnJournalID": 0, "FirstUsn": 8, "NextUsn": 16, "LowestValidUsn": 24,
        "MaxUsn": 32, "MaximumSize": 40, "AllocationDelta": 48,
    }
    for name, offset in journal_offsets.items():
        assert getattr(R.USN_JOURNAL_DATA_V0, name).offset == offset, name

    # 44 bytes of fields, padded to 48 by 8-byte alignment on x64.
    assert ctypes.sizeof(R.READ_USN_JOURNAL_DATA_V1) == 48
    read_offsets = {
        "StartUsn": 0, "ReasonMask": 8, "ReturnOnlyOnClose": 12, "Timeout": 16,
        "BytesToWaitFor": 24, "UsnJournalID": 32, "MinMajorVersion": 40,
        "MaxMajorVersion": 42,
    }
    for name, offset in read_offsets.items():
        assert getattr(R.READ_USN_JOURNAL_DATA_V1, name).offset == offset, name


@check("api: both control codes recompute from the CTL_CODE macro")
def _api2():
    assert R.FSCTL_QUERY_USN_JOURNAL == 0x000900F4
    assert R.FSCTL_READ_USN_JOURNAL == 0x000900BB
    assert R.ctl_code(61, R.METHOD_BUFFERED) == R.FSCTL_QUERY_USN_JOURNAL
    assert R.ctl_code(46, R.METHOD_NEITHER) == R.FSCTL_READ_USN_JOURNAL
    assert R.ctl_code(46, R.METHOD_BUFFERED) != R.FSCTL_READ_USN_JOURNAL


@check("text: capability descriptions say only what the status says")
def _txt1():
    d1178 = R.OUTCOME_DESCRIPTIONS[R.I_JOURNAL_DELETE_IN_PROGRESS]
    d1179 = R.OUTCOME_DESCRIPTIONS[R.I_JOURNAL_INACTIVE]
    d1181 = R.OUTCOME_DESCRIPTIONS[R.I_JOURNAL_ENTRY_DELETED]
    assert "being deleted" in d1178 and "cannot proceed" in d1178
    assert "not active" in d1179 and "cannot proceed" in d1179
    assert "requested journal entry has been deleted from the journal" in d1181
    # No inference beyond the documented status.
    for text in (d1178, d1179, d1181):
        lowered = text.lower()
        for forbidden in ("replacement", "no journal exists", "guarantee", "discarded"):
            assert forbidden not in lowered, (forbidden, text)


@check("validate: a record with no matching reason bit invalidates the premise")
def _v2():
    _, records, _ = R.parse_records(
        buf(1100, v2(usn=1000, reason=R.REASON_FLAGS["DATA_OVERWRITE"])), 1000
    )
    problem = R.validate_record(records[0], mask=0)
    assert problem and problem.code == R.I_MASK_MISMATCH


@check("validate: undocumented reason bits are rejected")
def _v3():
    _, records, _ = R.parse_records(
        buf(1100, v2(usn=1000, reason=R.REASON_FLAGS["DATA_OVERWRITE"] | 0x01000000)), 1000
    )
    problem = R.validate_record(records[0])
    assert problem and problem.code == R.I_UNDOCUMENTED_REASON


@check("validate: reserved low bits are rejected")
def _v4():
    _, records, _ = R.parse_records(
        buf(1100, v2(usn=1000, reason=R.REASON_FLAGS["DATA_OVERWRITE"] | 0x08)), 1000
    )
    problem = R.validate_record(records[0])
    assert problem and problem.code == R.I_UNDOCUMENTED_REASON


@check("validate: a FILE_ID_128 high half that is not zero is malformed on NTFS")
def _v5():
    _, records, _ = R.parse_records(buf(2100, v3(usn=2000, high=1)), 2000)
    problem = R.validate_record(records[0])
    assert problem and problem.code == R.I_MALFORMED_RECORD


@check("mask: the constant is the OR of every documented flag and no reserved bit")
def _v6():
    expected = 0
    for bit in R.REASON_FLAGS.values():
        expected |= bit
    assert expected == R.COMPLETE_REASON_MASK, hex(expected)
    assert R.COMPLETE_REASON_MASK & R.RESERVED_REASON_BITS == 0
    for bit in (0x08, 0x80, 0x01000000, 0x40000000):
        assert not R.COMPLETE_REASON_MASK & bit


# --------------------------------------------------------------------------
# 3. V2/V3/V4 identity
# --------------------------------------------------------------------------

@check("identity: V2 compares its 8-byte FRN directly")
def _i1():
    _, records, _ = R.parse_records(buf(1100, v2(usn=1000)), 1000)
    assert records[0].frn64() == S1B_FRN and records[0].high_half() == 0


@check("identity: V3 compares the low 64 bits of FILE_ID_128")
def _i2():
    _, records, _ = R.parse_records(buf(2100, v3(usn=2000)), 2000)
    assert records[0].frn64() == S1B_FRN and records[0].high_half() == 0


@check("identity: a different file's records are not attributed to the target")
def _i3():
    result, _ = run([state(), state()], [buf(5000, v2(frn=OTHER_FRN, usn=1000))], None)
    assert result.matches == [], result.matches
    assert result.outcome == R.OUTCOME_OK


@check("volume: 64-bit serial equality is the identity check")
def _i4():
    out = R.verify_volume_identity(VOLUME_SERIAL_64, VOLUME_SERIAL_64)
    assert out["identity_established"] is True
    assert out["serial_32_is_not_identity_proof"] is True


@check("volume: a 64-bit serial mismatch is not rescued by 32-bit agreement")
def _i5():
    out = R.verify_volume_identity(VOLUME_SERIAL_64, 0x1111111111111111, 0x0ECB9918, 0x0ECB9918)
    assert out["identity_established"] is False
    assert out["serial_32_consistency_check"] is True
    assert "not established" in out["statement"]


@check("volume: the 32-bit serials are only a consistency check")
def _i6():
    out = R.verify_volume_identity(VOLUME_SERIAL_64, VOLUME_SERIAL_64, 1, 2)
    assert out["identity_established"] is True
    assert out["serial_32_consistency_check"] is False


@check("identity: from_file_id derives the documented FRN")
def _i7():
    ident = R.FileIdentity.from_file_id(bytes.fromhex("914f0100000003000000000000000000"), VOLUME_SERIAL_64)
    assert ident.frn == S1B_FRN
    try:
        R.FileIdentity.from_file_id(b"\x00" * 8, 1)
    except ValueError:
        return
    raise AssertionError("an 8-byte file id was accepted as FILE_ID_128")


# --------------------------------------------------------------------------
# 4. Cursor traversal
# --------------------------------------------------------------------------

@check("scan: a completed traversal yields Claim R' and one matching record")
def _c1():
    result, vol = run(
        [state(), state()],
        [buf(5000, v2(usn=1000))],
        None,
    )
    assert result.outcome == R.OUTCOME_OK, result.detail
    assert result.claim_r_prime is True
    assert result.traversal["calls"] == 1 and vol.queries == 2
    assert len(result.matches) == 1
    assert result.matches[0]["reason_names"] == ["DATA_OVERWRITE"]


@check("scan: multiple calls advance the cursor to the fixed boundary")
def _c2():
    result, vol = run(
        [state(first=1000, next=3000), state(first=1000, next=3000)],
        [buf(2000, v2(usn=1000)), buf(3000, v2(usn=2000))],
        None,
    )
    assert result.outcome == R.OUTCOME_OK
    assert result.traversal["calls"] == 2
    assert [s for s, _ in vol.reads] == [1000, 2000]
    assert result.validity["records_returned"] == 2


@check("scan: an unchanged cursor with no records is a no-match result, not a completed traversal")
def _c3():
    result, _ = run([state(first=1000, next=5000), state()], [buf(1000)], None)
    assert result.outcome == R.OUTCOME_NO_MATCH_AT_CURSOR, result.outcome
    assert result.claim_r_prime is False, "no-match must not be reported as a completed traversal"
    assert result.validity["records_returned"] == 0
    assert "no records matching" in result.detail.lower()


@check("scan: an unchanged cursor with records is INCONCLUSIVE")
def _c4():
    result, _ = run([state(), state()], [buf(1000, v2(usn=1000))], None)
    assert result.outcome == R.I_NO_PROGRESS_WITH_RECORDS, result.outcome
    assert result.claim_r_prime is False


@check("scan: an advancing cursor with no records is INCONCLUSIVE")
def _c5():
    result, _ = run([state(), state()], [buf(2000)], None)
    assert result.outcome == R.I_ADVANCED_WITHOUT_RECORDS, result.outcome


@check("scan: a backwards cursor is INCONCLUSIVE")
def _c6():
    result, _ = run([state(first=2000, next=5000), state()], [buf(1000)], None)
    assert result.outcome == R.I_CURSOR_REGRESSED, result.outcome


@check("scan: records past the fixed boundary are excluded, not analysed")
def _c7():
    result, _ = run(
        [state(first=1000, next=2000), state(first=1000, next=2000)],
        [buf(3500, v2(usn=1000), v2(usn=2500))],
        None,
    )
    assert result.validity["records_returned"] == 1, result.validity
    assert result.over_read_records == 1
    assert result.claim_r_prime is True


@check("scan: an empty retained range completes without inventing records")
def _c8():
    result, vol = run([state(first=5000, next=5000), state(first=5000, next=5000)], [], None)
    assert result.outcome == R.OUTCOME_OK
    assert vol.reads == []
    assert result.validity["records_returned"] == 0


# --------------------------------------------------------------------------
# 5. Retention and journal identity across the scan
# --------------------------------------------------------------------------

@check("retention: a changed journal id is INCONCLUSIVE")
def _r1():
    result, _ = run([state(), state(journal_id=8)], [buf(5000, v2(usn=1000))], None)
    assert result.outcome == R.I_JOURNAL_ID_CHANGED, result.outcome
    assert result.claim_r_prime is False


@check("retention: FirstUsn advancing past the start is INCONCLUSIVE")
def _r2():
    result, _ = run([state(first=1000), state(first=1500)], [buf(5000, v2(usn=1000))], None)
    assert result.outcome == R.I_RETENTION_ADVANCED, result.outcome


@check("retention: NextUsn moving backwards is INCONCLUSIVE")
def _r3():
    result, _ = run([state(next=5000), state(next=4000)], [buf(4000, v2(usn=1000))], None)
    assert result.outcome == R.I_NEXT_USN_REGRESSED, result.outcome


@check("retention: a completed cursor does not override a moving boundary")
def _r4():
    result, _ = run([state(), state(first=1200)], [buf(5000, v2(usn=1000))], None)
    assert result.traversal["reached_boundary"] is True
    assert result.outcome == R.I_RETENTION_ADVANCED
    assert result.claim_r_prime is False


@check("retention: growth of NextUsn alone is tolerated")
def _r5():
    result, _ = run([state(next=5000), state(next=9000)], [buf(5000, v2(usn=1000))], None)
    assert result.outcome == R.OUTCOME_OK


@check("retention: a failed closing query cannot establish retention stability")
def _r6():
    class Failing(Vol):
        def query(self):
            self.queries += 1
            if self.queries == 1:
                return state()
            raise OSError(1178, "journal entry deleted")

    vol = Failing([state()], [buf(5000, v2(usn=1000))])
    result = R.scan(IDENTITY, vol.query, vol.read)
    assert result.claim_r_prime is False
    assert result.outcome.startswith("INCONCLUSIVE"), result.outcome


# --------------------------------------------------------------------------
# 6. Aborts and budget
# --------------------------------------------------------------------------

def _failing_read(exc):
    """A volume whose first read raises the supplied exception."""

    class Failing(Vol):
        def read(self, start_usn, buffer_bytes):
            raise exc

    vol = Failing([state()], [])
    return R.scan(IDENTITY, vol.query, vol.read)


@check("abort: 1181 entry-deleted is its own outcome, with diagnostic text")
def _a1():
    result = _failing_read(OSError(1181, "journal entry deleted"))
    assert result.outcome == R.I_JOURNAL_ENTRY_DELETED, result.outcome
    assert "1181" in result.detail and "ERROR_JOURNAL_ENTRY_DELETED" in result.detail, result.detail
    assert "has been deleted from the journal" in result.detail
    assert result.claim_r_prime is False


@check("abort: 1178 delete-in-progress is NOT classified as inactive")
def _a2():
    result = _failing_read(OSError(1178, "volume change journal is being deleted"))
    assert result.outcome == R.I_JOURNAL_DELETE_IN_PROGRESS, result.outcome
    assert result.outcome != R.I_JOURNAL_INACTIVE, "1178 was conflated with 1179"
    assert "1178" in result.detail and "ERROR_JOURNAL_DELETE_IN_PROGRESS" in result.detail
    assert "being deleted" in result.detail


@check("abort: 1179 not-active is its own outcome, with diagnostic text")
def _a2b():
    result = _failing_read(OSError(1179, "volume change journal is not active"))
    assert result.outcome == R.I_JOURNAL_INACTIVE, result.outcome
    assert result.outcome != R.I_JOURNAL_DELETE_IN_PROGRESS
    assert "1179" in result.detail and "ERROR_JOURNAL_NOT_ACTIVE" in result.detail
    assert "cannot proceed" in result.detail


@check("abort: the three journal statuses are three distinct named outcomes")
def _a2c():
    names = {
        R.I_JOURNAL_DELETE_IN_PROGRESS,
        R.I_JOURNAL_INACTIVE,
        R.I_JOURNAL_ENTRY_DELETED,
    }
    assert len(names) == 3, names
    for name in names:
        assert name in R.INCONCLUSIVE_OUTCOMES, name
        assert name in R.OUTCOME_DESCRIPTIONS, name
    assert R.WINERROR_JOURNAL_DELETE_IN_PROGRESS == 1178
    assert R.WINERROR_JOURNAL_NOT_ACTIVE == 1179
    assert R.WINERROR_JOURNAL_ENTRY_DELETED == 1181


@check("abort: the numeric status wins over misleading message text")
def _a2d():
    # A Windows error whose formatted text is unhelpful must still map by code.
    result = _failing_read(OSError(1178, "unspecified error"))
    assert result.outcome == R.I_JOURNAL_DELETE_IN_PROGRESS, result.outcome
    result = _failing_read(OSError(1181, "unspecified error"))
    assert result.outcome == R.I_JOURNAL_ENTRY_DELETED, result.outcome


@check("abort: an exception carrying only winerror is mapped by that attribute")
def _a2e():
    class Win32(Exception):
        winerror = 1181

    result = _failing_read(Win32("operation failed"))
    assert result.outcome == R.I_JOURNAL_ENTRY_DELETED, result.outcome


@check("abort: message-only failures map narrowly, and unknown text stays generic")
def _a2f():
    result = _failing_read(RuntimeError("the volume change journal is being deleted"))
    assert result.outcome == R.I_JOURNAL_DELETE_IN_PROGRESS, result.outcome
    result = _failing_read(RuntimeError("the volume change journal is not active"))
    assert result.outcome == R.I_JOURNAL_INACTIVE, result.outcome
    result = _failing_read(RuntimeError("something else entirely"))
    assert result.outcome == R.I_BAD_INPUT, result.outcome
    assert "1178" not in result.detail and "1181" not in result.detail


@check("abort: a mask mismatch is reachable through the scan, not only at unit level")
def _a3b():
    result, _ = run(
        [state(), state()],
        [buf(5000, v2(usn=1000, reason=0x01000000))],
        None,
    )
    assert result.outcome == R.I_MASK_MISMATCH, result.outcome
    assert result.claim_r_prime is False


@check("abort: a malformed record aborts the scan")
def _a3():
    blob = v2(usn=1000, length=40)
    result, _ = run([state(), state()], [buf(5000, blob)], None)
    assert result.outcome == R.I_MALFORMED_RECORD, result.outcome
    assert result.claim_r_prime is False


@check("abort: an unsupported version aborts the scan")
def _a4():
    result, _ = run([state(), state()], [buf(5000, v2(usn=1000, version=7))], None)
    assert result.outcome == R.I_UNSUPPORTED_VERSION, result.outcome


@check("abort: an undocumented reason bit aborts the scan")
def _a5():
    result, _ = run(
        [state(), state()],
        [buf(5000, v2(usn=1000, reason=R.REASON_FLAGS["DATA_OVERWRITE"] | 0x02000000))],
        None,
    )
    assert result.outcome == R.I_UNDOCUMENTED_REASON, result.outcome


@check("budget: the call cap stops the scan and reports INCONCLUSIVE")
def _b1():
    caps = R.Caps(max_calls=3, max_bytes=64 * 1024 * 1024, buffer_bytes=1024 * 1024)
    result, vol = run(
        [state(first=1000, next=100000), state(first=1000, next=100000)],
        [buf(2000, v2(usn=1000)), buf(3000, v2(usn=2000)), buf(4000, v2(usn=3000))],
        None,
        caps=caps,
    )
    assert result.outcome == R.I_BUDGET_CALLS, result.outcome
    assert len(vol.reads) == 3


@check("budget: the byte cap stops the scan and reports INCONCLUSIVE")
def _b2():
    caps = R.Caps(max_calls=10, max_bytes=2 * 1024 * 1024 + 1, buffer_bytes=1024 * 1024)
    result, _ = run(
        [state(first=1000, next=100000), state(first=1000, next=100000)],
        [buf(2000, v2(usn=1000)), buf(3000, v2(usn=2000)), buf(4000, v2(usn=3000))],
        None,
        caps=caps,
    )
    assert result.outcome in (R.I_BUDGET_BYTES, R.I_BUDGET_CALLS), result.outcome


@check("budget: a span beyond the cap is refused before any read")
def _b3():
    big = state(first=0, next=1024 * 1024 * 1024)
    vol = Vol([big], [])
    result = R.scan(IDENTITY, vol.query, vol.read, R.Caps(max_calls=128, max_bytes=1024))
    assert result.outcome == R.I_BUDGET_PREFLIGHT, result.outcome
    assert vol.reads == [], "a failed preflight must read nothing"


@check("budget: preflight never promises completion")
def _b4():
    plan = R.preflight(state(first=0, next=10 * 1024 * 1024), R.Caps())
    assert plan.budget_appears_sufficient is True
    assert plan.completion_guaranteed is False
    assert "not guaranteed" in plan.to_dict()["statement"]


@check("budget: an insufficient estimate says so, in those words")
def _b5():
    plan = R.preflight(state(first=0, next=200 * 1024 * 1024), R.Caps())
    assert plan.budget_appears_sufficient is False
    assert plan.to_dict()["statement"] == "budget feasibility not proven"


@check("budget: caps reject an impossible buffer size")
def _c9():
    try:
        R.Caps(buffer_bytes=16)
    except ValueError:
        return
    raise AssertionError("a buffer smaller than header + cursor was accepted")


# --------------------------------------------------------------------------
# 7. Output bound
# --------------------------------------------------------------------------

@check("output: the cap summarises rather than growing without bound")
def _o1():
    records = [v2(usn=1000 + i * 100, name=b"x" * 40) for i in range(40)]
    cursor = 1000 + 40 * 100
    caps = R.Caps(max_calls=8, max_bytes=8 * 1024 * 1024, max_output_bytes=1200)
    result, _ = run(
        [state(first=1000, next=cursor), state(first=1000, next=cursor)],
        [buf(cursor, *records)],
        None,
        caps=caps,
    )
    assert result.output_summarised is True
    assert result.matches[-1].get("summarised") is True
    assert result.matches[-1]["omitted_records"] > 0


# --------------------------------------------------------------------------
# 8. V4 sequences
# --------------------------------------------------------------------------

@check("v4: a single-record sequence terminates and finds its V3 companion")
def _q1():
    recs = parsed(v4(usn=1000, remaining=0, extents=((0, 4096),)), v3(usn=1100, reason=R.REASON_FLAGS["CLOSE"]))
    seq, index = R.collect_v4_sequence(recs, 0)
    assert seq.outcome == R.SEQ_CONSISTENT, seq.detail
    assert seq.companion is not None and seq.companion.usn == 1100
    assert index == 1


@check("v4: a multi-record sequence is grouped and validated")
def _q2():
    recs = parsed(
        v4(usn=1000, remaining=2, extents=((0, 4096), (4096, 4096))),
        v4(usn=1100, remaining=1, extents=((8192, 4096),)),
        v4(usn=1200, remaining=0, extents=((12288, 4096),)),
        v3(usn=1300, reason=R.REASON_FLAGS["CLOSE"]),
    )
    seq, index = R.collect_v4_sequence(recs, 0)
    assert seq.outcome == R.SEQ_CONSISTENT, seq.detail
    assert len(seq.records) == 3 and seq.total_extents == 4
    assert seq.companion is not None
    assert index == 3


@check("v4: an underflowing recurrence is INVALID")
def _q3():
    recs = parsed(
        v4(usn=1000, remaining=1, extents=((0, 4096), (4096, 4096))),
        v4(usn=1100, remaining=0, extents=((8192, 4096),)),
    )
    seq, _ = R.collect_v4_sequence(recs, 0)
    assert seq.outcome == R.SEQ_INVALID, seq.outcome


@check("v4: a recurrence that skips a count is INVALID, not merely non-increasing")
def _q4():
    recs = parsed(
        v4(usn=1000, remaining=4, extents=((0, 4096),)),
        v4(usn=1100, remaining=1, extents=((4096, 4096),)),
        v4(usn=1200, remaining=0, extents=((8192, 4096),)),
    )
    seq, _ = R.collect_v4_sequence(recs, 0)
    assert seq.outcome == R.SEQ_INVALID, seq.detail


@check("v4: a non-zero terminal record inside the scan is INVALID/BOUNDARY-TRUNCATED")
def _q5():
    recs = parsed(v4(usn=1000, remaining=2, extents=((0, 4096),)))
    seq, _ = R.collect_v4_sequence(recs, 0)
    assert seq.outcome in (R.SEQ_BOUNDARY_TRUNCATED, R.SEQ_INVALID), seq.outcome
    assert seq.companion is None


@check("v4: an extent array disagreeing with RecordLength is INVALID")
def _q6():
    blob = v4(usn=1000, remaining=0, extents=((0, 4096), (4096, 4096)))
    # Corrupt NumberOfExtents upward without growing the record.
    offset = 60  # NumberOfExtents in a V4 header
    blob = blob[:offset] + struct.pack("<H", 3) + blob[offset + 2 :]
    _, records, problems = R.parse_records(buf(2000, blob), 1000)
    assert problems, "a record whose extent array overruns its length must be rejected"
    assert problems[0].code == R.I_MALFORMED_RECORD


@check("v4: a non-16-byte ExtentSize is INVALID")
def _q7():
    blob = v4(usn=1000, remaining=0, extents=((0, 4096),))
    offset = 62
    blob = blob[:offset] + struct.pack("<H", 8) + blob[offset + 2 :]
    _, records, problems = R.parse_records(buf(2000, blob), 1000)
    assert problems and problems[0].code == R.I_MALFORMED_RECORD


@check("v4: a following V4 for another file is not treated as the companion")
def _q8():
    recs = parsed(
        v4(usn=1000, remaining=0, extents=((0, 4096),)),
        v4(frn=OTHER_FRN, usn=1100, remaining=0, extents=((0, 4096),)),
        v3(usn=1200, reason=R.REASON_FLAGS["CLOSE"]),
    )
    seq, _ = R.collect_v4_sequence(recs, 0)
    assert seq.outcome == R.SEQ_CONSISTENT
    assert seq.companion is not None and seq.companion.usn == 1200


@check("v4: a companion that is not V3+CLOSE is reported as absent")
def _q9():
    recs = parsed(
        v4(usn=1000, remaining=0, extents=((0, 4096),)),
        v3(usn=1100, reason=R.REASON_FLAGS["DATA_OVERWRITE"]),
    )
    seq, _ = R.collect_v4_sequence(recs, 0)
    assert seq.outcome == R.SEQ_CONSISTENT
    assert seq.companion is None
    assert "no companion" in seq.detail.lower()


@check("v4: the companion timestamp is labelled as a close-record timestamp")
def _q10():
    recs = parsed(v4(usn=1000), v3(usn=1100, reason=R.REASON_FLAGS["CLOSE"]))
    seq, _ = R.collect_v4_sequence(recs, 0)
    payload = seq.to_dict()
    assert payload["undated"] is True
    assert payload["companion_usn"] == 1100
    assert "close-record" in payload["companion_timestamp_label"]


@check("v4: a V4 record is never given a timestamp")
def _q11():
    _, records, _ = R.parse_records(buf(2000, v4(usn=1000)), 1000)
    assert records[0].dated is False
    assert records[0].timestamp is None


# --------------------------------------------------------------------------
# 9. SourceInfo
# --------------------------------------------------------------------------

@check("source-info: DATA_MANAGEMENT is decoded and qualified")
def _s1():
    out = R.decode_source_info(R.SOURCE_INFO_FLAGS["DATA_MANAGEMENT"])
    assert "DATA_MANAGEMENT" in out["flags"]
    assert out["changes_application_data"] is False
    assert "need not be an application-visible change" in out["weight"]


@check("source-info: AUXILIARY_DATA is decoded and qualified")
def _s2():
    out = R.decode_source_info(R.SOURCE_INFO_FLAGS["AUXILIARY_DATA"])
    assert "AUXILIARY_DATA" in out["flags"]
    assert out["changes_application_data"] is False
    assert "not changing application data" in out["weight"]


@check("source-info: replication flags are identified, not folded into 'ordinary'")
def _s3():
    out = R.decode_source_info(R.SOURCE_INFO_FLAGS["REPLICATION_MANAGEMENT"])
    assert "REPLICATION_MANAGEMENT" in out["flags"]
    assert "replication" in out["weight"]
    assert out["weight"] != "ordinary"


@check("source-info: raw value and unknown bits are preserved")
def _s4():
    out = R.decode_source_info(0x10)
    assert out["raw"] == 0x10 and out["unknown_bits"] == 0x10
    assert out["flags"] == []


@check("source-info: combined flags keep the raw value and both names")
def _s5():
    mask = R.SOURCE_INFO_FLAGS["DATA_MANAGEMENT"] | R.SOURCE_INFO_FLAGS["AUXILIARY_DATA"]
    out = R.decode_source_info(mask)
    assert out["raw"] == mask
    assert set(out["flags"]) == {"DATA_MANAGEMENT", "AUXILIARY_DATA"}


@check("source-info: a flagged record is reported, not discarded")
def _s6():
    result, _ = run(
        [state(), state()],
        [buf(5000, v2(usn=1000, source=R.SOURCE_INFO_FLAGS["AUXILIARY_DATA"]))],
        None,
    )
    assert len(result.matches) == 1
    assert result.matches[0]["source_info"]["changes_application_data"] is False


# --------------------------------------------------------------------------
# 10. Claim ceiling and report language
# --------------------------------------------------------------------------

@check("ceiling: the scan report refuses to upgrade Claim T'")
def _n1():
    result, _ = run([state(), state()], [buf(5000, v2(usn=1000))], None)
    payload = result.to_dict()
    assert payload["claim_t_prime"] == "not established by this reader"
    assert "does not establish that every journal record" in payload["completeness"]


@check("ceiling: an absent match is reported without becoming evidence of no change")
def _n2():
    result, _ = run([state(), state()], [buf(5000, v3(frn=OTHER_FRN, usn=1000))], None)
    payload = result.to_dict()
    assert payload["matches"] == []
    assert "not evidence that the file was not modified" in payload["absence_means"]


@check("gaps: non-adjacent returned USNs are reported, not interpreted")
def _g1():
    result, _ = run(
        [state(), state()],
        [buf(5000, v2(usn=1000), v2(usn=2000))],
        None,
    )
    assert result.gaps["count"] == 1
    assert "not interpreted" in result.gaps["interpretation"]


@check("gaps: adjacent records produce no gap report")
def _g2():
    first = v2(usn=1000)
    second = v2(usn=1000 + len(first))
    result, _ = run([state(), state()], [buf(5000, first, second)], None)
    assert result.gaps["count"] == 0


@check("report: every INCONCLUSIVE constant is a named outcome")
def _n3():
    assert len(R.INCONCLUSIVE_OUTCOMES) >= 14
    for name in (R.I_JOURNAL_ID_CHANGED, R.I_RETENTION_ADVANCED, R.I_NEXT_USN_REGRESSED,
                 R.I_JOURNAL_ENTRY_DELETED, R.I_JOURNAL_INACTIVE, R.I_JOURNAL_DELETE_IN_PROGRESS,
                 R.I_NO_PROGRESS_WITH_RECORDS,
                 R.I_ADVANCED_WITHOUT_RECORDS, R.I_CURSOR_REGRESSED, R.I_UNSUPPORTED_VERSION,
                 R.I_UNDOCUMENTED_REASON, R.I_MALFORMED_RECORD, R.I_MASK_MISMATCH,
                 R.I_BUDGET_CALLS, R.I_BUDGET_BYTES, R.I_BUDGET_PREFLIGHT):
        assert name in R.INCONCLUSIVE_OUTCOMES, name
        assert "INCONCLUSIVE" in name


# --------------------------------------------------------------------------
# 11. No live access
# --------------------------------------------------------------------------

@check("safety: volume access refuses without explicit authorisation")
def _z1():
    try:
        R.volume_io("\\\\.\\G:")
    except PermissionError as exc:
        assert "separate gates" in str(exc)
        return
    except RuntimeError:
        return  # not Windows: also a refusal
    raise AssertionError("volume_io opened or planned a volume without authorisation")


@check("safety: a non-raw volume path is refused")
def _z2():
    try:
        R.volume_io(r"C:\models\x.gguf", allow_open=True)
    except (ValueError, RuntimeError):
        return
    raise AssertionError("a non-volume path was accepted")


@check("safety: the module cannot construct a volume handle at import or scan time")
def _z3():
    src = open(R.__file__, encoding="utf-8").read()
    for forbidden in ("CreateFileW(", "DeviceIoControl("):
        # Both appear only inside volume_io, which is gated.
        index = src.find(forbidden)
        assert index > 0, forbidden
        gate = src.rfind("def volume_io", 0, index)
        assert gate > 0 and gate < index, f"{forbidden} appears outside volume_io"


@check("safety: the Win32 signatures are declared (static source check)")
def _z5():
    src = open(R.__file__, encoding="utf-8").read()
    for needle in (
        "kernel32.CreateFileW.restype = wintypes.HANDLE",
        "kernel32.CreateFileW.argtypes",
        "kernel32.DeviceIoControl.restype = wintypes.BOOL",
        "kernel32.DeviceIoControl.argtypes",
        "c_void_p(-1).value",
    ):
        assert needle in src, f"missing declaration: {needle}"
    # A default (int) return would truncate a 64-bit HANDLE.
    assert "restype = wintypes.HANDLE" in src


@check("safety: scan() needs injected functions and holds no handle")
def _z4():
    vol = Vol([state(), state()], [buf(5000, v2(usn=1000))])
    result = R.scan(IDENTITY, vol.query, vol.read)
    assert result.outcome == R.OUTCOME_OK
    assert vol.reads == [(1000, R.DEFAULT_BUFFER_BYTES)]


# --------------------------------------------------------------------------
# Runner
# --------------------------------------------------------------------------


def main(argv: list[str] | None = None) -> int:
    del argv
    for name, fn in CHECKS:
        try:
            fn()
        except AssertionError as exc:
            RESULTS.append((name, False, str(exc) or "assertion failed"))
        except Exception as exc:  # noqa: BLE001 - a crashing check is a failing check
            RESULTS.append((name, False, f"{type(exc).__name__}: {exc}"))
        else:
            RESULTS.append((name, True, ""))

    passed = sum(1 for _, ok, _ in RESULTS if ok)
    for name, ok, detail in RESULTS:
        print(f"{'PASS' if ok else 'FAIL'}  {name}" + (f"  — {detail}" if detail else ""))
    print(f"\n{passed}/{len(RESULTS)} checks passed")
    return 0 if passed == len(RESULTS) else 1


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
