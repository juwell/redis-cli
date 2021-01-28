package client

var (
	// Cfg redis连接配置
	Cfg *Config
)

func init() {
	Cfg = NewConfig()
}

// Config redis连接配置
type Config struct {
	HostIP      string
	HostPort    string
	HostSocket  string
	UserName    string
	PassWord    string
	DBNum       int
	ClusterMode int
	SlaveMode   int
	ShutDown    int
	Eval        string
}

// NewConfig 新建一个配置
func NewConfig() *Config {
	return &Config{
		HostIP:     `127.0.0.1`,
		HostPort:   `6379`,
		HostSocket: ``,
	}
}

// Version 返回当前版本号
func Version() string {
	return `0.0.1`
}

// RedisVersion 返回支持的最新的redis版本号
func RedisVersion() string {
	return `6.0.10`
}
