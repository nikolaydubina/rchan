## 🌸 rchan: Go channel through Redis List

[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/rchan)](https://goreportcard.com/report/github.com/nikolaydubina/rchan)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/rchan.svg)](https://pkg.go.dev/github.com/nikolaydubina/rchan)

* 40 LOC
* 30 _thousand_ RPS (send/receive individual message)
* 1.4 _million_ RPS (send/receive batch pipeline)
* graceful stop (no messages lost, unless error)
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
BenchmarkSendReceive_string-10 29498 34788 ns/op   100.0 receive/%  9.003 receive_B/op  0.2468 receive_MB/s  1068 B/op   36 allocs/op
BenchmarkSendReceive_bytes-10  30074 34886 ns/op   100.0 receive/%  9.003 receive_B/op  0.2461 receive_MB/s  1060 B/op   35 allocs/op
BenchmarkBatchSendReceive/batch_string_10-10     334543  4325 ns/op   100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  1.98 receive_MB/s  361 B/op 10 allocs/op
BenchmarkBatchSendReceive/batch_bytes__10-10     330771  4344 ns/op   100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  1.97 receive_MB/s  385 B/op 11 allocs/op
BenchmarkBatchSendReceive/batch_string_100-10   1000000  1085 ns/op   100.0 receive+queue/%   93.6 receive/%  8.428 receive_B/op  7.40 receive_MB/s  281 B/op  7 allocs/op
BenchmarkBatchSendReceive/batch_bytes__100-10   1000000  1084 ns/op   100.0 receive+queue/%   93.4 receive/%  8.411 receive_B/op  7.40 receive_MB/s  304 B/op  8 allocs/op
BenchmarkBatchSendReceive/batch_string_1000-10  1797055   667 ns/op   100.1 receive+queue/%   96.7 receive/%  8.709 receive_B/op 12.45 receive_MB/s  283 B/op  7 allocs/op
BenchmarkBatchSendReceive/batch_bytes__1000-10  1787865   678 ns/op   100.0 receive+queue/%  100.0 receive/%  9.001 receive_B/op 12.65 receive_MB/s  307 B/op  8 allocs/op
BenchmarkBatchSendReceive/batch_string_10000-10 1860106   661 ns/op   100.5 receive+queue/%  100.5 receive/%  9.048 receive_B/op 13.04 receive_MB/s  302 B/op  7 allocs/op
BenchmarkBatchSendReceive/batch_bytes__10000-10 1787569   685 ns/op   100.1 receive+queue/%  100.1 receive/%  9.012 receive_B/op 12.54 receive_MB/s  325 B/op  8 allocs/op
PASS
ok   github.com/nikolaydubina/rchan  17.187s
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
