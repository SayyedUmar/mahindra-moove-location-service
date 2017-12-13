package web

import (
	"context"
	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/identity"
	"net/http"
	"os"
)

func TokenAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authWithToken(w, req, fn)
	}
}

func SessionAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authWithSession(w, req, fn)
	}
}
func authWithToken(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	uid := req.Header.Get("uid")
	accessToken := req.Header.Get("access-token")
	client := req.Header.Get("client")
	if uid == "" || accessToken == "" || client == "" {
		ErrorWithMessage("Invalid Credentials").Respond(w, 401)
		return
	}
	ident := identity.FetchIdentityByUID(db.CurrentDB(), uid)
	if !ident.IsValid(client, accessToken) {
		ErrorWithMessage("Invalid Credentials").Respond(w, 401)
		return
	}
	ctx := context.WithValue(req.Context(), "identity", ident)
	newReq := req.WithContext(ctx)
	next(w, newReq)
}
func authWithSession(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	cookie, err := req.Cookie("_moove_session")
	if err != nil {
		ErrorWithMessage("Invalid cookie").Respond(w, 401)
	}
	railsCookie := cookie.Value
	sessionInfo := DecodeRailsSession(railsCookie, getRailsKeyBase())
	userId, err := ExtractUserId(sessionInfo)
	if err != nil {
		ErrorWithMessage("Invalid cookie").Respond(w, 401)
	}
	ident := identity.FetchIdentityByID(db.CurrentDB(), userId)
	ctx := context.WithValue(req.Context(), "identity", ident)
	newReq := req.WithContext(ctx)
	next(w, newReq)
}
func Auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Try session Auth
		cookie, err := req.Cookie("_moove_session")
		if cookie != nil && err == nil {
			authWithSession(w, req, fn)
		}

		// Try Token Auth
		uid := req.Header.Get("uid")
		accessToken := req.Header.Get("access-token")
		client := req.Header.Get("client")
		if uid != "" && accessToken != "" && client != "" {
			authWithToken(w, req, fn)
		}
		ErrorWithMessage("Invalid Credentials").Respond(w, 401)
	}
}

func getRailsKeyBase() string {
	keyBase := os.Getenv("RAILS_KEY_BASE")
	if keyBase != "" {
		return keyBase
	}
	return "f1d186616befd0912ed643cdc621377baa17368402970cfca9eaaf75f93286da121c22f1576ac5399a0d4c9ab3026849ebb67cd617437d73835c136e1c40a946"
}
