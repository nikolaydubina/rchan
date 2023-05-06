## rchan: Go channel through Redis List

[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/rchan)](https://goreportcard.com/report/github.com/nikolaydubina/rchan)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/rchan.svg)](https://pkg.go.dev/github.com/nikolaydubina/rchan)

> Disclaimer: this is experiment to see channel based API for distributed queues. Channel technically is a function call that does not error out. So effectively sending to or reading from distributed queue through channel is _delaying_ or _buffering_ calls to distributied message broker. It is possible that this can lead to better performance and higher throughput. However, this project is more of an experiment of API design with goal to make workable solution with minimal code. Very likely you may want to use function based API in production.

* 100 LOC
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
BenchmarkBatchSendReceive/buff-10-batch-10-string--10          287922     4221 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     2.033 receive_MB/s    310 B/op    8 allocs/op   
BenchmarkBatchSendReceive/buff-10-batch-10-bytes---10          249921     4009 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     2.141 receive_MB/s    350 B/op   10 allocs/op   
BenchmarkBatchSendReceive/buff-100-batch-10-string--10         280068     4131 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     2.078 receive_MB/s    310 B/op    8 allocs/op   
BenchmarkBatchSendReceive/buff-100-batch-10-bytes---10         282595     4182 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     2.052 receive_MB/s    350 B/op   10 allocs/op   
BenchmarkBatchSendReceive/buff-100-batch-100-string--10       1021838     1108 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     7.745 receive_MB/s    256 B/op    6 allocs/op   
BenchmarkBatchSendReceive/buff-100-batch-100-bytes---10        939798     1172 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     7.320 receive_MB/s    296 B/op    8 allocs/op   
BenchmarkBatchSendReceive/buff-1000-batch-10-string--10        286066     4069 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     2.109 receive_MB/s    310 B/op    8 allocs/op   
BenchmarkBatchSendReceive/buff-1000-batch-10-bytes---10        284227     4230 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     2.029 receive_MB/s    350 B/op   10 allocs/op   
BenchmarkBatchSendReceive/buff-1000-batch-100-string--10       894435     1172 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     7.326 receive_MB/s    256 B/op    6 allocs/op   
BenchmarkBatchSendReceive/buff-1000-batch-100-bytes---10       964352     1156 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     7.423 receive_MB/s    296 B/op    8 allocs/op   
BenchmarkBatchSendReceive/buff-1000-batch-1000-string--10     1444040      767.7 ns/op    100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op    11.18 receive_MB/s     259 B/op    6 allocs/op   
BenchmarkBatchSendReceive/buff-1000-batch-1000-bytes---10     1451476      762.2 ns/op    100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op    11.26 receive_MB/s     299 B/op    8 allocs/op   
BenchmarkBatchSendReceive/buff-10000-batch-10-string--10       267046     4113 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     2.087 receive_MB/s    310 B/op    8 allocs/op   
BenchmarkBatchSendReceive/buff-10000-batch-10-bytes---10       281936     4043 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     2.123 receive_MB/s    352 B/op   10 allocs/op   
BenchmarkBatchSendReceive/buff-10000-batch-100-string--10     1039837     1095 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     7.835 receive_MB/s    256 B/op    6 allocs/op   
BenchmarkBatchSendReceive/buff-10000-batch-100-bytes---10      974930     1144 ns/op      100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op     7.503 receive_MB/s    296 B/op    8 allocs/op   
BenchmarkBatchSendReceive/buff-10000-batch-1000-string--10    1462742      767.1 ns/op    100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op    11.19 receive_MB/s     259 B/op    6 allocs/op   
BenchmarkBatchSendReceive/buff-10000-batch-1000-bytes---10    1363063      811.4 ns/op    100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op    10.58 receive_MB/s     299 B/op    8 allocs/op   
BenchmarkBatchSendReceive/buff-10000-batch-10000-string--10   1622870      752.2 ns/op    100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op    11.41 receive_MB/s     277 B/op    6 allocs/op   
BenchmarkBatchSendReceive/buff-10000-batch-10000-bytes---10   1630532      743.0 ns/op    100.0 receive+queue/%     100.0 receive/%     9.000 receive_B/op    11.55 receive_MB/s     317 B/op    8 allocs/op 
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
