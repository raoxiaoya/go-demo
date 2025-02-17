package progressbar

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

func Run() {
	f2()
}

func f1() {
	bar := progressbar.Default(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(40 * time.Millisecond)
	}
}

func f2() {
	req, _ := http.NewRequest("GET", "https://dl.google.com/go/go1.14.2.src.tar.gz", nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	f, _ := os.OpenFile("go1.14.2.src.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"下载中",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
}
