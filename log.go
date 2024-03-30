package log

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	FieldSourceFile = Field("source_file")
	FieldSourceLine = Field("source_line")
	FieldCaller     = Field("caller")
	FieldStackTrace = Field("stack_trace")
)

var global *Logger
var Factory WrapperFactoryFunc = NewLogrusWrapper(logrus.StandardLogger())

func init() {
	global = New()
	global.callerFrameToSkip = 3
}

type Logger struct {
	registeredFields      []Field
	registeredFieldsMutex sync.RWMutex
	excludeRules          []ExcludeRule
	excludeRulesMutex     sync.RWMutex
	factory               WrapperFactoryFunc
	callerFrameToSkip     int
}

func New() *Logger {
	return NewWithFactory(nil)
}

func NewWithFactory(factory WrapperFactoryFunc) *Logger {
	logger := &Logger{factory: factory, callerFrameToSkip: 2}
	logger.RegisterDefaultFields()
	return logger
}

func (l *Logger) GetFramesToSkip() int {
	return l.callerFrameToSkip
}

func (l *Logger) SetFramesToSkip(s int) {
	l.callerFrameToSkip = s
}

func (l *Logger) RegisterField(fields ...Field) {
	l.registeredFieldsMutex.Lock()
	defer l.registeredFieldsMutex.Unlock()

	for _, f := range fields {
		var exist bool
		for _, existingF := range l.registeredFields {
			if f == existingF {
				exist = true
				break
			}
		}
		if !exist {
			l.registeredFields = append(l.registeredFields, f)
		}
	}

	sort.Slice(l.registeredFields, func(i, j int) bool {
		return l.registeredFields[i] < l.registeredFields[j]
	})
}

func (l *Logger) UnregisterField(fields ...Field) {
	l.registeredFieldsMutex.Lock()
	defer l.registeredFieldsMutex.Unlock()

loop:
	for _, f := range fields {
		for i, existingF := range l.registeredFields {
			if f == existingF {
				l.registeredFields = append(l.registeredFields[:i], l.registeredFields[i+1:]...)
				goto loop
			}
		}
	}

	sort.Slice(l.registeredFields, func(i, j int) bool {
		return l.registeredFields[i] < l.registeredFields[j]
	})
}

func (l *Logger) GetRegisteredFields() []Field {
	l.registeredFieldsMutex.RLock()
	defer l.registeredFieldsMutex.RUnlock()
	fields := make([]Field, len(l.registeredFields))
	copy(fields, l.registeredFields)
	return fields
}

func (l *Logger) GetExcludeRules() []ExcludeRule {
	l.excludeRulesMutex.RLock()
	defer l.excludeRulesMutex.RUnlock()
	excludeRules := make([]ExcludeRule, len(l.excludeRules))
	copy(excludeRules, l.excludeRules)
	return excludeRules
}

func (l *Logger) RegisterDefaultFields() {
	l.RegisterField(FieldSourceFile, FieldSourceLine, FieldCaller, FieldStackTrace)
}

func (l *Logger) Skip(field Field, value interface{}) {
	l.excludeRulesMutex.Lock()
	defer l.excludeRulesMutex.Unlock()
	for i := range l.excludeRules {
		if l.excludeRules[i].Field == field {
			l.excludeRules[i].Value = value
			return
		}
	}
	l.excludeRules = append(l.excludeRules, ExcludeRule{field, value})
}

func (l *Logger) Debug(ctx context.Context, format string, args ...interface{}) {
	l.call(ctx, LevelDebug, format, args...)
}

func (l *Logger) Info(ctx context.Context, format string, args ...interface{}) {
	l.call(ctx, LevelInfo, format, args...)
}

func (l *Logger) Warn(ctx context.Context, format string, args ...interface{}) {
	l.call(ctx, LevelWarn, format, args...)
}

func (l *Logger) Error(ctx context.Context, format string, args ...interface{}) {
	l.call(ctx, LevelError, format, args...)
}

func (l *Logger) Fatal(ctx context.Context, format string, args ...interface{}) {
	l.call(ctx, LevelFatal, format, args...)
}

func (l *Logger) Panic(ctx context.Context, format string, args ...interface{}) {
	l.call(ctx, LevelPanic, format, args...)
}

func (l *Logger) call(ctx context.Context, level Level, format string, args ...interface{}) {
	var entry Wrapper
	if l.factory == nil {
		entry = Factory()
	} else {
		entry = l.factory()
	}

	if level < entry.GetLevel() {
		return
	}

	pc, file, line, ok := runtime.Caller(l.callerFrameToSkip)
	if ok {
		ctx = context.WithValue(ctx, FieldSourceFile, file)
		ctx = context.WithValue(ctx, FieldSourceLine, line)
		details := runtime.FuncForPC(pc)
		if details != nil {
			ctx = context.WithValue(ctx, FieldCaller, details.Name())
		}
	}

	mExcludeRules := make(map[Field]any)
	for _, rule := range l.GetExcludeRules() {
		mExcludeRules[rule.Field] = rule.Value
	}

	for _, k := range l.GetRegisteredFields() {
		v := ctx.Value(k)
		if v != nil {
			if excludeValue, has := mExcludeRules[k]; has {
				if v == excludeValue {
					return
				}
			}
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

func (l *Logger) ErrorWithStackTrace(ctx context.Context, err error) {
	ctx = ContextWithStackTrace(ctx, err)
	l.call(ctx, LevelError, err.Error())
}

func (l *Logger) FieldValues(ctx context.Context) map[Field]interface{} {
	res := make(map[Field]interface{}, 10)
	for _, k := range l.registeredFields {
		v := ctx.Value(k)
		if v != nil {
			res[k] = v
		}
	}
	return res
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

func GetFramesToSkip() int {
	return global.GetFramesToSkip()
}

func SetFramesToSkip(s int) {
	global.SetFramesToSkip(s)
}

func RegisterField(fields ...Field) {
	global.RegisterField(fields...)
}

func UnregisterField(fields ...Field) {
	global.UnregisterField(fields...)
}

func GetRegisteredFields() []Field {
	return global.GetRegisteredFields()
}

func RegisterDefaultFields() {
	global.RegisterDefaultFields()
}

func Skip(field Field, value interface{}) {
	global.Skip(field, value)
}

func Debug(ctx context.Context, format string, args ...interface{}) {
	global.Debug(ctx, format, args...)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	global.Info(ctx, format, args...)
}

func Warn(ctx context.Context, format string, args ...interface{}) {
	global.Warn(ctx, format, args...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	global.Error(ctx, format, args...)
}

func Fatal(ctx context.Context, format string, args ...interface{}) {
	global.Fatal(ctx, format, args...)
}

func Panic(ctx context.Context, format string, args ...interface{}) {
	global.Panic(ctx, format, args...)
}

func ErrorWithStackTrace(ctx context.Context, err error) {
	global.ErrorWithStackTrace(ctx, err)
}

func FieldValues(ctx context.Context) map[Field]interface{} {
	return global.FieldValues(ctx)
}
