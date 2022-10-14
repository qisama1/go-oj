package test

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var MakeRandInt1 = func() func() int {
	rand.Seed(time.Now().UnixNano())
	min := 240
	max := 800
	return func() int {
		return rand.Intn(max-min) + min
	}
}()

var MakeRandInt2 = func() int {
	rand.Seed(time.Now().UnixNano())
	min := 240
	max := 800
	return rand.Intn(max-min) + min
}
var wg = sync.WaitGroup{}

func TestSample(t *testing.T) {
	//randInt := MakeRandInt()
	//for {
	//	fmt.Printf("%d %d \n", MakeRandInt1(), MakeRandInt2())
	////}
	//wg.Add(1)
	//getFib(10)
	//wg.Wait()
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:123456@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

}

func main() {

}
func getFib(n int) int {
	cal := cloForFib()
	res := 1
	go func() {
		for i := 3; i <= n; i++ {
			res = cal()
			fmt.Printf("%d, %d, %p \n", i, res, &res)
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Printf("%d, %p \n", res, &res)
	return res
}

func cloForFib() func() int {
	a, b := 1, 1
	return func() int {
		a, b = a+b, a
		return a
	}
}
