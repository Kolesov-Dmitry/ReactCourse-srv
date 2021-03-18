package net

import (
	"net/http"

	"robochat.org/chat-srv/src/chat"
)

// Handler wraps storage to have access from handlers
type Handler struct {
	st chat.Storage

	handlerFunc func(st chat.Storage, w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP implementation of http.Handler interface
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.handlerFunc(h.st, w, r); err != nil {
		switch e := err.(type) {
		case Error:
			http.Error(w, e.Error(), e.Status())

		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
