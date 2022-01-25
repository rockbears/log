package log

import (
	"context"
	"fmt"
	"runtime"
	"sort"

	"github.com/pkg/errors"
)

func RegisterField(fields ...Field) {
	registeredFieldsMutex.Lock()
	defer registeredFieldsMutex.Unlock()

	for _, f := range fields {
		var exist bool
		for _, existingF := range registeredFields {
			if f == existingF {
				exist = true
				break
			}
		}
		if !exist {
			registeredFields = append(registeredFields, f)
		}
	}

	sort.Slice(registeredFields, func(i, j int) bool {
		return registeredFields[i] < registeredFields[j]
	})
}

func UnregisterField(fields ...Field) {
	registeredFieldsMutex.Lock()
	defer registeredFieldsMutex.Unlock()

loop:
	for _, f := range fields {
		for i, existingF := range registeredFields {
			if f == existingF {
				registeredFields = append(registeredFields[:i], registeredFields[i+1:]...)
				goto loop
			}
		}
	}

	sort.Slice(registeredFields, func(i, j int) bool {
		return registeredFields[i] < registeredFields[j]
	})
}

func Debug(ctx context.Context, format string, args ...interface{}) {
	call(ctx, LevelDebug, 2, format, args...)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	call(ctx, LevelInfo, 2, format, args...)
}

func Warn(ctx context.Context, format string, args ...interface{}) {
	call(ctx, LevelWarn, 2, format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	call(ctx, LevelError, 2, format, args...)
}

func Fatal(ctx context.Context, format string, args ...interface{}) {
	call(ctx, LevelFatal, 2, format, args...)
}

var (
	FieldSourceFile = Field("source_file")
	FieldSourceLine = Field("source_line")
	FieldCaller     = Field("caller")
	FieldStackTrace = Field("stack_trace")
)

func init() {
	RegisterField(FieldSourceFile, FieldSourceLine, FieldCaller, FieldStackTrace)
}

func call(ctx context.Context, level Level, callerOffset int, format string, args ...interface{}) {
	entry := Factory()

	if level < entry.GetLevel() {
		return
	}

	pc, file, line, ok := runtime.Caller(callerOffset)
	if ok {
		ctx = context.WithValue(ctx, FieldSourceFile, file)
		ctx = context.WithValue(ctx, FieldSourceLine, line)
		details := runtime.FuncForPC(pc)
		if details != nil {
			ctx = context.WithValue(ctx, FieldCaller, details.Name())
		}
	}

	for _, k := range registeredFields {
		v := ctx.Value(k)
		if v != nil {
			entry.WithField(string(k), v)
		}
	}

	switch level {
	case LevelInfo:
		entry.Infof(format, args...)

	case LevelWarn:
		entry.Warnf(format, args...)

	case LevelError:
		entry.Errorf(format, args...)

	case LevelFatal:
		entry.Fatalf(format, args...)

	case LevelPanic:
		entry.Panicf(format, args...)

	default:
		entry.Debugf(format, args...)
	}
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

func ContextWithStackTrace(ctx context.Context, err error) context.Context {
	errWithStracktrace, ok := err.(StackTracer)
	if ok {
		ctx = context.WithValue(ctx, FieldStackTrace, fmt.Sprintf("%+v", errWithStracktrace))
	}
	return ctx
}

func ErrorWithStackTrace(ctx context.Context, err error) {
	ctx = ContextWithStackTrace(ctx, err)
	call(ctx, LevelError, 2, err.Error())
}
