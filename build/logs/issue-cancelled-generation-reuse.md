# Issue: a cancelled generation poisons the next request on the same engine

Date: 2026-11-02
Status: **reproduced once, root cause not proven**. Separate from Stage 3.
Product code changed: none.

## Symptom

On an sd-server instance that had already served a **cancelled** generation, the
next Z-Image request returned HTTP 500 from the server:

```
[I] decode_first_stage completed, taking 13.47s --- image.cpp:675
[I] generate_image completed in 230.82s --- image.cpp:1140
{"error":"generate_image returned no results"}
```

The engine did the work — it generated and decoded — then answered with no
results.

## Fresh-engine comparison (same model, same parameters, same machine)

```
restart both binaries -> generate z-image-turbo-Q4_K_M, 8 steps, cfg 1, 1024x1024
[I] decode_first_stage completed, taking 13.52s
[I] generate_image completed in 94.55s
images/gallery/img_20261102_000805_600.png   791,480 bytes   PNG 1024x1024
```

Success. So the model, prompt, parameters, binary and machine are all fine; the
difference is the engine process's prior history.

## Exact sequence that preceded the failure

```
17:51:xx  Z-Image generation running; user pressed Abort in the UI
          -> CancelImageJob ok 0s
17:56:xx  app + harness restarted (new engine, new sd-server)
17:57:33  StartImageJob ok 0s  (z-image)
17:57:33  sd-server pid 9908 created
17:57:34  load complete: "total params memory size = 12619.23MB", "listening on"
17:57:39  CancelImageJob ok 0s   (cancel landed 5s AFTER the load finished)
18:01:30  next generate issued on the same pid 9908
18:05:50  UI shows: 500 generate_image returned no results
          (server log: generate_image completed in 230.82s, decode_first_stage ok)
```

## Why the cancel did not actually stop anything (source-derived)

`client.txt2img` is the synchronous A1111-style request. `client.cancelJob("")`
under that path:

```go
if jobID == "" {
    jobID = c.activeJobID   // empty for synchronous txt2img
}
c.cancelled = true
...
if jobID == "" {
    return nil              // no HTTP request is sent at all
}
```

So in synchronous mode a cancel is **local bookkeeping only**: no interrupt is
sent to sd-server, and the in-flight generation keeps running to completion on
the server. The `cancelled` flag is consumed by `pollJob` (the async path) and by
`generateSequentialBatch`, not by the synchronous `txt2img` request.

## Hypothesis (unverified)

Two generations overlap on one sd-server: the abandoned one keeps running, the
next request starts on the same process, and sd.cpp returns an empty result set
to one of them. The 230.82 s duration of the request that failed matches the
cancelled job's request still occupying the server, which is consistent with the
overlap — but the timing is inferred from log timestamps, not traced.

Alternative hypothesis: sd.cpp's own state after an abandoned request (not our
client) produces the empty result. Only a controlled experiment distinguishes
them.

## Suggested repair stage (not started, needs its own authorization)

1. Give the synchronous generation request a cancellable context, so a cancel
   closes the HTTP connection instead of leaving the server working.
2. Find out what sd-server actually offers for interruption: probe
   `/sdapi/v1/interrupt`, `/sdcpp/v1/jobs/{id}` DELETE, and whether the sync
   endpoint honours either. Do not assume a cancel endpoint exists.
3. If no interrupt exists, the engine must refuse to start a new generation while
   an abandoned one is still running, or restart the child before reusing it.
4. Acceptance: cancel a generation, immediately start another on the same engine,
   and get an artifact — no 500.

## Why this is separate from Stage 3

Stage 3 moved loading into the job and made cancellation reach a *load*. It did
not change how a cancel reaches a *generation request*, and the failing sequence
above involves a generation, not a load. The defect exists on the synchronous
`txt2img` path regardless of the loading lifecycle change.
