package child_process_stdio

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"time"
)

func Run() {
	// 定义要执行的命令（例如：运行一个 Python 脚本）
	// -u : 取消缓冲，实时输出
	// -c : 执行一段python代码
	cmd := exec.Command("python", "-u", "-c", `user_input = input(); print("child get message: " + user_input)`)

	// 获取子进程的标准输入、输出和错误流
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		return
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}

	// 启动子进程，异步的方式执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// 使用 bufio 读取子进程的输出
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				fmt.Println("Error reading from stdout:", err)
				break
			}
			if line != "" {
				fmt.Print(line)
			}
			if err == io.EOF {
				break
			}
		}
	}()

	// 向子进程的标准输入写入数据
	input := "Hello from Go!\n"
	fmt.Fprintf(stdin, input)

	// 等待子进程结束
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command:", err)
	}

	time.Sleep(3 * time.Second)
}
