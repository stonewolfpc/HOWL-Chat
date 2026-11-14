# HOWL Chat — Memory Systems Source Register

**Status:** Research bibliography. Companion to `docs/research/MEMORY-01-RESEARCH-SYNTHESIS-2026-11-12.md`.
**Date:** 2026-11-12
**Revision:** 3 — expanded from S1–S45 to S1–S50 for the automated-lorebook research phase of 2026-11-12. New entries: S46 (SillyTavern vector storage, chat vectorization, Data Bank), S47 (KoboldAI context injection), S48 (Layla two-mode lorebook selection), S49 (llama.cpp GBNF / JSON-Schema constrained decoding), S50 (a reported constrained-decoding defect). Revision 2 added prompt-format sources (S28–S37), local creative-tool lore mechanics (S38–S40), runtime documentation (S41–S42), attention/KV-cache engineering (S43–S44), and prompt compression (S45). S1–S27 are unchanged.
**Access date for all online sources:** 2026-11-12
**Evidence labels:** `Source-derived` means the result is reported by the cited source; it does not mean this investigation reproduced it. `Preprint` and `Vendor claim` are called out explicitly.

This register records the specific source, method, result used, and limitation. It does not repeat the analysis in MEMORY-01.

---

## S1 — LongMemEval

- **Authors:** Di Wu, Hongwei Wang, Wenhao Yu, Yuwei Zhang, Kai-Wei Chang, Dong Yu
- **Title:** LongMemEval: Benchmarking Chat Assistants on Long-Term Interactive Memory
- **Date:** 2024-10-15 (arXiv v1); ICLR 2025
- **Type:** Peer-reviewed conference paper
- **URL / DOI:** https://arxiv.org/abs/2410.10813 · https://doi.org/10.48550/arXiv.2410.10813
- **Methodology:** 500 curated questions embedded in scalable user-assistant chat histories; evaluates information extraction, multi-session reasoning, knowledge update, temporal reasoning, and abstention; decomposes memory into indexing, retrieval, and reading; uses GPT-4o as judge with >97% human agreement.
- **Result used:** 30–60% long-context drop vs oracle evidence; Llama-3.1-8B full-history 0.454 vs oracle 0.710; GPT-4o 0.606 vs 0.870; 8B performance drops sharply beyond ~3k retrieved tokens, GPT-4o improves past 20k; 40–50% of errors are retrieval-correct/generation-wrong; 1B/3B models near-collapse on full history.
- **Limitations:** QA benchmark, not creative writing; 500 questions; LLM judge; finite histories (115k / 1.5M tokens); no multi-character isolation or deletion test.
- **Access:** 2026-11-12

## S2 — LoCoMo

- **Authors:** Adyasha Maharana, Dong-Ho Lee, Sergey Tulyakov, Mohit Bansal, Francesco Barbieri, Yuwei Fang
- **Title:** Evaluating Very Long-Term Conversational Memory of LLM Agents
- **Date:** 2024; ACL 2024
- **Type:** Peer-reviewed conference paper
- **URL / DOI:** https://doi.org/10.18653/v1/2024.acl-long.747 · https://arxiv.org/abs/2402.17753
- **Methodology:** 50 long multi-session dialogues, avg ~300 turns / ~9K tokens, up to 35 sessions; QA across single-hop, multi-hop, temporal, open-domain, and adversarial types, plus event summarization and multimodal dialogue; evaluates base, long-context, and RAG variants.
- **Result used:** Human QA F1 ~87.9; GPT-4-turbo ~32.1; Mistral-7B ~13.9; summary-based RAG can have high evidence recall but lower answer F1 than observation-level retrieval.
- **Limitations:** Generated/edited dialogues; F1/BLEU and LLM-judge scores are not comparable across papers; no creative continuity, character boundary, or deletion test.
- **Access:** 2026-11-12

## S3 — Zep / Graphiti

- **Authors:** Preston Rasmussen, Pavlo Paliychuk, Travis Beauvais, Jack Ryan, Daniel Chalef
- **Title:** Zep: A Temporal Knowledge Graph Architecture for Agent Memory
- **Date:** 2025-01-20
- **Type:** Preprint; authors affiliated with Zep
- **URL / DOI:** https://arxiv.org/abs/2501.13956 · https://doi.org/10.48550/arXiv.2501.13956
- **Methodology:** Bi-temporal knowledge graph with episodes, semantic facts with validity intervals, entity/community summaries, and hybrid retrieval; evaluated on DMR and LongMemEval.
- **Result used:** LongMemEval-S 71.2% vs 60.2% full-context with GPT-4o; 63.8% vs 55.4% with GPT-4o-mini; ~1.6k vs ~115k tokens; 2.58s vs 28.9s; DMR 94.8% vs 94.4% full-context.
- **Limitations:** Vendor-authored preprint; graph construction requires an LLM and a graph database; DMR is short/saturated; later vendor benchmark page (S23) reports different, higher numbers without protocol reconciliation.
- **Access:** 2026-11-12

## S4 — Mem0

- **Authors:** Prateek Chhikara, Dev Khant, Saket Aryan, Taranjeet Singh, Deshraj Yadav
- **Title:** Mem0: Building Production-Ready AI Agents with Scalable Long-Term Memory
- **Date:** 2025-04-28
- **Type:** Preprint; authors affiliated with Mem0
- **URL / DOI:** https://arxiv.org/abs/2504.19413 · https://doi.org/10.48550/arXiv.2504.19413
- **Methodology:** LLM extraction plus ADD/UPDATE/DELETE/NOOP consolidation over a vector store; optional graph memory; LoCoMo evaluation against full-context, RAG, MemGPT, A-Mem, LangMem, Zep, and OpenAI memory.
- **Result used:** Overall LLM-judge score 66.88 (Mem0) and 68.44 (Mem0g) vs 72.90 full-context; 1,764 / 3,616 memory tokens vs 26,031 full-context; 91–92% lower p95 latency; Mem0 loses to full-context on overall quality.
- **Limitations:** Vendor-authored preprint; GPT-4o-mini as extraction engine; LoCoMo only; LLM judge; later managed-platform claims (S24) are a different product/protocol; no character isolation or deletion test.
- **Access:** 2026-11-12

## S5 — APEX-MEM

- **Authors:** Pratyay Banerjee, Masud Moshtaghi, Shivashankar Subramanian, Amita Misra, Ankit Chadha
- **Title:** APEX-MEM: Agentic Semi-Structured Memory with Temporal Reasoning for Long-Term Conversational AI
- **Date:** 2026-07; ACL 2026
- **Type:** Peer-reviewed conference paper
- **URL / DOI:** https://aclanthology.org/2026.acl-long.749/ · https://doi.org/10.18653/v1/2026.acl-long.749
- **Methodology:** Domain-agnostic property graph of temporally grounded events, append-only storage, and a multi-tool retrieval agent that resolves conflicting/evolving information at query time; evaluated on LoCoMo and LongMemEval.
- **Result used:** 88.88% LoCoMo and 86.2% LongMemEval, reported with frontier QA agents; the result does not establish that a graph alone is sufficient.
- **Limitations:** Frontier models, not small local models; LLM-heavy graph construction/retrieval; no independent replication; no character isolation or deletion test.
- **Access:** 2026-11-12

## S6 — Memora / FAMA

- **Authors:** Md Nayem Uddin, Kumar Shubham, Eduardo Blanco, Chitta Baral, Gengyu Wang
- **Title:** From Recall to Forgetting: Benchmarking Long-Term Memory for Personalized Agents
- **Date:** 2026; ACL 2026 Findings
- **Type:** Peer-reviewed conference paper (Findings)
- **URL / DOI:** https://aclanthology.org/2026.findings-acl.1337/ · https://doi.org/10.18653/v1/2026.findings-acl.1337
- **Methodology:** Proposes Forgetting-Aware Memory Accuracy (FAMA) to penalize reliance on obsolete/invalidated memory; evaluates four LLMs and six memory agents.
- **Result used:** Memory agents frequently reuse invalid memories and fail to reconcile evolving memories under FAMA.
- **Limitations:** QA benchmark; no creative continuity; no product-scale data; results are source-reported.
- **Access:** 2026-11-12

## S7 — ReverieMem / KBF-QA

- **Authors:** Xushuo Tang, et al.
- **Title:** ReverieMem: Building Character-Centric Long-Form Memory for Role-Playing Agents
- **Date:** 2026-06-24
- **Type:** Preprint (arXiv)
- **URL / DOI:** https://arxiv.org/abs/2606.25632 · https://doi.org/10.48550/arXiv.2606.25632
- **Methodology:** Defines Factual Overreach; proposes first-person episodic memory, visibility-tagged semantic facts, and a presence roster; introduces KBF-QA (4,386 questions over eight novels) and BookWorld narrative evaluation.
- **Result used:** +34.6 percentage points knowledge-boundary fidelity over the strongest prior method; ~79% win rate on BookWorld; shared retrieval causes spoilers/omniscient knowledge.
- **Limitations:** Preprint; no code; no independent replication; novel-based rather than open creative writing; does not establish local small-model behavior.
- **Access:** 2026-11-12

## S8 — MemLeak

- **Authors:** Priyanka Mudgal, et al.
- **Title:** MemLeak: Incidental Memory Leakage in LLM Agents
- **Date:** 2026-10-03
- **Type:** Preprint; NeurIPS 2026 workshop
- **URL / DOI:** https://arxiv.org/abs/2610.04195 · https://doi.org/10.48550/arXiv.2610.04195
- **Methodology:** Multi-user shared-memory benchmark with pooled retrieval; measures incidental cross-user leakage and tests metadata filtering vs hard ownership gating.
- **Result used:** 70–100% incidental leakage; 90–100% top-k placement for crafted memories; metadata filtering alone insufficient; hard post-retrieval ownership gating restores clean baseline.
- **Limitations:** Synthetic, small-n; multi-user rather than multi-character; magnitudes indicative; no local model/hardware test.
- **Access:** 2026-11-12

## S9 — Authorization-First Retrieval

- **Authors:** Rohith Namboothiri
- **Title:** Authorization-First Retrieval: Enforcing Least Privilege in Multi-Agent RAG Systems
- **Date:** 2026; TrustNLP 2026
- **Type:** Peer-reviewed workshop paper
- **URL / DOI:** https://doi.org/10.18653/v1/2026.trustnlp-main.15 · https://aclanthology.org/2026.trustnlp-main.15/
- **Methodology:** Compares retrieve-then-filter with authorization-constrained candidate retrieval on a synthetic multi-department corpus; measures structural exposure and behavioral leakage.
- **Result used:** 86.1% structural exposure with retrieve-then-filter; 0% structural exposure with authorization-first; behavioral leakage still 41.3%/29.5% depending on model.
- **Limitations:** Synthetic corpus; two reader models; workshop paper; no character/presence model; no local model test.
- **Access:** 2026-11-12

## S10 — GateMem

- **Authors:** Zhe Ren, Yibo Yang, Yimeng Chen, Zijun Zhao, Benshuo Fu, Zhihao Shu, Bingjie Zhang, Yangyang Xu, Dandan Guo, Shuicheng Yan
- **Title:** GateMem: Benchmarking Memory Governance in Multi-Principal Shared-Memory Agents
- **Date:** 2026-06-17
- **Type:** Preprint (arXiv)
- **URL / DOI:** https://arxiv.org/abs/2606.18829 · https://doi.org/10.48550/arXiv.2606.18829
- **Methodology:** Jointly evaluates utility for legitimate long-horizon requests, access control across authorization boundaries, and active forgetting after deletion; spans medical, office, education, and household domains with hidden checkpoints and leak annotations.
- **Result used:** No method simultaneously achieves strong utility, robust access control, and reliable forgetting; long-context prompting has the best governance at high token cost; retrieval/external memory still leaks unauthorized or deleted information.
- **Limitations:** Synthetic institutional episodes; preprint; no creative-writing domain; no local desktop constraints.
- **Access:** 2026-11-12

## S11 — MemGPT / Letta

- **Authors:** Charles Packer, Sarah Wooders, Kevin Lin, Vivian Fang, Shishir G. Patil, Ion Stoica, Joseph E. Gonzalez
- **Title:** MemGPT: Towards LLMs as Operating Systems
- **Date:** 2023-10-12
- **Type:** Preprint (arXiv)
- **URL / DOI:** https://arxiv.org/abs/2310.08560 · https://doi.org/10.48550/arXiv.2310.08560
- **Methodology:** OS-inspired virtual context management; main/recall/archival memory; model edits memory via tools; evaluated on DMR and document QA.
- **Result used:** DMR 93.4% with GPT-4-turbo; in Zep's comparison, full-context DMR is 94.4%, so the task is saturated.
- **Limitations:** Short DMR conversations; tool-call discipline is a failure mode; no small-model benchmark; no character/deletion test.
- **Access:** 2026-11-12

## S12 — Generative Agents

- **Authors:** Joon Sung Park, Joseph C. O'Brien, Carrie J. Cai, Meredith Ringel Morris, Percy Liang, Michael S. Bernstein
- **Title:** Generative Agents: Interactive Simulacra of Human Behavior
- **Date:** 2023; UIST 2023
- **Type:** Peer-reviewed conference paper
- **URL / DOI:** https://doi.org/10.1145/3586183.3606763 · https://arxiv.org/abs/2304.03442
- **Methodology:** Memory stream with recency/importance/relevance retrieval and reflection; 25-agent sandbox evaluation with human judges.
- **Result used:** Demonstrates believable multi-agent behavior from episodic memory + reflection; not a factual-fidelity benchmark.
- **Limitations:** Small sandbox; no correction, deletion, cross-character privacy, or local-model evaluation; believability is not truth.
- **Access:** 2026-11-12

## S13 — BEIR

- **Authors:** Nandan Thakur, Nils Reimers, Andreas Rücklé, Abhishek Srivastava, Iryna Gurevych
- **Title:** BEIR: A Heterogeneous Benchmark for Zero-shot Evaluation of Information Retrieval Models
- **Date:** 2021; NeurIPS 2021 Datasets and Benchmarks
- **Type:** Peer-reviewed benchmark paper
- **URL / DOI:** https://arxiv.org/abs/2104.08663 · https://doi.org/10.48550/arXiv.2104.08663
- **Methodology:** 18 zero-shot IR datasets across diverse tasks/domains; compares lexical, dense, sparse, and reranking methods.
- **Result used:** No single architecture wins; BM25 is a strong zero-shot baseline; cross-encoder reranking improves average nDCG by ~11% over BM25 at ~450 ms CPU/query in the reported setup.
- **Limitations:** Not chat/personal memory; heterogeneous qrels; latency hardware-dependent; no temporal/state modeling.
- **Access:** 2026-11-12

## S14 — RAPTOR

- **Authors:** Parth Sarthi, Salman Abdullah, Aditi Tuli, Shubh Khanna, Anna Goldie, Christopher D. Manning
- **Title:** RAPTOR: Recursive Abstractive Processing for Tree-Organized Retrieval
- **Date:** 2024; ICLR 2024
- **Type:** Peer-reviewed conference paper
- **URL / DOI:** https://arxiv.org/abs/2401.18059 · https://doi.org/10.48550/arXiv.2401.18059
- **Methodology:** Recursively clusters and summarizes chunks into a tree; retrieves across summary levels; evaluated on QA datasets including QuALITY.
- **Result used:** Reports gains from hierarchical retrieval, including a large QuALITY gain with GPT-4; comparison changes the reader, so the retrieval-effect size is confounded.
- **Limitations:** Summary construction cost; reader confound; no independent replication of the retrieval-effect size; no chat memory or character scope.
- **Access:** 2026-11-12

## S15 — RULER

- **Authors:** Cheng-Ping Hsieh, Simeng Sun, Samuel Kriman, Shantanu Acharya, Dima Rekesh, Fei Jia, Yang Zhang, Boris Ginsburg
- **Title:** RULER: What's the Real Context Size of Your Long-Context Language Models?
- **Date:** 2024; COLM 2024
- **Type:** Peer-reviewed conference paper
- **URL / DOI:** https://arxiv.org/abs/2404.06654 · https://doi.org/10.48550/arXiv.2404.06654
- **Methodology:** Synthetic long-context tasks (retrieval, multi-hop tracing, aggregation, QA) at varying lengths; measures effective context vs claimed context.
- **Result used:** Effective context is often far below the claimed window; aggregation and multi-hop tracing degrade most.
- **Limitations:** Synthetic tasks; not memory architecture; model set predates 2026; no creative writing.
- **Access:** 2026-11-12

## S16 — LightMem

- **Authors:** Jizhan Fang, Xinle Deng, Haoming Xu, Ziyi Tang, Zican Dong, Huihui Xu, Yuanhao Wu, Junjie Wu, Xin Chen, Runyuan He, et al.
- **Title:** LightMem: Lightweight and Efficient Memory-Augmented Generation
- **Date:** 2025-10-21 (arXiv); ICLR 2026
- **Type:** Peer-reviewed conference paper
- **URL / DOI:** https://arxiv.org/abs/2510.18866 · https://doi.org/10.48550/arXiv.2510.18866
- **Methodology:** Sensory/short-term/long-term memory stages with sleep-time offline consolidation and topic-aware compression; evaluated on LongMemEval and LoCoMo.
- **Result used:** Up to 10.9% accuracy gain, 117× token reduction, 159× fewer API calls, >12× runtime reduction; reported with GPT/Qwen backbones.
- **Limitations:** Efficiency results do not prove a small local model can perform the same consolidation; no character/deletion test; source-reported.
- **Access:** 2026-11-12

## S17 — Memory-R1

- **Authors:** Sikuan Yan, Xiufeng Yang, Zuchao Huang, Ercong Nie, Zifeng Ding, Zonggen Li, Xiaowen Ma, Jinhe Bi, Kristian Kersting, Jeff Z. Pan, Hinrich Schuetze, Volker Tresp, Yunpu Ma
- **Title:** Memory-R1: Enhancing Large Language Model Agents to Manage and Utilize Memories via Reinforcement Learning
- **Date:** 2025-08-28 (arXiv); ACL 2026
- **Type:** Peer-reviewed conference paper
- **URL / DOI:** https://doi.org/10.18653/v1/2026.acl-long.583 · https://arxiv.org/abs/2508.19828
- **Methodology:** Trains a memory manager with reinforcement learning to perform ADD/UPDATE/DELETE/NOOP and a memory reader to use retrieved memories; evaluates at 3B–14B scale.
- **Result used:** Learned memory management can improve over heuristic pipelines at small-model scale; supports the idea that memory policy can be trained rather than prompted.
- **Limitations:** Benchmark-specific training; no independent replication; no creative continuity, character scope, or deletion audit; training cost.
- **Access:** 2026-11-12

## S18 — When Facts Change

- **Authors:** Jonas Wallat, Wolfgang Nejdl, Sandipan Sikdar
- **Title:** When Facts Change: Temporal Knowledge Conflict Resolution in LLMs
- **Date:** 2026; ACL 2026 Findings
- **Type:** Peer-reviewed conference paper (Findings)
- **URL / DOI:** https://doi.org/10.18653/v1/2026.findings-acl.103 · https://aclanthology.org/2026.findings-acl.103/
- **Methodology:** Evaluates whether models detect and act on changing facts, including temporal conflicts, across model sizes and memory conditions.
- **Result used:** Smaller models often fail to detect temporal conflicts; larger models may detect conflict but still fail to act on mutability.
- **Limitations:** QA-style temporal factuality, not creative retcons; no character perspective; no local deployment study.
- **Access:** 2026-11-12

## S19 — Lost in the Middle

- **Authors:** Nelson F. Liu, Kevin Lin, John Hewitt, Ashwin Paranjape, Michele Bevilacqua, Fabio Petroni, Percy Liang
- **Title:** Lost in the Middle: How Language Models Use Long Contexts
- **Date:** 2024; TACL 12:157–173
- **Type:** Peer-reviewed journal article
- **URL / DOI:** https://doi.org/10.1162/tacl_a_00638 · https://arxiv.org/abs/2307.03172
- **Methodology:** Controlled multi-document QA and key-value retrieval at varying context positions; measures accuracy by position.
- **Result used:** Performance is often U-shaped: information at the beginning/end is used better than information in the middle; more context is not uniformly better.
- **Limitations:** Synthetic/controlled retrieval tasks; model set predates 2026; not a memory architecture.
- **Access:** 2026-11-12

## S20 — Character-LLM

- **Authors:** Yunfan Shao, Linyang Li, Junqi Dai, Xipeng Qiu
- **Title:** Character-LLM: A Trainable Agent for Role-Playing
- **Date:** 2023; EMNLP 2023
- **Type:** Peer-reviewed conference paper
- **URL / DOI:** https://aclanthology.org/2023.emnlp-main.814/ · https://doi.org/10.18653/v1/2023.emnlp-main.814
- **Methodology:** Fine-tunes a separate agent per character on curated role-specific experiences; evaluates character fidelity and knowledge boundaries.
- **Result used:** Supports separate per-role agents to avoid cross-role knowledge collision; demonstrates that role isolation is a modeling choice, not only a prompt.
- **Limitations:** One model per character is impractical for many dynamic characters; no entry/exit/death lifecycle; no local small-model study.
- **Access:** 2026-11-12

## S21 — AgentPoison

- **Authors:** Zhaorun Chen, Zhen Xiang, Chaowei Xiao, Dawn Song, Bo Li
- **Title:** AgentPoison: Red-teaming LLM Agents via Poisoning Memory or Knowledge Bases
- **Date:** 2024; NeurIPS 2024
- **Type:** Peer-reviewed conference paper
- **URL / DOI:** https://doi.org/10.52202/079017-4136 · https://arxiv.org/abs/2407.12784
- **Methodology:** Optimizes poisoned memories/knowledge entries that are retrieved for attacker-chosen queries; evaluates across agent tasks.
- **Result used:** Very low insertion rates can persistently steer agent behavior through retrieval; poisoning is a memory-integrity problem, not only a prompt-injection problem.
- **Limitations:** Specific agent setups; not creative writing; no local desktop app; defense evidence limited.
- **Access:** 2026-11-12

## S22 — MINJA / memory injection

- **Authors:** Shen Dong, Shaochen Xu, Pengfei He, Yige Li, Jiliang Tang, Tianming Liu, Hui Liu, Zhen Xiang
- **Title:** Hijacking Memory: Prompt Injection Attack to Long-Term Memory of LLM Agents
- **Date:** 2025; NeurIPS 2025
- **Type:** Peer-reviewed conference paper
- **URL / DOI:** https://doi.org/10.52202/085713-1554 · https://arxiv.org/abs/2503.03704
- **Methodology:** Query-only interaction that induces an agent to write malicious content into long-term memory, then triggers later misuse; evaluates multiple memory systems.
- **Result used:** Persistent memory can be hijacked through ordinary interaction; provenance/trust labels and write review are needed.
- **Limitations:** Specific agent/memory systems; no creative-writing domain; defense evidence limited.
- **Access:** 2026-11-12

## S23 — Zep vendor benchmark page

- **Authors / organization:** Zep
- **Title:** Memory Benchmarks
- **Date:** 2026 (page)
- **Type:** Vendor claim / marketing page
- **URL:** https://www.getzep.com/research/memory-benchmarks/
- **Methodology:** Vendor-run benchmark comparisons; no full protocol/harness in the page.
- **Result used:** Reports 90.2% LongMemEval and 94.7% LoCoMo for Zep; not reconciled with the 2025 paper's 71.2%/60.2% protocol.
- **Limitations:** Vendor-authored; protocol differences; not independently verified; should not be compared to the paper's numbers.
- **Access:** 2026-11-12

## S24 — Mem0 vendor documentation

- **Authors / organization:** Mem0
- **Title:** Memory Evaluation
- **Date:** 2026 (docs)
- **Type:** Vendor documentation / claim
- **URL:** https://docs.mem0.ai/core-concepts/memory-evaluation
- **Methodology:** Vendor benchmark documentation; describes LoCoMo/LongMemEval results for the managed platform.
- **Result used:** Reports 94.4% LongMemEval and 92.5% LoCoMo for the managed platform; different from the 2025 paper's 66.88/68.44 LoCoMo scores.
- **Limitations:** Vendor-run; managed platform vs open-source paper; protocol differences; not independently verified.
- **Access:** 2026-11-12

## S25 — LangMem documentation

- **Authors / organization:** LangChain
- **Title:** LangMem SDK
- **Date:** 2025–2026 (docs)
- **Type:** Vendor documentation; no benchmark
- **URL:** https://langchain-ai.github.io/langmem/
- **Methodology:** Library documentation for memory tools/store; no controlled evaluation.
- **Result used:** Confirms a lightweight library approach with hot-path/background memory tools; no vendor benchmark to cite.
- **Limitations:** No quality evidence; third-party numbers are not attributable to the vendor.
- **Access:** 2026-11-12

## S26 — Cognee documentation

- **Authors / organization:** Cognee
- **Title:** Cognee — Memory for AI Agents / open-source repository
- **Date:** 2026
- **Type:** Vendor documentation / repository; vendor self-describes as non-neutral
- **URL:** https://github.com/topoteretes/cognee · https://www.cognee.ai/
- **Methodology:** ECL (extract-cognify-load) pipeline with relational, vector, graph, and summary layers; vendor reports LoCoMo ~91% and BEAM-100K 0.79.
- **Result used:** Illustrates the operational complexity and vendor-reported performance of a full graph+vector+summary stack; the 10M-token run is not reproducible from public material.
- **Limitations:** Vendor-reported; no independent replication; local GLiNER path exists but its accuracy is not publicly established; heavy infrastructure.
- **Access:** 2026-11-12

## S27 — GraphRAG

- **Authors:** Darren Edge, Ha Trinh, Newman Cheng, Joshua Bradley, Alex Chao, Apurva Mody, Steven Truitt, Jonathan Larson
- **Title:** From Local to Global: A Graph RAG Approach to Query-Focused Summarization
- **Date:** 2024-04-24
- **Type:** Technical report (Microsoft Research); not peer-reviewed
- **URL / DOI:** https://arxiv.org/abs/2404.16130 · https://doi.org/10.48550/arXiv.2404.16130
- **Methodology:** LLM-generated entity graph and community summaries; map-reduce query-focused summarization; evaluated with LLM judges on comprehensiveness/diversity.
- **Result used:** Graph/community summaries can improve global thematic answers, but construction is LLM-heavy and the report does not establish retrieval recall/nDCG or local-model feasibility.
- **Limitations:** Vendor technical report; LLM-judge evaluation; no independent replication; high indexing/query cost.
- **Access:** 2026-11-12

---

## Known source conflicts preserved

| Conflict | What differs | Handling |
|---|---|---|
| **Zep paper vs Zep vendor page** | 2025 paper: 71.2% LongMemEval, 60.2% full-context; vendor page: 90.2% LongMemEval, 94.7% LoCoMo | Both reported; vendor page treated as an unreconciled claim (S23) |
| **Mem0 paper vs Mem0 docs** | Paper: 66.88 LoCoMo, below 72.90 full-context; docs: 92.5 LoCoMo, 94.4 LongMemEval | Both reported; managed platform vs paper treated as different protocols (S4, S24) |
| **LoCoMo metrics across papers** | F1/BLEU vs LLM-judge; different reader models and prompts | No cross-paper numerical comparison |
| **APEX-MEM vs small-model evidence** | APEX reports frontier QA results; LongMemEval shows 8B limits | APEX treated as architecture evidence, not small-model evidence |
| **ReverieMem vs product reality** | Preprint reports +34.6 pp KBF; no code/independent replication | Treated as the strongest direct evidence, with preprint caveat |
| **MemLeak/AFR vs Howl** | Multi-user/multi-principal synthetic settings; Howl is single-user but multi-character | Direction (hard scoping) treated as robust; magnitudes indicative |
| **Summary RAG** | Some papers report summary gains; LoCoMo/LongMemEval show detail loss | Summary treated as index/routing, not sole evidence |

---

## S28 — llama.cpp server flags and behaviour

- **Authors:** ggml-org / llama.cpp contributors
- **Title:** `tools/server/README.md`
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official technical documentation
- **URL:** https://github.com/ggml-org/llama.cpp/blob/master/tools/server/README.md
- **Methodology:** Reference documentation for server flags.
- **Result used:** `--chat-template` accepts a custom template, but the documented text states that only commonly used templates are accepted **unless `--jinja` is set before the flag**; `--reasoning on|off|auto`, `--reasoning-format`, `--reasoning-budget`, `--reasoning-budget-message`; `--fit` / `--fit off`; `--flash-attn`; `-b`/`-ub` batch sizes; parallel slots share the context budget; `/apply-template` endpoint exists for rendering a conversation without inference.
- **Limitations:** Version-dependent; flags change between builds. A specific claim must be checked against the pinned backend commit.
- **Access:** 2026-11-12

## S29 — llama.cpp supported chat templates

- **Authors:** ggml-org / llama.cpp contributors
- **Title:** Templates supported by `llama_chat_apply_template` (wiki)
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official technical documentation
- **URL:** https://github.com/ggml-org/llama.cpp/wiki/Templates-supported-by-llama_chat_apply_template
- **Methodology:** Lists built-in template names and their rendered output, and documents the bundled `minja` Jinja parser.
- **Result used:** Built-in names include `chatml`, `llama2`, `llama3`, `gemma`, `command-r`, `chatglm3`, `chatglm4`, `phi3`, `phi4`, and others; rendered examples show exact role markers and turn terminators per family.
- **Limitations:** Wiki, not a versioned spec; the rendered examples are illustrative and may lag the model cards. Does not enumerate every template family.
- **Access:** 2026-11-12

## S30 — llama.cpp function calling

- **Authors:** ggml-org / llama.cpp contributors
- **Title:** `docs/function-calling.md`
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official technical documentation
- **URL:** https://github.com/ggml-org/llama.cpp/blob/master/docs/function-calling.md
- **Methodology:** Documents native tool-use support and the requirement for the right per-model template; ships example templates for DeepSeek-R1-Distill, Functionary, Hermes, FireFunction.
- **Result used:** Native tool support requires the model's own template; a generic fallback (`--chat-template chatml`) "may consume more tokens and be less efficient" and is explicitly caveated ("YMMV").
- **Limitations:** Documentation, not a measurement; per-model template quality is not benchmarked.
- **Access:** 2026-11-12

## S31 — llama.cpp multi-GPU

- **Authors:** ggml-org / llama.cpp contributors
- **Title:** `docs/multi-gpu.md`
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official technical documentation
- **URL:** https://github.com/ggml-org/llama.cpp/blob/master/docs/multi-gpu.md
- **Methodology:** Describes row-split and layer-split strategies, tensor overrides, and MoE expert placement.
- **Result used:** Multi-GPU is not a single configuration; row vs layer split have different balance characteristics; `-ot`/expert overrides allow expert tensors on CPU while attention stays on GPU.
- **Limitations:** Hardware-dependent; no performance table that transfers across GPUs.
- **Access:** 2026-11-12

## S32 — llama.cpp VRAM fitting and KV-cache constraints

- **Authors:** ggml-org / llama.cpp contributors
- **Title:** `tools/fit-params/README.md` and server `--fit-target` documentation
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official technical documentation
- **URL:** https://github.com/ggml-org/llama.cpp/blob/master/tools/fit-params/README.md
- **Methodology:** Documents the automatic fitting path, its target free-memory value, and KV-cache type constraints.
- **Result used:** `--fit-target` defaults to a **fixed 1024 MiB**, not a percentage; quantized V-cache requires Flash Attention; context shift is off by default since PR #15416.
- **Limitations:** Version-specific; the default has changed historically. Verify against the pinned commit.
- **Access:** 2026-11-12

## S33 — Gemma prompt structure and system role

- **Authors:** Google (Gemma documentation)
- **Title:** Gemma formatting and system instructions
- **Date:** 2025-03-21 (page marked deprecated in favour of gemma.ai)
- **Type:** Official technical documentation
- **URL:** https://ai.google.dev/gemma/docs/core/prompt-structure
- **Methodology:** Specifies the Gemma control tokens and role handling for instruction-tuned models.
- **Result used:** Gemma IT models are trained with **only two roles, `user` and `model`**; a separate `system` role is **not supported**, and system-level instructions must be placed inside the initial user prompt.
- **Limitations:** Documentation for an older doc location; Gemma 3 specifics should be re-checked against the current model card.
- **Access:** 2026-11-12

## S34 — DeepSeek V3.x chat and tool wire format

- **Authors:** Compilation of the DeepSeek model-card chat template and the vLLM tool-template reference
- **Title:** DeepSeek tool-calling wire format (model card "Chat Template"/"ToolCall" sections; `tool_chat_template_deepseekv31.jinja`)
- **Date:** 2025–2026; accessed 2026-11-12
- **Type:** Technical reference compiled from primary model artifacts
- **URL:** https://github.com/atyrode/omp/blob/main/docs/toolconv/deepseek.md (mirrors the model card and tokenizer data)
- **Methodology:** Token-level description of the DeepSeek envelope, with token IDs from `tokenizer.json` and the reference chat template.
- **Result used:** Fullwidth-pipe special tokens for user/assistant roles; **no OpenAI-style system token** — system text is concatenated into a flat prefix; dedicated tool-call open/close/separator tokens; a thinking channel with ` thinking`/`</think>`; fill-in-the-middle tokens distinct from chat.
- **Limitations:** Secondary compilation; the model card and tokenizer are the authority. Marker bytes must be read from the tokenizer rather than transcribed.
- **Access:** 2026-11-12

## S35 — Cohere Command R / R+ prompting

- **Authors:** Cohere
- **Title:** Prompting Command R and R+; Command R+ model card; system-message overview
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official vendor documentation
- **URL:** https://docs.cohere.com/docs/prompting-command-r · https://huggingface.co/CohereLabs/c4ai-command-r-plus-08-2024
- **Methodology:** Documents the Command-R turn tokens, preamble sections, and tool-use prompt structure.
- **Result used:** Turns begin with `<|START_OF_TURN_TOKEN|>` followed by `<|USER_TOKEN|>`, `<|SYSTEM_TOKEN|>`, or `<|CHATBOT_TOKEN|>` and end with `<|END_OF_TURN_TOKEN|>`; the preamble is structured (safety preamble, system preamble, user preamble, available tools) and tool output uses a `<results>` block; deviating from the template is documented to reduce performance.
- **Limitations:** Vendor documentation; behaviour under a mismatched template is warned about but not quantified.
- **Access:** 2026-11-12

## S36 — Qwen3 chat template and thinking control

- **Authors:** Alibaba Qwen; Hugging Face template; Qwen documentation
- **Title:** Qwen3 `chat_template.jinja`; "Run Qwen with llama.cpp"
- **Date:** 2025–2026; accessed 2026-11-12
- **Type:** Official model artifact + vendor documentation
- **URL:** https://huggingface.co/unsloth/Qwen3-32B/blob/main/chat_template.jinja · https://qwen.readthedocs.io/en/latest/run_locally/llama.cpp.html
- **Methodology:** Template source and run documentation.
- **Result used:** ChatML-based markers with an `enable_thinking` flag; the generation prompt can close the thinking block when disabled; Qwen documents that the **hard thinking switch in the chat template is not exposed in llama.cpp** and recommends a custom template file to force non-thinking; `--jinja` is required to use the embedded template.
- **Limitations:** Template evolves between Qwen versions (3, 3.5, 3.6 differ); verify per model file.
- **Access:** 2026-11-12

## S37 — Llama 3 / 3.1 chat template

- **Authors:** Meta; reproduced in the llama.cpp template wiki and HF tokenizer configs
- **Title:** Llama 3 / 3.1 chat template
- **Date:** 2024; accessed 2026-11-12
- **Type:** Official model artifact
- **URL:** https://github.com/ggml-org/llama.cpp/wiki/Templates-supported-by-llama_chat_apply_template
- **Methodology:** Template source rendered by llama.cpp.
- **Result used:** `<|start_header_id|>role<|end_header_id|>` blocks terminated by `<|eot_id|>`; Llama 3.1 adds tool-call support with an `ipython` role for tool results; system content is a normal header block.
- **Limitations:** 3.2/3.3 variants differ in tool handling; verify per file.
- **Access:** 2026-11-12

## S38 — SillyTavern World Info

- **Authors:** SillyTavern project
- **Title:** SillyTavern documentation — World Info
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official project documentation
- **URL:** https://github.com/SillyTavern/SillyTavern-Docs
- **Methodology:** Documents lorebook entry fields and activation behaviour.
- **Result used:** Keyword keys, secondary keys with AND/OR/NOT logic, regex keys, vectorized matching, recursion between entries, inclusion groups, probability, position/depth, and sticky/cooldown/delay; entry and global token budgets with budget-based dropping.
- **Limitations:** Documentation of mechanics only; no published precision/recall measurement.
- **Access:** 2026-11-12

## S39 — KoboldCpp

- **Authors:** KoboldAI / LostRuins
- **Title:** KoboldCpp documentation
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official project documentation
- **URL:** https://koboldcpp.net/docs/
- **Methodology:** Documents generation settings, context shift, and the separate long-term memory (TextDB) feature.
- **Result used:** Context shift keeps generation running past the window; world-info equivalent plus a separate vector/text memory store; soft/hard token limits.
- **Limitations:** Documentation; no published quality measurement.
- **Access:** 2026-11-12

## S40 — NovelAI Lorebook

- **Authors:** NovelAI
- **Title:** Lorebook documentation
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official vendor documentation
- **URL:** https://docs.novelai.net/en/text/lorebook/
- **Methodology:** Documents lorebook entry activation and context insertion.
- **Result used:** Key lists, regex, insertion order/position, search range, forced activation, activation probability, and context budget with over-budget cutting.
- **Limitations:** Vendor documentation; no quality measurement.
- **Access:** 2026-11-12

## S41 — LM Studio

- **Authors:** LM Studio
- **Title:** Documentation — context overflow and GPU offload
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official vendor documentation
- **URL:** https://lmstudio.ai/docs
- **Methodology:** Documents runtime configuration for local models.
- **Result used:** Context overflow policies (stop at limit, truncate middle, roll over); GPU offload layer control; no lorebook subsystem.
- **Limitations:** Product documentation, not a measurement.
- **Access:** 2026-11-12

## S42 — Ollama

- **Authors:** Ollama
- **Title:** Documentation — Modelfile, templates, context
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official vendor documentation
- **URL:** https://docs.ollama.com
- **Methodology:** Documents template handling and request behaviour.
- **Result used:** Chat formatting is template-driven via the Modelfile; the runtime does **not** maintain memory across requests and does not auto-summarize history; context is bounded by the configured window.
- **Limitations:** Vendor documentation; behaviour changes across versions.
- **Access:** 2026-11-12

## S43 — FlashAttention

- **Authors:** Tri Dao, Daniel Y. Fu, Stefano Ermon, Atri Rudra, Christopher Ré
- **Title:** FlashAttention: Fast and Memory-Efficient Exact Attention with IO-Awareness
- **Date:** 2022; NeurIPS 2022
- **Type:** Peer-reviewed conference paper
- **URL / DOI:** https://arxiv.org/abs/2205.14135 · https://doi.org/10.48550/arXiv.2205.14135
- **Methodology:** IO-aware exact attention kernel; benchmarks attention memory and speed.
- **Result used:** Exact attention with substantially lower memory traffic and memory footprint than a naive implementation; the basis for the engine's memory savings and for quantized-V support.
- **Limitations:** Speedups are hardware- and sequence-dependent; not a memory-architecture result.
- **Access:** 2026-11-12

## S44 — KV-cache compression family (H2O, SnapKV, KIVI)

- **Authors:** Zhang et al. (H2O); Li et al. (SnapKV); Liu et al. (KIVI)
- **Title:** H2O: Heavy-Hitter Oracle for Efficient Generative Inference; SnapKV: LLM Knows What You Are Looking For Before Generation; KIVI: A Tuning-Free Asymmetric 2bit Quantization for KV Cache
- **Date:** 2023–2024; NeurIPS 2023 (H2O), NeurIPS 2024 (SnapKV), ICML 2024 (KIVI)
- **Type:** Peer-reviewed conference papers
- **URL / DOI:** https://arxiv.org/abs/2306.14048 · https://arxiv.org/abs/2404.14469 · https://arxiv.org/abs/2402.02750
- **Methodology:** KV-cache eviction (H2O), observation-window selection (SnapKV), and low-bit KV quantization (KIVI).
- **Result used:** Large KV-memory reductions with reported limited quality loss on the evaluated tasks; the trade is task- and model-dependent.
- **Limitations:** Evaluated on specific models/tasks; no creative-continuity evaluation; eviction can drop exactly the detail a later question needs.
- **Access:** 2026-11-12

## S45 — Prompt/context compression (LLMLingua, Gist tokens)

- **Authors:** Jiang et al. (LLMLingua); Mu et al. (Gist tokens)
- **Title:** LLMLingua: Compressing Prompts for Accelerated Inference; Learning to Compress Prompts with Gist Tokens
- **Date:** 2023; EMNLP 2023 (LLMLingua), NeurIPS 2023 (Gist)
- **Type:** Peer-reviewed conference papers
- **URL / DOI:** https://arxiv.org/abs/2310.05736 · https://arxiv.org/abs/2304.08467
- **Methodology:** Token-level prompt compression using a small language model; learned soft-prompt compression.
- **Result used:** Substantial prompt-token reduction with reported task performance retention on the evaluated benchmarks.
- **Limitations:** Compression can remove exact strings and numbers; requires an extra model; not evaluated for narrative continuity or provenance.
- **Access:** 2026-11-12

## S46 — SillyTavern Vector Storage, chat vectorization, and Data Bank

- **Authors:** SillyTavern project
- **Title:** Chat Vectorization; Data Bank (RAG); World Info
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official project documentation
- **URL:** https://docs.sillytavern.app/extensions/chat-vectorization · https://docs.sillytavern.app/usage/core-concepts/data-bank · https://docs.sillytavern.app/usage/core-concepts/worldinfo
- **Methodology:** Documents the vector extension, its activation rules, and its relationship to World Info.
- **Result used:** Vector matching for world-info entries requires enabling the extension and marking entries vectorized, with a `Max Entries` cap; lorebooks can be bound to a character, a persona, or an **individual chat**, and chat-bound entries do not carry to other chats; dynamic prompt sources (World Info, summarization, chat vectorization) restructure the prompt prefix and **can cause frequent cache misses**; vector summarization is explicitly marked experimental.
- **Limitations:** Documentation of mechanics; no precision/recall measurement; explicit disclaimer that vectorization "does not guarantee a better chatting experience or improved memory".
- **Access:** 2026-11-12

## S47 — KoboldAI: Memory, Author's Note and World Info

- **Authors:** KoboldAI project
- **Title:** Memory, Author's Note and World Info (wiki)
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official project documentation
- **URL:** https://github.com/KoboldAI/KoboldAI-Client/wiki/Memory%2C-Author%27s-Note-and-World-Info
- **Methodology:** Explains the three context-injection features and their token economics.
- **Result used:** World Info entries are inserted only when their keywords appear in the story text, to conserve tokens; the model does not see the entry key or title, so each entry must be a **self-contained description**; suggested entry sizes are roughly 50 tokens for minor material, 100 for significant characters, and no more than 150 for major world elements.
- **Limitations:** Community wiki; guidance rather than measurement.
- **Access:** 2026-11-12

## S48 — Layla AI lorebook documentation

- **Authors:** Layla Network
- **Title:** Import Tavern PNG Cards & SillyTavern Lorebooks; "Layla supports SillyTavern Lorebooks"
- **Date:** 2024-09-14 (blog); docs rolling
- **Type:** Official vendor documentation / product blog
- **URL:** https://www.layla-network.ai/character-cards · https://www.layla-network.ai/post/layla-supports-sillytavern-lorebooks
- **Methodology:** Describes two lore mechanisms in one product and how imported lorebooks behave.
- **Result used:** Imported SillyTavern lorebooks activate by keyword rules "the same way SillyTavern does"; Layla's own lorebooks (PDF/ePub imports) select content by **embedding similarity**; the vendor states rule-based lorebooks are "more rigid but predictable" while similarity-based selection is less controllable.
- **Limitations:** Vendor description, not a measurement; no quality data.
- **Access:** 2026-11-12

## S49 — llama.cpp GBNF and JSON-Schema constrained decoding

- **Authors:** ggml-org / llama.cpp contributors
- **Title:** GBNF Guide (`grammars/README.md`); LLGuidance support (`docs/llguidance.md`)
- **Date:** rolling; accessed 2026-11-12
- **Type:** Official technical documentation
- **URL:** https://github.com/ggml-org/llama.cpp/blob/master/grammars/README.md · https://github.com/ggml-org/llama.cpp/blob/master/docs/llguidance.md
- **Methodology:** Documents grammar-constrained generation and JSON-Schema-to-GBNF conversion.
- **Result used:** A subset of JSON Schema converts to GBNF; tokens violating the grammar are given zero probability, so output **parses**; the schema is **not injected into the prompt** and the model has no visibility into it; the general `json.gbnf` accepts any object and lets the model invent field names, so the shape must be constrained explicitly; grammars are enforced in the sampler, and a grammar token that spans a boundary is legal in full or unavailable; documented performance gotchas exist; LLGuidance is an alternative with high JSON-Schema coverage and measured token-mask costs.
- **Limitations:** Documentation; grammar enforcement does not guarantee semantic correctness, and truncated generation can still yield incomplete output.
- **Access:** 2026-11-12

## S50 — llama.cpp constrained-decoding defect report

- **Authors:** llama.cpp issue reporter and maintainers
- **Title:** JSON Schema grammar fails on chat-template control tokens while equivalent GBNF succeeds (Granite 3.1 / Qwen3) — issue #29006
- **Date:** 2026-09-17
- **Type:** Public issue report (implementation evidence)
- **URL:** https://github.com/ggml-org/llama.cpp/issues/29006
- **Methodology:** Reproducible failure comparison between the chat endpoint and the raw completion endpoint.
- **Result used:** JSON-Schema-constrained decoding failed on the `/v1/chat/completions` path for two model families while an equivalent hand-written GBNF succeeded, indicating the schema-to-grammar conversion can produce false negatives around chat-template control tokens.
- **Limitations:** Single report, version-specific; not a systematic study. Cited as evidence that constrained decoding has real implementation defects, not that it is unusable.
- **Access:** 2026-11-12

## Source coverage check

| Area | Sources |
|---|---|
| Long-term conversational memory | S1, S2, S3, S4, S5, S6 |
| Structured/event/temporal | S3, S5, S6, S18 |
| Episodic/semantic/consolidation | S1, S2, S11, S12, S16 |
| Retrieval/reranking/graph/hierarchy | S13, S14, S27, S1, S3 |
| Agent frameworks | S3, S4, S11, S24, S25, S26 |
| Character/multi-principal isolation | S7, S8, S9, S10, S20 |
| Small/local models | S1, S15, S16, S17, S18, S19 |
| Provenance/poisoning | S10, S21, S22 |
| Automated lore (separate) | Project F3-C; S38, S39, S40 (mechanics); S1, S2 (summary/detail evidence); no direct benchmark |
| Prompt formats and templates | S28, S29, S30, S33, S34, S35, S36, S37 |
| Local runtime behaviour | S28, S31, S32, S41, S42 |
| Attention / KV-cache engineering | S32, S43, S44 |
| Prompt compression | S45, S16, S14 |
| Context flush / rollover | S28, S31, S39 (context shift), S41 (overflow policies); no direct benchmark |
| Hardware / VRAM / offload | S28, S31, S32, S41, S43, S44 |
| Lorebook activation mechanics (products) | S38, S40, S46, S47, S48 |
| Constrained / structured output | S49, S50 |
| Lore maintenance & deletion | S3, S5, S6, S10, S44 |

**No source in this register is a Howl runtime measurement.** All Howl-specific statements remain project records or inference, as labeled in MEMORY-01.
