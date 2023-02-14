package godog

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/cucumber/messages-go/v10"

	"github.com/cucumber/godog/colors"
)

func init() {
	Format("pretty", "Prints every feature with runtime statuses.", prettyFunc)
}

func prettyFunc(suite string, out io.Writer) Formatter {
	return &pretty{basefmt: newBaseFmt(suite, out)}
}

var outlinePlaceholderRegexp = regexp.MustCompile("<[^>]+>")

// a built in default pretty formatter
type pretty struct {
	*basefmt
	firstFeature *bool
}

func (f *pretty) TestRunStarted() {
	f.basefmt.TestRunStarted()

	f.lock.Lock()
	defer f.lock.Unlock()

	firstFeature := true
	f.firstFeature = &firstFeature
}

func (f *pretty) Feature(gd *messages.GherkinDocument, p string, c []byte) {
	f.lock.Lock()
	if !*f.firstFeature {
		fmt.Fprintln(f.out, "")
	}

	*f.firstFeature = false
	f.lock.Unlock()

	f.basefmt.Feature(gd, p, c)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.printFeature(gd.Feature)
}

// Pickle takes a gherkin node for formatting
func (f *pretty) Pickle(pickle *messages.Pickle) {
	f.basefmt.Pickle(pickle)

	f.lock.Lock()
	defer f.lock.Unlock()

	if len(pickle.Steps) == 0 {
		f.printUndefinedPickle(pickle)
		return
	}
}

func (f *pretty) Passed(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Passed(pickle, step, match)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.printStep(pickle, step)
}

func (f *pretty) Skipped(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Skipped(pickle, step, match)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.printStep(pickle, step)
}

func (f *pretty) Undefined(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Undefined(pickle, step, match)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.printStep(pickle, step)
}

func (f *pretty) Failed(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition, err error) {
	f.basefmt.Failed(pickle, step, match, err)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.printStep(pickle, step)
}

func (f *pretty) Pending(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Pending(pickle, step, match)

	f.lock.Lock()
	defer f.lock.Unlock()

	f.printStep(pickle, step)
}

func (f *pretty) printFeature(feature *messages.GherkinDocument_Feature) {
	fmt.Fprintln(f.out, keywordAndName(feature.Keyword, feature.Name))
	if strings.TrimSpace(feature.Description) != "" {
		for _, line := range strings.Split(feature.Description, "\n") {
			fmt.Fprintln(f.out, s(f.indent)+strings.TrimSpace(line))
		}
	}
}

func keywordAndName(keyword, name string) string {
	title := whiteb(keyword + ":")
	if len(name) > 0 {
		title += " " + name
	}
	return title
}

func (f *pretty) scenarioLengths(pickle *messages.Pickle) (scenarioHeaderLength int, maxLength int) {
	feature := f.storage.mustGetFeature(pickle.Uri)
	astScenario := feature.findScenario(pickle.AstNodeIds[0])
	astBackground := feature.findBackground(pickle.AstNodeIds[0])

	scenarioHeaderLength = f.lengthPickle(astScenario.Keyword, astScenario.Name)
	maxLength = f.longestStep(astScenario.Steps, scenarioHeaderLength)

	if astBackground != nil {
		maxLength = f.longestStep(astBackground.Steps, maxLength)
	}

	return scenarioHeaderLength, maxLength
}

func (f *pretty) printScenarioHeader(pickle *messages.Pickle, astScenario *messages.GherkinDocument_Feature_Scenario, spaceFilling int) {
	feature := f.storage.mustGetFeature(pickle.Uri)
	text := s(f.indent) + keywordAndName(astScenario.Keyword, astScenario.Name)
	text += s(spaceFilling) + line(feature.Uri, astScenario.Location)
	fmt.Fprintln(f.out, "\n"+text)
}

func (f *pretty) printUndefinedPickle(pickle *messages.Pickle) {
	feature := f.storage.mustGetFeature(pickle.Uri)
	astScenario := feature.findScenario(pickle.AstNodeIds[0])
	astBackground := feature.findBackground(pickle.AstNodeIds[0])

	scenarioHeaderLength, maxLength := f.scenarioLengths(pickle)

	if astBackground != nil {
		fmt.Fprintln(f.out, "\n"+s(f.indent)+keywordAndName(astBackground.Keyword, astBackground.Name))
		for _, step := range astBackground.Steps {
			text := s(f.indent*2) + cyan(strings.TrimSpace(step.Keyword)) + " " + cyan(step.Text)
			fmt.Fprintln(f.out, text)
		}
	}

	//  do not print scenario headers and examples multiple times
	if len(astScenario.Examples) > 0 {
		exampleTable, exampleRow := feature.findExample(pickle.AstNodeIds[1])
		firstExampleRow := exampleTable.TableBody[0].Id == exampleRow.Id
		firstExamplesTable := astScenario.Examples[0].Location.Line == exampleTable.Location.Line

		if !(firstExamplesTable && firstExampleRow) {
			return
		}
	}

	f.printScenarioHeader(pickle, astScenario, maxLength-scenarioHeaderLength)

	for _, examples := range astScenario.Examples {
		max := longestExampleRow(examples, cyan, cyan)

		fmt.Fprintln(f.out, "")
		fmt.Fprintln(f.out, s(f.indent*2)+keywordAndName(examples.Keyword, examples.Name))

		f.printTableHeader(examples.TableHeader, max)

		for _, row := range examples.TableBody {
			f.printTableRow(row, max, cyan)
		}
	}
}

// Summary sumarize the feature formatter output
func (f *pretty) Summary() {
	failedStepResults := f.storage.mustGetPickleStepResultsByStatus(failed)
	if len(failedStepResults) > 0 {
		fmt.Fprintln(f.out, "\n--- "+red("Failed steps:")+"\n")

		sort.Sort(sortPickleStepResultsByPickleStepID(failedStepResults))

		for _, fail := range failedStepResults {
			pickle := f.storage.mustGetPickle(fail.PickleID)
			pickleStep := f.storage.mustGetPickleStep(fail.PickleStepID)
			feature := f.storage.mustGetFeature(pickle.Uri)

			astScenario := feature.findScenario(pickle.AstNodeIds[0])
			scenarioDesc := fmt.Sprintf("%s: %s", astScenario.Keyword, pickle.Name)

			astStep := feature.findStep(pickleStep.AstNodeIds[0])
			stepDesc := strings.TrimSpace(astStep.Keyword) + " " + pickleStep.Text

			fmt.Fprintln(f.out, s(f.indent)+red(scenarioDesc)+line(feature.Uri, astScenario.Location))
			fmt.Fprintln(f.out, s(f.indent*2)+red(stepDesc)+line(feature.Uri, astStep.Location))
			fmt.Fprintln(f.out, s(f.indent*3)+red("Error: ")+redb(fmt.Sprintf("%+v", fail.err))+"\n")
		}
	}

	f.basefmt.Summary()
}

func (f *pretty) printOutlineExample(pickle *messages.Pickle, backgroundSteps int) {
	var errorMsg string
	var clr = green

	feature := f.storage.mustGetFeature(pickle.Uri)
	astScenario := feature.findScenario(pickle.AstNodeIds[0])
	scenarioHeaderLength, maxLength := f.scenarioLengths(pickle)

	exampleTable, exampleRow := feature.findExample(pickle.AstNodeIds[1])
	printExampleHeader := exampleTable.TableBody[0].Id == exampleRow.Id
	firstExamplesTable := astScenario.Examples[0].Location.Line == exampleTable.Location.Line

	pickleStepResults := f.storage.mustGetPickleStepResultsByPickleID(pickle.Id)

	firstExecutedScenarioStep := len(pickleStepResults) == backgroundSteps+1
	if firstExamplesTable && printExampleHeader && firstExecutedScenarioStep {
		f.printScenarioHeader(pickle, astScenario, maxLength-scenarioHeaderLength)
	}

	if len(exampleTable.TableBody) == 0 {
		// do not print empty examples
		return
	}

	lastStep := len(pickleStepResults) == len(pickle.Steps)
	if !lastStep {
		// do not print examples unless all steps has finished
		return
	}

	for _, result := range pickleStepResults {
		// determine example row status
		switch {
		case result.Status == failed:
			errorMsg = result.err.Error()
			clr = result.Status.clr()
		case result.Status == undefined || result.Status == pending:
			clr = result.Status.clr()
		case result.Status == skipped && clr == nil:
			clr = cyan
		}

		if firstExamplesTable && printExampleHeader {
			// in first example, we need to print steps

			pickleStep := f.storage.mustGetPickleStep(result.PickleStepID)
			astStep := feature.findStep(pickleStep.AstNodeIds[0])

			var text = ""
			if result.def != nil {
				if m := outlinePlaceholderRegexp.FindAllStringIndex(astStep.Text, -1); len(m) > 0 {
					var pos int
					for i := 0; i < len(m); i++ {
						pair := m[i]
						text += cyan(astStep.Text[pos:pair[0]])
						text += cyanb(astStep.Text[pair[0]:pair[1]])
						pos = pair[1]
					}
					text += cyan(astStep.Text[pos:len(astStep.Text)])
				} else {
					text = cyan(astStep.Text)
				}

				_, maxLength := f.scenarioLengths(pickle)
				stepLength := f.lengthPickleStep(astStep.Keyword, astStep.Text)

				text += s(maxLength - stepLength)
				text += " " + blackb("# "+result.def.definitionID())
			}

			// print the step outline
			fmt.Fprintln(f.out, s(f.indent*2)+cyan(strings.TrimSpace(astStep.Keyword))+" "+text)

			if table := pickleStep.Argument.GetDataTable(); table != nil {
				f.printTable(table, cyan)
			}

			if docString := astStep.GetDocString(); docString != nil {
				f.printDocString(docString)
			}
		}
	}

	max := longestExampleRow(exampleTable, clr, cyan)

	// an example table header
	if printExampleHeader {
		fmt.Fprintln(f.out, "")
		fmt.Fprintln(f.out, s(f.indent*2)+keywordAndName(exampleTable.Keyword, exampleTable.Name))

		f.printTableHeader(exampleTable.TableHeader, max)
	}

	f.printTableRow(exampleRow, max, clr)

	if errorMsg != "" {
		fmt.Fprintln(f.out, s(f.indent*4)+redb(errorMsg))
	}
}

func (f *pretty) printTableRow(row *messages.GherkinDocument_Feature_TableRow, max []int, clr colors.ColorFunc) {
	cells := make([]string, len(row.Cells))

	for i, cell := range row.Cells {
		val := clr(cell.Value)
		ln := utf8.RuneCountInString(val)
		cells[i] = val + s(max[i]-ln)
	}

	fmt.Fprintln(f.out, s(f.indent*3)+"| "+strings.Join(cells, " | ")+" |")
}

func (f *pretty) printTableHeader(row *messages.GherkinDocument_Feature_TableRow, max []int) {
	f.printTableRow(row, max, cyan)
}

func (f *pretty) printStep(pickle *messages.Pickle, pickleStep *messages.Pickle_PickleStep) {
	feature := f.storage.mustGetFeature(pickle.Uri)
	astBackground := feature.findBackground(pickle.AstNodeIds[0])
	astScenario := feature.findScenario(pickle.AstNodeIds[0])
	astStep := feature.findStep(pickleStep.AstNodeIds[0])

	var astBackgroundStep bool
	var firstExecutedBackgroundStep bool
	var backgroundSteps int
	if astBackground != nil {
		backgroundSteps = len(astBackground.Steps)

		for idx, step := range astBackground.Steps {
			if step.Id == pickleStep.AstNodeIds[0] {
				astBackgroundStep = true
				firstExecutedBackgroundStep = idx == 0
				break
			}
		}
	}

	firstPickle := feature.pickles[0].Id == pickle.Id

	if astBackgroundStep && !firstPickle {
		return
	}

	if astBackgroundStep && firstExecutedBackgroundStep {
		fmt.Fprintln(f.out, "\n"+s(f.indent)+keywordAndName(astBackground.Keyword, astBackground.Name))
	}

	if !astBackgroundStep && len(astScenario.Examples) > 0 {
		f.printOutlineExample(pickle, backgroundSteps)
		return
	}

	scenarioHeaderLength, maxLength := f.scenarioLengths(pickle)
	stepLength := f.lengthPickleStep(astStep.Keyword, pickleStep.Text)

	firstExecutedScenarioStep := astScenario.Steps[0].Id == pickleStep.AstNodeIds[0]
	if !astBackgroundStep && firstExecutedScenarioStep {
		f.printScenarioHeader(pickle, astScenario, maxLength-scenarioHeaderLength)
	}

	pickleStepResult := f.storage.mustGetPickleStepResult(pickleStep.Id)
	text := s(f.indent*2) + pickleStepResult.Status.clr()(strings.TrimSpace(astStep.Keyword)) + " " + pickleStepResult.Status.clr()(pickleStep.Text)
	if pickleStepResult.def != nil {
		text += s(maxLength - stepLength + 1)
		text += blackb("# " + pickleStepResult.def.definitionID())
	}
	fmt.Fprintln(f.out, text)

	if table := pickleStep.Argument.GetDataTable(); table != nil {
		f.printTable(table, cyan)
	}

	if docString := astStep.GetDocString(); docString != nil {
		f.printDocString(docString)
	}

	if pickleStepResult.err != nil {
		fmt.Fprintln(f.out, s(f.indent*2)+redb(fmt.Sprintf("%+v", pickleStepResult.err)))
	}

	if pickleStepResult.Status == pending {
		fmt.Fprintln(f.out, s(f.indent*3)+yellow("TODO: write pending definition"))
	}
}

func (f *pretty) printDocString(docString *messages.GherkinDocument_Feature_Step_DocString) {
	var ct string

	if len(docString.MediaType) > 0 {
		ct = " " + cyan(docString.MediaType)
	}

	fmt.Fprintln(f.out, s(f.indent*3)+cyan(docString.Delimiter)+ct)

	for _, ln := range strings.Split(docString.Content, "\n") {
		fmt.Fprintln(f.out, s(f.indent*3)+cyan(ln))
	}

	fmt.Fprintln(f.out, s(f.indent*3)+cyan(docString.Delimiter))
}

// print table with aligned table cells
// @TODO: need to make example header cells bold
func (f *pretty) printTable(t *messages.PickleStepArgument_PickleTable, c colors.ColorFunc) {
	maxColLengths := maxColLengths(t, c)
	var cols = make([]string, len(t.Rows[0].Cells))

	for _, row := range t.Rows {
		for i, cell := range row.Cells {
			val := c(cell.Value)
			colLength := utf8.RuneCountInString(val)
			cols[i] = val + s(maxColLengths[i]-colLength)
		}

		fmt.Fprintln(f.out, s(f.indent*3)+"| "+strings.Join(cols, " | ")+" |")
	}
}

// longest gives a list of longest columns of all rows in Table
func maxColLengths(t *messages.PickleStepArgument_PickleTable, clrs ...colors.ColorFunc) []int {
	if t == nil {
		return []int{}
	}

	longest := make([]int, len(t.Rows[0].Cells))
	for _, row := range t.Rows {
		for i, cell := range row.Cells {
			for _, c := range clrs {
				ln := utf8.RuneCountInString(c(cell.Value))
				if longest[i] < ln {
					longest[i] = ln
				}
			}

			ln := utf8.RuneCountInString(cell.Value)
			if longest[i] < ln {
				longest[i] = ln
			}
		}
	}

	return longest
}

func longestExampleRow(t *messages.GherkinDocument_Feature_Scenario_Examples, clrs ...colors.ColorFunc) []int {
	if t == nil {
		return []int{}
	}

	longest := make([]int, len(t.TableHeader.Cells))
	for i, cell := range t.TableHeader.Cells {
		for _, c := range clrs {
			ln := utf8.RuneCountInString(c(cell.Value))
			if longest[i] < ln {
				longest[i] = ln
			}
		}

		ln := utf8.RuneCountInString(cell.Value)
		if longest[i] < ln {
			longest[i] = ln
		}
	}

	for _, row := range t.TableBody {
		for i, cell := range row.Cells {
			for _, c := range clrs {
				ln := utf8.RuneCountInString(c(cell.Value))
				if longest[i] < ln {
					longest[i] = ln
				}
			}

			ln := utf8.RuneCountInString(cell.Value)
			if longest[i] < ln {
				longest[i] = ln
			}
		}
	}

	return longest
}

func (f *pretty) longestStep(steps []*messages.GherkinDocument_Feature_Step, pickleLength int) int {
	max := pickleLength

	for _, step := range steps {
		length := f.lengthPickleStep(step.Keyword, step.Text)
		if length > max {
			max = length
		}
	}

	return max
}

// a line number representation in feature file
func line(path string, loc *messages.Location) string {
	return " " + blackb(fmt.Sprintf("# %s:%d", path, loc.Line))
}

func (f *pretty) lengthPickleStep(keyword, text string) int {
	return f.indent*2 + utf8.RuneCountInString(strings.TrimSpace(keyword)+" "+text)
}

func (f *pretty) lengthPickle(keyword, name string) int {
	return f.indent + utf8.RuneCountInString(strings.TrimSpace(keyword)+": "+name)
}
