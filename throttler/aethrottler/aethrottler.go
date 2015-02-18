package aethrottler

import (
	"strconv"
	"time"

	"appengine"
	"appengine/memcache"
)

type AppengineThrottler struct {
	context appengine.Context
}

func (t AppengineThrottler) GetCount(key string) int64 {
	item, err := memcache.Get(t.context, key)
	if err != nil {
		return 0
	}

	count, err := strconv.ParseInt(string(item.Value), 10, 64)
	if err != nil {
		return 0
	}

	return count
}

func (t AppengineThrottler) AddCount(key string, expiry time.Duration) {
	err := memcache.Set(t.context, &memcache.Item{
		Key:        key,
		Value:      []byte("1"),
		Expiration: expiry,
	})

	if err != nil {
		t.context.Errorf("Throttler: AddCount failed: %v", err)
	}
}

func (t AppengineThrottler) IncrementCount(key string) {
	_, err := memcache.IncrementExisting(t.context, key, 1)
	if err != nil {
		t.context.Errorf("Throttler: IncrementCount failed: %v", err)
	}
}

func (t AppengineThrottler) Ban(key string, maxTries int64, expiry time.Duration) {
	val := strconv.FormatInt(maxTries, 10)
	err := memcache.Set(t.context, &memcache.Item{
		Key:        key,
		Value:      []byte(val),
		Expiration: expiry,
	})
	if err != nil {
		t.context.Errorf("Throttler: Ban failed: %v", err)
	}
}
