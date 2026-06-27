package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"column:username;type:varchar(50);not null;uniqueIndex" json:"username"`
	Password  string         `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Nickname  string         `gorm:"column:nickname;type:varchar(50);default:''" json:"nickname"`
	Email     string         `gorm:"column:email;type:varchar(100);default:'';uniqueIndex" json:"email"`
	Phone     string         `gorm:"column:phone;type:varchar(20);default:'';uniqueIndex" json:"phone"`
	Avatar    string         `gorm:"column:avatar;type:varchar(255);default:''" json:"avatar"`
	Gender    uint8          `gorm:"column:gender;type:tinyint unsigned;default:0" json:"gender"`
	Status    uint8          `gorm:"column:status;type:tinyint unsigned;default:1" json:"status"`
	IsDeleted uint8          `gorm:"column:is_deleted;type:tinyint unsigned;default:0" json:"is_deleted"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;default:current_timestamp;autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}