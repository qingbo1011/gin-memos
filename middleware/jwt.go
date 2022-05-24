package middleware

import (
	"gin-memos/pkg/util"
	"gin-memos/serializer"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// JWTAuth 定义一个JWTAuth的中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 通过http header中的token解析来认证
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusNotFound, serializer.Response{
				Status: http.StatusNotFound,
				Msg:    "请求未携带token，无权访问！",
			})
			c.Abort() // Abort 函数在被调用的函数中阻止后续中间件的执行(这里都没有携带token，后续就不用执行了)
			return
		}
		// 解析token中包含的相关信息（有效载荷）
		claims, err := util.ParserToken(token)
		if err != nil {
			c.JSON(http.StatusForbidden, serializer.Response{ // token无权限（是假的）
				Status: http.StatusForbidden,
				Msg:    "token解析失败！",
				Error:  err.Error(),
			})
			c.Abort()
			return
		}
		if time.Now().Unix() > claims.ExpiresAt { // token过期了
			c.JSON(http.StatusUnauthorized, serializer.Response{
				Status: http.StatusForbidden,
				Msg:    "token已过期！",
			})
			c.Abort()
			return
		}
		//// 将解析后的有效载荷claims重新写入gin.Context引用对象中（gin的上下文）
		//c.Set("claims",claims)
	}
}
