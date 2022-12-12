package log

type Config struct {
	Online     bool   `json:"online"`     // 是否是生产环境
	Service    string `json:"service"`    // 服务名
	Filename   string `json:"filename"`   // 日志文件名（含路径）
	MaxSize    int    `json:"maxSize"`    // 单文件大小限制（单位：MB）
	MaxBackups int    `json:"maxBackups"` // 最大滚动备份数
	MaxAge     int    `json:"maxAge"`     // 最大保留天数
	Compress   bool   `json:"compress"`   // 是否压缩
}

// DefaultConfig 默认配置
var DefaultConfig = &Config{
	Filename:   "./logs/log.txt",
	MaxSize:    100,
	MaxBackups: 10,
	MaxAge:     30,
	Online:     false,
}
