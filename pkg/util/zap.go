package util

import "go.uber.org/zap/zapcore"

// ZapEnablerFunc implements zapcore.LevelEnabler.
type ZapEnablerFunc func(zapcore.Level) bool

// Enabled ...
func (z ZapEnablerFunc) Enabled(level zapcore.Level) bool {
	if z != nil {
		return z(level)
	}

	return false
}

var _ zapcore.LevelEnabler = ZapEnablerFunc(nil)

// ZapEnableDebug enables all zap log levels.
func ZapEnableDebug() zapcore.LevelEnabler {
	return ZapEnablerFunc(func(zapcore.Level) bool { return true })
}
