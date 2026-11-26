"""Output validity classification — instrument v2.

Audit finding F-02: in v1 an unparseable response could collect the dimensions
that reward absence, so a table of column labels scored 98 on a restraint case.
v2 makes output validity an explicit, first-class outcome, and no absence-based
credit is available unless the output is a *valid* form.

The taxonomy is deliberately wider than "parsed / not parsed":

    VALID_KEEP          >=1 contract-valid KEEP row
    VALID_HOLD          contract-valid rows, no KEEP, >=1 HOLD
    VALID_REJECT        contract-valid rows, all REJECT (added: EC-1)
    VALID_NONE          the literal contract token NONE
    VALID_EMPTY         an explicit empty-extraction declaration that is not the
                        literal NONE token ("No durable memories.")
    MALFORMED           structural scaffolding only — headers, column labels,
                        separators — with no data row at all
    UNPARSEABLE         data rows exist but none satisfies the four-field contract
    MISSING             no content at all
    INCOMPLETE          the token cap was reached before any usable row existed
    EXECUTION_FAILURE   the run itself failed; there is no output to judge

Truncation is a flag, not only a status: a run can be VALID_KEEP *and*
truncated. `quality_eligible` is False for any truncated run, and aggregations
must report both the all-runs and eligible-only figures (F-08).

No historical file is read or written here.
"""

from __future__ import annotations

import re
from dataclasses import dataclass, field
from enum import Enum

DECISIONS = ("KEEP", "REJECT", "HOLD")


class Status(str, Enum):
    VALID_KEEP = "VALID_KEEP"
    VALID_HOLD = "VALID_HOLD"
    VALID_REJECT = "VALID_REJECT"
    VALID_NONE = "VALID_NONE"
    VALID_EMPTY = "VALID_EMPTY"
    MALFORMED = "MALFORMED"
    UNPARSEABLE = "UNPARSEABLE"
    MISSING = "MISSING"
    INCOMPLETE = "INCOMPLETE"
    EXECUTION_FAILURE = "EXECUTION_FAILURE"


#: Any status that represents a conforming output.
VALID = {Status.VALID_KEEP, Status.VALID_HOLD, Status.VALID_REJECT,
         Status.VALID_NONE, Status.VALID_EMPTY}

#: Statuses that produced at least one usable decision row.
DECISIONAL = {Status.VALID_KEEP, Status.VALID_HOLD, Status.VALID_REJECT}

#: Statuses on which "you stored nothing" is a valid answer rather than a failure.
RESTRAINT_VALID = {Status.VALID_NONE, Status.VALID_EMPTY, Status.VALID_REJECT}


@dataclass
class Record:
    """One candidate row, contract-valid or not.

    `text` is always the best-effort memory+reason text, even for a contract-
    invalid row, so that semantic credit can be separated from format (F-03).
    """

    raw: str
    decision: str | None
    memory: str
    reason: str
    index: int
    contract_valid: bool
    separator: str = "|"
    issues: list[str] = field(default_factory=list)

    @property
    def text(self) -> str:
        return f"{self.memory} {self.reason}".strip()

    @property
    def field_count(self) -> int:
        return 4 if self.reason else 3


@dataclass
class ParsedOutput:
    status: Status
    records: list[Record] = field(default_factory=list)
    header_rows: int = 0
    malformed_rows: int = 0
    truncated: bool = False
    cut_mid_record: bool = False
    notes: list[str] = field(default_factory=list)

    @property
    def valid(self) -> bool:
        return self.status in VALID

    @property
    def decisional(self) -> bool:
        return self.status in DECISIONAL

    @property
    def quality_eligible(self) -> bool:
        """Usable for an unqualified quality comparison (F-08)."""
        return self.status in DECISIONAL and not self.truncated

    @property
    def kept(self) -> list[Record]:
        return [r for r in self.records if r.decision == "KEEP"]

    def decisions(self) -> list[str]:
        return [f"{r.index}:{r.decision}" for r in self.records]


_FENCE = re.compile(r"^```[a-zA-Z0-9_-]*$")
_SEPARATOR_ROW = re.compile(r"^[\s|:\-+=]+$")
_DECOR = "*`#>•-– \t"


def _strip(line: str) -> str:
    return line.strip().strip(_DECOR)


def _none_token(text: str) -> bool:
    return re.fullmatch(r"(?i)\s*none\s*[.!]?\s*", text or "") is not None


_EMPTY_DECLARATIONS = (
    "none.", "none", "no memories.", "no memories", "no durable memories.",
    "no durable memories", "nothing to store.", "nothing to store",
    "nothing worth storing.", "nothing worth storing", "no candidates.",
    "no candidates", "no memory candidates.", "no memory candidates",
    "no durable information.", "no durable information",
)


def _explicit_empty(text: str) -> bool:
    t = _strip(text).lower()
    t = re.sub(r"\s+", " ", t).strip()
    return t in _EMPTY_DECLARATIONS


def _row_fields(line: str) -> tuple[list[str], str]:
    if "|" in line:
        return [p.strip() for p in line.split("|")], "|"
    if ";" in line and len(line) > 12:
        return [p.strip() for p in line.split(";")], ";"
    return [line], ""


def _is_header(fields: list[str]) -> bool:
    head = [f.strip().strip("*`[]").upper() for f in fields[:4]]
    if not head:
        return False
    if head[0] in ("ID", "CANDIDATE", "MEMORY ID", "MEMORY_ID"):
        return True
    # A column-label row such as  C1 | DECISION | MEMORY | REASON
    return len(head) >= 3 and "DECISION" in head[1] and "MEMORY" in head[2]


def _decision(field: str) -> str | None:
    d = field.upper().strip("[]*`\"' ")
    d = re.sub(r"\s+", " ", d)
    for cand in DECISIONS:
        if d == cand:
            return cand
    return None


def parse_output(text: str | None, finish_reason: str | None = None) -> ParsedOutput:
    """Classify one model response under the v2 contract."""
    truncated = finish_reason == "length"
    if text is None:
        return ParsedOutput(Status.MISSING, truncated=truncated,
                            notes=["no content field"])

    raw_text = text.replace("\r\n", "\n")
    stripped = raw_text.strip()
    if not stripped:
        if truncated:
            return ParsedOutput(Status.INCOMPLETE, truncated=True,
                                notes=["cap reached with empty content; "
                                       "likely reasoning consumed the budget"])
        return ParsedOutput(Status.MISSING, notes=["empty content"])

    if _none_token(stripped):
        return ParsedOutput(Status.VALID_NONE, truncated=truncated,
                            notes=["literal NONE token"])

    if _explicit_empty(stripped):
        return ParsedOutput(Status.VALID_EMPTY, truncated=truncated,
                            notes=["explicit empty-extraction declaration; "
                                   "not the literal NONE token"])

    lines = [l for l in raw_text.splitlines()
             if l.strip() and not _FENCE.match(l.strip())]

    records: list[Record] = []
    header_rows = 0
    malformed = 0
    saw_data = False
    cut_mid = False

    for raw in lines:
        line = _strip(raw)
        if not line or _SEPARATOR_ROW.match(line):
            continue
        fields, sep = _row_fields(line)
        if len(fields) < 2:
            # prose / free text: keep it as a text-bearing, contract-invalid row
            saw_data = True
            malformed += 1
            records.append(Record(raw=line, decision=None, memory=line, reason="",
                                  index=len(records), contract_valid=False,
                                  separator="", issues=["not a delimited row"]))
            continue
        if _is_header(fields):
            header_rows += 1
            continue
        saw_data = True
        dec = _decision(fields[1]) if len(fields) > 1 else None
        memory = fields[2].strip() if len(fields) > 2 else ""
        reason = " | ".join(fields[3:]).strip() if len(fields) > 3 else ""
        issues = [] if dec else [f"invalid decision token {fields[1].strip()!r}"]
        if not reason:
            issues.append("missing REASON field")
        if not memory:
            issues.append("missing MEMORY field")
            cut_mid = True
        records.append(Record(raw=line, decision=dec, memory=memory, reason=reason,
                              index=len(records),
                              contract_valid=bool(dec and memory and reason),
                              separator=sep, issues=issues))

    notes: list[str] = []
    contract_rows = [r for r in records if r.contract_valid]

    if not saw_data:
        notes.append(f"{header_rows} header/label row(s), no data row")
        return ParsedOutput(Status.MALFORMED, records=records,
                            header_rows=header_rows, malformed_rows=malformed,
                            truncated=truncated, notes=notes)

    if not contract_rows:
        if header_rows:
            notes.append("header present, but no row satisfies the contract")
            status = Status.MALFORMED
        else:
            notes.append("data present, no row satisfies the four-field contract")
            status = Status.UNPARSEABLE
        return ParsedOutput(status, records=records, header_rows=header_rows,
                            malformed_rows=malformed, truncated=truncated,
                            cut_mid_record=cut_mid, notes=notes)

    decs = {r.decision for r in contract_rows}
    if "KEEP" in decs:
        status = Status.VALID_KEEP
    elif "HOLD" in decs:
        status = Status.VALID_HOLD
    else:
        status = Status.VALID_REJECT

    if truncated:
        # A truncated run is still classified, but is never quality-eligible.
        if cut_mid or not contract_rows:
            status = Status.INCOMPLETE
            notes.append("truncated mid-record")
        else:
            notes.append("truncated: valid rows present, cap reached")

    if malformed:
        notes.append(f"{malformed} row(s) failed the contract")
    if header_rows:
        notes.append(f"{header_rows} header/label row(s) not counted as records")

    return ParsedOutput(status, records=records, header_rows=header_rows,
                        malformed_rows=malformed, truncated=truncated,
                        cut_mid_record=cut_mid, notes=notes)


def classify_run(record: dict) -> ParsedOutput:
    """Classify a stored run dict from results/*.json."""
    if not isinstance(record, dict):
        return ParsedOutput(Status.EXECUTION_FAILURE,
                            notes=["record is not an object"])
    if record.get("error"):
        return ParsedOutput(Status.EXECUTION_FAILURE,
                            notes=[f"run error: {record['error']}"])
    if "content" not in record:
        return ParsedOutput(Status.EXECUTION_FAILURE,
                            notes=["record has no content field"])
    return parse_output(record.get("content"), record.get("finish_reason"))
