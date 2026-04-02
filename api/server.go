package api

import (
	"jahitin_be/config"
	db "jahitin_be/database/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	config config.Config
	store  db.Store
	router *gin.Engine
}

func NewApiServer(c config.Config, store db.Store) (*ApiServer, error) {
	server := &ApiServer{
		config: c,
		store:  store,
	}

	server.SetupRoutes()

	return server, nil
}

func (s *ApiServer) Start() error {
	address := "0.0.0.0:" + strconv.Itoa(s.config.AppPort)

	return s.router.Run(address)
}
