package host

import (
	"context"
	stderrors "errors"
	"fmt"
	"sync"
	"time"
)

// Stage is one named, ordered step of the application sequence.
//
// Run does the work — constructing a domain, wiring a dependency, saving state
// or tearing a resource down — and whatever it builds belongs to the domain
// that built it. Nothing produced here is stored in the lifecycle.
type Stage struct {
	Name string
	Run  func(context.Context) error
}

// StageReport records what one stage did in the most recent sequence.
type StageReport struct {
	Name     string
	Duration time.Duration
	Err      error
	Panicked bool
}

// ShutdownOptions carries the application's failure sink. Ordering and
// accumulation belong to the lifecycle; reporting a failure is an application
// concern, so the sink is supplied rather than assumed.
type ShutdownOptions struct {
	OnFailure func(stage string, err error, panicked bool)
}

// Lifecycle owns the application's ordered sequence: boot (construct domains
// and wire dependencies in dependency order), readiness, orderly shutdown, and
// the restart that re-runs the same sequence.
//
// It stores no domain value. The stage sets it keeps are the composition the
// composition root handed it, which is what makes a restart identical to a
// first boot; readiness, epoch and the per-stage reports are coordination
// facts, not state any domain reads.
type Lifecycle struct {
	gate *Gate

	mu         sync.Mutex
	ready      bool
	bootErr    error
	epoch      uint32
	bootStages []Stage
	stopStages []Stage
	bootReport []StageReport
	stopReport []StageReport
}

// NewLifecycle binds the sequence to the gate that admits work.
func NewLifecycle(gate *Gate) *Lifecycle {
	if gate == nil {
		gate = NewGate()
	}
	return &Lifecycle{gate: gate}
}

// Gate exposes the admission gate this lifecycle closes on shutdown.
func (l *Lifecycle) Gate() *Gate { return l.gate }

// Boot runs the stages in order and marks the application ready only when every
// stage completed.
//
// The first failing stage stops the sequence: a stage that depends on an
// unwired domain must not run against a half-built application. Booting after a
// shutdown reopens admission first, so a restart runs through the same path
// instead of a second, parallel one.
func (l *Lifecycle) Boot(ctx context.Context, stages []Stage) error {
	if ctx == nil {
		ctx = context.Background()
	}

	l.mu.Lock()
	restart := l.gate.Stopping()
	l.bootStages = stages
	l.bootReport = nil
	l.bootErr = nil
	l.ready = false
	l.epoch++
	epoch := l.epoch
	l.mu.Unlock()

	if restart {
		l.gate.Reopen()
		fmt.Printf("[lifecycle] admission reopened for boot %d\n", epoch)
	}

	started := time.Now()
	reports := make([]StageReport, 0, len(stages))
	var firstErr error
	for _, stage := range stages {
		report := runStage(ctx, stage)
		reports = append(reports, report)
		if report.Err != nil {
			firstErr = fmt.Errorf("%s: %w", stage.Name, report.Err)
			break
		}
	}

	l.mu.Lock()
	l.bootReport = reports
	l.bootErr = firstErr
	l.ready = firstErr == nil
	l.mu.Unlock()

	if firstErr != nil {
		fmt.Printf("[lifecycle] startup stopped at %q: %v\n", reports[len(reports)-1].Name, firstErr)
		return firstErr
	}
	fmt.Printf("[lifecycle] ready: %d stages in %s (boot %d)\n",
		len(reports), time.Since(started).Round(time.Millisecond), epoch)
	return nil
}

func runStage(ctx context.Context, stage Stage) (report StageReport) {
	report.Name = stage.Name
	started := time.Now()
	defer func() {
		report.Duration = time.Since(started)
		if recovered := recover(); recovered != nil {
			report.Panicked = true
			report.Err = fmt.Errorf("%s panicked: %v", stage.Name, recovered)
		}
	}()
	if stage.Run == nil {
		report.Err = fmt.Errorf("%s has no body", stage.Name)
		return report
	}
	report.Err = stage.Run(ctx)
	return report
}

// Shutdown closes admission exactly once and runs the shutdown stages in
// dependency order. Every stage runs even when an earlier one failed, each
// failure is reported with the stage that produced it, and concurrent callers
// share the one result instead of running cleanup twice.
func (l *Lifecycle) Shutdown(ctx context.Context, stages []Stage, opts ShutdownOptions) error {
	l.mu.Lock()
	l.stopStages = stages
	l.stopReport = nil
	l.mu.Unlock()

	return l.gate.Once(ctx, func(ctx context.Context) error {
		l.gate.Stop()
		reports, err := l.runShutdownStages(ctx, stages, opts)

		l.mu.Lock()
		l.stopReport = reports
		l.ready = false
		l.mu.Unlock()

		return err
	})
}

func (l *Lifecycle) runShutdownStages(ctx context.Context, stages []Stage, opts ShutdownOptions) ([]StageReport, error) {
	runner := NewStages(opts.OnFailure)
	reports := make([]StageReport, 0, len(stages))

	for _, stage := range stages {
		report := StageReport{Name: stage.Name}
		started := time.Now()
		runner.Run(stage.Name, func() (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					report.Panicked = true
					// Stages converts the panic to a stage failure at critical
					// severity; re-panic so that conversion still happens.
					panic(recovered)
				}
			}()
			err = stage.Run(ctx)
			report.Err = err
			return err
		})
		report.Duration = time.Since(started)
		if report.Panicked && report.Err == nil {
			report.Err = fmt.Errorf("%s panicked", stage.Name)
		}
		reports = append(reports, report)
	}

	return reports, runner.Err()
}

// Ready reports whether the most recent boot completed every stage.
func (l *Lifecycle) Ready() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.ready
}

// BootError reports why the most recent boot did not reach readiness.
func (l *Lifecycle) BootError() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.bootErr
}

// Epoch counts how many times this application has been booted, so a restart is
// distinguishable from a first start in diagnostics.
func (l *Lifecycle) Epoch() uint32 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.epoch
}

// Report returns the reports of the most recent boot and shutdown sequences.
func (l *Lifecycle) Report() (boot, shutdown []StageReport) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]StageReport(nil), l.bootReport...), append([]StageReport(nil), l.stopReport...)
}

// StageNames lists the stages of the most recent boot, in order.
func (l *Lifecycle) StageNames() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	names := make([]string, 0, len(l.bootStages))
	for _, stage := range l.bootStages {
		names = append(names, stage.Name)
	}
	return names
}

// ShutdownStageNames lists the stages of the most recent shutdown, in order.
func (l *Lifecycle) ShutdownStageNames() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	names := make([]string, 0, len(l.stopStages))
	for _, stage := range l.stopStages {
		names = append(names, stage.Name)
	}
	return names
}

// JoinFailures is the one place a caller composes stage failures, kept small so
// diagnostics do not need errors.Join directly.
func JoinFailures(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	return stderrors.Join(errs...)
}
