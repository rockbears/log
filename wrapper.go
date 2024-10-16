package log

import (
	"fmt"
	"log"
	"os"
	"sort"
	"testing"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

/* Zap wrapper */

func NewZapWrapper(logger *zap.Logger) WrapperFactoryFunc {
	return func() Wrapper {
		return &ZapWrapper{logger: logger, sugar: logger.Sugar()}
	}
}

type ZapWrapper struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func (l *ZapWrapper) GetLevel() Level {
	switch {
	case l.logger.Core().Enabled(zap.DebugLevel):
		return LevelDebug
	case l.logger.Core().Enabled(zap.InfoLevel):
		return LevelInfo
	case l.logger.Core().Enabled(zap.WarnLevel):
		return LevelWarn
	case l.logger.Core().Enabled(zap.ErrorLevel):
		return LevelError
	case l.logger.Core().Enabled(zap.FatalLevel):
		return LevelFatal
	case l.logger.Core().Enabled(zap.PanicLevel):
		return LevelError
	default:
		panic(fmt.Errorf("zap level is not handled"))
	}
}
func (l *ZapWrapper) WithField(key string, value interface{}) {
	l.sugar = l.sugar.With(key, value)
}
func (l *ZapWrapper) format(format string, args ...interface{}) string {
	var msg = format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	return msg
}
func (l *ZapWrapper) Debugf(format string, args ...interface{}) {
	l.sugar.Debugw(l.format(format, args...))
}
func (l *ZapWrapper) Infof(format string, args ...interface{}) {
	l.sugar.Infow(l.format(format, args...))

}
func (l *ZapWrapper) Warnf(format string, args ...interface{}) {
	l.sugar.Warnw(l.format(format, args...))

}
func (l *ZapWrapper) Fatalf(format string, args ...interface{}) {
	l.sugar.Fatalw(l.format(format, args...))
}
func (l *ZapWrapper) Errorf(format string, args ...interface{}) {
	l.sugar.Errorw(l.format(format, args...))
}
func (l *ZapWrapper) Panicf(format string, args ...interface{}) {
	l.sugar.Panicw(l.format(format, args...))
}

/* Logrus wrapper */

func NewLogrusWrapper(logger *logrus.Logger) WrapperFactoryFunc {
	return func() Wrapper {
		return &LogrusWrapper{entry: logrus.NewEntry(logger)}
	}
}

type LogrusWrapper struct {
	entry *logrus.Entry
}

func (l *LogrusWrapper) GetLevel() Level {
	switch l.entry.Logger.Level {
	case logrus.DebugLevel:
		return LevelDebug
	case logrus.InfoLevel:
		return LevelInfo
	case logrus.WarnLevel:
		return LevelWarn
	case logrus.ErrorLevel:
		return LevelError
	case logrus.FatalLevel:
		return LevelFatal
	case logrus.PanicLevel:
		return LevelPanic
	case logrus.TraceLevel:
		return LevelTrace
	default:
		panic(fmt.Errorf("logrus level %q is not handled", l.entry.Level))
	}
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
func (l *LogrusWrapper) Panicf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.entry.Panic(format)
	} else {
		l.entry.Panicf(format, args...)
	}
}

/* testing.T wrapper */

func NewTestingWrapper(t testing.TB) WrapperFactoryFunc {
	return func() Wrapper {
		return &TestingWrapper{t: t}
	}
}

type TestingWrapper struct {
	ctx map[string]string
	t   testing.TB
}

func (l *TestingWrapper) GetLevel() Level {
	return LevelDebug
}
func (l *TestingWrapper) WithField(key string, value interface{}) {
	if l.ctx == nil {
		l.ctx = map[string]string{}
	}
	l.ctx[key] = fmt.Sprintf("%v", value)
}
func formatCtx(ctx map[string]string) string {
	var keys []string
	for k := range ctx {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var s string
	for _, k := range keys {
		v := ctx[k]
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
func (l *TestingWrapper) Panicf(format string, args ...interface{}) {
	l.log("PANIC", format, args...)
}

/* golang log package wrapper */

type StdWrapperOptions struct {
	Level            Level
	DisableTimestamp bool
}

type StdWrapper struct {
	opts StdWrapperOptions
	ctx  map[string]string
}

func NewStdWrapper(opts StdWrapperOptions) WrapperFactoryFunc {
	return func() Wrapper {
		return &StdWrapper{opts: opts}
	}
}

func (l *StdWrapper) GetLevel() Level {
	return l.opts.Level
}
func (l *StdWrapper) WithField(key string, value interface{}) {
	if l.ctx == nil {
		l.ctx = map[string]string{}
	}
	l.ctx[key] = fmt.Sprintf("%v", value)
}

func (l *StdWrapper) Print(s string) {
	if l.opts.DisableTimestamp {
		fmt.Println(s)
	} else {
		log.Println(s)
	}
}

func (l *StdWrapper) Debugf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.Print("[DEBUG] " + formatCtx(l.ctx) + " " + format)
	} else {
		l.Print("[DEBUG] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
}

func (l *StdWrapper) Infof(format string, args ...interface{}) {
	if len(args) == 0 {
		l.Print("[INFO] " + formatCtx(l.ctx) + " " + format)
	} else {
		l.Print("[INFO] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
}

func (l *StdWrapper) Warnf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.Print("[WARN] " + formatCtx(l.ctx) + " " + format)
	} else {
		l.Print("[WARN] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
}

func (l *StdWrapper) Fatalf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.Print("[FATAL] " + formatCtx(l.ctx) + " " + format)
	} else {
		l.Print("[FATAL] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
	os.Exit(1)
}

func (l *StdWrapper) Errorf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.Print("[ERROR] " + formatCtx(l.ctx) + " " + format)
	} else {
		l.Print("[ERROR] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
}

func (l *StdWrapper) Panicf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.Print("[PANIC] " + formatCtx(l.ctx) + " " + format)
	} else {
		l.Print("[PANIC] " + formatCtx(l.ctx) + " " + fmt.Sprintf(format, args...))
	}
}
