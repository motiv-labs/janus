package godog

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/cucumber/godog/colors"
	"github.com/cucumber/messages-go/v10"
)

const (
	exitSuccess int = iota
	exitFailure
	exitOptionError
)

type initializer func(*Suite)
type testSuiteInitializer func(*TestSuiteContext)
type scenarioInitializer func(*ScenarioContext)

type runner struct {
	randomSeed            int64
	stopOnFailure, strict bool

	features []*feature

	initializer          initializer
	testSuiteInitializer testSuiteInitializer
	scenarioInitializer  scenarioInitializer

	storage *storage
	fmt     Formatter
}

func (r *runner) concurrent(rate int, formatterFn func() Formatter) (failed bool) {
	var useFmtCopy bool
	var copyLock sync.Mutex

	// special mode for concurrent-formatter
	if _, ok := r.fmt.(ConcurrentFormatter); ok {
		useFmtCopy = true
	}

	if fmt, ok := r.fmt.(storageFormatter); ok {
		fmt.setStorage(r.storage)
	}

	testRunStarted := testRunStarted{StartedAt: timeNowFunc()}
	r.storage.mustInsertTestRunStarted(testRunStarted)
	r.fmt.TestRunStarted()

	queue := make(chan int, rate)
	for i, ft := range r.features {
		queue <- i // reserve space in queue
		ft := *ft

		go func(fail *bool, feat *feature) {
			var fmtCopy Formatter

			defer func() {
				<-queue // free a space in queue
			}()

			if r.stopOnFailure && *fail {
				return
			}

			suite := &Suite{
				fmt:           r.fmt,
				randomSeed:    r.randomSeed,
				stopOnFailure: r.stopOnFailure,
				strict:        r.strict,
				features:      []*feature{feat},
				storage:       r.storage,
			}

			if useFmtCopy {
				fmtCopy = formatterFn()
				suite.fmt = fmtCopy

				concurrentDestFmt, dOk := fmtCopy.(ConcurrentFormatter)
				concurrentSourceFmt, sOk := r.fmt.(ConcurrentFormatter)

				if dOk && sOk {
					concurrentDestFmt.Sync(concurrentSourceFmt)
				}
			}

			if fmt, ok := suite.fmt.(storageFormatter); ok {
				fmt.setStorage(r.storage)
			}

			r.initializer(suite)

			suite.run()

			if suite.failed {
				copyLock.Lock()
				*fail = true
				copyLock.Unlock()
			}

			if useFmtCopy {
				copyLock.Lock()

				concurrentDestFmt, dOk := r.fmt.(ConcurrentFormatter)
				concurrentSourceFmt, sOk := fmtCopy.(ConcurrentFormatter)

				if dOk && sOk {
					concurrentDestFmt.Copy(concurrentSourceFmt)
				} else if !dOk {
					panic("cant cast dest formatter to progress-typed")
				} else if !sOk {
					panic("cant cast source formatter to progress-typed")
				}

				copyLock.Unlock()
			}
		}(&failed, &ft)
	}
	// wait until last are processed
	for i := 0; i < rate; i++ {
		queue <- i
	}
	close(queue)

	// print summary
	r.fmt.Summary()
	return
}

func (r *runner) scenarioConcurrent(rate int) (failed bool) {
	var copyLock sync.Mutex

	if fmt, ok := r.fmt.(storageFormatter); ok {
		fmt.setStorage(r.storage)
	}

	testSuiteContext := TestSuiteContext{}
	if r.testSuiteInitializer != nil {
		r.testSuiteInitializer(&testSuiteContext)
	}

	testRunStarted := testRunStarted{StartedAt: timeNowFunc()}
	r.storage.mustInsertTestRunStarted(testRunStarted)
	r.fmt.TestRunStarted()

	// run before suite handlers
	for _, f := range testSuiteContext.beforeSuiteHandlers {
		f()
	}

	queue := make(chan int, rate)
	for _, ft := range r.features {
		pickles := make([]*messages.Pickle, len(ft.pickles))
		if r.randomSeed != 0 {
			r := rand.New(rand.NewSource(r.randomSeed))
			perm := r.Perm(len(ft.pickles))
			for i, v := range perm {
				pickles[v] = ft.pickles[i]
			}
		} else {
			copy(pickles, ft.pickles)
		}

		for i, p := range pickles {
			pickle := *p

			queue <- i // reserve space in queue

			if i == 0 {
				r.fmt.Feature(ft.GherkinDocument, ft.Uri, ft.content)
			}

			go func(fail *bool, pickle *messages.Pickle) {
				defer func() {
					<-queue // free a space in queue
				}()

				if r.stopOnFailure && *fail {
					return
				}

				suite := &Suite{
					fmt:        r.fmt,
					randomSeed: r.randomSeed,
					strict:     r.strict,
					storage:    r.storage,
				}

				if r.scenarioInitializer != nil {
					sc := ScenarioContext{suite: suite}
					r.scenarioInitializer(&sc)
				}

				err := suite.runPickle(pickle)
				if suite.shouldFail(err) {
					copyLock.Lock()
					*fail = true
					copyLock.Unlock()
				}
			}(&failed, &pickle)
		}
	}

	// wait until last are processed
	for i := 0; i < rate; i++ {
		queue <- i
	}

	close(queue)

	// run after suite handlers
	for _, f := range testSuiteContext.afterSuiteHandlers {
		f()
	}

	// print summary
	r.fmt.Summary()
	return
}

// RunWithOptions is same as Run function, except
// it uses Options provided in order to run the
// test suite without parsing flags
//
// This method is useful in case if you run
// godog in for example TestMain function together
// with go tests
//
// The exit codes may vary from:
//  0 - success
//  1 - failed
//  2 - command line usage error
//  128 - or higher, os signal related error exit codes
//
// If there are flag related errors they will
// be directed to os.Stderr
//
// Deprecated: The current Suite initializer will be removed and replaced by
// two initializers, one for the Test Suite and one for the Scenarios.
// Use:
//   godog.TestSuite{
//     Name: name,
//     TestSuiteInitializer: testSuiteInitializer,
//     ScenarioInitializer: scenarioInitializer,
//     Options: &opts,
//   }.Run()
// instead.
func RunWithOptions(suite string, initializer func(*Suite), opt Options) int {
	return runWithOptions(suite, runner{initializer: initializer}, opt)
}

func runWithOptions(suite string, runner runner, opt Options) int {
	var output io.Writer = os.Stdout
	if nil != opt.Output {
		output = opt.Output
	}

	if opt.NoColors {
		output = colors.Uncolored(output)
	} else {
		output = colors.Colored(output)
	}

	if opt.ShowStepDefinitions {
		s := &Suite{}
		runner.initializer(s)
		printStepDefinitions(s.steps, output)
		return exitOptionError
	}

	if len(opt.Paths) == 0 {
		inf, err := os.Stat("features")
		if err == nil && inf.IsDir() {
			opt.Paths = []string{"features"}
		}
	}

	if opt.Concurrency < 1 {
		opt.Concurrency = 1
	}

	formatter := FindFmt(opt.Format)
	if nil == formatter {
		var names []string
		for name := range AvailableFormatters() {
			names = append(names, name)
		}
		fmt.Fprintln(os.Stderr, fmt.Errorf(
			`unregistered formatter name: "%s", use one of: %s`,
			opt.Format,
			strings.Join(names, ", "),
		))
		return exitOptionError
	}
	runner.fmt = formatter(suite, output)

	var err error
	if runner.features, err = parseFeatures(opt.Tags, opt.Paths); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitOptionError
	}

	runner.storage = newStorage()
	for _, feat := range runner.features {
		runner.storage.mustInsertFeature(feat)

		for _, pickle := range feat.pickles {
			runner.storage.mustInsertPickle(pickle)
		}
	}

	// user may have specified -1 option to create random seed
	runner.randomSeed = opt.Randomize
	if runner.randomSeed == -1 {
		runner.randomSeed = makeRandomSeed()
	}

	runner.stopOnFailure = opt.StopOnFailure
	runner.strict = opt.Strict

	// store chosen seed in environment, so it could be seen in formatter summary report
	os.Setenv("GODOG_SEED", strconv.FormatInt(runner.randomSeed, 10))
	// determine tested package
	_, filename, _, _ := runtime.Caller(1)
	os.Setenv("GODOG_TESTED_PACKAGE", runsFromPackage(filename))

	var failed bool
	if runner.initializer != nil {
		failed = runner.concurrent(opt.Concurrency, func() Formatter { return formatter(suite, output) })
	} else {
		failed = runner.scenarioConcurrent(opt.Concurrency)
	}

	// @TODO: should prevent from having these
	os.Setenv("GODOG_SEED", "")
	os.Setenv("GODOG_TESTED_PACKAGE", "")
	if failed && opt.Format != "events" {
		return exitFailure
	}
	return exitSuccess
}

func runsFromPackage(fp string) string {
	dir := filepath.Dir(fp)
	for _, gp := range gopaths {
		gp = filepath.Join(gp, "src")
		if strings.Index(dir, gp) == 0 {
			return strings.TrimLeft(strings.Replace(dir, gp, "", 1), string(filepath.Separator))
		}
	}
	return dir
}

// Run creates and runs the feature suite.
// Reads all configuration options from flags.
// uses contextInitializer to register contexts
//
// the concurrency option allows runner to
// initialize a number of suites to be run
// separately. Only progress formatter
// is supported when concurrency level is
// higher than 1
//
// contextInitializer must be able to register
// the step definitions and event handlers.
//
// The exit codes may vary from:
//  0 - success
//  1 - failed
//  2 - command line usage error
//  128 - or higher, os signal related error exit codes
//
// If there are flag related errors they will
// be directed to os.Stderr
//
// Deprecated: The current Suite initializer will be removed and replaced by
// two initializers, one for the Test Suite and one for the Scenarios.
// Use:
//   godog.TestSuite{
//     Name: name,
//     TestSuiteInitializer: testSuiteInitializer,
//     ScenarioInitializer: scenarioInitializer,
//   }.Run()
// instead.
func Run(suite string, initializer func(*Suite)) int {
	var opt Options
	opt.Output = colors.Colored(os.Stdout)

	flagSet := FlagSet(&opt)
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitOptionError
	}

	opt.Paths = flagSet.Args()

	return RunWithOptions(suite, initializer, opt)
}

// TestSuite allows for configuration
// of the Test Suite Execution
type TestSuite struct {
	Name                 string
	TestSuiteInitializer func(*TestSuiteContext)
	ScenarioInitializer  func(*ScenarioContext)
	Options              *Options
}

// Run will execute the test suite.
//
// If options are not set, it will reads
// all configuration options from flags.
//
// The exit codes may vary from:
//  0 - success
//  1 - failed
//  2 - command line usage error
//  128 - or higher, os signal related error exit codes
//
// If there are flag related errors they will be directed to os.Stderr
func (ts TestSuite) Run() int {
	if ts.Options == nil {
		ts.Options = &Options{}
		ts.Options.Output = colors.Colored(os.Stdout)

		flagSet := FlagSet(ts.Options)
		if err := flagSet.Parse(os.Args[1:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitOptionError
		}

		ts.Options.Paths = flagSet.Args()
	}

	r := runner{testSuiteInitializer: ts.TestSuiteInitializer, scenarioInitializer: ts.ScenarioInitializer}
	return runWithOptions(ts.Name, r, *ts.Options)
}
