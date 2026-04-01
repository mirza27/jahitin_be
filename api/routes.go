package api

import "github.com/gin-gonic/gin"

func (s *Server) SetupRoutes() {

	router := gin.Default()

	router.POST("/register", s.registerUser)

	s.router = router
}
