# Reference index — model identity

**Series:** Round 3 model evidence, independent reading packets
**Prepared:** 2026-11-19

Every reading packet shows the stored identifier and the model file verbatim.
This index adds the recorded file size, the SHA-256 from
`results/model_manifest.json`, and the repository string where one was recorded.

`OBSERVED` limitations:

- Quantisation is not a separate field. It is visible only inside the file name
  (`Q4_K_M`, `QAD-Q4_0`) and is reproduced there, not normalised.
- Repository identifiers were recorded only for ablation conditions, in the
  `notes` field. For baseline conditions the repository is not stored with the
  run.
- `model_manifest.json` was written at 02:12 on 2026-11-19, after the runs it
  describes. A hash taken afterwards cannot prove the file was unchanged during
  the experiment; it identifies the file present when the manifest was written.

| Model file | Identifier(s) | Kind(s) | Size on disk (bytes) | SHA-256 | Repository or preset string recorded |
|---|---|---|---:|---|---|
| `LFM2.5-1.2B-Instruct-QAD-Q4_0.gguf` | `corpus-lfm120-instruct`, `corpus-lfm120-t08`, `lfm120-ins-r2cap`, `lfm120-ins-t03`, `lfm120-ins-t05`, `lfm120-ins-t08`, `lfm120-ins-t08-rep10`, `lfm120-instruct-official`, `lfm120-instruct-official-gpufull` | corpus, triage | 695755488 | `bb741ebb106d543e9de114b843a3d3d73d51c74b5801e69da2abde821a0cb3e1` | Liquid AI card: T0.1 / top_k50 / rep1.05; Stage B: repetition penalty only, at the better temperature; Stage B: reproduces round 2 temperature under the tuned chain; Stage B: temperature only; T0.3 (one factor off Liquid's 0.1), top_k50, rep1.05; T0.5, top_k50, rep1.05; T0.8 (round-2 shared setting), top_k50, rep1.05; T0.8 with repetition penalty disabled; round-2 chain + round-2 cap, LFM at the shared sampler |
| `LFM2.5-1.2B-Thinking-Q4_K_M.gguf` | `lfm120-think-official`, `lfm120-think-official-rb1024` | triage | 730899072 | `4a1e1824f71091f368137b6966094461019f6af3a5aadd96033bee59fc9eaeea` | Liquid AI card: T0.05 / top_k50 / rep1.05; thinking always on, no effort levels, no off switch; thinking always on, no effort levels, no off switch / runtime budget 1024 (unbounded did not terminate) |
| `LFM2.5-350M-QAD-Q4_0.gguf` | `lfm350-official` | triage | 219312832 | `3d10b6ab8fc91a919534b9558e266255aca0bbc7f6d015963599aa9e74e05b1d` | Liquid AI card: T0.1 / top_k50 / rep1.05 |
| `LFM3-1.7B-Q4_K_M.gguf` | `corpus-lfm3-nonthink`, `corpus-lfm3-think-512`, `lfm3-17b-nonthink`, `lfm3-17b-think-rb1024`, `lfm3-17b-think-rb512` | corpus, triage | 1099627808 | `3793e1e6e2266919d6e015287221b41100175a98eddaba92e111aed3b05c93c4` | LFM3 card: T0.1 / top_k50 / rep1.1; reasoning default on; enable_thinking=False is vendor-flagged as inconsistent; reasoning default on; enable_thinking=False is vendor-flagged as inconsistent / budget forced by runtime, not a model effort level; reasoning default on; enable_thinking=False is vendor-flagged as inconsistent / runtime budget 1024 (unbounded did not terminate) |
| `Qwen3-1.7B-Q4_K_M.gguf` | `qwen3-17b-nonthink`, `qwen3-17b-think-rb1024` | triage | 1107409472 | `b139949c5bd74937ad8ed8c8cf3d9ffb1e99c866c823204dc42c0d91fa181897` | Qwen3 non-thinking preset T0.7/top_p0.8; Qwen3 thinking preset T0.6/top_p0.95; runtime budget 1024 (unbounded did not terminate) |
| `Qwen3.5-0.8B-Q4_K_M.gguf` | `qwen35-08b-nonthink` | triage | 532517120 | `bd258782e35f7f458f8aced1adc053e6e92e89bc735ba3be89d38a06121dc517` | Qwen3.5 small non-thinking text preset |
| `Qwen3.5-2B-Q4_K_M.gguf` | `qwen35-2b-nonthink`, `qwen35-2b-nonthink-pres0`, `qwen35-2b-nonthink-r2cap`, `qwen35-2b-nonthink-round2`, `qwen35-2b-think`, `qwen35-2b-think-rb1024` | corpus, triage | 1280835840 | `aaf42c8b7c3cab2bf3d69c355048d4a0ee9973d48f16c731c0520ee914699223` | Qwen3.5 small non-thinking text preset; Qwen3.5 small thinking text preset; Stage B: direct reproduction of the 9.33/12 round-2 condition; Stage B: isolates the presence_penalty=2.0 in Qwen's own non-thinking preset; official T1.0/top_p1.0/top_k20 with presence 0.0; round-2 chain with the round-2 400-token cap and 2048 ctx; round-2 shared chain: T0.8 / top_p0.95 / top_k40 / min_p0.05; runtime budget 1024 (unbounded did not terminate) |
| `Qwen3.5-4B-Q4_K_M.gguf` | `corpus-4b-nonthink`, `corpus-4b-think-512`, `qwen35-4b-nonthink`, `qwen35-4b-nonthink-gpufull`, `qwen35-4b-think-rb1024`, `qwen35-4b-think-rb256`, `qwen35-4b-think-rb512` | corpus, triage | 2740937888 | `00fe7986ff5f6b463e62455821146049db6f9313603938a70800d1fb69ef11a4` | Qwen3.5 4B/9B non-thinking general preset; Qwen3.5 4B/9B thinking general preset; budget forced by runtime, not a model effort level; runtime budget 1024 (unbounded did not terminate) |
| `Qwen3.5-9B-Q4_K_M.gguf` | `qwen35-9b-think-rb1024` | triage | 5680522464 | `03b74727a860a56338e042c4420bb3f04b2fec5734175f4cb9fa853daf52b7e8` | Qwen3.5 4B/9B thinking general preset; quality-ceiling candidate / runtime budget 1024 (unbounded did not terminate) |
| `SmolLM3-Q4_K_M.gguf` | `smollm3-3b-nonthink`, `smollm3-3b-think-rb1024` | triage | 1915305312 | `8334b850b7bd46238c16b0c550df2138f0889bf433809008cc17a8b05761863e` | SmolLM3 extended thinking, T0.6/top_p0.95; SmolLM3 no_think, T0.6/top_p0.95; runtime budget 1024 (unbounded did not terminate) |
| `abl-2b-q4km.gguf` | `abl-2b-nonthink` | corpus, triage | 1270808032 | `be3ccca13a9d1bc8b67165ea80bd103cd6151a37ad131183cc7ca388f78d9517` | HauhauCS Aggressive; Qwen3.5-2B non-thinking preset (T1.0/top_p1.0/pres2.0); HauhauCS/Qwen3.5-2B-Uncensored-HauhauCS-Aggressive Q4_K_M |
| `abl-4b-q4km.gguf` | `abl-4b-nonthink`, `abl-4b-think-512` | corpus, triage | 2707513696 | `79e28ecacf84e75b6056cf4059636d435aa9eb67795780f7b7dbc7d32a962741` | HauhauCS Aggressive; Qwen3.5-4B non-thinking preset (T0.7/top_p0.8/pres1.5); HauhauCS Aggressive; Qwen3.5-4B thinking preset, runtime budget 512; HauhauCS/Qwen3.5-4B-Uncensored-HauhauCS-Aggressive Q4_K_M; same weights as abl-4b-nonthink, thinking on |
| `abl-lfm120-q4km.gguf` | `abl-lfm120` | corpus, triage | 730896000 | `11731cb404b6267b63dac25710379425e2c1f42b41d0f24b7002fa6f5175530a` | Huihui abliterated; Liquid preset T0.1/top_k50/rep1.05; Huihui-LFM2.5-1.2B-Instruct-abliterated Q4_K_M |
| `abl-lfm3-q4km.gguf` | `abl-lfm3-think-512` | triage | 1099627040 | `b145feb7fcc28fd0b4d0c2ef10a0f9ba4cbde17920ef2eb717afc8e8c6209e25` | Heretic abliterated; LFM3 T0.1/top_k50/rep1.1, runtime budget 512; Vortecks/LFM3-1.7B-Heretic-Abliterated Q4_K_M |
| `tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf` | `corpus-tinyllama`, `tinyllama-11b-chat` | corpus, triage | 668788096 | `9fecc3b3cd76bba89d504f29b616eedf7da85b96540e490ca5824d3f7d2776a0` | TinyLlama card: T0.7 / top_k50 / top_p0.95; latest official TinyLlama chat (2024-01); 2K ctx, no reasoning mode |

## Labels used more than once

`qwen35-2b-nonthink` is stored as the identifier for both six-item runs and
corpus runs, in different files and with different output caps. The identifier
alone is therefore not unique across prompt types; `kind` and the source file
disambiguate. The same is true of `abl-2b-nonthink`, `abl-4b-nonthink` and
`abl-lfm120`, which appear in both the corpus and six-item sets.
