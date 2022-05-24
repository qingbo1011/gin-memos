package mysql

import (
	"fmt"
	"gin-memos/model"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB // 设置全局DB

func DBInit(dsn string) {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	db.LogMode(true)                             // 开启 Logger, 以展示详细的日志
	db.SingularTable(true)                       // 如果设置禁用表名复数形式属性为 true，`User` 的表名将是 `user`(因为gorm默认表名是复数)
	db.DB().SetMaxIdleConns(20)                  // 设置空闲连接池中的最大连接数
	db.DB().SetMaxOpenConns(100)                 // 设置数据库连接最大打开数。
	db.DB().SetConnMaxLifetime(time.Second * 30) // 设置可重用连接的最长时间

	DB = db
	migration() // 自动迁移
}

// 自动迁移：使用 migrate 来维持你的表结构一直处于最新状态。
// migrate仅支持创建表、增加表中没有的字段和索引。为了保护你的数据，它并不支持改变已有的字段类型或删除未被使用的字段。
func migration() {
	DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.User{}, &model.Task{}) // 数据迁移
	// 外键关联（添加外键）
	DB.Model(&model.Task{}).AddForeignKey("uid", "user(id)", "RESTRICT", "RESTRICT")
}
