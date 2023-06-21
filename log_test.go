package log_test

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rockbears/log"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	fieldComponent = log.Field("component")
	fieldAsset     = log.Field("asset")
)

func init() {
	log.RegisterField(fieldComponent, fieldAsset)
}

func ExampleNewLogrusWrapper() {
	// Init the wrapper
	log.Factory = log.NewLogrusWrapper
	log.UnregisterField(log.FieldSourceLine, log.FieldSourceFile)

	// Init the logrus logger
	logrus.StandardLogger().SetLevel(logrus.InfoLevel)
	logrus.StandardLogger().SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.StandardLogger().Out = os.Stdout

	// Init the context
	ctx := context.Background()
	ctx = context.WithValue(ctx, fieldComponent, "rockbears/log")
	ctx = context.WithValue(ctx, fieldAsset, "ExampleWithLogrus")
	log.Debug(ctx, "this log should not be displayed")
	log.Info(ctx, "this is %q", "info")
	log.Warn(ctx, "this is warn")
	log.Error(ctx, "this is error")

	// Output:
	// level=info msg="this is \"info\"" asset=ExampleWithLogrus caller=github.com/rockbears/log_test.ExampleNewLogrusWrapper component=rockbears/log
	// level=warning msg="this is warn" asset=ExampleWithLogrus caller=github.com/rockbears/log_test.ExampleNewLogrusWrapper component=rockbears/log
	// level=error msg="this is error" asset=ExampleWithLogrus caller=github.com/rockbears/log_test.ExampleNewLogrusWrapper component=rockbears/log
}

func ExampleNewStdWrapper() {
	// Init the wrapper
	log.Factory = log.NewStdWrapper(log.StdWrapperOptions{Level: log.LevelInfo, DisableTimestamp: true})
	log.UnregisterField(log.FieldSourceLine, log.FieldSourceFile)
	// Init the context
	ctx := context.Background()
	ctx = context.WithValue(ctx, fieldComponent, "rockbears/log")
	ctx = context.WithValue(ctx, fieldAsset, "ExampleNewStdWrapper")
	log.Debug(ctx, "this log should not be displayed")
	log.Info(ctx, "this is %q", "info")
	log.Warn(ctx, "this is warn")
	log.Error(ctx, "this is error")
	// Output:
	// [INFO] [asset=ExampleNewStdWrapper][caller=github.com/rockbears/log_test.ExampleNewStdWrapper][component=rockbears/log] this is "info"
	// [WARN] [asset=ExampleNewStdWrapper][caller=github.com/rockbears/log_test.ExampleNewStdWrapper][component=rockbears/log] this is warn
	// [ERROR] [asset=ExampleNewStdWrapper][caller=github.com/rockbears/log_test.ExampleNewStdWrapper][component=rockbears/log] this is error
}

func ExampleNewZapWrapper() {
	// Init the wrapper
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, zap.InfoLevel)
	log.Factory = log.NewZapWrapper(zap.New(core))
	log.UnregisterField(log.FieldSourceLine, log.FieldSourceFile)
	// Init the context
	ctx := context.Background()
	ctx = context.WithValue(ctx, fieldComponent, "rockbears/log")
	ctx = context.WithValue(ctx, fieldAsset, "ExampleNewZapWrapper")
	log.Debug(ctx, "this log should not be displayed")
	log.Info(ctx, "this is %q", "info")
	log.Warn(ctx, "this is warn")
	log.Error(ctx, "this is error")
	// Output:
	// {"level":"info","msg":"this is \"info\"","asset":"ExampleNewZapWrapper","caller":"github.com/rockbears/log_test.ExampleNewZapWrapper","component":"rockbears/log"}
	// {"level":"warn","msg":"this is warn","asset":"ExampleNewZapWrapper","caller":"github.com/rockbears/log_test.ExampleNewZapWrapper","component":"rockbears/log"}
	// {"level":"error","msg":"this is error","asset":"ExampleNewZapWrapper","caller":"github.com/rockbears/log_test.ExampleNewZapWrapper","component":"rockbears/log"}
}

func TestErrorWithStackTrace() {
	// Init the wrapper
	log.Factory = log.NewStdWrapper(log.StdWrapperOptions{Level: log.LevelInfo, DisableTimestamp: true})
	log.UnregisterField(log.FieldSourceLine, log.FieldSourceFile)
	// Init the context
	ctx := context.Background()
	ctx = context.WithValue(ctx, fieldComponent, "rockbears/log")
	ctx = context.WithValue(ctx, fieldAsset, "ExampleErrorWithStackTrace")
	log.ErrorWithStackTrace(ctx, fmt.Errorf("this is an error"))
	log.ErrorWithStackTrace(ctx, errors.WithStack(fmt.Errorf("this is an error")))
}

func ExampleNewStdWrapperAndSkip() {
	// Init the wrapper
	log.Factory = log.NewStdWrapper(log.StdWrapperOptions{Level: log.LevelInfo, DisableTimestamp: true})
	log.UnregisterField(log.FieldSourceLine, log.FieldSourceFile)
	// Init the context
	ctx := context.Background()
	ctx = context.WithValue(ctx, fieldComponent, "rockbears/log")
	ctx = context.WithValue(ctx, fieldAsset, "ExampleNewStdWrapper")
	log.Debug(ctx, "this log should not be displayed")
	log.Info(ctx, "this is %q", "info")
	log.Warn(ctx, "this is warn")
	log.Error(ctx, "this is error")
	log.Skip(fieldAsset, "ExampleNewStdWrapper")
	log.Info(ctx, "this log should not be displayed because is should be skipped")
	// Output:
	// [INFO] [asset=ExampleNewStdWrapper][caller=github.com/rockbears/log_test.ExampleNewStdWrapperAndSkip][component=rockbears/log] this is "info"
	// [WARN] [asset=ExampleNewStdWrapper][caller=github.com/rockbears/log_test.ExampleNewStdWrapperAndSkip][component=rockbears/log] this is warn
	// [ERROR] [asset=ExampleNewStdWrapper][caller=github.com/rockbears/log_test.ExampleNewStdWrapperAndSkip][component=rockbears/log] this is error
}
