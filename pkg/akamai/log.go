package akamai

import (
	"github.com/apex/log"
	"github.com/hashicorp/go-hclog"
)

type (
	logHandler struct {
		l hclog.Logger
	}
)

func (h *logHandler) HandleLog(e *log.Entry) error {
	fields := make([]interface{}, 0)

	for k, v := range e.Fields {
		fields = append(fields, k, v)
	}

	switch e.Level {
	case log.DebugLevel:
		h.l.Debug(e.Message, fields...)
	case log.InfoLevel:
		h.l.Info(e.Message, fields...)
	case log.WarnLevel:
		h.l.Warn(e.Message, fields...)
	case log.ErrorLevel:
		h.l.Error(e.Message, fields...)
	case log.FatalLevel:
		panic(e.Message)
	}

	return nil
}
