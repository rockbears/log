package log

import "sync"

type (
	Field string
	Level int
)

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
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

var Factory WrapperFactoryFunc = NewLogrusWrapper

var registeredFields []Field

var registeredFieldsMutex sync.Mutex
