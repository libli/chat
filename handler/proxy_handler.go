package handler

import (
	"io"
	"log"
	"net/http"
	"strings"

	"chat/repo"
)

const openAIURL = "https://api.openai.com"

type ProxyHandler struct {
	OpenAIKey string
	user      *repo.UserRepo
}

// NewProxyHandler creates a new ProxyHandler.
func NewProxyHandler(key string, user *repo.UserRepo) *ProxyHandler {
	return &ProxyHandler{
		OpenAIKey: key,
		user:      user,
	}
}

// Proxy is the handler for the openai proxy.
func (p *ProxyHandler) Proxy(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")

	log.Println("auth: ", auth)
	if auth == "" {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if !p.checkToken(token) {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}

	uri := openAIURL + r.RequestURI
	destReq, err := http.NewRequest(r.Method, uri, r.Body)
	if err != nil {
		log.Println("new request error: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}
	p.copyHeaders(r.Header, &destReq.Header)

	var transport http.Transport
	destResp, err := transport.RoundTrip(destReq)
	if err != nil {
		log.Println("do request error: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}
	defer destResp.Body.Close()
	body, err := io.ReadAll(destResp.Body)
	if err != nil {
		log.Println("read response error: ", err)
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}
	log.Println("success")
	respHeader := w.Header()
	p.copyHeaders(destResp.Header, &respHeader)
	w.Write(body)
}

func (p *ProxyHandler) checkToken(token string) bool {
	user := p.user.GetByToken(token)
	if user == nil || user.Token != token {
		return false
	}
	log.Println("user name: ", user.Username)
	p.user.UpdateCount(user)
	return true
}

// copyHeaders copies the headers from source to dest.
func (p *ProxyHandler) copyHeaders(source http.Header, dest *http.Header) {
	// Header结构体为 map[string][]string，因为头部参数中可以使用相同的key
	for key, values := range source {
		// 把Authorization替换成OpenAI的key
		if key == "Authorization" {
			dest.Add("Authorization", "Bearer "+p.OpenAIKey)
			continue
		}
		for _, value := range values {
			dest.Add(key, value)
		}
	}
}
