package types

import "context"

// TTSProcessState is the supervisor's view of the helper process. It is
// deliberately separate from provider readiness: a live process can exist
// while its provider is unloaded or unavailable.
type TTSProcessState string

const (
	TTSProcessAbsent   TTSProcessState = "absent"
	TTSProcessStarting TTSProcessState = "starting"
	TTSProcessAlive    TTSProcessState = "alive"
	TTSProcessStopping TTSProcessState = "stopping"
	TTSProcessExited   TTSProcessState = "exited"
)

// TTSProviderState is the aggregate supervisor/provider state from the frozen
// helper contract.
type TTSProviderState string

const (
	TTSStateAbsent       TTSProviderState = "absent"
	TTSStateStarting     TTSProviderState = "starting"
	TTSStateHandshaking  TTSProviderState = "handshaking"
	TTSStateUnloaded     TTSProviderState = "unloaded"
	TTSStateLoading      TTSProviderState = "loading"
	TTSStateReady        TTSProviderState = "ready"
	TTSStateBusy         TTSProviderState = "busy"
	TTSStateDegraded     TTSProviderState = "degraded"
	TTSStateStopping     TTSProviderState = "stopping"
	TTSStateCrashed      TTSProviderState = "crashed"
	TTSStateIncompatible TTSProviderState = "incompatible"
	TTSStateUnavailable  TTSProviderState = "unavailable"
	TTSStateBackoff      TTSProviderState = "backoff"
)

// TTSFailure is the precise internal failure classification. Its mapping to
// AudioErrorCode is only a public compatibility projection.
type TTSFailure struct {
	Code      string
	Message   string
	Retryable bool
	Field     string
	Detail    map[string]string
}

// TTSCapability is one advertised helper/provider capability.
type TTSCapability struct {
	ID      string
	Version int
}

// TTSProviderStatus is the additive inspector status. It does not change
// TTSSynthesizer and never implies that a model was loaded.
type TTSProviderStatus struct {
	Available       bool
	ProcessState    TTSProcessState
	ProcessAlive    bool
	Ready           bool
	State           TTSProviderState
	Engine          string
	ProviderID      string
	ProviderVersion string
	HelperVersion   string
	ProtocolVersion string
	Capabilities    []TTSCapability
	Runtime         map[string]string
	ActiveRequestID string
	LastError       *TTSFailure
}

// TTSProviderCatalogRequest is the provider-neutral catalog request.
type TTSProviderCatalogRequest struct {
	Engine             string
	IncludeUnavailable bool
}

// TTSProviderCatalog is the provider-neutral catalog result. It is not the
// renderer-facing TTSVoiceCatalog JSON shape.
type TTSProviderCatalog struct {
	Engine          string
	ProviderID      string
	CatalogVersion  int
	Available       bool
	Detail          string
	DefaultVoice    string
	Languages       []string
	Voices          []TTSVoiceInfo
	Speed           *TTSSpeedRange
	SSMLSupported   bool
	StreamSupported bool
}

// TTSVoiceInfo extends the existing Voice with availability diagnostics.
type TTSVoiceInfo struct {
	Voice
	Available bool
	Detail    string
}

// TTSSpeedRange is the catalog's supported speed interval.
type TTSSpeedRange struct {
	Supported bool
	Min       float64
	Max       float64
}

// TTSProviderInspector is the additive status/catalog surface. Implementations
// may start and handshake the helper for status, but must not load a model
// solely for status.
type TTSProviderInspector interface {
	ProviderStatus(ctx context.Context) (TTSProviderStatus, error)
	ProviderCatalog(ctx context.Context, req TTSProviderCatalogRequest) (TTSProviderCatalog, error)
}
