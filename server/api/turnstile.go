package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const turnstileSiteverifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

// verifyTurnstile performs the canonical server-side siteverify check.
//
// It requires success === true, the expected action, and an approved frontend
// hostname. It fails closed: when Turnstile is not configured (empty secret or
// no hostnames), or the response token is missing/oversized, the request is
// refused so no data is ever written without the check.
func verifyTurnstile(r *http.Request, secret string, hostnames []string, token, expectedAction string) bool {
	if secret == "" || token == "" || len(token) > 2048 || len(hostnames) == 0 {
		return false
	}
	allowed := make(map[string]bool, len(hostnames))
	for _, h := range hostnames {
		allowed[strings.TrimSpace(h)] = true
	}

	form := url.Values{}
	form.Set("secret", secret)
	form.Set("response", token)
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		form.Set("remoteip", ip)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, turnstileSiteverifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}

	var result struct {
		Success  bool   `json:"success"`
		Action   string `json:"action"`
		Hostname string `json:"hostname"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}
	return result.Success && result.Action == expectedAction && allowed[result.Hostname]
}
