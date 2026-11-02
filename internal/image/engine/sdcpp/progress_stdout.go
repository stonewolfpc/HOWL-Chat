package sdcpp

// sd-server does not implement /sdapi/v1/progress: the only place the sampling
// step exists is its stdout ("10/20 - 5.32s/it"). Without reading that stream the
// job's step count never advances, so the UI sits at 0% for the whole generation
// and looks hung even while the image is being produced.

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"time"

	"howl-chat/internal/image/types"
)

// resetServerProgress clears the last observed step so a new generation starts
// from zero instead of showing the previous run's final value.
func (e *Engine) resetServerProgress() {
	e.srvMu.Lock()
	e.srvProgress = types.ProgressInfo{}
	e.srvUpdated = time.Time{}
	e.srvMu.Unlock()
}

// recordServerProgress parses one sd-server log line for a sampling step.
func (e *Engine) recordServerProgress(line string) {
	m := stepRe.FindStringSubmatch(line)
	if m == nil {
		return
	}
	step, errStep := strconv.Atoi(m[1])
	total, errTotal := strconv.Atoi(m[2])
	if errStep != nil || errTotal != nil || total <= 0 || step > total {
		return
	}
	e.srvMu.Lock()
	e.srvProgress = types.ProgressInfo{
		Progress:   float64(step) / float64(total),
		Step:       step,
		TotalSteps: total,
	}
	e.srvUpdated = time.Now()
	e.srvMu.Unlock()
}

// recentServerProgress reports the last parsed step if it is fresh enough that
// it can only belong to the run in progress.
func (e *Engine) recentServerProgress() (types.ProgressInfo, bool) {
	e.srvMu.Lock()
	defer e.srvMu.Unlock()
	if e.srvProgress.Step == 0 || time.Since(e.srvUpdated) > 10*time.Minute {
		return types.ProgressInfo{}, false
	}
	return e.srvProgress, true
}

// pumpServerOutput echoes sd-server output and feeds the progress parser.
// sd.cpp rewrites its progress line with carriage returns rather than newlines,
// so the scanner splits on either one; splitting on newline alone would buffer
// every step until the run finished.
func (e *Engine) pumpServerOutput(r io.Reader) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	sc.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		for i, b := range data {
			if b == '\n' || b == '\r' {
				return i + 1, data[:i], nil
			}
		}
		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	for sc.Scan() {
		line := sc.Text()
		if line != "" {
			fmt.Println(line)
		}
		e.recordServerProgress(line)
		e.recordPlacement(line)
	}
}
