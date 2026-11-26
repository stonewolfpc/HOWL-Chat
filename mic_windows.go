//go:build windows

package main

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

// Windows endpoint-role constants. Howl describes and probes the default
// *console* capture endpoint because that is the Windows Settings default input
// and the role Chromium uses for a default microphone. eCommunications can point
// at a different device, so the probe deliberately never falls back to it: doing
// so would silently validate a microphone the user did not select for Howl.
const (
	eCapture            = 1 // EDataFlow
	eConsole            = 0 // ERole: games, notifications, voice commands
	eMultimedia         = 1
	eCommunications     = 2
	defaultCaptureRole  = eConsole
	clsctxAll           = 0x17
	coinitMultithreaded = 0x0
	audclntSharedMode   = 0x0
	deviceStateActive   = 0x00000001
	deviceStateDisabled = 0x00000002
	deviceStateNotPres  = 0x00000004
	deviceStateUnplug   = 0x00000008
)

// HRESULTs. Only codes documented for the exact API stage are mapped to a
// distinct user state; everything else is conservatively an initialization
// failure. There is no winmm-style "NOTENABLED" inference here.
const (
	sOK       = 0x00000000
	sFalse    = 0x00000001
	ePointer  = 0x80004003
	eNoIntf   = 0x80004002
	eInval    = 0x80070057
	eAccess   = 0x80070005
	eNotFound = 0x80070490
	rpcChg    = 0x80010106

	// AUDCLNT_FACILITY errors (audioclient.h).
	audclntDeviceInvalidated  = 0x88890004
	audclntUnsupportedFormat  = 0x88890008
	audclntDeviceInUse        = 0x8889000A
	audclntEndpointCreateFail = 0x8889000F
	audclntServiceNotRunning  = 0x88890010
	audclntResourcesInvalid   = 0x88890026
)

const (
	stageCOM        wasapiStage = "com"
	stageGetDefault             = "get-default-endpoint"
	stageState                  = "device-state"
	// stagePresence marks a presence-only probe: the endpoint exists and is
	// active, and no capture stream was opened.
	stagePresence   wasapiStage = "presence"
	stageActivate               = "activate"
	stageMixFormat              = "mix-format"
	stageInitialize             = "initialize"
)

type wasapiStage string

// wasapiProbeOutcome is the raw result of one WASAPI preflight. Keeping it
// separate from MicrophoneStatus makes every classification testable without a
// COM object or physical device.
type wasapiProbeOutcome struct {
	Stage         wasapiStage
	HRESULT       uint32
	DeviceState   uint32
	EndpointFound bool
	Initialized   bool
	Format        waveFormatEx
	Role          uint32
}

type guid struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

var (
	clsidMMDeviceEnumerator = guid{0xBCDE0395, 0xE52F, 0x467C, [8]byte{0x8E, 0x3D, 0xC4, 0x57, 0x92, 0x91, 0x69, 0x2E}}
	iidIMMDeviceEnumerator  = guid{0xA95664D2, 0x9614, 0x4F35, [8]byte{0xA7, 0x46, 0xDE, 0x8D, 0xB6, 0x36, 0x17, 0xE6}}
	iidIAudioClient         = guid{0x1CB9AD4C, 0xDBFA, 0x4C32, [8]byte{0xB1, 0x78, 0xC2, 0xF5, 0x68, 0xA7, 0x03, 0xB2}}
)

// waveFormatEx mirrors the fields of WAVEFORMATEX that the status detail needs.
// It is deliberately not passed back to WASAPI: the COM-allocated pointer from
// GetMixFormat (which may be WAVEFORMATEXTENSIBLE) is passed to Initialize
// unchanged so its cbSize and extension are preserved.
type waveFormatEx struct {
	FormatTag      uint16
	Channels       uint16
	SamplesPerSec  uint32
	AvgBytesPerSec uint32
	BlockAlign     uint16
	BitsPerSample  uint16
	CBSize         uint16
}

var (
	ole32                = syscall.NewLazyDLL("ole32.dll")
	procCoInitializeEx   = ole32.NewProc("CoInitializeEx")
	procCoUninitialize   = ole32.NewProc("CoUninitialize")
	procCoCreateInstance = ole32.NewProc("CoCreateInstance")
	procCoTaskMemFree    = ole32.NewProc("CoTaskMemFree")
)

// probeMicrophone runs a WASAPI preflight against the Windows default console
// capture endpoint. It initializes a shared-mode IAudioClient with that
// endpoint's own mix format, then releases everything. It never calls Start and
// never receives samples, so "ready" means only that endpoint initialization
// succeeded; it does not mean a recording or transcription succeeded.
//
// Initializing a capture client is a system-wide audio event: Windows treats an
// active capture stream as communications activity, switches Bluetooth headsets
// from A2DP to HFP, and resamples other applications' playback on the device.
// This probe is therefore diagnostics- and test-only: it must never be called
// from a status read, from the capture port, or from an apply path. Those use
// probeCapturePresence.
func probeMicrophone() MicrophoneStatus {
	outcome := wasapiProbe(defaultCaptureRole)
	state, detail := classifyWASAPIOutcome(outcome)
	status := newMicrophoneStatus(state, detail, boolToCount(outcome.EndpointFound))
	status.Role = wasapiRoleName(outcome.Role)
	return status
}

// probeCapturePresence asks only whether the default console capture endpoint
// exists and is active. It never activates the endpoint, initializes a client,
// or creates a stream, so it leaves the Windows capture state and every other
// application's audio untouched.
func probeCapturePresence() MicrophoneStatus {
	outcome := wasapiPresenceProbe(defaultCaptureRole)
	state, detail := classifyWASAPIOutcome(outcome)
	status := newMicrophoneStatus(state, detail, boolToCount(outcome.EndpointFound))
	status.Role = wasapiRoleName(outcome.Role)
	return status
}

var (
	wasapiProbe         = probeWASAPIDefaultCapture
	wasapiPresenceProbe = probeWASAPIDefaultCapturePresence
)

// wasapiBackend is the narrow COM surface used by the probe. It exists so the
// cleanup contract can be tested deterministically on every failure path; the
// real implementation is the thin comWASAPIBackend below.
type wasapiBackend interface {
	createEnumerator() (unsafe.Pointer, uint32)
	getDefaultEndpoint(enumerator unsafe.Pointer, role uint32) (unsafe.Pointer, uint32)
	getState(device unsafe.Pointer) (uint32, uint32)
	activate(device unsafe.Pointer) (unsafe.Pointer, uint32)
	getMixFormat(client unsafe.Pointer) (*waveFormatEx, uint32)
	initialize(client unsafe.Pointer, format *waveFormatEx) uint32
	release(object unsafe.Pointer)
	freeFormat(format *waveFormatEx)
}

// runWASAPIProbe acquires and releases every object with defers, including all
// early-return paths. It never starts the stream.
func runWASAPIProbe(role uint32, backend wasapiBackend) wasapiProbeOutcome {
	return runWASAPIProbeLevel(role, backend, true)
}

// runWASAPIPresence stops as soon as the endpoint's state is known. It is the
// only capture check safe to run repeatedly: nothing is activated or opened.
func runWASAPIPresence(role uint32, backend wasapiBackend) wasapiProbeOutcome {
	return runWASAPIProbeLevel(role, backend, false)
}

func runWASAPIProbeLevel(role uint32, backend wasapiBackend, openStream bool) wasapiProbeOutcome {
	out := wasapiProbeOutcome{Role: role}

	enumerator, hr := backend.createEnumerator()
	if hr != sOK {
		out.Stage = stageCOM
		out.HRESULT = hr
		return out
	}
	defer backend.release(enumerator)

	device, hr := backend.getDefaultEndpoint(enumerator, role)
	if hr != sOK {
		out.Stage = stageGetDefault
		out.HRESULT = hr
		return out
	}
	out.EndpointFound = true
	defer backend.release(device)

	state, hr := backend.getState(device)
	out.DeviceState = state
	if hr != sOK {
		out.Stage = stageState
		out.HRESULT = hr
		return out
	}
	if state != deviceStateActive {
		out.Stage = stageState
		out.HRESULT = sOK
		return out
	}

	if !openStream {
		// Presence only: an active endpoint is the whole answer. IMMDevice
		// is released without ever being activated.
		out.Stage = stagePresence
		return out
	}

	client, hr := backend.activate(device)
	if hr != sOK {
		out.Stage = stageActivate
		out.HRESULT = hr
		return out
	}

	var format *waveFormatEx
	// Release the client before freeing the format it returned, then always free
	// the COM-allocated format on every path after GetMixFormat.
	defer func() {
		backend.release(client)
		if format != nil {
			backend.freeFormat(format)
		}
	}()

	format, hr = backend.getMixFormat(client)
	if hr != sOK || format == nil {
		out.Stage = stageMixFormat
		if hr == sOK {
			hr = ePointer
		}
		out.HRESULT = hr
		return out
	}
	out.Format = *format

	// AUDCLNT_SHAREMODE_SHARED; no stream flags; default buffer duration and
	// periodicity; the endpoint's own mix format; no session GUID. This
	// initializes the stream but does not start it.
	hr = backend.initialize(client, format)
	if hr != sOK {
		out.Stage = stageInitialize
		out.HRESULT = hr
		return out
	}
	out.Initialized = true
	return out
}

func probeWASAPIDefaultCapture(role uint32) wasapiProbeOutcome {
	return withCOMProbe(role, func(backend wasapiBackend) wasapiProbeOutcome {
		return runWASAPIProbe(role, backend)
	})
}

func probeWASAPIDefaultCapturePresence(role uint32) wasapiProbeOutcome {
	return withCOMProbe(role, func(backend wasapiBackend) wasapiProbeOutcome {
		return runWASAPIPresence(role, backend)
	})
}

func withCOMProbe(role uint32, run func(wasapiBackend) wasapiProbeOutcome) wasapiProbeOutcome {
	// COM apartments are thread-scoped. Keep CoInitializeEx, every interface
	// call, and CoUninitialize on the same OS thread so the probe never calls
	// an interface from a thread with a different (or no) COM initialization.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	hr := coInitializeEx(0, coinitMultithreaded)
	switch hr {
	case sOK, sFalse:
		// S_FALSE means COM was already initialized on this thread and, like
		// S_OK, must be balanced by CoUninitialize.
		defer coUninitialize()
	case rpcChg:
		// COM is already initialized in another apartment. Proceed without
		// uninitializing it; this call did not acquire COM.
	default:
		return wasapiProbeOutcome{Stage: stageCOM, HRESULT: hr, Role: role}
	}
	return run(comWASAPIBackend{})
}

type comWASAPIBackend struct{}

func (comWASAPIBackend) createEnumerator() (unsafe.Pointer, uint32) {
	var enumerator unsafe.Pointer
	hr := coCreateInstance(&clsidMMDeviceEnumerator, 0, clsctxAll, &iidIMMDeviceEnumerator, &enumerator)
	if hr == sOK && enumerator == nil {
		hr = ePointer
	}
	return enumerator, hr
}

func (comWASAPIBackend) getDefaultEndpoint(enumerator unsafe.Pointer, role uint32) (unsafe.Pointer, uint32) {
	var device unsafe.Pointer
	hr := comCall(enumerator, 4, uintptr(eCapture), uintptr(role), uintptr(unsafe.Pointer(&device)))
	if hr == sOK && device == nil {
		hr = ePointer
	}
	return device, hr
}

func (comWASAPIBackend) getState(device unsafe.Pointer) (uint32, uint32) {
	var state uint32
	hr := comCall(device, 6, uintptr(unsafe.Pointer(&state)))
	return state, hr
}

func (comWASAPIBackend) activate(device unsafe.Pointer) (unsafe.Pointer, uint32) {
	var client unsafe.Pointer
	hr := comCall(device, 3, uintptr(unsafe.Pointer(&iidIAudioClient)), uintptr(clsctxAll), 0, uintptr(unsafe.Pointer(&client)))
	if hr == sOK && client == nil {
		hr = ePointer
	}
	return client, hr
}

func (comWASAPIBackend) getMixFormat(client unsafe.Pointer) (*waveFormatEx, uint32) {
	var format *waveFormatEx
	hr := comCall(client, 8, uintptr(unsafe.Pointer(&format)))
	return format, hr
}

func (comWASAPIBackend) initialize(client unsafe.Pointer, format *waveFormatEx) uint32 {
	return comCall(client, 3,
		uintptr(audclntSharedMode),
		0,
		0,
		0,
		uintptr(unsafe.Pointer(format)),
		0,
	)
}

func (comWASAPIBackend) release(object unsafe.Pointer) {
	releaseCOM(&object)
}

func (comWASAPIBackend) freeFormat(format *waveFormatEx) {
	coTaskMemFree(unsafe.Pointer(format))
}

// classifyWASAPIOutcome is pure so every status and HRESULT path is testable
// without COM or hardware.
func classifyWASAPIOutcome(out wasapiProbeOutcome) (MicrophoneState, string) {
	role := wasapiRoleName(out.Role)
	if out.Stage == stagePresence {
		// Presence is all this probe proved: the endpoint exists and is active.
		// No capture client was created, so do not claim initialization.
		return MicStateReady, fmt.Sprintf("the default %s capture endpoint is present and active in Windows; capture was not opened (no stream, no recording)", role)
	}
	if out.Initialized {
		return MicStateReady, fmt.Sprintf(
			"the default %s capture endpoint initialized in WASAPI shared mode with its mix format (tag 0x%04X, %d Hz, %d channel(s), %d-bit); recording and transcription were not started",
			role, out.Format.FormatTag, out.Format.SamplesPerSec, out.Format.Channels, out.Format.BitsPerSample)
	}

	switch out.Stage {
	case stageCOM:
		return MicStateCaptureInitFailed, fmt.Sprintf("COM initialization failed (HRESULT 0x%08X)", out.HRESULT)
	case stageGetDefault:
		if out.HRESULT == eNotFound {
			return MicStateNoDevice, "GetDefaultAudioEndpoint(eCapture, eConsole) returned E_NOTFOUND: no default console capture endpoint exists"
		}
		if out.HRESULT == eAccess {
			return MicStateAccessDenied, "GetDefaultAudioEndpoint returned E_ACCESSDENIED"
		}
		return MicStateCaptureInitFailed, fmt.Sprintf("GetDefaultAudioEndpoint failed (HRESULT 0x%08X)", out.HRESULT)
	case stageState:
		if out.HRESULT == eAccess {
			return MicStateAccessDenied, "IMMDevice::GetState returned E_ACCESSDENIED"
		}
		if out.HRESULT == eNotFound {
			return MicStateUnavailable, fmt.Sprintf("the default %s capture endpoint disappeared during the probe (E_NOTFOUND)", role)
		}
		if out.HRESULT != sOK {
			return MicStateCaptureInitFailed, fmt.Sprintf("IMMDevice::GetState failed (HRESULT 0x%08X)", out.HRESULT)
		}
		return MicStateUnavailable, fmt.Sprintf("the default %s capture endpoint is %s", role, deviceStateName(out.DeviceState))
	case stageActivate:
		return classifyWASAPIEndpointError(out, "IMMDevice::Activate")
	case stageMixFormat:
		return classifyWASAPIEndpointError(out, "IAudioClient::GetMixFormat")
	case stageInitialize:
		return classifyWASAPIEndpointError(out, "IAudioClient::Initialize")
	default:
		return MicStateCaptureInitFailed, "the WASAPI preflight did not complete"
	}
}

func classifyWASAPIEndpointError(out wasapiProbeOutcome, operation string) (MicrophoneState, string) {
	switch out.HRESULT {
	case eAccess:
		return MicStateAccessDenied, fmt.Sprintf("%s returned E_ACCESSDENIED", operation)
	case audclntDeviceInvalidated, audclntResourcesInvalid, audclntDeviceInUse, audclntServiceNotRunning, audclntEndpointCreateFail, eNotFound:
		return MicStateUnavailable, fmt.Sprintf("%s returned HRESULT 0x%08X: the default capture endpoint is no longer available or is in use", operation, out.HRESULT)
	case audclntUnsupportedFormat:
		return MicStateCaptureInitFailed, fmt.Sprintf("%s rejected the endpoint's own mix format (HRESULT 0x%08X)", operation, out.HRESULT)
	default:
		return MicStateCaptureInitFailed, fmt.Sprintf("%s failed (HRESULT 0x%08X)", operation, out.HRESULT)
	}
}

func deviceStateName(state uint32) string {
	switch state {
	case deviceStateActive:
		return "active"
	case deviceStateDisabled:
		return "disabled"
	case deviceStateNotPres:
		return "not present"
	case deviceStateUnplug:
		return "unplugged"
	default:
		return fmt.Sprintf("in an unknown state (0x%08X)", state)
	}
}

func wasapiRoleName(role uint32) string {
	switch role {
	case eConsole:
		return "console"
	case eMultimedia:
		return "multimedia"
	case eCommunications:
		return "communications"
	default:
		return "unknown"
	}
}

func boolToCount(found bool) int {
	if found {
		return 1
	}
	return 0
}

// comCall invokes one stdcall COM vtable method. Method indices are IUnknown
// (0-2) followed by the interface's declared methods.
func comCall(this unsafe.Pointer, method int, args ...uintptr) uint32 {
	if this == nil {
		return ePointer
	}
	// The first word of a COM object is a pointer to its vtable. Read it as an
	// unsafe.Pointer, then index the vtable without converting a uintptr back
	// into unsafe.Pointer.
	vtable := *(*unsafe.Pointer)(this)
	if vtable == nil {
		return ePointer
	}
	fn := (*[16]uintptr)(vtable)[method]
	callArgs := make([]uintptr, 0, len(args)+1)
	callArgs = append(callArgs, uintptr(this))
	callArgs = append(callArgs, args...)
	r, _, _ := syscall.SyscallN(fn, callArgs...)
	return uint32(r)
}

// releaseCOM releases an acquired interface. IUnknown::Release has no failure
// HRESULT; it returns the remaining reference count, so there is no cleanup error
// to ignore. CoTaskMemFree and CoUninitialize are void as well.
func releaseCOM(ptr *unsafe.Pointer) {
	if ptr == nil || *ptr == nil {
		return
	}
	comCall(*ptr, 2)
	*ptr = nil
}

func coInitializeEx(reserved uintptr, coInit uint32) uint32 {
	r, _, _ := syscall.SyscallN(procCoInitializeEx.Addr(), reserved, uintptr(coInit))
	return uint32(r)
}

func coUninitialize() {
	procCoUninitialize.Call()
}

func coCreateInstance(clsid *guid, outer uintptr, context uint32, iid *guid, out *unsafe.Pointer) uint32 {
	r, _, _ := syscall.SyscallN(
		procCoCreateInstance.Addr(),
		uintptr(unsafe.Pointer(clsid)),
		outer,
		uintptr(context),
		uintptr(unsafe.Pointer(iid)),
		uintptr(unsafe.Pointer(out)),
	)
	return uint32(r)
}

func coTaskMemFree(ptr unsafe.Pointer) {
	procCoTaskMemFree.Call(uintptr(ptr))
}
