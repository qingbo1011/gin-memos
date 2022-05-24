package api

import (
	"gin-memos/pkg/util"
	"gin-memos/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateTask 创建一条备忘录
func CreateTask(c *gin.Context) {
	var createTaskService service.CreateTaskService
	err := c.ShouldBind(&createTaskService)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "json解析失败",
			"error":  err.Error(),
		})
		return
	}
	// 先解析一下token，看看是哪个用户在进行操作
	claims, err := util.ParserToken(c.GetHeader("token"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "token解析失败",
			"error":  err.Error(),
		})
		return
	}
	res := createTaskService.CreateTask(claims.Uid)
	c.JSON(http.StatusOK, res)
}

// GetTaskById 根据id查询备忘录
func GetTaskById(c *gin.Context) {
	var getTaskByIdService service.GetTaskByIdService
	err := c.ShouldBindUri(&getTaskByIdService)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "url参数绑定失败！",
			"error": err.Error(),
		})
		return
	}
	res := getTaskByIdService.GetTaskById(getTaskByIdService.Id)
	c.JSON(http.StatusOK, res)
}

// GetAllTask 分页查询所有备忘录
func GetAllTask(c *gin.Context) {
	var getAllService service.GetAllService
	err := c.ShouldBind(&getAllService)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "json解析失败",
			"error":  err.Error(),
		})
		return
	}
	// 先解析一下token，看看是哪个用户在进行操作
	claims, err := util.ParserToken(c.GetHeader("token"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "token解析失败",
			"error":  err.Error(),
		})
		return
	}
	res := getAllService.GetAllTask(claims.Uid)
	c.JSON(http.StatusOK, res)
}

// UpdateTask 更新一条备忘录
func UpdateTask(c *gin.Context) {
	var updateTaskService service.UpdateTaskService
	err := c.ShouldBind(&updateTaskService)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "json解析失败",
			"error":  err.Error(),
		})
		return
	}
	res := updateTaskService.UpdateTask(c.Param("id"))
	c.JSON(http.StatusOK, res)
}

// DeleteTask 删除一条备忘录
func DeleteTask(c *gin.Context) {
	res := service.DeleteTaskById(c.Param("id"))
	c.JSON(http.StatusOK, res)
}
