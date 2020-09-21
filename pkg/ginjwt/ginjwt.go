package ginjwt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GinWebToken gin json-web-token
type GinWebToken struct {
	*Token

	Lifetime time.Duration
	domain   string
	key      string
}

// NewGinWebToken new token
func NewGinWebToken(domain, key, secret string, lifetime time.Duration) *GinWebToken {
	if lifetime == 0 {
		lifetime = time.Hour * 24
	}

	return &GinWebToken{
		Token:    NewToken(secret),
		Lifetime: lifetime,
		key:      key,
		domain:   domain,
	}
}

// Set 设置token
func (T *GinWebToken) Set(ctx *gin.Context, token string) {
	if ctx == nil {
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     T.key,
		Value:    token,
		HttpOnly: false,
		Path:     "/",
		MaxAge:   int(T.Lifetime.Seconds()),
		//Expires:  time.Now().Add(time.Hour * 3),
		Domain: T.domain,
	})
}

// Get 获取token
func (T *GinWebToken) Get(ctx *gin.Context) (token string) {
	token, _ = ctx.Cookie(T.key)

	return
}

// Expire 终止token
func (T *GinWebToken) Expire(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     T.key,
		Value:    "delete",
		HttpOnly: false,
		Expires:  time.Unix(0, 0),
		Path:     "/",
		Domain:   T.domain,
	})

	ctx.Next()
}
