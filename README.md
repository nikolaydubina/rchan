## 🌸 rchan: Go channel through Redis List

[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/rchan)](https://goreportcard.com/report/github.com/nikolaydubina/rchan)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/rchan.svg)](https://pkg.go.dev/github.com/nikolaydubina/rchan)

> Disclaimer: this is experiment to see channel based API for distributed queues. Channel technically is a function call that does not error out. So effectively sending to or reading from distributed queue through channel is _delaying_ or _buffering_ calls to distributied message broker. It is possible that this can lead to better performance and higher throughput. However, this project is more of experiment of API design and trying out definig APIs in terms of channels. Very likely you may want to use function based API in production.

* 50 LOC
* 30 _thousand_ RPS (send/receive individual message)
* 1.4 _million_ RPS (send/receive batch pipeline)
* graceful stop (no messages lost, unless error)
* integration test

```go
rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

// Alice
w := rchan.NewWriter[string](rdb, "my-queue", 10000, 10, 1)
w <- "hello world 🌏🤍✨"

// ... 🌏 ...

// Bob
r, _ := rchan.NewReader[string](rdb, "my-queue", 10000, 10, 1, time.Millisecond*100)
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
BenchmarkBatchSendReceive/buff-10-batch-10-string--10          281578  4043 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  2.12 receive_MB/s 310 B/op  8 allocs/op
BenchmarkBatchSendReceive/buff-10-batch-10-bytes---10          279426  4081 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  2.10 receive_MB/s 350 B/op 10 allocs/op
BenchmarkBatchSendReceive/buff-100-batch-10-string--10         282370  4131 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  2.07 receive_MB/s 310 B/op  8 allocs/op
BenchmarkBatchSendReceive/buff-100-batch-10-bytes---10         280876  4049 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  2.12 receive_MB/s 350 B/op 10 allocs/op
BenchmarkBatchSendReceive/buff-100-batch-100-string--10       1024000  1092 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  7.85 receive_MB/s 256 B/op  6 allocs/op
BenchmarkBatchSendReceive/buff-100-batch-100-bytes---10        988862  1129 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  7.60 receive_MB/s 296 B/op  8 allocs/op
BenchmarkBatchSendReceive/buff-1000-batch-10-string--10        281828  4041 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  2.12 receive_MB/s 310 B/op  8 allocs/op
BenchmarkBatchSendReceive/buff-1000-batch-10-bytes---10        280272  4212 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  2.03 receive_MB/s 350 B/op 10 allocs/op
BenchmarkBatchSendReceive/buff-1000-batch-100-string--10      1018701  1095 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  7.84 receive_MB/s 256 B/op  6 allocs/op
BenchmarkBatchSendReceive/buff-1000-batch-100-bytes---10      1008822  1109 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  7.74 receive_MB/s 296 B/op  8 allocs/op
BenchmarkBatchSendReceive/buff-1000-batch-1000-string--10     1566339   713 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op 12.02 receive_MB/s 259 B/op  6 allocs/op
BenchmarkBatchSendReceive/buff-1000-batch-1000-bytes---10     1428948   708 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op 12.12 receive_MB/s 299 B/op  8 allocs/op
BenchmarkBatchSendReceive/buff-10000-batch-10-string--10       271543  4116 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  2.08 receive_MB/s 311 B/op  8 allocs/op
BenchmarkBatchSendReceive/buff-10000-batch-10-bytes---10       274959  4100 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  2.09 receive_MB/s 352 B/op 10 allocs/op
BenchmarkBatchSendReceive/buff-10000-batch-100-string--10     1026082  1087 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  7.89 receive_MB/s 256 B/op  6 allocs/op
BenchmarkBatchSendReceive/buff-10000-batch-100-bytes---10     1002680  1115 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op  7.69 receive_MB/s 296 B/op  8 allocs/op
BenchmarkBatchSendReceive/buff-10000-batch-1000-string--10    1438033   704 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op 12.18 receive_MB/s 259 B/op  6 allocs/op
BenchmarkBatchSendReceive/buff-10000-batch-1000-bytes---10    1436354   701 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op 12.23 receive_MB/s 299 B/op  8 allocs/op
BenchmarkBatchSendReceive/buff-10000-batch-10000-string--10   1480132   687 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op 12.49 receive_MB/s 277 B/op  6 allocs/op
BenchmarkBatchSendReceive/buff-10000-batch-10000-bytes---10   1585638   646 ns/op 100.0 receive+queue/%  100.0 receive/%  9.000 receive_B/op 13.29 receive_MB/s 317 B/op  8 allocs/op
PASS
ok   github.com/nikolaydubina/rchan  40.939s
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

## Appendix A: Paths Not Taken

### Single function that returns writer and reader channels
Some call-sites do not use reader at all.
Usually it is either only reader or writer, if you have both, you can use normal channel.
Returning both reader and writer causes to read from Redis List into anonymous reader channel that is not used.
It keeps data there up to buffer size and it is not processed by anyone.
So to keep logic simple, instantiate what you need — either reader or writer.
This also helps to avoid duplication of parameters of buffering for reader vs writer.

### Single channel return for reader
We need to tell reader goroutine that no more reads will ever be performed and send anything in channel back to Redis List.
For example, for graceful termination.
