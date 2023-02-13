package godog

import (
	"time"

	"github.com/cucumber/godog/colors"
)

type testRunStarted struct {
	StartedAt time.Time
}

type pickleResult struct {
	PickleID  string
	StartedAt time.Time
}

type pickleStepResult struct {
	Status     stepResultStatus
	finishedAt time.Time
	err        error

	PickleID     string
	PickleStepID string

	def *StepDefinition
}

func newStepResult(pickleID, pickleStepID string, match *StepDefinition) pickleStepResult {
	return pickleStepResult{finishedAt: timeNowFunc(), PickleID: pickleID, PickleStepID: pickleStepID, def: match}
}

type sortPickleStepResultsByPickleStepID []pickleStepResult

func (s sortPickleStepResultsByPickleStepID) Len() int { return len(s) }
func (s sortPickleStepResultsByPickleStepID) Less(i, j int) bool {
	iID := mustConvertStringToInt(s[i].PickleStepID)
	jID := mustConvertStringToInt(s[j].PickleStepID)
	return iID < jID
}
func (s sortPickleStepResultsByPickleStepID) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type stepResultStatus int

const (
	passed stepResultStatus = iota
	failed
	skipped
	undefined
	pending
)

func (st stepResultStatus) clr() colors.ColorFunc {
	switch st {
	case passed:
		return green
	case failed:
		return red
	case skipped:
		return cyan
	default:
		return yellow
	}
}

func (st stepResultStatus) String() string {
	switch st {
	case passed:
		return "passed"
	case failed:
		return "failed"
	case skipped:
		return "skipped"
	case undefined:
		return "undefined"
	case pending:
		return "pending"
	default:
		return "unknown"
	}
}
