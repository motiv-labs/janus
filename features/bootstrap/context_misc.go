package bootstrap

import (
	"time"

	"github.com/DATA-DOG/godog"
)

const durationAWhile = time.Second

// RegisterMiscContext registers godog suite context for handling misc steps
func RegisterMiscContext(s *godog.Suite) {
	ctx := &miscContext{}

	s.Step(`^I wait for a while$`, ctx.iWaitForAWhile)
	s.Step(`^I wait for "([^"]*)"$`, ctx.iWaitFor)
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
