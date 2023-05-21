package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/phantranhieunhan/authen-author/db/redis"
	db "github.com/phantranhieunhan/authen-author/db/sqlc"
	"github.com/phantranhieunhan/authen-author/token"
	"github.com/phantranhieunhan/authen-author/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	session    redis.SessionStore
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store, session redis.SessionStore) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		session:    session,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authenMiddleware(server.tokenMaker, server.session), authorMiddleware(server.store))
	authRoutes.POST("/users/logout", server.logoutUser)
	authRoutes.GET("/accounts/:id", server.getAccount)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
