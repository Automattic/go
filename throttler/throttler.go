package throttler

import "time"

type Throttler interface {
	GetCount(key string) int64
	AddCount(key string, expiry time.Duration)
	IncrementCount(key string)
	Ban(key string, maxTries int64, expiry time.Duration)
}

func IsThrottled(t Throttler, key string, maxTries int64, expiry, ban time.Duration) bool {
	if count := t.GetCount(key); count == 0 {
		t.AddCount(key, expiry)
		return false
	} else if count < maxTries {
		t.IncrementCount(key)
		return false
	}

	// TODO: keep incrementing ban on repeat requests?
	t.Ban(key, maxTries, ban)

	return true
}
