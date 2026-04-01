package api

import (
	"jahitin_be/config"
	db "jahitin_be/database/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.Config
	store  db.Store
	router *gin.Engine
}

func NewServer(c config.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: c,
		store:  store,
	}

	server.SetupRoutes()

	return server, nil
}

func (s *Server) Start() error {
	address := "0.0.0.0:" + strconv.Itoa(s.config.AppPort)

	return s.router.Run(address)
}
