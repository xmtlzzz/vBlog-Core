package middleware

import (
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	restful "github.com/emicklei/go-restful/v3"
)

// rateLimiter is a fixed-window in-memory limiter keyed by string.
type rateLimiter struct {
	mu      sync.Mutex
	window  time.Duration
	limit   int
	counts  map[string]*rateBucket
	lastGC  time.Time
}

type rateBucket struct {
	count int
	reset time.Time
}

func newRateLimiter(window time.Duration, limit int) *rateLimiter {
	return &rateLimiter{window: window, limit: limit, counts: map[string]*rateBucket{}, lastGC: time.Now()}
}

// Allow reports whether one more request under key is permitted in the current window.
func (r *rateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	// Occasional GC to bound memory
	if now.Sub(r.lastGC) > time.Hour {
		for k, b := range r.counts {
			if now.After(b.reset) {
				delete(r.counts, k)
			}
		}
		r.lastGC = now
	}
	b, ok := r.counts[key]
	if !ok || now.After(b.reset) {
		r.counts[key] = &rateBucket{count: 1, reset: now.Add(r.window)}
		return true
	}
	if b.count >= r.limit {
		return false
	}
	b.count++
	return true
}

// RateLimitFilter returns a go-restful filter limiting requests per client IP.
// X-Forwarded-For is honored only behind a trusted proxy (TRUST_PROXY=1);
// direct exposure lets clients spoof it to rotate buckets.
func RateLimitFilter(window time.Duration, limit int) restful.FilterFunction {
	limiter := newRateLimiter(window, limit)
	trustProxy := os.Getenv("TRUST_PROXY") == "1"
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		ip := req.Request.RemoteAddr
		if trustProxy {
			if forwarded := req.Request.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = strings.TrimSpace(strings.Split(forwarded, ",")[0])
			}
		} else if host, _, err := net.SplitHostPort(ip); err == nil {
			ip = host // RemoteAddr is host:port; key on host so new connections share the bucket
		}
		if !limiter.Allow(ip) {
			resp.AddHeader("Retry-After", "60")
			resp.WriteHeaderAndEntity(http.StatusTooManyRequests, map[string]string{"error": "请求过于频繁，请稍后再试"})
			return
		}
		chain.ProcessFilter(req, resp)
	}
}
