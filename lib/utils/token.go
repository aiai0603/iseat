package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const SecretKey = "Lib_SecretKey"
const issuer = "lib_issuer"

func Md5V2(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

type jwtCustomClaims struct {
	jwt.StandardClaims
	ID   uint `json:"id"` // 用户的编号
	Admin bool `json:"admin"` // 是否为管理员
}

func CreateToken(id uint, isAdmin bool) (tokenString string, err error) {
	claims := &jwtCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Hour * 72).Unix()),
			Issuer:    issuer,
		},
		id,
		isAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(SecretKey))
	return
}

func ParseToken(tokenSrt string) (claims jwt.Claims, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenSrt, func(*jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	claims = token.Claims
	return
}
