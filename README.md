## 🌸 rchan: Go channel through Redis List

[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/rchan)](https://goreportcard.com/report/github.com/nikolaydubina/rchan)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/rchan.svg)](https://pkg.go.dev/github.com/nikolaydubina/rchan)

* 40 LOC
* 30 _thousand_ RPS (send individual message)
* 1.4 _million_ RPS (send batch pipeline)
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
BenchmarkSendReceive_string-10    28257    35575 ns/op     9.003 receive_B/op     0.2414 receive_MB/s  1068 B/op   36 allocs/op
BenchmarkSendReceive_bytes-10     28677    35678 ns/op     9.003 receive_B/op     0.2407 receive_MB/s  1060 B/op   35 allocs/op
BenchmarkBatchSendReceive/batch_string_10-10     333405     3653 ns/op     9.843 receive_B/op     2.57 receive_MB/s     344 B/op    9 allocs/op
BenchmarkBatchSendReceive/batch_bytes_10-10      331542     3692 ns/op     9.134 receive_B/op     2.35 receive_MB/s     366 B/op   10 allocs/op
BenchmarkBatchSendReceive/batch_string_100-10   1000000     1085 ns/op     9.162 receive_B/op     8.05 receive_MB/s     280 B/op    7 allocs/op
BenchmarkBatchSendReceive/batch_bytes_100-10    1000000     1087 ns/op     8.476 receive_B/op     7.43 receive_MB/s     301 B/op    8 allocs/op
BenchmarkBatchSendReceive/batch_string_1000-10  1761710      674 ns/op     8.956 receive_B/op    12.65 receive_MB/s     283 B/op    7 allocs/op
BenchmarkBatchSendReceive/batch_bytes_1000-10   1750096      699 ns/op     9.467 receive_B/op    12.91 receive_MB/s     310 B/op    8 allocs/op
BenchmarkBatchSendReceive/batch_string_10000-10 1899872      649 ns/op     8.953 receive_B/op    13.14 receive_MB/s     300 B/op    6 allocs/op
BenchmarkBatchSendReceive/batch_bytes_10000-10  1598560      663 ns/op     8.839 receive_B/op    12.71 receive_MB/s     324 B/op    7 allocs/op
PASS
ok   github.com/nikolaydubina/rchan  16.378s
```
