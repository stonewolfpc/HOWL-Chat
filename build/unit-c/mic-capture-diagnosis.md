# Microphone capture diagnosis — evidence record

Date: 2026-11-05. Scope: authorized diagnosis of the post-D3 Unit C empty
transcript. Status: **DIAGNOSIS COMPLETE — ENVIRONMENTAL CAUSE, NO CODE
CHANGE.** Unit C remains open; chat-model work remains separate.

No application code, frontend asset, settings, menu, D1 or D3 record was
modified. All tests ran from the isolated root described in section 6.

## 1. Capture path (source trace)

```
microphone
  └─ Chromium/WebView2 default capture endpoint (getUserMedia, no deviceId)
     frontend/dist/js/audio-service.js:175
       constraints: channelCount 1, echoCancellation, noiseSuppression,
                    autoGainControl = true        (no explicit device)
  └─ AudioContext (requested 16 kHz, else device default)   :190
  └─ createMediaStreamSource + createScriptProcessor(4096,1,1)  :196-197
       Float32 chunks pushed while _recording                :203-206
  └─ stopRecording(): mergeSamples                          :250
       hasSignal(): RMS > 0.002  ← only silence gate        :332-338
  └─ AudioWav.encodeWav16(): resample to 16 kHz, PCM16 mono  audio-wav.js
  └─ bytesToDataURL(): FileReader → data:audio/wav;base64
  └─ app.TranscribeAudioData(dataURL)                       :350
     app-audio.go:89
       base64 decode, RIFF/WAVE check, 64 MiB cap
       temp file  <TMP>/howl-asr/recording-*.wav
       manager operation + lease → Pipeline.ProcessAudioMessage
  └─ llama.Recognizer.Transcribe  internal/audio/asr/llama/llama.go
       multipart POST /v1/audio/transcriptions
       fields: file=recording-*.wav, model=qwen3-asr, response_format=json
  └─ JSON {"text": ...} → CleanTranscript (strips "language X<asr_text>")
       empty text → ErrCodeFormatNotSupported
                    "llama-server returned an empty transcript"
       → app-audio.go logs ASR-TRANS-FAIL (high)
  └─ renderer inserts non-empty text into #message-input; sending stays manual
```

Two consequences matter for the earlier failure:

1. The renderer calls the backend **only after** `hasSignal` passes, so the
   previous run's payload was not digital silence. It was non-silent audio.
2. The backend has one error for "recognizer returned no text", which is
   logged at high severity as a transcription failure. Noise-only audio
   therefore surfaces exactly like an ASR fault.

The application never enumerates or pins a capture device; it uses the
WebView2/Chromium default (communications role). That is why the earlier run
recorded whichever endpoint Windows had defaulted, regardless of whether the
headset was being worn.

## 2. Capture endpoints (observed)

Default at test time, console and communications roles identical:

```
capture/console : Headset (WH-1000XM4)  state=1 (active)
                  {0.0.1.00000000}.{6b9264d1-a6cc-48b2-84cd-082bacac1599}
capture/comms   : same endpoint
render/console  : Headphones (WH-1000XM4) state=1
                  {0.0.0.00000000}.{9460d508-20ce-4ade-8b76-e405f9ab758a}
render/comms    : same endpoint
```

- The endpoint is Bluetooth HFP (`BTHHFENUM\BTHHFPAudio`,
  `WH-1000XM4 Hands-Free`), 1 channel, 16-bit.
- Other capture endpoints were absent, disabled, or unplugged: Corsair VOID
  mic `state=4`, Voicemod `state=2`, Realtek mic `state=4`/`state=8`.
- Microphone consent (HKCU ConsentStore\microphone) = `Allow`.
- Application-side endpoint is **inferred**, not directly logged: the code
  supplies no `deviceId`, and Windows communications default equals the
  headset above. There is no app-side device enumeration to quote.

## 3. Controlled test

Condition: operator wearing the WH-1000XM4 with its boom microphone, quiet
room, headphones as render endpoint, isolated app instance. Cue: audible
beeps from the host; recording state confirmed in the UI before the cue.

| Run | Result | Notes |
|---|---|---|
| A (worn) | **Transcribed** | transcript reached composer |
| B | invalid | mic toggle fired before the button existed (`index 0 out of range`); no recording, no transcript; discarded |
| C (worn) | **Transcribed + WAV captured** | definitive measurement run |

Run A transcript (composer, verbatim):
`The quick brown fox jumped over the whatever. Is it picking any of this up? You said a telemessage that that that just anything that it picks up works, right?`

Run C transcript (composer, verbatim):
`Testing. Testing one, two, three. The quick brown fox jumped over the something or other, and this is the message to test it.`

Run C notification: `Transcription complete`.

### Run C capture measurements (exact payload sent to ASR)

```
file      mic3-recording-3767929088.wav
sha256    74bf6216c7350bc6d8d89c397f8e4ef59dce9d91585d939a1469f838568fd36e
bytes     434,220  (44-byte header + 217,088 samples)
format    16000 Hz, 1 channel, 16-bit PCM
duration  13.57 s
samples   217,088 total; 161,616 nonzero; 66,283 above 2% FS
peak      1.0000  (7 samples at full scale)
RMS       0.10658  (about -19.4 dBFS)
```

This is healthy speech-level audio: strong RMS, no digital silence, and only
incidental clipping from Bluetooth HFP with AGC enabled. It was recovered from
`<TEMP>/howl-asr/` while the transcription request was in flight, so it is the
same byte sequence the recognizer received.

ASR request/response evidence: the app's own POST returned the Run C transcript
above and the renderer inserted it. A separate replay of the captured file with
`curl` returned `HTTP=000` (no response). The cause was not established (the
request path/transport was not captured), so the replay is **INCONCLUSIVE**,
not a second data point. No new ASR error was logged in any worn-headset run.

## 4. Root cause

**Environmental, not an application capture defect.** The previous empty
transcript came from an unworn microphone: room noise passed the renderer's
weak RMS gate, was encoded and sent, and the recognizer correctly returned no
text. The headset was confirmed not worn by the operator.

Grading:

- **Observed:** with the headset worn, two recordings produced nonempty
  transcripts and composer insertion; Run C's WAV has 0.107 RMS, 1.0 peak,
  13.57 s at 16 kHz mono.
- **Observed:** Windows default capture and render endpoints were the headset
  pair, mic consent `Allow`, and the app entered recording state.
- **Source-derived:** the only silence gate is `RMS > 0.002` in
  `audio-service.js`; the backend maps empty recognizer text to a
  high-severity `ASR-TRANS-FAIL`. Noise-only input therefore reaches the ASR
  and is reported as an ASR failure.
- **Source-derived:** the app does not select or pin a capture device; it
  relies on the WebView2 default.
- **Inferred:** the earlier payload was noise-only room audio from the
  unworn headset. The WAV from that run was not retained, so this is inferred
  from the gate behavior plus the confirmed unworn condition.
- **Unestablished:** whether a stronger gate/VAD would materially improve
  robustness; no such change was attempted.

No defect in capture, resampling, WAV encoding, transport, or ASR handling was
demonstrated. **No repair was made and no file was changed.**

## 5. Verdicts

| Check | Verdict | Evidence |
|---|---|---|
| Read-only capture-path trace | **PASS** | source trace in section 1 |
| Exact Windows capture endpoint | **PASS** | default console/comms = WH-1000XM4 HFP |
| Application endpoint identity | **PASS (inferred)** | no `deviceId`; default comms = headset |
| Headset mic active, unmuted, consented | **PASS** | `state=1`, consent `Allow` |
| Worn-mic capture produces nonempty audio | **PASS** | Run C WAV, RMS 0.107 |
| Audio reaches ASR adapter | **PASS** | temp WAV intercepted during request; app transcript |
| Qwen3-ASR returns nonempty transcript | **PASS** | two transcripts, section 3 |
| Transcript inserted into editable composer | **PASS** | composer read via UIA; `Transcription complete` |
| Sending remains user-controlled | **NOT TESTED** | no send performed; source inserts text only |
| Ordinary typed chat unchanged | **NOT TESTED** (no code touched) | — |
| ASR startup/shutdown D3 ownership | **PASS** | every run: normal close, app + ASR child gone, `none` left |
| Replay of captured bytes | **INCONCLUSIVE** | `HTTP=000` after release |
| Protected data preserved | **PASS** | section 6 |

## 6. Isolation and protected data

Same discipline as the post-D3 rerun: copied accepted build in
`build/tmp/unit-c-rerun/`, junctioned `runtime/`, own `config/`, and per-process

```
USERPROFILE  ...\build\tmp\unit-c-rerun\profile
APPDATA      %USERPROFILE%\AppData\Roaming
LOCALAPPDATA %USERPROFILE%\AppData\Local
TEMP/TMP     %USERPROFILE%\Temp
```

The redirected TEMP is what made the exact WAV recoverable inside the isolated
root instead of the real user temp.

After testing:

- `config/settings.json` =
  `e786f53f2fd7f91b794936dad1930e2b09594de8a81167179a857ecb011757bf`
  (unchanged).
- Real `%APPDATA%\HOWLChat` file hashes: **IDENTICAL** to the pre-test
  baseline (`real-before.txt`).
- All launch/close cycles: `CLOSE=True`, `EXITED=True`, zero `howl-chat` or
  `llama-server` processes remaining.
- Evidence artifacts are under `build/unit-c/mic-diagnosis/`.

## 7. Commands (representative)

```
# endpoints
powershell -File build/unit-c/tools/AudioEndpoints.ps1

# isolated launch with redirected TEMP (recorded APP_PID)
Start-Process howl-chat-unit-c.exe -WorkingDirectory <isolated root>

# record / stop (UIA toggle; class changes to "mic-button recording")
uia.ps1 -Action invoke -Class "mic-button"

# exact-payload capture while the request is in flight
cp profile/Temp/howl-asr/recording-*.wav build/unit-c/mic-diagnosis/

# measurements
analyze.ps1 -Path <wav>     (RIFF fields, duration, peak, RMS, nonzero, clip)
```

## 8. Limitations and next decision

- The earlier failing WAV was not retained; the environmental cause is
  inferred from the confirmed unworn headset plus the source-derived gate, not
  from that file.
- Replay to the ASR server did not complete (`HTTP=000`, cause not
  established); only the app's own request/response is evidence for the
  transcript.
- Run B was invalid due to UIA timing and is excluded.
- No send, no typed-chat regression, and no chat LLM were exercised; those
  remain Unit C gaps.
- The weak silence gate and the generic high-severity "empty transcript" log
  are real robustness/observability weaknesses, but they did not cause this
  failure. Changing them needs separate authorization.

Next decision required: proceed to the separate chat-model blocker (as the
owner sequenced), and decide separately whether the silence-gate / error
classification weakness should ever be scheduled as its own scoped change.
