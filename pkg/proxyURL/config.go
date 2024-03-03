package proxyURL

type Config struct {
	// SchemeはHTTPSで確定しているため、Hostのみ指定させる

	// 前段に立っているHTTPS終端を行うリバースプロキシを表すホスト名
	Host string
}
