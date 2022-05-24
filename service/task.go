package service

import (
	"gin-memos/db/mysql"
	"gin-memos/model"
	"gin-memos/serializer"
	"net/http"
	"time"
)

type CreateTaskService struct {
	Tittle  string `json:"tittle" form:"tittle" binding:"required"`
	Content string `json:"content" form:"content"`
}

type GetTaskByIdService struct {
	Id string `uri:"id" binding:"required"`
}

type GetAllService struct {
	PageSize int    `json:"page_size" form:"page_size" default:"10"` // 每一页展示多少条结果(不加binding:"required"，默认是10)
	PageNum  int    `json:"page_num" form:"page_num" default:"1"`    // 第几页(不加binding:"required"，默认是第一页)
	KeyWord  string `json:"key_word" form:"key_word"`                // 支持可能的模糊查询
	Desc     bool   `json:"desc" form:"desc"`                        // 是否反向搜索（比如需要寻找较早数据时）
}

type UpdateTaskService struct {
	Tittle  string `json:"tittle" form:"tittle"`
	Content string `json:"content" form:"content"`
	Status  int    `json:"status"  form:"status"`
}

func (s *CreateTaskService) CreateTask(uid uint) serializer.Response {
	var user model.User
	mysql.DB.First(&user, uid) // 通过主键进行查询 (仅适用于主键是数字类型)（行内条件查询，类似where）
	// 原生SQL：SELECT * FROM user WHERE id = uid LIMIT 1;
	task := model.Task{
		User:       user,
		Uid:        user.ID,
		Tittle:     s.Tittle,
		Status:     0,
		Context:    s.Content,
		StartTime:  time.Now().Unix(),
		FinishTime: 0,
	}
	err := mysql.DB.Create(&task).Error // 向数据库中插入记录
	if err != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "创建备忘录失败！",
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: http.StatusOK,
		Msg:    user.UserName + "创建备忘录成功！",
	}
}

func (s *GetTaskByIdService) GetTaskById(tid string) serializer.Response {
	var task model.Task
	var user model.User
	err := mysql.DB.First(&task, tid).Error
	err = mysql.DB.First(&user, task.Uid).Error
	if err != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "查询task或user失败！",
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: http.StatusOK,
		Msg:    "查询成功！",
		Data:   serializer.BuildTask(task, user.UserName),
	}
}

func (s *GetAllService) GetAllTask(uid uint) serializer.Response {
	var tasks []model.Task
	var user model.User
	count := 0
	// 对PageSize和PageNum没有传入的情况的处理（因为没有加binding:"required"）
	if s.PageSize <= 0 {
		s.PageSize = 10
	}
	if s.PageNum <= 0 {
		s.PageNum = 1
	}
	if s.KeyWord == "" { // 不加关键词进行模糊查询
		if s.Desc {
			err := mysql.DB.Table("task").Order("created_at desc").Where("uid = ?", uid).Count(&count).Limit(s.PageSize).Offset(s.PageSize * (s.PageNum - 1)).Find(&tasks).Error
			if err != nil {
				return serializer.Response{
					Status: http.StatusInternalServerError,
					Msg:    "数据库查询失败！",
					Error:  err.Error(),
				}
			}
		} else {
			err := mysql.DB.Table("task").Where("uid = ?", uid).Count(&count).Limit(s.PageSize).Offset(s.PageSize * (s.PageNum - 1)).Find(&tasks).Error
			if err != nil {
				return serializer.Response{
					Status: http.StatusInternalServerError,
					Msg:    "数据库查询失败！",
					Error:  err.Error(),
				}
			}
		}
	} else { // 加关键词进行模糊查询
		if s.Desc {
			err := mysql.DB.Table("task").Order("created_at desc").Where("uid = ? and tittle like ?", uid, "%"+s.KeyWord+"%").Or("uid = ? and context like ?", uid, "%"+s.KeyWord+"%").Count(&count).Limit(s.PageSize).Offset(s.PageSize * (s.PageNum - 1)).Find(&tasks).Error
			if err != nil {
				return serializer.Response{
					Status: http.StatusInternalServerError,
					Msg:    "数据库查询失败！",
					Error:  err.Error(),
				}
			}
		} else {
			err := mysql.DB.Table("task").Where("uid = ?", uid).Where("tittle like ? or context like ?", "%"+s.KeyWord+"%", "%"+s.KeyWord+"%").
				Count(&count).Limit(s.PageSize).Offset(s.PageSize * (s.PageNum - 1)).Find(&tasks).Error
			if err != nil {
				return serializer.Response{
					Status: http.StatusInternalServerError,
					Msg:    "数据库查询失败！",
					Error:  err.Error(),
				}
			}
		}
	}
	err := mysql.DB.First(&user, uid).Error
	if err != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库查询失败！",
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: http.StatusOK,
		Msg:    "查询成功！",
		Data:   serializer.BuildListResponse(serializer.BuildTasks(tasks, user.UserName), uint(count), user.UserName),
	}
}

func (s *UpdateTaskService) UpdateTask(tid string) serializer.Response {
	var task model.Task
	var user model.User
	err := mysql.DB.First(&task, tid).Error // 通过主键进行查询 (仅适用于主键是数字类型)
	mysql.DB.First(&user, task.Uid)
	if err != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "指定的tid查询失败！",
			Error:  err.Error(),
		}
	}
	task.Context = s.Content
	task.Tittle = s.Tittle
	task.Status = s.Status
	err = mysql.DB.Save(&task).Error // Save 方法在执行 SQL 更新操作时将包含所有字段，即使这些字段没有被修改。
	if err != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "DB.Save(task)出现错误，更新失败！",
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: http.StatusOK,
		Msg:    "更新备忘录成功！",
		Data:   serializer.BuildTask(task, user.UserName),
	}
}

func DeleteTaskById(tid string) serializer.Response {
	var task model.Task
	var user model.User
	err := mysql.DB.First(&task, tid).Error
	if err != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "DB.First(&task, tid)出现错误，根据tid查询task失败！",
			Error:  err.Error(),
		}
	}
	err = mysql.DB.First(&user, task.Uid).Error
	if err != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "DB.First(&user, task.Uid)出现错误，根据uid查询user失败！",
			Error:  err.Error(),
		}
	}
	err = mysql.DB.Delete(&task).Error
	if err != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "DB.Delete(&task)出现错误，删除task失败！",
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: http.StatusOK,
		Msg:    "删除成功！",
		Data:   serializer.BuildTask(task, user.UserName),
	}
}
