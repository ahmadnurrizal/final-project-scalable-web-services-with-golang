package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Photo struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Caption   string    `gorm:"size:255;not null;" json:"caption"`
	PhotoURL  string    `gorm:"size:255;not null;" json:"photo_url"`
	User      User      `json:"user"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Photo) Prepare() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Caption = html.EscapeString(strings.TrimSpace(p.Caption))
	p.User = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Photo) Validate() map[string]string {

	var err error

	var errorMessages = make(map[string]string)

	if p.Title == "" {
		err = errors.New("Required Title")
		errorMessages["Required_title"] = err.Error()

	}
	if p.Caption == "" {
		err = errors.New("Required Caption")
		errorMessages["Required_caption"] = err.Error()

	}
	if p.UserID < 1 {
		err = errors.New("Required User")
		errorMessages["Required_user"] = err.Error()
	}
	if p.PhotoURL == "" {
		err = errors.New("Required Photo URL")
		errorMessages["Required_photo"] = err.Error()
	}
	return errorMessages
}

func (p *Photo) SavePhoto(db *gorm.DB) (*Photo, error) {
	var err error
	err = db.Debug().Model(&Photo{}).Create(&p).Error
	if err != nil {
		return &Photo{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Photo{}, err
		}
	}
	return p, nil
}

func (p *Photo) FindAllPhotos(db *gorm.DB) (*[]Photo, error) {
	var err error
	photos := []Photo{}
	err = db.Debug().Model(&Photo{}).Limit(100).Order("created_at desc").Find(&photos).Error
	if err != nil {
		return &[]Photo{}, err
	}
	if len(photos) > 0 {
		for i, _ := range photos {
			err := db.Debug().Model(&User{}).Where("id = ?", photos[i].UserID).Take(&photos[i].User).Error
			if err != nil {
				return &[]Photo{}, err
			}
		}
	}
	return &photos, nil
}

func (p *Photo) FindPhotoByID(db *gorm.DB, pid uint64) (*Photo, error) {
	var err error
	err = db.Debug().Model(&Photo{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Photo{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Photo{}, err
		}
	}
	return p, nil
}

func (p *Photo) UpdateAPhoto(db *gorm.DB) (*Photo, error) {

	var err error

	err = db.Debug().Model(&Photo{}).Where("id = ?", p.ID).Updates(Photo{Title: p.Title, Caption: p.Caption, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Photo{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Photo{}, err
		}
	}
	return p, nil
}

func (p *Photo) DeleteAPhoto(db *gorm.DB) (int64, error) {

	db = db.Debug().Model(&Photo{}).Where("id = ?", p.ID).Take(&Photo{}).Delete(&Photo{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *Photo) FindUserPhotos(db *gorm.DB, uid uint32) (*[]Photo, error) {

	var err error
	photos := []Photo{}
	err = db.Debug().Model(&Photo{}).Where("user_id = ?", uid).Limit(100).Order("created_at desc").Find(&photos).Error
	if err != nil {
		return &[]Photo{}, err
	}
	if len(photos) > 0 {
		for i, _ := range photos {
			err := db.Debug().Model(&User{}).Where("id = ?", photos[i].UserID).Take(&photos[i].User).Error
			if err != nil {
				return &[]Photo{}, err
			}
		}
	}
	return &photos, nil
}

// When a user is deleted, we also delete the photo that the user had
func (c *Photo) DeleteUserPhotos(db *gorm.DB, uid uint32) (int64, error) {
	photos := []Photo{}
	db = db.Debug().Model(&Photo{}).Where("user_id = ?", uid).Find(&photos).Delete(&photos)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
