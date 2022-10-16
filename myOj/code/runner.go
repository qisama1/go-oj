package main

import (
	"bytes"
	"io"
	"log"
	"os/exec"
)

func main() {
	// go run code-user/main.go
	cmd := exec.Command("go", "run", "code-user/main.go")
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out

	pipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	io.WriteString(pipe, "23 11 \n")
	// 根据测试案例运行
	if err := cmd.Run(); err != nil {
		log.Fatalln(err, stderr.String())
	}
	log.Println(out.String())
}
