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
// To terminate, close the write channel,
// this will terminate threads and will send to Redis everything buffered, both to-read and to-send.
func NewRedisListChannel[T string | []byte](rdb *redis.Client, key string, size uint, buff int, batch int, poll time.Duration) (<-chan T, chan<- T) {
	r, w, wi, stop := make(chan T, buff), make(chan T, buff), make(chan T, buff), make(chan bool, 1)

	go func() {
		t := time.NewTicker(poll)
		for {
			select {
			case <-t.C:
				for has := true; has; {
					m, err := rdb.LPop(context.Background(), key).Bytes()
					if has = len(m) > 0 && err == nil; !has {
						if !errors.Is(err, redis.Nil) {
							log.Printf("receive message(%v) error: %s\n", m, err)
						}
						continue
					}
					r <- T(m)
				}
			case <-stop:
				close(r)
				t.Stop()
				for m := range r {
					wi <- m
				}
				return
			}
		}
	}()

	go func() {
		for m := range wi {
			if err := send(context.Background(), rdb, key, size, nil, m); err != nil {
				log.Printf("send message(%v) error: %s\n", m, err)
			}
		}
	}()

	go func() {
		ms := make([]T, 0, batch)
		for m := range w {
			if err := send(context.Background(), rdb, key, size, wi, m); err != nil {
				log.Printf("send message(%v) error: %s\n", m, err)
			}
		}
		stop <- true
	}()

	return r, w
}

func send[T string | []byte](ctx context.Context, rdb *redis.Client, key string, size uint, w chan<- T, ms ...T) error {
	tx := rdb.TxPipeline()
	for _, m := range ms {
		tx.RPush(ctx, key, m)
	}
	tx.LTrim(ctx, key, 0, int64(size))
	rs, err := tx.Exec(ctx)
	if err != nil && w != nil {
		for i, r := range rs {
			if r.Err() != nil {
				w <- ms[i]
			}
		}
	}
	return err
}
