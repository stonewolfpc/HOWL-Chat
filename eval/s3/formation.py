"""Layer 2 — memory formation (deterministic path).

Turns evidence spans into candidate records with provenance. This is the
deterministic extractor: the architecture guarantees retrieval works without a
successful AI extraction, so this path must stand on its own.

AI-assisted formation is a separate, optional method (`Q-MF-01` remains open).
Nothing here selects it as the final architecture; it is `formation_method =
"deterministic"` on every record it produces.

Rules are deliberately narrow and inspectable. A rule that cannot cite the span
it came from cannot produce a record.
"""
from __future__ import annotations

from dataclasses import dataclass, field
import re

from eval.s3.records import RecordStore, ACTIVE, REJECTED, SUPERSEDED
from eval.s3.evidence import EvidenceLog

MIN_CONFIDENCE = 0.55

# (name, regex, kind, predicate group, value group, confidence)
# Groups: 1 = predicate (optional), 2 = value.
_RULES: list[tuple[str, re.Pattern, str, float]] = [
    ("trained_by", re.compile(r"\bI was trained by (.+?)[.!?]?$", re.I), "fact", 0.92),
    ("name", re.compile(r"\bmy name is ([A-Z][\w'’-]*(?: [A-Z][\w'’-]*)*)", re.I), "fact", 0.9),
    ("keeps", re.compile(r"\bI keep (?:a |an |the )?(.+?)[.!?]?$", re.I), "fact", 0.7),
    ("ledger", re.compile(r"\bI keep a (ledger|journal|record) (.+?)[.!?]?$", re.I), "fact", 0.72),
    ("never", re.compile(r"\bI have never (.+?)[.!?]?$", re.I), "absence", 0.78),
    ("cannot", re.compile(r"\bI (?:cannot|can't|can not) (.+?)[.!?]?$", re.I), "absence", 0.72),
    ("preference", re.compile(r"\bI (?:like|love|prefer|enjoy|hate|dislike) (.+?)[.!?]?$", re.I),
     "preference", 0.7),
    ("promise", re.compile(r"\bI (?:promise|swear|give you my word) (.+?)[.!?]?$", re.I),
     "commitment", 0.8),
    ("served", re.compile(r"\bI served (.+?)[.!?]?$", re.I), "fact", 0.75),
]

_CORRECTION = re.compile(r"\b(?:actually|correction|I was wrong|to be clear)[,:]?\s+(.+?)[.!?]?$", re.I)
_CODE = re.compile(r"\bthe (?:vault|door|gate) code is ([a-z0-9-]+(?:\s+[a-z0-9-]+)*)", re.I)


@dataclass
class FormationStats:
    captured: int = 0
    duplicates: int = 0
    superseded: int = 0
    rejected_low_confidence: int = 0
    corrected: int = 0
    per_rule: dict[str, int] = field(default_factory=dict)

    def as_dict(self) -> dict:
        return {"captured": self.captured, "duplicates": self.duplicates,
                "superseded": self.superseded,
                "rejected_low_confidence": self.rejected_low_confidence,
                "corrected": self.corrected, "per_rule": dict(self.per_rule)}


def _normalise_value(text: str) -> str:
    return re.sub(r"\s+", " ", text.strip().rstrip(".!?")).strip()


def _candidates(span_text: str, speaker: str) -> list[dict]:
    out: list[dict] = []
    plain = span_text.strip()

    m = _CODE.search(plain)
    if m:
        out.append({"kind": "fact", "subject": speaker, "predicate": "vault_code",
                    "value": _normalise_value(m.group(1)), "confidence": 0.88,
                    "rule": "code"})

    correction = bool(_CORRECTION.search(plain))
    body = plain
    if correction:
        cm = _CORRECTION.search(plain)
        body = cm.group(1)

    for name, pattern, kind, conf in _RULES:
        m = pattern.search(body)
        if not m:
            continue
        value = _normalise_value(m.group(1) if m.lastindex else m.group(0))
        if not value:
            continue
        predicate = name
        if name == "ledger":
            predicate, value = "keeps_" + m.group(1).lower(), _normalise_value(m.group(2))
        out.append({"kind": "state_change" if correction else kind, "subject": speaker,
                    "predicate": predicate, "value": value,
                    "confidence": conf + (0.05 if correction else 0.0), "rule": name})
    # A specific rule outranks a greedy one on the same span: "I keep a ledger…"
    # must not also produce a vague `keeps` record.
    if any(c["rule"] == "ledger" for c in out):
        out = [c for c in out if c["rule"] != "keeps"]
    # A correction that produced a specific record must not also emit a vague
    # duplicate: two records for one fact is the failure the architecture names.
    if correction and not any(c["rule"] != "correction" for c in out):
        out.append({"kind": "event", "subject": speaker, "predicate": "correction",
                    "value": _normalise_value(body), "confidence": 0.6, "rule": "correction"})
    return out


class Formatter:
    """Deterministic formation over evidence spans, appended incrementally."""

    def __init__(self, store: RecordStore, log: EvidenceLog) -> None:
        self.store = store
        self.log = log
        self.stats = FormationStats()
        self.formed_spans: set[str] = set()

    def run_pending(self, log: EvidenceLog, owner_of: dict[str, str]) -> list[str]:
        """Form every span not yet processed. A checkpoint must not skip turns."""
        pending = [s for t in log.turns() for s in t.span_ids if s not in self.formed_spans]
        return self.run(pending, owner_of)

    def run(self, span_ids: list[str], owner_of: dict[str, str]) -> list[str]:
        """Form records from spans not yet processed. Returns new record ids."""
        new_ids: list[str] = []
        for sid in span_ids:
            if sid in self.formed_spans:
                continue
            self.formed_spans.add(sid)
            span = self.log.span(sid)
            owner_id = owner_of.get(span.speaker)
            if owner_id is None:
                continue  # user/persona speech is not character knowledge
            # Subject is the character name (natural key); owner_id is the stable
            # identity. Keeping them distinct is what makes history queries work.
            for cand in _candidates(span.text, span.speaker):
                rec_id = self._admit(cand, sid, span.turn, owner_id)
                if rec_id:
                    new_ids.append(rec_id)
                    self.stats.captured += 1
                    self.stats.per_rule[cand["rule"]] = self.stats.per_rule.get(cand["rule"], 0) + 1
        return new_ids

    def _admit(self, cand: dict, span_id: str, turn: int, owner_id: str) -> str | None:
        if cand["confidence"] < MIN_CONFIDENCE:
            self.store.add(kind=cand["kind"], owner_id=owner_id, owner_kind="character",
                           text=f"{cand['subject']} {cand['predicate']}: {cand['value']}",
                           subject=cand["subject"], predicate=cand["predicate"],
                           value=cand["value"], source_span=span_id, established_turn=turn,
                           confidence=cand["confidence"], status=REJECTED)
            self.stats.rejected_low_confidence += 1
            return None

        prior = [r for r in self.store.all()
                 if r.owner_id == owner_id and r.subject == cand["subject"]
                 and r.predicate == cand["predicate"]]
        active = [r for r in prior if r.status == ACTIVE]

        for r in active:
            if r.value.strip().lower() == cand["value"].strip().lower():
                self.stats.duplicates += 1
                return None  # same claim: no second record, no drift

        rec = self.store.add(kind=cand["kind"], owner_id=owner_id, owner_kind="character",
                            text=f"{cand['subject']} {cand['predicate']}: {cand['value']}",
                            subject=cand["subject"], predicate=cand["predicate"],
                            value=cand["value"], source_span=span_id, established_turn=turn,
                            confidence=cand["confidence"])
        for old in active:
            self.store.supersede(old.record_id, rec.record_id)
            self.stats.superseded += 1
            if cand["kind"] == "state_change":
                self.stats.corrected += 1
        return rec.record_id
