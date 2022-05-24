package util

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWT_sercet = []byte("qingbo1011.top") // 声明签名信息

// Claims 自定义有效载荷
type Claims struct {
	Uid                uint   `json:"uid"`
	UserName           string `json:"user_name"`
	jwt.StandardClaims        // StandardClaims结构体实现了Claims接口(Valid()函数)
}

// GenerateToken 签发token（调用jwt-go库生成token）
func GenerateToken(uid uint, userName string) (string, error) {
	notTime := time.Now()
	expireTime := notTime.Add(24 * time.Hour)
	claims := Claims{
		Uid:      uid,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			NotBefore: notTime.Unix(),    // 签名生效时间
			ExpiresAt: expireTime.Unix(), // 签名过期时间
			Issuer:    "qingbo1011.top",  // 签名颁发者
		},
	}
	// 指定编码算法为jwt.SigningMethodHS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 返回一个token结构体指针
	tokenString, err := token.SignedString(JWT_sercet)
	return tokenString, err
}

// ParserToken token解码
func ParserToken(tokenString string) (*Claims, error) {
	// 输入用户token字符串,自定义的Claims结构体对象,以及自定义函数来解析token字符串为jwt的Token结构体指针
	//Keyfunc是匿名函数类型: type Keyfunc func(*Token) (interface{}, error)
	//func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWT_sercet, nil
	})
	if err != nil {
		return nil, err
	}
	// 将token中的claims信息解析出来,并断言成用户自定义的有效载荷结构
	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token不可用")
}
