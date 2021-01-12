package log

import (
	"context"
)

func RegisterField(fields ...Field) {
	for _, f := range fields {
		registeredFields[f] = struct{}{}
	}
}

func Debug(ctx context.Context, format string, args ...interface{}) {
	call(ctx, debug, format, args...)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	call(ctx, info, format, args...)
}

func Warn(ctx context.Context, format string, args ...interface{}) {
	call(ctx, warn, format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	call(ctx, error, format, args...)
}

func Fatal(ctx context.Context, format string, args ...interface{}) {
	call(ctx, fatal, format, args...)
}

func call(ctx context.Context, level level, format string, args ...interface{}) {
	entry := Factory()
	for k := range registeredFields {
		v := ctx.Value(k)
		if v != nil {
			entry.WithField(string(k), v)
		}
	}
	switch level {
	case info:
		entry.Infof(format, args...)

	case warn:
		entry.Warnf(format, args...)

	case error:
		entry.Errorf(format, args...)

	case fatal:
		entry.Fatalf(format, args...)

	default:
		entry.Debugf(format, args...)
	}
}
