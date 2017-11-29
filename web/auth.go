package web

import (
	"context"
	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/identity"
	"net/http"
)

func TokenAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
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
		fn(w, newReq)
	}
}
