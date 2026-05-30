package cache

import (
	"strings"
	"sync"
	"time"
	"unicode"
)

type entry struct {
	value     string
	expiresAt time.Time
}

type ResponseCache struct {
	mu      sync.RWMutex
	entries map[string]entry
	ttl     time.Duration
}

func New(ttl time.Duration) *ResponseCache {
	c := &ResponseCache{
		entries: make(map[string]entry),
		ttl:     ttl,
	}
	go c.evict()
	return c
}

func (c *ResponseCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok || time.Now().After(e.expiresAt) {
		return "", false
	}
	return e.value, true
}

func (c *ResponseCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry{value: value, expiresAt: time.Now().Add(c.ttl)}
}

// NormalizeKey lowercases, strips punctuation, and collapses whitespace
// so "When do you open??" and "when do you open" map to the same key.
func NormalizeKey(msg string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(msg) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func (c *ResponseCache) evict() {
	ticker := time.NewTicker(15 * time.Minute)
	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		for k, e := range c.entries {
			if now.After(e.expiresAt) {
				delete(c.entries, k)
			}
		}
		c.mu.Unlock()
	}
}
