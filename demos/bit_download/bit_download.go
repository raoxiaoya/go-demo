package bitdownload

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
)

/*
magnet:?xt=urn:btih:46be2d71cd2d690fdce026299eca164dbdaefe01&dn=[javdb.com]SONE-201
magnet:?xt=urn:btih:<infohash>&dn=文件名&tr=tracker地址

torrent



go install github.com/anacrolix/torrent/cmd/...@latest
torrent download magnet:?xt=urn:btih:KRWPCX3SJUM4IMM4YF5RPHL6ANPYTQPU


*/

func test() {
	var url string
	flag.StringVar(&url, "url", "https://download.virtualbox.org/virtualbox/6.1.6/VirtualBox-6.1.6.tar.bz2", "下载URL")
	flag.Parse()

	if url == "" {
		log.Println("请提供下载URL")
		return
	}
}

func Run() {
	// 1. 定义你的 magnet 链接
	magnetURI := "magnet:?xt=urn:btih:46be2d71cd2d690fdce026299eca164dbdaefe01&dn=[javdb.com]SONE-201" // 替换为真实的 magnet 链接
	outputFile := "demos/bit_download/downloaded.torrent" // 保存的 .torrent 文件名
	file , err := os.OpenFile(outputFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("创建文件失败: ", err)
	}
	defer file.Close()

	// 2. 创建 torrent 客户端配置
	config := torrent.NewDefaultClientConfig()
	// 可以调整监听端口等，这里使用默认
	client, err := torrent.NewClient(config)
	if err != nil {
		log.Fatal("创建客户端失败: ", err)
	}
	defer client.Close()

	// 3. 添加 magnet 链接到客户端
	tor, err := client.AddMagnet(magnetURI)
	if err != nil {
		log.Fatal("添加 magnet 失败: ", err)
	}

	// 4. 等待元数据下载完成
	fmt.Println("正在通过 DHT 和 peers 下载 .torrent 元数据...")
	metaCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute) // 设置超时
	defer cancel()

	select {
	case <-tor.GotInfo():
		// 元数据已成功下载
		fmt.Println("✅ 成功获取 .torrent 元数据!")

		// tor.DownloadAll()
		// client.WaitAll()

		// 5. 获取元数据
		mi := tor.Metainfo()
		// 6. 保存为 .torrent 文件
		err = mi.Write(file)
		// err = mi.WriteFile(outputFile)
		if err != nil {
			log.Fatal("保存 .torrent 文件失败: ", err)
		}

		fmt.Printf("🎉 .torrent 文件已保存为: %s\n", outputFile)

		// 可选：打印一些元数据信息
		info, _ := mi.UnmarshalInfo()
		fmt.Printf("文件名: %s\n", info.Name)
		fmt.Printf("文件大小: %d 字节\n", info.TotalLength())

	case <-metaCtx.Done():
		log.Fatal("❌ 下载元数据超时或失败: ", metaCtx.Err())
	}
}

// MagnetLink represents the parsed components of a magnet link
type MagnetLink struct {
	InfoHash string            // infohash (without urn:btih:)
	Name     string            // display name (dn)
	Trackers []string          // list of trackers (tr)
	Keywords []string          // keywords (kt)
	Other    map[string]string // other parameters
}

// ParseMagnetLink 解析 magnet 链接
func ParseMagnetLink(magnetURI string) (*MagnetLink, error) {
	// 必须以 magnet:? 开头
	if !strings.HasPrefix(magnetURI, "magnet:?") {
		return nil, fmt.Errorf("invalid magnet URI: does not start with 'magnet:?'")
	}

	// 使用 net/url 解析查询参数
	u, err := url.Parse(magnetURI)
	if err != nil {
		return nil, err
	}

	// 获取查询参数
	params := u.Query()

	// 初始化结果结构
	magnet := &MagnetLink{
		Other: make(map[string]string),
	}

	// 解析 xt (exact topic) 参数，通常是 urn:btih:<infohash>
	for _, xt := range params["xt"] {
		if strings.HasPrefix(xt, "urn:btih:") {
			magnet.InfoHash = xt[9:] // 去掉 urn:btih: 前缀
			break // 通常只取第一个 btih
		}
	}

	// 如果没有找到 infohash，返回错误
	if magnet.InfoHash == "" {
		return nil, fmt.Errorf("no valid infohash found in magnet link")
	}

	// 解析 dn (display name)
	if dn := params.Get("dn"); dn != "" {
		magnet.Name = dn
	}

	// 解析 tr (tracker)
	if trs, exists := params["tr"]; exists {
		magnet.Trackers = trs
	}

	// 解析 kt (keywords)
	if kts, exists := params["kt"]; exists {
		magnet.Keywords = kts
	}

	// 其他参数存入 Other
	for k := range params {
		if k != "xt" && k != "dn" && k != "tr" && k != "kt" {
			magnet.Other[k] = params.Get(k)
		}
	}

	return magnet, nil
}

func main() {
	// 示例 magnet 链接
	magnetURI := "magnet:?xt=urn:btih:ABCDEF1234567890ABCDEF1234567890ABCDEF12&dn=Example+File&tr=http%3A%2F%2Ftracker.example.com%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&kt=golang+torrent"

	magnet, err := ParseMagnetLink(magnetURI)
	if err != nil {
		log.Fatalf("Error parsing magnet link: %v", err)
	}

	// 输出解析结果
	fmt.Printf("InfoHash: %s\n", magnet.InfoHash)
	fmt.Printf("Name: %s\n", magnet.Name)
	fmt.Printf("Trackers:\n")
	for _, tr := range magnet.Trackers {
		fmt.Printf("  - %s\n", tr)
	}
	fmt.Printf("Keywords: %s\n", strings.Join(magnet.Keywords, ", "))
	if len(magnet.Other) > 0 {
		fmt.Printf("Other parameters: %v\n", magnet.Other)
	}
}

type TorrentFile struct {
	Announce string            `bencode:"announce"`
	Info     map[string]interface{} `bencode:"info"`
}

func parseTorrentFile(filePath string) (*TorrentFile, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content, _ := io.ReadAll(f)

	var torrent TorrentFile
	err = bencode.Unmarshal(content, &torrent)
	if err != nil {
		return nil, err
	}

	return &torrent, nil
}