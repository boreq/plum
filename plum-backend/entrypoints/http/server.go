package http

import (
	"net/http"

	"github.com/boreq/plum/plum-backend/logging"
	"github.com/klauspost/compress/gzhttp"
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

func (s *Server) wrap() http.Handler {
	handler := cors.AllowAll().Handler(s.handler)
	return gzhttp.GzipHandler(handler)
}

func (s *Server) Serve(address string) error {
	handler := s.wrap()

	s.log.Info("starting listening", "address", address)
	return http.ListenAndServe(address, handler)
}
