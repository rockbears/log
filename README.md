# rockbears/log

Log with context values as fields.

Compatible with [zap](https://github.com/uber-go/zap), [logrus](https://github.com/sirupsen/logrus), [std](https://pkg.go.dev/log) logger and [testing.T](https://pkg.go.dev/testing#T) logger.

It supports [pkg/errors](https://github.com/pkg/errors) to add a `stack_trace` field if the handled error `error` implements `StackTracer`interface.

It offers a convenient way to keep your logs when you are running unit tests.

## Install

```golang
    go get github.com/rockbears/log
```

## How to use

By default, it is initialized to wrap `logrus` package. You can override `log.Factory` to use the logger library you want.

First register `fields`

```golang
    const myField = log.Field("component")

    func init() {
        log.RegisterField(myField)
    }
```

Then add `fields` as values to you current context.

```golang
    ctx = context.WithValue(ctx, myField, "myComponent")
```

Finally log as usual.
```golang
    log.Info(ctx, "this is a log")
```

## Examples

```golang
    const myField = log.Field("component")

    func init() {
        log.RegisterField(myField)
    }

    func foo(ctx context.Context) {
        ctx = context.WithValue(ctx, myField, "myComponent")
        log.Info(ctx, "this is a log")
    }
```

Preserve your log in unit tests.

```golang
    func TestFoo(t *testing.T) {
        log.Factory = log.NewTestingWrapper(t)
        foo(context.TODO())
    }
```

Log errors easily.

```golang
    import "github.com/pkg/errors"

    func foo() {
        ctx := context.Background()
        err := errors.New("this is an error") // from package "github.com/pkg/errors"
        log.ErrorWithStackTrace(ctx, err) // will produce a nice stack_trace field 
    )
```


