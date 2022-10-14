package models

import "gorm.io/gorm"

type ProblemCategory struct {
	gorm.Model
	ProblemId     string         `gorm:"column:problem_id;type:varchar(36);" json:"problem_id"`   // 问题的唯一标识
	CategoryId    string         `gorm:"column:category_id;type:varchar(38);" json:"category_id"` // 分类的唯一标识
	CategoryBasic *CategoryBasic `gorm:"foreignKey:id;references:category_id"`                    // 关联分类基础信息表
}

func (table *ProblemCategory) TableName() string {
	return "problem_category"
}
