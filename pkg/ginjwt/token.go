package ginjwt

import (
	"errors"
	"time"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

// Token provides a Json-Web-Token authentication implementation.
type Token struct {
	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	Key []byte
}

// NewToken 新建token
func NewToken(key string) *Token {
	return &Token{
		SigningAlgorithm: "HS256",
		Key:              []byte(key),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Make method that clients can use to get a jwt token.
func (T *Token) Make(expireAt time.Time, payloads map[string]interface{}) (token string, err error) {

	jToken := jwt.New(jwt.GetSigningMethod(T.SigningAlgorithm))
	claims := jToken.Claims.(jwt.MapClaims)

	claims["exp"] = expireAt.Unix()
	claims["orig_iat"] = time.Now().Unix()
	for k, v := range payloads {
		claims[k] = v
	}

	jToken.Claims = claims

	return jToken.SignedString(T.Key)
}

// Parse 解析token
func (T *Token) Parse(token string) (payloads map[string]interface{}, err error) {
	jToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(T.SigningAlgorithm) != token.Method {
			return nil, errors.New("invalid signing algorithm")
		}

		return T.Key, nil
	})

	if err != nil {
		return
	}

	d, ok := jToken.Claims.(jwt.MapClaims)
	if !ok {
		return payloads, errors.New("invalid token")
	}

	return map[string]interface{}(d), nil
}

// Check 解析token并验证token是否失效
func (T *Token) Check(token string) (payloads map[string]interface{}, err error) {
	payloads, err = T.Parse(token)
	if err != nil {
		return
	}

	exp := int64(payloads["exp"].(float64))

	expireAt := time.Unix(int64(exp), 0)
	if expireAt.Before(time.Now()) {
		return payloads, errors.New("expired")
	}

	return
}
