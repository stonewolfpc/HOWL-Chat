# `MEMORY-09` — Inference Route Policy

> **The objective is not to remember everything. It is to remember the right things, at the right times, at the right places.**

**Identity:** `MEMORY-09-IRP-1` · **Version:** 1.1 · **Date:** 2026-11-13 · **Status:** **POLICY — enforced in code; no call made under it.**
**Revision 1.1:** the owner confirmed there is **no payment method**, so paid inference is *unreachable* rather than prohibited pending authorisation. The attestation requirement in §3 is superseded by §7; the rest stands.
**Authority:** owner direction, 2026-11-13: **maximum authorised paid spend $0.00.**
**Implementation:** `eval/harness/budget.py`, `eval/harness/local_adapter.py`; checks 26–35.

---

## 1. The rule

> **No paid inference. No automatic credit purchases. No surprise charges.** Hugging Face free credits are an optional testing route, not a dependency maintained by payment. If the free allowance is insufficient, the route changes to **local**, not to the wallet.

## 2. Route table

| Route | Permission |
|---|---|
| Hugging Face free credits | **Allowed, only when confirmed free** |
| Paid Hugging Face inference | *(blank)* — **no permission** |
| Automatic paid fallback | *(blank)* — **no permission** |
| Local SillyTavern / KoboldCpp | **Preferred fallback** |
| Other genuinely free inference | Eligible **after verification** |

The two blank cells are read as **no permission**, not "unspecified". The policy table encodes the strings `NO PERMISSION` for both, and `budget.check_provider` has no code path that returns a paid route.

## 3. Four hard requirements, and where each is enforced

| # | Requirement | Enforcement |
|---|---|---|
| 1 | Maximum authorised paid spend: **$0.00** | `MAX_PAID_SPEND_USD = 0.00`; `check_spend(0.01)` refuses (check 26) |
| 2 | Never enable automatic top-up or paid-credit consumption | `AUTO_TOPUP_PERMITTED = False`; attestation with `auto_topup_enabled` refuses (check 29) |
| 3 | If free-credit eligibility or balance **cannot be verified**, refuse the call | no attestation → refuse (check 27); paid tier → refuse (check 28) |
| 4 | Stop when free credits are exhausted; **do not silently switch providers** | `remaining_units <= 0` → refuse with "stop, do not switch" (check 30); any route change → `RouteChangeRefused` (check 31) |

**On requirement 3, stated plainly:** the harness cannot read a provider's balance through a chat-completions endpoint. So "unverified" is the *default* state, and the default is refusal. A hosted call becomes possible only with an owner attestation naming the route, the tier, the remaining balance, the verification date and the source. No attestation exists, so **the hosted route is currently refused in code**, which is the intended state, not a gap.

## 4. Why local is a fallback and not a provider

SillyTavern and KoboldCpp are **interfaces and runtimes**, not hosted model services. They serve a model running on this machine: no per-call charge, no balance to verify, and no vendor between the harness and the weights.

The trade is explicit and returns the constraint that hosted inference had deferred:

| Gain | Cost |
|---|---|
| No spend, no credit dependency, no vendor availability risk | Hardware becomes binding again |
| Full control of quantisation and sampler values | Throughput and context limits are the machine's |
| Offline reproducibility | The 40B-class question moves from "S3 someday" to "the live question" |

That is why S3 exists. Choosing local is choosing to measure the hardware.

## 5. Replaceability — the property this policy depends on

The harness tests a **memory architecture**, not a billing relationship. So the adapter is replaceable and the tests are not rewritten:

- One request type (`GenerateRequest`), one response type, one provenance schema.
- Route chosen by provider string: `local:<runtime>` selects the loopback adapter; anything else must clear the budget policy first.
- **Loopback only.** `localhost`, `127.0.0.1`, `::1`. A non-loopback host is refused, because a remote host wearing a local name is a hosted provider and must face the budget policy (check 32).
- **No endpoint is guessed.** A local route without an explicit endpoint is an error, not a default port.
- **No fallback.** A failed call stays on its route (row 477 F-1); the factory never substitutes.

## 6. What this policy does not establish

- That any hosted free allowance exists, is sufficient, or will persist.
- That local inference is adequate for the tests — that is S3's question.
- That a free route is safe to use indefinitely; a free tier is permission to *try*, verified before each series.
- Anything about model behaviour. This is a spending and traceability rule, nothing more.

## 7. Revision 1.1 — the boundary got simpler because there is nothing to charge

> *"I won't be putting my card info in."* — owner, 2026-11-13

With no payment method on the account, a paid call cannot complete. That changes the character of the rule rather than its direction:

| Before | After |
|---|---|
| Paid routes **prohibited**; required an attestation proving a route was free | Paid routes **unreachable**; nothing can incur a charge |
| Free-credit balance had to be verified before a call | The provider's own response is the signal: success was free, a demand is a stop |
| Grant carried a cost ceiling that had to be $0.00 | The ceiling remains, but nothing hinges on it |

**The three-line rule now in force:**

1. Hosted is usable **only while inside the free allowance**.
2. If it asks for payment — HTTP 402, "insufficient credits", "quota exceeded", "add a payment method" — **stop**. No retry, no top-up, no automatic switch. `budget.PaymentDemanded` is raised at the dispatch boundary, and the call still consumes a ledger slot.
3. The fallback is **local**, chosen by a person, not by the harness.

**What was removed, deliberately:** the `FreeCreditAttestation` type, the attestation file requirement, balance verification, and the spend-authorisation flow. They existed to make a *paid* route safe. There is no paid route, so they were machinery guarding a door that has no key.

**What was kept:** the `$0.00` ceiling, the prohibition on automatic fallback, route-change refusal, and `cost.py` — which is now belt-and-braces, kept because it costs nothing and was an inspection prerequisite, not because anything depends on it.

**What this simplification costs:** if a payment method is ever added, `PAYMENT_METHOD_ON_FILE` must change to `True`, and every hosted route is then refused until the owner restates the rule. The harness does not assume that adding a card makes spending acceptable.

## 8. Related

Row 477 (adapter interface, provider pinning, fail-closed) · row 479 (stage ladder, S3 = local) · row 482 (P1/P2 split). Registry row 484.
