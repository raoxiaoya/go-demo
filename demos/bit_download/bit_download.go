package bit_download

/*
magnet:?xt=urn:btih:46be2d71cd2d690fdce026299eca164dbdaefe01&dn=[javdb.com]SONE-201
magnet:?xt=urn:btih:ZOCMZQIPFFW7OLLMIC5HUB6BPCSDEOQU
magnet:?xt=urn:btih:KRWPCX3SJUM4IMM4YF5RPHL6ANPYTQPU

torrent

下载进度存储在本地 SQLite 数据库中

go install github.com/anacrolix/torrent/cmd/...@latest

torrent download magnet:?xt=urn:btih:KRWPCX3SJUM4IMM4YF5RPHL6ANPYTQPU
torrent download magnet:?xt=urn:btih:46be2d71cd2d690fdce026299eca164dbdaefe01

下载资源的同时保存 .torrent 文件
torrent download --save-metainfos=1 magnet:?xt=urn:btih:46be2d71cd2d690fdce026299eca164dbdaefe01

*/

func Run() {
	// TestParseTorrentFile()
	// downloadTorrent()
	downloadTorrentByMagnet()
}
