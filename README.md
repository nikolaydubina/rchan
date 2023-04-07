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

### localhost

8.7 GB/s

```bash
iperf3 -s -p 3000  
```

```
-----------------------------------------------------------
Server listening on 3000 (test #1)
-----------------------------------------------------------
Accepted connection from ::1, port 59020
[  5] local ::1 port 3000 connected to ::1 port 59021
[ ID] Interval           Transfer     Bitrate
[  5]   0.00-1.00   sec  7.43 GBytes  63.8 Gbits/sec                  
[  5]   1.00-2.00   sec  8.39 GBytes  72.0 Gbits/sec                  
[  5]   2.00-3.00   sec  8.79 GBytes  75.5 Gbits/sec                  
[  5]   3.00-4.00   sec  8.71 GBytes  74.8 Gbits/sec                  
[  5]   4.00-5.00   sec  8.75 GBytes  75.2 Gbits/sec                  
[  5]   5.00-6.00   sec  8.76 GBytes  75.3 Gbits/sec                  
[  5]   6.00-7.00   sec  8.26 GBytes  70.9 Gbits/sec                  
[  5]   7.00-8.00   sec  8.62 GBytes  74.0 Gbits/sec                  
[  5]   8.00-9.00   sec  8.49 GBytes  72.9 Gbits/sec                  
[  5]   9.00-10.00  sec  8.85 GBytes  76.0 Gbits/sec                  
[  5]  10.00-10.00  sec  2.06 MBytes  56.7 Gbits/sec                  
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bitrate
[  5]   0.00-10.00  sec  85.0 GBytes  73.0 Gbits/sec                  receiver
-----------------------------------------------------------
```

```bash
iperf3 -c localhost -p 3000 -f M
```

```
Connecting to host localhost, port 3000
[  7] local ::1 port 59021 connected to ::1 port 3000
[ ID] Interval           Transfer     Bitrate
[  7]   0.00-1.00   sec  7.43 GBytes  7608 MBytes/sec                  
[  7]   1.00-2.00   sec  8.39 GBytes  8588 MBytes/sec                  
[  7]   2.00-3.00   sec  8.79 GBytes  8998 MBytes/sec                  
[  7]   3.00-4.00   sec  8.71 GBytes  8916 MBytes/sec                  
[  7]   4.00-5.00   sec  8.75 GBytes  8961 MBytes/sec                  
[  7]   5.00-6.00   sec  8.76 GBytes  8972 MBytes/sec                  
[  7]   6.00-7.00   sec  8.26 GBytes  8456 MBytes/sec                  
[  7]   7.00-8.00   sec  8.62 GBytes  8823 MBytes/sec                  
[  7]   8.00-9.00   sec  8.49 GBytes  8691 MBytes/sec                  
[  7]   9.00-10.00  sec  8.85 GBytes  9064 MBytes/sec                  
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bitrate
[  7]   0.00-10.00  sec  85.0 GBytes  8708 MBytes/sec                  sender
[  7]   0.00-10.00  sec  85.0 GBytes  8708 MBytes/sec                  receiver

iperf Done.
```
