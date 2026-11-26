# M4 — Observation precedence over metadata estimates

Date: 2026-11-26
Isolated root: `C:\HOWL-Test\hra-image-20261125` (`.howl-isolated-root`)

M3 established that Howl's image projection is 586 MB below the device's own
measurement. M4 makes that measurement decide the next admission, without
turning it into a new estimate and without erasing the evidence that produced it.

## What a number covers

A residency figure is only meaningful with the question it answers attached.
These four are related and are not interchangeable:

| Category | What it is | What it protects against | Measured today? |
|---|---|---|---|
| Persistent model residency | Weights and modules the engine keeps between jobs | Running out of VRAM while a model is merely loaded | Engine placement for llama.cpp; not reported by the pinned sd-server |
| Peak generation allocation | Highest simultaneous use during one job | OOM during a generation, not before it | Yes, as the device delta across a load plus generation |
| Temporary execution buffers | Decode/compute workspace allocated per stage | A stage of a different size than the one measured | Only through the workload shape attached to the peak |
| Allocator-retained memory | Freed-but-not-returned blocks the process still holds | The next process finding less free memory than expected | Included in the delta; not separable on this build |

Howl's `5690 MB` image figure is a **post-generation device delta**: it covers
persistent weights, that generation's activation, its temporary buffers and
whatever the allocator kept. It is an upper bound on the peak of that workload
shape. It is deliberately *not* written into the resident requirement, because
charging a peak as residency would double-count activation on the next job.

## The precedence rule

`ResidencyObservation` now carries `Coverage` and the `WorkloadShape` (width,
height, batch, weight precision) that produced the delta. A measurement is
reused only when it covers the request: equal or larger in every dimension and
of the same precision. An observation with no measured shape is never reused,
and a q8_0 measurement does not bound an f16 request.

Admission gained `MinimumCommitmentVRAMMB`, fed from the recorded observation
for the same model and a covering shape. It can only raise the commitment:

- the projection stays exactly what it was (`ProjectedPeakVRAMMB` unchanged);
- `CommittedVRAMMB()` returns the larger of projection, required and floor;
- the grant and refusal reasons name the floor and its coverage;
- the resident requirement is untouched.

## Evidence (browser dev host, real sd-server, real generation)

Same session, Illustrious-XL-v2.0, 512×512, 8 steps, q8_0:

```
[Residency] kind=image model="Illustrious-XL-v2.0" planned=5093MB engine=0MB device_delta=5690MB ledger=5690MB above plan
```

Unload, then load the same model again. The f16 attempt first asks at the stored
precision and correctly ignores the q8_0 observation (`observed_floor=0MB`).
The q8_0 attempt:

```
[App] image request: ... measured_resident=0MB observed_floor=5690MB model="Illustrious-XL-v2.0"
[App] image admission final: admitted=true projected=4613MB resident=4101MB activation=0MB required=5084MB ...
      admitted: 5690 MB VRAM required, 6262 MB available after 819 MB reserve;
      measured floor 5690 MB (post-generation-device-delta) raised the commitment above the 4613 MB projection
```

The second admission commits 5690 MB where the first committed 5084 MB. That is
the M3 disagreement changing a decision, not decorating a log.

## Owners

- **Resident ledger**: `resource.Manager` (`RegisterModel`) stays the single
  owner. The measurement path rewrites entries only through
  `reconcileRegisteredResidency`, and only with an observed value.
- **Admission**: `resource.Authority` stays the single decision owner. It
  receives a floor; it decides.
- **Observation**: `App.residency` is the only recorder. It is read back through
  `GetResidencyReport` and consulted only by `imageAdmissionRequest`.
- **Model intelligence**: `modelintel` profiles remain metadata/artifact
  classification. They are not modified by measurement and do not carry
  residency claims; that separation is what lets a measurement outrank them
  without silently editing the catalog.

No competing writers were added. No reserve was consolidated: R1 admission, R2
allocator, R3 image reclaim and R4 `--max-vram` remain four computations.

## Verified limitations

- **Packaged-host image UI still unexercised.** The packaged image panel exposes
  no accessible nav targets, so this evidence again comes from the browser dev
  host driving the same `App` bindings; the packaged image *workflow* was not
  clicked through in M4.
- **No per-category breakdown.** This sd-server build reports no components, so
  Howl cannot yet say how much of the 5690 MB is weights, activation, buffers or
  allocator retention. The categories are recorded as distinct concepts; only
  the aggregate is measured.
- **Chat does not consume the floor.** Its engine reports placement and its plan
  already exceeds the device delta, so no floor is applied there yet.
- **Observations are per session.** Nothing is persisted or calibrated; a
  restart begins with the projection again. Persistent calibration remains
  prohibited by the phase-one plan.
- **Shape reuse is a bound, not a model.** A larger measured shape is assumed to
  bound a smaller request; that assumption is conservative but unmeasured.
