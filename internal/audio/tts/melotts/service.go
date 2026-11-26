package melotts

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/audio/types"
	"howl-chat/internal/paths"
	"howl-chat/internal/process"
)

const (
	// readinessTimeout bounds how long a started service may take to answer on
	// its loopback endpoint. It is a real bound: expiry terminates the service
	// and fails the load.
	readinessTimeout = 30 * time.Second
	readinessPoll    = 200 * time.Millisecond
	// terminateGrace is the long-lived service's graceful shutdown interval
	// before escalation to process-tree termination.
	terminateGrace = 5 * time.Second
)

// Service is the MeloTTS process Howl owns. It is never inferred from an open
// port: a listening endpoint without this ownership is reported as unowned and
// is not adopted.
type Service struct {
	registry *process.Registry

	mu      sync.Mutex
	child   *process.Child
	baseURL string
}

// NewService returns a service owner bound to registry.
func NewService(registry *process.Registry) *Service {
	return &Service{registry: registry}
}

// Ensure starts the service if this owner is not already running it, waits for
// readiness, and reports a real failure when it cannot start or become ready.
func (s *Service) Ensure(ctx context.Context, baseURL string) error {
	if s == nil {
		return types.NewAudioError(types.ErrCodeModelLoadFailed, "melotts service owner is nil", nil)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	baseURL = NormalizeURL(baseURL)
	if err := requireLoopback(baseURL); err != nil {
		return err
	}

	s.mu.Lock()
	if s.child != nil && !s.child.Exited() {
		s.mu.Unlock()
		return nil
	}
	s.baseURL = baseURL
	s.mu.Unlock()

	// An endpoint that already answers was not started by this process. Howl
	// will not claim an unowned service as its own.
	if reachable(baseURL) {
		return types.NewAudioError(types.ErrCodeModelLoadFailed,
			fmt.Sprintf("a MeloTTS endpoint is already listening at %s but this process did not start it; stop it or let Howl manage the service", baseURL), nil)
	}

	launcher, args, err := resolveServiceCommand()
	if err != nil {
		return err
	}

	env := append(os.Environ(), "HOWL_MELOTTS_URL="+baseURL)
	child, err := process.StartChild(s.registry, process.ChildSpec{
		Role:     process.RoleTTSService,
		Instance: "melotts|" + baseURL,
		Path:     launcher,
		Args:     args,
		Dir:      filepath.Dir(launcher),
		Env:      env,
		Capture:  true,
	})
	if err != nil {
		return types.NewAudioError(types.ErrCodeModelLoadFailed,
			fmt.Sprintf("MeloTTS service could not be started under ownership: %v", err), err)
	}

	s.mu.Lock()
	s.child = child
	s.mu.Unlock()

	deadline := time.Now().Add(readinessTimeout)
	for {
		if reachable(baseURL) {
			return nil
		}
		if child.Exited() {
			_ = child.Terminate(context.Background(), terminateGrace)
			return types.NewAudioError(types.ErrCodeModelLoadFailed,
				fmt.Sprintf("MeloTTS service exited before becoming ready: %s", strings.TrimSpace(child.Output())), child.Err())
		}
		if time.Now().After(deadline) {
			termErr := s.Terminate(context.Background(), terminateGrace)
			detail := fmt.Sprintf("MeloTTS service did not become ready at %s within %s", baseURL, readinessTimeout)
			if termErr != nil {
				detail += fmt.Sprintf("; termination: %v", termErr)
			}
			return types.NewAudioError(types.ErrCodeModelLoadFailed, detail, nil)
		}
		select {
		case <-ctx.Done():
			termErr := s.Terminate(context.Background(), terminateGrace)
			detail := "MeloTTS startup cancelled"
			if termErr != nil {
				detail += fmt.Sprintf("; termination: %v", termErr)
			}
			return types.NewAudioError(types.ErrCodeModelLoadFailed, detail, ctx.Err())
		case <-time.After(readinessPoll):
		}
	}
}

// Terminate stops the owned service gracefully first, escalates to forced
// process-tree termination, and reports failure when exit cannot be confirmed.
func (s *Service) Terminate(ctx context.Context, grace time.Duration) error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	child := s.child
	s.mu.Unlock()
	if child == nil {
		return nil
	}
	if grace <= 0 {
		grace = terminateGrace
	}
	if err := child.Terminate(ctx, grace); err != nil {
		return fmt.Errorf("melotts service termination: %w", err)
	}
	s.mu.Lock()
	if s.child == child {
		s.child = nil
	}
	s.mu.Unlock()
	return nil
}

// Owned reports whether this process is running the service.
func (s *Service) Owned() bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.child != nil && !s.child.Exited()
}

// Identity reports the registered identity of the owned service.
func (s *Service) Identity() (process.Identity, bool) {
	if s == nil {
		return process.Identity{}, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.child == nil {
		return process.Identity{}, false
	}
	return s.child.Identity(), true
}

// ---------------------------------------------------------------------------
// Shared, reference-counted ownership
// ---------------------------------------------------------------------------

type sharedEntry struct {
	service *Service
	refs    int
}

var sharedServices = struct {
	mu   sync.Mutex
	byID map[string]*sharedEntry
}{byID: make(map[string]*sharedEntry)}

// AcquireService returns the process's single owned service for baseURL,
// starting it on first use. Every provider instance that uses the endpoint holds
// one reference; the last release terminates the service.
func AcquireService(registry *process.Registry, baseURL string) (*Service, error) {
	baseURL = NormalizeURL(baseURL)
	sharedServices.mu.Lock()
	defer sharedServices.mu.Unlock()
	if sharedServices.byID == nil {
		sharedServices.byID = make(map[string]*sharedEntry)
	}
	entry := sharedServices.byID[baseURL]
	if entry == nil {
		entry = &sharedEntry{service: NewService(registry)}
		sharedServices.byID[baseURL] = entry
	}
	if err := entry.service.Ensure(context.Background(), baseURL); err != nil {
		if entry.refs == 0 {
			delete(sharedServices.byID, baseURL)
		}
		return nil, err
	}
	entry.refs++
	return entry.service, nil
}

// ReleaseService drops one reference. The last one terminates the service.
func ReleaseService(service *Service) error {
	if service == nil {
		return nil
	}
	service.mu.Lock()
	baseURL := NormalizeURL(service.baseURL)
	service.mu.Unlock()
	if baseURL == "" {
		return nil
	}

	sharedServices.mu.Lock()
	entry := sharedServices.byID[baseURL]
	if entry == nil || entry.service != service {
		sharedServices.mu.Unlock()
		return nil
	}
	entry.refs--
	if entry.refs > 0 {
		sharedServices.mu.Unlock()
		return nil
	}
	delete(sharedServices.byID, baseURL)
	sharedServices.mu.Unlock()

	return service.Terminate(context.Background(), terminateGrace)
}

// ServiceOwned reports whether baseURL is already owned by this process.
func ServiceOwned(baseURL string) bool {
	baseURL = NormalizeURL(baseURL)
	sharedServices.mu.Lock()
	entry := sharedServices.byID[baseURL]
	sharedServices.mu.Unlock()
	return entry != nil && entry.service.Owned()
}

// ---------------------------------------------------------------------------
// Discovery
// ---------------------------------------------------------------------------

// resolveServiceCommand finds a local launcher Howl can own. It never invents a
// dependency: without an explicit command or a packaged launcher the load fails.
func resolveServiceCommand() (string, []string, error) {
	if command := strings.TrimSpace(os.Getenv("HOWL_MELOTTS_COMMAND")); command != "" {
		args := strings.Fields(strings.TrimSpace(os.Getenv("HOWL_MELOTTS_ARGS")))
		return command, args, nil
	}
	root := paths.ProjectRoot()
	candidates := []string{
		filepath.Join(root, "bin", "melotts-serve.exe"),
		filepath.Join(root, "bin", "melotts-serve"),
		filepath.Join(root, "runtime", "tts", "melotts-serve.exe"),
		filepath.Join(root, "runtime", "tts", "melotts-serve"),
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil, nil
		}
	}
	return "", nil, types.NewAudioError(types.ErrCodeModelNotFound,
		"MeloTTS is not installed for Howl ownership: set HOWL_MELOTTS_COMMAND to a local launcher that serves HOWL_MELOTTS_URL, or install bin/melotts-serve.exe", nil)
}

// NormalizeURL canonicalises an endpoint so ownership keys and readiness lookups
// cannot disagree over trailing slashes or spacing.
func NormalizeURL(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Host == "" {
		return strings.TrimSpace(raw)
	}
	return parsed.String()
}

// ConfiguredURL is the loopback endpoint the application is configured to use.
func ConfiguredURL() string {
	if raw := strings.TrimSpace(os.Getenv("HOWL_MELOTTS_URL")); raw != "" {
		return NormalizeURL(raw)
	}
	return "http://127.0.0.1:8888"
}

// Readiness reports whether MeloTTS is a ready Howl provider, not merely a
// listening socket. Readiness requires that this process started and still owns
// the service AND that the owned service answers. An endpoint Howl did not start
// is never advertised, and an owned service that stopped answering is not ready.
func Readiness(rawURL string) (bool, string) {
	baseURL := strings.TrimSpace(rawURL)
	if baseURL == "" {
		baseURL = ConfiguredURL()
	}
	baseURL = NormalizeURL(baseURL)
	if err := requireLoopback(baseURL); err != nil {
		return false, err.Error()
	}
	if !ServiceOwned(baseURL) {
		if reachable(baseURL) {
			return false, fmt.Sprintf("MeloTTS is listening at %s but this process did not start it; an unowned endpoint is not adopted as a Howl provider", baseURL)
		}
		return false, "MeloTTS is not running as a Howl-managed local service at " + baseURL
	}
	if !reachable(baseURL) {
		return false, "the Howl-managed MeloTTS service stopped answering at " + baseURL
	}
	return true, "Howl-managed MeloTTS service is running at " + baseURL
}

func requireLoopback(raw string) error {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Host == "" {
		return types.NewAudioError(types.ErrCodeModelNotFound, fmt.Sprintf("MeloTTS URL is invalid: %q", raw), err)
	}
	host := parsed.Hostname()
	if host == "localhost" {
		return nil
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		return nil
	}
	return types.NewAudioError(types.ErrCodeFormatNotSupported,
		fmt.Sprintf("MeloTTS must be a Howl-managed local service on a loopback address; %q is not loopback", raw), nil)
}

func reachable(baseURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil || parsed.Host == "" {
		return false
	}
	address := parsed.Host
	if parsed.Port() == "" {
		port := "80"
		if parsed.Scheme == "https" {
			port = "443"
		}
		address = net.JoinHostPort(parsed.Hostname(), port)
	}
	conn, err := net.DialTimeout("tcp", address, 300*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
