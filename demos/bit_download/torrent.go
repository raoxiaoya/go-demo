package bit_download

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
)

/*

解析 torrent 文件

github.com/jackpal/bencode-go
github.com/anacrolix/torrent/bencode

*/

func TestParseTorrentFile() {
	parseTorrentFile("46be2d71cd2d690fdce026299eca164dbdaefe01.torrent")
	// parseTorrentFile("demos/bit_download/downloaded.torrent")
	// parseTorrentFile("demos/bit_download/46be2d71cd2d690fdce026299eca164dbdaefe01.torrent")
}

// TorrentFile represents the structure of a .torrent file.
type TorrentFile struct {
	Announce     string      `bencode:"announce"`
	AnnounceList [][]string  `bencode:"announce-list"`
	UrlList      []string    `bencode:"url-list"`
	Nodes        []string    `bencode:"nodes,omitempty,ignore_unmarshal_type_error"`
	Info         TorrentInfo `bencode:"info"`
}

// TorrentInfo contains details about the files to be downloaded.
type TorrentInfo struct {
	Name         string             `bencode:"name"`
	PieceLength  int                `bencode:"piece length"`
	Pieces       string             `bencode:"pieces"`
	Length       int                `bencode:"length,omitempty"`
	Files        []TorrentFileEntry `bencode:"files,omitempty"`
	Publisher    string             `bencode:"publisher,omitempty"`
	PublisherUrl string             `bencode:"publisher-url,omitempty"`
}

// TorrentFileEntry represents individual file information in a multi-file torrent.
type TorrentFileEntry struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

// 解析 torrent 文件，得到文件列表
func parseTorrentFile(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open torrent file: %v", err)
	}
	defer f.Close()
	content, _ := io.ReadAll(f)
	var torrent TorrentFile
	err = bencode.Unmarshal(content, &torrent)
	if err != nil {
		log.Fatalf("Failed to decode torrent file: %v", err)
	}
	fmt.Printf("Torrent Name: %s\n", torrent.Info.Name)

	fmt.Println(torrent.Nodes)
	fmt.Println(torrent.Info.Publisher)
	fmt.Println(torrent.Info.PublisherUrl)

	if len(torrent.Info.Files) > 0 {
		for _, file := range torrent.Info.Files {
			fmt.Printf("File with path '%s' and length %d bytes.\n",
				strings.Join(file.Path, "/"),
				file.Length)
		}
	} else if torrent.Info.Length > 0 {
		fmt.Printf("Single file with length %d bytes.\n", torrent.Info.Length)
	} else {
		fmt.Println("No files found in the torrent.")
	}
}

// 根据 magnet 下载 torrent 文件
func downloadTorrent() {
	file, _ := os.OpenFile("demos/bit_download/46be2d71cd2d690fdce026299eca164dbdaefe01.torrent", os.O_CREATE|os.O_RDWR, 0644)
	defer file.Close()

	c, _ := torrent.NewClient(nil)
	defer c.Close()

	t, _ := c.AddMagnet("magnet:?xt=urn:btih:46be2d71cd2d690fdce026299eca164dbdaefe01")
	<-t.GotInfo()

	mi := t.Metainfo()
	mi.Write(file)
	fmt.Println("成功获取 .torrent 元数据!")
}

// 下载文件
func downloadTorrentByMagnet() {
	dirname := "demos/bit_download/data2"
	torrentFile, _ := os.OpenFile(dirname + "/46be2d71cd2d690fdce026299eca164dbdaefe01.torrent", os.O_CREATE|os.O_RDWR, 0644)
	defer torrentFile.Close()

	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = dirname
	c, _ := torrent.NewClient(cfg)
	defer c.Close()

	t, _ := c.AddMagnet("magnet:?xt=urn:btih:46be2d71cd2d690fdce026299eca164dbdaefe01")

	<-t.GotInfo()

	log.Println("成功获取 torrent 元数据!")

	// 保存 torrent 文件
	mi := t.Metainfo()
	mi.Write(torrentFile)

	downloadAll := false

	if downloadAll {
		// 下载全部文件
		t.DownloadAll()
	} else {
		// 指定文件下载
		fileList := []string{
			"SONE-201a.jpg",
			"文宣/xp1024.com -1024核工厂.url",
		}

		for _, f := range t.Files() {
			for _, fileArg := range fileList {
				if f.DisplayPath() == fileArg {
					f.Download()
				}
			}
		}
	}
}
