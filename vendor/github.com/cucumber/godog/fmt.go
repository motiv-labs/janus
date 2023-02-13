package godog

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/cucumber/messages-go/v10"
)

type registeredFormatter struct {
	name        string
	description string
	fmt         FormatterFunc
}

var formatters []*registeredFormatter

// FindFmt searches available formatters registered
// and returns FormaterFunc matched by given
// format name or nil otherwise
func FindFmt(name string) FormatterFunc {
	for _, el := range formatters {
		if el.name == name {
			return el.fmt
		}
	}

	return nil
}

// Format registers a feature suite output
// formatter by given name, description and
// FormatterFunc constructor function, to initialize
// formatter with the output recorder.
func Format(name, description string, f FormatterFunc) {
	formatters = append(formatters, &registeredFormatter{
		name:        name,
		fmt:         f,
		description: description,
	})
}

// AvailableFormatters gives a map of all
// formatters registered with their name as key
// and description as value
func AvailableFormatters() map[string]string {
	fmts := make(map[string]string, len(formatters))

	for _, f := range formatters {
		fmts[f.name] = f.description
	}

	return fmts
}

// Formatter is an interface for feature runner
// output summary presentation.
//
// New formatters may be created to represent
// suite results in different ways. These new
// formatters needs to be registered with a
// godog.Format function call
type Formatter interface {
	TestRunStarted()
	Feature(*messages.GherkinDocument, string, []byte)
	Pickle(*messages.Pickle)
	Defined(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition)
	Failed(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition, error)
	Passed(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition)
	Skipped(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition)
	Undefined(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition)
	Pending(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition)
	Summary()
}

// ConcurrentFormatter is an interface for a Concurrent
// version of the Formatter interface.
//
// Deprecated: The formatters need to handle concurrency
// instead of being copied and synchronized for each thread.
type ConcurrentFormatter interface {
	Formatter
	Copy(ConcurrentFormatter)
	Sync(ConcurrentFormatter)
}

type storageFormatter interface {
	setStorage(*storage)
}

// FormatterFunc builds a formatter with given
// suite name and io.Writer to record output
type FormatterFunc func(string, io.Writer) Formatter

func isLastStep(pickle *messages.Pickle, step *messages.Pickle_PickleStep) bool {
	return pickle.Steps[len(pickle.Steps)-1].Id == step.Id
}

func printStepDefinitions(steps []*StepDefinition, w io.Writer) {
	var longest int
	for _, def := range steps {
		n := utf8.RuneCountInString(def.Expr.String())
		if longest < n {
			longest = n
		}
	}

	for _, def := range steps {
		n := utf8.RuneCountInString(def.Expr.String())
		location := def.definitionID()
		spaces := strings.Repeat(" ", longest-n)
		fmt.Fprintln(w, yellow(def.Expr.String())+spaces, blackb("# "+location))
	}

	if len(steps) == 0 {
		fmt.Fprintln(w, "there were no contexts registered, could not find any step definition..")
	}
}
