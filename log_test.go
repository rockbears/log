package log_test

import (
	"context"
	"fmt"
	"os"
	"testing"

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
	registerDefaultFields()
}

func registerDefaultFields() {
	log.RegisterField(fieldComponent, fieldAsset)
}

func ExampleNewLogrusWrapper() {
	// Init the wrapper
	log.Factory = log.NewLogrusWrapper(logrus.StandardLogger())
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

func TestErrorWithStackTrace(t *testing.T) {
	// Init the wrapper
	log.Factory = log.NewTestingWrapper(t)
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

func ExampleGetRegisteredFields() {
	log.RegisterDefaultFields()
	registerDefaultFields()
	fmt.Println(log.GetRegisteredFields())
	// Output:
	// [asset caller component source_file source_line stack_trace]
}

func TestGetRegisteredFields(t *testing.T) {
	log.RegisterDefaultFields()
	defer log.RegisterDefaultFields()
	registerDefaultFields()
	defer registerDefaultFields()

	log.UnregisterField(log.GetRegisteredFields()...)

	got := log.GetRegisteredFields()
	if len(got) != 0 {
		t.Fatalf("want empty slice, got %s", got)
	}

	log.RegisterField(log.FieldCaller)
	got = log.GetRegisteredFields()
	if len(got) != 1 {
		t.Fatalf("want slice with length 1, got %s", got)
	}

	gotField := got[0]
	wantField := log.FieldCaller
	if gotField != wantField {
		t.Fatalf("want field %s, got %s", wantField, gotField)
	}
}

func ExampleNewWithFactory() {
	// Init the wrapper
	lrus := logrus.New()
	logger := log.NewWithFactory(log.NewLogrusWrapper(lrus))
	logger.RegisterField(fieldComponent, fieldAsset)
	logger.UnregisterField(log.FieldSourceLine, log.FieldSourceFile)

	// Init the logrus logger
	lrus.SetLevel(logrus.InfoLevel)
	lrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	lrus.Out = os.Stdout

	// Init the context
	ctx := context.Background()
	ctx = context.WithValue(ctx, fieldComponent, "rockbears/log")
	ctx = context.WithValue(ctx, fieldAsset, "ExampleWithLogrus")
	logger.Debug(ctx, "this log should not be displayed")
	logger.Info(ctx, "this is %q", "info")
	logger.Warn(ctx, "this is warn")
	logger.Error(ctx, "this is error")

	// Output:
	// level=info msg="this is \"info\"" asset=ExampleWithLogrus caller=github.com/rockbears/log_test.ExampleNewWithFactory component=rockbears/log
	// level=warning msg="this is warn" asset=ExampleWithLogrus caller=github.com/rockbears/log_test.ExampleNewWithFactory component=rockbears/log
	// level=error msg="this is error" asset=ExampleWithLogrus caller=github.com/rockbears/log_test.ExampleNewWithFactory component=rockbears/log
}

func TestGlobalDefaultLog(t *testing.T) {
	log.Info(context.Background(), "hello world")
}
