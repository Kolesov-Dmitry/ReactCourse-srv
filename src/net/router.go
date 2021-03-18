package net

import (
	"net/http"

	"github.com/gorilla/mux"
	"robochat.org/chat-srv/src/chat"
)

func initRouter(st chat.Storage) *mux.Router {
	router := mux.NewRouter()

	// CORS OPTION request handlers
	router.HandleFunc("/api/chats", optionsHandler).Methods("OPTIONS")
	router.HandleFunc("/api/messages", optionsHandler).Methods("OPTIONS")
	router.HandleFunc("/api/profile", optionsHandler).Methods("OPTIONS")

	// storage handlers
	router.Handle("/api/chats", Handler{st, chatsGetHandler}).Methods("GET")
	router.Handle("/api/chats", Handler{st, chatsPostHandler}).Methods("POST")
	router.Handle("/api/chats", Handler{st, chatsDeleteHandler}).Methods("DELETE")

	router.Handle("/api/messages", Handler{st, messagesGetHandler}).Methods("GET")
	router.Handle("/api/messages", Handler{st, messagesPostHandler}).Methods("POST")
	router.Handle("/api/messages", Handler{st, messagesDeleteHandler}).Methods("DELETE")

	router.Handle("/api/profile", Handler{st, profileGetHandler}).Methods("GET")
	router.Handle("/api/profile", Handler{st, profilePostHandler}).Methods("POST")

	// static files handler
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return router
}
