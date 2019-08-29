package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Title 登录接口
// @Param userName
// @Param password
// @Param rePassword
// @Param captcha
func Login(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"name": "mark",
	})
}
