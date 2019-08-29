package router

import (
	. "curry/controllers"
	. "curry/tools"
	_ "curry/tools/config"
	_ "curry/tools/db/mysql"
	redis2 "curry/tools/db/redis"
	"curry/tools/jwt"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
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
	// public api
	r.GET("/test-redis", func(c *gin.Context) {
		redis := redis2.NewClient()
		key := "name"
		err := redis.Set(key, "mark", 0).Err()
		go CheckError(err)
		str, _ := redis.Get(key).Result()
		c.JSON(http.StatusOK, str)
	})

	// token api
	api := r.Group("/api")
	api.Use(AuthRequired())
	{
		api.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"name": "mark",
			})
		})
	}

	// Login & Register
	r.POST("/login", Login)
}

// auth middleware
func AuthRequired() gin.HandlerFunc {
	var (
		userId  string
		isValid bool
		err     error
	)
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		fmt.Println("token:", token)
		userId, _, isValid = jwt.CheckToken(token)
		if !isValid {
			err = c.AbortWithError(403, errors.New("token验证失败"))
			go CheckError(err)
			return
		}

		if res := e.Enforce(userId, c.Request.URL.Path, c.Request.Method); !res {
			err = c.AbortWithError(403, errors.New("没有权限"))
			go CheckError(err)
			return
		}
		c.Next()
	}
}
