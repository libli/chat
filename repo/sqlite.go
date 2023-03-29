package repo

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// SQLiteRepo SQLite数据库
type SQLiteRepo struct {
	db   *gorm.DB
	User *UserRepo
}

// NewSQLiteRepo 创建SQLite数据库
func NewSQLiteRepo(dbname string) (*SQLiteRepo, error) {
	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &SQLiteRepo{
		db:   db,
		User: newUserRepo(db),
	}, nil
}

// AutoMigrate 创建用户表
func (r *SQLiteRepo) AutoMigrate() error {
	return r.db.AutoMigrate(&User{})
}
