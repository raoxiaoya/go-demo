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


初始的引导节点：
router.utorrent.com:6881      ------> 108.160.166.42
router.bittorrent.com:6881    ------> 157.240.2.36
dht.transmissionbt.com:6881   ------> 212.129.33.59
dht.aelitis.com:6881          ------> 34.229.89.117
router.silotis.us:6881        ------> 27.19.50.246
dht.libtorrent.org:25401      ------> 103.252.114.11
dht.anacrolix.link:42069      ------> 51.15.69.20
router.bittorrent.cloud:42069 ------> 51.15.69.20


*/

func Run() {
	// TestParseTorrentFile()
	// downloadTorrent()
	downloadTorrentByMagnet()
}
