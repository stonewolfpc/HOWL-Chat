"""MEMORY-09 v9 engine package (Gate 2 implementation, REV 13 snapshot contract).

Implements the 57 ratified architecture rules as executable operations. It does not
execute probes and contains no conformance suite: Gate 2 authorises implementation
only, and testing remains separately unauthorised.

The G1 experiment harness in the parent package is untouched. That harness builds
records for the ablation arms; this package implements the architecture under test.
They share no state.

Absence is part of the contract: no transfer, share or carryover operation exists here
(rule 45 / K27), and no operation promotes automatically (§16). The absence is
structural -- there is no `transfer` symbol to call, not even a refusing one.

Two capabilities, one regime (rule 39 / U-8): `Grant.operations` distinguishes
`promotion` from `authorship`, `AuthorityRegistry.require` takes the operation it is
gating, and `author_ambient` is gated on `authorship` alone. Neither capability confers
the other, and revoking one leaves the other intact.
"""
from __future__ import annotations

from .authority import (BASIS_KINDS, DEFAULT_MODEL_OPERATIONS, HUMAN, MODEL, OPERATIONS,
                        AuthorityRegistry, BasisEntry, Grant, validate_basis)
from .engine import DISCLAIMER, KNOWN_VIEWS, Engine
from .records import (ACTIVE, BELIEF_DECLARERS, CONTRIBUTION_KINDS, CONTRIBUTOR_KINDS,
                      NON_TRUTH_CLASSES, ORIGIN_CLASSES, PROPOSITION_KINDS, RECORD_CLASSES,
                      REVOKED, SUPERSEDED, TRUTH_CLASSES, VIAS, Record)
from . import retrieval

ENGINE_VERSION = "9.3.0-gate2-rev13"

__all__ = ["ENGINE_VERSION", "Engine", "DISCLAIMER", "KNOWN_VIEWS",
           "AuthorityRegistry", "BasisEntry", "Grant", "validate_basis",
           "HUMAN", "MODEL", "BASIS_KINDS", "OPERATIONS",
           "DEFAULT_MODEL_OPERATIONS", "Record", "ACTIVE", "SUPERSEDED",
           "REVOKED", "TRUTH_CLASSES", "NON_TRUTH_CLASSES", "RECORD_CLASSES",
           "ORIGIN_CLASSES", "VIAS", "CONTRIBUTION_KINDS", "CONTRIBUTOR_KINDS",
           "BELIEF_DECLARERS", "PROPOSITION_KINDS", "retrieval"]
