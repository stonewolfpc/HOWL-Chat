package audio

import (
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	"howl-chat/internal/audio/capability"
	"howl-chat/internal/audio/tts/melotts"
	audiotypes "howl-chat/internal/audio/types"
	"howl-chat/internal/process"
)

func meloCapabilityRecord(t *testing.T) capability.Record {
	t.Helper()
	for _, record := range ProductionCapabilityRecords() {
		if record.Task == capability.TaskSpeechSynthesis && record.ID == string(audiotypes.TTSMeloTTS) {
			return record
		}
	}
	t.Fatal("capability catalog has no MeloTTS record")
	return capability.Record{}
}

// TestMeloTTSCapabilityRequiresOwnedService is the readiness boundary: a TCP
// listener Howl did not start is never advertised as a ready provider, and an
// owned service is advertised only while it is actually running.
func TestMeloTTSCapabilityRequiresOwnedService(t *testing.T) {
	// 1. An unowned loopback listener. This is exactly the case a TCP probe
	// used to report as available.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	unownedURL := "http://" + listener.Addr().String()
	t.Setenv("HOWL_MELOTTS_URL", unownedURL)

	record := meloCapabilityRecord(t)
	if record.ModelState.OK() {
		t.Fatalf("an unowned listener was advertised as ready: %+v", record.ModelState)
	}
	if !strings.Contains(record.ModelState.Detail, "did not start it") {
		t.Fatalf("unowned readiness reason does not explain ownership: %q", record.ModelState.Detail)
	}
	if record.Selectable() {
		t.Fatal("an unowned listener is selectable")
	}

	// 2. A Howl-managed service on a free port becomes ready.
	portAddr := freeLoopbackPortForReadiness(t)
	managedURL := "http://127.0.0.1:" + portAddr
	t.Setenv("HOWL_MELOTTS_URL", managedURL)
	t.Setenv("HOWL_MELOTTS_COMMAND", os.Args[0])
	t.Setenv("HOWL_MELOTTS_ARGS", "-test.run=TestAudioMeloServiceHelper")
	t.Setenv("HOWL_MELO_HELPER", "serve")
	t.Setenv("HOWL_MELO_HELPER_PORT", portAddr)

	service, err := melotts.AcquireService(process.NewRegistry(), managedURL)
	if err != nil {
		t.Skipf("managed MeloTTS helper could not start in this environment: %v", err)
	}

	record = meloCapabilityRecord(t)
	if !record.ModelState.OK() {
		t.Fatalf("an owned, running service is not ready: %+v", record.ModelState)
	}
	if !strings.Contains(record.ModelState.Detail, "Howl-managed") {
		t.Fatalf("ready reason does not establish ownership: %q", record.ModelState.Detail)
	}

	// 3. Releasing ownership removes readiness; nothing lingers.
	if err := melotts.ReleaseService(service); err != nil {
		t.Fatalf("release service: %v", err)
	}
	record = meloCapabilityRecord(t)
	if record.ModelState.OK() {
		t.Fatalf("a terminated service is still advertised as ready: %+v", record.ModelState)
	}
}

// TestMeloTTSCapabilityRefusesNonLoopback proves a remote endpoint can never be
// reported as a local Howl provider.
func TestMeloTTSCapabilityRefusesNonLoopback(t *testing.T) {
	t.Setenv("HOWL_MELOTTS_URL", "http://203.0.113.10:8888")
	record := meloCapabilityRecord(t)
	if record.ModelState.OK() {
		t.Fatalf("non-loopback MeloTTS was advertised as ready: %+v", record.ModelState)
	}
	if !strings.Contains(record.ModelState.Detail, "loopback") {
		t.Fatalf("reason does not name the loopback requirement: %q", record.ModelState.Detail)
	}
}

func freeLoopbackPortForReadiness(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve port: %v", err)
	}
	defer listener.Close()
	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected listener address %T", listener.Addr())
	}
	return strconv.Itoa(addr.Port)
}
