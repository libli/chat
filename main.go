package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"chat/config"
	"chat/handler"
	"chat/repo"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func main() {
	log.Printf("chatapi (%s)", runtime.Version())

	// 读取配置文件
	c, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}
	sqliteRepo, err := repo.NewSQLiteRepo(c.DBName, c.InitUsers)
	if err != nil {
		log.Fatal(err)
	}
	proxy := handler.NewProxyHandler(c.RoundRobinKey, sqliteRepo.User)

	http.HandleFunc("/", proxy.Proxy)
	http.HandleFunc("/healthz", healthz)
	log.Println("Serveing", ":"+c.GinPort)
	log.Fatal(http.ListenAndServe(":"+c.GinPort, nil))
}
