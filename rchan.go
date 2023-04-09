package rchan

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewReader creates long polling reader from Redis List into channel.
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

	go send(context.Background(), rdb, key, size, batch, back)

	return r, func() { close(cstop) }
}

// NewWriter from chanel into Redis List.
// To terminate, close write channel.
func NewWriter[T string | []byte](rdb *redis.Client, key string, size uint, buff int, batch int) chan<- T {
	w := make(chan T, buff)
	go send(context.Background(), rdb, key, size, batch, w)
	return w
}

func send[T string | []byte](ctx context.Context, rdb *redis.Client, key string, size uint, batch int, w <-chan T) {
	tx := rdb.Pipeline()
	n := 0
	for m := range w {
		tx.RPush(ctx, key, m)
		n++
		if n == batch {
			tx.LTrim(ctx, key, 0, int64(size))
			if _, err := tx.Exec(ctx); err != nil {
				log.Printf("send error: %s\n", err)
			}
			tx = rdb.Pipeline()
			n = 0
		}
	}
	tx.LTrim(ctx, key, 0, int64(size))
	if _, err := tx.Exec(ctx); err != nil {
		log.Printf("send error: %s\n", err)
	}
}
