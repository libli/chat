package repo

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySQLRepo MySQL数据库
type MySQLRepo struct {
	db   *gorm.DB
	User *UserRepo
}

// NewMySQLRepo 创建MySQL数据库
func NewMySQLRepo(dsn string) (*MySQLRepo, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &MySQLRepo{
		db:   db,
		User: newUserRepo(db),
	}, nil
}

// AutoMigrate 创建用户表
func (r *MySQLRepo) AutoMigrate() error {
	return r.db.AutoMigrate(&User{})
}
