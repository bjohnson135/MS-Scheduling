// Faraday proxies all requests to Staffjoy
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"v2.staffjoy.com/faraday/services"

	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/errorpages"
	"v2.staffjoy.com/healthcheck"
	"v2.staffjoy.com/middlewares"
)

type contextKey int // Used for gorilla/context

const (
	// ServiceName is how this app is identified in logs and error handlers
	ServiceName          string     = "faraday"
	userID               contextKey = iota // Used for gorilla/context
	userSudo             contextKey = iota // Used for gorilla/context
	requestAuthenticated contextKey = iota // Used for gorilla/context
	requestedService     contextKey = iota // Used for gorilla/context
	requestedPathPrefix  contextKey = iota // matched prefix from ServiceMiddleware
)

var (
	logger       *logrus.Entry
	config       environments.Config
	signingToken = os.Getenv("SIGNING_SECRET")
	bannedUsers  = map[string]string{ // Use a map for constant time lookups. Value doesn't matter
		// Hypothetically these should be universally unique, so we don't have to limit by env
		"d7b9dbed-9719-4856-5f19-23da2d0e3dec": "hidden",
	}
)

// Setup environment, logger, etc
func init() {
	// Set the ENV environment variable to control dev/stage/prod behavior
	var err error
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine configuration")
	}
	logger = config.GetLogger(ServiceName)
}

// Listen for incoming requests, then validate, sanitize, and route them.
func main() {
	logger.Infof("Initialized environment %s", config.Name)

	if signingToken == "" {
		logger.Fatal("SIGNING_SECRET is required; refusing to boot. Set it in .env (see .env.schema).")
	}

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":80"
	}

	r := NewRouter(config, logger)
	s := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Infof("Faraday listening on %s", addr)
	logger.Panicf("%v", s.ListenAndServe())
}

// NewRouter returns a router composed of internal and external parts
func NewRouter(config environments.Config, logger *logrus.Entry) http.Handler {
	// Create a new router. We use Gorilla instead of stdlib because it handles
	// memory clean up for the 'context' package correctly
	externalRouter := mux.NewRouter()
	internalRouter := mux.NewRouter().PathPrefix("/").Subrouter().StrictSlash(true)

	// Faraday's /health is an aggregator: it probes every backend and
	// returns 200 only if all are healthy. `make doctor` hits this from
	// outside the docker network. Backends still expose their own
	// /health endpoints (handled by healthcheck.Handler) on the
	// internal network.
	externalRouter.HandleFunc(healthcheck.HEALTHPATH, healthzAggregator(services.StaffjoyServices, config.Name))
	externalRouter.HandleFunc(MobileConfigPath, MobileConfigHandler)

	sentryPublicDSN, err := environments.GetPublicSentryDSN(config.GetSentryDSN())
	if err != nil {
		logger.Fatalf("Cannot get sentry info - %s", err)
	}

	//traceMW, err := NewTraceMiddleware(logger, config)
	//if err != nil {
	//	logger.Fatalf("Unable to load trace middleware - %v", err)
	//}

	// only apply security to the internal routes
	externalRouter.PathPrefix("/").Handler(negroni.New(
		middlewares.NewRecovery(ServiceName, config, sentryPublicDSN),
		NewSecurityMiddleware(config),
		NewServiceMiddleware(config, services.StaffjoyServices),
		//traceMW,
		NewRobotstxtMiddleware(config),
		negroni.Wrap(internalRouter),
	))
	internalRouter.PathPrefix("/").HandlerFunc(proxyHandler)

	return externalRouter
}

// HTTP function that handles proxying after all of the middlewares
func proxyHandler(res http.ResponseWriter, req *http.Request) {
	service := req.Context().Value(requestedService).(services.Service)
	// No security on backend right now :-(
	destination := "http://" + service.BackendDomain + req.URL.RequestURI()
	logger.Debugf("Proxying to %s", destination)
	b, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		panic(fmt.Sprintf("Could not read request body - %s", err))
	}

	internalReq, err := http.NewRequest(req.Method, destination, bytes.NewReader(b))
	if err != nil {
		panic(fmt.Sprintf("Unable to create request - %s", err))
	}

	auth.SetInternalHeaders(req, internalReq.Header)

	currentUserUUID, err := auth.GetCurrentUserUUIDFromHeader(internalReq.Header)
	if err == nil {
		// authenticated request
		if _, isBanned := bannedUsers[currentUserUUID]; isBanned {
			logger.Warningf("Banned user accessing service - user %s", currentUserUUID)
			errorpages.Forbidden(res)
			return
		}
	}

	// Right here - check response Authorization and see if it's ok
	// with the requested service

	// Check perimeter authorization
	switch internalReq.Header.Get(auth.AuthorizationHeader) {
	case auth.AuthorizationAnonymousWeb:
		if service.Security != services.Public {
			// Send to /login on the same host (path-based routing — see
			// ADR-0004). The ServiceMiddleware may have stripped the
			// matched prefix; we have to put it back so the user lands
			// where they tried to go after they log in.
			prefix, _ := req.Context().Value(requestedPathPrefix).(string)
			returnTo := req.URL.EscapedPath()
			if service.StripPrefix && prefix != "" && prefix != "/" {
				returnTo = prefix + returnTo
				if returnTo == prefix+"/" && !strings.HasSuffix(req.URL.Path, "/") {
					// Original was "/app" not "/app/"; honor that.
					returnTo = prefix
				}
			}
			if req.URL.RawQuery != "" {
				returnTo += "?" + req.URL.RawQuery
			}
			http.Redirect(res, req, "/login/?return_to="+returnTo, http.StatusTemporaryRedirect)
			return
		}
	case auth.AuthorizationAuthenticatedUser:
		if service.Security == services.Admin {
			errorpages.Forbidden(res)
			return
		}
	case auth.AuthorizationSupportUser:
		// no restrictions
	default:
		logger.Panicf("unknown authorization header")
	}

	client := http.Client{
		// RETURN a redirect, do not FOLLOW it (which ends up causing relative redirect issues)
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	internalRes, err := client.Do(internalReq)
	if err != nil {
		logger.Warningf("Unable to query backend - %s", err)
		errorpages.GatewayTimeout(res)
		return
	}
	// Copy headers from service to user
	auth.ProxyHeaders(internalRes.Header, res.Header())

	if service.NoCacheHTML {
		if strings.Contains(strings.Join(res.Header()["Content-Type"], ""), "text/html") {
			// insert header to prevent caching
			res.Header().Set("Cache-Control", "no-cache")
		}
	}

	res.WriteHeader(internalRes.StatusCode)
	io.Copy(res, internalRes.Body)
	internalRes.Body.Close()

}
