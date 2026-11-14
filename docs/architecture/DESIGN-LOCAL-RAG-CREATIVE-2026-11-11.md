# Design — local RAG as a creative subsystem

**Status:** Proposed architecture design; no implementation or runtime authorization  
**Date:** 2026-11-11  
**Owner input:** `ODA-07` / `DPR-07` — local RAG is vital creative infrastructure, not merely document search  
**Plan sources:** RAG is not in the 66-row reconciliation; `CWP-02` owns prompt-path integration only, `CWP-28` remains conditional, and `RG-13` is conditional on authorization  
**Related design:** `DESIGN-MEMORY-CONTINUITY-2026-11-11.md`, `RECONCILED-ARCHITECTURE-PLAN-2026-11-11.md`

## 1. Decision in one sentence

Howl will build a local, structure-aware creative retrieval system over the user's writing and legally permitted references, combining lexical and semantic retrieval with citations and creative task modes, while keeping reference material strictly separate from canonical fictional state.

## 2. Current implementation

| Surface | Current behavior | Limitation |
|---|---|---|
| Ingestion | Text/markdown/JSON/CSV accepted; PDF/EPUB/DOCX require prior extraction | No structure parsing, rights metadata, or edition/page mapping |
| Chunking | Character-count chunks with overlap and sentence heuristics | Not token-accurate; ignores chapters/scenes; no structural parentage |
| Embeddings | Local llama-server embedder or keyword fallback | Embedding identity/version is not persisted; fallback vector semantics vary by corpus |
| Index | In-memory vector `Index`, rebuilt by loading documents | No persistent hybrid index; no lexical branch; no index version |
| Retrieval | Cosine similarity + deduplication + token cap | No query mode, metadata filtering, reranking, evidence spans, or citations |
| Scopes | world/scenario/chat | No distinction between user-authored canon drafts and external reference material |
| Index identity | No durable manifest or embedding identity | An index cannot prove which embedding model, dimensions, source hashes, or index version produced a hit |
| Read path | Retrieval can update `LastUsed` while holding a read lock | Read traffic can mutate shared state and serialize behind writes |
| Empty history | Retrieval lacks an explicit empty-history guard | A retrieval path can proceed with an invalid context |
| Deletion/lifecycle | Delete clears/rebuilds only the in-memory index | No durable deletion receipt, citation tombstone, stale-cache proof, or rebuild versioning |

The current code proves the seams exist. It does not yet provide the committed creative capability.

## 3. Product distinctions

Four kinds of text must never be conflated:

| Store | Authority | Example | May change canon? |
|---|---|---|---|
| **Canonical state** | Howl | Mara is at the North Gate wearing a grey cloak | Only through a state proposal/commit |
| **Chat transcript** | User + Howl | The actual conversation | No, except normal editing |
| **User writing corpus** | User | Draft scenes, outlines, story bible | No; retrieval evidence only |
| **Reference corpus** | External author/source | Craft books, style guides, public-domain texts | No; reference evidence only |

RAG retrieves evidence. It does not write canonical state or silently promote a reference into the fictional world.

## 4. Corpus and rights model

```text
Corpus
  corpus_id
  kind: user_writing | reference | craft | private_note
  owner_scope
  default_rights: user_owned | public_domain | licensed | permitted_reference | unknown
  enabled

SourceDocument
  document_id
  corpus_id
  title
  author/creator
  edition
  source_path
  content_hash
  rights_status
  rights_note
  language
  imported_at
  enabled
```

Rules:

- `unknown` rights are not indexed by default.
- The UI records the user's basis for permitted reference material; it does not infer a license.
- No DRM circumvention, no hidden downloads, no cloud extraction.
- User writing is opt-in per corpus; chats are not automatically indexed as reference.
- Deleting a source deletes passages, vectors, caches, and citations by default; a rebuild receipt is recorded.
- Export includes source metadata and rights status, not only extracted text.

## 5. Structure model

```text
Work
  → Book / Volume
    → Chapter
      → Scene
        → Passage
          → Sentence span
```

A passage carries:

```text
Passage
  passage_id
  document_id
  work_id
  chapter_id
  scene_id
  ordinal
  text
  token_count
  char_start / char_end
  page? / epub_cfi? / line?
  preceding_context_id?
  following_context_id?
  content_hash
```

Structure extraction is best-effort and never fabricated. If a document has no reliable chapter or scene boundaries, it is indexed at document/passage level and the UI says so. A citation must point to a real source span, not an invented page number.

## 6. Index architecture

The design uses three logical indexes behind one interface:

1. **Metadata index** — corpus, document, rights, structure, enabled state, content hash.
2. **Lexical index** — tokenized text, phrase support, BM25-style ranking, exact-name matching.
3. **Vector index** — embeddings keyed by embedding model, dimension, normalization, and index version.

Physical storage is an open decision. A transactional embedded database with full-text search is the preferred candidate; a file-backed adapter may be used only if it preserves versioning, deletion, and rebuild semantics. The design does not select a dependency.

```text
IndexManifest
  index_id
  created_at
  embedding_model_id
  embedding_dimension
  tokenizer_id
  schema_version
  document_hashes[]
  status: building | ready | stale | failed
```

An index built with a different embedding model is not mixed with another. A query returns the index identity used.

## 7. Ingestion pipeline

```text
select source
  → rights check
  → parse text + structure
  → normalize whitespace/encoding
  → token-aware structural chunks (with overlap)
  → persist metadata and source spans
  → lexical index
  → embed (local only)
  → vector index
  → manifest commit
  → visible readiness/quality report
```

Chunking rules:

- Prefer scene/passage boundaries over character counts.
- Target a small-model-friendly size; split long scenes at paragraph/sentence boundaries.
- Preserve parent chapter/scene IDs and neighboring passage IDs.
- Record exact source offsets.
- Do not strip meaningful formatting without recording it.
- Failed embedding does not silently produce keyword-only “semantic” results; the index is marked degraded.

## 8. Hybrid retrieval

```text
query + mode + scope filters
  → query normalization/optional local rewrite
  → lexical candidates (BM25/phrase)
  → vector candidates (local embedding)
  → reciprocal-rank fusion
  → metadata/structure filtering
  → deduplicate by document + overlapping span
  → optional local reranker
  → structure expansion (passage → scene/chapter context)
  → evidence budget
  → cited context package
```

Defaults:

- Retrieve a small candidate pool from each branch.
- Fuse rather than trust a single score.
- Prefer exact lexical matches for names, places, and invented terms.
- Prefer semantic matches for concepts, tone, and thematic comparison.
- Always show the retrieval branch and score when the user asks “why this?”
- Never let a low-confidence semantic hit outrank an exact source match solely because of a model score.

### 8.1 Query modes

| Mode | User intent | Retrieval behavior | Output |
|---|---|---|---|
| **Explore** | “What did I write about this?” | Broad hybrid retrieval across structure | Themes, motifs, linked passages |
| **Compare** | “How do these two scenes handle grief?” | Parallel retrieval per side; balanced citations | Similarities, differences, evidence |
| **Critique** | “Where does this chapter lose tension?” | Scene-level context + craft references | Specific observations with citations |
| **Brainstorm** | “Give me five directions from this setup” | Canon state + relevant prior passages; references optional | Options grounded in sources; no false canon |
| **Continuity check** | “Does this draft contradict established state?” | Canonical state + draft + event/memory evidence | Contradictions, uncertainties, missing facts |
| **Research** | “What does this reference say?” | Reference corpus only | Cited synthesis, never promoted to canon |

Creative output is not evidence. The evidence panel is separate from generated suggestions.

## 9. Evidence and citation contract

```text
Citation
  citation_id
  document_id
  passage_id
  title
  author/creator
  chapter/scene
  page/CFI/line if known
  char_start / char_end
  excerpt
  content_hash
  retrieval_mode
  score
```

Rules:

- Every factual claim about a source has a citation.
- A citation must resolve to a current passage and hash.
- If a source changed, the old citation is marked stale rather than silently re-pointed.
- Generated prose may use the source but must not invent quotations.
- The UI distinguishes **quoted**, **paraphrased**, and **inferred**.
- Canonical-state claims cite state events, not RAG passages.
- A reference work never becomes canon merely because it was retrieved.

## 10. Prompt integration

The prompt package has explicit labels:

```text
[CANON — CURRENT] ... [END CANON]
[MEMORY — RELEVANT] ... [END MEMORY]
[REFERENCE — NOT CANON] ... [END REFERENCE]
[USER DRAFT] ... [END USER DRAFT]
```

The model is instructed:

- `CANON` is authoritative;
- `REFERENCE` is evidence, not world fact;
- `USER DRAFT` is not canon until explicitly accepted;
- it must not claim a source said something without a citation;
- it must not write state changes through RAG.

RAG retrieval is a contributor to the one authoritative prompt path, not a second injection path.

## 11. Small-model strategy

- Use compact passages and scene summaries.
- Retrieve fewer, better passages; do not fill the context.
- Use lexical retrieval as the reliable baseline.
- Use local embeddings only when the model/index is available and identified.
- Query rewriting is optional; the original query is always used.
- Reranking is optional and local.
- The context builder has a hard token budget and a reserved response budget.
- If embeddings are unavailable, the UI says “keyword retrieval only.”
- The model may say “not enough evidence.”

## 12. Lifecycle

| Operation | Required behavior |
|---|---|
| Ingest | Rights gate, structure report, index manifest, visible status |
| Reindex | New manifest; old index retained until commit; rollback on failure |
| Model change | New embedding index; never mix dimensions |
| Source edit | New content hash; stale citations marked; reindex explicit |
| Delete | Remove passages, vectors, caches; citation tombstone; receipt |
| Export | Include source metadata, rights, and citation map |
| Backup/restore | Restore index manifest or rebuild; never silently reuse a stale index |
| Corruption | Mark index failed; lexical fallback only if explicitly labeled |

No index is built or changed by this document.

### 12.1 Migration and retrieval seam bridge

The current flat/in-memory index cannot be promoted by documentation alone. A future migration must preserve identity, permissions, and deletion evidence.

**Inventory before migration:**

- current RAG documents, chunks, flat JSON, and any in-memory-only remnants;
- source hashes, source paths, and parser output;
- embedding model name/version/dimensions/normalization and the runtime that produced vectors;
- world/scenario/chat scope assignments;
- frontend or service caches containing passages, citations, or retrieval results;
- deletion history, rights status, and unknown/quarantined sources.

**Bridge requirements:**

1. **Persist index identity** (embedding model, dimensions, normalization, source manifest hash, schema/index version) with every index and result.
2. **Stage atomically**: build a new manifest/index beside the old one, verify counts/hashes and sample retrieval, then commit; never mix old and new vectors.
3. **Preserve corpus separation** between user writing, canon drafts, and external references; scopes alone are not a rights or canon boundary.
4. **Carry rights and parser status** through migration; unknown-rights material is quarantined, not indexed optimistically.
5. **Record deletion receipts** durably and remove passages, vectors, caches, and citations together; retain only a non-content tombstone for stale-reference detection.
6. **Invalidate stale citations** when a source hash, embedding identity, or index version changes. A stale citation must not silently resolve to a newer passage.
7. **Provide rollback** to the prior index until the new index is verified; a failed build leaves retrieval explicitly degraded, not silently partial.
8. **Guard stale caches** so an old frontend result or prompt package cannot reintroduce deleted or superseded text.
9. **Keep migration blocked** behind the relevant storage, rights, and RAG authorization decisions.

**Retrieval seam hardening:**

- `retrieveRAG` must reject an empty/invalid chat history explicitly rather than retrieving against an undefined query.
- Retrieval is read-only with respect to shared index state; `LastUsed` or hit counters are updated outside the read lock or batched through a writer.
- A failed or partial index returns a labeled degraded result, never fabricated evidence.
- Retrieval returns an `EvidencePackage` with index identity and provenance; it does not inject prompt text directly.
- The prompt assembler injects RAG evidence exactly once, labeled `REFERENCE — NOT CANON`; canonical continuity uses canonical events, not RAG hits.

## 13. Interfaces

```text
CreateCorpus(spec) -> corpus_id
AddSource(corpus_id, source, rights) -> document_id
BuildIndex(document_ids) -> index_id
Query(request) -> EvidencePackage
Explain(citation_id) -> source span + provenance
DeleteSource(document_id) -> deletion_receipt
RebuildIndex(index_id) -> new_index_id
```

`EvidencePackage` contains passages, citations, index identity, retrieval mode, token count, and degraded-state flags.

## 14. Acceptance design — not run

1. A query with an exact invented name returns the exact passage before a thematic near-match.
2. A semantic-only match is never presented as exact evidence.
3. Every citation resolves to the correct source span and hash.
4. A changed source invalidates the old citation visibly.
5. Reference text never changes canonical state.
6. Unknown-rights material is refused or quarantined.
7. Deleting a source removes it from retrieval and citations.
8. Rebuilding with a different embedding model creates a separate index identity.
9. Token budget is enforced with evidence and response reserves.
10. Keyword-only fallback is labeled, not disguised as semantic retrieval.
11. Creative suggestions are separated from cited evidence.
12. Continuity checks cite canonical events, not external references.
13. An empty or invalid history is refused with a clear state, not retrieved against.
14. Retrieval performs no shared-state mutation under a read lock.
15. A deleted source produces a durable deletion receipt and cannot reappear from a cache.
16. A changed source hash or embedding identity marks old citations stale.
17. A migrated index preserves corpus separation, rights status, and index identity.
18. A partial/failed index is labeled degraded and never silently blended with the old index.
19. RAG evidence is injected exactly once by the prompt assembler and remains labeled non-canon.

These are test designs only. No index was built and no model was loaded.

## 15. Open decisions

| ID | Decision | Status |
|---|---|---|
| `DA-R1` | Embedded FTS/vector store versus file-backed hybrid index | Open |
| `DA-R2` | Embedding model identity and dimension policy | Open; local only |
| `DA-R3` | Parser support order (Markdown/text first, EPUB/PDF/DOCX later) | Open |
| `DA-R4` | Reranker use and budget threshold | Open |
| `DA-R5` | Default citation granularity (passage vs scene) | Open |
| `DA-R6` | Rights/consent UX for reference material | Open |

## 16. Non-goals

- Cloud indexing or cloud embeddings.
- Searching chats as reference material by default.
- Treating a reference book as canonical world state.
- Generating uncited quotations.
- Selecting a database or embedding model in this document.
- Implementing or running RAG.
