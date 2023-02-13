// +build !windows

package logging

import (
	"fmt"
	"log/syslog"

	log "github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

var (
	severitiesMap = map[string]syslog.Priority{
		"LOG_EMERG":   syslog.LOG_EMERG,
		"LOG_ALERT":   syslog.LOG_ALERT,
		"LOG_CRIT":    syslog.LOG_CRIT,
		"LOG_ERR":     syslog.LOG_ERR,
		"LOG_WARNING": syslog.LOG_WARNING,
		"LOG_NOTICE":  syslog.LOG_NOTICE,
		"LOG_INFO":    syslog.LOG_INFO,
		"LOG_DEBUG":   syslog.LOG_DEBUG,
	}

	facilitiesMap = map[string]syslog.Priority{
		"LOG_KERN":     syslog.LOG_KERN,
		"LOG_USER":     syslog.LOG_USER,
		"LOG_MAIL":     syslog.LOG_MAIL,
		"LOG_DAEMON":   syslog.LOG_DAEMON,
		"LOG_AUTH":     syslog.LOG_AUTH,
		"LOG_SYSLOG":   syslog.LOG_SYSLOG,
		"LOG_LPR":      syslog.LOG_LPR,
		"LOG_NEWS":     syslog.LOG_NEWS,
		"LOG_UUCP":     syslog.LOG_UUCP,
		"LOG_CRON":     syslog.LOG_CRON,
		"LOG_AUTHPRIV": syslog.LOG_AUTHPRIV,
		"LOG_FTP":      syslog.LOG_FTP,
		"LOG_LOCAL0":   syslog.LOG_LOCAL0,
		"LOG_LOCAL1":   syslog.LOG_LOCAL1,
		"LOG_LOCAL2":   syslog.LOG_LOCAL2,
		"LOG_LOCAL3":   syslog.LOG_LOCAL3,
		"LOG_LOCAL4":   syslog.LOG_LOCAL4,
		"LOG_LOCAL5":   syslog.LOG_LOCAL5,
		"LOG_LOCAL6":   syslog.LOG_LOCAL6,
		"LOG_LOCAL7":   syslog.LOG_LOCAL7,
	}
)

func (c LogConfig) initSyslogHook(h LogHook) error {
	if err := c.validateRequiredHookSettings(h, []string{"network"}); err != nil {
		return err
	}
	network, _ := h.Settings["network"]

	priority, err := getSyslogPriority(h.Settings)
	if nil != err {
		log.WithError(err).WithField("hook", h.Format).Error("Failed to configure hook")
		return ErrFailedToConfigureLogHook
	}

	tag, _ := h.Settings["tag"]
	hook, err := logrus_syslog.NewSyslogHook(network, fmt.Sprintf("%s:%s", h.Settings["host"], h.Settings["port"]), priority, tag)
	if nil != err {
		log.WithError(err).WithField("hook", h.Format).Error("Failed to configure hook")
		return ErrFailedToConfigureLogHook
	}

	log.AddHook(hook)

	return nil
}

func getSyslogPriority(settings map[string]string) (syslog.Priority, error) {
	severity, ok := settings["severity"]
	if !ok {
		return 0, fmt.Errorf("Syslog severity setting is not set")
	}

	facility, ok := settings["facility"]
	if !ok {
		return 0, fmt.Errorf("Syslog facility setting is not set")
	}

	severityPriority, ok := severitiesMap[severity]
	if !ok {
		return 0, fmt.Errorf("Unknown syslog severity value")
	}

	facilityPriority, ok := facilitiesMap[facility]
	if !ok {
		return 0, fmt.Errorf("Unknown syslog facility value")
	}

	return severityPriority | facilityPriority, nil
}
