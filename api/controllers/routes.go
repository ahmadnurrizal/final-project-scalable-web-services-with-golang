package controllers

import (
	"github.com/ahmadnurrizal/final-project-scalable-web-services-with-golang/api/middlewares"
)

func (s *Server) initializeRoutes() {
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
}
