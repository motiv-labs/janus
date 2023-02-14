package godog

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

func init() {
	Format("junit", "Prints junit compatible xml to stdout", junitFunc)
}

func junitFunc(suite string, out io.Writer) Formatter {
	return &junitFormatter{basefmt: newBaseFmt(suite, out)}
}

type junitFormatter struct {
	*basefmt
}

func (f *junitFormatter) Summary() {
	suite := f.buildJUNITPackageSuite()

	_, err := io.WriteString(f.out, xml.Header)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to write junit string:", err)
	}

	enc := xml.NewEncoder(f.out)
	enc.Indent("", s(2))
	if err = enc.Encode(suite); err != nil {
		fmt.Fprintln(os.Stderr, "failed to write junit xml:", err)
	}
}

func junitTimeDuration(from, to time.Time) string {
	return strconv.FormatFloat(to.Sub(from).Seconds(), 'f', -1, 64)
}

func (f *junitFormatter) buildJUNITPackageSuite() junitPackageSuite {
	features := f.storage.mustGetFeatures()
	sort.Sort(sortFeaturesByName(features))

	testRunStartedAt := f.storage.mustGetTestRunStarted().StartedAt

	suite := junitPackageSuite{
		Name:       f.suiteName,
		TestSuites: make([]*junitTestSuite, len(features)),
		Time:       junitTimeDuration(testRunStartedAt, timeNowFunc()),
	}

	for idx, feature := range features {
		pickles := f.storage.mustGetPickles(feature.Uri)
		sort.Sort(sortPicklesByID(pickles))

		ts := junitTestSuite{
			Name:      feature.Feature.Name,
			TestCases: make([]*junitTestCase, len(pickles)),
		}

		var testcaseNames = make(map[string]int)
		for _, pickle := range pickles {
			testcaseNames[pickle.Name] = testcaseNames[pickle.Name] + 1
		}

		firstPickleStartedAt := testRunStartedAt
		lastPickleFinishedAt := testRunStartedAt

		var outlineNo = make(map[string]int)
		for idx, pickle := range pickles {
			tc := junitTestCase{}

			pickleResult := f.storage.mustGetPickleResult(pickle.Id)

			if idx == 0 {
				firstPickleStartedAt = pickleResult.StartedAt
			}

			lastPickleFinishedAt = pickleResult.StartedAt

			if len(pickle.Steps) > 0 {
				lastStep := pickle.Steps[len(pickle.Steps)-1]
				lastPickleStepResult := f.storage.mustGetPickleStepResult(lastStep.Id)
				lastPickleFinishedAt = lastPickleStepResult.finishedAt
			}

			tc.Time = junitTimeDuration(pickleResult.StartedAt, lastPickleFinishedAt)

			tc.Name = pickle.Name
			if testcaseNames[tc.Name] > 1 {
				outlineNo[tc.Name] = outlineNo[tc.Name] + 1
				tc.Name += fmt.Sprintf(" #%d", outlineNo[tc.Name])
			}

			ts.Tests++
			suite.Tests++

			pickleStepResults := f.storage.mustGetPickleStepResultsByPickleID(pickle.Id)
			for _, stepResult := range pickleStepResults {
				pickleStep := f.storage.mustGetPickleStep(stepResult.PickleStepID)

				switch stepResult.Status {
				case passed:
					tc.Status = passed.String()
				case failed:
					tc.Status = failed.String()
					tc.Failure = &junitFailure{
						Message: fmt.Sprintf("Step %s: %s", pickleStep.Text, stepResult.err),
					}
				case skipped:
					tc.Error = append(tc.Error, &junitError{
						Type:    "skipped",
						Message: fmt.Sprintf("Step %s", pickleStep.Text),
					})
				case undefined:
					tc.Status = undefined.String()
					tc.Error = append(tc.Error, &junitError{
						Type:    "undefined",
						Message: fmt.Sprintf("Step %s", pickleStep.Text),
					})
				case pending:
					tc.Status = pending.String()
					tc.Error = append(tc.Error, &junitError{
						Type:    "pending",
						Message: fmt.Sprintf("Step %s: TODO: write pending definition", pickleStep.Text),
					})
				}
			}

			switch tc.Status {
			case failed.String():
				ts.Failures++
				suite.Failures++
			case undefined.String(), pending.String():
				ts.Errors++
				suite.Errors++
			}

			ts.TestCases[idx] = &tc
		}

		ts.Time = junitTimeDuration(firstPickleStartedAt, lastPickleFinishedAt)

		suite.TestSuites[idx] = &ts
	}

	return suite
}

type junitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr,omitempty"`
}

type junitError struct {
	XMLName xml.Name `xml:"error,omitempty"`
	Message string   `xml:"message,attr"`
	Type    string   `xml:"type,attr"`
}

type junitTestCase struct {
	XMLName xml.Name      `xml:"testcase"`
	Name    string        `xml:"name,attr"`
	Status  string        `xml:"status,attr"`
	Time    string        `xml:"time,attr"`
	Failure *junitFailure `xml:"failure,omitempty"`
	Error   []*junitError
}

type junitTestSuite struct {
	XMLName   xml.Name `xml:"testsuite"`
	Name      string   `xml:"name,attr"`
	Tests     int      `xml:"tests,attr"`
	Skipped   int      `xml:"skipped,attr"`
	Failures  int      `xml:"failures,attr"`
	Errors    int      `xml:"errors,attr"`
	Time      string   `xml:"time,attr"`
	TestCases []*junitTestCase
}

type junitPackageSuite struct {
	XMLName    xml.Name `xml:"testsuites"`
	Name       string   `xml:"name,attr"`
	Tests      int      `xml:"tests,attr"`
	Skipped    int      `xml:"skipped,attr"`
	Failures   int      `xml:"failures,attr"`
	Errors     int      `xml:"errors,attr"`
	Time       string   `xml:"time,attr"`
	TestSuites []*junitTestSuite
}
