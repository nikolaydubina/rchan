package rchan

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisListChannel passes data between input and output channels through Redis list.
// Fills out input buffer channel on poll interval.
// Output is buffered too.
// When goroutine is terminated, then output buffered messages are lost.
// Beyond queue size messages are dropped.
// No retries.
// To terminate, close write channel.
func NewRedisListChannel[T string | []byte](rdb *redis.Client, key string, size uint, buff int, poll time.Duration) (<-chan T, chan<- T) {
	r, w := make(chan T, buff), make(chan T, buff)

	go func() {
		ctx := context.Background()
		t := time.NewTicker(poll)
		for {
			select {
			case <-t.C:
				for has := true; has; {
					m, err := rdb.LPop(ctx, key).Bytes()
					if len(m) == 0 || err != nil {
						if err != nil && !errors.Is(err, redis.Nil) {
							log.Printf("receive message(%v) error: %s\n", m, err)
						}
						has = false
						continue
					}
					r <- T(m)
				}
			case m, ok := <-w:
				if !ok {
					close(r)
					return
				}
				tx := rdb.TxPipeline()
				tx.RPush(ctx, key, m)
				tx.LTrim(ctx, key, 0, int64(size))
				if _, err := tx.Exec(ctx); err != nil {
					log.Printf("send message(%v) error: %s\n", m, err)
				}
			}
		}
	}()

	return r, w
}
