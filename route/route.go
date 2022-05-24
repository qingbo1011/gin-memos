package route

import (
	"gin-memos/api"
	"gin-memos/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()
	user := r.Group("/api/user")
	{
		user.POST("/register", api.UserRegister)
		user.POST("/login", api.Login)
	}
	task := r.Group("/api/task")
	task.Use(middleware.JWTAuth()) // 加载自定义JWTAuth()中间件，在整合task路由都生效
	{
		// POST是不幂等的， GET、PUT、DELETE等方法都是幂等的
		task.POST("/create", api.CreateTask)
		task.GET("/get/:id", api.GetTaskById)
		task.POST("/getall", api.GetAllTask)
		task.PUT("/update/:id", api.UpdateTask)
		task.DELETE("/delete/:id", api.DeleteTask)
	}
	return r
}
