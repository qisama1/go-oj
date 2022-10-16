package models

import (
	"gorm.io/gorm"
)

type ProblemBasic struct {
	gorm.Model
	Identity          string             `gorm:"column:identity;type:varchar(36);" json:"identity"` // 问题表的唯一标识
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id"`               // 关联问题分类表
	Title             string             `gorm:"column:title;type:varchar(255);" json:"title"`      // 文章标题
	Content           string             `gorm:"columns:content;type:varchar(255);" json:"content"` // 文章正文
	MaxRuntime        int                `gorm:"column:max_runtime;type:int(11);" json:"max_runtime"`
	MaxMem            int                `gorm:"column:max_mem;type:int(11);" json:"max_mem" `
	TestCases         []*TestCase        `gorm:"foreignKey:problem_id;references:id"`
}

func (table *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	tx := DB.Model(new(ProblemBasic)).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").Preload("TestCases").
		Where("title like ? OR content like ?", "%"+keyword+"%",
			"%"+keyword+"%") // 指定查询的表
	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problem_basic.id").
			Where("pc.category_id = (SELECT cb.id FROM category_basic cb WHERE cb.identity = ?)", categoryIdentity)
	}
	//DB.Raw("SELECT id, identity, problem_categories, title, content, max_runtime," +
	//	"max_mem from problem_basic left join problem_category pc on pc.problem_id = problem_basic.id" +
	//	"where pc.category_id = (SELECT cb.id FROM category_basic cb WHERE cb.identity = ?)", categoryIdentity)
	return tx
}
