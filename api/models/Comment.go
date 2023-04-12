package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Message   string    `gorm:"size:255;not null" json:"message"`
	User      User      `json:"user"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	PhotoID   uint64    `gorm:"not null" json:"photo_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Comment) Prepare() {
	p.Message = html.EscapeString(strings.TrimSpace(p.Message))
	p.User = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Comment) Validate() map[string]string {

	var err error

	var errorMessages = make(map[string]string)

	if p.Message == "" {
		err = errors.New("Required Message")
		errorMessages["Required_message"] = err.Error()

	}
	if p.UserID < 1 {
		err = errors.New("Required User")
		errorMessages["Required_user"] = err.Error()
	}
	return errorMessages
}

func (p *Comment) SaveComment(db *gorm.DB) (*Comment, error) {
	var err error
	err = db.Debug().Model(&Comment{}).Create(&p).Error
	if err != nil {
		return &Comment{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Comment{}, err
		}
	}

	return p, nil
}

func (p *Comment) FindAllComments(db *gorm.DB) (*[]Comment, error) {
	var err error
	comments := []Comment{}
	err = db.Debug().Model(&Comment{}).Limit(100).Order("created_at desc").Find(&comments).Error
	if err != nil {
		return &[]Comment{}, err
	}
	if len(comments) > 0 {
		for i, _ := range comments {
			err := db.Debug().Model(&User{}).Where("id = ?", comments[i].UserID).Take(&comments[i].User).Error
			if err != nil {
				return &[]Comment{}, err
			}
		}
	}
	return &comments, nil
}

func (p *Comment) FindCommentByID(db *gorm.DB, pid uint64) (*Comment, error) {
	var err error
	err = db.Debug().Model(&Comment{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Comment{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return p, nil
}

func (p *Comment) UpdateAComment(db *gorm.DB) (*Comment, error) {

	var err error

	err = db.Debug().Model(&Comment{}).Where("id = ?", p.ID).Updates(Comment{Message: p.Message, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Comment{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return p, nil
}

func (p *Comment) DeleteAComment(db *gorm.DB) (int64, error) {

	db = db.Debug().Model(&Comment{}).Where("id = ?", p.ID).Take(&Comment{}).Delete(&Comment{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (p *Comment) FindUserComments(db *gorm.DB, uid uint32) (*[]Comment, error) {

	var err error
	comments := []Comment{}
	err = db.Debug().Model(&Comment{}).Where("user_id = ?", uid).Limit(100).Order("created_at desc").Find(&comments).Error
	if err != nil {
		return &[]Comment{}, err
	}
	if len(comments) > 0 {
		for i, _ := range comments {
			err := db.Debug().Model(&User{}).Where("id = ?", comments[i].UserID).Take(&comments[i].User).Error
			if err != nil {
				return &[]Comment{}, err
			}
		}
	}
	return &comments, nil
}

// When a user is deleted, we also delete the comment that the user had
func (c *Comment) DeleteUserComments(db *gorm.DB, uid uint32) (int64, error) {
	comments := []Comment{}
	db = db.Debug().Model(&Comment{}).Where("user_id = ?", uid).Find(&comments).Delete(&comments)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
