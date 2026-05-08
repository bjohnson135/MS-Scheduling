package services

// Configuration for back-end services.
//
// Faraday uses path-based routing on a single host (default: localhost:8080).
// Path prefixes here are matched against the incoming request URL.Path; the
// service to which the request is proxied is determined by the longest
// matching prefix.
//
// See ADR-0004 for the rationale (path-based routing replaces the
// subdomain model that depended on `*.staffjoy-v2.local` DNS).

const (
	// Public security means a user may be logged out or logged in
	Public = iota
	// Authenticated security means a user must be logged in
	Authenticated = iota
	// Admin security means a user must be both logged in and have sudo flag
	Admin = iota
)

// ServiceDirectory maps a path prefix to a backend Service.
type ServiceDirectory map[string]Service

// Service describes a backend that Faraday proxies requests to.
type Service struct {
	// Security is the access policy applied before the proxy is hit.
	Security int

	// RestrictDev=true suppresses the service in staging and production.
	RestrictDev bool

	// BackendDomain is the in-network hostname of the backend (matches the
	// service name in docker-compose.yml).
	BackendDomain string

	// NoCacheHTML, when true, instructs the browser not to cache HTML
	// responses for this service.
	NoCacheHTML bool

	// StripPrefix, when true, strips the matched path prefix before
	// proxying. Used for backends that aren't aware of their mount point
	// (e.g. /api/account/v1/... -> the backend sees /v1/...).
	StripPrefix bool
}

// StaffjoyServices is the live route table.
//
// Path prefixes are exact-prefix matches and are evaluated in longest-prefix-
// first order at runtime (see PathToService in service_mw.go).
//
// Order in this map does not matter; the runtime sort handles it. The
// "/" entry is the catch-all and goes to the marketing site.
var StaffjoyServices = ServiceDirectory{
	"/api/account": {
		Security:      Public, // gateway translates REST to gRPC; authn happens further in
		BackendDomain: "accountapi-service",
		StripPrefix:   true,
	},
	"/api/company": {
		Security:      Public,
		BackendDomain: "companyapi-service",
		StripPrefix:   true,
	},
	"/app": {
		Security:      Authenticated,
		BackendDomain: "app-service",
		NoCacheHTML:   true,
	},
	"/myaccount": {
		Security:      Authenticated,
		BackendDomain: "myaccount-service",
		NoCacheHTML:   true,
	},
	"/whoami": {
		Security:      Public,
		BackendDomain: "whoami-service",
	},
	"/ical": {
		Security:      Public,
		BackendDomain: "ical-service",
	},
	"/superpowers": {
		Security:      Authenticated,
		RestrictDev:   true,
		BackendDomain: "superpowers-service",
	},
	"/": {
		// www (marketing + login + signup + logout) is the catch-all.
		Security:      Public,
		BackendDomain: "www-service",
	},
}
