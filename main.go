package main

import (
	"github.com/ahmadnurrizal/final-project-scalable-web-services-with-golang/api"
)

// @title          MyGram
// @version        1.0
// @description    Final Project Scalable Web Services With Golang

// @contact.name  Ahmad Nur Rizal
// @contact.url   https://lynk.id/ahmadnurrizal
// @contact.email ahmadnur.rizal45@gmail.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host     localhost:8080
// @BasePath /api/v1

// @schemes http
func main() {
	api.Run()
}
