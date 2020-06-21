package xds

import xdslog "github.com/envoyproxy/go-control-plane/pkg/log"

// LoggerFuncs implements the go-control-plane Logger, allowing the
// caller to specify only the logging functions that are desired.
type LoggerFuncs struct {
	DebugFunc func(string, ...interface{})
	InfoFunc  func(string, ...interface{})
	WarnFunc  func(string, ...interface{})
	ErrorFunc func(string, ...interface{})
}

var _ xdslog.Logger = &LoggerFuncs{}

// Debugf logs a formatted debugging message.
func (f LoggerFuncs) Debugf(format string, args ...interface{}) {
	if f.DebugFunc != nil {
		f.DebugFunc(format, args...)
	}
}

// Infof logs a formatted informational message.
func (f LoggerFuncs) Infof(format string, args ...interface{}) {
	if f.InfoFunc != nil {
		f.InfoFunc(format, args...)
	}
}

// Warnf logs a formatted warning message.
func (f LoggerFuncs) Warnf(format string, args ...interface{}) {
	if f.WarnFunc != nil {
		f.WarnFunc(format, args...)
	}
}

// Errorf logs a formatted error message.
func (f LoggerFuncs) Errorf(format string, args ...interface{}) {
	if f.ErrorFunc != nil {
		f.ErrorFunc(format, args...)
	}
}
