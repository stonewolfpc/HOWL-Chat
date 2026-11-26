"""Arm D — placebo construction and per-turn fidelity (plan §2.5).

D is a control only if it matches A1 **per turn**: same record count, same
serialised shape (same wrapper, bullets, provenance and uncertainty fields,
same ordering), and memory tokens within a declared tolerance. Aggregate-only
equality is explicitly not accepted.

The decoy pool is fixed and declared before the run — never improvised per case —
and every decoy must be demonstrably irrelevant to the probe: it may neither
answer it nor contradict it. A decoy that leaks relevance is an instrument
failure, because it would make D a weak A1 rather than a placebo.
"""
from __future__ import annotations

import hashlib
import itertools
from dataclasses import dataclass

from .assemble import AssembledBlock, Wrapper, shape_signature
from .errors import DecoyLeakError, PlaceboFidelityError
from .retrieve import terms

DECOY_PREFIX = "dec-"


@dataclass(frozen=True)
class Decoy:
    decoy_id: str
    text: str
    tokens: int

    def as_dict(self) -> dict:
        return {"decoy_id": self.decoy_id, "text": self.text, "tokens": self.tokens}


class DecoyPool:
    """Fixed before the run; hashed so a mid-run edit is visible."""

    def __init__(self, entries: list[dict]) -> None:
        if not entries:
            raise PlaceboFidelityError("decoy pool is empty: arm D cannot be constructed")
        self._decoys: list[Decoy] = []
        seen: set[str] = set()
        for i, entry in enumerate(entries, start=1):
            did = entry.get("decoy_id") or f"{DECOY_PREFIX}{i:04d}"
            text = (entry.get("text") or "").strip()
            tokens = entry.get("tokens")
            if not text:
                raise PlaceboFidelityError(f"decoy {did} has no text")
            if not isinstance(tokens, int) or tokens <= 0:
                raise PlaceboFidelityError(f"decoy {did} has no positive declared token count")
            if did in seen:
                raise PlaceboFidelityError(f"duplicate decoy id {did}")
            seen.add(did)
            self._decoys.append(Decoy(did, text, tokens))

    def decoys(self) -> list[Decoy]:
        return list(self._decoys)

    def sha256(self) -> str:
        blob = "\n".join(f"{d.decoy_id}|{d.tokens}|{d.text}" for d in self._decoys)
        return hashlib.sha256(blob.encode("utf-8")).hexdigest()

    def rendered_tokens(self, decoy: Decoy, counter, wrapper: Wrapper) -> int:
        """Count the decoy's *rendered bullet*, the same quantity A1 is measured on."""
        block = decoy_block([decoy], wrapper)
        return counter.count(block.bullets_text)

    def subset(self, count: int, target_tokens: int, tolerance: int, counter,
               wrapper: Wrapper) -> list[Decoy]:
        """Smallest-deviation subset of exactly `count` decoys matching A1's rendered tokens."""
        if count <= 0:
            raise PlaceboFidelityError("A1 selected no records: arm D cannot be matched per turn")
        rendered = {d.decoy_id: self.rendered_tokens(d, counter, wrapper) for d in self._decoys}
        best: tuple[int, list[Decoy]] | None = None
        for combo in itertools.combinations(self._decoys, count):
            total = sum(rendered[d.decoy_id] for d in combo)
            diff = abs(total - target_tokens)
            if diff <= tolerance and (best is None or diff < best[0]):
                best = (diff, list(combo))
                if diff == 0:
                    break
        if best is None:
            raise PlaceboFidelityError(
                f"no {count}-decoy subset matches {target_tokens} tokens within +/-{tolerance}: "
                f"arm D cannot be token-matched this turn"
            )
        return best[1]


def decoy_block(selected: list[Decoy], wrapper: Wrapper) -> AssembledBlock:
    """Render decoys through the *same* wrapper, so only content differs."""
    bullets = []
    for i, decoy in enumerate(selected, start=1):
        line = f"{wrapper.bullet}{decoy.text}"
        if wrapper.provenance:
            line += f" (source: {decoy.decoy_id})"
        if wrapper.uncertainty:
            # Same field presence and format; the numbers are the placebo's own.
            line += f" [confidence 0.00, established turn {i}]"
        bullets.append(line)
    bullets_text = "\n".join(bullets) + "\n"
    wrapper_text = wrapper.header + "\n"
    text = wrapper_text + bullets_text
    return AssembledBlock(text=text, citation_ids=[d.decoy_id for d in selected],
                          record_count=len(selected), shape=shape_signature(text, wrapper),
                          arm="D", bullets_text=bullets_text, wrapper_text=wrapper_text)


def build_placebo_measured(a1: AssembledBlock, pool: DecoyPool, wrapper: Wrapper,
                           counter, tolerance: int) -> tuple[AssembledBlock, dict]:
    """Match D to A1 using **measured** memory tokens for both sides."""
    if a1.record_count <= 0:
        raise PlaceboFidelityError("A1 block is empty; a placebo cannot be matched to nothing")
    a1_tokens = counter.count(a1.bullets_text) if a1.bullets_text else 0
    if a1_tokens <= 0:
        raise PlaceboFidelityError("A1 memory tokens are not measured; refusing to approximate")
    selected = pool.subset(a1.record_count, a1_tokens, tolerance, counter, wrapper)
    d_block = decoy_block(selected, wrapper)
    d_tokens = counter.count(d_block.bullets_text)
    if abs(d_tokens - a1_tokens) > tolerance:
        raise PlaceboFidelityError(
            f"D memory tokens {d_tokens} vs A1 {a1_tokens} outside +/-{tolerance}"
        )
    if d_block.shape != a1.shape:
        raise PlaceboFidelityError(f"shape differs: A1={list(a1.shape)} D={list(d_block.shape)}")
    report = {"a1_tokens": a1_tokens, "d_tokens": d_tokens, "tolerance": tolerance,
              "record_count": a1.record_count, "decoy_ids": d_block.citation_ids,
              "shape": list(d_block.shape), "within_tolerance": True}
    return d_block, report


def assert_irrelevant(decoy_texts: list[str], probe) -> None:
    """A decoy must neither answer nor contradict the probe (plan §2.5).

    The forbidden set is the probe's key value (the answer itself), its declared key
    terms, and its contradiction terms. Checking only `key_terms` would miss the case
    that matters most: a decoy that contains the answer.
    """
    forbidden = {t.lower() for t in
                 list(terms(probe.key_value)) + list(probe.key_terms)
                 + list(probe.contradiction_terms)
                 + [t.lower() for t in terms(probe.question)]}
    forbidden.discard("the")
    for text in decoy_texts:
        overlap = forbidden & set(terms(text))
        if overlap:
            raise DecoyLeakError(
                f"decoy leaks relevance for probe {probe.probe_id!r}: terms {sorted(overlap)}"
            )


__all__ = ["Decoy", "DecoyPool", "decoy_block", "build_placebo_measured", "assert_irrelevant"]
