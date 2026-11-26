"""Semantic branch — client for a llama-server embedding endpoint.

Optional by construction: if the server is not reachable, `SemanticBranch.unavailable()`
returns a branch that reports itself absent. Nothing here ever falls back to
keywords and calls it semantics.
"""
from __future__ import annotations

import json
import time
import urllib.error
import urllib.request

QUERY_PREFIX = "search_query: "
DOCUMENT_PREFIX = "search_document: "
METRIC = "cosine"


class SemanticUnavailable(RuntimeError):
    pass


class SemanticBranch:
    def __init__(self, base_url: str, model_id: str, timeout: float = 30.0) -> None:
        self.base_url = base_url.rstrip("/")
        self.model_id = model_id
        self.timeout = timeout
        self.calls = 0

    # -- transport -----------------------------------------------------------
    def _embed(self, text: str) -> list[float]:
        body = json.dumps({"input": text, "model": self.model_id}).encode()
        req = urllib.request.Request(f"{self.base_url}/v1/embeddings", data=body,
                                     headers={"Content-Type": "application/json"})
        with urllib.request.urlopen(req, timeout=self.timeout) as r:
            payload = json.loads(r.read().decode())
        self.calls += 1
        return payload["data"][0]["embedding"]

    def health(self) -> bool:
        try:
            with urllib.request.urlopen(f"{self.base_url}/health", timeout=3) as r:
                return r.status == 200
        except Exception:
            return False

    # -- branch --------------------------------------------------------------
    def score_documents(self, texts: dict[str, str]) -> dict[str, float]:
        docs = {rid: self._embed(DOCUMENT_PREFIX + t) for rid, t in texts.items()}
        return docs

    def score_query(self, query: str, docs: dict[str, list[float]]) -> dict[str, float]:
        q = self._embed(QUERY_PREFIX + query)
        out = {}
        for rid, d in docs.items():
            out[rid] = sum(a * b for a, b in zip(q, d))
        return out

    def identity(self) -> dict:
        return {"model_id": self.model_id, "endpoint": self.base_url, "metric": METRIC,
                "query_prefix": QUERY_PREFIX, "document_prefix": DOCUMENT_PREFIX,
                "dimension": None, "calls_made": self.calls}


def try_connect(base_url: str, model_id: str) -> SemanticBranch | None:
    b = SemanticBranch(base_url, model_id)
    return b if b.health() else None


def timed_embed(branch: SemanticBranch, text: str) -> tuple[list[float], float]:
    t0 = time.perf_counter()
    v = branch._embed(text)
    return v, (time.perf_counter() - t0) * 1000
