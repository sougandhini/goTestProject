package main

import (
	"fmt"
	"net/http"

	"github.com/sougandhini/rssagg/internal/auth"
	"github.com/sougandhini/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// Note: this type does not match authentic http method handler, so we write another function

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	// we're returning a closure http to use it with chi router, the returned function is a new anonymous function with the same params as our http normal func
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			responseWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := cfg.DB.GetUserByAPI(r.Context(), apiKey) //context package which gives way to track something that's happening across mutilple go routine, and we can also cancel the context - which kills http request
		if err != nil {
			responseWithError(w, 400, fmt.Sprintf("Could not get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}
