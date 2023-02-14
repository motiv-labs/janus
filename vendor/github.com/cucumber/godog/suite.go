package godog

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"

	"github.com/cucumber/messages-go/v10"
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()
var typeOfBytes = reflect.TypeOf([]byte(nil))

// ErrUndefined is returned in case if step definition was not found
var ErrUndefined = fmt.Errorf("step is undefined")

// ErrPending should be returned by step definition if
// step implementation is pending
var ErrPending = fmt.Errorf("step implementation is pending")

// Suite allows various contexts
// to register steps and event handlers.
//
// When running a test suite, the instance of Suite
// is passed to all functions (contexts), which
// have it as a first and only argument.
//
// Note that all event hooks does not catch panic errors
// in order to have a trace information. Only step
// executions are catching panic error since it may
// be a context specific error.
//
// Deprecated: The current Suite initializer will be removed and replaced by
// two initializers, one for the Test Suite and one for the Scenarios.
// This struct will therefore not be exported in the future.
type Suite struct {
	steps    []*StepDefinition
	features []*feature

	fmt     Formatter
	storage *storage

	failed        bool
	randomSeed    int64
	stopOnFailure bool
	strict        bool

	// suite event handlers
	beforeSuiteHandlers    []func()
	beforeScenarioHandlers []func(*messages.Pickle)
	beforeStepHandlers     []func(*messages.Pickle_PickleStep)
	afterStepHandlers      []func(*messages.Pickle_PickleStep, error)
	afterScenarioHandlers  []func(*messages.Pickle, error)
	afterSuiteHandlers     []func()
}

// Step allows to register a *StepDefinition in Godog
// feature suite, the definition will be applied
// to all steps matching the given Regexp expr.
//
// It will panic if expr is not a valid regular
// expression or stepFunc is not a valid step
// handler.
//
// Note that if there are two definitions which may match
// the same step, then only the first matched handler
// will be applied.
//
// If none of the *StepDefinition is matched, then
// ErrUndefined error will be returned when
// running steps.
//
// Deprecated: The current Suite initializer will be removed and replaced by
// two initializers, one for the Test Suite and one for the Scenarios. Use
// func (ctx *ScenarioContext) Step instead.
func (s *Suite) Step(expr interface{}, stepFunc interface{}) {
	var regex *regexp.Regexp

	switch t := expr.(type) {
	case *regexp.Regexp:
		regex = t
	case string:
		regex = regexp.MustCompile(t)
	case []byte:
		regex = regexp.MustCompile(string(t))
	default:
		panic(fmt.Sprintf("expecting expr to be a *regexp.Regexp or a string, got type: %T", expr))
	}

	v := reflect.ValueOf(stepFunc)
	typ := v.Type()
	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("expected handler to be func, but got: %T", stepFunc))
	}

	if typ.NumOut() != 1 {
		panic(fmt.Sprintf("expected handler to return only one value, but it has: %d", typ.NumOut()))
	}

	def := &StepDefinition{
		Handler: stepFunc,
		Expr:    regex,
		hv:      v,
	}

	typ = typ.Out(0)
	switch typ.Kind() {
	case reflect.Interface:
		if !typ.Implements(errorInterface) {
			panic(fmt.Sprintf("expected handler to return an error, but got: %s", typ.Kind()))
		}
	case reflect.Slice:
		if typ.Elem().Kind() != reflect.String {
			panic(fmt.Sprintf("expected handler to return []string for multistep, but got: []%s", typ.Kind()))
		}
		def.nested = true
	default:
		panic(fmt.Sprintf("expected handler to return an error or []string, but got: %s", typ.Kind()))
	}

	s.steps = append(s.steps, def)
}

// BeforeSuite registers a function or method
// to be run once before suite runner.
//
// Use it to prepare the test suite for a spin.
// Connect and prepare database for instance...
//
// Deprecated: The current Suite initializer will be removed and replaced by
// two initializers, one for the Test Suite and one for the Scenarios. Use
// func (ctx *TestSuiteContext) BeforeSuite instead.
func (s *Suite) BeforeSuite(fn func()) {
	s.beforeSuiteHandlers = append(s.beforeSuiteHandlers, fn)
}

// BeforeScenario registers a function or method
// to be run before every pickle.
//
// It is a good practice to restore the default state
// before every scenario so it would be isolated from
// any kind of state.
//
// Deprecated: The current Suite initializer will be removed and replaced by
// two initializers, one for the Test Suite and one for the Scenarios. Use
// func (ctx *ScenarioContext) BeforeScenario instead.
func (s *Suite) BeforeScenario(fn func(*messages.Pickle)) {
	s.beforeScenarioHandlers = append(s.beforeScenarioHandlers, fn)
}

// BeforeStep registers a function or method
// to be run before every step.
//
// Deprecated: The current Suite initializer will be removed and replaced by
// two initializers, one for the Test Suite and one for the Scenarios. Use
// func (ctx *ScenarioContext) BeforeStep instead.
func (s *Suite) BeforeStep(fn func(*messages.Pickle_PickleStep)) {
	s.beforeStepHandlers = append(s.beforeStepHandlers, fn)
}

// AfterStep registers an function or method
// to be run after every step.
//
// It may be convenient to return a different kind of error
// in order to print more state details which may help
// in case of step failure
//
// In some cases, for example when running a headless
// browser, to take a screenshot after failure.
//
// Deprecated: The current Suite initializer will be removed and replaced by
// two initializers, one for the Test Suite and one for the Scenarios. Use
// func (ctx *ScenarioContext) AfterStep instead.
func (s *Suite) AfterStep(fn func(*messages.Pickle_PickleStep, error)) {
	s.afterStepHandlers = append(s.afterStepHandlers, fn)
}

// AfterScenario registers an function or method
// to be run after every pickle.
//
// Deprecated: The current Suite initializer will be removed and replaced by
// two initializers, one for the Test Suite and one for the Scenarios. Use
// func (ctx *ScenarioContext) AfterScenario instead.
func (s *Suite) AfterScenario(fn func(*messages.Pickle, error)) {
	s.afterScenarioHandlers = append(s.afterScenarioHandlers, fn)
}

// AfterSuite registers a function or method
// to be run once after suite runner
//
// Deprecated: The current Suite initializer will be removed and replaced by
// two initializers, one for the Test Suite and one for the Scenarios. Use
// func (ctx *TestSuiteContext) AfterSuite instead.
func (s *Suite) AfterSuite(fn func()) {
	s.afterSuiteHandlers = append(s.afterSuiteHandlers, fn)
}

func (s *Suite) run() {
	// run before suite handlers
	for _, f := range s.beforeSuiteHandlers {
		f()
	}
	// run features
	for _, f := range s.features {
		s.runFeature(f)
		if s.failed && s.stopOnFailure {
			// stop on first failure
			break
		}
	}
	// run after suite handlers
	for _, f := range s.afterSuiteHandlers {
		f()
	}
}

func (s *Suite) matchStep(step *messages.Pickle_PickleStep) *StepDefinition {
	def := s.matchStepText(step.Text)
	if def != nil && step.Argument != nil {
		def.args = append(def.args, step.Argument)
	}
	return def
}

func (s *Suite) runStep(pickle *messages.Pickle, step *messages.Pickle_PickleStep, prevStepErr error) (err error) {
	// run before step handlers
	for _, f := range s.beforeStepHandlers {
		f(step)
	}

	match := s.matchStep(step)
	s.fmt.Defined(pickle, step, match)

	// user multistep definitions may panic
	defer func() {
		if e := recover(); e != nil {
			err = &traceError{
				msg:   fmt.Sprintf("%v", e),
				stack: callStack(),
			}
		}

		if prevStepErr != nil {
			return
		}

		if err == ErrUndefined {
			return
		}

		sr := newStepResult(pickle.Id, step.Id, match)

		switch err {
		case nil:
			sr.Status = passed
			s.storage.mustInsertPickleStepResult(sr)

			s.fmt.Passed(pickle, step, match)
		case ErrPending:
			sr.Status = pending
			s.storage.mustInsertPickleStepResult(sr)

			s.fmt.Pending(pickle, step, match)
		default:
			sr.Status = failed
			sr.err = err
			s.storage.mustInsertPickleStepResult(sr)

			s.fmt.Failed(pickle, step, match, err)
		}

		// run after step handlers
		for _, f := range s.afterStepHandlers {
			f(step, err)
		}
	}()

	if undef, err := s.maybeUndefined(step.Text, step.Argument); err != nil {
		return err
	} else if len(undef) > 0 {
		if match != nil {
			match = &StepDefinition{
				args:      match.args,
				hv:        match.hv,
				Expr:      match.Expr,
				Handler:   match.Handler,
				nested:    match.nested,
				undefined: undef,
			}
		}

		sr := newStepResult(pickle.Id, step.Id, match)
		sr.Status = undefined
		s.storage.mustInsertPickleStepResult(sr)

		s.fmt.Undefined(pickle, step, match)
		return ErrUndefined
	}

	if prevStepErr != nil {
		sr := newStepResult(pickle.Id, step.Id, match)
		sr.Status = skipped
		s.storage.mustInsertPickleStepResult(sr)

		s.fmt.Skipped(pickle, step, match)
		return nil
	}

	err = s.maybeSubSteps(match.run())
	return
}

func (s *Suite) maybeUndefined(text string, arg interface{}) ([]string, error) {
	step := s.matchStepText(text)
	if nil == step {
		return []string{text}, nil
	}

	var undefined []string
	if !step.nested {
		return undefined, nil
	}

	if arg != nil {
		step.args = append(step.args, arg)
	}

	for _, next := range step.run().(Steps) {
		lines := strings.Split(next, "\n")
		// @TODO: we cannot currently parse table or content body from nested steps
		if len(lines) > 1 {
			return undefined, fmt.Errorf("nested steps cannot be multiline and have table or content body argument")
		}
		if len(lines[0]) > 0 && lines[0][len(lines[0])-1] == ':' {
			return undefined, fmt.Errorf("nested steps cannot be multiline and have table or content body argument")
		}
		undef, err := s.maybeUndefined(next, nil)
		if err != nil {
			return undefined, err
		}
		undefined = append(undefined, undef...)
	}
	return undefined, nil
}

func (s *Suite) maybeSubSteps(result interface{}) error {
	if nil == result {
		return nil
	}

	if err, ok := result.(error); ok {
		return err
	}

	steps, ok := result.(Steps)
	if !ok {
		return fmt.Errorf("unexpected error, should have been []string: %T - %+v", result, result)
	}

	for _, text := range steps {
		if def := s.matchStepText(text); def == nil {
			return ErrUndefined
		} else if err := s.maybeSubSteps(def.run()); err != nil {
			return fmt.Errorf("%s: %+v", text, err)
		}
	}
	return nil
}

func (s *Suite) matchStepText(text string) *StepDefinition {
	for _, h := range s.steps {
		if m := h.Expr.FindStringSubmatch(text); len(m) > 0 {
			var args []interface{}
			for _, m := range m[1:] {
				args = append(args, m)
			}

			// since we need to assign arguments
			// better to copy the step definition
			return &StepDefinition{
				args:    args,
				hv:      h.hv,
				Expr:    h.Expr,
				Handler: h.Handler,
				nested:  h.nested,
			}
		}
	}
	return nil
}

func (s *Suite) runSteps(pickle *messages.Pickle, steps []*messages.Pickle_PickleStep) (err error) {
	for _, step := range steps {
		stepErr := s.runStep(pickle, step, err)
		switch stepErr {
		case ErrUndefined:
			// do not overwrite failed error
			if err == ErrUndefined || err == nil {
				err = stepErr
			}
		case ErrPending:
			err = stepErr
		case nil:
		default:
			err = stepErr
		}
	}
	return
}

func (s *Suite) shouldFail(err error) bool {
	if err == nil {
		return false
	}

	if err == ErrUndefined || err == ErrPending {
		return s.strict
	}

	return true
}

func (s *Suite) runFeature(f *feature) {
	s.fmt.Feature(f.GherkinDocument, f.Uri, f.content)

	pickles := make([]*messages.Pickle, len(f.pickles))
	if s.randomSeed != 0 {
		r := rand.New(rand.NewSource(s.randomSeed))
		perm := r.Perm(len(f.pickles))
		for i, v := range perm {
			pickles[v] = f.pickles[i]
		}
	} else {
		copy(pickles, f.pickles)
	}

	for _, pickle := range pickles {
		err := s.runPickle(pickle)
		if s.shouldFail(err) {
			s.failed = true
			if s.stopOnFailure {
				return
			}
		}
	}
}

func isEmptyFeature(pickles []*messages.Pickle) bool {
	for _, pickle := range pickles {
		if len(pickle.Steps) > 0 {
			return false
		}
	}

	return true
}

func (s *Suite) runPickle(pickle *messages.Pickle) (err error) {
	if len(pickle.Steps) == 0 {
		pr := pickleResult{PickleID: pickle.Id, StartedAt: timeNowFunc()}
		s.storage.mustInsertPickleResult(pr)

		s.fmt.Pickle(pickle)
		return ErrUndefined
	}

	// run before scenario handlers
	for _, f := range s.beforeScenarioHandlers {
		f(pickle)
	}

	pr := pickleResult{PickleID: pickle.Id, StartedAt: timeNowFunc()}
	s.storage.mustInsertPickleResult(pr)

	s.fmt.Pickle(pickle)

	// scenario
	err = s.runSteps(pickle, pickle.Steps)

	// run after scenario handlers
	for _, f := range s.afterScenarioHandlers {
		f(pickle, err)
	}

	return
}
