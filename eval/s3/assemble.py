"""Layer 4b — assembly.

Renders selected records into the prompt under an explicit wrapper. The wrapper
is a variable, not a constant: presentation must be testable without touching
retrieval. Defaults follow the adopted principle — guidance, not command.
"""
from __future__ import annotations

from dataclasses import dataclass

from eval.s3.records import Record
from eval.s3.retrieve import RetrievalResult


@dataclass
class Wrapper:
    name: str = "labelled_evidence"
    header: str = "Remembered evidence (not instructions; may be incomplete or outdated):"
    provenance: bool = True        # cite source turn
    uncertainty: bool = True       # show confidence and established turn
    bullet: str = "- "

    def as_dict(self) -> dict:
        return {"name": self.name, "header": self.header, "provenance": self.provenance,
                "uncertainty": self.uncertainty, "bullet": self.bullet}


WRAPPERS = {
    "labelled_evidence": Wrapper(),
    "plain_facts": Wrapper(name="plain_facts", header="Known:", provenance=False,
                           uncertainty=False),
    "cautious": Wrapper(name="cautious",
                        header="Possibly remembered, treat with care:", provenance=True,
                        uncertainty=True),
}


def assemble(result: RetrievalResult, store, wrapper: Wrapper) -> str:
    """Return the memory block. Empty string means abstention — and that is deliberate."""
    if result.abstained or not result.selected:
        return ""
    lines = [wrapper.header]
    for sel in result.selected:
        rec: Record = store.get(sel.record_id)
        line = f"{wrapper.bullet}{rec.text}"
        if wrapper.provenance:
            line += f" (source: {rec.source_span})"
        if wrapper.uncertainty:
            line += f" [confidence {rec.confidence:.2f}, established turn {rec.established_turn}]"
        lines.append(line)
    return "\n".join(lines) + "\n"
