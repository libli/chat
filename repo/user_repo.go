package repo

import (
	"chat/model"
	"log"
	"time"

	"gorm.io/gorm"
)

// UserRepo 用户表DAO
type UserRepo struct {
	db *gorm.DB
}

func newUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Add 添加用户
func (u *UserRepo) Add(user *model.User) {
	log.Println(user)
	u.db.Create(&user)
}

// GetByToken 根据token获取用户
func (u *UserRepo) GetByToken(token string) *model.User {
	var user model.User
	u.db.Where("token = ?", token).First(&user)
	return &user
}

// GetByName 根据username获取用户
func (u *UserRepo) GetByName(name string) *model.User {
	var user model.User
	u.db.Where("username = ?", name).First(&user)
	return &user
}

// UpdateCount 更新用户使用次数
func (u *UserRepo) UpdateCount(user *model.User) {
	user.Count++
	user.UpdateTime = time.Now()
	u.db.Save(&user)
}
