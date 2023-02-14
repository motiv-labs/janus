package godog

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cucumber/messages-go/v10"
)

const nanoSec = 1000000
const spec = "0.1.0"

func init() {
	Format("events", fmt.Sprintf("Produces JSON event stream, based on spec: %s.", spec), eventsFunc)
}

func eventsFunc(suite string, out io.Writer) Formatter {
	return &events{basefmt: newBaseFmt(suite, out)}
}

type events struct {
	*basefmt
}

func (f *events) event(ev interface{}) {
	data, err := json.Marshal(ev)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal stream event: %+v - %v", ev, err))
	}
	fmt.Fprintln(f.out, string(data))
}

func (f *events) Pickle(pickle *messages.Pickle) {
	f.basefmt.Pickle(pickle)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.event(&struct {
		Event     string `json:"event"`
		Location  string `json:"location"`
		Timestamp int64  `json:"timestamp"`
	}{
		"TestCaseStarted",
		f.scenarioLocation(pickle),
		timeNowFunc().UnixNano() / nanoSec,
	})

	if len(pickle.Steps) == 0 {
		// @TODO: is status undefined or passed? when there are no steps
		// for this scenario
		f.event(&struct {
			Event     string `json:"event"`
			Location  string `json:"location"`
			Timestamp int64  `json:"timestamp"`
			Status    string `json:"status"`
		}{
			"TestCaseFinished",
			f.scenarioLocation(pickle),
			timeNowFunc().UnixNano() / nanoSec,
			"undefined",
		})
	}
}

func (f *events) TestRunStarted() {
	f.basefmt.TestRunStarted()

	f.lock.Lock()
	defer f.lock.Unlock()

	f.event(&struct {
		Event     string `json:"event"`
		Version   string `json:"version"`
		Timestamp int64  `json:"timestamp"`
		Suite     string `json:"suite"`
	}{
		"TestRunStarted",
		spec,
		timeNowFunc().UnixNano() / nanoSec,
		f.suiteName,
	})
}

func (f *events) Feature(ft *messages.GherkinDocument, p string, c []byte) {
	f.basefmt.Feature(ft, p, c)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.event(&struct {
		Event    string `json:"event"`
		Location string `json:"location"`
		Source   string `json:"source"`
	}{
		"TestSource",
		fmt.Sprintf("%s:%d", p, ft.Feature.Location.Line),
		string(c),
	})
}

func (f *events) Summary() {
	// @TODO: determine status
	status := passed

	f.storage.mustGetPickleStepResultsByStatus(failed)

	if len(f.storage.mustGetPickleStepResultsByStatus(failed)) > 0 {
		status = failed
	} else if len(f.storage.mustGetPickleStepResultsByStatus(passed)) == 0 {
		if len(f.storage.mustGetPickleStepResultsByStatus(undefined)) > len(f.storage.mustGetPickleStepResultsByStatus(pending)) {
			status = undefined
		} else {
			status = pending
		}
	}

	snips := f.snippets()
	if len(snips) > 0 {
		snips = "You can implement step definitions for undefined steps with these snippets:\n" + snips
	}

	f.event(&struct {
		Event     string `json:"event"`
		Status    string `json:"status"`
		Timestamp int64  `json:"timestamp"`
		Snippets  string `json:"snippets"`
		Memory    string `json:"memory"`
	}{
		"TestRunFinished",
		status.String(),
		timeNowFunc().UnixNano() / nanoSec,
		snips,
		"", // @TODO not sure that could be correctly implemented
	})
}

func (f *events) step(pickle *messages.Pickle, pickleStep *messages.Pickle_PickleStep) {
	feature := f.storage.mustGetFeature(pickle.Uri)
	pickleStepResult := f.storage.mustGetPickleStepResult(pickleStep.Id)
	step := feature.findStep(pickleStep.AstNodeIds[0])

	var errMsg string
	if pickleStepResult.err != nil {
		errMsg = pickleStepResult.err.Error()
	}
	f.event(&struct {
		Event     string `json:"event"`
		Location  string `json:"location"`
		Timestamp int64  `json:"timestamp"`
		Status    string `json:"status"`
		Summary   string `json:"summary,omitempty"`
	}{
		"TestStepFinished",
		fmt.Sprintf("%s:%d", pickle.Uri, step.Location.Line),
		timeNowFunc().UnixNano() / nanoSec,
		pickleStepResult.Status.String(),
		errMsg,
	})

	if isLastStep(pickle, pickleStep) {
		var status string

		pickleStepResults := f.storage.mustGetPickleStepResultsByPickleID(pickle.Id)
		for _, stepResult := range pickleStepResults {
			switch stepResult.Status {
			case passed:
				status = passed.String()
			case failed:
				status = failed.String()
			case undefined:
				status = undefined.String()
			case pending:
				status = pending.String()
			}
		}

		f.event(&struct {
			Event     string `json:"event"`
			Location  string `json:"location"`
			Timestamp int64  `json:"timestamp"`
			Status    string `json:"status"`
		}{
			"TestCaseFinished",
			f.scenarioLocation(pickle),
			timeNowFunc().UnixNano() / nanoSec,
			status,
		})
	}
}

func (f *events) Defined(pickle *messages.Pickle, pickleStep *messages.Pickle_PickleStep, def *StepDefinition) {
	f.basefmt.Defined(pickle, pickleStep, def)

	f.lock.Lock()
	defer f.lock.Unlock()

	feature := f.storage.mustGetFeature(pickle.Uri)
	step := feature.findStep(pickleStep.AstNodeIds[0])

	if def != nil {
		m := def.Expr.FindStringSubmatchIndex(pickleStep.Text)[2:]
		var args [][2]int
		for i := 0; i < len(m)/2; i++ {
			pair := m[i : i*2+2]
			var idxs [2]int
			idxs[0] = pair[0]
			idxs[1] = pair[1]
			args = append(args, idxs)
		}

		if len(args) == 0 {
			args = make([][2]int, 0)
		}

		f.event(&struct {
			Event    string   `json:"event"`
			Location string   `json:"location"`
			DefID    string   `json:"definition_id"`
			Args     [][2]int `json:"arguments"`
		}{
			"StepDefinitionFound",
			fmt.Sprintf("%s:%d", pickle.Uri, step.Location.Line),
			def.definitionID(),
			args,
		})
	}

	f.event(&struct {
		Event     string `json:"event"`
		Location  string `json:"location"`
		Timestamp int64  `json:"timestamp"`
	}{
		"TestStepStarted",
		fmt.Sprintf("%s:%d", pickle.Uri, step.Location.Line),
		timeNowFunc().UnixNano() / nanoSec,
	})
}

func (f *events) Passed(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Passed(pickle, step, match)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.step(pickle, step)
}

func (f *events) Skipped(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Skipped(pickle, step, match)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.step(pickle, step)
}

func (f *events) Undefined(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Undefined(pickle, step, match)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.step(pickle, step)
}

func (f *events) Failed(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition, err error) {
	f.basefmt.Failed(pickle, step, match, err)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.step(pickle, step)
}

func (f *events) Pending(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Pending(pickle, step, match)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.step(pickle, step)
}

func (f *events) scenarioLocation(pickle *messages.Pickle) string {
	feature := f.storage.mustGetFeature(pickle.Uri)
	scenario := feature.findScenario(pickle.AstNodeIds[0])

	line := scenario.Location.Line
	if len(pickle.AstNodeIds) == 2 {
		_, row := feature.findExample(pickle.AstNodeIds[1])
		line = row.Location.Line
	}

	return fmt.Sprintf("%s:%d", pickle.Uri, line)
}
