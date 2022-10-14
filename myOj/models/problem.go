package models

import (
	"gorm.io/gorm"
)

type Problem struct {
	gorm.Model
	Identity   string `gorm:"column:identity;type:varchar(36);" json:"identity"`        // 问题表的唯一标识
	CategoryId string `gorm:"column:category_id;type:varchar(255);" json:"category_id"` // 分类ID， 逗号分割
	Title      string `gorm:"column:title;type:varchar(255);" json:"title"`             // 文章标题
	Content    string `gorm:"columns:content;type:varchar(255);" json:"content"`        // 文章正文
	MaxRuntime int    `gorm:"column:max_runtime;type:int(11);" json:"max_runtime"`
	MaxMem     int    `gorm:"column:max_mem;type:int(11);" json:"max_mem" `
}

func (table *Problem) TableName() string {
	return "problem"
}

func GetProblemList(keyword string) *gorm.DB {
	return DB.Model(new(Problem)).
		Where("title like ? OR content like ?", "%"+keyword+"%",
			"%"+keyword+"%") // 指定查询的表
}
