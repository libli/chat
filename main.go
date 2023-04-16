package main

import (
	"log"
	"net/http"

	"chat/config"
	"chat/handler"
	"chat/repo"
)

func main() {
	// 读取配置文件
	c, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}

	if c.Driver == "mysql" {
		// 创建 MySQL 数据库
		mysqlRepo, err := repo.NewMySQLRepo(c.MySQLDSN())
		if err != nil {
			log.Fatal(err)
		}
		mysqlRepo.AutoMigrate()
		proxy := handler.NewProxyHandler(c.OpenAIKey, mysqlRepo.User)
		http.HandleFunc("/", proxy.Proxy)
		log.Fatal(http.ListenAndServe(":"+c.GinPort, nil))
		return
	} else {
		// 创建 SQLite 数据库
		sqliteRepo, err := repo.NewSQLiteRepo(c.DBName)
		if err != nil {
			log.Fatal(err)
		}
		sqliteRepo.AutoMigrate()
		proxy := handler.NewProxyHandler(c.OpenAIKey, sqliteRepo.User)
		http.HandleFunc("/", proxy.Proxy)
		log.Fatal(http.ListenAndServe(":"+c.GinPort, nil))
		return
	}
}
