package rchan

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisListChannel passes data between input and output channels through Redis list.
// Buffers send and received messages.
// Beyond queue size messages are dropped.
// No retries.
// To terminate, close write channel which will send whatever is buffered to send and send back whatever buffered to read.
func NewRedisListChannel[T string | []byte](rdb *redis.Client, key string, size uint, buff int, poll time.Duration) (<-chan T, chan<- T) {
	r, w, stop := make(chan T, buff), make(chan T, buff), make(chan bool, 1)

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
					if err := send(context.Background(), rdb, key, size, m); err != nil {
						log.Printf("send message(%v) error: %s\n", m, err)
					}
				}
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
	}()

	return r, w
}

func send[T string | []byte](ctx context.Context, rdb *redis.Client, key string, size uint, m T) error {
	tx := rdb.TxPipeline()
	tx.RPush(ctx, key, m)
	tx.LTrim(ctx, key, 0, int64(size))
	_, err := tx.Exec(ctx)
	return err
}
