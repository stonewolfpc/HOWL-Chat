"""Local-endpoint adapter: KoboldCpp / SillyTavern / OpenAI-compatible servers.

These are *runtimes*, not hosted providers: they serve a model running on this
machine, so there is no per-call charge and no credit balance to verify. They
are the preferred fallback when the free hosted allowance is exhausted, and
they reintroduce the hardware constraints that S3 exists to measure.

Same `GenerateRequest` contract as the hosted adapter, so the memory tests do
not change when the route does. Loopback only: a non-loopback host would be a
remote provider wearing a local name, and is refused.
"""

from __future__ import annotations

import json
import time
import urllib.error
import urllib.request
from dataclasses import dataclass
from urllib.parse import urlparse

from .adapter import (
    UNAVAILABLE,
    AdapterConfigError,
    ExecutionNotAuthorised,
    GenerateRequest,
    GenerateResponse,
)
from .local_surfaces import SurfaceError, transcribe_any

LOOPBACK_HOSTS = frozenset({"127.0.0.1", "localhost", "::1", "[::1]"})

LOCAL_PROVIDERS = {
    "local:koboldcpp": "KoboldCpp OpenAI-compatible server",
    "local:sillytavern": "SillyTavern (proxying a local model)",
    "local:llamacpp": "llama.cpp server",
    "local:ollama": "Ollama",
}


class LocalEndpointError(AdapterConfigError):
    """Raised when a local endpoint is malformed or not actually local."""


@dataclass(frozen=True)
class LocalEndpoint:
    base_url: str
    chat_path: str = "/v1/chat/completions"

    def __post_init__(self) -> None:
        parsed = urlparse(self.base_url)
        if parsed.scheme not in ("http", "https"):
            raise LocalEndpointError(f"base_url scheme must be http(s), got {parsed.scheme!r}")
        if (parsed.hostname or "").lower() not in LOOPBACK_HOSTS:
            raise LocalEndpointError(
                f"host {parsed.hostname!r} is not loopback: a remote host is a hosted provider "
                "and must go through the budget policy"
            )
        if not self.chat_path.startswith("/"):
            raise LocalEndpointError("chat_path must begin with /")

    @property
    def is_loopback(self) -> bool:
        return (urlparse(self.base_url).hostname or "").lower() in LOOPBACK_HOSTS

    def url(self) -> str:
        return self.base_url.rstrip("/") + self.chat_path


class LocalHTTPAdapter:
    """Loopback-only, no per-call charge, still gated: local execution is execution."""

    capabilities = {
        "seed": False,
        "logprobs": False,
        "stop_sequences": True,
        "streaming": False,
        "max_context": None,
        "network": True,
        "loopback_only": True,
        "paid": False,
    }

    def __init__(self, grant: object = None, endpoint: LocalEndpoint | None = None, transport=None,
                 post_transport=None, timeout: float = 120.0) -> None:
        if grant is None or not getattr(grant, "grants_execution", False):
            raise ExecutionNotAuthorised(
                "local inference is still execution: P1 forbids it, and P2 needs a grant"
            )
        if endpoint is None:
            raise LocalEndpointError("a local adapter needs an explicit loopback endpoint")
        if not endpoint.is_loopback:
            raise LocalEndpointError("endpoint is not loopback")
        self._endpoint = endpoint
        self._grant = grant
        self._transport = transport
        self._post_transport = post_transport or default_post_transport
        self._timeout = timeout

    def payload(self, request: GenerateRequest) -> dict:
        if not isinstance(request, GenerateRequest):
            raise LocalEndpointError("local adapter accepts GenerateRequest only")
        return {
            "model": request.model_id,
            "messages": [dict(m) for m in request.messages],
            "stream": False,
            **request.params,
        }

    def generate(self, request: GenerateRequest) -> GenerateResponse:
        """One attempt. No retry, no repair, no parameter pruning.

        A failed call is a recorded result, not a reason to try again: the
        ledger already counted it (row 482 §6).
        """
        if not getattr(self._grant, "grants_execution", False):
            raise ExecutionNotAuthorised("no grant")
        payload = self.payload(request)
        started = time.monotonic()
        status, raw_text = self._post_transport(self._endpoint.url(), payload, self._timeout)
        latency_ms = int(round((time.monotonic() - started) * 1000))

        try:
            body = json.loads(raw_text)
        except Exception:
            # Keep the bytes: an unparsable body is still the evidence.
            return GenerateResponse(
                text=UNAVAILABLE,
                raw_body={"_unparsed": raw_text, "_http_status": status},
                request_echo=payload,
                provider_requested=request.provider,
                provider_used=self._endpoint.base_url,
                model_requested=request.model_id,
                latency_ms=latency_ms,
                http_status=status,
                surface="unparsed",
                error=f"non-JSON body (http {status})",
            )

        try:
            surface, response = transcribe_any(body)
        except SurfaceError as exc:
            return GenerateResponse(
                text=UNAVAILABLE,
                raw_body=body,
                request_echo=payload,
                provider_requested=request.provider,
                provider_used=self._endpoint.base_url,
                model_requested=request.model_id,
                latency_ms=latency_ms,
                http_status=status,
                surface="unknown",
                error=f"unrecognised body: {exc}",
            )

        response.latency_ms = latency_ms
        response.http_status = status
        response.surface = surface
        response.request_echo = payload
        response.provider_requested = request.provider
        response.provider_used = self._endpoint.base_url
        response.model_requested = request.model_id
        if status != 200 and response.error is None:
            response.error = f"http {status}"
        return response


def default_post_transport(url: str, payload: dict, timeout: float) -> tuple[int, str]:
    """Single POST, no retries, no redirects followed silently."""
    data = json.dumps(payload).encode("utf-8")
    request = urllib.request.Request(
        url, data=data, headers={"Content-Type": "application/json"}, method="POST"
    )
    try:
        with urllib.request.urlopen(request, timeout=timeout) as response:
            return int(response.status), response.read().decode("utf-8", "replace")
    except urllib.error.HTTPError as exc:
        return int(exc.code), exc.read().decode("utf-8", "replace")
