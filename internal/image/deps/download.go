package deps

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ErrOffline marks a download that failed because the host could not be
// reached. It is separate from a server error so the UI can say "no connection"
// rather than blaming the file or the disk, and so the app can keep running
// with manual instructions.
var ErrOffline = errors.New("no internet connection")

// Progress reports download progress for one file.
type Progress func(filename string, received, total int64)

// Download fetches one artifact into dir and returns its final path.
//
// The file is written as "<name>.part" and renamed only after the whole body
// arrived, so an interrupted download can never be mistaken for a usable model.
func Download(ctx context.Context, client *http.Client, art Artifact, dir string, progress Progress) (string, error) {
	if art.URL == "" {
		return "", fmt.Errorf("%s has no direct download link", art.FileName)
	}
	if art.FileName == "" {
		return "", fmt.Errorf("dependency has no file name")
	}
	if client == nil {
		client = &http.Client{Timeout: 6 * time.Hour}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create %s: %w", dir, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, art.URL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "HOWL_Chat/dependency-fetch")

	resp, err := client.Do(req)
	if err != nil {
		if IsOffline(err) {
			return "", fmt.Errorf("%w: could not reach the download host (%v)", ErrOffline, err)
		}
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return "", fmt.Errorf("download %s: HTTP %s %s", art.FileName, resp.Status, strings.TrimSpace(string(body)))
	}

	dest := filepath.Join(dir, art.FileName)
	part := dest + ".part"
	out, err := os.Create(part)
	if err != nil {
		return "", fmt.Errorf("create %s: %w", part, err)
	}

	total := resp.ContentLength
	var received int64
	buf := make([]byte, 1<<20)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := out.Write(buf[:n]); werr != nil {
				out.Close()
				os.Remove(part)
				return "", fmt.Errorf("write %s: %w", part, werr)
			}
			received += int64(n)
			if progress != nil {
				progress(art.FileName, received, total)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			out.Close()
			os.Remove(part)
			if IsOffline(readErr) {
				return "", fmt.Errorf("%w: connection lost while downloading %s", ErrOffline, art.FileName)
			}
			return "", fmt.Errorf("download %s: %w", art.FileName, readErr)
		}
	}
	if err := out.Close(); err != nil {
		os.Remove(part)
		return "", err
	}
	// A short body means the transfer was truncated; keep nothing.
	if total > 0 && received != total {
		os.Remove(part)
		return "", fmt.Errorf("download %s: got %d of %d bytes", art.FileName, received, total)
	}
	if err := os.Rename(part, dest); err != nil {
		os.Remove(part)
		return "", fmt.Errorf("finish %s: %w", dest, err)
	}
	return dest, nil
}

// IsOffline reports whether an error means the network was unreachable, as
// opposed to the server answering with a failure.
func IsOffline(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrOffline) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}
	msg := strings.ToLower(err.Error())
	for _, marker := range []string{
		"no such host", "connection refused", "network is unreachable",
		"connection timed out", "tls handshake timeout", "no route to host",
		"server misbehaving", "proxyconnect", "connection reset",
	} {
		if strings.Contains(msg, marker) {
			return true
		}
	}
	return false
}
