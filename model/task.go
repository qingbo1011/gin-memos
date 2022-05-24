package model

import (
	"github.com/jinzhu/gorm"
)

type Task struct {
	gorm.Model // gorm.Model 是一个包含一些基本字段的结构体, 包含的字段有 ID，CreatedAt， UpdatedAt， DeletedAt。
	// gorm 用 tag 的方式来标识 mysql 里面的约束
	User User `gorm:"foreignkey:Uid"` // Task和User是一对一关系，这里使用外键关联
	Uid  uint `gorm:"not null"`       // 设置不为空

	Tittle     string `gorm:"index,not null"` // tittle字段为索引且不为空
	Status     int    `gorm:"default:0"`      // 0表示List未完成，1为已完成
	Context    string // List内容（可以通过`gorm:"type:longtext"`设置该字段类型，也可以后续在Navicat中设置该字段类型）
	StartTime  int64  // List开始时间(使用时间戳)
	FinishTime int64  `gorm:"default:0"` // List完成时间(使用时间戳)
}
