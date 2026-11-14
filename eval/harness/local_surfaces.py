"""Local reply surfaces: KoboldCpp native and OpenAI-compatible chat.

Two surfaces, one response type. This module maps *documented* shapes only —
nothing here has been checked against a running server, because none was
contacted (see the local adapter plan, L1).

The rule the hosted route taught us: **the parser follows the provider.** If a
real body disagrees with a fixture at L3, the fixture is corrected and this
module follows it — never the reverse.
"""

from __future__ import annotations

from .adapter import (
    UNAVAILABLE,
    AdapterConfigError,
    GenerateRequest,
    GenerateResponse,
    transcribe,
)

SURFACE_CHAT = "chat_completions"
SURFACE_NATIVE = "kobold_native"
SURFACE_UNKNOWN = "unknown"

# Native surface names the token limit differently and uses its own sampler keys.
NATIVE_PARAM_MAP = {
    "max_tokens": "max_length",
    "top_p": "top_p",
    "top_k": "top_k",
    "temperature": "temperature",
    "presence_penalty": "rep_pen",
    "frequency_penalty": "rep_pen_range",
    "min_p": "min_p",
    "stop": "stop_sequence",
}

# Flattening messages into a prompt is a *semantic* change, not a formatting one.
FLATTEN_TEMPLATE = "system\\n\\nuser"


class SurfaceError(AdapterConfigError):
    """Raised when a body matches neither documented shape."""


def flatten_prompt(request: GenerateRequest) -> str:
    """Join messages for the native surface.

    Recorded as a template name so a cross-surface comparison can see that the
    prompt text differed. It is a different condition, not a re-run.
    """
    parts: list[str] = []
    for message in request.messages:
        role = message["role"]
        if role == "system":
            parts.insert(0, message["content"])
        else:
            parts.append(message["content"])
    return "\n\n".join(p for p in parts if p)


def build_chat_payload(request: GenerateRequest) -> dict:
    """OpenAI-compatible request body."""
    if not isinstance(request, GenerateRequest):
        raise SurfaceError("chat payload requires a GenerateRequest")
    return {
        "model": request.model_id,
        "messages": [dict(m) for m in request.messages],
        "stream": False,
        **request.params,
    }


def build_native_payload(request: GenerateRequest, stop_sequences: list[str] | None = None) -> dict:
    """KoboldCpp native request body, with recorded parameter translation."""
    if not isinstance(request, GenerateRequest):
        raise SurfaceError("native payload requires a GenerateRequest")
    payload: dict = {"prompt": flatten_prompt(request), "stream": False}
    translated: dict = {}
    for key, value in sorted(request.params.items()):
        target = NATIVE_PARAM_MAP.get(key)
        if target is None:
            raise SurfaceError(f"no native equivalent for parameter {key!r}; refusing to drop it silently")
        payload[target] = value
        translated[key] = target
    if stop_sequences:
        payload["stop_sequence"] = list(stop_sequences)
    payload["sampler_translation"] = translated
    return payload


def detect_surface(body: object) -> str:
    """Which documented shape is this?"""
    if not isinstance(body, dict):
        raise SurfaceError(f"body must be a mapping, got {type(body).__name__}")
    if "choices" in body:
        return SURFACE_CHAT
    if "results" in body:
        return SURFACE_NATIVE
    raise SurfaceError(
        "body matches neither surface: expected 'choices' (chat) or 'results' (native)"
    )


def _error_response(body: dict) -> GenerateResponse:
    detail = body.get("error") or body.get("detail") or body.get("message") or "unknown error"
    return GenerateResponse(
        text=UNAVAILABLE,
        raw_body=body,
        provider_used=body.get("provider") or UNAVAILABLE,
        model_reported=body.get("model") or UNAVAILABLE,
        system_fingerprint=UNAVAILABLE,
        usage=UNAVAILABLE,
        timestamp=str(body.get("created") or UNAVAILABLE),
        request_id=body.get("id") or UNAVAILABLE,
        error=f"server error: {detail}",
    )


def transcribe_native(body: dict) -> GenerateResponse:
    """KoboldCpp native reply -> one response record."""
    results = body.get("results")
    if not isinstance(results, list) or not results:
        return GenerateResponse(
            text=UNAVAILABLE,
            raw_body=body,
            provider_used=body.get("provider") or UNAVAILABLE,
            model_reported=body.get("model") or UNAVAILABLE,
            system_fingerprint=UNAVAILABLE,
            usage=UNAVAILABLE,
            timestamp=str(body.get("created") or UNAVAILABLE),
            request_id=body.get("id") or UNAVAILABLE,
            error="empty results: server returned no completion",
        )
    first = results[0] if isinstance(results[0], dict) else {}
    usage = UNAVAILABLE
    if first.get("tokens_predicted") is not None or first.get("tokens_evaluated") is not None:
        usage = {
            "completion_tokens": first.get("tokens_predicted"),
            "prompt_tokens": first.get("tokens_evaluated"),
        }
    return GenerateResponse(
        text=first.get("text") if first.get("text") is not None else UNAVAILABLE,
        raw_body=body,
        provider_used=body.get("provider") or UNAVAILABLE,
        model_reported=first.get("model") or body.get("model") or UNAVAILABLE,
        system_fingerprint=UNAVAILABLE,
        usage=usage,
        timestamp=str(body.get("created") or UNAVAILABLE),
        request_id=body.get("id") or UNAVAILABLE,
        finish_reason=first.get("finish_reason") or UNAVAILABLE,
    )


def transcribe_any(body: object) -> tuple[str, GenerateResponse]:
    """Detect, then transcribe. Error bodies return a record, they do not raise.

    Error bodies are checked **first**: a 500 from a local server typically
    carries `error`/`detail` and no `choices`/`results`, so detecting the shape
    before handling the error would raise on the most common failure there is.
    """
    if not isinstance(body, dict):
        raise SurfaceError(f"body must be a mapping, got {type(body).__name__}")
    if body.get("error") or body.get("detail"):
        try:
            surface = detect_surface(body)
        except SurfaceError:
            surface = SURFACE_UNKNOWN
        return surface, _error_response(body)
    surface = detect_surface(body)
    if surface == SURFACE_CHAT:
        return surface, transcribe(body)
    return surface, transcribe_native(body)


def all_surfaces() -> tuple[str, str]:
    return (SURFACE_CHAT, SURFACE_NATIVE)
