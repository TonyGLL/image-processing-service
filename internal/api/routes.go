package api

import (
	"net/http"

	"github.com/TonyGLL/image-processing-service/internal/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func (server *Server) SetupRoutes(version string) http.Handler {
	r := gin.Default()

	docs.SwaggerInfo.Title = "Swagger Image processing service Documentation"
	docs.SwaggerInfo.Description = "This is an Image processing service with Golang."
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = "localhost:3000"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	/* CORS */
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "OPTIONS", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
	}))

	v1 := r.Group("/api/v1")
	{
		// No need token validation
		v1.POST("/login", server.loginHandler)

		// Token Validation, everything down here need token validation
		v1.Use(middlewares.ValidateJWT)

		v1.POST("/register", server.createUserHandler)
		v1.POST("/images", server.uploadImageHandler)
		v1.GET("/images", server.listImagesHandler)
		v1.GET("/images/:id", server.getImageByIdHandler)
	}

	return r
}
