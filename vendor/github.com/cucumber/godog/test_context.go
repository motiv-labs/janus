package godog

import "github.com/cucumber/messages-go/v10"

// Scenario represents the executed scenario
type Scenario = messages.Pickle

// Step represents the executed step
type Step = messages.Pickle_PickleStep

// DocString represents the DocString argument made to a step definition
type DocString = messages.PickleStepArgument_PickleDocString

// Table represents the Table argument made to a step definition
type Table = messages.PickleStepArgument_PickleTable

// TestSuiteContext allows various contexts
// to register event handlers.
//
// When running a test suite, the instance of TestSuiteContext
// is passed to all functions (contexts), which
// have it as a first and only argument.
//
// Note that all event hooks does not catch panic errors
// in order to have a trace information
type TestSuiteContext struct {
	beforeSuiteHandlers []func()
	afterSuiteHandlers  []func()
}

// BeforeSuite registers a function or method
// to be run once before suite runner.
//
// Use it to prepare the test suite for a spin.
// Connect and prepare database for instance...
func (ctx *TestSuiteContext) BeforeSuite(fn func()) {
	ctx.beforeSuiteHandlers = append(ctx.beforeSuiteHandlers, fn)
}

// AfterSuite registers a function or method
// to be run once after suite runner
func (ctx *TestSuiteContext) AfterSuite(fn func()) {
	ctx.afterSuiteHandlers = append(ctx.afterSuiteHandlers, fn)
}

// ScenarioContext allows various contexts
// to register steps and event handlers.
//
// When running a scenario, the instance of ScenarioContext
// is passed to all functions (contexts), which
// have it as a first and only argument.
//
// Note that all event hooks does not catch panic errors
// in order to have a trace information. Only step
// executions are catching panic error since it may
// be a context specific error.
type ScenarioContext struct {
	suite *Suite
}

// BeforeScenario registers a function or method
// to be run before every scenario.
//
// It is a good practice to restore the default state
// before every scenario so it would be isolated from
// any kind of state.
func (ctx *ScenarioContext) BeforeScenario(fn func(sc *Scenario)) {
	ctx.suite.BeforeScenario(fn)
}

// AfterScenario registers an function or method
// to be run after every scenario.
func (ctx *ScenarioContext) AfterScenario(fn func(sc *Scenario, err error)) {
	ctx.suite.AfterScenario(fn)
}

// BeforeStep registers a function or method
// to be run before every step.
func (ctx *ScenarioContext) BeforeStep(fn func(st *Step)) {
	ctx.suite.BeforeStep(fn)
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
func (ctx *ScenarioContext) AfterStep(fn func(st *Step, err error)) {
	ctx.suite.AfterStep(fn)
}

// Step allows to register a *StepDefinition in the
// Godog feature suite, the definition will be applied
// to all steps matching the given Regexp expr.
//
// It will panic if expr is not a valid regular
// expression or stepFunc is not a valid step
// handler.
//
// The expression can be of type: *regexp.Regexp, string or []byte
//
// The stepFunc may accept one or several arguments of type:
// - int, int8, int16, int32, int64
// - float32, float64
// - string
// - []byte
// - *godog.DocString
// - *godog.Table
//
// The stepFunc need to return either an error or []string for multistep
//
// Note that if there are two definitions which may match
// the same step, then only the first matched handler
// will be applied.
//
// If none of the *StepDefinition is matched, then
// ErrUndefined error will be returned when
// running steps.
func (ctx *ScenarioContext) Step(expr, stepFunc interface{}) {
	ctx.suite.Step(expr, stepFunc)
}
