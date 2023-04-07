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

	r, w := rchan.NewRedisListChannel[string](rdb, "my-queue", 10000, 10, 1, time.Millisecond*100)

	w <- "hello world 🌏🤍✨"

	// ...

	fmt.Println(<-r)
	// Output: hello world 🌏🤍✨
}

// 1. start redis server
// 2. REDIS_HOST=localhost REDIS_PORT=6379 go test -coverprofile rchan.cover ./...
func TestRedisListChannel_SendAndReceive(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	in, out := rchan.NewRedisListChannel[string](rdb, "my-queue-test", 10000, 10, 1, time.Millisecond*100)

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
			out <- k
		}
	}

	time.Sleep(time.Second)

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
	in, out := rchan.NewRedisListChannel[T](rdb, name, 1000000, buff, batch, time.Millisecond*100)

	wg := &sync.WaitGroup{}
	sum := NewCounter()
	count := NewCounter()

	rdb.Del(context.Background(), name)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c, cnt := 0, 0
			for q := range in {
				c += len(q)
				cnt++
			}
			sum.Incr(c)
			count.Incr(cnt)
		}()
	}

	b.ResetTimer()
	for n := 0; n < b.N; n += batch {
		out <- T("blueberry")
	}

	close(out)
	wg.Wait()

	queueLenAfter := rdb.LLen(context.Background(), name).Val()

	b.ReportMetric(float64(count.Get()+int(queueLenAfter))/float64(b.N)*100, "receive+queue/%")
	b.ReportMetric(float64(count.Get())/float64(b.N)*100, "receive/%")
	b.ReportMetric(float64(sum.Get())/float64(b.N), "receive_B/op")
	b.ReportMetric(float64(sum.Get())/b.Elapsed().Seconds()/float64(1<<20), "receive_MB/s")
}

func BenchmarkBatchSendReceive(b *testing.B) {
	for _, buff := range []int{10, 100, 1000, 10000} {
		for _, batch := range []int{10, 100, 1000, 10000} {
			b.Run(fmt.Sprintf("buff_%d_batch_string_%d", buff, batch), func(b *testing.B) {
				bench[string](b, "batch-string", buff, batch)
			})

			b.Run(fmt.Sprintf("buff_%d_batch_%d", buff, batch), func(b *testing.B) {
				bench[[]byte](b, "batch-bytes", buff, batch)
			})
		}
	}
}
