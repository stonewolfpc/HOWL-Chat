//go:build windows

package main

import (
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func useWASAPIProbe(t *testing.T, probe func(uint32) wasapiProbeOutcome) {
	t.Helper()
	prev := wasapiProbe
	wasapiProbe = probe
	t.Cleanup(func() { wasapiProbe = prev })
}

func successfulWASAPIOutcome(role uint32, hz uint32, channels uint16, bits uint16) wasapiProbeOutcome {
	return wasapiProbeOutcome{
		EndpointFound: true,
		Initialized:   true,
		Role:          role,
		Format: waveFormatEx{
			FormatTag:     0xFFFE,
			Channels:      channels,
			SamplesPerSec: hz,
			BitsPerSample: bits,
		},
	}
}

type fakeWASAPIBackend struct {
	enumHR          uint32
	defaultHR       uint32
	stateHR         uint32
	activateHR      uint32
	formatHR        uint32
	initHR          uint32
	state           uint32
	format          waveFormatEx
	roles           []uint32
	defaultHRByRole map[uint32]uint32
	enumToken       byte
	deviceToken     byte
	clientToken     byte
	acquired        int
	released        int
	freed           int
}

func (f *fakeWASAPIBackend) createEnumerator() (unsafe.Pointer, uint32) {
	if f.enumHR != sOK {
		return nil, f.enumHR
	}
	f.acquired++
	return unsafe.Pointer(&f.enumToken), sOK
}

func (f *fakeWASAPIBackend) getDefaultEndpoint(_ unsafe.Pointer, role uint32) (unsafe.Pointer, uint32) {
	f.roles = append(f.roles, role)
	if hr, ok := f.defaultHRByRole[role]; ok && hr != sOK {
		return nil, hr
	}
	if f.defaultHR != sOK {
		return nil, f.defaultHR
	}
	f.acquired++
	return unsafe.Pointer(&f.deviceToken), sOK
}

func (f *fakeWASAPIBackend) getState(_ unsafe.Pointer) (uint32, uint32) {
	if f.stateHR != sOK {
		return 0, f.stateHR
	}
	return f.state, sOK
}

func (f *fakeWASAPIBackend) activate(_ unsafe.Pointer) (unsafe.Pointer, uint32) {
	if f.activateHR != sOK {
		return nil, f.activateHR
	}
	f.acquired++
	return unsafe.Pointer(&f.clientToken), sOK
}

func (f *fakeWASAPIBackend) getMixFormat(_ unsafe.Pointer) (*waveFormatEx, uint32) {
	if f.formatHR != sOK {
		return nil, f.formatHR
	}
	return &f.format, sOK
}

func (f *fakeWASAPIBackend) initialize(_ unsafe.Pointer, _ *waveFormatEx) uint32 {
	return f.initHR
}

func (f *fakeWASAPIBackend) release(object unsafe.Pointer) {
	if object != nil {
		f.released++
	}
}

func (f *fakeWASAPIBackend) freeFormat(format *waveFormatEx) {
	if format != nil {
		f.freed++
	}
}

func TestDefaultCaptureRoleIsConsole(t *testing.T) {
	if defaultCaptureRole != eConsole {
		t.Fatalf("default capture role = %d, want eConsole (%d)", defaultCaptureRole, eConsole)
	}
	if defaultCaptureRole == eCommunications {
		t.Fatal("the probe must not silently use the communications default")
	}
}

func TestRunWASAPIProbeUsesConsoleRoleOnly(t *testing.T) {
	backend := &fakeWASAPIBackend{state: deviceStateActive}
	outcome := runWASAPIProbe(defaultCaptureRole, backend)

	if !reflect.DeepEqual(backend.roles, []uint32{eConsole}) {
		t.Fatalf("GetDefaultAudioEndpoint roles = %v, want only [eConsole]", backend.roles)
	}
	if outcome.Role != eConsole || !outcome.Initialized {
		t.Fatalf("outcome = %#v, want an initialized console endpoint", outcome)
	}
}

func TestConsoleDefaultDoesNotFallBackToCommunications(t *testing.T) {
	// Simulate Windows having a usable communications default while the console
	// default is absent. The probe must report no-device, not validate comms.
	backend := &fakeWASAPIBackend{
		state: deviceStateActive,
		defaultHRByRole: map[uint32]uint32{
			eConsole:        eNotFound,
			eCommunications: sOK,
		},
	}
	outcome := runWASAPIProbe(eConsole, backend)
	state, detail := classifyWASAPIOutcome(outcome)

	if state != MicStateNoDevice {
		t.Fatalf("state = %q, want no-device when the console default is absent", state)
	}
	if len(backend.roles) != 1 || backend.roles[0] != eConsole {
		t.Fatalf("roles probed = %v, want exactly [eConsole] and no communications fallback", backend.roles)
	}
	if !strings.Contains(detail, "console") {
		t.Fatalf("detail = %q, want the console role named", detail)
	}
}

func TestClassifyWASAPIOutcome(t *testing.T) {
	tests := []struct {
		name       string
		outcome    wasapiProbeOutcome
		wantState  MicrophoneState
		wantDetail []string
	}{
		{
			name:       "no default endpoint",
			outcome:    wasapiProbeOutcome{Stage: stageGetDefault, HRESULT: eNotFound, Role: eConsole},
			wantState:  MicStateNoDevice,
			wantDetail: []string{"E_NOTFOUND", "console"},
		},
		{
			name:       "default endpoint disabled",
			outcome:    wasapiProbeOutcome{Stage: stageState, HRESULT: sOK, DeviceState: deviceStateDisabled, EndpointFound: true, Role: eConsole},
			wantState:  MicStateUnavailable,
			wantDetail: []string{"disabled"},
		},
		{
			name:       "default endpoint unplugged",
			outcome:    wasapiProbeOutcome{Stage: stageState, HRESULT: sOK, DeviceState: deviceStateUnplug, EndpointFound: true, Role: eConsole},
			wantState:  MicStateUnavailable,
			wantDetail: []string{"unplugged"},
		},
		{
			name:       "access denied from initialize",
			outcome:    wasapiProbeOutcome{Stage: stageInitialize, HRESULT: eAccess, EndpointFound: true, Role: eConsole},
			wantState:  MicStateAccessDenied,
			wantDetail: []string{"E_ACCESSDENIED"},
		},
		{
			name:       "access denied from get state",
			outcome:    wasapiProbeOutcome{Stage: stageState, HRESULT: eAccess, EndpointFound: true, Role: eConsole},
			wantState:  MicStateAccessDenied,
			wantDetail: []string{"E_ACCESSDENIED"},
		},
		{
			name:       "winmm not enabled is not access denied",
			outcome:    wasapiProbeOutcome{Stage: stageInitialize, HRESULT: 3, EndpointFound: true, Role: eConsole},
			wantState:  MicStateCaptureInitFailed,
			wantDetail: []string{"HRESULT 0x00000003"},
		},
		{
			name:       "unsupported format",
			outcome:    wasapiProbeOutcome{Stage: stageInitialize, HRESULT: audclntUnsupportedFormat, EndpointFound: true, Role: eConsole},
			wantState:  MicStateCaptureInitFailed,
			wantDetail: []string{"mix format"},
		},
		{
			name:       "endpoint invalidated during initialization",
			outcome:    wasapiProbeOutcome{Stage: stageInitialize, HRESULT: audclntDeviceInvalidated, EndpointFound: true, Role: eConsole},
			wantState:  MicStateUnavailable,
			wantDetail: []string{"no longer available"},
		},
		{
			name:       "endpoint in use",
			outcome:    wasapiProbeOutcome{Stage: stageInitialize, HRESULT: audclntDeviceInUse, EndpointFound: true, Role: eConsole},
			wantState:  MicStateUnavailable,
			wantDetail: []string{"no longer available"},
		},
		{
			name:       "audio service not running",
			outcome:    wasapiProbeOutcome{Stage: stageInitialize, HRESULT: audclntServiceNotRunning, EndpointFound: true, Role: eConsole},
			wantState:  MicStateUnavailable,
			wantDetail: []string{"no longer available"},
		},
		{
			name:       "generic initialization failure",
			outcome:    wasapiProbeOutcome{Stage: stageInitialize, HRESULT: eNoIntf, EndpointFound: true, Role: eConsole},
			wantState:  MicStateCaptureInitFailed,
			wantDetail: []string{"0x80004002"},
		},
		{
			name:       "mix format acquisition failure",
			outcome:    wasapiProbeOutcome{Stage: stageMixFormat, HRESULT: ePointer, EndpointFound: true, Role: eConsole},
			wantState:  MicStateCaptureInitFailed,
			wantDetail: []string{"GetMixFormat"},
		},
		{
			name:       "initialized with endpoint mix format",
			outcome:    successfulWASAPIOutcome(eConsole, 48000, 2, 32),
			wantState:  MicStateReady,
			wantDetail: []string{"console", "48000", "2 channel", "recording and transcription were not started"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state, detail := classifyWASAPIOutcome(test.outcome)
			if state != test.wantState {
				t.Fatalf("state = %q, want %q (detail %q)", state, test.wantState, detail)
			}
			for _, want := range test.wantDetail {
				if !strings.Contains(detail, want) {
					t.Fatalf("detail = %q, want it to contain %q", detail, want)
				}
			}
		})
	}
}

func TestRunWASAPIProbeCleanupOnEveryPath(t *testing.T) {
	tests := []struct {
		name         string
		backend      *fakeWASAPIBackend
		wantStage    wasapiStage
		wantState    MicrophoneState
		wantAcquired int
		wantReleased int
		wantFreed    int
	}{
		{
			name:      "enumerator creation fails",
			backend:   &fakeWASAPIBackend{enumHR: eAccess},
			wantStage: stageCOM,
			wantState: MicStateCaptureInitFailed,
		},
		{
			name:         "no default endpoint",
			backend:      &fakeWASAPIBackend{defaultHR: eNotFound},
			wantStage:    stageGetDefault,
			wantState:    MicStateNoDevice,
			wantAcquired: 1,
			wantReleased: 1,
		},
		{
			name:         "state lookup access denied",
			backend:      &fakeWASAPIBackend{stateHR: eAccess},
			wantStage:    stageState,
			wantState:    MicStateAccessDenied,
			wantAcquired: 2,
			wantReleased: 2,
		},
		{
			name:         "endpoint inactive",
			backend:      &fakeWASAPIBackend{state: deviceStateDisabled},
			wantStage:    stageState,
			wantState:    MicStateUnavailable,
			wantAcquired: 2,
			wantReleased: 2,
		},
		{
			name:         "activate fails",
			backend:      &fakeWASAPIBackend{state: deviceStateActive, activateHR: eAccess},
			wantStage:    stageActivate,
			wantState:    MicStateAccessDenied,
			wantAcquired: 2,
			wantReleased: 2,
		},
		{
			name:         "mix format fails",
			backend:      &fakeWASAPIBackend{state: deviceStateActive, formatHR: eNoIntf},
			wantStage:    stageMixFormat,
			wantState:    MicStateCaptureInitFailed,
			wantAcquired: 3,
			wantReleased: 3,
		},
		{
			name:         "initialize fails",
			backend:      &fakeWASAPIBackend{state: deviceStateActive, initHR: eInval},
			wantStage:    stageInitialize,
			wantState:    MicStateCaptureInitFailed,
			wantAcquired: 3,
			wantReleased: 3,
			wantFreed:    1,
		},
		{
			name:         "success",
			backend:      &fakeWASAPIBackend{state: deviceStateActive},
			wantState:    MicStateReady,
			wantAcquired: 3,
			wantReleased: 3,
			wantFreed:    1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outcome := runWASAPIProbe(eConsole, test.backend)
			if outcome.Stage != test.wantStage {
				t.Fatalf("stage = %q, want %q", outcome.Stage, test.wantStage)
			}
			state, _ := classifyWASAPIOutcome(outcome)
			if state != test.wantState {
				t.Fatalf("state = %q, want %q", state, test.wantState)
			}
			if test.backend.acquired != test.wantAcquired {
				t.Fatalf("acquired objects = %d, want %d", test.backend.acquired, test.wantAcquired)
			}
			if test.backend.released != test.wantReleased {
				t.Fatalf("released objects = %d, want %d", test.backend.released, test.wantReleased)
			}
			if test.backend.acquired != test.backend.released {
				t.Fatalf("acquired %d object(s) but released %d", test.backend.acquired, test.backend.released)
			}
			if test.backend.freed != test.wantFreed {
				t.Fatalf("freed format pointers = %d, want %d", test.backend.freed, test.wantFreed)
			}
		})
	}
}

func TestNoDefaultEndpointProducesExactNoDeviceStatus(t *testing.T) {
	useWASAPIProbe(t, func(role uint32) wasapiProbeOutcome {
		return wasapiProbeOutcome{Stage: stageGetDefault, HRESULT: eNotFound, Role: role}
	})

	status := microphoneStatusNow()
	if status.State != MicStateNoDevice {
		t.Fatalf("state = %q, want no-device", status.State)
	}
	if status.Message != MicrophoneNoDeviceMessage {
		t.Fatalf("message = %q, want exact %q", status.Message, MicrophoneNoDeviceMessage)
	}
	if status.Role != "console" {
		t.Fatalf("role = %q, want console", status.Role)
	}
	if status.DeviceCount != 0 {
		t.Fatalf("device count = %d, want 0", status.DeviceCount)
	}
}

func TestProbeStatusReportsConsoleRoleAndMixFormat(t *testing.T) {
	useWASAPIProbe(t, func(role uint32) wasapiProbeOutcome {
		return successfulWASAPIOutcome(role, 44100, 1, 16)
	})

	status := probeMicrophone()
	if status.State != MicStateReady {
		t.Fatalf("state = %q, want ready", status.State)
	}
	if status.Role != "console" {
		t.Fatalf("role = %q, want console", status.Role)
	}
	if status.DeviceCount != 1 {
		t.Fatalf("device count = %d, want 1", status.DeviceCount)
	}
	for _, want := range []string{"44100", "1 channel", "16-bit", "recording and transcription were not started"} {
		if !strings.Contains(status.Detail, want) {
			t.Fatalf("detail = %q, want it to contain %q", status.Detail, want)
		}
	}
	if !strings.Contains(status.Message, "Recording has not started") {
		t.Fatalf("message = %q, want explicit recording-not-started wording", status.Message)
	}
}

// TestProbeMicrophoneOnHost exercises the real COM/WASAPI probe. It is the only
// test here that can touch a physical endpoint; all other WASAPI tests are
// deterministic mocks. It asserts only that a known state is returned. A ready
// result still does not exercise recording or transcription.
func TestProbeMicrophoneOnHost(t *testing.T) {
	for i := 0; i < 3; i++ {
		status := probeMicrophone()
		t.Logf("live WASAPI probe %d: state=%s role=%s count=%d detail=%q", i+1, status.State, status.Role, status.DeviceCount, status.Detail)

		switch status.State {
		case MicStateReady, MicStateNoDevice, MicStateUnavailable, MicStateAccessDenied, MicStateCaptureInitFailed:
		default:
			t.Fatalf("probe returned an unknown state %q", status.State)
		}
		if status.State == MicStateReady && status.Role != "console" {
			t.Fatalf("ready status probed role %q, want console", status.Role)
		}
		if !status.Usable() && status.Message == "" {
			t.Fatal("an unusable probe returned no user-facing message")
		}
	}
}
