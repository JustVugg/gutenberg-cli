package forge

import (
	"errors"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Resilience: token-bucket rate limiter + retry-with-backoff + circuit breaker.
// Configurable via env:
//   GITHUB_RATE_LIMIT_RPS   tokens per second (0 disables; default 0)
//   GITHUB_RATE_BURST       burst size (default rps or 1)
//   GITHUB_RETRY_MAX        max attempts including the first (default 3)
//   GITHUB_RETRY_BASE_MS    base backoff in ms (default 200)
//   GITHUB_CB_THRESHOLD     consecutive failures to open the breaker (default 5)
//   GITHUB_CB_COOLDOWN_MS   cooldown when open (default 30000)

type tokenBucket struct {
	mu       sync.Mutex
	rate     float64
	burst    float64
	tokens   float64
	last     time.Time
}

func newTokenBucket(rps, burst float64) *tokenBucket {
	if burst < 1 {
		burst = 1
	}
	return &tokenBucket{rate: rps, burst: burst, tokens: burst, last: time.Now()}
}

func (b *tokenBucket) Take() {
	if b == nil || b.rate <= 0 {
		return
	}
	for {
		b.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(b.last).Seconds()
		b.tokens = math.Min(b.burst, b.tokens+elapsed*b.rate)
		b.last = now
		if b.tokens >= 1 {
			b.tokens--
			b.mu.Unlock()
			return
		}
		needed := (1 - b.tokens) / b.rate
		b.mu.Unlock()
		time.Sleep(time.Duration(needed * float64(time.Second)))
	}
}

type circuitBreaker struct {
	mu        sync.Mutex
	threshold int
	cooldown  time.Duration
	failures  int
	openedAt  time.Time
}

func (c *circuitBreaker) Allow() error {
	if c == nil || c.threshold <= 0 {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.failures >= c.threshold {
		if time.Since(c.openedAt) < c.cooldown {
			return errors.New("github: circuit breaker open")
		}
		c.failures = 0
	}
	return nil
}

func (c *circuitBreaker) onSuccess() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures = 0
}

func (c *circuitBreaker) onFailure() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures++
	if c.failures >= c.threshold {
		c.openedAt = time.Now()
	}
}

var (
	resilienceOnce   sync.Once
	resilienceBucket *tokenBucket
	resilienceCB     *circuitBreaker
	resilienceMax    int
	resilienceBaseMs int
)

func initResilience() {
	rps := envFloat("GITHUB_RATE_LIMIT_RPS", 0)
	burst := envFloat("GITHUB_RATE_BURST", rps)
	if burst <= 0 {
		burst = 1
	}
	resilienceBucket = newTokenBucket(rps, burst)
	resilienceMax = envInt("GITHUB_RETRY_MAX", 3)
	if resilienceMax < 1 {
		resilienceMax = 1
	}
	resilienceBaseMs = envInt("GITHUB_RETRY_BASE_MS", 200)
	cbThreshold := envInt("GITHUB_CB_THRESHOLD", 5)
	cbCooldownMs := envInt("GITHUB_CB_COOLDOWN_MS", 30000)
	resilienceCB = &circuitBreaker{threshold: cbThreshold, cooldown: time.Duration(cbCooldownMs) * time.Millisecond}
}

// DoWithResilience wraps an HTTP call with rate limiting, retries, and a circuit breaker.
func DoWithResilience(req *http.Request) (*http.Response, error) {
	resilienceOnce.Do(initResilience)
	if resp, ok, err := snapshotIntercept(req); ok {
		return resp, err
	}
	if err := resilienceCB.Allow(); err != nil {
		return nil, err
	}
	var lastErr error
	for attempt := 1; attempt <= resilienceMax; attempt++ {
		resilienceBucket.Take()
		resp, err := http.DefaultClient.Do(req)
		if resp != nil && err == nil && !shouldRetry(resp.StatusCode) {
			snapshotPersist(req, resp)
		}
		if resp != nil && resp.StatusCode == 429 {
			wait := parseRetryAfter(resp.Header.Get("Retry-After"))
			if wait > 0 && wait <= 5*time.Minute {
				resp.Body.Close()
				time.Sleep(wait)
				continue
			}
		}
		if err == nil && !shouldRetry(resp.StatusCode) {
			resilienceCB.onSuccess()
			return resp, nil
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		lastErr = err
		resilienceCB.onFailure()
		if attempt == resilienceMax {
			break
		}
		sleep := backoff(attempt, resilienceBaseMs)
		time.Sleep(sleep)
	}
	if lastErr == nil {
		lastErr = errors.New("github: request failed after retries")
	}
	return nil, lastErr
}

func shouldRetry(status int) bool {
	return status == 429 || (status >= 500 && status <= 599)
}

// parseRetryAfter accepts the HTTP Retry-After header: delta-seconds or HTTP-date.
func parseRetryAfter(value string) time.Duration {
	if value == "" {
		return 0
	}
	if secs, err := strconv.Atoi(value); err == nil {
		return time.Duration(secs) * time.Second
	}
	if when, err := http.ParseTime(value); err == nil {
		diff := time.Until(when)
		if diff > 0 {
			return diff
		}
	}
	return 0
}

func backoff(attempt, baseMs int) time.Duration {
	exp := math.Pow(2, float64(attempt-1))
	jitter := rand.Float64() * 0.5
	return time.Duration(float64(baseMs)*exp*(1+jitter)) * time.Millisecond
}

func envFloat(key string, fallback float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	if parsed, err := strconv.ParseFloat(value, 64); err == nil {
		return parsed
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return fallback
}
