package log

import (
	"fmt"
	l "log"
	"os"
	"strings"
)

// Handler defines logger type to override stats debug and error logging
type Handler func(msg string, fields map[string]interface{}, err error)

var handler Handler

func init() {
	logger := l.New(os.Stderr, "[StatsGo] ", l.LstdFlags)

	SetHandler(func(msg string, fields map[string]interface{}, err error) {
		if nil != err {
			if nil == fields {
				fields = make(map[string]interface{})
			}
			fields["error"] = err.Error()
		}

		msgParts := make([]string, len(fields)+1)

		msgParts[0] = msg
		idx := 1
		for k, v := range fields {
			msgParts[idx] = fmt.Sprintf("%s=%v", k, v)
			idx++
		}

		logger.Println(strings.Join(msgParts, "\t"))
	})
}

// SetHandler sets log handler to use for stats debug and error logging
func SetHandler(h Handler) {
	handler = h
}

// Log calls log handler
func Log(msg string, fields map[string]interface{}, err error) {
	handler(msg, fields, err)
}
