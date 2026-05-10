package middleware

import (
	"log"
	"net"
	"net/http"
	"strings"

	"vblog-core/service"
)

// PVMiddleware returns an HTTP middleware that records page views for non-API routes.
func PVMiddleware(pvSvc *service.PageViewService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if pvSvc != nil {
				ip := r.RemoteAddr
				if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
					ip = strings.Split(forwarded, ",")[0]
				} else {
					if host, _, err := net.SplitHostPort(ip); err == nil {
						ip = host
					}
				}
				ip = strings.TrimSpace(ip)
				go func() {
					if err := pvSvc.Record(ip, r.URL.Path, r.UserAgent()); err != nil {
						log.Printf("PV record failed: %v", err)
					}
				}()
			}
			next.ServeHTTP(w, r)
		})
	}
}
