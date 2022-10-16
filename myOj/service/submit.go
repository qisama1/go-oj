package service

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"myOj/define"
	"myOj/models"
	"myOj/utils"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "请输入当前页，默认第一页"
// @Param size query int false "size"
// @Param user_identity query string false "user identity"
// @Param problem_identity query string false "problem identity"
// @Param status query int false "状态"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		fmt.Println("Get ProblemList Size strconv err ", err)
	}

	// page是1，其实是从0开始的offset
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage)) // 设置默认值
	if err != nil {
		fmt.Println("Get ProblemList Page strconv err ", err)
	}
	page = (page - 1) * size // 起始位置
	var count int64

	userIdentity := c.Query("user_identity")
	problemIdentity := c.Query("problem_identity")
	status, err := strconv.Atoi(c.Query("status"))

	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	submits := make([]*models.SubmitBasic, 0)
	data := make(map[string]interface{})
	err = tx.Count(&count).Offset(page).Limit(20).Find(&submits).Error
	if err != nil {
		c.JSON(500, gin.H{
			"code":    -1,
			"message": "GetSubmitList err, " + err.Error(),
		})
		return
	}
	data["submits"] = submits
	data["count"] = count
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// Submit
// @Tags 内部方法
// @Summary 创建问题
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/submit [post]
func Submit(c *gin.Context) {
	problemIdentity := c.Query("problem_identity")
	code, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Read code error: " + err.Error(),
		})
	}
	path, err := utils.CodeSave(code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "code save error:" + err.Error(),
		})
		return
	}
	u, _ := c.Get("user")
	user := u.(*utils.UserClaims)
	sb := &models.SubmitBasic{
		Identity:        utils.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    user.Identity,
		Path:            path,
	}
	// 代码判断
	pb := new(models.ProblemBasic)
	err = models.DB.Where("identity = ?", problemIdentity).Preload("TestCases").First(pb).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Problem error:" + err.Error(),
		})
		return
	}

	// 运行判断
	// 答案错误的判断
	WA := make(chan int)
	// 超内存的channel
	OOM := make(chan int)
	// 编译错误的channel
	CE := make(chan int)
	// 通过的个数
	var passCount int32
	// 互斥锁
	//var lock sync.Mutex
	// 结果信息
	var msg string
	for _, testCase := range pb.TestCases {
		testCase := testCase
		go func() {
			// 执行代码逻辑
			cmd := exec.Command("go", "run", path)
			var out, stderr bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = &out

			pipe, err := cmd.StdinPipe()
			if err != nil {
				log.Fatalln(err)
			}
			io.WriteString(pipe, testCase.Input)
			// 根据测试案例运行
			var bm runtime.MemStats
			runtime.ReadMemStats(&bm)
			if err := cmd.Run(); err != nil {
				log.Println(err, stderr.String())
				msg = "编译出错"
				CE <- 1
				return
			}
			var em runtime.MemStats
			runtime.ReadMemStats(&em)

			if em.Alloc/1024-bm.Alloc/1024 >= uint64(pb.MaxMem) {
				msg = "运行超内存了"
				OOM <- 1
				return
			}

			// 答案错误
			if testCase.Output != out.String() {
				msg = "答案错误"
				WA <- 1
				return
			}
			// 运行超内存
			//lock.Lock()
			atomic.AddInt32(&passCount, 1)
			//lock.Unlock()
		}()
	}
	//[-1-待判断，1-答案正确，2-答案错误，3-运行超时，4-运行超内存]
	select {
	case <-WA:
		sb.Status = 2
	case <-OOM:
		sb.Status = 4
	case <-time.After(time.Millisecond * time.Duration(pb.MaxRuntime)):
		if int(atomic.LoadInt32(&passCount)) == len(pb.TestCases) {
			sb.Status = 1
		} else {
			sb.Status = 3
		}
	case <-CE:
		sb.Status = 5
	}

	err = models.DB.Create(&sb).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "submit error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "提交成功",
		"data": map[string]interface{}{
			"status": sb.Status,
			"msg":    msg,
		},
	})
	return
}
