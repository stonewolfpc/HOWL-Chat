package host

import (
	stderrors "errors"
	"fmt"
)

// Stages runs an ordered shutdown sequence where one broken component cannot
// prevent the remaining cleanup from running, and where every failure is
// reported with the stage that produced it.
//
// The escalation sink is supplied by the caller, because error reporting is an
// application concern; this type owns only ordering and accumulation.
type Stages struct {
	failures  []error
	onFailure func(stage string, err error, panicked bool)
}

// NewStages returns a sequence with no failures recorded. onFailure may be nil.
func NewStages(onFailure func(stage string, err error, panicked bool)) *Stages {
	return &Stages{onFailure: onFailure}
}

// Run executes one stage, records its failure, and continues. A panic is
// converted to a stage failure at critical severity rather than escaping.
func (s *Stages) Run(stage string, fn func() error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			s.fail(stage, fmt.Errorf("%s panicked: %v", stage, recovered), true)
		}
	}()
	if err := fn(); err != nil {
		s.fail(stage, fmt.Errorf("%s: %w", stage, err), false)
	}
}

func (s *Stages) fail(stage string, err error, panicked bool) {
	s.failures = append(s.failures, err)
	if s.onFailure != nil {
		s.onFailure(stage, err, panicked)
	}
}

// Err returns every recorded failure as one joined error, or nil.
func (s *Stages) Err() error {
	if len(s.failures) == 0 {
		return nil
	}
	return stderrors.Join(s.failures...)
}
