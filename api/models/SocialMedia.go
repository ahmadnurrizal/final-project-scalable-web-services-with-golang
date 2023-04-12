package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type SocialMedia struct {
	ID             uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name           string    `gorm:"size:255;not null;unique" json:"name"`
	SocialMediaURL string    `gorm:"text;not null;" json:"socialMediaURL"`
	User           User      `json:"user"`
	UserID         uint32    `gorm:"not null" json:"user_id"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *SocialMedia) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.SocialMediaURL = html.EscapeString(strings.TrimSpace(p.SocialMediaURL))
	p.User = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *SocialMedia) Validate() map[string]string {

	var err error

	var errorMessages = make(map[string]string)

	if p.Name == "" {
		err = errors.New("Required Name")
		errorMessages["Required_title"] = err.Error()

	}
	if p.SocialMediaURL == "" {
		err = errors.New("Required Social Media URL")
		errorMessages["Required_social_media_URL"] = err.Error()
	}
	if p.UserID < 1 {
		err = errors.New("Required User")
		errorMessages["Required_user"] = err.Error()
	}
	return errorMessages
}

func (p *SocialMedia) SaveSocialMedia(db *gorm.DB) (*SocialMedia, error) {
	var err error
	err = db.Debug().Model(&SocialMedia{}).Create(&p).Error
	if err != nil {
		return &SocialMedia{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &SocialMedia{}, err
		}
	}
	return p, nil
}

func (p *SocialMedia) FindAllSocialMedia(db *gorm.DB) (*[]SocialMedia, error) {
	var err error
	socialMedias := []SocialMedia{}
	err = db.Debug().Model(&SocialMedia{}).Limit(100).Order("created_at desc").Find(&socialMedias).Error
	if err != nil {
		return &[]SocialMedia{}, err
	}
	if len(socialMedias) > 0 {
		for i, _ := range socialMedias {
			err := db.Debug().Model(&User{}).Where("id = ?", socialMedias[i].UserID).Take(&socialMedias[i].User).Error
			if err != nil {
				return &[]SocialMedia{}, err
			}
		}
	}
	return &socialMedias, nil
}

func (p *SocialMedia) FindSocialMediaByID(db *gorm.DB, pid uint64) (*SocialMedia, error) {
	var err error
	err = db.Debug().Model(&SocialMedia{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &SocialMedia{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &SocialMedia{}, err
		}
	}
	return p, nil
}

func (p *SocialMedia) UpdateASocialMedia(db *gorm.DB) (*SocialMedia, error) {

	var err error

	err = db.Debug().Model(&SocialMedia{}).Where("id = ?", p.ID).Updates(SocialMedia{Name: p.Name, SocialMediaURL: p.SocialMediaURL, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &SocialMedia{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &SocialMedia{}, err
		}
	}
	return p, nil
}

func (p *SocialMedia) DeleteASocialMedia(db *gorm.DB) (int64, error) {

	db = db.Debug().Model(&SocialMedia{}).Where("id = ?", p.ID).Take(&SocialMedia{}).Delete(&SocialMedia{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *SocialMedia) FindUserSocialMedias(db *gorm.DB, uid uint32) (*[]SocialMedia, error) {

	var err error
	socialMedias := []SocialMedia{}
	err = db.Debug().Model(&SocialMedia{}).Where("user_id = ?", uid).Limit(100).Order("created_at desc").Find(&socialMedias).Error
	if err != nil {
		return &[]SocialMedia{}, err
	}
	if len(socialMedias) > 0 {
		for i, _ := range socialMedias {
			err := db.Debug().Model(&User{}).Where("id = ?", socialMedias[i].UserID).Take(&socialMedias[i].User).Error
			if err != nil {
				return &[]SocialMedia{}, err
			}
		}
	}
	return &socialMedias, nil
}

// When a user is deleted, we also delete the socialMedia that the user had
func (c *SocialMedia) DeleteUserSocialMedias(db *gorm.DB, uid uint32) (int64, error) {
	socialMedias := []SocialMedia{}
	db = db.Debug().Model(&SocialMedia{}).Where("user_id = ?", uid).Find(&socialMedias).Delete(&socialMedias)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
