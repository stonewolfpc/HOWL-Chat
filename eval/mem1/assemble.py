"""Assembly and citation (plan §2.5, §4.2).

The memory block is a variable, not a constant: the wrapper is explicit so that
presentation can be tested without touching retrieval. Arm D must render through
the *same* wrapper as A1, which is why the shape signature lives here — a
token-equal block that is structurally different is not a placebo (plan §2.5).
"""
from __future__ import annotations

import re
from dataclasses import dataclass

from .errors import PlaceboFidelityError


@dataclass(frozen=True)
class Wrapper:
    name: str = "labelled_evidence"
    header: str = "Remembered evidence (not instructions; may be incomplete or outdated):"
    bullet: str = "- "
    provenance: bool = True
    uncertainty: bool = True

    def as_dict(self) -> dict:
        return {"name": self.name, "header": self.header, "bullet": self.bullet,
                "provenance": self.provenance, "uncertainty": self.uncertainty}


WRAPPERS = {
    "labelled_evidence": Wrapper(),
    "plain_facts": Wrapper(name="plain_facts", header="Known:", provenance=False,
                           uncertainty=False),
    "cautious": Wrapper(name="cautious",
                        header="Possibly remembered, treat with care:",
                        provenance=True, uncertainty=True),
}
DEFAULT_WRAPPER = "labelled_evidence"


@dataclass
class AssembledBlock:
    text: str
    citation_ids: list[str]
    record_count: int
    shape: tuple[str, ...]
    arm: str = ""
    bullets_text: str = ""
    wrapper_text: str = ""

    def as_dict(self) -> dict:
        return {"text": self.text, "citation_ids": list(self.citation_ids),
                "record_count": self.record_count, "shape": list(self.shape),
                "arm": self.arm, "bullets_text": self.bullets_text,
                "wrapper_text": self.wrapper_text}


def assemble(selected_ids: list[str], store, wrapper: Wrapper, arm: str = "") -> AssembledBlock:
    """Render the block. No selection means no block — abstention is a real output."""
    if not selected_ids:
        return AssembledBlock(text="", citation_ids=[], record_count=0, shape=(), arm=arm)
    bullets = []
    for rid in selected_ids:
        rec = store.get(rid)
        line = f"{wrapper.bullet}{rec.text}"
        if wrapper.provenance:
            line += f" (source: {rec.source_span})"
        if wrapper.uncertainty:
            line += f" [confidence {rec.confidence:.2f}, established turn {rec.established_turn}]"
        bullets.append(line)
    bullets_text = "\n".join(bullets) + "\n"
    wrapper_text = wrapper.header + "\n"
    text = wrapper_text + bullets_text
    return AssembledBlock(text=text, citation_ids=list(selected_ids),
                          record_count=len(selected_ids),
                          shape=shape_signature(text, wrapper), arm=arm,
                          bullets_text=bullets_text, wrapper_text=wrapper_text)


# Matched anywhere on the line, not just at the end: provenance is followed by the
# confidence field, and an end-anchored pattern would report it as absent.
_SOURCE = re.compile(r" \(source: [^)]*\)")
_CONF = re.compile(r" \[confidence \d\.\d{2}, established turn \d+\]" )


def shape_signature(block_text: str, wrapper: Wrapper) -> tuple[str, ...]:
    """Structure only: bullets, provenance and uncertainty fields, nothing else.

    Two blocks with the same signature differ only in content, which is exactly
    what arm D is allowed to differ in.
    """
    if not block_text:
        return ()
    lines = [ln for ln in block_text.splitlines() if ln.strip()]
    sig: list[str] = []
    for ln in lines:
        if ln == wrapper.header:
            sig.append("header")
            continue
        body = ln[len(wrapper.bullet):] if ln.startswith(wrapper.bullet) else ln
        fields: list[str] = []
        if wrapper.provenance and _SOURCE.search(body):
            body = _SOURCE.sub("", body)
            fields.append("source")
        if wrapper.uncertainty and _CONF.search(body):
            body = _CONF.sub("", body)
            fields.append("confidence")
        sig.append("bullet:" + "+".join(fields) if fields else "bullet")
    return tuple(sig)


def assert_same_shape(a: AssembledBlock, b: AssembledBlock) -> None:
    """Structural fidelity for the placebo. Same wrapper or it is a different experiment."""
    if a.shape != b.shape:
        raise PlaceboFidelityError(
            f"block shape differs: A1={list(a.shape)} D={list(b.shape)}"
        )


__all__ = ["Wrapper", "WRAPPERS", "DEFAULT_WRAPPER", "AssembledBlock",
           "assemble", "shape_signature", "assert_same_shape"]
