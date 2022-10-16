package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Identity         string `gorm:"column:identity;type:varchar(36);" json:"identity"` // 用户表的唯一标识
	Name             string `gorm:"column:name;type:varchar(100);" json:"name"`        // 用户名
	Password         string `gorm:"column:password;type:varchar(32);" json:"name""`
	Phone            string `gorm:"column:phone;type:varchar(20);" json:"phone"`
	Mail             string `gorm:"column:mail;type:varchar(100);" json:"mail"`
	FinishProblemNum int64  `gorm:"column:finish_problem_num;type:int(64);" json:"finish_problem_num"`
	SubmitProblemNum int64  `gorm:"column:submit_problem_num;type:int(64);" json:"submit_problem_num"`
	IsAdmin          int    `gorm:"column:is_admin;type:tinyint(1);" json:"is_admin"` //[1-管理员，0-不是]
}

func (user *UserBasic) TableName() string {
	return "user_basic"
}

//
//func GetUserDetail(identity string) *gorm.DB {
//	tx := DB.Omit("password").Where("identity = ?", identity)
//	return tx
//}
