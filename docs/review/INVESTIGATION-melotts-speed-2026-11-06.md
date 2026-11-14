# MeloTTS Speed Investigation

**Date:** 2026-11-06  
**Scope:** Read-only investigation of MeloTTS speed controls for Group 2A.  
**Status:** Investigation complete. No production or test source changed.

## Verdict

`speed` is the authoritative MeloTTS control. `length_scale` is not an independent public MeloTTS parameter; official MeloTTS derives it internally as `1 / speed`. The adapter's legacy `/synthesize` payload currently sends both, which is unnecessary and can only create an extra-field rejection or, in a nonstandard wrapper, double application.

No live MeloTTS deployment was available to inspect. The repository has no configured `HOWL_MELOTTS_URL`, no process is listening on the default port `8888`, the Python `melo` package is not installed, and Docker is not available on this host. Therefore the adapter's actual remote schema and the deployed server's real accepted range remain unverified.

## 1. Local Environment Evidence

Read-only checks performed:

- `HOWL_MELOTTS_URL`: unset
- `HOWL_MELOTTS_API_KEY`: unset
- `HOWL_MELOTTS_OPENAI_MODEL`: unset
- Listening ports checked: `8888`, `8000`, `8080`, `7860`, `5000` — no listener
- Python import `melo`: `ModuleNotFoundError`
- Docker: command not available
- Working config: `ttsEnabled=false`, `ttsEngine=piper`, `ttsSpeed=1`

There is no live deployment in this environment to query for `/tts/openapi.json`, `/tts/docs`, or `/synthesize` behavior. The findings below therefore distinguish source-established facts from deployment-dependent behavior.

## 2. Current Adapter Behavior

Source: `internal/audio/tts/melotts/client.go`.

### Legacy `POST /synthesize`

The adapter sends:

```json
{
  "text": "...",
  "language": "EN",
  "speaker_id": "EN-US",
  "speed": 1.0,
  "length_scale": 1.0
}
```

`length_scale` is calculated as `1 / clampSpeechSpeed(speed)`. The values are always consistent with each other, but the field is not part of the official MeloTTS API and is not needed once `speed` is sent.

The adapter falls back to `/v1/audio/speech` only on HTTP `404` or `405`. A `422` caused by an unexpected `length_scale` field would not trigger fallback; it would fail the request.

### OpenAI-compatible `POST /v1/audio/speech`

The adapter sends `speed` only, clamped to `0.25–4.0`:

```json
{
  "model": "tts-1",
  "input": "...",
  "voice": "EN-US",
  "response_format": "wav",
  "speed": 1.0
}
```

This path has no competing field.

### Shared validation

`BaseSynthesizer.ValidateRequest` rejects nonzero speeds outside `0.5–2.0`. This runs before either HTTP request. The effective adapter range is therefore `0.5–2.0` for both routes, regardless of the wider `0.25–4.0` clamp in the OpenAI payload builder.

## 3. Official MeloTTS Behavior

Source: <https://github.com/myshell-ai/MeloTTS/blob/main/melo/api.py>

The public Python method is:

```python
tts_to_file(..., speed=1.0, ...)
```

Inside that method, the model receives:

```python
length_scale=1. / speed
```

The official implementation therefore treats `speed` as the user-facing multiplier and computes `length_scale` internally. There is no `length_scale` argument on the public method and no HTTP `length_scale` field in the official API.

The official Web UI exposes a speed slider from `0.1` to `10.0`, default `1.0`:

<https://github.com/myshell-ai/MeloTTS/blob/main/melo/app.py>

The official docs demonstrate `--speed 1.5` and `speed=1.0` but do not state a hard API range:

<https://github.com/myshell-ai/MeloTTS/blob/main/docs/install.md>

This means the official library does not enforce the repository's `0.8–1.2` profile bounds. It also does not define `length_scale` as an HTTP control.

## 4. Common HTTP Deployments

### Hangry Labs port-8888 fork

Sources:

- <https://github.com/hangry-labs/MeloTTS>
- <https://github.com/hangry-labs/MeloTTS/blob/main/melo/schemas.py>
- <https://github.com/hangry-labs/MeloTTS/blob/main/melo/openai_compat.py>

Its native request schema has:

```python
speed: float = Field(1.0, ge=0.5, le=2.0)
```

It has no `length_scale` field. Pydantic's default extra-field behavior applies to the native schema, so an unknown `length_scale` is ignored rather than treated as a control.

Its OpenAI-compatible schema is explicit:

```python
model_config = ConfigDict(extra="forbid")
speed: float = Field(1.0, ge=0.25, le=4.0)
```

A `length_scale` field sent to that endpoint would be rejected as an extra field. The adapter does not send one there.

The OpenAI endpoint maps the full `0.25–4.0` range to native model speed plus tempo:

```python
def openai_speed_controls(speed: float) -> tuple[float, float]:
    model_speed = min(2.0, max(0.5, speed))
    return model_speed, speed / model_speed
```

This fork does not expose `/synthesize`; it exposes `/tts/generate`, `/tts/stream`, and `/v1/audio/speech`.

### Clore.ai `/synthesize` deployment

Source: <https://docs.clore.ai/guides/audio-and-voice/melotts>

Its documented request is:

```json
{
  "text": "Hello from Clore.ai!",
  "language": "EN",
  "speaker_id": "EN-Default",
  "speed": 1.0,
  "format": "wav"
}
```

The documented schema uses `speed`, not `length_scale`. The service is started with `python -m melo.api_server`, which is not part of the official repository's documented public API, so its exact validation model is not available here. Its examples use `speed=1.0` and do not establish a hard range.

### Other wrappers

- `timhagel/MeloTTS-Docker-API-Server` exposes `/convert/tts` and accepts only `text`, `speed`, `language`, and `speaker_id`; it passes `speed=speed` to `tts_to_file`.
- The OpenAI-compatible PR #280 defines `speed: float = 0.9` and passes it directly to `tts_to_file`.
- These wrappers also expose no `length_scale` control.

## 5. Answers to the Specific Questions

### Which control governs synthesis?

`speed`. Official MeloTTS uses `speed` as the public multiplier and derives `length_scale = 1 / speed` internally.

### What happens when both are sent and disagree?

There is no known MeloTTS server that treats both as independent controls.

- If the server ignores unknown fields (the likely FastAPI/Pydantic default), `speed` governs and `length_scale` has no effect.
- If the server forbids extra fields, the request is rejected, usually with HTTP `422`.
- A custom server that accepted both and applied both would receive contradictory or duplicate rate controls. The adapter's values are inverse-consistent, but applying both would still be redundant and could double-adjust if the server uses them separately.

The adapter should not rely on this ambiguity. It should send `speed` only.

### What is the real accepted range?

There is no single universal MeloTTS HTTP range. The known implementations differ:

| Implementation | Range |
|---|---|
| Official Python/UI | No hard API check; UI `0.1–10.0`, default `1.0` |
| Hangry native schema | `0.5–2.0`, default `1.0` |
| Hangry OpenAI-compatible schema | `0.25–4.0`, default `1.0` |
| Current adapter validation | `0.5–2.0` |
| Repository MeloTTS profile | `0.8–1.2` |

The repository's `0.8–1.2` bounds are not established by the official implementation, the common wrapper schema, or the documented endpoint. `git blame` attributes them to the initial audio-foundation commit (`d72d5cd`); the source tree provides no citation or explanation for the narrower range. They appear to be a conservative product guard rather than a server-enforced limit.

## 6. Recommendation

### 6.1 Send `speed` only on the legacy endpoint

Remove `length_scale` from the `/synthesize` payload. Keep `speed` as the single rate control. This matches the official API, the Clore documented schema, the Hangry native schema, and every other inspected wrapper.

This is the narrow, evidence-based correction. It avoids the unnecessary extra field and removes the only place where a deployed server could see conflicting rate controls.

### 6.2 Keep the OpenAI-compatible payload unchanged

It already sends `speed` only. Do not add `length_scale` there.

### 6.3 Keep the global setting at `0.5–2.0`

Do not narrow Piper or Kokoro. Their profiles and adapters support `0.5–2.0`.

### 6.4 Keep MeloTTS at `0.8–1.2` only as a product guard, not as a server fact

The investigation supports the user's conservative product policy, but the rationale must be stated truthfully:

- `0.8–1.2` is the existing repository profile.
- It is not an upstream MeloTTS limit.
- No live deployment is available to establish whether `0.5–0.8` or `1.2–2.0` is safe for the specific server in use.
- Keeping the narrower range avoids claiming unverified behavior while the deployment is unavailable.

If a live deployment is later verified against the native `0.5–2.0` schema, the profile can be widened with evidence rather than guesswork.

### 6.5 Make application engine-aware

For Group 2A:

- Piper and Kokoro: accept `0.5–2.0`.
- MeloTTS: accept only `0.8–1.2` until deployment verification.
- `auto`: do not claim speed applied unless the selected engine is known and the value is within that engine's effective range. If the engine is unresolved, report the setting as not applied/unavailable rather than applying it optimistically.
- Out-of-range or unsupported values: report honestly; never silently clamp.
- A failed replacement must preserve the previous effective speed and status.

### 6.6 Add request-level tests

Before implementation is accepted, capture the exact outbound payloads:

- `/synthesize` contains `speed` and does not contain `length_scale`.
- `/v1/audio/speech` contains `speed` only.
- Non-default `1.25` reaches the server payload unchanged.
- Default `1.0` remains neutral.
- `0.8` and `1.2` are accepted for MeloTTS; values outside that range are rejected or marked unavailable according to the selected policy.
- A server returning `422` produces a truthful failure without fallback or silent clamping.

## 7. What a Live Verification Would Still Need

If a MeloTTS deployment is later started, the remaining questions can be settled with read-only/controlled checks:

1. `GET /tts/status` or `GET /tts/openapi.json` to identify the actual server and schema.
2. `POST /synthesize` at `0.8`, `1.0`, and `1.2` and compare returned audio duration.
3. Repeat at `0.5` and `2.0` only if the server schema permits them.
4. Send an extra `length_scale` field and inspect whether it is ignored, rejected, or changes the result.
5. Confirm whether the service exposes one control or a tempo/length post-processing layer.

Do not run these against a production service with user data.

## 8. Confidence and Limits

| Finding | Confidence |
|---|---|
| Official MeloTTS public control is `speed` | High — upstream source |
| Official `length_scale` is derived internally as `1/speed` | High — upstream source |
| Common native schemas do not expose `length_scale` | High — inspected schemas |
| `speed` governs when both are sent to a normal wrapper | Medium-high — no known server accepts both; Pydantic behavior is schema-dependent |
| The deployed server's exact range | Unverified — no live service |
| `0.8–1.2` is a server limit | Rejected — it is a repository product bound with no cited source |
| Audible speed behavior | Unverified — no synthesized acoustic comparison |

## 9. Conclusion

Group 2A should proceed with a single MeloTTS control: `speed`. The profile can retain `0.8–1.2` as a conservative product guard, but the implementation must not describe that range as a MeloTTS API limit or claim it has been verified against the live deployment. The cache already includes speed and needs no redesign.
