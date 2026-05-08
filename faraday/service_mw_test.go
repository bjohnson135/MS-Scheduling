package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/faraday/services"
)

func TestPathToService_LongestPrefixWins(t *testing.T) {
	dir := services.ServiceDirectory{
		"/api/account": {Security: services.Public, BackendDomain: "accountapi-service", StripPrefix: true},
		"/api/company": {Security: services.Public, BackendDomain: "companyapi-service", StripPrefix: true},
		"/app":         {Security: services.Authenticated, BackendDomain: "app-service"},
		"/whoami":      {Security: services.Public, BackendDomain: "whoami-service"},
		"/":            {Security: services.Public, BackendDomain: "www-service"},
	}

	tests := []struct {
		path        string
		wantPrefix  string
		wantBackend string
	}{
		{"/api/account/v1/users/abc", "/api/account", "accountapi-service"},
		{"/api/company/v1/shifts", "/api/company", "companyapi-service"},
		{"/app", "/app", "app-service"},
		{"/app/", "/app", "app-service"},
		{"/app/team/123", "/app", "app-service"},
		{"/whoami", "/whoami", "whoami-service"},
		{"/whoami/", "/whoami", "whoami-service"},
		{"/", "/", "www-service"},
		{"/login/", "/", "www-service"},
		{"/anything-else", "/", "www-service"},
	}

	for _, tc := range tests {
		gotPrefix, gotSvc, err := PathToService(tc.path, dir)
		if err != nil {
			t.Fatalf("PathToService(%q) unexpected err: %v", tc.path, err)
		}
		if gotPrefix != tc.wantPrefix {
			t.Errorf("PathToService(%q) prefix = %q, want %q", tc.path, gotPrefix, tc.wantPrefix)
		}
		if gotSvc.BackendDomain != tc.wantBackend {
			t.Errorf("PathToService(%q) backend = %q, want %q", tc.path, gotSvc.BackendDomain, tc.wantBackend)
		}
	}
}

func TestServiceMiddleware_StripsPrefix(t *testing.T) {
	dir := services.ServiceDirectory{
		"/api/account": {Security: services.Public, BackendDomain: "accountapi-service", StripPrefix: true},
		"/":            {Security: services.Public, BackendDomain: "www-service"},
	}
	mw := NewServiceMiddleware(environments.Config{Name: "development"}, dir)

	req := httptest.NewRequest(http.MethodGet, "/api/account/v1/users/abc?x=1", nil)
	rec := httptest.NewRecorder()

	var nextReq *http.Request
	mw.ServeHTTP(rec, req, func(_ http.ResponseWriter, r *http.Request) {
		nextReq = r
	})

	if nextReq == nil {
		t.Fatalf("middleware did not call next; status=%d body=%q", rec.Code, rec.Body.String())
	}
	if got, want := nextReq.URL.Path, "/v1/users/abc"; got != want {
		t.Errorf("after StripPrefix, URL.Path = %q, want %q", got, want)
	}
	svc, ok := nextReq.Context().Value(requestedService).(services.Service)
	if !ok {
		t.Fatal("requestedService missing from context")
	}
	if svc.BackendDomain != "accountapi-service" {
		t.Errorf("backend = %q, want accountapi-service", svc.BackendDomain)
	}
}

func TestServiceMiddleware_DoesNotStripWhenDisabled(t *testing.T) {
	dir := services.ServiceDirectory{
		"/app": {Security: services.Public, BackendDomain: "app-service"},
		"/":    {Security: services.Public, BackendDomain: "www-service"},
	}
	mw := NewServiceMiddleware(environments.Config{Name: "development"}, dir)

	req := httptest.NewRequest(http.MethodGet, "/app/team/123", nil)
	rec := httptest.NewRecorder()

	var nextReq *http.Request
	mw.ServeHTTP(rec, req, func(_ http.ResponseWriter, r *http.Request) {
		nextReq = r
	})

	if nextReq == nil {
		t.Fatalf("middleware did not call next; status=%d body=%q", rec.Code, rec.Body.String())
	}
	if got, want := nextReq.URL.Path, "/app/team/123"; got != want {
		t.Errorf("URL.Path = %q, want %q (no strip)", got, want)
	}
}

func TestServiceMiddleware_RestrictDevHidesInProd(t *testing.T) {
	dir := services.ServiceDirectory{
		"/superpowers": {Security: services.Public, BackendDomain: "superpowers-service", RestrictDev: true},
		"/":            {Security: services.Public, BackendDomain: "www-service"},
	}

	for _, env := range []string{"production", "staging"} {
		mw := NewServiceMiddleware(environments.Config{Name: env}, dir)
		req := httptest.NewRequest(http.MethodGet, "/superpowers/", nil)
		rec := httptest.NewRecorder()
		nextCalled := false
		mw.ServeHTTP(rec, req, func(http.ResponseWriter, *http.Request) { nextCalled = true })

		if nextCalled {
			t.Errorf("env=%s: dev-only service should not have called next", env)
		}
		if rec.Code != http.StatusNotFound {
			t.Errorf("env=%s: expected 404, got %d", env, rec.Code)
		}
	}

	mw := NewServiceMiddleware(environments.Config{Name: "development"}, dir)
	req := httptest.NewRequest(http.MethodGet, "/superpowers/", nil)
	rec := httptest.NewRecorder()
	nextCalled := false
	mw.ServeHTTP(rec, req, func(http.ResponseWriter, *http.Request) { nextCalled = true })
	if !nextCalled {
		t.Error("development env should reach dev-only service")
	}
}

func TestServiceMiddleware_PrefixSetInContext(t *testing.T) {
	dir := services.ServiceDirectory{
		"/whoami": {Security: services.Public, BackendDomain: "whoami-service"},
		"/":       {Security: services.Public, BackendDomain: "www-service"},
	}
	mw := NewServiceMiddleware(environments.Config{Name: "development"}, dir)

	req := httptest.NewRequest(http.MethodGet, "/whoami", nil).
		WithContext(context.Background())
	rec := httptest.NewRecorder()
	var nextReq *http.Request
	mw.ServeHTTP(rec, req, func(_ http.ResponseWriter, r *http.Request) {
		nextReq = r
	})
	if nextReq == nil {
		t.Fatal("middleware did not call next")
	}
	prefix, ok := nextReq.Context().Value(requestedPathPrefix).(string)
	if !ok {
		t.Fatal("requestedPathPrefix missing from context")
	}
	if prefix != "/whoami" {
		t.Errorf("prefix = %q, want %q", prefix, "/whoami")
	}
}
