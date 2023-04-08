package rchan_test

import (
	"context"
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

	// Alice
	w := rchan.NewWriter[string](rdb, "my-queue", 10000, 10, 1)
	w <- "hello world 🌏🤍✨"

	// ... 🌏 ...

	// Bob
	r, _ := rchan.NewReader[string](rdb, "my-queue", 10000, 10, 1, time.Millisecond*100)
	fmt.Println(<-r)
	// Output: hello world 🌏🤍✨
}

// 1. start redis server
// 2. REDIS_HOST=localhost REDIS_PORT=6379 go test -coverprofile rchan.cover ./...
func TestRedisListChannel_SendAndReceive(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	name := "my-queue-test"

	rdb.Del(context.Background(), name)

	r, _ := rchan.NewReader[string](rdb, name, 10000, 10, 1, time.Millisecond*100)
	w := rchan.NewWriter[string](rdb, name, 10000, 10, 1)

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
			for q := range r {
				lock.Lock()
				vs[string(q)]++
				lock.Unlock()
			}
		}()
	}

	for k, v := range counters {
		for i := 0; i < v; i++ {
			w <- k
		}
	}
	close(w)

	time.Sleep(time.Millisecond * 500)

	for k, v := range counters {
		if vs[k] != v {
			t.Errorf("key(%s) exp(%d) != got(%d)", k, v, vs[k])
		}
	}
}

type Counter struct {
	C chan int
	v *int
}

func NewCounter() *Counter { return &Counter{C: make(chan int, 100)} }

func (s *Counter) Incr(v int) { s.C <- v }

func (s *Counter) Get() int {
	if s.v == nil {
		close(s.C)
		var v int
		for q := range s.C {
			v += q
		}
		s.v = &v
	}
	return *s.v
}

func bench[T string | []byte](b *testing.B, name string, buff int, batch int) {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	name = name + "-" + strconv.Itoa(batch)

	wg := &sync.WaitGroup{}
	numBytes := NewCounter()
	numMessages := NewCounter()

	rdb.Del(context.Background(), name)

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			r, stop := rchan.NewReader[T](rdb, name, 1000000, buff, batch, time.Millisecond*100)
			defer stop()

			b, c := 0, 0

			for q := range r {
				if string(q) == "stop" {
					numBytes.Incr(b)
					numMessages.Incr(c)
					return
				}
				b += len(q)
				c++
			}
		}()
	}

	w := rchan.NewWriter[T](rdb, name, 1000000, buff, batch)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		w <- T("blueberry")
	}
	for i := 0; i < 1; i++ {
		w <- T("stop")
	}
	close(w)

	wg.Wait()

	queueLenAfter := rdb.LLen(context.Background(), name).Val()

	b.ReportMetric(float64(numMessages.Get()+int(queueLenAfter))/float64(b.N)*100, "receive+queue/%")
	b.ReportMetric(float64(numMessages.Get())/float64(b.N)*100, "receive/%")
	b.ReportMetric(float64(numBytes.Get())/float64(b.N), "receive_B/op")
	b.ReportMetric(float64(numBytes.Get())/b.Elapsed().Seconds()/float64(1<<20), "receive_MB/s")
}

func BenchmarkBatchSendReceive(b *testing.B) {
	for _, buff := range []int{10, 100, 1000, 10000} {
		for _, batch := range []int{10, 100, 1000, 10000} {
			if buff < batch {
				continue
			}

			b.Run(fmt.Sprintf("buff-%d-batch-%d-string-", buff, batch), func(b *testing.B) {
				bench[string](b, "batch-string", buff, batch)
			})

			b.Run(fmt.Sprintf("buff-%d-batch-%d-bytes--", buff, batch), func(b *testing.B) {
				bench[[]byte](b, "batch-bytes", buff, batch)
			})
		}
	}
}
