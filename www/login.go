package main

import (
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/grpc/metadata"

	"v2.staffjoy.com/account"
	"v2.staffjoy.com/auth"
	"v2.staffjoy.com/helpers"

	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
)

type loginPage struct {
	Title         string // Used in <title>
	CSSId         string // e.g. 'careers'
	Version       string // e.g. master-1, for cachebusting
	CsrfField     template.HTML
	Denied        bool
	Description   string
	TemplateName  string
	PreviousEmail string
	ReturnTo      string
}

// isValidReturnPath returns true if returnTo is a safe in-app path that we
// can redirect a freshly-logged-in user to. With path-based routing on a
// single host (ADR-0004), the only valid return_to is a path beginning
// with "/" that matches a registered Faraday prefix. Absolute URLs (with
// scheme + host) are rejected so an attacker can't redirect off-host.
func isValidReturnPath(p string) bool {
	if p == "" || !strings.HasPrefix(p, "/") || strings.HasPrefix(p, "//") {
		return false
	}
	// Reject anything that looks like it has a host component.
	if strings.Contains(p, "://") {
		return false
	}
	allowed := []string{"/app", "/myaccount", "/whoami", "/ical", "/superpowers"}
	for _, a := range allowed {
		if p == a || strings.HasPrefix(p, a+"/") {
			return true
		}
	}
	return false
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
	// if logged in - go away
	if req.Header.Get(auth.AuthorizationHeader) == auth.AuthorizationAuthenticatedUser {
		http.Redirect(res, req, "/myaccount/", http.StatusFound)
		return
	}

	// for GET
	returnTo := req.URL.Query().Get("return_to")

	p := loginPage{
		Title:        "Staffjoy Log in",
		CSSId:        "login",
		CsrfField:    csrf.TemplateField(req),
		Version:      config.GetDeployVersion(),
		TemplateName: "login.tmpl",
		Description:  "Log in to Staffjoy to start scheduling your workers. All you’ll need is your email and password.",
		ReturnTo:     returnTo,
	}

	if req.Method == http.MethodPost {
		email := req.FormValue("email")
		password := req.FormValue("password")
		// POST and GET are in the same handler - reset
		returnTo = req.FormValue("return_to")
		// rememberMe=True means that the session is set for a month instead of a day
		rememberMe := len(req.FormValue("remember-me")) > 0

		md := metadata.New(map[string]string{auth.AuthorizationMetadata: auth.AuthorizationWWWService})
		ctx, cancel := context.WithCancel(metadata.NewOutgoingContext(context.Background(), md))
		defer cancel()

		accountClient, close, err := account.NewClient()
		if err != nil {
			panic(err)
		}
		defer close()

		user, err := accountClient.VerifyPassword(ctx, &account.VerifyPasswordRequest{Email: email, Password: password})
		if err == nil {
			auth.LoginUser(user.Uuid, user.Support, rememberMe, res)
			go helpers.TrackEvent(user.Uuid, "login")
			go accountClient.SyncUser(ctx, &account.SyncUserRequest{Uuid: user.Uuid})

			logger.WithFields(logrus.Fields{"user_uuid": user.Uuid}).Info("Logging in user")

			// Path-based routing: return_to must already be an in-app path
			// like "/app/team/123". Default landing is /app/. Any unsafe
			// value falls back to /myaccount/.
			if !isValidReturnPath(returnTo) {
				if returnTo == "" {
					returnTo = "/app/"
				} else {
					returnTo = "/myaccount/"
				}
			}

			http.Redirect(res, req, returnTo, http.StatusFound)
			return
		}

		logger.WithFields(logrus.Fields{"email": email}).Infof("Login attempt denied")
		p.Denied = true
		p.PreviousEmail = email
	}
	err := tmpl.ExecuteTemplate(res, p.TemplateName, p)
	if err != nil {
		panic(err)
	}

}
