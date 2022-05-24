package conf

import (
	"fmt"
	"gin-memos/db/mysql"
	"strings"

	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string

	RedisAddr     string
	RedisPassWord string
	RedisDbName   string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
)

func Init() {
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	LoadServer(file) // 加载服务器配置
	LoadMysql(file)  // 加载MySQL配置
	LoadRedis(file)  // 加载Redis配置
	// 综合易用性和性能，一般推荐使用 strings.Builder 来拼接字符串。(使用+效率很低)
	var builder strings.Builder
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	s := []string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=True&loc=Local"}
	for _, str := range s {
		builder.WriteString(str)
	}
	dsn := builder.String()
	mysql.DBInit(dsn) // 连接MySQL数据库
}

func LoadRedis(file *ini.File) {
	RedisAddr = file.Section("mysql").Key("RedisAddr").String()
	RedisPassWord = file.Section("mysql").Key("RedisPassWord").String()
	RedisDbName = file.Section("mysql").Key("RedisDbName").String()
}

func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}
