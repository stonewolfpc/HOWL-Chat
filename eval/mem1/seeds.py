"""Matched seed derivation and enforcement (prereg §4, Amend. 1 §5).

`seed(c, r) = H(conversation_id || rep_index)`. One value per cell, identical in
every arm. Enforcement is mechanical, not documentary:

  * the table is built before the run and hashed;
  * every call receipt records `seed` and `seed_honoured`;
  * `verify_cell_seeds` asserts the seed multiset per cell is identical across
    arms, and raises for a missing seed or an unhonoured one.

A cell that fails cannot be "replaced" by another conversation — that would
change the independent unit (prereg §4).
"""
from __future__ import annotations

import hashlib


def seed_for(conversation_id: str, rep_index: int) -> int:
    """Deterministic, declared before the run. Positive 31-bit for API sanity."""
    if not conversation_id:
        raise ValueError("conversation_id must not be empty")
    if rep_index < 0:
        raise ValueError("rep_index must be >= 0")
    digest = hashlib.sha256(f"{conversation_id}|{rep_index}".encode("utf-8")).hexdigest()
    return int(digest[:8], 16) % (2**31 - 1)


def build_seed_table(conversation_ids: list[str], reps: int) -> dict[tuple[str, int], int]:
    """One table, committed and hashed before the first call."""
    table: dict[tuple[str, int], int] = {}
    for cid in conversation_ids:
        for rep in range(reps):
            table[(cid, rep)] = seed_for(cid, rep)
    return table


def table_sha256(table: dict[tuple[str, int], int]) -> str:
    blob = "\n".join(f"{cid}|{rep}|{seed}" for (cid, rep), seed in sorted(table.items()))
    return hashlib.sha256(blob.encode("utf-8")).hexdigest()


def verify_cell_seeds(
    table: dict[tuple[str, int], int],
    receipts: list[dict],
) -> list[dict]:
    """Assert every arm used the table's seed for each cell, and honoured it.

    `receipts` items need: conversation_id, rep_index, arm, seed, seed_honoured.
    Raises SeedMissingError / SeedNotHonouredError; returns the verified cells.
    """
    from .errors import SeedMissingError, SeedNotHonouredError

    if not receipts:
        raise SeedMissingError("no receipts to verify: seed matching cannot be assumed")

    by_cell: dict[tuple[str, int], dict[str, dict]] = {}
    for rec in receipts:
        cid, rep = rec.get("conversation_id"), rec.get("rep_index")
        arm = rec.get("arm")
        if cid is None or rep is None or arm is None:
            raise SeedMissingError(f"receipt lacks cell identity: {sorted(rec)[:6]}")
        by_cell.setdefault((cid, int(rep)), {})[arm] = rec

    verified: list[dict] = []
    for (cid, rep), arms in sorted(by_cell.items()):
        expected = table.get((cid, rep))
        if expected is None:
            raise SeedMissingError(f"cell ({cid}, {rep}) has no seed in the committed table")
        seeds: dict[str, int] = {}
        for arm, rec in arms.items():
            seed = rec.get("seed")
            if seed is None or seed != expected:
                raise SeedMissingError(
                    f"cell ({cid}, {rep}) arm {arm}: seed {seed!r} != table {expected}"
                )
            if rec.get("seed_honoured") is not True:
                raise SeedNotHonouredError(
                    f"cell ({cid}, {rep}) arm {arm}: seed_honoured={rec.get('seed_honoured')!r}"
                )
            seeds[arm] = seed
        if len(set(seeds.values())) != 1:
            raise SeedMissingError(f"cell ({cid}, {rep}) seeds differ across arms: {seeds}")
        verified.append({"conversation_id": cid, "rep_index": rep,
                         "seed": expected, "arms": sorted(arms)})
    return verified


__all__ = ["seed_for", "build_seed_table", "table_sha256", "verify_cell_seeds"]
