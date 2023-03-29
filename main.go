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
	sqliteRepo, err := repo.NewSQLiteRepo(c.DBName)
	if err != nil {
		log.Fatal(err)
	}
	sqliteRepo.AutoMigrate()
	proxy := handler.NewProxyHandler(c.OpenAIKey, sqliteRepo.User)
	http.HandleFunc("/", proxy.Proxy)
	log.Fatal(http.ListenAndServe(":"+c.GinPort, nil))
}
