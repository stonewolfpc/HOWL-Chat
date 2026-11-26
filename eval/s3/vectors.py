"""Subword similarity branch.

This is NOT a semantic model. It is a hashed character-n-gram vector space: it
captures morphological overlap ("trainer" ~ "trained_by" share `tra/rai/ain/ine`)
and nothing about meaning. It is labelled `subword` everywhere so a fused score
can never be explained by a component pretending to be something else.

A true semantic branch needs an embedding model; see hybrid.SEMANTIC_STATUS.
"""
from __future__ import annotations

from dataclasses import dataclass
import hashlib
import math
import re

DIM = 2048
NGRAMS = (3, 4, 5)
SEED = "s3-subword-v1"

SEMANTIC_STATUS = {
    "available": False,
    "reason": "no embedding model authorized or configured; branch reports unavailable rather than faking a fallback",
    "metric": "cosine",
}


# Content-only n-grams. Including function words made an unrelated question match
# any record containing "the", which destroyed abstention (observed in P2/P3b).
_STOP = set("""a an the and or of to in on at for with is was were be been it its
 i you he she they we my your his her their me him them this that these those
 who what when where why how do does did have has had not no yes if then than
 about into over under after before as by from again know think""".split())


def _content_words(text: str) -> list[str]:
    return [w for w in re.findall(r"[a-z0-9']+", text.lower()) if w not in _STOP]


def _ngrams(text: str) -> list[str]:
    words = _content_words(text)
    if not words:
        return []
    t = " " + " ".join(words) + " "
    out: list[str] = []
    for n in NGRAMS:
        for i in range(max(0, len(t) - n + 1)):
            gram = t[i:i + n]
            if gram.strip():  # drop grams that are only separators
                out.append(gram)
    return out


def _bucket(gram: str) -> int:
    h = hashlib.blake2b((SEED + gram).encode(), digest_size=8).digest()
    return int.from_bytes(h, "big") % DIM


def counts(text: str) -> dict[str, int]:
    out: dict[str, int] = {}
    for g in _ngrams(text):
        out[g] = out.get(g, 0) + 1
    return out


@dataclass
class SubwordIndex:
    """IDF-weighted coverage of the query's recognizable n-grams.

    Symmetric cosine is wrong here: a short question against a long record is
    penalised by document length, which made this branch return nothing at all.
    Coverage asks the question the branch is actually good at: how much of what
    the query says, in subword terms, does the record contain?
    """

    docs: dict[str, dict[str, int]]
    idf: dict[str, float]

    @classmethod
    def build(cls, texts: dict[str, str]) -> "SubwordIndex":
        docs = {rid: counts(t) for rid, t in texts.items()}
        df: dict[str, int] = {}
        for c in docs.values():
            for g in c:
                df[g] = df.get(g, 0) + 1
        n = max(1, len(docs))
        idf = {g: math.log(1 + (n - c + 0.5) / (c + 0.5)) for g, c in df.items()}
        return cls(docs=docs, idf=idf)

    def score(self, query: str) -> dict[str, float]:
        q = counts(query)
        if not q:
            return {rid: 0.0 for rid in self.docs}
        # Unseen grams score at maximum IDF so they *penalise* coverage. Dropping
        # them made a tiny vocabulary collapse the denominator and let a single
        # accidental gram score a perfect match (observed on P3b).
        unseen = max(self.idf.values()) if self.idf else 1.0
        mass = sum(self.idf.get(g, unseen) * c for g, c in q.items())
        out: dict[str, float] = {}
        for rid, doc in self.docs.items():
            shared = sum(self.idf.get(g, unseen) * min(c, doc.get(g, 0)) for g, c in q.items())
            out[rid] = shared / mass
        return out
