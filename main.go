package main

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/casbin/casbin"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	key = []byte("qwertyuiop")
)

type ApiClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

type Auth struct {
	Token    string `json:"token"`
	ExpireIn string `json:"expire_at"`
}

type Merchant struct {
	Id   int
	Name string
}

var client *redis.Client
var e *casbin.Enforcer

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@/api_base?charset=utf8")

	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	e = casbin.NewEnforcer("./rbac_model.conf", "./policy.csv")
}

func main() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	// o := orm.NewOrm()
	//var t []*Merchant
	//r.GET("/", func(c *gin.Context) {
	//	s, err := client.Get("merchants").Result()
	//	if err == redis.Nil {
	//		_, err := o.Raw("SELECT * from merchants").QueryRows(&t)
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//		var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//		str,_ := json.Marshal(&t)
	//		s = string(str)
	//		client.Set("merchants", s, 0)
	//	} else if err != nil  {
	//		panic(err)
	//	}
	//	c.JSON(http.StatusOK, jsoniter.Unmarshal([]byte(s), &t))
	//})
	needAuth := r.Group("/api")
	needAuth.Use(AuthRequired())
	{
		needAuth.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"name": "mark",
			})
		})
	}
	r.GET("/login/:id", func(c *gin.Context) {
		id := c.Param("id")
		auth := GetToken(id)
		c.JSON(http.StatusOK, auth)
	})
	log.Fatal(r.Run(":8888"))
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		uid, isValid := CheckToken(token)
		if !isValid {
			c.AbortWithError(403, errors.New("token验证失败"))
		}
		sub := uid
		obj := "product"
		act := "read"
		if res := e.Enforce(sub, obj, act); !res {
			c.AbortWithError(403, errors.New("没有权限"))
		}
	}
}

func GetToken(Userid string) Auth {
	claims := &ApiClaims{
		Userid,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix()),
			ExpiresAt: int64(time.Now().Unix() + (60 * 60 * 24)),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		logs.Error(err)
		return Auth{}
	}
	return Auth{Token: ss, ExpireIn: strconv.Itoa(60 * 60 * 24)}
}

// 校验token是否有效
func CheckToken(tokenStr string) (string, bool) {
	token, err := jwt.ParseWithClaims(tokenStr, &ApiClaims{}, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		fmt.Println("parase with claims failed.", err)
		return "", false
	}
	if claims, ok := token.Claims.(*ApiClaims); ok && token.Valid {
		return claims.UserId, true
	} else {
		return "", false
	}
}
