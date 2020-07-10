# TDD with Golang channels

Writing unit tests for code that uses channels has always been a stumbling block for me.

But today I finally did it. Please see [./internal/cats_test](https://github.com/johnmackenzie91/implementing-go-channels/blob/master/internal/cats/cats_test.go) directory for code and test.
The key take aways are;

### The time.After channel
Using time.After channel to stop unit test code from possibly slipping into a infinite loop and thus crashing your IDE

```go
	for {
		select {
		case <-time.After(2 * time.Second):
			t.Fatal("assumed that test is stuck in infinite loop")
		case v, ok := <-outCh:
			if !ok {
				closeLoop = true
				break
			}
			output = append(output, v)
		}
		if closeLoop {
			break
		}
	}
```

### Always use a context

Using a context will allow you to close go routines and associated channels gracefully.
Also make sure the first thing you do is check whether the context has been cancelled. This will prevent doing some work before checking whether it has expired.
Sounds obvious but in most cases I have seen the `work` been done, then a check to the context.Done()

```go
func (c CatsAPI) Fetch(ctx context.Context) (chan *Cat, chan error) {
	outCh := make(chan *Cat)
	errCh := make(chan error)
	go func() {
		for {
            // before we do anything, has the context been cancelled
			select {
			case <-ctx.Done():
				close(outCh)
				return
			default:
			}
    ...
        }
    ...
    }
...
}
```