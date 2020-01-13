package main

import (
	"fmt"
	"net/http"

	"common/fatality"
)

func handle_POST_SMS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		fatality.LogMsg("unsupported method", fmt.Errorf(r.Method))
		return
	}

	// try to publish message to MQ
	if err := publish(r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fatality.LogMsg("failed to publish message", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func registerHandlers(router *http.ServeMux) *http.ServeMux {
	router.Handle("/sms", recoveryMiddleware(http.HandlerFunc(handle_POST_SMS)))
	return router
}
