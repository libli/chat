package repo

import (
	"log"
	"time"

	"chat/model"

	"github.com/glebarez/sqlite"
	// "github.com/ysmood/gop"
	"gorm.io/gorm"
)

// SQLiteRepo SQLite数据库
type SQLiteRepo struct {
	db   *gorm.DB
	User *UserRepo
}

// NewSQLiteRepo 创建SQLite数据库
func NewSQLiteRepo(dbname string, users []model.User) (*SQLiteRepo, error) {
	log.Println("Init:", dbname)
	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	repo := SQLiteRepo{
		db:   db,
		User: newUserRepo(db),
	}
	repo.AutoMigrate()
	for _, it := range users {
		it.CreateTime = time.Now()
		u := repo.User.GetByName(it.Username)
		// gop.P(u)
		if u.Username == it.Username && u.Token == it.Token {
			continue
		} else {
			if u.Username == it.Username {
				u.Token = it.Token
				repo.User.db.Save(u)
			} else {
				repo.User.Add(&it)
				if repo.User.GetByName(it.Username).Username == it.Username {
					log.Println("Inited:", it.Username, it.Token)
				}
			}
		}
	}
	log.Println("Inited:", "ALL")
	return &repo, nil
}

// AutoMigrate 创建用户表
func (r *SQLiteRepo) AutoMigrate() error {
	return r.db.AutoMigrate(&model.User{})
}
