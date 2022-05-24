package api

import (
	"gin-memos/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	var userRegister service.UserService
	err := c.ShouldBind(&userRegister) // 绑定json（前端传来的json数据绑定到userRegister中）
	// c.JSON向前端返回json数据
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "json数据解析失败！",
			"error":  err.Error(),
		})
		return
	}
	res := userRegister.Register()
	c.JSON(http.StatusOK, res)
}

// Login 用户登陆
func Login(c *gin.Context) {
	var userRegister service.UserService
	// 将request的body中的数据，自动按照json格式解析到结构体
	err := c.ShouldBind(&userRegister) // 绑定json（前端传来的json数据绑定到userRegister中）
	// c.JSON向前端返回json数据
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ // gin.H封装了生成json数据的工具
			"status": -1,
			"msg":    "json数据解析失败！",
			"error":  err.Error(),
		})
		return
	}
	res := userRegister.Login()
	c.JSON(http.StatusOK, res)
}
