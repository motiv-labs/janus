package bootstrap

import (
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/pkg/errors"

	"github.com/hellofresh/janus/pkg/api"
)

// RegisterAPIContext registers godog suite context for handling API related steps
func RegisterAPIContext(ctx *godog.ScenarioContext, apiRepo api.Repository, ch chan<- api.ConfigurationMessage) {
	scenarioCtx := &apiContext{apiRepo: apiRepo, ch: ch}

	ctx.BeforeScenario(scenarioCtx.clearAPI)
}

type apiContext struct {
	apiRepo api.Repository
	ch      chan<- api.ConfigurationMessage
}

func (c *apiContext) clearAPI(*messages.Pickle) {
	data, err := c.apiRepo.FindAll()
	if err != nil {
		panic(errors.Wrap(err, "Failed to get all registered route specs"))
	}

	for _, definition := range data {
		c.ch <- api.ConfigurationMessage{
			Operation:     api.RemovedOperation,
			Configuration: definition,
		}
	}

	time.Sleep(time.Second)
}
