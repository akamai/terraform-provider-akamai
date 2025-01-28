// Package log contains types and functions related to logging in terraform provider
package log

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/log"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

// Logger is a wrapper over logger.Interface and hclog.Logger
type Logger struct {
	log.Interface
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
	args := hclog.ImpliedArgs()
	result := getFieldsFromHClog(args)

	rval := &Logger{
		Interface: log.Default().With("", result),
		hclog:     hclog,
	}

	lvlString := strings.ToLower(logging.LogLevel())

	var levelStrings = map[string]slog.Level{
		"trace": log.LevelTrace,
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	level := levelStrings[lvlString]

	handler := log.NewSlogHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	rval.Interface = log.NewSlogAdapter(handler).With("", result)

	return rval
}

func getFieldsFromHClog(args []interface{}) log.Fields {

	var extra interface{}

	if len(args)%2 != 0 {
		extra = args[len(args)-1]
		args = args[:len(args)-1]
	}
	result := make(map[string]interface{}, len(args))

	for i := 0; i < len(args); i += 2 {
		key := args[i].(string)
		result[key] = args[i+1]
	}

	if extra != nil {
		result["EXTRA_VALUE_AT_END"] = extra

	}
	return result
}

// Get returns a global log object, there is no context like operation id
func Get(args ...interface{}) *Logger {
	return FromHCLog(hclog.Default().With(args...))
}

// FromContext returns the logger from the context
func FromContext(ctx context.Context, args ...interface{}) *Logger {
	return FromHCLog(hclog.FromContext(ctx).With(args...))
}
