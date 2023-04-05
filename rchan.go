package rchan

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisListChannel passes data between input and output channels through Redis list.
// Always reads values from the Redis and fills out the read channel buffer.
// When sending fills out out channel buffer.
// If process is restarted or value is not read, then values are lost.
// If nothing to read, then blocks by timer for poll duration.
// Beyond size messages are dropped.
func NewRedisListChannel[T string | []byte](r *redis.Client, key string, size uint, buff int, poll time.Duration) (<-chan T, chan<- T) {
	out, in := make(chan T, buff), make(chan T, buff)

	go func() {
		ctx := context.Background()
		t := time.NewTicker(poll)
		for {
			select {
			case <-t.C:
				for has := true; has; {
					m, err := r.LPop(ctx, key).Bytes()
					if len(m) == 0 || err != nil {
						if err != nil && !errors.Is(err, redis.Nil) {
							log.Printf("receive error: %s\n", err)
						}
						has = false
						continue
					}
					in <- T(m)
				}
			case m, ok := <-out:
				if !ok {
					close(in)
					return
				}
				pipe := r.TxPipeline()
				pipe.RPush(ctx, key, m)
				pipe.LTrim(ctx, key, 0, int64(size))
				if _, err := pipe.Exec(ctx); err != nil {
					log.Printf("send error: %s\n", err)
				}
			}
		}
	}()

	return in, out
}
