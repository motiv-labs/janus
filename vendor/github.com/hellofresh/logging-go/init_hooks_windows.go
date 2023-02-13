package logging

func (c LogConfig) initHooks() error {
	for _, h := range c.Hooks {
		if err := c.validateRequiredHookSettings(h, []string{"host", "port"}); err != nil {
			return err
		}

		switch h.Format {
		case HookLogstash:
			if err := c.initLogstashHook(h); err != nil {
				return err
			}

		case HookGraylog:
			if err := c.initGraylogHook(h); err != nil {
				return err
			}

		default:
			return ErrUnknownLogHookFormat
		}
	}

	return nil
}
