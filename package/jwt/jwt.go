package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 本文件的文档参考博客：https://www.liwenzhou.com/posts/Go/json-web-token/

var InvaildToken = "invaild token"

// 加密字符串
var signKey = []byte("are you ok?")

// 自定义jwt结构体，需要存储用户id
type SelfClaims struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// 用于生成jwt的token
func GenerateJwtToken(user_id int64) (string, error) {
	claims := SelfClaims{
		user_id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)), // 过期时间，设定2小时
			Issuer:    "lx",                                              // 签发人
		},
	}

	// 指定签名方法进行签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signKey) // 获得完整的签名字符串
}

// 解析token
func ParseJwtToken(tokenStr string) (*SelfClaims, error) {
	// 自定义的claim需要使用这个方法，官方提供的则直接使用Parse
	sc := new(SelfClaims)
	token, err := jwt.ParseWithClaims(tokenStr, sc, func(tk *jwt.Token) (i interface{}, err error) {
		return signKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return sc, nil
	}

	return nil, errors.New(InvaildToken)
}
