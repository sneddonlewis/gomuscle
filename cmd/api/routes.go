package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *ApiServer) Start() {
	router := s.routes(mux.NewRouter())

	addr := fmt.Sprintf(":%d", s.cfg.port)

	s.logger.Printf("Starting %s API server on %s", s.cfg.env, addr)
	http.ListenAndServe(addr, router)
}

func (s *ApiServer) routes(router *mux.Router) *mux.Router {
	router.HandleFunc("/v1/health", makeHttpHandlerFunc(s.health))

	router.HandleFunc("/account", makeHttpHandlerFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHttpHandlerFunc(s.handleGetAccountById)))
	router.HandleFunc("/transfer", makeHttpHandlerFunc(s.handleTransfer))

	return router
}
