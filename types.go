package log

type (
	Field string
	Level int
)

type ExcludeRule struct {
	Field Field
	Value any
}

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
	LevelTrace
)

type Wrapper interface {
	GetLevel() Level
	WithField(key string, value interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}

type WrapperFactoryFunc func() Wrapper
