package rchan

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisListChannel passes data between input and output channels through Redis List.
// Buffers to-sent and to-read messages.
// Beyond queue size messages are dropped.
// No retries.
// To terminate, close the write channel, this will terminate threads and send back to Redis everything buffered, both to-read and to-send.
func NewRedisListChannel[T string | []byte](rdb *redis.Client, key string, size uint, buff int, batch int, poll time.Duration) (<-chan T, chan<- T) {
	r, w, stop := make(chan T, buff), make(chan T, buff), make(chan bool)

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
			case <-stop:
				close(r)
				return
			}
		}
	}()

	go func() {
		send(context.Background(), rdb, key, size, batch, w)
		stop <- true
		send(context.Background(), rdb, key, size, batch, r)
	}()

	return r, w
}

func send[T string | []byte](ctx context.Context, rdb *redis.Client, key string, size uint, batch int, w <-chan T) {
	tx := rdb.Pipeline()
	n := 0
	for m := range w {
		if n == batch {
			tx.LTrim(ctx, key, 0, int64(size))
			if _, err := tx.Exec(ctx); err != nil {
				log.Printf("send error: %s\n", err)
			}
			tx = rdb.Pipeline()
			n = 0
		}
		tx.RPush(ctx, key, m)
		n++
	}
	tx.LTrim(ctx, key, 0, int64(size))
	if _, err := tx.Exec(ctx); err != nil {
		log.Printf("send error: %s\n", err)
	}
}
