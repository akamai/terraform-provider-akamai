// Package logger contains types and functions related to logging in terraform provider
package logger

import (
	"context"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

// Logger is a wrapper over log.Logger and hclog.Logger
type Logger struct {
	log.Logger
	hclog hclog.Logger
}

const defaultTimestampFormat = "2006/01/02 03:04:05"

func init() {
	if fmt, ok := os.LookupEnv("AKAMAI_TS_FORMAT"); ok {
		hclog.DefaultOptions.TimeFormat = fmt
	} else {
		hclog.DefaultOptions.TimeFormat = defaultTimestampFormat
	}
}

// HCLog returns Logger's hclog
func (l *Logger) HCLog() hclog.Logger {
	return l.hclog
}

// FromHCLog returns a new Logger from a hclog.Logger
func FromHCLog(hclog hclog.Logger) *Logger {
	const (
		defaultLevel = log.InfoLevel
	)

	rval := &Logger{
		Logger: log.Logger{
			Level: defaultLevel,
		},
		hclog: hclog,
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

// Get returns a global log object, there is no context like operation id
func Get(args ...interface{}) *Logger {
	return FromHCLog(hclog.Default().With(args...))
}

// FromContext returns the logger from the context
func FromContext(ctx context.Context, args ...interface{}) *Logger {
	return FromHCLog(hclog.FromContext(ctx).With(args...))
}

// HandleLog implements the logic for handling log events
func (l *Logger) HandleLog(e *log.Entry) error {
	fields := make([]interface{}, 0)

	for k, v := range e.Fields {
		fields = append(fields, k, v)
	}

	switch e.Level {
	case log.DebugLevel:
		l.hclog.Debug(e.Message, fields...)
	case log.InfoLevel:
		l.hclog.Info(e.Message, fields...)
	case log.WarnLevel:
		l.hclog.Warn(e.Message, fields...)
	case log.ErrorLevel:
		l.hclog.Error(e.Message, fields...)
	case log.FatalLevel:
		panic(e.Message)
	}

	return nil
}
