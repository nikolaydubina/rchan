package rchan_test

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/nikolaydubina/rchan"
	"github.com/redis/go-redis/v9"
)

func Example_simple() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	r, w := rchan.NewRedisListChannel[string](rdb, "my-queue", 10000, 10, time.Millisecond*100)

	w <- "hello world 🌏🤍✨"

	// ... 🗺️ ⏳ ...

	fmt.Println(<-r)
	// Output: hello world 🌏🤍✨
}

// 1. start redis server
// 2. REDIS_HOST=localhost REDIS_PORT=6379 go test -coverprofile rchan.cover ./...
func TestRedisListChannel_SendAndReceive(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	in, out := rchan.NewRedisListChannel[[]byte](rdb, "my-queue-test", 10000, 10, time.Millisecond*100)

	counters := map[string]int{
		"a":                     100,
		"b":                     5,
		"c":                     17,
		"some-very-long-string": 1111,
	}

	vs := map[string]int{}
	lock := sync.RWMutex{}

	for i := 0; i < 5; i++ {
		go func() {
			for q := range in {
				lock.Lock()
				vs[string(q)]++
				lock.Unlock()
			}
		}()
	}

	for k, v := range counters {
		for i := 0; i < v; i++ {
			out <- []byte(k)
		}
	}

	time.Sleep(time.Second)

	for k, v := range counters {
		if vs[k] != v {
			t.Errorf("key(%s) exp(%d) != got(%d)", k, v, vs[k])
		}
	}
}

func bench[T string | []byte](b *testing.B, name string) {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	in, out := rchan.NewRedisListChannel[T](rdb, name, 100000, 10000, time.Millisecond*100)

	sumch := make(chan int, 5)
	wg := &sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			c := 0
			defer wg.Done()
			for q := range in {
				c += len(q)
			}
			sumch <- c
		}()
	}

	v := T("blueberry")

	b.ResetTimer()
	for n := 0; n < b.N+10; n++ {
		out <- v
	}

	close(out)

	wg.Wait()

	close(sumch)
	sum := 0
	for q := range sumch {
		sum += q
	}

	b.ReportMetric(float64(sum)/float64(b.N), "receive_B/op")
	b.ReportMetric(float64(sum)/b.Elapsed().Seconds()/float64(1<<20), "receive_MB/s")
}

func BenchmarkSendReceive_string(b *testing.B) { bench[string](b, "bench-string") }

func BenchmarkSendReceive_bytes(b *testing.B) { bench[[]byte](b, "bench-bytes") }

func benchBatch[T string | []byte](b *testing.B, name string, batch int) {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	name = name + "-" + strconv.Itoa(batch)
	in, out := rchan.NewBatchRedisListChannel[T](rdb, name, 1000000, batch, time.Millisecond*100)

	sumch := make(chan int, 5)
	wg := &sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := 0
			for q := range in {
				c += len(q)
			}
			sumch <- c
		}()
	}

	var msgs []T
	for i := 0; i < batch; i++ {
		msgs = append(msgs, T("blueberry"))
	}

	b.ResetTimer()
	for n := 0; n < b.N; n += batch {
		out <- msgs
	}

	close(out)

	wg.Wait()

	close(sumch)
	sum := 0
	for q := range sumch {
		sum += q
	}

	b.ReportMetric(float64(sum)/float64(b.N), "receive_B/op")
	b.ReportMetric(float64(sum)/b.Elapsed().Seconds()/float64(1<<20), "receive_MB/s")
}

func BenchmarkBatchSendReceive(b *testing.B) {
	for _, n := range []int{10, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("batch_string_%d", n), func(b *testing.B) {
			benchBatch[string](b, "batch-string", n)
		})

		b.Run(fmt.Sprintf("batch_bytes__%d", n), func(b *testing.B) {
			benchBatch[[]byte](b, "batch-bytes", n)
		})
	}
}
