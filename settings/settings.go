package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 使用结构体来存储配置信息 - 实现配置信息的更清晰的展示
type PrjConfig struct {
	*AppConfig    `mapstructure:"app"`
	*SystemConfig `mapstructure:"system"`
	*AuthConfig   `mapstructure:"auth"`
	*LogConfig    `mapstructure:"log"`
	*MysqlConfig  `mapstructure:"mysql"`
	*RedisConfig  `mapstructure:"redis"`
}

type AppConfig struct {
	Name           string `mapstructure:"name"`
	Mode           string `mapstructure:"mode"`
	Version        string `mapstructure:"version"`
	Port           string `mapstructure:"port"`
	StartTime      string `mapstructure:"start_time"`
	SqlFile        string `mapstructure:"sqlfile"`
	AdminBasePath  string `mapstructure:"admin_base_path"`  // admin管理页访问的路径
	ClientBasePath string `mapstructure:"client_base_path"` // 客户端访问的基路径
	MachineId      int64  `mapstructure:"machine_id"`
}

type SystemConfig struct {
	PageSize int8 `mapstructure:"page_size"` // 文章页每页默认的显示数量
}

type AuthConfig struct {
	CodeNum     int   `mapstructure:"code_num"`
	ExpiredTime int64 `mapstructure:"verify_expired"`
}

type LogConfig struct {
	Level     string `mapstructure:"level"`
	FileName  string `mapstructure:"file_name"`
	MaxAge    int    `mapstructure:"max_age"`
	MaxSize   int    `mapstructure:"max_size"`
	MaxBackup int    `mapstructure:"max_backup"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Schema       string `mapstructure:"schema"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"poolsize"`
}

var Confg = new(PrjConfig)

// 初始化viper初始化
func Init() (err error) {
	viper.SetConfigFile("./config/config.yaml") //也可以指定配置文件路径
	//viper.SetConfigName("config") // 不加后缀
	//viper.SetConfigType("yaml") // 这个专用于远程获取配置时指定文件类型
	//viper.AddConfigPath(".")

	if err = viper.ReadInConfig(); err != nil {
		fmt.Printf("Read config error!!! Error is :#{err}\n")
		return
	}
	if err = viper.Unmarshal(Confg); err != nil {
		fmt.Printf("Unmarshal config error!!! Error is :#{err}\n")
		return
	}

	viper.WatchConfig() // 实时监控配置文件的更改
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file has changed!!!")
		if err = viper.Unmarshal(Confg); err != nil {
			fmt.Printf("Unmarshal config error!!! Error is :#{err}\n")
		}
	})
	return
}
