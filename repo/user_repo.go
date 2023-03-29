package repo

import (
	"time"

	"gorm.io/gorm"
)

// UserRepo 用户表DAO
type UserRepo struct {
	db *gorm.DB
}

// User 用户表
type User struct {
	ID         uint      `gorm:"primaryKey"`
	Username   string    `gorm:"size:50;not null;uniqueIndex:uk_username;comment:用户唯一标识号"`
	Token      string    `gorm:"size:50;not null;uniqueIndex:uk_token;comment:用户登录凭证"`
	Count      int       `gorm:"not null;default:0;comment:用户使用GPT次数"`
	Status     uint8     `gorm:"not null;default:1;comment:用户状态, 1: 正常, 2: 禁用"`
	CreateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdateTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:更新时间"`
}

func newUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetByToken 根据token获取用户
func (u *UserRepo) GetByToken(token string) *User {
	var user User
	u.db.Where("token = ?", token).First(&user)
	return &user
}

// UpdateCount 更新用户使用次数
func (u *UserRepo) UpdateCount(user *User) {
	user.Count++
	user.UpdateTime = time.Now()
	u.db.Save(&user)
}
