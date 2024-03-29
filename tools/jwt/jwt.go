package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

const EXPIRETIME = 60 * 60 * 24

var (
	key []byte
)

type ApiClaims struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type Auth struct {
	UserId   string `json:"user_id"`
	Role     string `json:"role"`
	Token    string `json:"token"`
	ExpireIn string `json:"expire_at"`
}

func GetToken(userId string, role string) Auth {
	claims := &ApiClaims{
		userId,
		role,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix()),
			ExpiresAt: int64(time.Now().Unix() + EXPIRETIME),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		return Auth{}
	}
	return Auth{UserId: userId, Role: role, Token: ss, ExpireIn: strconv.Itoa(EXPIRETIME)}
}

// 校验token是否有效
func CheckToken(tokenStr string) (string, string, bool) {
	token, err := jwt.ParseWithClaims(tokenStr, &ApiClaims{}, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		fmt.Println("parase with claims failed.", err)
		return "", "", false
	}
	if claims, ok := token.Claims.(*ApiClaims); ok && token.Valid {
		return claims.UserId, claims.Role, true
	} else {
		return "", "", false
	}
}

func init() {
	str := viper.GetString("jwtKey")
	key = []byte(str)
	fmt.Println("初始化jwt, key:", str)
}
