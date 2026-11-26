"""Local validation and enforcement checks (P1).

Runs entirely offline under a socket guard: any attempt to open a socket
during validation fails the run rather than succeeding quietly. Exits non-zero
if any check fails.
"""

from __future__ import annotations

import contextlib
import dataclasses
import io
import json
import shutil
import socket
import sys
import tempfile
import types
from pathlib import Path

from . import HARNESS_VERSION
from . import adapter as adapter_module
from . import authorization, budget, cli as cli_module, conditions, cost, discover as discovery, dispatch, fixtures, local_adapter, local_fixtures, local_surfaces, plan, provenance
from .adapter import (
    UNAVAILABLE,
    AdapterConfigError,
    ExecutionNotAuthorised,
    GenerateRequest,
    HuggingFaceAdapter,
    OfflineAdapter,
)
from .conditions import SeparationError
from .fixtures import FixtureError
from .state import State, StateError, diff
from .plan import PlanIncomplete, Slot

HARNESS_DIR = Path(__file__).resolve().parent

# Prefixes the suite creates. Anything matching one of these is disposable
# scratch, owned by the run that created it. Nothing else is ever deleted.
SCRATCH_PREFIXES = ("_selftest-", "_grant-")


def find_scratch(root: Path | None = None) -> list[Path]:
    """All scratch directories present under `root` (default: this directory)."""
    root = Path(root) if root is not None else HARNESS_DIR
    if not root.is_dir():
        raise AssertionError(f"scratch scan root is not a directory: {root}")
    found: list[Path] = []
    for prefix in SCRATCH_PREFIXES:
        found.extend(sorted(p for p in root.glob(prefix + "*") if p.is_dir()))
    return found


@contextlib.contextmanager
def scratch_dir(prefix: str):
    """Create a scratch directory we own, and remove exactly that one.

    Cleanup runs on success, on exception, and on early return. The path is held
    by the context, so no glob is used to decide what to delete — a directory we
    did not create cannot be removed by this code.
    """
    if not prefix.startswith("_"):
        raise AssertionError("scratch prefixes must begin with _ so they are discoverable")
    target = Path(tempfile.mkdtemp(prefix=prefix, dir=HARNESS_DIR))
    try:
        yield target
    finally:
        shutil.rmtree(target, ignore_errors=True)


CHECKS: list[tuple[str, callable]] = []
RESULTS: list[tuple[str, bool, str]] = []


def check(name: str):
    def register(fn):
        CHECKS.append((name, fn))
        return fn

    return register


# ---------------------------------------------------------------- adapter
@check("adapter: no network client imported at module scope")
def _c1():
    if "huggingface_hub" in sys.modules:
        raise AssertionError("huggingface_hub was imported at module scope")


@check("adapter: provider='auto' rejected (T-1)")
def _c2():
    try:
        GenerateRequest(messages=({"role": "user", "content": "x"},), model_id="m", provider="auto", params={})
    except AdapterConfigError:
        return
    raise AssertionError("provider='auto' was accepted")


@check("adapter: routing policies rejected (T-2)")
def _c3():
    for suffix in (":fastest", ":cheapest", ":preferred"):
        try:
            GenerateRequest(
                messages=({"role": "user", "content": "x"},),
                model_id="org/model" + suffix,
                provider="groq",
                params={},
            )
        except AdapterConfigError:
            continue
        raise AssertionError(f"routing policy {suffix} accepted")


@check("adapter: provider suffix in model id rejected (T-2)")
def _c4():
    try:
        GenerateRequest(
            messages=({"role": "user", "content": "x"},), model_id="org/model:groq", provider="groq", params={}
        )
    except AdapterConfigError:
        return
    raise AssertionError("model id with provider suffix accepted")


@check("adapter: streaming rejected")
def _c5():
    try:
        GenerateRequest(
            messages=({"role": "user", "content": "x"},), model_id="m", provider="groq", params={}, stream=True
        )
    except AdapterConfigError:
        return
    raise AssertionError("streaming accepted")


@check("adapter: unknown sampling parameter rejected")
def _c6():
    try:
        GenerateRequest(
            messages=({"role": "user", "content": "x"},), model_id="m", provider="groq", params={"wibble": 1}
        )
    except AdapterConfigError:
        return
    raise AssertionError("unknown parameter accepted")


@check("adapter: offline adapter refuses to generate")
def _c7():
    req = GenerateRequest(messages=({"role": "user", "content": "x"},), model_id="m", provider="groq", params={})
    try:
        OfflineAdapter().generate(req)
    except ExecutionNotAuthorised:
        return
    raise AssertionError("offline adapter produced a response")


@check("adapter: network adapter cannot be constructed without a grant")
def _c8():
    try:
        HuggingFaceAdapter(grant=None)
    except ExecutionNotAuthorised:
        return
    raise AssertionError("HuggingFaceAdapter constructed without a grant")


# ---------------------------------------------------------- authorisation
@check("authorisation: no P2 grant exists")
def _c9():
    if authorization.grant_status() != "ABSENT":
        raise AssertionError("a P2 grant file exists; P1 must not have one")
    try:
        authorization.require_grant()
    except ExecutionNotAuthorised:
        return
    raise AssertionError("require_grant() succeeded without a grant file")


# ------------------------------------------------------------------ state
@check("state: hash is deterministic and snapshot immutable")
def _c10():
    a = fixtures.build_f_st_01_s1().state
    b = fixtures.build_f_st_01_s1().state
    if a.hash() != b.hash():
        raise AssertionError("identical states hashed differently")
    snap = a.snapshot()
    a.put("item", "item-001", {"id": "item-001", "name": "lantern"})
    if snap.sha256 == a.hash():
        raise AssertionError("snapshot changed after a later mutation")


@check("state: diff detects an added entity")
def _c11():
    st = fixtures.build_f_st_01_s1().state
    before = st.snapshot()
    st.put("item", "item-001", {"id": "item-001", "name": "lantern"})
    changes = diff(before, st.snapshot())
    if not any(c["kind"] == "item" and c["change"] == "added" for c in changes):
        raise AssertionError(f"diff missed the addition: {changes}")


@check("state: deletion records removal, never erases")
def _c12():
    st = fixtures.build_f_st_01_s1().state
    before = st.snapshot()
    st.put("item", "item-002", {"id": "item-002", "name": "waterskin"})
    mid = st.snapshot()
    st.delete("item", "item-002", turn=12)
    changes = diff(mid, st.snapshot())
    rec = st.get("item", "item-002")
    if rec.get("state") != "removed" or rec.get("removed_at_turn") != 12:
        raise AssertionError(f"removal not recorded: {rec}")
    if not any(c["change"] == "removed" for c in changes):
        raise AssertionError("diff missed the removal")


@check("state: invalid fact source rejected")
def _c13():
    st = State()
    try:
        st.put("fact", "f", {"id": "f", "text": "t", "source": "guessed"})
    except StateError:
        return
    raise AssertionError("invalid fact source accepted")


@check("state: absent fact asserts as absent")
def _c14():
    st = fixtures.build_f_st_01_s1().state
    present, source, turn = st.assert_fact("fact-999")
    if present or source is not None or turn is not None:
        raise AssertionError("absent fact reported as present")


# --------------------------------------------------------------- fixture
@check("fixture: JSON round-trip preserves both hashes")
def _c15():
    fx = fixtures.build_f_st_01_s1()
    rt = fixtures.from_dict(fx.to_dict())
    if rt.sha256() != fx.sha256():
        raise AssertionError("fixture hash drifted through JSON round-trip")
    if rt.state.hash() != fx.state.hash():
        raise AssertionError("state hash drifted through JSON round-trip")


@check("fixture: declared guards pass, tampering is caught")
def _c16():
    fx = fixtures.build_f_st_01_s1()
    fixtures.assert_declared(fx)
    tampered = dataclasses.replace(fx, fact_text="Someone else trained me.")
    try:
        fixtures.assert_declared(tampered)
    except FixtureError:
        return
    raise AssertionError("tampered fixture passed its declared guards")


# ------------------------------------------------------------ conditions
@check("conditions: differ only by the state line")
def _c17():
    fx = fixtures.build_f_st_01_s1()
    a, b = conditions.build_conditions(fx)
    d = conditions.diff_conditions(a, b)
    indices = [entry["index"] for entry in d["differing_indices"]]
    if indices != [0]:
        raise AssertionError(f"unexpected differing indices: {indices}")
    if not d["params_equal"]:
        raise AssertionError("sampling params differ between conditions")
    if a.prompt_hash == b.prompt_hash:
        raise AssertionError("prompt hashes identical: the fact is not actually supplied in A")


@check("conditions: a widened difference is refused")
def _c18():
    fx = fixtures.build_f_st_01_s1()
    a, b = conditions.build_conditions(fx)
    widened = dataclasses.replace(
        b, messages=({"role": "system", "content": b.system_text() + "\nBe brief."}, b.messages[-1])
    )
    try:
        conditions.assert_separated(a, widened, fx)
    except SeparationError:
        return
    raise AssertionError("widened difference accepted as separated")


# ------------------------------------------------------------------ plan
@check("plan: eighteen calls, ceiling thirty")
def _c19():
    fx = fixtures.build_f_st_01_s1()
    conds = conditions.build_conditions(fx)
    slots = (
        Slot("S1-A", "org/model-a", "provider-a"),
        Slot("S1-B", "org/model-b", "provider-b"),
        Slot("S1-C", "org/model-c", "provider-c"),
    )
    calls = plan.build_plan(slots, conds)
    if len(calls) != 18:
        raise AssertionError(f"plan has {len(calls)} calls, expected 18")
    if plan.PLANNED_CALLS != 18 or plan.HARD_CEILING != 30:
        raise AssertionError("caps drifted from the authorised schedule")


@check("plan: unassigned slots refuse to become calls")
def _c20():
    fx = fixtures.build_f_st_01_s1()
    try:
        plan.build_plan(plan.default_slots(), conditions.build_conditions(fx))
    except PlanIncomplete:
        return
    raise AssertionError("unassigned slots produced a call plan")


@check("plan: overrun and inflated ceiling refused")
def _c21():
    for attempted, ceiling in ((31, 30), (30, 40)):
        try:
            plan.check_caps(attempted, ceiling)
        except PlanIncomplete:
            continue
        raise AssertionError(f"caps accepted attempted={attempted} ceiling={ceiling}")


# ----------------------------------------------------------- provenance
def _sample_request() -> GenerateRequest:
    fx = fixtures.build_f_st_01_s1()
    a, _ = conditions.build_conditions(fx)
    return GenerateRequest(
        messages=a.messages, model_id="org/model-a", provider="provider-a", params=a.params
    )


@check("provenance: absent fields are marked unavailable")
def _c22():
    rec = provenance.build_record(
        test_id="S1",
        fixture_id=fixtures.FIXTURE_ID,
        fixture_version=fixtures.FIXTURE_VERSION,
        request=_sample_request(),
        state_hash_before="a" * 64,
        state_hash_after="b" * 64,
    )
    if rec["raw_body"] != UNAVAILABLE or rec["text"] != UNAVAILABLE:
        raise AssertionError("absent response fields were not marked unavailable")
    if rec["model_revision"] != UNAVAILABLE or rec["system_fingerprint"] != UNAVAILABLE:
        raise AssertionError("absent provider metadata was not marked unavailable")


@check("provenance: a record with no raw body is inadmissible")
def _c23():
    rec = provenance.build_record(
        test_id="S1",
        fixture_id=fixtures.FIXTURE_ID,
        fixture_version=fixtures.FIXTURE_VERSION,
        request=_sample_request(),
        state_hash_before="a" * 64,
        state_hash_after="b" * 64,
        provider_used="provider-a",
    )
    problems = provenance.check_admissible(rec)
    if not any("raw_body" in p for p in problems):
        raise AssertionError(f"missing raw body not flagged: {problems}")
    if any("provider_used" in p for p in problems):
        raise AssertionError("provider_used present but flagged")


@check("provenance: record writes and re-reads unchanged")
def _c24():
    rec = provenance.build_record(
        test_id="S1",
        fixture_id=fixtures.FIXTURE_ID,
        fixture_version=fixtures.FIXTURE_VERSION,
        request=_sample_request(),
        state_hash_before="a" * 64,
        state_hash_after="b" * 64,
        provider_used="provider-a",
    )
    tmp = Path(tempfile.mkdtemp(prefix="_selftest-", dir=Path(__file__).resolve().parent))
    try:
        path = provenance.write_record(tmp, rec)
        import json

        back = json.loads(path.read_text(encoding="utf-8"))
        if back != rec:
            raise AssertionError("record changed on write/read")
        if (tmp / "raw_body.json").exists():
            raise AssertionError("a raw body file was written without a raw body")
    finally:
        shutil.rmtree(tmp, ignore_errors=True)


@check("cleanup: no test scratch directories remain (self-test or grant)")
def _c25():
    if HARNESS_DIR != Path(__file__).resolve().parent:
        raise AssertionError("scratch scan is not pointed at the harness directory")
    stray = find_scratch()
    if stray:
        raise AssertionError(f"leftover scratch directories: {[p.name for p in stray]}")


# ------------------------------------------------------------------ budget
@check("budget: paid spend ceiling is zero and nothing is permitted to charge")
def _c26():
    if budget.MAX_PAID_SPEND_USD != 0.00:
        raise AssertionError(f"paid ceiling is {budget.MAX_PAID_SPEND_USD}, not zero")
    if budget.AUTO_TOPUP_PERMITTED or budget.AUTOMATIC_FALLBACK_PERMITTED:
        raise AssertionError("safety flags are not False")
    if budget.PAYMENT_METHOD_ON_FILE is not False:
        raise AssertionError("policy assumes a payment method exists")
    table = budget.permission_table()
    if "no payment method" not in table["huggingface_paid"]:
        raise AssertionError("paid route not described as unreachable")
    if "permission" not in table["automatic_paid_fallback"]:
        raise AssertionError("automatic fallback is not marked prohibited")
    try:
        budget.check_spend(0.01)
    except budget.BudgetRefusal:
        return
    raise AssertionError("a one-cent spend was accepted")


@check("budget: hosted is refused because no free allowance exists")
def _c27():
    try:
        budget.check_provider("org/model")
    except budget.BudgetRefusal as exc:
        text = str(exc)
        if "no included monthly allowance" not in text:
            raise AssertionError(f"refusal does not cite the allowance: {text}")
        if "local route" not in text:
            raise AssertionError(f"refusal does not point at the local fallback: {text}")
    else:
        raise AssertionError("hosted route was permitted")
    if budget.FREE_ALLOWANCE_EXISTS is not False:
        raise AssertionError("policy still assumes a free allowance exists")


@check("budget: the local route never depends on payment state")
def _c28():
    for state in (True, False):
        if budget.check_provider("local:koboldcpp", payment_method_on_file=state) != "local":
            raise AssertionError(f"local route refused with payment_method_on_file={state}")


@check("budget: a 402 response is a payment demand, not a retry")
def _c29():
    class Err(Exception):
        status_code = 402

    if not budget.is_payment_demand(Err("nope")):
        raise AssertionError("402 not recognised as payment demand")
    try:
        budget.check_payment_demand(Err("nope"))
    except budget.PaymentDemanded as exc:
        text = str(exc)
        if "do not retry" not in text or "local route" not in text:
            raise AssertionError(f"stop instruction incomplete: {text}")
        return
    raise AssertionError("402 raised something other than PaymentDemanded")


@check("budget: exhausted-allowance text is a payment demand; other errors are not")
def _c30():
    for text in ("insufficient credits", "quota exceeded", "add a payment method"):
        if not budget.is_payment_demand(RuntimeError(text)):
            raise AssertionError(f"not recognised as payment demand: {text}")
    for text in ("connection reset", "invalid api token", "model not found"):
        if budget.is_payment_demand(RuntimeError(text)):
            raise AssertionError(f"misread as payment demand: {text}")


@check("budget: a route change is refused, not tolerated")
def _c31():
    try:
        budget.check_no_route_change("local:koboldcpp", "local:ollama")
    except budget.RouteChangeRefused:
        pass
    else:
        raise AssertionError("local-to-local route change accepted")
    try:
        budget.check_no_route_change("org/model", "local:koboldcpp")
    except budget.RouteChangeRefused:
        return
    raise AssertionError("hosted-to-local route change accepted")


# ------------------------------------------------------------ local route
@check("local: a non-loopback endpoint is refused")
def _c32():
    for bad in ("http://192.168.1.10:5001", "https://api.example.com", "http://0.0.0.0:8080"):
        try:
            local_adapter.LocalEndpoint(base_url=bad)
        except local_adapter.LocalEndpointError:
            continue
        raise AssertionError(f"non-loopback endpoint accepted: {bad}")
    ok = local_adapter.LocalEndpoint(base_url="http://127.0.0.1:5001")
    if not ok.is_loopback:
        raise AssertionError("loopback endpoint not recognised")


@check("local: adapter cannot be constructed without a grant")
def _c33():
    ep = local_adapter.LocalEndpoint(base_url="http://127.0.0.1:5001")
    try:
        local_adapter.LocalHTTPAdapter(grant=None, endpoint=ep)
    except ExecutionNotAuthorised:
        return
    raise AssertionError("local adapter constructed without a grant")


@check("local: same GenerateRequest contract as the hosted adapter")
def _c34():
    ep = local_adapter.LocalEndpoint(base_url="http://127.0.0.1:5001")
    grant = types.SimpleNamespace(grants_execution=True)
    adapter_obj = local_adapter.LocalHTTPAdapter(grant=grant, endpoint=ep)
    req = _sample_request()
    payload = adapter_obj.payload(req)
    if payload["messages"] != [dict(m) for m in req.messages]:
        raise AssertionError("local payload altered the messages")
    if payload["stream"] is not False:
        raise AssertionError("local payload enabled streaming")
    if set(req.params) - set(payload):
        raise AssertionError("local payload dropped sampling params")
    if adapter_obj.capabilities.get("paid") is not False:
        raise AssertionError("local route reported as paid")
    # L3 replaced the blanket refusal with a real call, so the contract is now
    # checked through an injected transport: exactly one send, transcribed body.
    sent = []

    def transport(url, payload, timeout):
        sent.append(url)
        return 200, json.dumps(
            {"choices": [{"message": {"content": "ok"}, "finish_reason": "stop"}]}
        )

    real = local_adapter.LocalHTTPAdapter(grant=grant, endpoint=ep, post_transport=transport)
    response = real.generate(req)
    if len(sent) != 1:
        raise AssertionError(f"local adapter must send exactly once, sent {len(sent)}")
    if response.text != "ok":
        raise AssertionError("local adapter did not transcribe the injected body")
    if json.loads(json.dumps(response.raw_body))["choices"][0]["message"]["content"] != "ok":
        raise AssertionError("raw body not preserved verbatim")


@check("route: the factory cannot select a hosted route under any grant")
def _c35():
    grant = types.SimpleNamespace(grants_execution=True)
    ep = "http://127.0.0.1:5001"
    obj = adapter_module.build_adapter("local:koboldcpp", grant=grant, local_endpoint=ep)
    if not isinstance(obj, local_adapter.LocalHTTPAdapter):
        raise AssertionError("local provider did not select the local adapter")
    try:
        adapter_module.build_adapter("local:koboldcpp", grant=grant)
    except AdapterConfigError:
        pass
    else:
        raise AssertionError("local route guessed an endpoint")
    for provider in ("org/model", "meta-llama/Llama-3-8B:fastest"):
        try:
            adapter_module.build_adapter(provider, grant=grant)
        except (budget.BudgetRefusal, AdapterConfigError):
            continue
        raise AssertionError(f"hosted route was buildable: {provider}")
    try:
        adapter_module.build_adapter("local:koboldcpp", grant=None, local_endpoint=ep)
    except ExecutionNotAuthorised:
        return
    raise AssertionError("local adapter constructed without a grant")


# ------------------------------------------------------------------ cost
@check("cost: no quote means no spend")
def _c36():
    try:
        cost.preflight(None, 18, 256)
    except cost.CostRefusal:
        return
    raise AssertionError("an unpriced plan was accepted")


@check("cost: a paid quote is refused before dispatch")
def _c37():
    q = cost.PriceQuote(
        model="org/m",
        provider="p",
        currency="USD",
        input_per_1k=0.5,
        output_per_1k=1.5,
        retrieved_at="2026-11-13",
        source="owner statement",
        is_free_tier=False,
    )
    try:
        cost.preflight(q, 18, 256)
    except cost.CostRefusal:
        return
    raise AssertionError("paid quote priced into the ceiling")


@check("cost: a verified free quote prices at zero")
def _c38():
    q = cost.PriceQuote(
        model="org/m",
        provider="p",
        currency="USD",
        input_per_1k=0.0,
        output_per_1k=0.0,
        retrieved_at="2026-11-13",
        source="owner statement",
        is_free_tier=True,
    )
    decision = cost.preflight(q, 18, 256)
    if decision["worst_case"] != 0.0 or decision["ceiling"] != 0.0:
        raise AssertionError(f"free quote priced non-zero: {decision}")


# -------------------------------------------------------------- dispatch
@check("dispatch: the ceiling bites at the call, not at the plan")
def _c39():
    ledger = dispatch.CallLedger()
    disp = dispatch.Dispatcher(ledger, OfflineAdapter())
    req = _sample_request()
    for _ in range(dispatch.CallLedger().ceiling):
        try:
            disp.dispatch("S1-A", req)
        except ExecutionNotAuthorised:
            pass
    if ledger.attempted != 30:
        raise AssertionError(f"ledger counted {ledger.attempted}, expected 30")
    try:
        disp.dispatch("S1-A", req)
    except dispatch.CallCeilingExceeded:
        return
    raise AssertionError("call 31 was dispatched")


@check("dispatch: failed calls consume the ledger")
def _c40():
    ledger = dispatch.CallLedger(ceiling=3)
    disp = dispatch.Dispatcher(ledger, OfflineAdapter())
    req = _sample_request()
    for _ in range(3):
        try:
            disp.dispatch("S1-B", req)
        except ExecutionNotAuthorised:
            pass
    if ledger.consumed_by_failure != 3 or ledger.remaining() != 0:
        raise AssertionError(
            f"failures did not consume the budget: {ledger.consumed_by_failure} / remaining {ledger.remaining()}"
        )
    try:
        ledger.reserve("S1-B")
    except dispatch.CallCeilingExceeded:
        return
    raise AssertionError("reserve succeeded after the budget was spent")


def _grant_payload(**over):
    payload = {
        "owner_signature": "Owner Signature 2026-11-13",
        "date": "2026-11-13",
        "call_ceiling": 18,
        "cost_ceiling_usd": 0.0,
        "models": {s: "org/model-" + s.lower() for s in plan.SLOTS},
        "providers": {s: "provider-" + s.lower() for s in plan.SLOTS},
    }
    recompute = over.pop("recompute", True)
    payload.update(over)
    payload["caps_sha256"] = authorization.caps_sha256(payload) if recompute else "0" * 64
    return payload


@contextlib.contextmanager
def scratch_grant(payload):
    """A disposable grant file inside a disposable directory we own."""
    with scratch_dir("_grant-") as tmp:
        path = tmp / "grant.json"
        path.write_text(json.dumps(payload), encoding="utf-8")
        yield path


# ----------------------------------------------------------------- grant
@check("grant: an incomplete grant is refused")
def _c41():
    payload = _grant_payload()
    payload.pop("caps_sha256")
    with scratch_grant(payload) as path:
        try:
            authorization.load_grant(path)
        except ExecutionNotAuthorised:
            return
    raise AssertionError("grant without caps_sha256 accepted")


@check("grant: wrong field types are refused")
def _c42():
    for over in ({"models": ["a", "b", "c"]}, {"call_ceiling": "18"}, {"date": "13/11/2026"}):
        with scratch_grant(_grant_payload(**over)) as path:
            try:
                authorization.load_grant(path)
            except ExecutionNotAuthorised:
                continue
            raise AssertionError(f"malformed grant accepted: {over}")


@check("grant: editing caps after signing breaks the hash")
def _c43():
    payload = _grant_payload()
    payload["call_ceiling"] = 30  # edited after signing, hash not recomputed
    with scratch_grant(payload) as path:
        try:
            authorization.load_grant(path)
        except ExecutionNotAuthorised as exc:
            if "caps_sha256" not in str(exc):
                raise AssertionError(f"wrong refusal reason: {exc}")
            return
    raise AssertionError("edited grant accepted")


@check("grant: a paid cost ceiling cannot be signed")
def _c44():
    with scratch_grant(_grant_payload(cost_ceiling_usd=5.0)) as path:
        try:
            authorization.load_grant(path)
        except ExecutionNotAuthorised:
            return
    raise AssertionError("grant with a paid ceiling accepted")


# ------------------------------------------------------- evidence integrity
@check("evidence: records are write-once and hash-verified")
def _c45():
    rec = provenance.build_record(
        test_id="S1",
        fixture_id=fixtures.FIXTURE_ID,
        fixture_version=fixtures.FIXTURE_VERSION,
        request=_sample_request(),
        state_hash_before="a" * 64,
        state_hash_after="b" * 64,
        provider_used="provider-a",
    )
    tmp = Path(tempfile.mkdtemp(prefix="_selftest-", dir=Path(__file__).resolve().parent))
    try:
        path = provenance.write_record(tmp, rec)
        ok, why = provenance.verify_record(tmp)
        if not ok:
            raise AssertionError(f"fresh record failed verification: {why}")
        try:
            provenance.write_record(tmp, rec)
        except provenance.RecordError:
            pass
        else:
            raise AssertionError("record was overwritten")
        path.write_text(path.read_text(encoding="utf-8").replace("S1", "S9"), encoding="utf-8")
        ok2, _ = provenance.verify_record(tmp)
        if ok2:
            raise AssertionError("tampered record still verified")
    finally:
        shutil.rmtree(tmp, ignore_errors=True)


@check("evidence: writing outside eval/ is refused")
def _c46():
    outside = Path(tempfile.mkdtemp(prefix="memory09-outside-"))
    try:
        try:
            provenance.write_record(outside, {"test_id": "S1"})
        except provenance.RecordError:
            return
        raise AssertionError("evidence written outside eval/")
    finally:
        shutil.rmtree(outside, ignore_errors=True)


@check("state: a snapshot cannot be mutated after hashing")
def _c47():
    st = fixtures.build_f_st_01_s1().state
    snap = st.snapshot()
    before_hash = snap.sha256
    data = snap.data
    for kind in list(data):
        data[kind].clear()
    if snap.sha256 != before_hash:
        raise AssertionError("snapshot hash changed after mutation")
    if snap.entity_count() == 0:
        raise AssertionError("snapshot contents were erased by a caller mutation")


# --------------------------------------------------------------- adapter
@check("adapter: unsupported client params are caught offline")
def _c48():
    bad = adapter_module.check_params_supported({"wibble_param": 1})
    if bad != ["wibble_param"]:
        raise AssertionError(f"unsupported param not detected: {bad}")
    accepted = adapter_module.client_chat_parameters()
    if accepted is not None and "temperature" in accepted:
        if adapter_module.check_params_supported({"temperature": 0.8}):
            raise AssertionError("a supported param was flagged unsupported")


@check("adapter: transcription works on a synthetic body")
def _c49():
    body = {
        "id": "req-1",
        "model": "org/model-a",
        "provider": "provider-a",
        "created": 1700000000,
        "system_fingerprint": "fp_test",
        "usage": {"total_tokens": 12},
        "choices": [{"message": {"content": "Corvin Ashgrave of Ossary trained me."}}],
    }
    resp = adapter_module.transcribe(body)
    if "Corvin Ashgrave" not in str(resp.text):
        raise AssertionError(f"transcription lost the text: {resp.text!r}")
    if resp.raw_body != body:
        raise AssertionError("raw body not preserved verbatim")
    ok, problems = resp.check_admissible()
    if not ok:
        raise AssertionError(f"transcribed response inadmissible: {problems}")
    try:
        adapter_module.transcribe("not a mapping")
    except AdapterConfigError:
        return
    raise AssertionError("non-mapping body accepted")


@check("authorization: no code path writes the grant file")
def _c50():
    src = Path(authorization.__file__).read_text(encoding="utf-8")
    for forbidden in ("write_text", "open(", "json.dump("):
        if forbidden in src:
            raise AssertionError(f"authorization.py contains a write path: {forbidden}")


@check("route: the CLI run path goes through the dispatcher, not around it")
def _c51():
    src = Path(cli_module.__file__).read_text(encoding="utf-8")
    if "Dispatcher(" not in src:
        raise AssertionError("cli does not construct a Dispatcher")
    if src.count("generate(") > 0:
        raise AssertionError("cli calls generate() directly, bypassing the ledger")


@check("dispatch: a payment demand stops the call and still consumes the slot")
def _c52():
    class Buyer(Exception):
        status_code = 402

    class Asker:
        def generate(self, request):
            raise Buyer("402 Payment Required")

    ledger = dispatch.CallLedger(ceiling=5)
    disp = dispatch.Dispatcher(ledger, Asker())
    try:
        disp.dispatch("S1-A", _sample_request())
    except budget.PaymentDemanded:
        if ledger.attempted != 1 or ledger.consumed_by_failure != 1:
            raise AssertionError(f"ledger after demand: {ledger.attempted}/{ledger.consumed_by_failure}")
        return
    raise AssertionError("payment demand was swallowed instead of raised")


# ------------------------------------------------- scratch lifecycle (P1.1)
@check("scratch: a leaked _grant- directory is detected")
def _c53():
    injected = HARNESS_DIR / "_grant-injected-probe"
    injected.mkdir()
    try:
        if injected not in find_scratch():
            raise AssertionError("grants scratch not detected by the cleanliness scan")
    finally:
        shutil.rmtree(injected, ignore_errors=True)
    if find_scratch():
        raise AssertionError("scan did not return to clean after removal")


@check("scratch: a leaked _selftest- directory is detected")
def _c54():
    injected = HARNESS_DIR / "_selftest-injected-probe"
    injected.mkdir()
    try:
        if injected not in find_scratch():
            raise AssertionError("selftest scratch not detected by the cleanliness scan")
    finally:
        shutil.rmtree(injected, ignore_errors=True)


@check("scratch: owned directories are removed on success and on failure")
def _c55():
    with scratch_dir("_grant-") as made:
        pass
    if made.exists():
        raise AssertionError("scratch directory survived the success path")
    made_on_failure = None
    try:
        with scratch_dir("_grant-") as made_on_failure:
            raise RuntimeError("simulated failure inside the with-block")
    except RuntimeError:
        pass
    if made_on_failure is None or made_on_failure.exists():
        raise AssertionError("scratch directory survived the failure path")


@check("scratch: cleanup does not touch unrelated files or directories")
def _c56():
    keep_dir = HARNESS_DIR / "_keepme-dir"
    keep_file = HARNESS_DIR / "_keepme.file"
    keep_dir.mkdir(exist_ok=True)
    keep_file.write_text("unrelated\n", encoding="utf-8")
    try:
        with scratch_dir("_grant-") as made:
            pass
        if keep_dir.exists() is False or keep_file.exists() is False:
            raise AssertionError("cleanup removed a file or directory it did not create")
        if made.exists():
            raise AssertionError("owned scratch directory was not removed")
    finally:
        shutil.rmtree(keep_dir, ignore_errors=True)
        keep_file.unlink(missing_ok=True)


@check("scratch: a failing grant validation leaves no new directories")
def _c57():
    before = {p.name for p in find_scratch()}
    with scratch_grant(_grant_payload()) as path:
        payload = json.loads(path.read_text(encoding="utf-8"))
        payload["call_ceiling"] = 99
        path.write_text(json.dumps(payload), encoding="utf-8")
        try:
            authorization.load_grant(path)
        except ExecutionNotAuthorised:
            pass
        else:
            raise AssertionError("tampered grant accepted")
    after = {p.name for p in find_scratch()}
    if after - before:
        raise AssertionError(f"scratch accumulated across a failing validation: {after - before}")


@check("cleanup: a complete run leaves no scratch directories behind")
def _c58():
    stray = find_scratch()
    if stray:
        raise AssertionError(f"run leaked scratch directories: {[p.name for p in stray]}")


# ------------------------------------------------- hosted refusal (P1.2)
@check("hosted: the adapter object cannot be built even with a valid grant")
def _c59():
    grant = types.SimpleNamespace(grants_execution=True)
    try:
        adapter_module.HuggingFaceAdapter(grant=grant)
    except budget.BudgetRefusal:
        return
    raise AssertionError("hosted adapter was constructed with a valid grant")


@check("hosted: refusal names the local fallback and makes no call")
def _c60():
    class Counting(OfflineAdapter):
        calls = 0

        def generate(self, request):
            Counting.calls += 1
            return super().generate(request)

    adapter_obj = Counting()
    try:
        adapter_obj.generate(_sample_request())
    except ExecutionNotAuthorised:
        pass
    # The refusal happens before any adapter exists, so nothing can be dispatched.
    try:
        adapter_module.build_adapter("org/model", grant=types.SimpleNamespace(grants_execution=True))
    except budget.BudgetRefusal as exc:
        if "local" not in str(exc):
            raise AssertionError("refusal does not recommend a local route")
        if Counting.calls != 1:
            raise AssertionError("an unexpected number of adapter calls occurred")
        return
    raise AssertionError("hosted route was built")


@check("hosted: the CLI cannot reach a hosted route")
def _c61():
    buf = io.StringIO()
    with contextlib.redirect_stdout(buf):
        rc = cli_module.main(["routes"])
    out = buf.getvalue()
    if rc != 0:
        raise AssertionError(f"routes command failed: exit {rc}")
    lowered = out.lower()
    if "refused" not in lowered or "no included monthly allowance" not in lowered:
        raise AssertionError(f"routes output does not show the refusal: {out}")
    if "allowed while inside the free allowance" in lowered:
        raise AssertionError("routes output still advertises a free hosted allowance")


@check("hosted: no module constructs an adapter outside the factory")
def _c62():
    for module in (cli_module, dispatch):
        src = Path(module.__file__).read_text(encoding="utf-8")
        for name in ("HuggingFaceAdapter(", "LocalHTTPAdapter("):
            if name in src:
                raise AssertionError(f"{module.__name__} constructs {name} directly")


# ------------------------------------------- local surfaces (plan L1)
@check("local L1: chat payload serialises model, messages, params, no stream")
def _c63():
    req = _sample_request()
    payload = local_surfaces.build_chat_payload(req)
    if payload["model"] != req.model_id:
        raise AssertionError("model not serialised")
    if payload["messages"] != [dict(m) for m in req.messages]:
        raise AssertionError("messages altered")
    if payload["stream"] is not False:
        raise AssertionError("streaming enabled")
    for key, value in req.params.items():
        if payload.get(key) != value:
            raise AssertionError(f"param {key} lost in chat serialisation")


@check("local L1: native payload renames limits and records the translation")
def _c64():
    req = _sample_request()
    payload = local_surfaces.build_native_payload(req, stop_sequences=["</s>"])
    if "max_length" not in payload or "max_tokens" in payload:
        raise AssertionError(f"token limit not mapped: {sorted(payload)}")
    if payload["max_length"] != req.params["max_tokens"]:
        raise AssertionError("token limit value changed during mapping")
    if payload["sampler_translation"].get("max_tokens") != "max_length":
        raise AssertionError("translation not recorded")
    if payload["stop_sequence"] != ["</s>"]:
        raise AssertionError("stop sequences not carried")
    if "\n\n" not in payload["prompt"]:
        raise AssertionError("prompt was not flattened from the message pair")


@check("local L1: an unmappable parameter is refused, not dropped")
def _c65():
    req = adapter_module.GenerateRequest(
        messages=({"role": "user", "content": "x"},),
        model_id="m",
        provider="local:koboldcpp",
        params={},
    )
    object.__setattr__(req, "params", {"temperature": 0.8})
    try:
        local_surfaces.build_native_payload(req, stop_sequences=None)
    except local_surfaces.SurfaceError:
        raise AssertionError("a valid mapped param was refused")
    # A parameter with no native equivalent must raise rather than vanish.
    bad = adapter_module.GenerateRequest(
        messages=({"role": "user", "content": "x"},),
        model_id="m",
        provider="local:koboldcpp",
        params={},
    )
    object.__setattr__(bad, "params", {"logit_bias": {}})
    try:
        local_surfaces.build_native_payload(bad)
    except local_surfaces.SurfaceError:
        return
    raise AssertionError("unmappable parameter was silently dropped")


@check("local L1: surface detection identifies both shapes and rejects others")
def _c66():
    if local_surfaces.detect_surface(local_fixtures.CHAT_OK) != local_surfaces.SURFACE_CHAT:
        raise AssertionError("chat body misdetected")
    if local_surfaces.detect_surface(local_fixtures.NATIVE_OK) != local_surfaces.SURFACE_NATIVE:
        raise AssertionError("native body misdetected")
    try:
        local_surfaces.detect_surface(local_fixtures.UNRECOGNISED)
    except local_surfaces.SurfaceError:
        pass
    else:
        raise AssertionError("unrecognised body accepted")
    try:
        local_surfaces.detect_surface("not a mapping")
    except local_surfaces.SurfaceError:
        return
    raise AssertionError("non-mapping body accepted")


@check("local L1: chat body transcribes text, usage and finish reason")
def _c67():
    surface, resp = local_surfaces.transcribe_any(local_fixtures.CHAT_OK)
    if surface != local_surfaces.SURFACE_CHAT:
        raise AssertionError("wrong surface label")
    if "Corvin Ashgrave" not in str(resp.text):
        raise AssertionError(f"text lost: {resp.text!r}")
    if resp.usage == adapter_module.UNAVAILABLE:
        raise AssertionError("usage lost")
    if resp.finish_reason != "stop":
        raise AssertionError(f"finish reason not captured: {resp.finish_reason}")
    ok, problems = resp.check_admissible()
    if not ok:
        raise AssertionError(f"chat record inadmissible: {problems}")


@check("local L1: a truncated chat reply is incomplete, not a failure")
def _c68():
    _, resp = local_surfaces.transcribe_any(local_fixtures.CHAT_TRUNCATED)
    if not resp.incomplete:
        raise AssertionError("truncation not flagged")
    if resp.error not in (None, adapter_module.UNAVAILABLE):
        raise AssertionError(f"truncation recorded as an error: {resp.error}")
    if not resp.text:
        raise AssertionError("truncated text discarded")


@check("local L1: native body transcribes text and token counts")
def _c69():
    surface, resp = local_surfaces.transcribe_any(local_fixtures.NATIVE_OK)
    if surface != local_surfaces.SURFACE_NATIVE:
        raise AssertionError("wrong surface label")
    if "Corvin Ashgrave" not in str(resp.text):
        raise AssertionError(f"text lost: {resp.text!r}")
    usage = resp.usage
    if not isinstance(usage, dict) or usage.get("completion_tokens") != 9:
        raise AssertionError(f"token counts not mapped: {usage}")
    if resp.finish_reason != "stop":
        raise AssertionError("finish reason lost")


@check("local L1: native truncation and empty results are distinguished")
def _c70():
    _, trunc = local_surfaces.transcribe_any(local_fixtures.NATIVE_TRUNCATED)
    if not trunc.incomplete or trunc.error not in (None, adapter_module.UNAVAILABLE):
        raise AssertionError("native truncation mishandled")
    _, empty = local_surfaces.transcribe_any(local_fixtures.NATIVE_EMPTY)
    if empty.text != adapter_module.UNAVAILABLE:
        raise AssertionError("empty results produced text")
    if not empty.error or "empty results" not in str(empty.error):
        raise AssertionError(f"empty results not recorded as such: {empty.error}")


@check("local L1: unknown fields are preserved in the raw body")
def _c70b():
    _, resp = local_surfaces.transcribe_any(local_fixtures.NATIVE_UNKNOWN_FIELDS)
    if resp.raw_body != local_fixtures.NATIVE_UNKNOWN_FIELDS:
        raise AssertionError("raw body not preserved verbatim")
    if "some_future_field" not in resp.raw_body["results"][0]:
        raise AssertionError("unknown nested field dropped from raw body")
    if "Corvin Ashgrave" not in str(resp.text):
        raise AssertionError("known field lost while tolerating unknown ones")


@check("local L1: error bodies become records, not exceptions")
def _c71():
    for fixture in (local_fixtures.CHAT_ERROR, local_fixtures.NATIVE_ERROR):
        surface, resp = local_surfaces.transcribe_any(fixture)
        if surface not in (*local_surfaces.all_surfaces(), local_surfaces.SURFACE_UNKNOWN):
            raise AssertionError(f"error body produced an impossible surface: {surface}")
        if resp.text != adapter_module.UNAVAILABLE:
            raise AssertionError("error body produced text")
        if not resp.error or "server error" not in str(resp.error):
            raise AssertionError(f"error not recorded: {resp.error}")
    # A shapeless error body (the common 500) must be handled, not rejected.
    surface, resp = local_surfaces.transcribe_any({"error": "boom"})
    if surface != local_surfaces.SURFACE_UNKNOWN:
        raise AssertionError("shapeless error body mislabelled")
    if resp.raw_body != {"error": "boom"}:
        raise AssertionError("shapeless error body not preserved")


@check("local L1: both surfaces yield one identical record schema")
def _c72():
    _, chat = local_surfaces.transcribe_any(local_fixtures.CHAT_OK)
    _, native = local_surfaces.transcribe_any(local_fixtures.NATIVE_OK)
    if type(chat) is not type(native):
        raise AssertionError("surfaces yield different record types")
    if set(vars(chat)) != set(vars(native)):
        raise AssertionError("surfaces yield different record fields")
    if user := [f for f in vars(chat) if f not in vars(native)]:
        raise AssertionError(f"field mismatch: {user}")


@check("local L1: prompt flattening is recorded as a distinct condition")
def _c73():
    req = _sample_request()
    chat = local_surfaces.build_chat_payload(req)
    native = local_surfaces.build_native_payload(req)
    if native["prompt"] == chat["messages"][0]["content"]:
        raise AssertionError("native prompt silently equals the system message alone")
    if local_surfaces.FLATTEN_TEMPLATE != "system\\n\\nuser":
        raise AssertionError("flatten template name drifted")
    joined = local_surfaces.flatten_prompt(req)
    if "Corvin Ashgrave" not in joined or req.messages[-1]["content"] not in joined:
        raise AssertionError("flattening dropped content")


# --------------------------------------------------- L2 discovery (identity)
@check("L2: discovery refuses a non-loopback endpoint")
def _c75():
    try:
        discovery.probe("http://192.168.1.20:5001", transport=lambda u, t: (200, "{}"))
    except Exception as exc:
        if type(exc).__name__ not in ("DiscoveryRefused", "LocalEndpointError"):
            raise AssertionError(f"wrong refusal type: {type(exc).__name__}")
    else:
        raise AssertionError("non-loopback endpoint was probed")


@check("L2: only identity paths may be requested")
def _c76():
    paths = discovery.allowed_paths()
    if "/api/v1/generate" in paths:
        raise AssertionError("generate path is in the discovery allowlist")
    if len(paths) != 3:
        raise AssertionError(f"allowlist drifted: {paths}")
    calls = []

    def transport(url, timeout):
        calls.append(url)
        return (200, '{"result": "x"}')

    discovery.probe("http://127.0.0.1:5001", transport=transport)
    if len(calls) != 3:
        raise AssertionError(f"discovery made {len(calls)} requests, expected 3")
    if any("/generate" in url for url in calls):
        raise AssertionError("discovery hit a generation endpoint")


@check("L2: absent identity fields are recorded as unavailable")
def _c77():
    block = discovery.discover("http://127.0.0.1:5001", transport=lambda u, t: (200, "{}"))
    if block["server_version"] != adapter_module.UNAVAILABLE:
        raise AssertionError(f"missing version not marked: {block['server_version']}")
    if block["loaded_model"] != adapter_module.UNAVAILABLE:
        raise AssertionError("missing model not marked")
    if block["max_context_length"] != adapter_module.UNAVAILABLE:
        raise AssertionError("missing context length not marked")


@check("L2: a silent endpoint is recorded, not raised")
def _c78():
    block = discovery.discover("http://127.0.0.1:5001", transport=lambda u, t: (0, "not listening: refused"))
    if block["listening"] is not False:
        raise AssertionError("silent endpoint reported as listening")
    if "does not start servers" not in block["conclusion"]:
        raise AssertionError(f"conclusion missing: {block.get('conclusion')}")
    if block["server_version"] != adapter_module.UNAVAILABLE:
        raise AssertionError("version invented for a silent endpoint")


@check("L2: raw identity bodies are preserved verbatim")
def _c79():
    body = '{"result": "1.999", "extra": {"nested": true}}'
    block = discovery.discover("http://127.0.0.1:5001", transport=lambda u, t: (200, body))
    raw = block["probe"]["/api/extra/version"]["body"]
    if raw != body:
        raise AssertionError("identity body altered")
    if block["server_version"] != "1.999":
        raise AssertionError(f"version not parsed: {block['server_version']}")
    if block["listening"] is not True:
        raise AssertionError("listening not detected")


@check("L2: discovery is GET-only and has no generation path")
def _c80():
    src = Path(discovery.__file__).read_text(encoding="utf-8")
    if 'method="GET"' not in src:
        raise AssertionError("discovery does not pin the HTTP method to GET")
    for forbidden in ("POST", "generate(", "Dispatcher", "dispatch("):
        if forbidden in src:
            raise AssertionError(f"discovery source contains a non-identity path: {forbidden}")
    if "IDENTITY_PATHS" not in src:
        raise AssertionError("discovery has no allowlist")


@check("L2: discovery needs no execution grant and cannot imply one")
def _c81():
    if authorization.grant_status() != "ABSENT":
        raise AssertionError("a grant exists; discovery must be independent of, and no excuse for, execution")
    block = discovery.discover("http://127.0.0.1:5001", transport=lambda u, t: (0, "not listening"))
    if block["listening"] is not False:
        raise AssertionError("silent probe reported as listening")
    # Discovery must not create, write or require a grant.
    src = Path(discovery.__file__).read_text(encoding="utf-8")
    for forbidden in ("require_grant", "ExecutionGrant", "grants_execution", "P2-GRANT"):
        if forbidden in src:
            raise AssertionError(f"discovery references execution authorisation: {forbidden}")


@check("L2: the version field wins over result, and the server name is kept")
def _c82():
    real = '{"result": "KoboldCpp", "version": "1.120", "protected": false, "llm": false}'
    block = discovery.discover("http://127.0.0.1:5001", transport=lambda u, t: (200, real))
    if block["server_version"] != "1.120":
        raise AssertionError(f"version not taken from the version field: {block['server_version']}")
    if block["server_name"] != "KoboldCpp":
        raise AssertionError(f"server name lost: {block['server_name']}")
    # A bare numeric result (older/fixture shape) is still accepted as a version.
    block2 = discovery.discover("http://127.0.0.1:5001", transport=lambda u, t: (200, '{"result": "1.999"}'))
    if block2["server_version"] != "1.999":
        raise AssertionError("numeric result no longer accepted as a version")


@check("L3: the local adapter still refuses without a grant")
def _c83():
    from .local_adapter import LocalEndpoint, LocalHTTPAdapter

    endpoint = LocalEndpoint("http://127.0.0.1:5001")
    try:
        LocalHTTPAdapter(grant=None, endpoint=endpoint)
    except Exception:
        pass
    else:
        raise AssertionError("adapter was constructed without a grant")


@check("L3: one POST, body preserved, latency and status recorded")
def _c84():
    from .local_adapter import LocalEndpoint, LocalHTTPAdapter

    calls = []

    def transport(url, payload, timeout):
        calls.append((url, payload))
        return 200, json.dumps(
            {
                "choices": [{"message": {"content": "Blue."}, "finish_reason": "stop"}],
                "model": "qwen3.5-9b",
                "id": "cmpl-1",
                "usage": {"prompt_tokens": 12, "completion_tokens": 2},
            }
        )

    class G:
        grants_execution = True

    adapter = LocalHTTPAdapter(grant=G(), endpoint=LocalEndpoint("http://127.0.0.1:5001"), post_transport=transport)
    request = GenerateRequest(
        messages=({"role": "user", "content": "sky colour?"},),
        model_id="qwen3.5-9b",
        provider="local:koboldcpp",
        params={"max_tokens": 8},
    )
    response = adapter.generate(request)
    if len(calls) != 1:
        raise AssertionError(f"expected exactly one POST, got {len(calls)}")
    if response.text != "Blue.":
        raise AssertionError(f"text not transcribed: {response.text!r}")
    if not isinstance(response.raw_body, dict) or "choices" not in response.raw_body:
        raise AssertionError("raw body not preserved")
    if response.latency_ms is None or response.http_status != 200:
        raise AssertionError("latency/status not recorded")
    if response.provider_used != "http://127.0.0.1:5001":
        raise AssertionError("provider_used not recorded")


@check("L3: a malformed or error body is recorded, not retried")
def _c85():
    from .local_adapter import LocalEndpoint, LocalHTTPAdapter

    class G:
        grants_execution = True

    request = GenerateRequest(
        messages=({"role": "user", "content": "x"},),
        model_id="m",
        provider="local:koboldcpp",
        params={},
    )
    calls = []

    def bad_json(url, payload, timeout):
        calls.append(1)
        return 200, "not json at all"

    adapter = LocalHTTPAdapter(grant=G(), endpoint=LocalEndpoint("http://127.0.0.1:5001"), post_transport=bad_json)
    response = adapter.generate(request)
    if len(calls) != 1:
        raise AssertionError("unparsable body caused a retry")
    if response.error is None or "_unparsed" not in response.raw_body:
        raise AssertionError("unparsable body was not preserved as evidence")

    def server_error(url, payload, timeout):
        calls.append(1)
        return 500, json.dumps({"detail": "model not loaded"})

    adapter2 = LocalHTTPAdapter(grant=G(), endpoint=LocalEndpoint("http://127.0.0.1:5001"), post_transport=server_error)
    response2 = adapter2.generate(request)
    if response2.error is None or response2.http_status != 500:
        raise AssertionError("server error not recorded as a record")
    if len(calls) != 2:
        raise AssertionError("a failed call was retried")


@check("L3: the dispatcher stops at the grant ceiling without sending again")
def _c86():
    from .local_adapter import LocalEndpoint, LocalHTTPAdapter

    class G:
        grants_execution = True

    sent = []

    def transport(url, payload, timeout):
        sent.append(1)
        return 200, json.dumps({"choices": [{"message": {"content": "ok"}, "finish_reason": "stop"}]})

    adapter = LocalHTTPAdapter(grant=G(), endpoint=LocalEndpoint("http://127.0.0.1:5001"), post_transport=transport)
    ledger = dispatch.CallLedger(ceiling=1)
    dispatcher = dispatch.Dispatcher(ledger, adapter)
    request = GenerateRequest(
        messages=({"role": "user", "content": "x"},), model_id="m", provider="local:koboldcpp", params={}
    )
    dispatcher.dispatch("L3-01", request)
    try:
        dispatcher.dispatch("L3-01", request)
    except dispatch.CallCeilingExceeded:
        pass
    else:
        raise AssertionError("second dispatch was allowed past a ceiling of 1")
    if len(sent) != 1:
        raise AssertionError(f"the blocked call still reached the server: {len(sent)} sends")


@check("L3: verification reads a loopback endpoint and tolerates absent state")
def _c87():
    from . import l3 as l3mod
    from .adapter import canonical, sha256_text

    if not l3mod.is_loopback("http://127.0.0.1:5001"):
        raise AssertionError("loopback endpoint not recognised")
    if l3mod.is_loopback("https://api.example.com/v1"):
        raise AssertionError("a remote endpoint was accepted as loopback")
    record = {
        "raw_body": {"choices": [{"message": {"content": "ok"}, "finish_reason": "stop"}]},
        "text": "ok",
        "response_hash": sha256_text(canonical({"choices": [{"message": {"content": "ok"}, "finish_reason": "stop"}]})),
        "provider_used": "http://127.0.0.1:5001",
        "state_hash_before": "unavailable",
        "state_hash_after": "unavailable",
        "completeness": {
            "http_status": 200,
            "latency_ms": 12,
            "usage": {"completion_tokens": 1},
            "finish_reason": "stop",
            "incomplete": False,
        },
        "ledger": {"attempted": 1, "ceiling": 1},
        "runtime": {"model_file": {"sha256": "0" * 64}},
    }
    rows = l3mod.verify_l3(record)
    failed = [r for r in rows if r["result"] != "PASS"]
    if failed:
        raise AssertionError(f"clean record failed verification: {failed}")
    record["provider_used"] = "https://api.example.com/v1"
    rows2 = l3mod.verify_l3(record)
    if not any(r["requirement"] == "no hosted route used" and r["result"] == "FAIL" for r in rows2):
        raise AssertionError("a hosted endpoint passed the L3 verification")


@check("S1: the frozen readiness baseline carries full digests and matches the code")
def _c88():
    """Abbreviated hashes cannot verify identity; the manifest must carry all 64 chars.

    Fixture/state/prompt values are re-derived here, so a code change that shifts a
    prompt fails the suite instead of quietly invalidating the pinned baseline.
    """
    import hashlib as _hashlib
    import json as _json
    from pathlib import Path as _Path

    from . import preflight as _preflight

    def need_digest(value, label):
        if not isinstance(value, str) or len(value) != 64:
            raise AssertionError(f"{label}: not a full 64-character digest: {value!r}")
        if any(c not in "0123456789abcdef" for c in value):
            raise AssertionError(f"{label}: not lowercase hex: {value!r}")

    fixture = fixtures.build_f_st_01_s1()
    a, b = conditions.build_conditions(fixture)
    expected = {
        "fixture": fixture.sha256(),
        "state_sha256": fixture.state.hash(),
        "prompts.A": a.prompt_hash(),
        "prompts.B": b.prompt_hash(),
    }
    def check_manifest(manifest, origin):
        """Every manifest on disk, not just the pinned one: the guard follows the file."""
        actual = {
            "fixture": manifest["fixture"]["sha256"],
            "state_sha256": manifest["state_sha256"],
            "prompts.A": manifest["prompts"]["A"],
            "prompts.B": manifest["prompts"]["B"],
        }
        for label, value in actual.items():
            need_digest(value, f"{origin}:{label}")
            if value != expected[label]:
                raise AssertionError(
                    f"{origin}: {label} drifted from the code: manifest {value} != computed {expected[label]}")

        if manifest["params"] != dict(conditions.FROZEN_PARAMS):
            raise AssertionError(
                f"{origin}: frozen params drifted: {manifest['params']} != {conditions.FROZEN_PARAMS}")

        plan = manifest["plan"]
        need_digest(plan["sha256"], f"{origin}:plan.sha256")
        plan_bytes = _Path(plan["path"]).read_bytes()
        if len(plan_bytes) != plan["bytes"]:
            raise AssertionError(f"{origin}: plan byte count drifted: {len(plan_bytes)} != {plan['bytes']}")
        if _hashlib.sha256(plan_bytes).hexdigest() != plan["sha256"]:
            raise AssertionError(f"{origin}: plan digest drifted from the manifest")

        slots = [m["slot"] for m in manifest["models"]]
        if slots != ["S1-A", "S1-B", "S1-C"]:
            raise AssertionError(f"{origin}: baseline must name three slots, got {slots}")
        for model in manifest["models"]:
            need_digest(model["sha256"], f"{origin}:{model['slot']}.sha256")
            if not isinstance(model["bytes"], int) or model["bytes"] <= 0:
                raise AssertionError(f"{origin}:{model['slot']}: byte count must be a positive integer")
        if len({m["path"] for m in manifest["models"]}) != 3:
            raise AssertionError(f"{origin}: baseline slots must name three distinct files")

        if manifest["qualification"]["generations"] != 0:
            raise AssertionError(f"{origin}: qualification must not generate anything")

    paths = sorted(_Path("eval").glob("s1-readiness-baseline*.json"))
    if not paths:
        raise AssertionError("no eval/s1-readiness-baseline*.json manifest on disk")
    def version_tuple(value):
        """Same ordering preflight uses: missing version ranks below every numbered one."""
        try:
            return tuple(int(x) for x in str(value).split("."))
        except (TypeError, ValueError):
            return (0,)

    by_version, by_id = {}, {}
    for path in paths:
        manifest = _json.loads(path.read_text(encoding="utf-8"))
        check_manifest(manifest, path.name)
        by_version[version_tuple(manifest.get("version", "0"))] = manifest
        by_id[manifest["baseline_id"]] = manifest

    # The guard must cover whatever preflight would actually execute. Resolving it through
    # preflight means a future v1.3 is covered by adding the file, not by editing this check.
    active = _json.loads(_Path(_preflight.manifest_path()).read_text(encoding="utf-8"))
    if version_tuple(active.get("version", "0")) != max(by_version):
        raise AssertionError(
            f"preflight does not select the newest manifest: {active.get('version')!r} vs {max(by_version)}")

    # A supersession must be auditable: every predecessor stays on disk, and a re-baseline
    # names both the value it replaced and the value it installed.
    node, seen = active, set()
    while "supersedes" in node:
        previous = by_id.get(node["supersedes"])
        if previous is None:
            raise AssertionError(
                f"{node['baseline_id']} supersedes {node['supersedes']}, which is not on disk")
        if previous["baseline_id"] in seen:
            raise AssertionError("supersession chain is cyclic")
        seen.add(previous["baseline_id"])
        was = {m["slot"]: m["sha256"] for m in previous["models"]}
        now = {m["slot"]: m["sha256"] for m in node["models"]}
        changed = sorted(slot for slot in now if was.get(slot) != now[slot])
        change = node.get("rebaseline")
        if changed and not change:
            raise AssertionError(
                f"{node['baseline_id']} changes {changed} but records no rebaseline reason")
        if change:
            before = [m for m in previous["models"] if m["slot"] == change["slot"]]
            after = [m for m in node["models"] if m["slot"] == change["slot"]]
            if not before or not after:
                raise AssertionError(f"{node['baseline_id']}: rebaseline names unknown slot {change['slot']}")
            need_digest(change["from"], f"{node['baseline_id']}:rebaseline.from")
            need_digest(change["to"], f"{node['baseline_id']}:rebaseline.to")
            if change["from"] == change["to"]:
                raise AssertionError("rebaseline replaces a digest with itself")
            if change["from"] != before[0]["sha256"]:
                raise AssertionError(
                    f"rebaseline.from {change['from']} is not the value it claims to replace")
            if change["to"] != after[0]["sha256"]:
                raise AssertionError(f"rebaseline.to {change['to']} is not the digest now installed")
            if change.get("generations_at_rebaseline") != 0:
                raise AssertionError("a baseline may not be re-based after generations")
            if changed != [change["slot"]]:
                raise AssertionError(
                    f"rebaseline covers {change['slot']} but the changed slots are {changed}")
        node = previous


@check("S1: qualification has no prompt path at all")
def _c89():
    """AST, not text: a docstring may *name* generate() without being able to call it."""
    import ast as _ast
    import inspect as _inspect
    from pathlib import Path as _Path

    from . import qualify as qmod

    source = _Path(_inspect.getsourcefile(qmod)).read_text(encoding="utf-8")
    tree = _ast.parse(source)

    for node in _ast.walk(tree):
        if isinstance(node, _ast.Call):
            func = node.func
            name = func.id if isinstance(func, _ast.Name) else getattr(func, "attr", None)
            if name in ("generate", "dispatch", "chat_completion"):
                raise AssertionError(f"qualification calls {name}(): it must be unable to prompt")
        if isinstance(node, _ast.Name) and node.id == "Dispatcher":
            raise AssertionError("qualification references Dispatcher")
        if isinstance(node, _ast.Attribute) and node.attr in ("generate", "dispatch"):
            raise AssertionError(f"qualification references .{node.attr}")

    if qmod.GENERATIONS != 0:
        raise AssertionError("qualification must declare zero generations")
    if set(qmod.ALLOWED_REQUESTS) != set(discovery.IDENTITY_PATHS):
        raise AssertionError("qualification may only issue the three identity GETs")
    if not hasattr(qmod, "write_qualification"):
        raise AssertionError("qualification writer missing")


@contextlib.contextmanager
def no_network():
    real_socket = socket.socket
    real_create = socket.create_connection

    class Blocked(real_socket):
        def __init__(self, *args, **kwargs):
            raise AssertionError("network access attempted during P1 validation")

    def blocked_create(*args, **kwargs):
        raise AssertionError("network access attempted during P1 validation")

    socket.socket = Blocked
    socket.create_connection = blocked_create
    try:
        yield
    finally:
        socket.socket = real_socket
        socket.create_connection = real_create


def main() -> int:
    print(f"MEMORY-09 harness {HARNESS_VERSION} — P1 validation (offline)")
    print(f"grant status: {authorization.grant_status()}")
    print(f"caps: {plan.caps_summary()}")
    print()
    with no_network():
        for name, fn in CHECKS:
            try:
                fn()
                RESULTS.append((name, True, ""))
            except AssertionError as exc:
                RESULTS.append((name, False, str(exc)))
            except Exception as exc:  # a crashing check is a failing check
                RESULTS.append((name, False, f"{type(exc).__name__}: {exc}"))
    passed = sum(1 for _, ok, _ in RESULTS if ok)
    for i, (name, ok, detail) in enumerate(RESULTS, 1):
        print(f"{'PASS' if ok else 'FAIL'} {i:2d}/{len(RESULTS)}  {name}")
        if not ok:
            print(f"           -> {detail}")
    print()
    print(f"RESULT: {passed}/{len(RESULTS)} passed, {len(RESULTS) - passed} failed")
    print(f"socket guard: no connection attempt observed across {len(RESULTS)} checks")
    return 0 if passed == len(RESULTS) else 1


if __name__ == "__main__":
    raise SystemExit(main())
