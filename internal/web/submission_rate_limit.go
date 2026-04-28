package web

import (
	"sync"
	"time"
)

type submissionRateLimiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	now    func() time.Time
	hits   map[string][]time.Time
}

func newSubmissionRateLimiter(limit int, window time.Duration) *submissionRateLimiter {
	return &submissionRateLimiter{
		limit:  limit,
		window: window,
		now:    time.Now,
		hits:   make(map[string][]time.Time),
	}
}

func (l *submissionRateLimiter) Allow(key string) bool {
	if l == nil || l.limit <= 0 || l.window <= 0 {
		return true
	}

	now := l.now
	if now == nil {
		now = time.Now
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	current := now()
	cutoff := current.Add(-l.window)
	timestamps := l.hits[key]
	kept := timestamps[:0]
	for _, ts := range timestamps {
		if !ts.Before(cutoff) {
			kept = append(kept, ts)
		}
	}
	if len(kept) >= l.limit {
		l.hits[key] = append([]time.Time(nil), kept...)
		return false
	}

	kept = append(kept, current)
	l.hits[key] = append([]time.Time(nil), kept...)
	return true
}
