package serializer

import "gin-memos/model"

type UserResponse struct {
	Uid      uint   `json:"uid"`
	UserName string `json:"user_name"`
}

func BuildUser(user model.User) *UserResponse {
	return &UserResponse{
		Uid:      user.ID,
		UserName: user.UserName,
	}
}
