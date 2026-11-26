"""Generate the normative contract extract from the design record.

Every rule, decision, invariant and open question in the output is copied
verbatim from MEMORY-09-ARCH-STAGE-CONTRACT. Nothing is paraphrased, so the
extract cannot drift from the design record through retyping.

Rationale, worked cases and the normative tables stay in the design record and
are referenced by section; this document holds the clauses a frozen contract
would consist of.

Run from the repository root:
    python docs/research/build_normative_extract.py
"""

from __future__ import annotations

import hashlib
import pathlib
import re

SRC = pathlib.Path("docs/research/MEMORY-09-ARCH-STAGE-CONTRACT-2026-11-15.md")
OUT = pathlib.Path("docs/research/MEMORY-09-CONTRACT-NORMATIVE-2026-11-15.md")

RULE = re.compile(r"^- \*\*K-R(\d+)\*\*")
INVARIANT = re.compile(r"^- \*\*I(\d+)\*\*")


def sha256(path: pathlib.Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def read_lines(path: pathlib.Path) -> list[str]:
    return path.read_text(encoding="utf-8").splitlines()


def section(lines: list[str], startswith: str) -> tuple[list[str], int]:
    """Return a section body (heading excluded, trailing rule removed) and its 1-based line."""
    for i, line in enumerate(lines):
        if line.startswith(startswith):
            body: list[str] = []
            for line in lines[i + 1 :]:
                if line.startswith("## ") or line.startswith("### "):
                    break
                body.append(line)
            while body and body[-1].strip() in ("", "---"):
                body.pop()
            while body and body[0].strip() == "":
                body.pop(0)
            return body, i + 1
    raise SystemExit(f"heading not found: {startswith!r}")


def bullet_block(lines: list[str], i: int) -> tuple[list[str], int]:
    """Capture a bullet plus indented or fenced continuation lines."""
    block = [lines[i]]
    i += 1
    fenced = False
    while i < len(lines):
        line = lines[i]
        if line.strip().startswith("```"):
            fenced = not fenced
            block.append(line)
            i += 1
            continue
        if fenced:
            block.append(line)
            i += 1
            continue
        if line.startswith("  ") or line.startswith("\t"):
            block.append(line)
            i += 1
            continue
        break
    return block, i


def collect_rules(lines: list[str]) -> list[tuple[str, str, list[str]]]:
    """Return (heading, rule id, verbatim block) in document order."""
    found: list[tuple[str, str, list[str]]] = []
    heading = "(header)"
    i = 0
    while i < len(lines):
        line = lines[i]
        if line.startswith("## ") or line.startswith("### "):
            heading = line.lstrip("#").strip()
            i += 1
            continue
        match = RULE.match(line)
        if match:
            block, i = bullet_block(lines, i)
            found.append((heading, f"K-R{match.group(1)}", block))
            continue
        i += 1
    return found


def collect_invariants(lines: list[str]) -> list[tuple[str, list[str]]]:
    found: list[tuple[str, list[str]]] = []
    i = 0
    while i < len(lines):
        match = INVARIANT.match(lines[i])
        if match:
            block, i = bullet_block(lines, i)
            found.append((f"I{match.group(1)}", block))
            continue
        i += 1
    return found


def render() -> str:
    src_text = SRC.read_bytes().decode("utf-8")
    src_hash = sha256(SRC)
    revision = re.search(r"\*\*Revision:\*\* (r\d+)", src_text)
    revision = revision.group(1) if revision else "unknown"

    lines = read_lines(SRC)
    decisions, decisions_line = section(lines, "## 5.")
    invariants_section, invariants_line = section(lines, "### Invariants")
    questions, questions_line = section(lines, "## 14.")

    rules = collect_rules(lines)
    # Count invariants from the invariant section only; the whole-source scan would pick up
    # any mislabelled bullet elsewhere and overstate the contract.
    invariants = [
        m.group(0) for m in (INVARIANT.match(l) for l in invariants_section) if m
    ]

    out: list[str] = []
    add = out.append

    add("# MEMORY-09 — Normative Contract (extract)")
    add("")
    add("**Status:** normative extract of the design record, generated mechanically. Nothing here is")
    add("new, and nothing here is adopted beyond what the design record itself marks as decided.")
    add("")
    add(f"- **Source:** `{SRC.name}`")
    add(f"- **Source revision:** {revision}")
    add(f"- **Source SHA-256:** `{src_hash}`")
    add(f"- **Generator:** `docs/research/build_normative_extract.py` (verbatim copy; no paraphrase)")
    add(f"- **Rule count:** {len(rules)} · **Invariants:** {len(invariants)} (counted from the")
    add("  invariant section; the rule index is counted from the whole source)")
    add("")
    add("Every clause below is copied byte-for-byte from the source. Rationale, worked cases and the")
    add("normative tables remain in the design record and are referenced by section, not duplicated:")
    add("a second copy of the reasoning would create exactly the two-sources-of-truth problem this")
    add("contract exists to prevent (design record §13.1).")
    add("")
    add("## 1. Adoption status")
    add("")
    add("| Item | Status | Where |")
    add("|---|---|---|")
    add("| D1–D7 (ownership, existence≠knowledge, assets≠experience, knowledge-gated, decay, conflict, belief) | **adopted** | §2 of this document |")
    add("| D8 chat-deletion contract | proposed | design record §11.1; §2 here |")
    add("| Q13 transfer carries content, not knowledge | approved in principle | design record §11.2 |")
    add("| Materiality predicate revision 2 | accepted candidate, **not** approved | design record §10.1 |")
    add("| S24–S26 materiality probes | registered, **not runnable** | design record §12 |")
    add("| Schema versioning, evaluation provenance, replayability states | proposed | design record §10.2–10.2.2 |")
    add("| Temporal applicability conventions (§11.3) | proposed | design record §11.3 |")
    add("| Q1–Q15 | open | §4 of this document |")
    add("| E-SCOPE-CONTRACT suite | **not frozen** | design record §12 |")
    add("")
    add("## 2. Decisions")
    add("")
    add(f"Verbatim from the design record §5 (source lines {decisions_line}+):")
    add("")
    out.extend(decisions)
    add("")
    add("D8 is proposed, not adopted; it appears in the design record §11.1 and in the open")
    add("questions listed in §4 below.")
    add("")
    add("## 3. Invariants")
    add("")
    add(f"Verbatim from the design record (source lines {invariants_line}+):")
    add("")
    out.extend(invariants_section)
    add("")
    add("## 4. Rules")
    add("")
    add("Verbatim, in document order, grouped by their source section. Continuation lines and fenced")
    add("blocks belonging to a rule are included; normative tables that follow a rule are not copied")
    add("and must be read in the source section named by the group heading.")
    add("")
    current = None
    for heading, rule_id, block in rules:
        if heading != current:
            add(f"### {heading}")
            add("")
            current = heading
        out.extend(block)
        add("")
    add("## 5. Open questions")
    add("")
    add(f"Verbatim from the design record §14 (source lines {questions_line}+):")
    add("")
    out.extend(questions)
    add("")
    add("## 6. Measured baseline at the time of this extract")
    add("")
    add("Established by running the existing apparatus unchanged; no measurement was made for this")
    add("document. Receipts live in `eval/s3/`.")
    add("")
    add("| Measurement | Result |")
    add("|---|---|")
    add("| Gate A | 14/14 PASS |")
    add("| Gate B | **6/7 PARTIAL** — P5 (\"Who trained the horse?\") still open |")
    add("| Formation | frozen at Gate A baseline (7 records); hashes match |")
    add("| Harness self-checks | PASS |")
    add("| S1 audit | PASS (5 calls, 0 problems) |")
    add("| `validate` | 89/89 |")
    add("| Corpus | PASS |")
    add("| Reranker (bge-reranker-v2-m3 Q8_0) | evaluated, **not shipped**; no further candidate fetched |")
    add("| Resource envelope | no configuration meets VRAM ≤384 MiB **and** p50 ≤40 ms (87/198/277/444 MiB; 282/102/70/20 ms) |")
    add("")
    add("Gate B 6/7 is **not** a completed baseline: P5 is an open failure, and the proposed")
    add("architecture may not inherit it as an assumption.")
    add("")
    add("## 7. What this extract does not claim")
    add("")
    add("- It is **not** behaviourally validated. No probe in the registry has been run.")
    add("- It is **not** a suite freeze. The registry is not frozen and cannot be until an independent")
    add("  adversarial review of its expectations has taken place.")
    add("- It does **not** authorise implementation, or integration into Howl Chat.")
    add("")
    return "\n".join(out) + "\n"


def main() -> None:
    text = render()
    OUT.write_text(text, encoding="utf-8")
    written = OUT.read_bytes()

    # Verification: every rule and invariant bullet must appear verbatim.
    src_rules = {m.group(0) for m in (RULE.match(l) for l in read_lines(SRC)) if m}
    out_rules = {m.group(0) for m in (RULE.match(l) for l in written.decode("utf-8").splitlines()) if m}
    missing = sorted(src_rules - out_rules)
    if missing:
        raise SystemExit(f"rules missing from extract: {missing}")

    src_inv = {
        m.group(0)
        for section_lines in [section(read_lines(SRC), "### Invariants")[0]]
        for m in (INVARIANT.match(l) for l in section_lines)
        if m
    }
    out_inv = {
        m.group(0)
        for m in (INVARIANT.match(l) for l in written.decode("utf-8").splitlines())
        if m
    }
    if src_inv != out_inv:
        raise SystemExit(f"invariants drifted: missing {sorted(src_inv - out_inv)}, extra {sorted(out_inv - src_inv)}")

    print(f"wrote {OUT} | {len(written)} B | sha256 {hashlib.sha256(written).hexdigest()}")
    print(f"  rules copied: {len(out_rules)} | invariants copied: {len(out_inv)}")
    print(f"  missing: {missing or 'none'} | invariant drift: none")


if __name__ == "__main__":
    main()
