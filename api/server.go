package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peienxie/url-shortener/storage"
)

// Server serves HTTP requests
type Server struct {
	store  storage.URLStore
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup its routing
func NewServer(store storage.URLStore) *Server {
	server := &Server{
		store:  store,
		router: gin.Default(),
	}

	// each IP can only create 30 shorten urls per hour
	limiter := NewRateLimiter(time.Hour, 30)

	// initilize routing
	server.router.POST("/shorturls", limiter, server.createShortURL)
	server.router.GET("/:short_url", server.redirectShortURL)

	return server
}

// Serve runs the http server on the provided address
func (s *Server) Serve(addr string) error {
	return s.router.Run(addr)
}
