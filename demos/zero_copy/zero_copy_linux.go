package zero_copy

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"
)

func Run() {
	test3()
}

func example(src, dst net.Conn) {
	// 由于 src 及 dst 都是 socket, 没法直接使用 splice, 因此先创建临时 pipe
	var r, w int
	var fds = []int{r, w} // readfd, writefd
	if err := syscall.Pipe(fds); err != nil {
		panic(err)
	}
	// 使用完后关闭 pipe
	defer syscall.Close(fds[0])
	defer syscall.Close(fds[1])
	// 获取 src fd
	srcfile, err := src.(*net.TCPConn).File()
	if err != nil {
		panic(err)
	}
	srcfd := int(srcfile.Fd())
	syscall.SetNonblock(srcfd, true)
	// 从 srcfd 读出, 写入 fds[1] (pipe write fd)
	DEFAULTSIZE := 1 << 10
	SPLICE_F_NONBLOCK := 0x2
	num, err := syscall.Splice(srcfd, nil, fds[1], nil, DEFAULTSIZE, SPLICE_F_NONBLOCK)
	fmt.Println(num, err)
}

func test2() {
	src, _ := os.Open("/tmp/source.txt")
	defer src.Close()

	target, _ := os.Create("/tmp/target.txt")
	defer target.Close()

	// 创建管道文件
	// 作为两个文件传输数据的中介
	pipeReader, pipeWriter, _ := os.Pipe()
	defer pipeReader.Close()
	defer pipeWriter.Close()

	// 设置文件读写模式
	// 笔者在标准库中没有找到对应的常量说明
	// 读者可以参考这个文档:
	//   https://pkg.go.dev/golang.org/x/sys/unix#pkg-constants
	//   SPLICE_F_NONBLOCK = 0x2
	spliceNonBlock := 0x02

	// 使用 Splice 将数据从源文件描述符移动到管道 writer
	syscall.Splice(int(src.Fd()), nil, int(pipeWriter.Fd()), nil, 1024, spliceNonBlock)

	// 使用 Splice 将数据从管道 reader 移动到目标文件描述符
	syscall.Splice(int(pipeReader.Fd()), nil, int(target.Fd()), nil, 1024, spliceNonBlock)
}

func test1() {
	src, _ := os.Open("/tmp/source.txt")
	defer src.Close()

	target, _ := os.Create("/tmp/target.txt")
	defer target.Close()

	syscall.Sendfile(int(src.Fd()), int(target.Fd()), nil, 1024)
}

func test3() {
	src, err := os.Open("./demos/zero_copy/source.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	target, err := os.Create("./demos/zero_copy/target.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer target.Close()

	_, err = io.Copy(target, src)
	if err != nil {
		log.Fatal(err)
	}
}
