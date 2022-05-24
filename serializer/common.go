package serializer

import "net/http"

// Response 通用返回体（基础序列化器）
type Response struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

// TokenData 带有token的Data数据
type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

// DataList 带有总数的Data结构（这里我加上用户名，专门用户分页查询时的返回体）
type DataList struct {
	Item     interface{} `json:"item"`
	Total    uint        `json:"total"`
	UserName string      `json:"user_name"`
}

func BuildListResponse(items interface{}, total uint, uname string) *Response {
	return &Response{
		Status: http.StatusOK,
		Msg:    "序列化成功",
		Data: DataList{
			Item:     items,
			Total:    total,
			UserName: uname,
		},
	}
}
