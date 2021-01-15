package log

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

/* Logrus wrapper */

func NewLogrusWrapper() Wrapper {
	return &LogrusWrapper{entry: logrus.NewEntry(logrus.StandardLogger())}
}

type LogrusWrapper struct {
	entry *logrus.Entry
}

func (l *LogrusWrapper) WithField(key string, value interface{}) {
	l.entry = l.entry.WithField(key, value)
}
func (l *LogrusWrapper) Debugf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.entry.Debug(format)
	} else {
		l.entry.Debugf(format, args...)
	}
}
func (l *LogrusWrapper) Infof(format string, args ...interface{}) {
	if len(args) == 0 {
		l.entry.Info(format)
	} else {
		l.entry.Infof(format, args...)
	}
}
func (l *LogrusWrapper) Warnf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.entry.Warn(format)
	} else {
		l.entry.Warnf(format, args...)
	}
}
func (l *LogrusWrapper) Fatalf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.entry.Fatal(format)
	} else {
		l.entry.Fatalf(format, args...)
	}
}
func (l *LogrusWrapper) Errorf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.entry.Error(format)
	} else {
		l.entry.Errorf(format, args...)
	}
}

func NewTestingWrapper(t *testing.T) WrapperFactoryFunc {
	return func() Wrapper {
		return &TestingWrapper{t: t}
	}
}

/* testing.T wrapper */

type TestingWrapper struct {
	ctx map[string]string
	t   *testing.T
}

func (l *TestingWrapper) WithField(key string, value interface{}) {
	if l.ctx == nil {
		l.ctx = map[string]string{}
	}
	l.ctx[key] = fmt.Sprintf("%v", value)
}
func formatCtx(ctx map[string]string) string {
	var s string
	for k, v := range ctx {
		s += fmt.Sprintf("[%s=%s]", k, v)
	}
	return s
}
func (l *TestingWrapper) log(level, format string, args ...interface{}) {
	// Recover function to avoid panic: Log in goroutine after test has completed
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[" + level + "] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
		}
	}()
	if len(args) == 0 {
		l.t.Log("[" + level + "] " + formatCtx(l.ctx) + " " + format)
	} else {
		l.t.Logf("[" + level + "] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
}
func (l *TestingWrapper) fatal(format string, args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[FATAL] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
			os.Exit(2)
		}
	}()
	if len(args) == 0 {
		l.t.Fatal("[FATAL] " + formatCtx(l.ctx) + " " + format)
	} else {
		l.t.Fatalf("[FATAL] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
}

func (l *TestingWrapper) Debugf(format string, args ...interface{}) {
	l.log("DEBUG", format, args...)
}
func (l *TestingWrapper) Infof(format string, args ...interface{}) {
	l.log("INFO", format, args...)

}
func (l *TestingWrapper) Warnf(format string, args ...interface{}) {
	l.log("WARN", format, args...)

}
func (l *TestingWrapper) Fatalf(format string, args ...interface{}) {
	l.fatal(format, args...)
}
func (l *TestingWrapper) Errorf(format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}

/* golang log package wrapper */

func NewStdWrapper() Wrapper {
	return &StdWrapper{}
}

type StdWrapper struct {
	ctx map[string]string
}

func (l *StdWrapper) WithField(key string, value interface{}) {
	if l.ctx == nil {
		l.ctx = map[string]string{}
	}
	l.ctx[key] = fmt.Sprintf("%v", value)
}

func (l *StdWrapper) Debugf(format string, args ...interface{}) {
	if len(args) == 0 {
		log.Print("[DEBUG] " + formatCtx(l.ctx) + " " + format)
	} else {
		log.Printf("[DEBUG] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
}

func (l *StdWrapper) Infof(format string, args ...interface{}) {
	if len(args) == 0 {
		log.Print("[INFO] " + formatCtx(l.ctx) + " " + format)
	} else {
		log.Printf("[INFO] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
}

func (l *StdWrapper) Warnf(format string, args ...interface{}) {
	if len(args) == 0 {
		log.Print("[WARN] " + formatCtx(l.ctx) + " " + format)
	} else {
		log.Printf("[WARN] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
}

func (l *StdWrapper) Fatalf(format string, args ...interface{}) {
	if len(args) == 0 {
		log.Fatal("[FATAL] " + formatCtx(l.ctx) + " " + format)
	} else {
		log.Fatalf("[FATAL] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
}

func (l *StdWrapper) Errorf(format string, args ...interface{}) {
	if len(args) == 0 {
		log.Print("[ERROR] " + formatCtx(l.ctx) + " " + format)
	} else {
		log.Printf("[ERROR] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
}
