package api

import (
	db "github.com/denizumutdereli/golang-pg-sqlc-cli/db/sqlc"
	"github.com/gin-gonic/gin"
)

//backend connection
type Server struct {
	store  *db.Store
	router *gin.Engine
}

//setup
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	//routes
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)

	server.router = router
	return server
}

//starting server - specific address for
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
