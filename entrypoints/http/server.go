package http

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/boreq/plum/logging"
	"github.com/rs/cors"
)

type Server struct {
	handler http.Handler
	log     logging.Logger
}

func NewServer(handler http.Handler) *Server {
	return &Server{
		handler: handler,
		log:     logging.New("entrypoints/http.Server"),
	}
}

func (s *Server) Serve(address string) error {
	handler := cors.AllowAll().Handler(s.handler)
	handler = gziphandler.GzipHandler(handler)

	s.log.Info("starting listening", "address", address)
	return http.ListenAndServe(address, handler)
}
