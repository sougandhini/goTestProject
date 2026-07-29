package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// this is a standard function signature that we need to use if we want to define an HTTP handler according to go std library

	responseWithJSON(w, 200, struct{}{}) // now hook this up handler with a router
}
