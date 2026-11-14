"""L2 endpoint discovery: identity only, no generation.

Authorised 2026-11-13 ("you may run it"). Scope, fixed in code:

* loopback only — a remote host is refused, not discovered;
* **GET on identity paths only** — the allowed path list is the whole API;
* no prompt, no model load, no server start, no retry storm;
* missing fields are recorded as `unavailable`, never inferred.

If nothing is listening, that is the result. The harness does not start servers.
"""

from __future__ import annotations

from .adapter import UNAVAILABLE, sha256_text
from .local_adapter import LocalEndpoint, LocalEndpointError

DEFAULT_ENDPOINT = "http://127.0.0.1:5001"

# The only paths this module may request. Anything else is a bug.
IDENTITY_PATHS = (
    "/api/extra/version",
    "/api/v1/model",
    "/api/extra/true_max_context_length",
)

# Discovery authorisation date. Set to None to disable the probe entirely.
L2_AUTHORISED_ON = "2026-11-13"

TIMEOUT_SECONDS = 3.0


class DiscoveryRefused(RuntimeError):
    """Raised when discovery is not authorised or the endpoint is not loopback."""


class DiscoveryError(RuntimeError):
    """Raised on a malformed discovery attempt."""


def default_transport(url: str, timeout: float = TIMEOUT_SECONDS) -> tuple[int, str]:
    """One GET. Returns (status, body_text). Never raises, never retries."""
    import urllib.error
    import urllib.request

    try:
        with urllib.request.urlopen(urllib.request.Request(url, method="GET"), timeout=timeout) as resp:
            return int(resp.status), resp.read().decode("utf-8", "replace")
    except urllib.error.HTTPError as exc:
        try:
            body = exc.read().decode("utf-8", "replace")
        except Exception:
            body = ""
        return int(exc.code), body
    except urllib.error.URLError as exc:
        return 0, f"not listening: {exc.reason}"
    except Exception as exc:  # timeout, reset, anything else
        return 0, f"no response: {type(exc).__name__}: {exc}"


def _check_path(path: str) -> str:
    if path not in IDENTITY_PATHS:
        raise DiscoveryError(
            f"path {path!r} is not an identity path; discovery may only GET {list(IDENTITY_PATHS)}"
        )
    return path


def probe(endpoint_url: str = DEFAULT_ENDPOINT, transport=None, timeout: float = TIMEOUT_SECONDS) -> dict:
    """GET the identity paths once each. Returns raw results, no interpretation."""
    if L2_AUTHORISED_ON is None:
        raise DiscoveryRefused("L2 discovery is not authorised (L2_AUTHORISED_ON is None)")
    try:
        endpoint = LocalEndpoint(base_url=endpoint_url)
    except LocalEndpointError:
        raise
    if not endpoint.is_loopback:
        raise DiscoveryRefused(f"endpoint {endpoint_url!r} is not loopback")

    send = transport or default_transport
    results: dict = {}
    for path in IDENTITY_PATHS:
        _check_path(path)
        url = endpoint.base_url.rstrip("/") + path
        status, body = send(url, timeout)
        results[path] = {"url": url, "status": int(status), "body": body}
    return results


def _parse_json(body: str):
    import json

    try:
        return json.loads(body)
    except Exception:
        return None


def identity_block(results: dict) -> dict:
    """Turn raw probe results into an identity record with explicit gaps."""
    version_raw = results.get("/api/extra/version", {})
    model_raw = results.get("/api/v1/model", {})
    ctx_raw = results.get("/api/extra/true_max_context_length", {})

    listening = any(int(r.get("status", 0)) > 0 for r in results.values())

    version_body = _parse_json(version_raw.get("body", ""))
    model_body = _parse_json(model_raw.get("body", ""))
    ctx_body = _parse_json(ctx_raw.get("body", ""))

    # Real /api/extra/version: {"result": "KoboldCpp", "version": "1.120", ...}
    # The fixture-era assumption that `result` *was* the version was wrong; the
    # parser follows the provider (row 487 §5), so `version` is preferred and a
    # non-numeric `result` is treated as the server name instead.
    server_name = None
    version = None
    if isinstance(version_body, dict):
        raw_result = version_body.get("result")
        if isinstance(raw_result, str) and raw_result[:1].isdigit():
            version = raw_result
        else:
            server_name = raw_result
        if version_body.get("version") is not None:
            version = version_body.get("version")
            server_name = raw_result if isinstance(raw_result, str) else server_name
    model = model_body.get("result") if isinstance(model_body, dict) else None
    context = None
    if isinstance(ctx_body, dict):
        context = ctx_body.get("value", ctx_body.get("result"))
    elif isinstance(ctx_body, int):
        context = ctx_body

    block = {
        "listening": listening,
        "endpoint": version_raw.get("url", "").rsplit("/api/", 1)[0] or UNAVAILABLE,
        "server_name": server_name if server_name else UNAVAILABLE,
        "server_version": version if version else UNAVAILABLE,
        "loaded_model": model if model else UNAVAILABLE,
        "max_context_length": context if context is not None else UNAVAILABLE,
        "probe": {path: {"status": r.get("status"), "body": r.get("body")} for path, r in results.items()},
    }
    if not listening:
        block["conclusion"] = (
            "no local inference server answered. The harness does not start servers; "
            "starting one is a separate decision."
        )
    return block


def discover(endpoint_url: str = DEFAULT_ENDPOINT, transport=None, timeout: float = TIMEOUT_SECONDS) -> dict:
    results = probe(endpoint_url, transport=transport, timeout=timeout)
    block = identity_block(results)
    block["block_sha256"] = sha256_text(str(sorted(block.items())))
    return block


def allowed_paths() -> tuple[str, ...]:
    return IDENTITY_PATHS
