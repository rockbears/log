package log

import (
	"context"
	"runtime"
)

func RegisterField(fields ...Field) {
	registeredFieldsMutex.Lock()
	defer registeredFieldsMutex.Unlock()

	for _, f := range fields {
		registeredFields[f] = struct{}{}
	}
}

func UnregisterField(fields ...Field) {
	registeredFieldsMutex.Lock()
	defer registeredFieldsMutex.Unlock()

	for _, f := range fields {
		delete(registeredFields, f)
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

var (
	FieldSourceFile = Field("source_file")
	FieldSourceLine = Field("source_line")
	FieldCaller     = Field("caller")
)

func init() {
	RegisterField(FieldSourceFile, FieldSourceLine, FieldCaller)
}

func call(ctx context.Context, level level, format string, args ...interface{}) {
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		ctx = context.WithValue(ctx, FieldSourceFile, file)
		ctx = context.WithValue(ctx, FieldSourceLine, line)
		details := runtime.FuncForPC(pc)
		if details != nil {
			ctx = context.WithValue(ctx, FieldCaller, details.Name())
		}
	}

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
