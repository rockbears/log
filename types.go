package log

import "sync"

type (
	Field string
	level int
)

const (
	debug level = iota
	info
	warn
	error
	fatal
)

type Wrapper interface {
	WithField(key string, value interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type WrapperFactoryFunc func() Wrapper

var Factory WrapperFactoryFunc = NewLogrusWrapper

var registeredFields = map[Field]struct{}{}

var registeredFieldsMutex sync.Mutex
