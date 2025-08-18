package bit_download

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func TestParseMagnetLink() {
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
			break                    // 通常只取第一个 btih
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
