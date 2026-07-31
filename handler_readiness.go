package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// this is a standard function signature that we need to use if we want to define an HTTP handler according to go std library
	healthmsg := "Server is ready"
	type heatlhCheckStruct struct {
		HStatus string `json:"health"`
	}
	//Note: I made a mistake here, previously i was writing the var name inside the struct as hStatus and it was throwing the error saying its unexported.... so write it in Capitalise format coz you have to export it
	responseWithJSON(w, 200, heatlhCheckStruct{
		HStatus: healthmsg,
	}) // now hook this up handler with a router
}
