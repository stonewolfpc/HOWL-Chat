# Universal Audio Architecture for HOWL Chat

## Executive Summary

A staged, consumer-ready audio system supporting speech-to-text (ASR), text-to-speech (TTS), and multimodal integration. Built for universal deployment across hardware constraints.

## Architecture Principles

### 1. Surgical Precision (SWEObeyMe Compliance)
- **500-line max** per implementation file
- **Single responsibility** per module
- **Zero digital debt** - no deprecated code paths
- **Interface-driven** design for swappable components

### 2. Universal Compatibility
- **CPU-first** design (GPU acceleration optional)
- **Memory tiers**: Tiny (<2GB), Small (4-8GB), Large (16GB+)
- **Cross-platform**: Windows, Linux, macOS
- **Model agnostic**: Works with Whisper, Qwen-ASR, Qwen-Omni, etc.

### 3. Consumer-Grade Reliability
- **Fail-soft**: Degrades gracefully on errors
- **Progressive enhancement**: Core features work everywhere
- **Observable**: Full logging and metrics
- **Configurable**: Runtime feature toggles

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    AUDIO SUBSYSTEM                           │
├─────────────────────────────────────────────────────────────┤
│  STAGE 1: Core Types & Interfaces (Foundation)              │
│  ├── audio/types/audio_types.go              (interfaces)   │
│  ├── audio/types/model_profiles.go           (config)       │
│  └── audio/types/audio_format.go             (utilities)  │
├─────────────────────────────────────────────────────────────┤
│  STAGE 2: ASR - Speech-to-Text                              │
│  ├── audio/asr/recognizer.go                 (interface)    │
│  ├── audio/asr/whisper/whisper_asr.go        (impl)         │
│  ├── audio/asr/qwen-asr/qwen_asr.go          (impl)         │
│  └── audio/asr/fallback/local_whisper.go   (fallback)     │
├─────────────────────────────────────────────────────────────┤
│  STAGE 3: TTS - Text-to-Speech                              │
│  ├── audio/tts/synthesizer.go                (interface)  │
│  ├── audio/tts/piper/piper_tts.go              (impl)       │
│  ├── audio/tts/melotts/melo_tts.go             (impl)       │
│  └── audio/tts/cache/tts_cache.go            (optional)   │
├─────────────────────────────────────────────────────────────┤
│  STAGE 4: Multimodal Integration                           │
│  ├── audio/multimodal/omni_bridge.go         (orchestrator)│
│  ├── audio/multimodal/audio_vision.go        (combined)   │
│  └── audio/multimodal/streaming_chunker.go   (streaming)  │
├─────────────────────────────────────────────────────────────┤
│  STAGE 5: Real-time Pipeline (Advanced)                     │
│  ├── audio/realtime/vad.go                     (detection)│
│  ├── audio/realtime/buffer.go                  (streaming)│
│  └── audio/realtime/webrtc_bridge.go           (network)  │
└─────────────────────────────────────────────────────────────┘
```

## Staged Implementation Roadmap

### Stage 1: Foundation (Week 1)
**Goal**: Define contracts and abstractions

**Deliverables**:
1. `internal/audio/types/interfaces.go` - Core abstractions
2. `internal/audio/types/models.go` - Model profiles and configs
3. `internal/audio/types/formats.go` - Audio format utilities

**Success Criteria**:
- [ ] All interfaces defined with clear contracts
- [ ] Model profiles for supported ASR/TTS models
- [ ] Audio format conversion utilities
- [ ] Zero dependencies on external libraries

### Stage 2: ASR Implementation (Week 2)
**Goal**: File-based speech-to-text

**Deliverables**:
1. `internal/audio/asr/interface.go` - ASR recognizer interface
2. `internal/audio/asr/whisper/local.go` - Whisper.cpp integration
3. `internal/audio/asr/qwen-asr/client.go` - Qwen3-ASR HTTP client
4. `internal/audio/asr/fallback/chain.go` - Fallback orchestrator

**Success Criteria**:
- [ ] Load audio file (WAV, MP3, OGG) → Text
- [ ] Support Whisper (local) and Qwen-ASR (remote)
- [ ] Automatic fallback on failures
- [ ] Progress callbacks for UI

### Stage 3: TTS Implementation (Week 3)
**Goal**: Text-to-speech with caching

**Deliverables**:
1. `internal/audio/tts/interface.go` - TTS synthesizer interface
2. `internal/audio/tts/piper/local.go` - Piper TTS integration
3. `internal/audio/tts/melotts/client.go` - MeloTTS integration
4. `internal/audio/tts/cache/manager.go` - Audio caching layer

**Success Criteria**:
- [ ] Text → WAV/MP3 audio file
- [ ] Multiple voice options
- [ ] Response caching for repeated phrases
- [ ] Streaming playback support

### Stage 4: Multimodal Bridge (Week 4)
**Goal**: Combine audio + vision + text

**Deliverables**:
1. `internal/audio/multimodal/bridge.go` - Omni model orchestrator
2. `internal/audio/multimodal/audio_vision.go` - AV combined input
3. `internal/audio/multimodal/context_manager.go` - Multi-modal context

**Success Criteria**:
- [ ] Audio input + Image → LLM response
- [ ] Audio input + Text → LLM response
- [ ] Unified context across modalities
- [ ] Token budget management across audio+vision

### Stage 5: Real-time Pipeline (Week 5-6)
**Goal**: Streaming voice conversation

**Deliverables**:
1. `internal/audio/realtime/vad.go` - Voice Activity Detection
2. `internal/audio/realtime/buffer.go` - Circular audio buffer
3. `internal/audio/realtime/streamer.go` - Chunked streaming
4. `internal/audio/realtime/session.go` - Conversation state machine

**Success Criteria**:
- [ ] Real-time microphone input
- [ ] Interruptible responses
- [ ] <500ms latency end-to-end
- [ ] Natural turn-taking

## Technical Specifications

### Audio Format Standards
```go
// Supported input formats
InputFormats = []string{"wav", "mp3", "ogg", "flac", "m4a"}

// Internal processing format (always convert to this)
ProcessingFormat = "wav"  // 16-bit PCM, 16kHz mono

// Output formats
OutputFormats = []string{"wav", "mp3", "ogg"}
```

### Model Profiles
```go
// ASR Models
type ASRModel string
const (
    ASRWhisperTiny   ASRModel = "whisper-tiny"    // ~39MB, CPU realtime
    ASRWhisperBase   ASRModel = "whisper-base"    // ~74MB, good quality
    ASRWhisperSmall  ASRModel = "whisper-small"   // ~244MB, best quality
    ASRQwen3ASR06B   ASRModel = "qwen3-asr-0.6b"  // ~600MB, multilingual
    ASRQwen3ASR17B   ASRModel = "qwen3-asr-1.7b"  // ~1.7GB, best multilingual
)

// TTS Models
type TTSModel string
const (
    TTSPiper      TTSModel = "piper"       // Fast, small, multiple voices
    TTSMeloTTS    TTSModel = "melotts"     // Better quality, slower
    TTSBark       TTSModel = "bark"        // Most natural, very slow
)
```

### Memory Budgets
```go
type MemoryTier string
const (
    TierTiny   MemoryTier = "tiny"    // <2GB VRAM/RAM
    TierSmall  MemoryTier = "small"   // 4-8GB
    TierLarge  MemoryTier = "large"   // 16GB+
)

// Recommended configurations per tier
var RecommendedConfigs = map[MemoryTier]AudioConfig{
    TierTiny: {
        ASRModel:  ASRWhisperTiny,
        TTSModel:  TTSPiper,
        EnableCache: true,
        MaxAudioLength: 30, // seconds
    },
    TierSmall: {
        ASRModel:  ASRWhisperBase,
        TTSModel:  TTSPiper,
        EnableCache: true,
        MaxAudioLength: 120,
    },
    TierLarge: {
        ASRModel:  ASRQwen3ASR17B,
        TTSModel:  TTSMeloTTS,
        EnableCache: true,
        MaxAudioLength: 300,
    },
}
```

## Integration Points

### 1. Wails Frontend Integration
```javascript
// Frontend will expose:
window.runtime.EventsOn('audio:transcription', (text) => { ... })
window.runtime.EventsOn('audio:synthesis-complete', (path) => { ... })
window.runtime.EventsEmit('audio:request-transcription', audioData)
window.runtime.EventsEmit('audio:request-synthesis', text)
```

### 2. LLM Backend Integration
```go
// Audio output feeds into existing chat pipeline
type AudioMessage struct {
    Type string // "audio_input" or "audio_output"
    Text string // Transcribed or to-be-spoken
    AudioPath string // File path to audio
    Duration float64 // Audio duration in seconds
}

// Integration with existing chat system
func (a *App) ProcessAudioMessage(msg AudioMessage) (string, error) {
    // 1. If audio_input: transcribe → send to LLM
    // 2. If LLM response: optionally synthesize → play
}
```

### 3. Model Management Integration
```go
// Reuse existing model loading infrastructure
func (m *ModelManager) LoadAudioModel(modelType string, profile ModelProfile) error {
    // Leverage existing download, caching, and lifecycle management
}
```

## Error Handling Strategy

### Graceful Degradation Chain
```
1. Try primary ASR (Qwen3-ASR)
   ↓ (if fails)
2. Fallback to Whisper local
   ↓ (if fails)
3. Fallback to Whisper tiny (always works)
   ↓ (if all fail)
4. Show error + manual text input
```

### TTS Fallback Chain
```
1. Try requested TTS model
   ↓ (if fails)
2. Fallback to Piper (lightweight)
   ↓ (if fails)
3. Disable audio output (text only)
```

## Configuration Schema

```yaml
# config/audio.yaml
audio:
  enabled: true
  
  asr:
    primary_model: "qwen3-asr-1.7b"
    fallback_model: "whisper-base"
    language: "auto"  # auto-detect or "en", "zh", etc.
    
  tts:
    enabled: true
    model: "piper"
    voice: "en_US-lessac-medium"
    speed: 1.0
    cache_enabled: true
    cache_size_mb: 100
    
  multimodal:
    enabled: false  # Stage 4
    omni_model: "qwen2.5-omni-7b"
    
  realtime:
    enabled: false  # Stage 5
    vad_threshold: 0.5
    buffer_duration_ms: 300
    
  hardware:
    memory_tier: "auto"  # auto, tiny, small, large
    prefer_gpu: true
```

## Testing Strategy

### Unit Tests (per module)
- Mock audio files for deterministic testing
- Interface compliance tests
- Error injection tests

### Integration Tests
- End-to-end audio pipeline tests
- Memory pressure tests
- Fallback chain verification

### Consumer Validation
- Test on: Low-end laptop, mid-range desktop, high-end workstation
- Validate on: Windows 10/11, Ubuntu 22.04, macOS 14

## Documentation Requirements

Each stage must include:
1. **README.md** - Module overview and usage
2. **API.md** - Interface documentation
3. **CONFIG.md** - Configuration options
4. **TROUBLESHOOTING.md** - Common issues and solutions

## Success Metrics

- [ ] Stage 1: All interfaces compile, tests pass
- [ ] Stage 2: ASR works on 3 different audio files
- [ ] Stage 3: TTS generates playable audio
- [ ] Stage 4: Audio + Image → Text response
- [ ] Stage 5: Real-time conversation <500ms latency

## Next Steps

1. **Review and approve** architecture document
2. **Begin Stage 1** implementation (Core Types)
3. **Set up CI/CD** for audio subsystem testing
4. **Create feature flags** for gradual rollout

---

**Author**: Senior Engineer  
**Date**: May 2026  
**Status**: Architecture Review  
**Compliance**: SWEObeyMe Surgical Standards
