package main

import (
	"log"
	"net/http"
	"runtime/debug"
)

func recoverHTTP(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Println("http recovered from ", r, "\nStack trace:", string(debug.Stack()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func recoveryMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer recoverHTTP(w)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
