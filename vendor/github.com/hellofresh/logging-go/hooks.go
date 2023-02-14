package logging

import (
	"fmt"
	"net"
	"strconv"

	"github.com/bshuster-repo/logrus-logstash-hook"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gemnasium/logrus-graylog-hook.v2"
)

func (c LogConfig) initLogstashHook(h LogHook) error {
	if err := c.validateRequiredHookSettings(h, []string{"network"}); err != nil {
		return err
	}
	network, _ := h.Settings["network"]

	conn, err := net.Dial(network, fmt.Sprintf("%s:%s", h.Settings["host"], h.Settings["port"]))
	if nil != err {
		log.WithError(err).WithField("hook", h.Format).Error("Failed to connect to logstash")
		return ErrFailedToConfigureLogHook
	}

	formatter := getLogstashFormatter(h.Settings).(*logrustash.LogstashFormatter)

	hook, err := logrustash.NewHookWithConn(conn, formatter.Type)
	if nil != err {
		log.WithError(err).WithField("hook", h.Format).Error("Failed to instantiate logstash hook")
		return ErrFailedToConfigureLogHook
	}
	hook.TimeFormat = formatter.TimestampFormat
	log.AddHook(hook)

	return nil
}

func (c LogConfig) initGraylogHook(h LogHook) error {
	var async bool
	var err error
	if asyncStr, ok := h.Settings["async"]; ok {
		async, err = strconv.ParseBool(asyncStr)
		if nil != err {
			log.WithError(err).WithField("hook", h.Format).Error("Failed to parse async setting")
			return ErrFailedToConfigureLogHook
		}
	}

	extra := make(map[string]interface{})
	for k, v := range h.Settings {
		if k != "host" && k != "port" && k != "async" {
			extra[k] = v
		}
	}
	var hook *graylog.GraylogHook
	if async {
		hook = graylog.NewAsyncGraylogHook(fmt.Sprintf("%s:%s", h.Settings["host"], h.Settings["port"]), extra)
		c.mustFlushHooks = append(c.mustFlushHooks, hook)
	} else {
		hook = graylog.NewGraylogHook(fmt.Sprintf("%s:%s", h.Settings["host"], h.Settings["port"]), extra)
	}

	log.AddHook(hook)

	return nil
}
