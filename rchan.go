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
func NewRedisListChannel[T string | []byte](redis *redis.Client, key string, size uint, buff int, poll time.Duration) (<-chan T, chan<- T) {
	s := rc[T]{
		redis: redis,
		key:   key,
		size:  size,
		poll:  poll,
		in:    make(chan T, buff),
		out:   make(chan T, buff),
	}
	go s.run()
	return s.out, s.in
}

type rc[T []byte | string] struct {
	redis *redis.Client
	key   string
	size  uint
	poll  time.Duration
	in    chan T
	out   chan T
}

func (s *rc[T]) run() {
	ctx := context.Background()
	t := time.NewTicker(s.poll)
	for {
		select {
		case <-t.C:
			for has := true; has; {
				m, err := s.redis.LPop(ctx, s.key).Bytes()
				if len(m) == 0 || err != nil {
					if err != nil && !errors.Is(err, redis.Nil) {
						log.Printf("receive error: %s\n", err)
					}
					has = false
					continue
				}
				s.out <- T(m)
			}
		case m, ok := <-s.in:
			if !ok {
				close(s.out)
				return
			}
			pipe := s.redis.TxPipeline()
			pipe.RPush(ctx, s.key, m)
			pipe.LTrim(ctx, s.key, 0, int64(s.size))
			if _, err := pipe.Exec(ctx); err != nil {
				log.Printf("send error: %s\n", err)
			}
		}
	}
}
