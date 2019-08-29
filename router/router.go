package router

import (
	. "curry/tools"
	_ "curry/tools/config"
	redis2 "curry/tools/db/redis"
	"curry/tools/jwt"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

var (
	r *gin.Engine
	e *casbin.Enforcer
)

func init() {
	e = casbin.NewEnforcer(viper.GetString("casbin.modelPath"), viper.GetString("casbin.policyPath"))
	logs.Debug("casbin获取所有对象:", e.GetAllObjects())
}

func NewRouter() *gin.Engine {
	r = gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	setupRouter()
	return r
}

func setupRouter() {
	r.GET("/test-redis", func(c *gin.Context) {
		redis := redis2.NewClient()
		str, _ := redis.Get("merchants").Result()
		c.JSON(http.StatusOK, str)
	})
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
		auth := jwt.GetToken(id, "admin")
		c.JSON(http.StatusOK, auth)
	})
}

func AuthRequired() gin.HandlerFunc {
	var err error
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		uid, isValid := jwt.CheckToken(token)
		if !isValid {
			err = c.AbortWithError(403, errors.New("token验证失败"))
			go CheckError(err)
		}
		sub := uid
		obj := "product"
		act := "read"
		if res := e.Enforce(sub, obj, act); !res {
			err = c.AbortWithError(403, errors.New("没有权限"))
			go CheckError(err)
		}
	}
}
