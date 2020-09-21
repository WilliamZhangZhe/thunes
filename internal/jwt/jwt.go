package jwtauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thunes/internal/model"
	"github.com/thunes/pkg/ginjwt"
)

var (
	JWTValueKey = "thunes-user"
	_ginJWT     = ginjwt.NewGinWebToken(".tutormeetplus.com", "token", "thunes", time.Hour*24)
)

// JWTCheck 检查请求中jwt-token是否有效
func JWTCheck(ctx *gin.Context) {
	token := _ginJWT.Get(ctx)
	if len(token) == 0 {
		UnAuthorized(ctx)
		ctx.Abort()
		return
	}

	if _, err := _ginJWT.Token.Check(token); err != nil {
		UnAuthorized(ctx)
		ctx.Abort()
		return
	}

	ctx.Next()
}

// JWTSet 设置返回新的jwt-token
func JWTSet(ctx *gin.Context) bool {
	var (
		err error
	)
	user := ctx.Value(JWTValueKey)
	if user == nil {
		fmt.Println("JWTSet: expect gin.Context user info ,but not found")
		return false
	}

	u, ok := user.(model.User)
	if !ok {
		fmt.Println("JWTSet: expect gin.Context user info of model.User,wrong type")
		return false
	}

	data := []byte{}
	payloads := func() (ret map[string]interface{}) {
		data, err = json.Marshal(&u)
		if err != nil {
			return
		}

		ret = map[string]interface{}{}
		err = json.Unmarshal(data, &ret)

		return ret
	}()

	if err != nil {
		fmt.Println("JWTSet: expect gin.Context user info of model.User,wrong type")
		return false
	}

	token, err := _ginJWT.Make(
		time.Now().Add(_ginJWT.Lifetime),
		payloads,
	)
	if err != nil {
		return false
	}

	_ginJWT.Set(ctx, token)
	return true
}

// JWTInvalid 设置JWT为无效
func JWTInvalid(ctx *gin.Context) {
	_ginJWT.Expire(ctx)
}

// UnAuthorized 授权失败或未授权结果通知
func UnAuthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, nil)
}
