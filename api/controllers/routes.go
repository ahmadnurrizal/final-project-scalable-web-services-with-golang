package controllers

import (
	"github.com/ahmadnurrizal/final-project-scalable-web-services-with-golang/api/middlewares"
	docs "github.com/ahmadnurrizal/final-project-scalable-web-services-with-golang/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title          MyGram
// @version        1.0
// @description    Final Project Scalable Web Services With Golang
// @termsOfService http://swagger.io/terms/

// @contact.name  Ahmad Nur Rizal
// @contact.url   https://lynk.id/ahmadnurrizal
// @contact.email ahmadnur.rizal45@gmail.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Token

// @host     localhost:8080
// @BasePath /api/v1

// @schemes http
func (s *Server) initializeRoutes() {
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := s.Router.Group("/api/v1")
	{
		// Login Route
		v1.POST("/login", s.Login)
		v1.POST("/users", s.Register)

		//Photos routes
		v1.GET("/photos", s.GetPhotos)
		v1.GET("/photos/:id", s.GetPhoto)
		v1.POST("/photos", middlewares.TokenAuthMiddleware(), s.CreatePhoto)
		v1.PUT("/photos/:id", middlewares.TokenAuthMiddleware(), s.UpdatePhoto)
		v1.DELETE("/photos/:id", middlewares.TokenAuthMiddleware(), s.DeletePhoto)

		//Comment routes
		v1.GET("/comments", s.GetComments)
		v1.GET("/comments/:id", s.GetComment)
		v1.POST("/comments/:id", middlewares.TokenAuthMiddleware(), s.CreateComment)
		v1.PUT("/comments/:id", middlewares.TokenAuthMiddleware(), s.UpdateComment)
		v1.DELETE("/comments/:id", middlewares.TokenAuthMiddleware(), s.DeleteComment)

		//SocialMedia routes
		v1.GET("/social-media-all", s.GetSocialMediaAll)
		v1.GET("/social-media/:id", s.GetSocialMedia)
		v1.POST("/social-media", middlewares.TokenAuthMiddleware(), s.CreateSocialMedia)
		v1.PUT("/social-media/:id", middlewares.TokenAuthMiddleware(), s.UpdateSocialMedia)
		v1.DELETE("/social-media/:id", middlewares.TokenAuthMiddleware(), s.DeleteSocialMedia)
	}
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
