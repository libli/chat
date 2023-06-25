package main

import (
	"chat/config"
	"chat/handler"
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/carlmjohnson/requests"
)

type balance struct {
	TotalUsage float64 `json:"total_usage"`
}

func main() {
	log.Printf("chatapi (%s)", runtime.Version())

	// 读取配置文件
	c, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	for _, v := range c.AllKeys() {
		totalUsage := 0.0
		t0, _ := time.Parse("2006-01-02", "2023-01-01")
		// log.Println(k, *v)
		for i := 0; i < 12; i++ {
			t1 := t0.AddDate(0, i, 0)
			t2 := t0.AddDate(0, i+1, 0)
			if t1.After(time.Now()) {
				break
			}

			var b balance
			err := requests.
				URL(fmt.Sprintf(`https://%s/dashboard/billing/usage?start_date=%s&end_date=%s`, handler.OpenAIURL, t1.Format("2006-01-02"), t2.Format("2006-01-02"))).
				Header("Authorization", "Bearer "+v).
				ToJSON(&b).
				Fetch(ctx)
			if err != nil {
				log.Println(v, err)
			} else {
				totalUsage += float64(b.TotalUsage)
			}
		}
		log.Println(v, totalUsage)
	}
}
