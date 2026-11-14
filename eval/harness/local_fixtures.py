"""Synthetic local reply fixtures (plan L1).

Built from **documented** shapes, not captured bodies. Each exists to force one
behaviour in the transcription path. When a real body arrives at L3 and differs,
the fixture is corrected first — the parser follows the provider.

These are engineering fixtures: no corpus family, no score, never cited as a
MEMORY-09 unit.
"""

from __future__ import annotations

CHAT_OK = {
    "id": "chatcmpl-local-1",
    "object": "chat.completion",
    "created": 1700000000,
    "model": "local-model",
    "choices": [
        {
            "index": 0,
            "message": {"role": "assistant", "content": "Corvin Ashgrave of Ossary trained me."},
            "finish_reason": "stop",
        }
    ],
    "usage": {"prompt_tokens": 34, "completion_tokens": 9, "total_tokens": 43},
}

CHAT_TRUNCATED = {
    "id": "chatcmpl-local-2",
    "created": 1700000001,
    "model": "local-model",
    "choices": [
        {
            "index": 0,
            "message": {"role": "assistant", "content": "Corvin Ashgrave of Oss"},
            "finish_reason": "length",
        }
    ],
    "usage": {"prompt_tokens": 34, "completion_tokens": 256, "total_tokens": 290},
}

CHAT_EMPTY = {
    "id": "chatcmpl-local-3",
    "created": 1700000002,
    "model": "local-model",
    "choices": [],
}

CHAT_ERROR = {
    "error": "The model has not been loaded yet",
    "detail": {"code": "model_not_loaded"},
}

NATIVE_OK = {
    "results": [
        {
            "text": "Corvin Ashgrave of Ossary trained me.",
            "tokens_predicted": 9,
            "tokens_evaluated": 34,
            "finish_reason": "stop",
            "model": "local-model",
        }
    ],
    "model": "local-model",
}

NATIVE_TRUNCATED = {
    "results": [
        {
            "text": "Corvin Ashgrave of Oss",
            "tokens_predicted": 256,
            "tokens_evaluated": 34,
            "finish_reason": "length",
        }
    ]
}

NATIVE_EMPTY = {"results": []}

NATIVE_ERROR = {"detail": "context length exceeded", "code": "context_exceeded"}

NATIVE_UNKNOWN_FIELDS = {
    "results": [
        {
            "text": "Corvin Ashgrave of Ossary trained me.",
            "tokens_predicted": 9,
            "tokens_evaluated": 34,
            "finish_reason": "stop",
            "seed": 42,
            "sampler_order": [6, 0, 1, 3, 4, 2, 5],
            "some_future_field": {"nested": True},
        }
    ],
    "koboldcpp_version": "unknown-synthetic",
}

ALL_FIXTURES = {
    "chat_ok": CHAT_OK,
    "chat_truncated": CHAT_TRUNCATED,
    "chat_empty": CHAT_EMPTY,
    "chat_error": CHAT_ERROR,
    "native_ok": NATIVE_OK,
    "native_truncated": NATIVE_TRUNCATED,
    "native_empty": NATIVE_EMPTY,
    "native_error": NATIVE_ERROR,
    "native_unknown_fields": NATIVE_UNKNOWN_FIELDS,
}

# A body matching neither documented shape.
UNRECOGNISED = {"message": "hello", "payload": {"text": "x"}}
