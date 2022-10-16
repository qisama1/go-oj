package models

import "gorm.io/gorm"

type TestCase struct {
	gorm.Model
	Identity  string `gorm:"column:identity;type:varchar(255)" json:"identity"`
	ProblemId uint   `gorm:"column:problem_id;type:varchar(255)" json:"problem_id"`
	Input     string `gorm:"column:input;type:text;" json:"input"`
	Output    string `gorm:"column:input;type:text;"json:"output"`
}

func (table *TestCase) TableName() string {
	return "test_case"
}
