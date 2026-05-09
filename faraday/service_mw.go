package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/faraday/services"
)

// ServiceMiddleware routes incoming requests to backend services based on the
// URL path prefix. See ADR-0004.
type ServiceMiddleware struct {
	Config       environments.Config
	Services     services.ServiceDirectory
	sortedRoutes []string // path prefixes sorted longest-first
}

// NewServiceMiddleware constructs a path-routing middleware. The route table
// is sorted by descending length once at construction so the request hot path
// stays cheap.
func NewServiceMiddleware(config environments.Config, svcs services.ServiceDirectory) *ServiceMiddleware {
	routes := make([]string, 0, len(svcs))
	for prefix := range svcs {
		routes = append(routes, prefix)
	}
	sort.Slice(routes, func(i, j int) bool {
		return len(routes[i]) > len(routes[j])
	})

	return &ServiceMiddleware{
		Config:       config,
		Services:     svcs,
		sortedRoutes: routes,
	}
}

func (svc *ServiceMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	prefix, service, ok := svc.match(req.URL.Path)
	if !ok {
		res.WriteHeader(http.StatusBadGateway)
		res.Write([]byte("No backend matches this path"))
		return
	}

	if service.RestrictDev && svc.Config.Name != "development" && svc.Config.Name != "test" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if service.StripPrefix && prefix != "/" {
		// Rewrite the request path so the backend doesn't have to know its mount point.
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
		if !strings.HasPrefix(req.URL.Path, "/") {
			req.URL.Path = "/" + req.URL.Path
		}
	}

	ctx := context.WithValue(req.Context(), requestedService, service)
	ctx = context.WithValue(ctx, requestedPathPrefix, prefix)
	req = req.WithContext(ctx)

	next(res, req)
}

// match returns the first registered prefix (in longest-first order) that the
// path begins with, the corresponding Service, and ok=true. If no prefix
// matches it returns ok=false.
func (svc *ServiceMiddleware) match(path string) (string, services.Service, bool) {
	for _, prefix := range svc.sortedRoutes {
		if prefix == "/" {
			// The root prefix is the catch-all and always matches last.
			service := svc.Services[prefix]
			return prefix, service, true
		}
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return prefix, svc.Services[prefix], true
		}
	}
	return "", services.Service{}, false
}

// PathToService is a convenience wrapper used in tests.
func PathToService(path string, dir services.ServiceDirectory) (string, services.Service, error) {
	mw := NewServiceMiddleware(environments.Config{}, dir)
	prefix, svc, ok := mw.match(path)
	if !ok {
		return "", services.Service{}, fmt.Errorf("no backend matches path %q", path)
	}
	return prefix, svc, nil
}
