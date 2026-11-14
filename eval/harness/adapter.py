"""Model adapter interface (row 477 v1.1), offline.

P1 rule: this module can *describe* and *validate* a hosted request, and can
refuse one. It cannot perform one. The network-capable adapter requires an
`ExecutionGrant` that no file in this repository can produce (authorization.py
looks for eval/runs/P2-GRANT.json, which does not exist and must not be
created during P1).

`huggingface_hub` is imported *inside a method*, never at module scope, so the
package imports and validates with no network-capable client present.
"""

from __future__ import annotations

import hashlib
import inspect
import json
from dataclasses import dataclass, field
from typing import Any, Protocol

HARNESS_ADAPTER_NAME = "memory09-harness"
HARNESS_ADAPTER_VERSION = "0.1.0-p1"

ROUTING_SUFFIXES = (":fastest", ":cheapest", ":preferred")
PARAM_KEYS = frozenset(
    {
        "temperature",
        "top_p",
        "top_k",
        "min_p",
        "presence_penalty",
        "frequency_penalty",
        "max_tokens",
        "stop",
    }
)

# Provenance fields the adapter must attempt to fill (row 477 v1.1 §8).
PROVENANCE_FIELDS = (
    "adapter_name",
    "adapter_version",
    "client_version",
    "model_requested",
    "model_revision",
    "provider_requested",
    "provider_used",
    "returned_model",
    "system_fingerprint",
    "provider_headers",
    "params_sent",
    "seed",
    "seed_honoured",
    "prompt_hash",
    "response_hash",
    "state_hash_before",
    "state_hash_after",
    "fixture_id",
    "fixture_version",
    "test_id",
    "timestamp",
    "raw_body",
    "text",
    "request_echo",
    "error",
)

UNAVAILABLE = "unavailable"


class ExecutionNotAuthorised(RuntimeError):
    """Raised when anything would reach a provider during P1."""


class AdapterConfigError(ValueError):
    """Raised when a request violates the interface rules (T-1, T-2, D-1)."""


def canonical(obj: Any) -> str:
    """Stable JSON used for every hash in the harness."""
    return json.dumps(obj, sort_keys=True, separators=(",", ":"), ensure_ascii=False)


def sha256_text(text: str) -> str:
    return hashlib.sha256(text.encode("utf-8")).hexdigest()


@dataclass(frozen=True)
class GenerateRequest:
    """One request, fully determined. Frozen so params cannot drift mid-run."""

    messages: tuple[dict, ...]
    model_id: str
    provider: str
    params: dict
    seed: int | None = None
    stream: bool = False

    def __post_init__(self) -> None:
        self.validate()

    def validate(self) -> None:
        # T-1: provider must be explicit.
        if not isinstance(self.provider, str) or not self.provider.strip():
            raise AdapterConfigError("provider must be named (T-1)")
        if self.provider.strip().lower() == "auto":
            raise AdapterConfigError("provider='auto' is prohibited for evaluation (T-1)")
        # T-2: model id unmodified; no router routing suffixes.
        if not isinstance(self.model_id, str) or not self.model_id.strip():
            raise AdapterConfigError("model_id must be named")
        for suffix in ROUTING_SUFFIXES:
            if suffix in self.model_id:
                raise AdapterConfigError(
                    f"model_id contains routing policy {suffix!r}; the routing interface is not used (T-2)"
                )
        if ":" in self.model_id:
            raise AdapterConfigError(
                "model_id must be passed unmodified; provider belongs in the provider argument (T-2)"
            )
        # Streaming complicates exact capture.
        if self.stream:
            raise AdapterConfigError("stream=True is prohibited for evaluation capture")
        if not self.messages:
            raise AdapterConfigError("messages must not be empty")
        for m in self.messages:
            if set(m) != {"role", "content"}:
                raise AdapterConfigError("each message must have exactly role and content")
            if not isinstance(m["content"], str):
                raise AdapterConfigError("message content must be text")
        unknown = set(self.params) - PARAM_KEYS
        if unknown:
            raise AdapterConfigError(f"unknown sampling parameters: {sorted(unknown)}")
        if self.seed is not None and not isinstance(self.seed, int):
            raise AdapterConfigError("seed must be an int or None")

    def echo(self) -> dict:
        """Exactly what was sent, in serialisable form."""
        return {
            "messages": [dict(m) for m in self.messages],
            "model": self.model_id,
            "provider": self.provider,
            "params": dict(self.params),
            "seed": self.seed,
            "stream": self.stream,
        }

    def prompt_hash(self) -> str:
        return sha256_text(canonical({"messages": [dict(m) for m in self.messages]}))


@dataclass
class GenerateResponse:
    """The record, not the text. `raw_body` is the evidence."""

    text: Any = None
    raw_body: Any = None
    request_echo: Any = None
    provider_requested: Any = None
    provider_used: Any = None
    model_requested: Any = None
    model_reported: Any = None
    system_fingerprint: Any = None
    usage: Any = None
    latency_ms: Any = None
    timestamp: Any = None
    request_id: Any = None
    client_version: Any = None
    finish_reason: Any = None
    error: Any = None
    http_status: Any = None
    surface: Any = None

    @property
    def incomplete(self) -> bool:
        """Truncated at the token limit: a truncated reply is not a model failure."""
        return self.finish_reason == "length"

    def check_admissible(self) -> tuple[bool, list[str]]:
        problems: list[str] = []
        if self.raw_body is None:
            problems.append("raw_body missing: text alone is not evidence")
        if self.raw_body == UNAVAILABLE:
            problems.append("raw_body marked unavailable")
        if self.provider_used is None:
            problems.append("provider_used not recorded")
        return (not problems, problems)


def transcribe(body: object) -> GenerateResponse:
    """Turn a provider body into a response record. Pure: no network, no client.

    Testable against a synthetic body, which is why P1 can close this gap
    without making a call. `raw_body` is preserved verbatim.
    """
    if not isinstance(body, dict):
        raise AdapterConfigError(f"provider body must be a mapping, got {type(body).__name__}")
    choices = body.get("choices")
    text: Any = UNAVAILABLE
    if isinstance(choices, list) and choices:
        first = choices[0] if isinstance(choices[0], dict) else {}
        message = first.get("message") if isinstance(first.get("message"), dict) else {}
        content = message.get("content")
        if content is not None:
            text = content
    created = body.get("created")
    return GenerateResponse(
        text=text,
        raw_body=body,
        provider_used=body.get("provider") or UNAVAILABLE,
        model_reported=body.get("model") or UNAVAILABLE,
        system_fingerprint=body.get("system_fingerprint") or UNAVAILABLE,
        usage=body.get("usage") if body.get("usage") is not None else UNAVAILABLE,
        timestamp=str(created) if created is not None else UNAVAILABLE,
        request_id=body.get("id") or UNAVAILABLE,
        finish_reason=(choices[0].get("finish_reason", UNAVAILABLE) if isinstance(choices, list) and choices and isinstance(choices[0], dict) else UNAVAILABLE),
    )


def client_chat_parameters() -> frozenset[str] | None:
    """Parameters the installed client's chat_completion actually accepts.

    Offline: signature inspection only. `None` means the client is not installed,
    which is reported rather than assumed.
    """
    try:
        from huggingface_hub import InferenceClient
    except Exception:
        return None
    try:
        return frozenset(inspect.signature(InferenceClient.chat_completion).parameters)
    except (TypeError, ValueError):
        return None


def check_params_supported(params: dict) -> list[str]:
    """Return params the client would reject. Empty list means all supported.

    A param the harness validates but the client cannot accept is a runtime
    TypeError waiting for P2, so it is caught offline instead.
    """
    accepted = client_chat_parameters()
    if accepted is None:
        return sorted(params)  # nothing can be shown to be supported
    return sorted(set(params) - set(accepted))


class ModelAdapter(Protocol):
    capabilities: dict

    def generate(self, request: GenerateRequest) -> GenerateResponse: ...


class OfflineAdapter:
    """The only adapter P1 can build. Capable of nothing but refusal."""

    capabilities = {
        "seed": False,
        "logprobs": False,
        "stop_sequences": False,
        "streaming": False,
        "max_context": None,
        "network": False,
    }

    def generate(self, request: GenerateRequest) -> GenerateResponse:
        raise ExecutionNotAuthorised(
            "P1 forbids generation. P2 execution requires a separate owner decision."
        )


def classify_route(provider: str) -> str:
    """Route label for a provider string: "local" or "hosted"."""
    from .budget import classify

    return classify(provider)


class HuggingFaceAdapter:
    """Transport: `huggingface_hub.InferenceClient` with explicit `provider=`.

    Written but **never constructed during P1**: the grant does not exist, so
    the constructor raises. It is included so P2 has an inspectable instrument
    rather than code written under time pressure after authorisation.
    """

    capabilities = {
        "seed": False,
        "logprobs": False,
        "stop_sequences": True,
        "streaming": True,
        "max_context": None,
        "network": True,
    }

    def __init__(self, grant: Any = None, client_factory: Any = None) -> None:
        if grant is None or not getattr(grant, "grants_execution", False):
            raise ExecutionNotAuthorised(
                "network-capable adapter construction requires a signed P2 grant"
            )
        # Belt and braces at the object level: even a valid grant does not make a
        # hosted route affordable. See budget.HOSTED_REFUSAL.
        from .budget import FREE_ALLOWANCE_EXISTS, HOSTED_REFUSAL, BudgetRefusal

        if not FREE_ALLOWANCE_EXISTS:
            raise BudgetRefusal(HOSTED_REFUSAL)
        self._grant = grant
        self._client_factory = client_factory
        self._client = None

    def _get_client(self, provider: str) -> Any:
        if self._client is None:
            from huggingface_hub import InferenceClient  # local, never module scope

            self._client = (
                self._client_factory(provider) if self._client_factory else InferenceClient(provider=provider)
            )
        return self._client

    def generate(self, request: GenerateRequest) -> GenerateResponse:
        if not getattr(self._grant, "grants_execution", False):
            raise ExecutionNotAuthorised("no grant")
        unsupported = check_params_supported(request.params)
        if unsupported:
            raise AdapterConfigError(
                f"the client does not accept these params: {unsupported} (fail before the wire)"
            )
        client = self._get_client(request.provider)
        body = client.chat_completion(
            model=request.model_id,
            messages=[dict(m) for m in request.messages],
            stream=False,
            **request.params,
        )
        return transcribe(body)


def build_adapter(provider: str, grant: Any = None, local_endpoint: Any = None):
    """Route-aware factory. The budget policy runs *before* anything is constructed.

    A provider string of `local:<runtime>` selects the loopback adapter; anything
    else is hosted-free-only, which the policy allows because no payment method
    exists. Nothing here falls back to another route (row 477 F-1).
    """
    from .budget import check_provider

    route = check_provider(provider)
    if route == "local":
        from .local_adapter import LocalEndpoint, LocalHTTPAdapter

        if local_endpoint is None:
            raise AdapterConfigError(
                "a local route needs an explicit loopback endpoint; no default is guessed"
            )
        endpoint = (
            local_endpoint
            if isinstance(local_endpoint, LocalEndpoint)
            else LocalEndpoint(base_url=str(local_endpoint))
        )
        return LocalHTTPAdapter(grant=grant, endpoint=endpoint)
    return HuggingFaceAdapter(grant=grant)
