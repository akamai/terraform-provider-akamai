package akamai

import (
	"context"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

type (
	logger struct {
		log.Logger
		l hclog.Logger
	}
)

// 2020/12/02 11:51:03
const (
	DefaultTimestampFormat = "2006/01/02 03:04:05"
)

func init() {
	if fmt, ok := os.LookupEnv("AKAMAI_TS_FORMAT"); ok {
		hclog.DefaultOptions.TimeFormat = fmt
	} else {
		hclog.DefaultOptions.TimeFormat = DefaultTimestampFormat
	}
}

// LogFromHCLog returns a new log.Interface from an hclog.Logger
func LogFromHCLog(l hclog.Logger) log.Interface {
	const (
		defaultLevel = log.InfoLevel
	)

	rval := &logger{
		Logger: log.Logger{
			Level: defaultLevel,
		},
		l: l,
	}

	// check for trace as the structured logger does not support trace
	// just make it debug to get everything from the provider
	lvlString := strings.ToLower(logging.LogLevel())
	if lvlString == "trace" {
		lvlString = "debug"
	}

	if lvl, err := log.ParseLevel(lvlString); err == nil {
		rval.Logger.Level = lvl
	}

	rval.Logger.Handler = rval

	return rval
}

// Log returns a global log object, there is no context like operation id
func Log(args ...interface{}) log.Interface {
	return LogFromHCLog(hclog.Default().With(args...))
}

// LogFromContext returns the logger from the context
func LogFromContext(ctx context.Context, args ...interface{}) log.Interface {
	return LogFromHCLog(hclog.FromContext(ctx).With(args...))
}

func (h *logger) HandleLog(e *log.Entry) error {
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
