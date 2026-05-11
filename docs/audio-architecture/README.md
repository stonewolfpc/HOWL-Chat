# HOWL Chat - Universal Audio Architecture

## Executive Summary

A **consumer-ready, staged audio system** designed for universal deployment across hardware tiers. Built with surgical precision following SWEObeyMe governance.

**Status**: Stage 1 Complete (Foundation)  
**Compliance**: SWEObeyMe Surgical Standards  
**Constraint**: 500-line max per file  
**Architecture**: Interface-driven, separation of concerns

---

## Quick Start

### Hardware Requirements

| Tier | RAM | VRAM | ASR | TTS | Use Case |
|------|-----|------|-----|-----|----------|
| **Tiny** | 4GB | None | Whisper Tiny | Piper | Laptops, Raspberry Pi |
| **Small** | 8GB | Optional | Whisper Base | Piper | Standard desktops |
| **Medium** | 16GB | 4GB | Qwen3-ASR 1.7B | MeloTTS | Enthusiast builds |
| **Large** | 32GB+ | 8GB+ | Qwen-Omni | MeloTTS | Power users |

### Integration

```go
// Initialize audio pipeline
import "howl-chat/internal/audio/types"

// Detect hardware and recommend config
config := types.AutoDetectConfig()

// Create pipeline
pipeline, err := audio.NewPipeline(config)
if err != nil {
    log.Fatal(err)
}

// Transcribe audio
result, err := pipeline.ASR.Transcribe(ctx, "input.wav", nil)
fmt.Println(result.Text)

// Synthesize speech
synthesis, err := pipeline.TTS.Synthesize(ctx, types.SynthesisRequest{
    Text: "Hello, world!",
    Voice: "en_US-lessac-medium",
}, nil)
```

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     AUDIO SUBSYSTEM                          │
├─────────────────────────────────────────────────────────────┤
│  Stage 1: Core Types (COMPLETED)                            │
│  ├── types/interfaces.go       (427 lines)                  │
│  ├── types/model_profiles.go   (498 lines)                  │
│  └── types/utils.go            (385 lines)                  │
├─────────────────────────────────────────────────────────────┤
│  Stage 2: ASR - Speech-to-Text (PENDING)                    │
│  ├── asr/interface.go                                        │
│  ├── asr/whisper/local.go                                    │
│  ├── asr/qwen-asr/client.go                                  │
│  └── asr/fallback/chain.go                                   │
├─────────────────────────────────────────────────────────────┤
│  Stage 3: TTS - Text-to-Speech (PENDING)                    │
│  ├── tts/interface.go                                        │
│  ├── tts/piper/local.go                                      │
│  ├── tts/melotts/client.go                                   │
│  └── tts/cache/manager.go                                    │
├─────────────────────────────────────────────────────────────┤
│  Stage 4: Multimodal (PENDING)                             │
│  ├── multimodal/bridge.go                                    │
│  └── multimodal/context.go                                   │
├─────────────────────────────────────────────────────────────┤
│  Stage 5: Real-time (PENDING)                              │
│  ├── realtime/vad.go                                         │
│  └── realtime/streamer.go                                    │
└─────────────────────────────────────────────────────────────┘
```

---

## Key Design Decisions

### 1. Interface-Driven Architecture

Every component implements a minimal interface:

```go
// ASRRecognizer - 6 methods, clear contract
type ASRRecognizer interface {
    Initialize(ctx context.Context, config ASRConfig) error
    Transcribe(ctx context.Context, audioPath string, callback ProgressCallback) (*RecognitionResult, error)
    TranscribeStream(ctx context.Context, chunks <-chan AudioChunk, results chan<- RecognitionResult) error
    SupportsLanguage(lang string) bool
    GetModelInfo() ASRModelInfo
    Release() error
}
```

**Benefits**:
- Swappable implementations (Whisper, Qwen, future models)
- Easy mocking for testing
- Clear dependency boundaries

### 2. Graceful Degradation

```
User requests ASR with Qwen3-ASR-1.7B
         ↓
    Check hardware (need 2GB RAM)
         ↓
    If insufficient:
         ↓
    Fallback to Whisper Base (200MB)
         ↓
    If that fails:
         ↓
    Fallback to Whisper Tiny (100MB)
         ↓
    If all fail:
         ↓
    Return error + suggest manual text input
```

### 3. Zero Digital Debt

- **No deprecated code paths**
- **No TODO/FIXME comments** (track in issues)
- **Clean separation** - no circular dependencies
- **Explicit error handling** - no silent failures

### 4. Consumer-Grade Reliability

- **Fail-soft**: Errors don't crash the app
- **Progress reporting**: Users see what's happening
- **Cancel-friendly**: Operations respect context cancellation
- **Resource cleanup**: Always call Release()

---

## File Organization

```
internal/audio/
├── types/                    # Core interfaces (Stage 1)
│   ├── interfaces.go         # All interface definitions
│   ├── model_profiles.go     # Model metadata & registry
│   └── utils.go              # Audio format utilities
│
├── asr/                      # Speech-to-text (Stage 2)
│   ├── interface.go          # ASR recognizer interface
│   ├── whisper/              # Whisper.cpp integration
│   │   └── local.go
│   ├── qwen-asr/             # Qwen3-ASR HTTP client
│   │   └── client.go
│   └── fallback/             # Fallback orchestrator
│       └── chain.go
│
├── tts/                      # Text-to-speech (Stage 3)
│   ├── interface.go          # TTS synthesizer interface
│   ├── piper/                # Piper TTS
│   │   └── local.go
│   ├── melotts/              # MeloTTS
│   │   └── client.go
│   └── cache/                # Audio caching
│       └── manager.go
│
├── multimodal/               # Combined audio+vision (Stage 4)
│   ├── bridge.go             # Omni model orchestrator
│   └── context.go            # Context management
│
├── realtime/                 # Streaming pipeline (Stage 5)
│   ├── vad.go                # Voice Activity Detection
│   ├── buffer.go             # Circular audio buffer
│   └── streamer.go           # Chunked streaming
│
└── pipeline.go               # High-level orchestrator

```

---

## Model Support Matrix

### ASR Models

| Model | Size | Languages | Speed | Quality | Tier |
|-------|------|-----------|-------|---------|------|
| Whisper Tiny | 39MB | 10 | 0.3x RT | ★★★ | Tiny |
| Whisper Base | 74MB | 12 | 0.5x RT | ★★★★ | Small |
| Whisper Small | 244MB | 15 | 1.0x RT | ★★★★★ | Medium |
| Qwen3-ASR 0.6B | 600MB | 10 | 0.8x RT | ★★★★ | Small |
| Qwen3-ASR 1.7B | 1.7GB | 18 | 1.2x RT | ★★★★★ | Medium |

### TTS Models

| Model | Size | Languages | Speed | Quality | Tier |
|-------|------|-----------|-------|---------|------|
| Piper | 100MB | 10+ | 0.1x RT | ★★★ | Tiny |
| MeloTTS | 300MB | 8 | 0.5x RT | ★★★★ | Small |
| Bark | 2GB+ | 10 | 5.0x RT | ★★★★★ | Large |

---

## Usage Examples

### Basic ASR (File → Text)

```go
import (
    "context"
    "howl-chat/internal/audio/types"
    "howl-chat/internal/audio/asr/whisper"
)

func main() {
    ctx := context.Background()
    
    // Initialize recognizer
    recognizer := whisper.New()
    err := recognizer.Initialize(ctx, types.ASRConfig{
        ModelPath: "models/whisper-base.bin",
        Language:  "en",
        Device:    "auto",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer recognizer.Release()
    
    // Transcribe
    result, err := recognizer.Transcribe(ctx, "meeting.wav", func(p int, msg string) {
        fmt.Printf("Progress: %d%% - %s\n", p, msg)
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Transcribed: %s\n", result.Text)
    fmt.Printf("Confidence: %.2f%%\n", result.Confidence*100)
}
```

### Basic TTS (Text → Audio)

```go
import (
    "context"
    "howl-chat/internal/audio/types"
    "howl-chat/internal/audio/tts/piper"
)

func main() {
    ctx := context.Background()
    
    // Initialize synthesizer
    synth := piper.New()
    err := synth.Initialize(ctx, types.TTSConfig{
        ModelPath:    "models/piper-en_US-lessac-medium.onnx",
        CacheEnabled: true,
        CacheDir:     "cache/tts",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer synth.Release()
    
    // Synthesize
    result, err := synth.Synthesize(ctx, types.SynthesisRequest{
        Text:  "Welcome to HOWL Chat!",
        Voice: "en_US-lessac-medium",
        Speed: 1.0,
    }, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Generated: %s\n", result.AudioPath)
    fmt.Printf("Duration: %.1f seconds\n", result.Duration)
}
```

### Hardware Auto-Detection

```go
// Automatically detect hardware and recommend models
tier := types.DetectMemoryTier()
config := types.RecommendedConfigs[tier]

fmt.Printf("Detected tier: %s\n", tier)
fmt.Printf("Recommended ASR: %s\n", config.ASRModel)
fmt.Printf("Recommended TTS: %s\n", config.TTSModel)
```

---

## Configuration

### YAML Config

```yaml
# config/audio.yaml
audio:
  enabled: true
  
  asr:
    primary_model: "qwen3-asr-1.7b"
    fallback_model: "whisper-base"
    language: "auto"
    
  tts:
    enabled: true
    model: "piper"
    voice: "en_US-lessac-medium"
    speed: 1.0
    cache_enabled: true
    cache_size_mb: 100
    
  hardware:
    memory_tier: "auto"  # auto-detect
    prefer_gpu: true
```

### Environment Variables

```bash
export HOWL_AUDIO_ENABLED=true
export HOWL_ASR_MODEL=whisper-base
export HOWL_TTS_MODEL=piper
export HOWL_AUDIO_CACHE_DIR=/var/cache/howl/tts
```

---

## Testing Strategy

### Unit Tests

```bash
# Run all audio tests
go test ./internal/audio/... -v

# Run specific stage tests
go test ./internal/audio/types/... -v
go test ./internal/audio/asr/... -v
```

### Integration Tests

```bash
# Test with real audio files
go test ./internal/audio/integration -tags=integration

# Test fallback chains
go test ./internal/audio/asr/fallback -v
```

### Benchmarks

```bash
# ASR performance benchmarks
go test -bench=BenchmarkASR ./internal/audio/asr/whisper

# TTS latency benchmarks
go test -bench=BenchmarkTTS ./internal/audio/tts/piper
```

---

## Troubleshooting

### Common Issues

**Q: ASR model fails to load**  
A: Check available RAM against model requirements. Use `types.EstimateMemoryUsage()` to calculate needs.

**Q: Audio format not supported**  
A: Ensure ffmpeg is installed. Run `ffmpeg -version` to verify.

**Q: TTS audio sounds robotic**  
A: Upgrade to MeloTTS or Bark for better quality. Check voice selection.

**Q: Real-time latency too high**  
A: Use a faster ASR model (Whisper Tiny/Base). Enable GPU acceleration.

### Debug Logging

```go
// Enable verbose logging
config := types.PipelineConfig{
    LogLevel: "debug",
    LogAudioStats: true,
}
```

---

## Development Roadmap

### Completed ✅

- [x] Architecture documentation
- [x] Core interfaces (Stage 1)
- [x] Model profiles & registry
- [x] Audio format utilities

### In Progress 🚧

- [ ] ASR implementations (Stage 2)
- [ ] TTS implementations (Stage 3)

### Planned 📅

- [ ] Multimodal integration (Stage 4)
- [ ] Real-time streaming (Stage 5)
- [ ] WebRTC bridge (Stage 5)

---

## Contributing

### Code Standards

1. **Follow SWEObeyMe**: Max 500 lines per file
2. **Interface-first**: Define interfaces before implementations
3. **Error handling**: Use `types.AudioError` for structured errors
4. **Documentation**: Every public function needs documentation
5. **Testing**: Minimum 80% coverage for new code

### Commit Messages

```
feat(audio): Add Whisper Tiny ASR support
fix(asr): Handle empty audio files gracefully
docs: Update ASR model compatibility matrix
test(tts): Add Piper voice selection tests
```

---

## License

MIT - See LICENSE file for details

## Contact

**Project Lead**: Senior Engineer  
**Status**: Active Development  
**Last Updated**: May 2026

---

## Acknowledgments

- Whisper.cpp by Georgi Gerganov
- Qwen models by Alibaba Cloud
- Piper TTS by Rhasspy
- MeloTTS by MyShell AI
