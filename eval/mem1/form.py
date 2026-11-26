"""Deterministic formation from genuine evidence spans (plan §3.1, §3.3).

This is the *only* formation path the instrument implements. `RQ-MF-01`
(model-assisted formation) stays open and untested: nothing here calls a model,
and every record it makes is `formation_method="deterministic"`.

Rules are deliberately narrow and inspectable. A rule that cannot cite the span
it came from cannot produce a record — the span is passed to `MemoryStore`, which
raises `ProvenanceFabricationError` if it does not resolve.

Formation is incremental and checkpoint-free: every span is formed exactly once,
whether or not a rollover happened that turn (the S3 self-check found that class
of bug the hard way; it is asserted here as a test).
"""
from __future__ import annotations

import re
from dataclasses import dataclass, field

from .model import CHAT_OWNER, MemoryStore, ACTIVE

MIN_CONFIDENCE = 0.55

# (rule name, pattern, kind, predicate, confidence)
_RULES: tuple[tuple[str, re.Pattern, str, str, float], ...] = (
    ("trained_by", re.compile(r"\bI was trained by (.+?)[.!?]?$", re.I), "fact", "trained_by", 0.92),
    ("name", re.compile(r"\bmy name is (.+?)[.!?]?$", re.I), "fact", "name", 0.90),
    ("keeps_ledger",
     re.compile(r"\bI keep (?:a |an |the )?(ledger|journal|record)\b.*?\bin the ([\w-]+(?: [\w-]+)?) "
                r"(cabinet|drawer|chest|shelf)", re.I),
     "fact", "keeps_ledger", 0.80),
    ("vault_code", re.compile(r"\bthe vault code is ([\w\s-]+?)[.!?]?$", re.I), "fact", "vault_code", 0.88),
    ("served", re.compile(r"\bI served (.+?)[.!?]?$", re.I), "fact", "served", 0.75),
    ("preference", re.compile(r"\bI (?:like|love|prefer|enjoy|hate|dislike) (.+?)[.!?]?$", re.I),
     "preference", "preference", 0.70),
    ("absence", re.compile(r"\bI have never (.+?)[.!?]?$", re.I), "absence", "never", 0.78),
    ("commitment", re.compile(r"\bI (?:promise|swear) (.+?)[.!?]?$", re.I), "commitment", "promise", 0.80),
    ("restricted_fact",
     re.compile(r"^(?:between us|just between us|do not tell anyone)\b[,:]?\s*(?:the )?"
                r"([\w-]+(?: [\w-]+)?) is (.+?)[.!?]?$", re.I),
     "fact", "restricted_location", 0.80),
    # World facts: spoken by the narrator, not owned by a character.
    ("world_fact", re.compile(r"^(?:The|A|An) ([\w-]+(?: [\w-]+)?) (is|are|was|were) (.+?)[.!?]?$", re.I),
     "fact", "world_state", 0.60),
)

_CORRECTION = re.compile(r"\b(?:actually|correction|I was wrong|to be clear)\b[,:]?\s*(.+?)[.!?]?$", re.I)

# A retraction must be *cued* as one. The bare word "forget" is not a cue: "I forget
# where the ledger is" is a memory lapse, and "Forget the weather — the ledger is in
# the green cabinet" asserts the very fact it mentions. Both would revoke a record if
# the matcher accepted any sentence containing "forget" plus two content tokens.
_RETRACTION = re.compile(
    r"\b(?:i was wrong|that was wrong|retract(?:ed|ing)?|scratch that|never mind"
    r"|forget (?:that|what i said|everything|it)\b)", re.I)
_RESTRICTED = re.compile(r"\b(?:between us|just between us|do not tell)\b", re.I)

_NARRATOR = "Narrator"
_CHAT_OWNER = CHAT_OWNER


@dataclass
class FormationStats:
    captured: int = 0
    duplicates: int = 0
    superseded: int = 0
    rejected_low_confidence: int = 0
    corrections: int = 0
    retractions: int = 0
    per_rule: dict[str, int] = field(default_factory=dict)

    def as_dict(self) -> dict:
        return {"captured": self.captured, "duplicates": self.duplicates,
                "superseded": self.superseded,
                "rejected_low_confidence": self.rejected_low_confidence,
                "corrections": self.corrections, "retractions": self.retractions,
                "per_rule": dict(self.per_rule)}


def _norm(text: str) -> str:
    return re.sub(r"\s+", " ", text.strip().rstrip(".!?")).strip()


_STOP = set("""a an the and or of to in on at for with is was were be been it its i you
he she they we my your his her their me him them this that these those who what when
where why how do does did have has had not no yes if then than about into over under
after before as by from forget retract scratch never mind wrong said say""".split())


def _content_tokens(text: str) -> set[str]:
    return {w for w in re.findall(r"[a-z0-9']+", text.lower())
            if len(w) > 2 and w not in _STOP}


def _candidates(text: str, speaker: str) -> list[dict]:
    """Deterministic candidates for one span. No model, no heuristics beyond these rules."""
    out: list[dict] = []
    plain = text.strip()
    body = plain
    correction = bool(_CORRECTION.search(plain))
    if correction:
        body = _CORRECTION.search(plain).group(1)

    for name, pattern, kind, predicate, conf in _RULES:
        # World facts belong to the narrator. A character saying "The vault code is…"
        # is character knowledge, and forming it as world state would leak it to
        # every perspective.
        if name == "world_fact" and speaker != _NARRATOR:
            continue
        m = pattern.search(body)
        if not m:
            continue
        if name == "keeps_ledger":
            value = f"{m.group(1).lower()} in the {m.group(2)} {m.group(3)}"
        elif name == "restricted_fact":
            value = _norm(m.group(2))
            predicate = f"{m.group(1).lower().replace(' ', '_')}_location"
        elif name == "world_fact":
            value = f"{m.group(1)} {m.group(2)} {m.group(3)}"
        else:
            value = _norm(m.group(1) if m.lastindex else m.group(0))
        if not value:
            continue
        out.append({"rule": name, "kind": "state_change" if correction else kind,
                    "predicate": predicate, "value": value,
                    "confidence": conf + (0.05 if correction else 0.0)})
    return out


class Determiner:
    """Forms records from spans, once each, with provenance checked by the store."""

    def __init__(self, store: MemoryStore, evidence) -> None:
        self.store = store
        self.evidence = evidence
        self.stats = FormationStats()
        self.formed_spans: set[str] = set()

    def run_pending(self, owner_of: dict[str, str], active_world: str | None = None) -> list[str]:
        pending = [sid for turn in self.evidence.turns() for sid in turn.span_ids
                   if sid not in self.formed_spans]
        return self.run(pending, owner_of, active_world)

    def run(self, span_ids: list[str], owner_of: dict[str, str],
            active_world: str | None = None) -> list[str]:
        """Form records from unformed spans.

        `active_world` is the world active at this turn, supplied by the driver. It
        is attribution, not inference: narrator facts belong to the world that is
        active when they are said, and to the chat container when no world is
        active. Narrator content is never scanned to guess a world.
        """
        new_ids: list[str] = []
        for sid in span_ids:
            if sid in self.formed_spans:
                continue
            self.formed_spans.add(sid)
            span = self.evidence.span(sid)
            owner_id = owner_of.get(span.speaker)
            if owner_id is None and span.speaker != _NARRATOR:
                continue        # user/persona speech is not character knowledge
            if span.speaker == _NARRATOR:
                if active_world:
                    owner_id, owner_kind = active_world, "world"
                    scope, audience = "world", ("all",)
                else:
                    owner_id, owner_kind = _CHAT_OWNER, "chat"
                    scope, audience = "chat_state", ()
            else:
                scope = "restricted" if _RESTRICTED.search(span.text) else "owner"
                audience = (owner_id,) if scope == "restricted" else ()
                owner_kind = "character"

            if _RETRACTION.search(span.text) and span.speaker != _NARRATOR:
                self._retract(span.text, owner_id, span.turn, sid)
                continue

            for cand in _candidates(span.text, span.speaker):
                rid = self._admit(cand, sid, span.turn, owner_id, owner_kind, span.speaker,
                                  scope, audience)
                if rid:
                    new_ids.append(rid)
                    self.stats.captured += 1
                    self.stats.per_rule[cand["rule"]] = self.stats.per_rule.get(cand["rule"], 0) + 1
        return new_ids

    def _retract(self, text: str, owner_id: str, turn: int, source_span: str) -> None:
        """A retraction removes a claim from current truth by revocation, keeping history.

        Matching is on **content-token overlap**, not whole-string containment: a
        retraction says "forget the Ninth Company", which shares no full value string
        with "four years in the Ninth Company" but plainly names it. Two shared
        content tokens are required, so a single common word cannot revoke a fact.
        """
        said = _content_tokens(text)
        revoked = 0
        for rec in list(self.store.current(owner_id=owner_id)):
            claimed = _content_tokens(f"{rec.predicate} {rec.value} {rec.subject}")
            if len(said & claimed) >= 2:
                self.store.retract(rec.record_id, turn, source_span,
                                    reason="retracted_in_conversation")
                revoked += 1
        if revoked:
            self.stats.retractions += revoked

    def _admit(self, cand: dict, span_id: str, turn: int, owner_id: str, owner_kind: str,
               speaker: str, scope: str, audience: tuple[str, ...]) -> str | None:
        if cand["confidence"] < MIN_CONFIDENCE:
            self.stats.rejected_low_confidence += 1
            return None

        prior = [r for r in self.store.all()
                 if r.owner_id == owner_id and r.subject == speaker
                 and r.predicate == cand["predicate"]]
        active = [r for r in prior if r.status == ACTIVE]

        for rec in active:
            if _norm(rec.value).lower() == _norm(cand["value"]).lower():
                self.stats.duplicates += 1
                return None     # same claim: no second record, no drift

        supersede = None
        if active and cand["kind"] == "state_change":
            supersede = active[0].record_id
        elif active and _CORRECTION.search(cand.get("value", "")) is not None:
            supersede = active[0].record_id

        rec = self.store.add_event(
            turn=turn, via="witnessed", owner_id=owner_id, owner_kind=owner_kind,
            kind=cand["kind"],
            subject=speaker, predicate=cand["predicate"], value=cand["value"],
            source_span=span_id, confidence=cand["confidence"], scope=scope,
            audience=audience, supersede=supersede,
        )
        if supersede:
            self.stats.superseded += 1
            if cand["kind"] == "state_change":
                self.stats.corrections += 1
        return rec.record_id


__all__ = ["MIN_CONFIDENCE", "FormationStats", "Determiner"]
