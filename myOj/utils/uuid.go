package utils

import (
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strconv"
	"time"
)

func GetUUID() string {
	uid := uuid.NewV4().String()
	return uid
}

// GetRandom 生成验证码
func GetRandom() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}
