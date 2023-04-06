## 🌸 rchan: Go channel thorugh Redis List

[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/rchan)](https://goreportcard.com/report/github.com/nikolaydubina/rchan)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/rchan.svg)](https://pkg.go.dev/github.com/nikolaydubina/rchan)

* 30 LOC
* integration test

```go
rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

r, w := rchan.NewRedisListChannel[string](rdb, "my-queue", 10000, 10, time.Millisecond*100)

w <- "hello world 🌏🤍✨"

// ...

fmt.Println(<-r)
// Output: hello world 🌏🤍✨
```

## Benchmarks

```bash
REDIS_HOST=localhost REDIS_PORT=6379 go test -bench=. -benchmem .
```

```
goos: darwin
goarch: arm64
pkg: github.com/nikolaydubina/rchan
BenchmarkSendReceive_string-10             20949             52242 ns/op            1043 B/op         35 allocs/op
BenchmarkSendReceive_bytes-10              25119             55099 ns/op            1073 B/op         36 allocs/op
PASS
ok      github.com/nikolaydubina/rchan  4.831s
```
