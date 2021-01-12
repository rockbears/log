# rockbears/log

Log with context values as fields.
Compatible with `logrus`, `std` logger and `testing.T` logger.

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



