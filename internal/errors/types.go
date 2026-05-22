package errors

import (
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------
// Error Codes - Categorized by subsystem and priority
// Format: SUBSYSTEM-CATEGORY-SEQUENCE
// ---------------------------------------------------------------------------

type Code string

// LLM Loading Pipeline (Critical)
const (
	LLMServerNotFound      Code = "LLM-SRV-NF"
	LLMPortConflict        Code = "LLM-PORT-CONF"
	LLMServerSpawnFail     Code = "LLM-SRV-SPAWN"
	LLMServerCrash         Code = "LLM-SRV-CRASH"
	LLMServerTimeout       Code = "LLM-SRV-TIME"
	LLMGGUFReadFail        Code = "LLM-GGUF-READ"
	LLMMemEstimateFail     Code = "LLM-MEM-EST"
	LLMResourceRegFail     Code = "LLM-RES-REG"
	LLMHTTPHealthFail      Code = "LLM-HTTP-HEALTH"
	LLMHTTPRequestFail     Code = "LLM-HTTP-REQ"
	LLMStreamInterrupted   Code = "LLM-STREAM-BRK"
	LLMModelAlreadyLoading Code = "LLM-LOAD-DUP"
	LLMConfigApplyFail     Code = "LLM-CFG-APPLY"
)

// Image Generation Pipeline (Critical/High)
const (
	IMGServerNotFound   Code = "IMG-SRV-NF"
	IMGServerSpawnFail  Code = "IMG-SRV-SPAWN"
	IMGModelLoadFail    Code = "IMG-MODEL-LOAD"
	IMGGGUFIncompatible Code = "IMG-GGUF-BAD"
	IMGCPUOffloadFail   Code = "IMG-CPU-OFF"
	IMGGenerateFail     Code = "IMG-GEN-FAIL"
	IMGServerTimeout    Code = "IMG-SRV-TIME"
	IMGRegistryFail     Code = "IMG-REG-FAIL"
	IMGGallerySaveFail  Code = "IMG-GLRY-SAVE"
)

// TTS / Audio Pipeline (High)
const (
	TTSPipelineInitFail Code = "TTS-PIPE-INIT"
	TTSModelLoadFail    Code = "TTS-MODEL-LOAD"
	TTSSynthesizeFail   Code = "TTS-SYNTH-FAIL"
	TTSVoiceSwitchFail  Code = "TTS-VOICE-SW"
	ASRModelLoadFail    Code = "ASR-MODEL-LOAD"
	ASRTranscribeFail   Code = "ASR-TRANS-FAIL"
	AudioHWDetectFail   Code = "AUD-HW-DET"
)

// Resource Management (High)
const (
	RESHWDetectFail   Code = "RES-HW-DET"
	RESGPUQueryFail   Code = "RES-GPU-QRY"
	RESAllocFail      Code = "RES-ALLOC"
	RESEvictFail      Code = "RES-EVICT"
	RESRegisterFail   Code = "RES-REG"
	RESUnregisterFail Code = "RES-UNREG"
)

// Chat Page (High/Medium)
const (
	CHATSaveFail        Code = "CHAT-SAVE"
	CHATLoadFail        Code = "CHAT-LOAD"
	CHATMessageFail     Code = "CHAT-MSG"
	CHATPersistInitFail Code = "CHAT-PERSIST"
	CHATExportFail      Code = "CHAT-EXPORT"
)

// Hardware Edge Cases (High/Medium)
const (
	SYSHWUnsupported Code = "SYS-HW-UNSUPPORTED"
	SYSHWMobile      Code = "SYS-HW-MOBILE"
	SYSHWPotato      Code = "SYS-HW-POTATO"
	SYSHWMacIncompat Code = "SYS-HW-MAC"
	SYSHWOldGPU      Code = "SYS-HW-OLD-GPU"
	SYSHWNoGPU       Code = "SYS-HW-NO-GPU"
)

// File CRUD - Characters (High/Medium)
const (
	CHARCreateFail Code = "CHAR-CREATE"
	CHARLoadFail   Code = "CHAR-LOAD"
	CHARSaveFail   Code = "CHAR-SAVE"
	CHARParseFail  Code = "CHAR-PARSE"
	CHARDeleteFail Code = "CHAR-DEL"
	CHARDirFail    Code = "CHAR-DIR"
)

// File CRUD - Worlds (High/Medium)
const (
	WORLDCreateFail Code = "WORLD-CREATE"
	WORLDLoadFail   Code = "WORLD-LOAD"
	WORLDSaveFail   Code = "WORLD-SAVE"
	WORLDParseFail  Code = "WORLD-PARSE"
	WORLDDeleteFail Code = "WORLD-DEL"
	WORLDDirFail    Code = "WORLD-DIR"
)

// File CRUD - Scenarios (High/Medium)
const (
	SCENCreateFail Code = "SCEN-CREATE"
	SCENLoadFail   Code = "SCEN-LOAD"
	SCENSaveFail   Code = "SCEN-SAVE"
	SCENParseFail  Code = "SCEN-PARSE"
	SCENDeleteFail Code = "SCEN-DEL"
	SCENDirFail    Code = "SCEN-DIR"
)

// File CRUD - Lorebooks (High/Medium)
const (
	LORECreateFail Code = "LORE-CREATE"
	LORELoadFail   Code = "LORE-LOAD"
	LORESaveFail   Code = "LORE-SAVE"
	LOREParseFail  Code = "LORE-PARSE"
	LOREDeleteFail Code = "LORE-DEL"
)

// Memory / RAG (High/Medium)
const (
	MEMGetFail    Code = "MEM-GET"
	MEMAddFail    Code = "MEM-ADD"
	MEMEditFail   Code = "MEM-EDIT"
	MEMDeleteFail Code = "MEM-DEL"
	MEMExportFail Code = "MEM-EXPORT"
	MEMInitFail   Code = "MEM-INIT"
	MEMIdxCorrupt Code = "MEM-IDX-CORRUPT"
	MEMEmbedFail  Code = "MEM-EMBED"
	RAGIdxCorrupt Code = "RAG-IDX-CORRUPT"
	RAGQueryFail  Code = "RAG-QUERY"
	RAGInitFail   Code = "RAG-INIT"
)

// Model Registry (Medium)
const (
	REGModelAddFail  Code = "REG-MODEL-ADD"
	REGModelListFail Code = "REG-MODEL-LIST"
	REGModelDelFail  Code = "REG-MODEL-DEL"
)

// Image Vision / Edge Cases (High/Medium)
const (
	IMGVisionOversize Code = "IMG-VISION-OVERSIZE"
	IMGVisionBlank    Code = "IMG-VISION-BLANK"
	IMGVisionFormat   Code = "IMG-VISION-FORMAT"
	IMGCorrupt        Code = "IMG-CORRUPT"
	IMGNoModelLoaded  Code = "IMG-NO-MODEL"
	IMGPromptEmpty    Code = "IMG-PROMPT-EMPTY"
)

// Audio Edge Cases (High/Medium)
const (
	AUDMicNotFound       Code = "AUD-MIC-NF"
	AUDFormatUnsupported Code = "AUD-FMT-UNSUPPORTED"
	AUDPipeNotInit       Code = "AUD-PIPE-NOTINIT"
	TTSEngineMissing     Code = "TTS-ENGINE-MISSING"
	ASREngineMissing     Code = "ASR-ENGINE-MISSING"
)

// LLM Inference Edge Cases (Critical/High)
const (
	LLMCtxOverflow   Code = "LLM-CTX-OVERFLOW"
	LLMJSONMalformed Code = "LLM-JSON-MALFORM"
	LLMNetTimeout    Code = "LLM-NET-TIME"
	LLMNetReset      Code = "LLM-NET-RESET"
	LLMNoModelLoaded Code = "LLM-NO-MODEL"
	LLMPromptEmpty   Code = "LLM-PROMPT-EMPTY"
)

// Settings / Migration (Medium)
const (
	SYSMigrateFail Code = "SYS-MIGRATE"
	SYSCfgCorrupt  Code = "SYS-CFG-CORRUPT"
	SYSCfgVersion  Code = "SYS-CFG-VERSION"
)

// System / Infrastructure (Medium)
const (
	SYSPanicRecovered   Code = "SYS-PANIC"
	SYSConfigReadFail   Code = "SYS-CFG-READ"
	SYSSettingsSaveFail Code = "SYS-SET-SAVE"
	SYSSettingsLoadFail Code = "SYS-SET-LOAD"
	FSReadFail          Code = "FS-READ"
	FSWriteFail         Code = "FS-WRITE"
	FSDirCreateFail     Code = "FS-DIR"
	NETRequestFail      Code = "NET-REQ"
	NETTimeout          Code = "NET-TIME"
)

// Frontend (Medium/Low)
const (
	FEEventBindFail    Code = "FE-EVT-BIND"
	FECallbackMissing  Code = "FE-CB-MISS"
	FESettingsSyncFail Code = "FE-SET-SYNC"
)

// ---------------------------------------------------------------------------
// Severity
// ---------------------------------------------------------------------------

type Severity int

const (
	Critical Severity = iota
	High
	Medium
	Low
)

func (s Severity) String() string {
	switch s {
	case Critical:
		return "critical"
	case High:
		return "high"
	case Medium:
		return "medium"
	case Low:
		return "low"
	}
	return "unknown"
}

// ---------------------------------------------------------------------------
// TrackedError - The primary error type used throughout the application
// ---------------------------------------------------------------------------

type TrackedError struct {
	Code      Code
	Severity  Severity
	Operation string
	Location  string
	Err       error
	Context   map[string]string
	Time      time.Time
}

func New(code Code, severity Severity, operation, location string, err error) *TrackedError {
	return &TrackedError{
		Code:      code,
		Severity:  severity,
		Operation: operation,
		Location:  location,
		Err:       err,
		Context:   make(map[string]string),
		Time:      time.Now().UTC(),
	}
}

func (e *TrackedError) WithContext(key, value string) *TrackedError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

func (e *TrackedError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}
	return string(e.Code)
}

func (e *TrackedError) Unwrap() error {
	return e.Err
}

// IsCritical returns true for errors that break subsystem functionality
func (e *TrackedError) IsCritical() bool {
	return e.Severity == Critical
}

// ShouldNotify returns true for errors that should surface to the user
func (e *TrackedError) ShouldNotify() bool {
	return e.Severity <= High
}

// ReportPayload returns the privacy-safe data sent in error reports
func (e *TrackedError) ReportPayload() map[string]interface{} {
	return map[string]interface{}{
		"code":      string(e.Code),
		"severity":  e.Severity.String(),
		"operation": e.Operation,
		"location":  e.Location,
		"time":      e.Time.Format(time.RFC3339),
		"ctx":       sanitizeContext(e.Context),
	}
}

func sanitizeContext(ctx map[string]string) map[string]string {
	clean := make(map[string]string, len(ctx))
	for k, v := range ctx {
		if isSensitive(k) {
			clean[k] = "[redacted]"
		} else {
			clean[k] = v
		}
	}
	return clean
}

func isSensitive(key string) bool {
	switch key {
	case "user", "username", "email", "ip", "host", "password", "token", "key":
		return true
	}
	return false
}
