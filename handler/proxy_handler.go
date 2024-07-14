package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"chat/repo"
)

const OpenAIURL = "api.openai.com"

type ProxyHandler struct {
	OpenAIKey func(*http.Request) string
	user      *repo.UserRepo
}

// NewProxyHandler creates a new ProxyHandler.
func NewProxyHandler(getKey func(*http.Request) string, user *repo.UserRepo) *ProxyHandler {
	return &ProxyHandler{
		OpenAIKey: getKey,
		user:      user,
	}
}

// Proxy is the handler for the openai proxy.
func (p *ProxyHandler) Proxy(w http.ResponseWriter, r *http.Request) {
	// CORS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusOK)
		return
	}

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
		req.URL.Host = OpenAIURL
		req.Host = OpenAIURL
		req.Header.Set("Authorization", "Bearer "+p.OpenAIKey(r))
		req.Header.Del("Sec-Ch-Ua")
		req.Header.Del("Sec-Ch-Ua-Mobile")
		req.Header.Del("Sec-Ch-Ua-Platform")
		req.Header.Del("Sec-Fetch-Dest")
		req.Header.Del("Sec-Fetch-Mode")
		req.Header.Del("Sec-Fetch-Site")
		req.Header.Set("User-Agent", "Darwin/23.4.0")
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
