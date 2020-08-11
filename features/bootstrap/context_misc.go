package bootstrap

import (
	"time"

	"github.com/cucumber/godog"
)

const durationAWhile = time.Second

// RegisterMiscContext registers godog suite context for handling misc steps
func RegisterMiscContext(ctx *godog.ScenarioContext) {
	scenarioCtx := &miscContext{}

	ctx.Step(`^I wait for a while$`, scenarioCtx.iWaitForAWhile)
	ctx.Step(`^I wait for "([^"]*)"$`, scenarioCtx.iWaitFor)
}

type miscContext struct{}

func (c *miscContext) iWaitForAWhile() error {
	time.Sleep(durationAWhile)
	return nil
}

func (c *miscContext) iWaitFor(duration string) error {
	parsedDuration, err := time.ParseDuration(duration)
	if nil != err {
		return err
	}
	time.Sleep(parsedDuration)
	return nil
}
