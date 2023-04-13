package api

import (
	"fmt"
	"log"
	"os"

	"github.com/ahmadnurrizal/final-project-scalable-web-services-with-golang/api/controllers"
	"github.com/ahmadnurrizal/final-project-scalable-web-services-with-golang/api/seed"
	_ "github.com/ahmadnurrizal/final-project-scalable-web-services-with-golang/api/seed"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGPORT"), os.Getenv("PGHOST"), os.Getenv("PGDATABASE"))

	// This is for testing, when done, do well to comment
	seed.Load(server.DB)

	apiPort := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("PORT"))
	fmt.Printf("Listening to port %s", apiPort)

	server.Run(apiPort)

}
