package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"chat/repo"
)

const openAIURL = "api.openai.com"

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

	director := func(req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = openAIURL
		req.Host = openAIURL
		req.Header.Set("Authorization", "Bearer "+p.OpenAIKey)
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
	log.Printf("[*] receive the destination website response header: %s\n", w.Header())
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
