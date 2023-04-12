package seed

import (
	"log"

	"github.com/ahmadnurrizal/final-project-scalable-web-services-with-golang/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Username: "admin",
		Email:    "admin@gmail.com",
		Password: "password",
		Age:      17,
	},
	models.User{
		Username: "udin",
		Email:    "udin@gmail.com",
		Password: "password",
		Age:      19,
	},
	models.User{
		Username: "rizal",
		Email:    "rizal@gmail.com",
		Password: "password",
		Age:      13,
	},
}

func Load(db *gorm.DB) {

	// make sure drop and migrate table in correct order
	// or can avoid error by remove foreign key constraint first
	// db.Model(&models.Comment{}).RemoveForeignKey("user_id", "users(id)")
	// db.Model(&models.Comment{}).RemoveForeignKey("photo_id", "photos(id)")
	err := db.Debug().DropTableIfExists(&models.SocialMedia{}, &models.Comment{}, &models.Photo{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Photo{}, &models.SocialMedia{}, &models.Comment{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Photo{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Comment{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.SocialMedia{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Photo{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Comment{}).AddForeignKey("photo_id", "photos(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}
