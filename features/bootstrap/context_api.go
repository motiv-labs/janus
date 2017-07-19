package bootstrap

import (
	"github.com/DATA-DOG/godog"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/pkg/errors"
)

// RegisterAPIContext registers godog suite context for handling API related steps
func RegisterAPIContext(s *godog.Suite, readOnly bool, apiRepo api.Repository) {
	ctx := &apiContext{readOnly: readOnly, apiRepo: apiRepo}

	s.BeforeScenario(ctx.clearAPI)
}

type apiContext struct {
	readOnly bool
	apiRepo  api.Repository
}

func (c *apiContext) clearAPI(arg interface{}) {
	if !c.readOnly {
		data, err := c.apiRepo.FindAll()
		if err != nil {
			panic(errors.Wrap(err, "Failed to get all registered route specs"))
		}

		for _, definition := range data {
			err := c.apiRepo.Remove(definition.Name)
			if nil != err {
				panic(errors.Wrapf(err, "Failed to remove route spec \"%s\"", definition.Name))
			}
		}
	}
}
