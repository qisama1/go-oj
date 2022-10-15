package test

import (
	"fmt"
	"myOj/models"
	"myOj/utils"
	"testing"
)

//func TestGormTest(t *testing.T) {
//	dsn := "root:123456@tcp(127.0.0.1:3306)/gin_gorm_oj?charset=utf8mb4&parseTime=True&loc=Local"
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		t.Fatal(err)
//	}
//	data := make([]*models.Problem, 0)
//	err = db.Find(&data).Error
//	if err != nil {
//		t.Fatal(err)
//	}
//	for _, v := range data {
//		fmt.Printf("Problem ===> %v /n", v)
//	}
//}

func TestGorm(t *testing.T) {
	data := &models.UserBasic{
		Identity: utils.GetUUID(),
		Name:     "abc",
		Password: utils.GetMd5("123456"),
		Phone:    "123",
		Mail:     "123@3333.com",
	}
	fmt.Printf("%v", data)
	models.DB.Create(data)
}
