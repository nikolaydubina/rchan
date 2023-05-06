package rchan

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewReader creates long polling reader from Redis List into channel.
// Every pooling interval tries to populate buffer.
// To terminate call stop function, which will send back to Redis List anything buffered.
func NewReader[T string | []byte](rdb *redis.Client, key string, size uint, buff int, batch int, poll time.Duration) (read <-chan T, stop func()) {
	r, back, cstop := make(chan T, buff), make(chan T, buff), make(chan bool)

	go func() {
		t := time.NewTicker(poll)
		for {
			select {
			case <-t.C:
				for has := len(r) < cap(r); has && len(r) < cap(r); {
					mb, err := rdb.LPopCount(context.Background(), key, batch).Result()
					if err != nil && !errors.Is(err, redis.Nil) {
						log.Printf("receive messages_count(%d) error: %s\n", len(mb), err)
					}
					if has = len(mb) > 0; !has {
						continue
					}
					for _, m := range mb {
						r <- T(m)
					}
				}
			case <-cstop:
				close(r)
				for m := range r {
					back <- m
				}
				close(back)
				return
			}
		}
	}()

	go send(context.Background(), rdb, key, size, batch, poll, back)

	return r, func() { close(cstop) }
}

// NewWriter from chanel into Redis List.
// Flushes pending items on interval or when batch size is reached.
// To terminate, close write channel.
func NewWriter[T string | []byte](rdb *redis.Client, key string, size uint, buff int, batch int, flushInterval time.Duration) chan<- T {
	w := make(chan T, buff)
	go send(context.Background(), rdb, key, size, batch, flushInterval, w)
	return w
}

func send[T string | []byte](ctx context.Context, rdb *redis.Client, key string, size uint, batch int, flushInterval time.Duration, w <-chan T) {
	n, tx, t := 0, rdb.Pipeline(), time.NewTicker(flushInterval)
	for {
		select {
		case <-t.C:
			if n > 0 {
				flushSend(ctx, tx, size, key)
				n, tx = 0, rdb.Pipeline()
			}
		case m, ok := <-w:
			if !ok {
				flushSend(ctx, tx, size, key)
				return
			}
			tx.RPush(ctx, key, m)
			n++
			if n == batch {
				flushSend(ctx, tx, size, key)
				n, tx = 0, rdb.Pipeline()
			}
		}
	}
}

func flushSend(ctx context.Context, tx redis.Pipeliner, size uint, key string) {
	tx.LTrim(ctx, key, 0, int64(size))
	if _, err := tx.Exec(ctx); err != nil {
		log.Printf("send error: %s\n", err)
	}
}
