package service

import (
	"gin-memos/db/mysql"
	"gin-memos/model"
	"gin-memos/pkg/util"
	"gin-memos/serializer"
	"net/http"

	"github.com/jinzhu/gorm"
)

type UserService struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (s *UserService) Register() serializer.Response {
	var user model.User
	var count int
	mysql.DB.Where(&model.User{UserName: s.UserName}).First(&user).Count(&count)
	// SELECT * FROM user WHERE name = s.UserName LIMIT 1;
	if count == 1 {
		return serializer.Response{
			Status: http.StatusBadRequest,
			Msg:    "用户名重复！",
		}
	}
	// 如果数据库中没有该用户，那么就开始注册
	user.UserName = s.UserName
	err := user.SetPassWord(s.Password)
	if err != nil {
		return serializer.Response{
			Status: http.StatusBadRequest,
			Msg:    "存储密码发生错误！",
			Error:  err.Error(),
		}
	}
	// 加密成功就可以创建用户了
	err = mysql.DB.Create(&user).Error
	if err != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库插入数据出错！",
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: http.StatusOK,
		Msg:    "用户创建成功！",
	}
}
func (s *UserService) Login() serializer.Response {
	var user model.User
	err := mysql.DB.Where(&model.User{UserName: s.UserName}).First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) { // 数据库中没有找到记录
			return serializer.Response{
				Status: http.StatusBadRequest,
				Msg:    "该用户不存在，请先注册！",
				Error:  err.Error(),
			}
		}
		// 不是用户不存在却还是出错，其他不可抗拒的因素
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库发生错误！",
			Error:  err.Error(),
		}
	}
	// 用户从数据库中找到了
	success, err := user.CheckPassword(s.Password)
	if err != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "登录失败！",
			Error:  err.Error(),
		}
	}
	if !success {
		return serializer.Response{
			Status: http.StatusBadRequest,
			Msg:    "密码错误，登录失败！",
		}
	}
	// 登录成功要分发token（其他功能需要身份验证，给前端存储的）
	// 创建一个备忘录是需要携带一个有效token的，不然都不知道是谁创建的备忘录
	token, err := util.GenerateToken(user.ID, user.UserName)
	if err != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "token签发错误！",
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: http.StatusOK,
		Msg:    "登录成功！",
		Data: serializer.TokenData{
			User:  serializer.BuildUser(user),
			Token: token,
		},
	}
}
