package api

import (
	"net/http"
	"time"

	db "github.com/TonyGLL/image-processing-service/internal/db/sqlc"
	util "github.com/TonyGLL/image-processing-service/internal/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	config util.Config
}

func NewServer(store db.Store, config util.Config) *http.Server {
	NewServer := &Server{
		store:  store,
		config: config,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         config.ServerAddress,
		Handler:      NewServer.SetupRoutes(config.Version),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
