package sessions

import (
	"github.com/boj/redistore"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/sessions"
)

type RediStore interface {
	store
}

// RedisStore instantiates a RediStore with a *redis.Pool passed in.
//
// Ref: https://godoc.org/github.com/boj/redistore#NewRediStoreWithPool
func RedisStore(pool *redis.Pool, keyPairs ...[]byte) (RediStore, error) {
	store, err := redistore.NewRediStoreWithPool(pool, keyPairs...)
	if err != nil {
		return nil, err
	}
	store.SetMaxAge(31556926)

	st := &rediStore{store}
	Store = st
	return st, nil
}

type rediStore struct {
	*redistore.RediStore
}

func (c *rediStore) Options(options Options) {
	c.RediStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
