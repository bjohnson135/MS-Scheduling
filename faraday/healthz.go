package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"v2.staffjoy.com/faraday/services"
	"v2.staffjoy.com/healthcheck"
)

// healthzAggregator returns an http.HandlerFunc that probes every backend's
// /health endpoint in parallel and writes a JSON roll-up. Status code is 200
// when every backend in scope is OK, 503 otherwise.
//
// Replaces the base healthcheck.Handler on faraday so `make doctor` (which
// hits faraday from outside the docker network) can verify the entire stack
// in one HTTP call.
func healthzAggregator(svcs services.ServiceDirectory, envName string) http.HandlerFunc {
	type result struct {
		Status string `json:"status"`
		Code   int    `json:"http_status"`
	}

	type summary struct {
		Status   string            `json:"status"`
		Services map[string]result `json:"services"`
	}

	probe := func(ctx context.Context, host string) result {
		c := &http.Client{Timeout: 800 * time.Millisecond}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+host+healthcheck.HEALTHPATH, nil)
		if err != nil {
			return result{Status: "error", Code: 0}
		}
		resp, err := c.Do(req)
		if err != nil {
			return result{Status: "unreachable", Code: 0}
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return result{Status: "ok", Code: resp.StatusCode}
		}
		return result{Status: "unhealthy", Code: resp.StatusCode}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 1500*time.Millisecond)
		defer cancel()

		out := summary{
			Status:   "ok",
			Services: make(map[string]result, len(svcs)),
		}

		var mu sync.Mutex
		var wg sync.WaitGroup

		for prefix, svc := range svcs {
			if prefix == "/" {
				// Catch-all is www; not interesting as a separate probe.
				continue
			}
			if svc.RestrictDev && envName != "development" && envName != "test" {
				continue
			}
			wg.Add(1)
			go func(prefix string, svc services.Service) {
				defer wg.Done()
				res := probe(ctx, svc.BackendDomain)
				mu.Lock()
				out.Services[prefix] = res
				// RestrictDev services are profile-gated in compose; an
				// unreachable one is "not in scope," not a degradation.
				if res.Status != "ok" && !svc.RestrictDev {
					out.Status = "degraded"
				}
				mu.Unlock()
			}(prefix, svc)
		}

		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		if out.Status == "ok" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		_ = json.NewEncoder(w).Encode(out)
	}
}
