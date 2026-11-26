"""Embedder adapter — interface only. No model, no download, no network.

The contract exists so that model identity is a recorded fact rather than an
assumption, and so that a missing model is a loud failure instead of a silent
keyword fallback. The app's existing `KeywordEmbedder` is exactly the failure
this adapter refuses to repeat: it is lexical frequency, not semantics, and it
must never be reported as an embedding.
"""
from __future__ import annotations

from dataclasses import dataclass, asdict
import hashlib
import pathlib


class EmbedderUnavailable(RuntimeError):
    """Raised instead of falling back. A degraded index must be visible."""


@dataclass(frozen=True)
class EmbedderIdentity:
    model_id: str
    revision: str
    file_sha256: str
    dimension: int
    max_tokens: int
    normalized: bool
    metric: str
    query_prefix: str
    document_prefix: str

    def as_dict(self) -> dict:
        return asdict(self)


def identity_for(path: pathlib.Path, **kw) -> EmbedderIdentity:
    """Pin identity from the artefact itself, never from a config string."""
    if not path.exists():
        raise EmbedderUnavailable(
            f"embedding model not present at {path}; no fallback will be used")
    digest = hashlib.sha256(path.read_bytes()).hexdigest()
    return EmbedderIdentity(file_sha256=digest, **kw)


class Embedder:
    """Implementations must produce L2-normalised vectors and declare their identity."""

    identity: EmbedderIdentity

    def embed_documents(self, texts: list[str]) -> list[list[float]]:
        raise NotImplementedError

    def embed_query(self, text: str) -> list[float]:
        raise NotImplementedError


def require(embedder: Embedder | None, index_state: str) -> None:
    """Called by evaluation. Absence degrades the index explicitly."""
    if embedder is None:
        raise EmbedderUnavailable(
            f"semantic branch disabled: {index_state}. Fused scores must not be "
            f"attributed to a semantic component while this is raised.")
