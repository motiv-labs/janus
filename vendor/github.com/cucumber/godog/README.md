[![CircleCI](https://circleci.com/gh/cucumber/godog/tree/master.svg?style=svg)](https://circleci.com/gh/cucumber/godog/tree/master)
[![GoDoc](https://godoc.org/github.com/cucumber/godog?status.svg)](https://godoc.org/github.com/cucumber/godog)
[![codecov](https://codecov.io/gh/cucumber/godog/branch/master/graph/badge.svg)](https://codecov.io/gh/cucumber/godog)

# Godog

<p align="center"><img src="/logo.png" alt="Godog logo" style="width:250px;" /></p>

**The API is likely to change a few times before we reach 1.0.0**

Please read the full README, you may find it very useful. And do not forget to peek into the [Release Notes](https://github.com/cucumber/godog/blob/master/release-notes) and the [CHANGELOG](https://github.com/cucumber/godog/blob/master/CHANGELOG.md) from time to time.

Package godog is the official Cucumber BDD framework for Golang, it merges specification and test documentation into one cohesive whole, using Gherkin formatted scenarios in the format of Given, When, Then.

**Godog** does not intervene with the standard **go test** command behavior. You can leverage both frameworks to functionally test your application while maintaining all test related source code in **_test.go** files.

**Godog** acts similar compared to **go test** command, by using go compiler and linker tool in order to produce test executable. Godog contexts need to be exported the same way as **Test** functions for go tests. Note, that if you use **godog** command tool, it will use `go` executable to determine compiler and linker.

The project was inspired by [behat][behat] and [cucumber][cucumber].

## Why Godog/Cucumber

### A single source of truth

Godog merges specification and test documentation into one cohesive whole.

### Living documentation

Because they're automatically tested by Godog, your specifications are
always bang up-to-date.

### Focus on the customer

Business and IT don't always understand each other. Godog's executable specifications encourage closer collaboration, helping teams keep the business goal in mind at all times.

### Less rework

When automated testing is this much fun, teams can easily protect themselves from costly regressions.

### Read more
- [Behaviour-Driven Development](https://cucumber.io/docs/bdd/)
- [Gherkin Reference](https://cucumber.io/docs/gherkin/reference/)

## Install
```
go get github.com/cucumber/godog/cmd/godog@v0.10.0
```
Adding `@v0.10.0` will install v0.10.0 specifically instead of master.

Running `within the $GOPATH`, you would also need to set `GO111MODULE=on`, like this:
```
GO111MODULE=on go get github.com/cucumber/godog/cmd/godog@v0.10.0
```

## Contributions

Godog is a community driven Open Source Project within the Cucumber organization, it is maintained by a handfull of developers, but we appreciate contributions from everyone.

If you are interested in developing Godog, we suggest you to visit one of our slack channels.

Feel free to open a pull request. Note, if you wish to contribute larger changes or an extension to the exported methods or types, please open an issue before and visit us in slack to discuss the changes.

Reach out to the community on our [Cucumber Slack Community](https://cucumberbdd.slack.com/).
Join [here](https://cucumberbdd-slack-invite.herokuapp.com/).

### Popular Cucumber Slack channels for Godog:
- [#help-godog](https://cucumberbdd.slack.com/archives/CTNL1JCVA) - General Godog Adoption Help
- [#committers-go](https://cucumberbdd.slack.com/archives/CA5NJPDJ4) - Golang focused Cucumber Contributors
- [#committers](https://cucumberbdd.slack.com/archives/C62D0FK0E) - General Cucumber Contributors

## Example

The following example can be [found here](/_examples/godogs).

### Step 1

Given we create a new go package **$GOPATH/src/godogs**. From now on, this is our work directory `cd $GOPATH/src/godogs`.

Imagine we have a **godog cart** to serve godogs for lunch. First of all, we describe our feature in plain text - `vim $GOPATH/src/godogs/features/godogs.feature`:

``` gherkin
# file: $GOPATH/src/godogs/features/godogs.feature
Feature: eat godogs
  In order to be happy
  As a hungry gopher
  I need to be able to eat godogs

  Scenario: Eat 5 out of 12
    Given there are 12 godogs
    When I eat 5
    Then there should be 7 remaining
```

**NOTE:** same as **go test** godog respects package level isolation. All your step definitions should be in your tested package root directory. In this case - `$GOPATH/src/godogs`

### Step 2

If godog is installed in your GOPATH. We can run `godog` inside the **$GOPATH/src/godogs** directory. You should see that the steps are undefined:

![Undefined step snippets](/screenshots/undefined.png?raw=true)

If we wish to vendor godog dependency, we can do it as usual, using tools you prefer:

    git clone https://github.com/cucumber/godog.git $GOPATH/src/godogs/vendor/github.com/cucumber/godog

It gives you undefined step snippets to implement in your test context. You may copy these snippets into your `godogs_test.go` file.

Our directory structure should now look like:

![Directory layout](/screenshots/dir-tree.png?raw=true)

If you copy the snippets into our test file and run godog again. We should see the step definition is now pending:

![Pending step definition](/screenshots/pending.png?raw=true)

You may change **ErrPending** to **nil** and the scenario will pass successfully.

Since we need a working implementation, we may start by implementing only what is necessary.

### Step 3

We only need a number of **godogs** for now. Lets keep it simple.

``` go
package main

// Godogs available to eat
var Godogs int

func main() { /* usual main func */ }
```

### Step 4

Now lets implement our step definitions, which we can copy from generated console output snippets in order to test our feature requirements:

``` go
package main

import (
	"fmt"
	
	messages "github.com/cucumber/messages-go/v10" // needed for godog v0.9.0 and earlier

	"github.com/cucumber/godog"
)

func thereAreGodogs(available int) error {
	Godogs = available
	return nil
}

func iEat(num int) error {
	if Godogs < num {
		return fmt.Errorf("you cannot eat %d godogs, there are %d available", num, Godogs)
	}
	Godogs -= num
	return nil
}

func thereShouldBeRemaining(remaining int) error {
	if Godogs != remaining {
		return fmt.Errorf("expected %d godogs to be remaining, but there is %d", remaining, Godogs)
	}
	return nil
}

// godog v0.9.0 and earlier
func FeatureContext(s *godog.Suite) {
	s.BeforeSuite(func() { Godogs = 0 })

	s.BeforeScenario(func(*messages.Pickle) {
		Godogs = 0 // clean the state before every scenario
	})

	s.Step(`^there are (\d+) godogs$`, thereAreGodogs)
	s.Step(`^I eat (\d+)$`, iEat)
	s.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}

// godog v0.10.0 (latest)
func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() { Godogs = 0 })
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
		Godogs = 0 // clean the state before every scenario
	})

	ctx.Step(`^there are (\d+) godogs$`, thereAreGodogs)
	ctx.Step(`^I eat (\d+)$`, iEat)
	ctx.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}
```

Now when you run the `godog` again, you should see:

![Passed suite](/screenshots/passed.png?raw=true)

We have hooked to **BeforeScenario** event in order to reset application state before each scenario. You may hook into more events, like **AfterStep** to print all state in case of an error. Or **BeforeSuite** to prepare a database.

By now, you should have figured out, how to use **godog**. Another advice is to make steps orthogonal, small and simple to read for a user. Whether the user is a dumb website user or an API developer, who may understand a little more technical context - it should target that user.

When steps are orthogonal and small, you can combine them just like you do with Unix tools. Look how to simplify or remove ones, which can be composed.

### References and Tutorials

- [cucumber-html-reporter](https://github.com/gkushang/cucumber-html-reporter),
  may be used in order to generate **html** reports together with **cucumber** output formatter. See the [following docker image](https://github.com/myie/cucumber-html-reporter) for usage details.
- [how to use godog by semaphoreci](https://semaphoreci.com/community/tutorials/how-to-use-godog-for-behavior-driven-development-in-go)
- see [examples](https://github.com/cucumber/godog/tree/master/_examples)
- see extension [AssistDog](https://github.com/hellomd/assistdog),
  which may have useful **gherkin.DataTable** transformations or comparison methods for assertions.

### Documentation

See [pkg documentation][godoc] for general API details.
See **[Circle Config](/.circleci/config.yml)** for supported **go** versions.
See `godog -h` for general command options.

See implementation examples:

- [rest API server](/_examples/api)
- [rest API with Database](/_examples/db)
- [godogs](/_examples/godogs)

## FAQ

### Running Godog with go test

You may integrate running **godog** in your **go test** command. You can run it using go [TestMain](https://golang.org/pkg/testing/#hdr-Main) func available since **go 1.4**. In this case it is not necessary to have **godog** command installed. See the following examples.

The following example binds **godog** flags with specified prefix `godog` in order to prevent flag collisions.

``` go
var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = flag.Args()

	// godog v0.9.0 and earlier
	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opts)

	// godog v0.10.0 (latest)
	status := godog.TestSuite{
		Name: "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
```

Then you may run tests with by specifying flags in order to filter features.

```
go test -v --godog.random --godog.tags=wip
go test -v --godog.format=pretty --godog.random -race -coverprofile=coverage.txt -covermode=atomic
```

The following example does not bind godog flags, instead manually configuring needed options.

``` go
func TestMain(m *testing.M) {
	opts := godog.Options{
		Format:    "progress",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	}

	// godog v0.9.0 and earlier
	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opts)

	// godog v0.10.0 (latest)
	status := godog.TestSuite{
		Name: "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
```

You can even go one step further and reuse **go test** flags, like **verbose** mode in order to switch godog **format**. See the following example:

``` go
func TestMain(m *testing.M) {
	format := "progress"
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}

	opts := godog.Options{
		Format: format,
		Paths:     []string{"features"},
	}

	// godog v0.9.0 and earlier
	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opts)

	// godog v0.10.0 (latest)
	status := godog.TestSuite{
		Name: "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
```

Now when running `go test -v` it will use **pretty** format.

### Tags

If you want to filter scenarios by tags, you can use the `-t=<expression>` or `--tags=<expression>` where `<expression>` is one of the following:

- `@wip` - run all scenarios with wip tag
- `~@wip` - exclude all scenarios with wip tag
- `@wip && ~@new` - run wip scenarios, but exclude new
- `@wip,@undone` - run wip or undone scenarios

### Using assertion packages like testify with Godog
A more extensive example can be [found here](/_examples/assert-godogs).

``` go
func thereShouldBeRemaining(remaining int) error {
	return assertExpectedAndActual(
		assert.Equal, Godogs, remaining,
		"Expected %d godogs to be remaining, but there is %d", remaining, Godogs,
	)
}

// assertExpectedAndActual is a helper function to allow the step function to call
// assertion functions where you want to compare an expected and an actual value.
func assertExpectedAndActual(a expectedAndActualAssertion, expected, actual interface{}, msgAndArgs ...interface{}) error {
	var t asserter
	a(&t, expected, actual, msgAndArgs...)
	return t.err
}

type expectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool

// asserter is used to be able to retrieve the error reported by the called assertion
type asserter struct {
	err error
}

// Errorf is used by the called assertion to report an error
func (a *asserter) Errorf(format string, args ...interface{}) {
	a.err = fmt.Errorf(format, args...)
}
```

### Configure common options for godog CLI

There are no global options or configuration files. Alias your common or project based commands: `alias godog-wip="godog --format=progress --tags=@wip"`

### Testing browser interactions

**godog** does not come with builtin packages to connect to the browser. You may want to look at [selenium](http://www.seleniumhq.org/) and probably [phantomjs](http://phantomjs.org/). See also the following components:

1. [browsersteps](https://github.com/llonchj/browsersteps) - provides basic context steps to start selenium and navigate browser content.
2. You may wish to have [goquery](https://github.com/PuerkitoBio/goquery)
   in order to work with HTML responses like with JQuery.

### Concurrency

When concurrency is configured in options, godog will execute the scenarios concurrently, which is support by all supplied formatters.

In order to support concurrency well, you should reset the state and isolate each scenario. They should not share any state. It is suggested to run the suite concurrently in order to make sure there is no state corruption or race conditions in the application.

It is also useful to randomize the order of scenario execution, which you can now do with **--random** command option.

## License
**Godog** and **Gherkin** are licensed under the [MIT][license] and developed as a part of the [cucumber project][cucumber]

[godoc]: https://pkg.go.dev/github.com/cucumber/godog "Documentation on godog"
[golang]: https://golang.org/  "GO programming language"
[behat]: http://docs.behat.org/ "Behavior driven development framework for PHP"
[cucumber]: https://cucumber.io/ "Behavior driven development framework"
[license]: https://en.wikipedia.org/wiki/MIT_License "The MIT license"
