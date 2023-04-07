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
func NewRedisListChannel[T string | []byte](rdb *redis.Client, key string, size uint, buff int, poll time.Duration) (<-chan T, chan<- T) {
	r, w, stop := make(chan T, buff), make(chan T, buff), make(chan bool, 1)

	go func() {
		t := time.NewTicker(poll)
		for {
			select {
			case <-t.C:
				for has := len(r) < cap(r); has && len(r) < cap(r); {
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
				return
			}
		}
	}()

	go func() {
		for m := range w {
			if err := send(context.Background(), rdb, key, size, m); err != nil {
				log.Printf("send message(%v) error: %s\n", m, err)
			}
		}
		stop <- true
		for m := range r {
			if err := send(context.Background(), rdb, key, size, m); err != nil {
				log.Printf("send message(%v) error: %s\n", m, err)
			}
		}
	}()

	return r, w
}

// NewBatchRedisListChannel same as NewRedisListChannel, but sends batches.
func NewBatchRedisListChannel[T string | []byte](rdb *redis.Client, key string, size uint, buff int, poll time.Duration) (<-chan T, chan<- []T) {
	r, w, stop := make(chan T, buff), make(chan []T, buff), make(chan bool, 1)

	go func() {
		t := time.NewTicker(poll)
		for {
			select {
			case <-t.C:
				for has := len(r) < cap(r); has && len(r) < cap(r); {
					mb, err := rdb.LPopCount(context.Background(), key, buff).Result()
					if has = len(mb) > 0 && err == nil; !has {
						if !errors.Is(err, redis.Nil) {
							log.Printf("receive messages_count(%d) error: %s\n", len(mb), err)
						}
						continue
					}
					for _, m := range mb {
						r <- T(m)
					}
				}
			case <-stop:
				close(r)
				t.Stop()
				return
			}
		}
	}()

	go func() {
		for mb := range w {
			if err := send(context.Background(), rdb, key, size, mb...); err != nil {
				log.Printf("send message(%v) error: %s\n", mb, err)
			}
		}
		stop <- true
		for m := range r {
			if err := send(context.Background(), rdb, key, size, m); err != nil {
				log.Printf("send message(%v) error: %s\n", m, err)
			}
		}
	}()

	return r, w
}

func send[T string | []byte](ctx context.Context, rdb *redis.Client, key string, size uint, ms ...T) error {
	tx := rdb.TxPipeline()
	for _, m := range ms {
		tx.RPush(ctx, key, m)
	}
	tx.LTrim(ctx, key, 0, int64(size))
	_, err := tx.Exec(ctx)
	return err
}
