package audio

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"howl-chat/internal/audio/tts/melotts"
	"howl-chat/internal/audio/types"
	"howl-chat/internal/process"
	"howl-chat/internal/resource"
)

// ttsProvider is one loaded synthesis engine and its residency. The admission
// lease becomes this registration when the provider becomes resident and is
// released at that moment, so the same bytes are never counted as both a
// pending reservation and a resident model.
//
// refs counts active executions (synthesis or enumeration). A provider with
// refs > 0 cannot be reclaimed; a resident provider is not unloaded just because
// one execution finished.
type ttsProvider struct {
	key        string
	synth      types.TTSSynthesizer
	engine     string
	voicePath  string
	loaderID   string
	modelName  string
	registered bool
	refs       int
}

// audioRegistration is one model this pipeline registered with the shared
// resource manager. Shutdown removes exactly these entries; it never sweeps the
// global ledger for models another owner registered.
type audioRegistration struct {
	kind resource.ResourceType
	name string
}

func audioRegKey(kind resource.ResourceType, name string) string {
	return string(kind) + "\x00" + name
}

// providerByKey returns a resident provider by map key.
func (p *Pipeline) providerByKey(key string) *ttsProvider {
	if p == nil || key == "" {
		return nil
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.ttsProviders[key]
}

// activeProvider returns the provider p.ttsSynth currently refers to, so an
// execution lease always lands on the provider that actually runs.
func (p *Pipeline) activeProvider() *ttsProvider {
	if p == nil {
		return nil
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.ttsProviders[p.activeTTSKey]
}

// acquireTTSExecution takes the provider's execution lease. The returned
// release is idempotent and must be called when the operation finishes.
func (p *Pipeline) acquireTTSExecution(provider *ttsProvider) func() {
	if p == nil || provider == nil {
		return func() {}
	}
	p.mu.Lock()
	provider.refs++
	p.mu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			p.mu.Lock()
			if provider.refs > 0 {
				provider.refs--
			}
			p.mu.Unlock()
		})
	}
}

// ReclaimIdleTTSProviders unloads every resident provider that is not executing
// and is not the active global engine. A provider that is in use is refused, and
// a failed unload keeps its residency and registration intact for a retry.
func (p *Pipeline) ReclaimIdleTTSProviders() error {
	if p == nil {
		return nil
	}
	p.ttsLoadMu.Lock()
	defer p.ttsLoadMu.Unlock()
	return p.reclaimIdleTTSProvidersLocked("")
}

// reclaimIdleTTSProvidersLocked is the load path's reclamation step; the caller
// already holds ttsLoadMu. keep names the provider that was just installed and
// must survive, in addition to the active engine and anything executing.
func (p *Pipeline) reclaimIdleTTSProvidersLocked(keep string) error {
	p.mu.RLock()
	activeKey := p.activeTTSKey
	globalKey := p.globalTTSKey
	idle := make([]*ttsProvider, 0, len(p.ttsProviders))
	for key, provider := range p.ttsProviders {
		if key == activeKey || key == globalKey || key == keep || provider.refs > 0 {
			continue
		}
		idle = append(idle, provider)
	}
	p.mu.RUnlock()

	var errs []error
	for _, provider := range idle {
		if err := p.unloadTTSProvider(provider); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// unloadTTSProvider releases one provider's residency. Ownership is retained on
// any failure: the loader entry, the registration and the map entry all stay, so
// a later retry targets the same resource instead of a phantom.
func (p *Pipeline) unloadTTSProvider(provider *ttsProvider) error {
	if p == nil || provider == nil {
		return nil
	}
	p.mu.RLock()
	refs := provider.refs
	p.mu.RUnlock()
	if refs > 0 {
		return fmt.Errorf("tts provider %q is executing (%d active); reclamation refused", provider.key, refs)
	}

	if provider.loaderID != "" {
		if err := p.loader.Unload(provider.loaderID); err != nil {
			return fmt.Errorf("unload tts provider %q: %w", provider.key, err)
		}
	}

	p.mu.Lock()
	if provider.registered {
		resource.GetManager().UnregisterModel(resource.ResourceTTS, provider.modelName)
		provider.registered = false
	}
	delete(p.ttsProviders, provider.key)
	if p.audioRegistered != nil {
		delete(p.audioRegistered, audioRegKey(resource.ResourceTTS, provider.modelName))
	}
	if p.activeTTSKey == provider.key {
		p.activeTTSKey = ""
	}
	p.mu.Unlock()
	return nil
}

// ttsResidencyRequest asks the shared authority for the provider's residency
// and returns the request used for registration. A MeloTTS service this process
// already owns is already resident: charging it again would double-count the
// same service, so no new lease is taken.
func (p *Pipeline) ttsResidencyRequest(engine, modelPath string) (resource.AdmissionRequest, *resource.AdmissionLease, error) {
	req, err := ttsAdmissionRequest(engine, modelPath)
	if err != nil {
		return resource.AdmissionRequest{}, nil, err
	}
	if engine == string(types.TTSMeloTTS) && melotts.ServiceOwned(modelPath) {
		return req, nil, nil
	}
	lease, err := p.admitAudio(req)
	if err != nil {
		return req, nil, err
	}
	return req, lease, nil
}

// registerTTSResidency records the provider and its registration. The caller
// holds p.mu.
func (p *Pipeline) registerTTSResidencyLocked(provider *ttsProvider) {
	if provider == nil {
		return
	}
	if p.ttsProviders == nil {
		p.ttsProviders = make(map[string]*ttsProvider)
	}
	p.ttsProviders[provider.key] = provider
	if provider.modelName != "" {
		if p.audioRegistered == nil {
			p.audioRegistered = make(map[string]audioRegistration)
		}
		p.audioRegistered[audioRegKey(resource.ResourceTTS, provider.modelName)] = audioRegistration{
			kind: resource.ResourceTTS,
			name: provider.modelName,
		}
	}
}

// releaseOwnedAudioRegistrations removes exactly the models this pipeline
// registered. It is used by Shutdown after every provider and child is settled.
func (p *Pipeline) releaseOwnedAudioRegistrations() {
	p.mu.Lock()
	retained := make(map[string]struct{}, len(p.ttsProviders))
	for _, provider := range p.ttsProviders {
		if provider.registered {
			retained[audioRegKey(resource.ResourceTTS, provider.modelName)] = struct{}{}
		}
	}
	registrations := make([]audioRegistration, 0, len(p.audioRegistered))
	for key, reg := range p.audioRegistered {
		if _, keep := retained[key]; keep {
			continue
		}
		registrations = append(registrations, reg)
		delete(p.audioRegistered, key)
	}
	p.mu.Unlock()

	for _, reg := range registrations {
		resource.GetManager().UnregisterModel(reg.kind, reg.name)
	}
}

// stopOwnedTTSChildren escalates every live TTS child and service to confirmed
// exit. It is the shutdown backstop: the manager cancels operation contexts
// first, and this ensures no owned process survives a cancelled or stuck call.
func (p *Pipeline) stopOwnedTTSChildren() error {
	if p == nil || p.ttsRegistry == nil {
		return nil
	}
	errs := []error{
		process.StopRole(p.ttsRegistry, process.RoleTTSChild, process.DefaultGrace),
		process.StopRole(p.ttsRegistry, process.RoleTTSService, process.DefaultGrace),
	}
	return errors.Join(errs...)
}

// voiceEnumerator is implemented by providers whose enumeration spawns a
// process. The pipeline gives it the caller's operation scope so cancellation
// reaches the child.
type voiceEnumerator interface {
	VoicesContext(context.Context) ([]types.Voice, error)
}

// GetTTSVoicesContext enumerates voices as a managed operation: it holds the
// active provider's execution lease for the duration and hands the caller's
// context to providers that support cancellation.
func (p *Pipeline) GetTTSVoicesContext(ctx context.Context) ([]types.Voice, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	p.mu.RLock()
	synth := p.ttsSynth
	provider := p.ttsProviders[ttsEngineKey(p.activeTTSEngine, "")]
	p.mu.RUnlock()
	if synth == nil {
		return nil, nil
	}
	release := p.acquireTTSExecution(provider)
	defer release()

	if enumerator, ok := synth.(voiceEnumerator); ok {
		return enumerator.VoicesContext(ctx)
	}
	return synth.GetVoices(), nil
}

// boundedVoiceContext gives the interface's context-free enumeration a scope
// instead of letting it run unbounded.
func boundedVoiceContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 15*time.Second)
}
