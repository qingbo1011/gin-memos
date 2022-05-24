package db

import (
	"fmt"
	"gin-memos/conf"
	"strings"
	"testing"

	"gopkg.in/ini.v1"
)

var (
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
)

func TestMysqlPath(t *testing.T) {
	file, err := ini.Load("../../conf/config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	LoadMysql(file) // 加载MySQL配置
	// 综合易用性和性能，一般推荐使用 strings.Builder 来拼接字符串。(使用+效率很低)
	var builder strings.Builder
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	s := []string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=True&loc=Local"}
	for _, str := range s {
		builder.WriteString(str)
	}
	dsn := builder.String()
	fmt.Println(dsn)
}

func TestMySqlInit(t *testing.T) {
	conf.Init()
	fmt.Println("MySQL连接成功！")
}

func LoadMysql(file *ini.File) {
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}
