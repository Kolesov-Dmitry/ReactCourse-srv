package net

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"robochat.org/chat-srv/src/chat"
)

// Server is a Catalog HTTP server
type Server struct {
	useTLS bool
	server *http.Server
	port   int
	st     chat.Storage
}

// NewServer returns new Server with TLS support
func NewServer(port int, st chat.Storage) *Server {
	http.Handle("/", initRouter(st))

	return &Server{
		useTLS: false,
		server: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		port: port,
		st:   st,
	}
}

// UseTLS ...
func (s *Server) UseTLS(certPath string, name string) error {
	cert, err := tls.LoadX509KeyPair(certPath+"/server.crt", certPath+"/server.key")
	if err != nil {
		return err
	}

	s.useTLS = true
	s.server.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   name,
	}

	return nil
}

// Start starts HTTP server
func (s *Server) Start() {
	go func() {
		var err error
		if s.useTLS {
			err = s.server.ListenAndServeTLS("", "")
		} else {
			err = s.server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

// Stop closes HTTP server
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Port returns number of the listened port
func (s *Server) Port() int {
	return s.port
}
