package parallel_download

/*

实现多 goroutine 并行下载大文件的功能

使用 HTTP 请求下载一个大文件；
将文件分成多个部分（分块）；
每个部分由一个独立的 goroutine 并行下载；
下载完成后合并成完整文件；
支持断点续传（通过 Range 请求）；
使用 sync.WaitGroup 协调多个 goroutine。

需要重试机制，因为服务器会有并发限制或者中断导致下载失败。

153M
https://download.virtualbox.org/virtualbox/6.1.6/VirtualBox-6.1.6.tar.bz2

*/

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

func Run() {
	var url string
	flag.StringVar(&url, "url", "https://download.virtualbox.org/virtualbox/6.1.6/VirtualBox-6.1.6.tar.bz2", "下载URL")
	flag.Parse()

	if url == "" {
		log.Println("请提供下载URL")
		return
	}

	fmt.Println("开始下载:", url)

	task, err := NewTask(url)
	if err != nil {
		fmt.Printf("创建任务失败: %v\n", err)
		return
	}

	task.Execute()
}

type DownloadState struct {
	FileSize int64         `json:"file_size"`
	Chunks   []*ChunkState `json:"chunks"`
}

type ChunkState struct {
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	DoneBytes int64  `json:"done_bytes"` // 已下载的字节数
	Done      bool   `json:"done"`       // 是否已完成
	Path      string `json:"path"`       // 分块路径
}

type Task struct {
	NumWorkers int            // 并发数
	Url        string         // 地址
	TempDir    string         // 临时目录
	StateFile  string         // 保存下载状态
	OutputFile string         // 输出文件路径
	State      *DownloadState // 下载状态
	Bar        *progressbar.ProgressBar
	BarMu      sync.Mutex
	Wg         sync.WaitGroup
}

func NewTask(url string) (*Task, error) {
	task := &Task{
		Url:        url,
		NumWorkers: 4,
	}
	task.TempDir = filepath.Base(url)
	task.StateFile = filepath.Join(task.TempDir, "progress.json")
	task.OutputFile = filepath.Join(task.TempDir, task.TempDir)

	// 创建临时目录
	if err := os.MkdirAll(task.TempDir, 0755); err != nil {
		return nil, fmt.Errorf("无法创建临时目录: %w", err)
	}

	// 加载或初始化下载状态
	if err := task.LoadOrCreateState(); err != nil {
		return nil, err
	}

	task.SaveState()

	fmt.Printf("文件总大小: %d bytes\n", task.State.FileSize)

	// 初始化进度条（只统计未完成的部分）
	bar := progressbar.NewOptions64(
		task.State.FileSize,
		progressbar.OptionSetDescription("下载中..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
	)
	bar.Set64(task.GetCurrentDownloaded())
	task.Bar = bar

	return task, nil
}

func (t *Task) Execute() {
	go func() {
		ti := time.NewTicker(100 * time.Millisecond)
		for range ti.C {
			t.SaveState()
		}
	}()

	if err := t.Download(); err != nil {
		fmt.Printf("下载失败: %v\n", err)
		return
	}

	// 所有块完成，合并文件
	fmt.Println("所有分块下载完成，正在合并...")
	if err := t.MergeChunks(); err != nil {
		fmt.Printf("合并失败: %v\n", err)
		return
	}

	fmt.Printf("下载完成: %s\n", t.OutputFile)

	// 清理状态文件和临时目录（可选）
	_ = os.RemoveAll(t.TempDir)
}

func (t *Task) Download() error {
	errs := make([]error, t.NumWorkers)

	// 为每个未完成的块启动 goroutine
	for i, chunk := range t.State.Chunks {
		if chunk.Done {
			continue
		}
		time.Sleep(1 * time.Second)
		t.Wg.Add(1)
		go func(index int) {
			defer t.Wg.Done()

			var downloadErr error
			for range 3 {
				downloadErr = t.DownloadChunkIfNotDone(index)
				if downloadErr == nil {
					// 标记为完成
					t.State.Chunks[index].Done = true
					break
				} else {
					time.Sleep(100 * time.Millisecond)
					continue
				}
			}
			if downloadErr != nil {
				e := fmt.Errorf("分块 %d 下载失败: %v", index, downloadErr)
				errs = append(errs, e)
			}
		}(i)
	}

	t.Wg.Wait()

	for _, err := range errs {
		if err != nil {
			fmt.Println(err)
		}
	}

	if errs[0] != nil {
		return errs[0]
	} else {
		return nil
	}
}

// 获取当前已下载的总字节数（用于恢复进度条）
func (t *Task) GetCurrentDownloaded() int64 {
	var total int64
	for _, c := range t.State.Chunks {
		if c.Done {
			total += (c.End - c.Start + 1)
		} else {
			total += c.DoneBytes
		}
	}
	return total
}

// 加载状态，如果不存在则创建
func (t *Task) LoadOrCreateState() error {
	// 获取文件大小
	fileSize, err := getFileSize(t.Url)
	if err != nil {
		return errors.New("无法获取文件大小: " + err.Error())
	}

	state := &DownloadState{
		FileSize: fileSize,
		Chunks:   make([]*ChunkState, t.NumWorkers),
	}

	// 如果状态文件存在，尝试加载
	data, err := os.ReadFile(t.StateFile)
	if err == nil {
		if json.Unmarshal(data, state) == nil {
			// 验证文件大小是否一致
			if state.FileSize == fileSize {
				log.Println("恢复之前的下载进度")
				// 校正进度数据
				for k, c := range state.Chunks {
					fi, err := os.Stat(c.Path)
					if err == nil {
						// fmt.Println(c.Path, fi.Size())
						if fi.Size() == (c.End - c.Start + 1) {
							c.Done = true
						} else {
							c.Done = false
						}
						c.DoneBytes = fi.Size()
					} else {
						c.Done = false
						c.DoneBytes = 0
						os.Remove(c.Path)
					}
					state.Chunks[k] = c
				}
				t.State = state
				return nil
			} else {
				log.Println("文件大小变化，重新开始下载")
			}
		}
	}

	// 否则创建新状态
	for i := 0; i < t.NumWorkers; i++ {
		start := int64(i) * (fileSize / int64(t.NumWorkers))
		end := start + (fileSize / int64(t.NumWorkers)) - 1
		if i == t.NumWorkers-1 {
			end = fileSize - 1
		}

		state.Chunks[i] = &ChunkState{
			Start: start,
			End:   end,
			Done:  false,
			Path:  filepath.Join(t.TempDir, fmt.Sprintf("part_%d.temp", i)),
		}
	}

	t.State = state
	if err := t.SaveState(); err != nil {
		return err
	}
	fmt.Println("开始新的下载任务")
	return nil
}

func (t *Task) SaveState() error {
	data, err := json.MarshalIndent(t.State, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(t.StateFile, data, 0644)
}

// 下载分块（如果未完成）
func (t *Task) DownloadChunkIfNotDone(index int) error {
	c := t.State.Chunks[index]

	// 否则下载
	client := &http.Client{}
	req, err := http.NewRequest("GET", t.Url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", c.Start+c.DoneBytes, c.End))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 206 {
		return fmt.Errorf("不支持 Range，状态码: %d", resp.StatusCode)
	}
	file, err := os.OpenFile(c.Path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := &progressReader{
		ReadCloser: resp.Body,
		Index:      index,
		Task:       t,
	}

	_, err = io.Copy(file, reader)
	return err
}

// 合并所有分块
func (t *Task) MergeChunks() error {
	outFile, err := os.Create(t.OutputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 按 index 排序合并
	var indices []int
	for i := range t.State.Chunks {
		indices = append(indices, i)
	}
	sort.Ints(indices)

	for _, i := range indices {
		c := t.State.Chunks[i]
		partFile, err := os.Open(c.Path)
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, partFile)
		partFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// 自定义带进度的 Reader
type progressReader struct {
	io.ReadCloser
	Index int
	Task  *Task
}

func (r *progressReader) Read(p []byte) (n int, err error) {
	n, err = r.ReadCloser.Read(p)
	if n > 0 {
		_ = r.Task.Bar.Add(n)

		r.Task.State.Chunks[r.Index].DoneBytes += int64(n)
	}
	return n, err
}

func getFileSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return resp.ContentLength, nil
}
