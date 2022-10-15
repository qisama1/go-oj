package test

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestUUID(t *testing.T) {
	uid := uuid.NewV4().String()
	fmt.Println(uid)
}
