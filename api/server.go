package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/spruce-cq/simplebank/db/sqlc"
)

// Server servers HTTP requests for banking service.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer create a new HTTP server and set up routers
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

// Start start the HTTP Server and listen on given address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
