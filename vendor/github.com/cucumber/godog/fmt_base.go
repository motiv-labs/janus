package godog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/cucumber/messages-go/v10"

	"github.com/cucumber/godog/colors"
)

func baseFmtFunc(suite string, out io.Writer) Formatter {
	return newBaseFmt(suite, out)
}

func newBaseFmt(suite string, out io.Writer) *basefmt {
	return &basefmt{
		suiteName: suite,
		indent:    2,
		out:       out,
		lock:      new(sync.Mutex),
	}
}

type basefmt struct {
	suiteName string
	out       io.Writer
	indent    int

	storage *storage
	lock    *sync.Mutex
}

func (f *basefmt) setStorage(st *storage) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.storage = st
}

func (f *basefmt) TestRunStarted()                                                        {}
func (f *basefmt) Feature(ft *messages.GherkinDocument, p string, c []byte)               {}
func (f *basefmt) Pickle(p *messages.Pickle)                                              {}
func (f *basefmt) Defined(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition) {}
func (f *basefmt) Passed(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
}
func (f *basefmt) Skipped(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
}
func (f *basefmt) Undefined(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
}
func (f *basefmt) Failed(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition, err error) {
}
func (f *basefmt) Pending(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
}

func (f *basefmt) Summary() {
	var totalSc, passedSc, undefinedSc int
	var totalSt, passedSt, failedSt, skippedSt, pendingSt, undefinedSt int

	pickleResults := f.storage.mustGetPickleResults()
	for _, pr := range pickleResults {
		var prStatus stepResultStatus
		totalSc++

		pickleStepResults := f.storage.mustGetPickleStepResultsByPickleID(pr.PickleID)

		if len(pickleStepResults) == 0 {
			prStatus = undefined
		}

		for _, sr := range pickleStepResults {
			totalSt++

			switch sr.Status {
			case passed:
				prStatus = passed
				passedSt++
			case failed:
				prStatus = failed
				failedSt++
			case skipped:
				skippedSt++
			case undefined:
				prStatus = undefined
				undefinedSt++
			case pending:
				prStatus = pending
				pendingSt++
			}
		}

		if prStatus == passed {
			passedSc++
		} else if prStatus == undefined {
			undefinedSc++
		}
	}

	var steps, parts, scenarios []string
	if passedSt > 0 {
		steps = append(steps, green(fmt.Sprintf("%d passed", passedSt)))
	}
	if failedSt > 0 {
		parts = append(parts, red(fmt.Sprintf("%d failed", failedSt)))
		steps = append(steps, red(fmt.Sprintf("%d failed", failedSt)))
	}
	if pendingSt > 0 {
		parts = append(parts, yellow(fmt.Sprintf("%d pending", pendingSt)))
		steps = append(steps, yellow(fmt.Sprintf("%d pending", pendingSt)))
	}
	if undefinedSt > 0 {
		parts = append(parts, yellow(fmt.Sprintf("%d undefined", undefinedSc)))
		steps = append(steps, yellow(fmt.Sprintf("%d undefined", undefinedSt)))
	} else if undefinedSc > 0 {
		// there may be some scenarios without steps
		parts = append(parts, yellow(fmt.Sprintf("%d undefined", undefinedSc)))
	}
	if skippedSt > 0 {
		steps = append(steps, cyan(fmt.Sprintf("%d skipped", skippedSt)))
	}
	if passedSc > 0 {
		scenarios = append(scenarios, green(fmt.Sprintf("%d passed", passedSc)))
	}
	scenarios = append(scenarios, parts...)

	testRunStartedAt := f.storage.mustGetTestRunStarted().StartedAt
	elapsed := timeNowFunc().Sub(testRunStartedAt)

	fmt.Fprintln(f.out, "")

	if totalSc == 0 {
		fmt.Fprintln(f.out, "No scenarios")
	} else {
		fmt.Fprintln(f.out, fmt.Sprintf("%d scenarios (%s)", totalSc, strings.Join(scenarios, ", ")))
	}

	if totalSt == 0 {
		fmt.Fprintln(f.out, "No steps")
	} else {
		fmt.Fprintln(f.out, fmt.Sprintf("%d steps (%s)", totalSt, strings.Join(steps, ", ")))
	}

	elapsedString := elapsed.String()
	if elapsed.Nanoseconds() == 0 {
		// go 1.5 and 1.6 prints 0 instead of 0s, if duration is zero.
		elapsedString = "0s"
	}
	fmt.Fprintln(f.out, elapsedString)

	// prints used randomization seed
	seed, err := strconv.ParseInt(os.Getenv("GODOG_SEED"), 10, 64)
	if err == nil && seed != 0 {
		fmt.Fprintln(f.out, "")
		fmt.Fprintln(f.out, "Randomized with seed:", colors.Yellow(seed))
	}

	if text := f.snippets(); text != "" {
		fmt.Fprintln(f.out, "")
		fmt.Fprintln(f.out, yellow("You can implement step definitions for undefined steps with these snippets:"))
		fmt.Fprintln(f.out, yellow(text))
	}
}

func (f *basefmt) snippets() string {
	undefinedStepResults := f.storage.mustGetPickleStepResultsByStatus(undefined)
	if len(undefinedStepResults) == 0 {
		return ""
	}

	var index int
	var snips []undefinedSnippet
	// build snippets
	for _, u := range undefinedStepResults {
		pickleStep := f.storage.mustGetPickleStep(u.PickleStepID)

		steps := []string{pickleStep.Text}
		arg := pickleStep.Argument
		if u.def != nil {
			steps = u.def.undefined
			arg = nil
		}
		for _, step := range steps {
			expr := snippetExprCleanup.ReplaceAllString(step, "\\$1")
			expr = snippetNumbers.ReplaceAllString(expr, "(\\d+)")
			expr = snippetExprQuoted.ReplaceAllString(expr, "$1\"([^\"]*)\"$2")
			expr = "^" + strings.TrimSpace(expr) + "$"

			name := snippetNumbers.ReplaceAllString(step, " ")
			name = snippetExprQuoted.ReplaceAllString(name, " ")
			name = strings.TrimSpace(snippetMethodName.ReplaceAllString(name, ""))
			var words []string
			for i, w := range strings.Split(name, " ") {
				switch {
				case i != 0:
					w = strings.Title(w)
				case len(w) > 0:
					w = string(unicode.ToLower(rune(w[0]))) + w[1:]
				}
				words = append(words, w)
			}
			name = strings.Join(words, "")
			if len(name) == 0 {
				index++
				name = fmt.Sprintf("StepDefinitioninition%d", index)
			}

			var found bool
			for _, snip := range snips {
				if snip.Expr == expr {
					found = true
					break
				}
			}
			if !found {
				snips = append(snips, undefinedSnippet{Method: name, Expr: expr, argument: arg})
			}
		}
	}

	sort.Sort(snippetSortByMethod(snips))

	var buf bytes.Buffer
	if err := undefinedSnippetsTpl.Execute(&buf, snips); err != nil {
		panic(err)
	}
	// there may be trailing spaces
	return strings.Replace(buf.String(), " \n", "\n", -1)
}
