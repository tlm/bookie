package api

import (
	"net/http"
)

type Server struct {
	httpServer http.Server
	Mux        *http.ServeMux
	Port       string
}

const (
	DEFAULT_PORT string = "3348"
)

func NewServer() *Server {
	serv := &Server{
		Port: DEFAULT_PORT,
	}
	serv.Mux = http.NewServeMux()
	serv.httpServer = http.Server{
		Addr:    ":" + DEFAULT_PORT,
		Handler: serv.Mux,
	}
	return serv
}
