// configs 用于对配置文件的定义以及读取
// 目前使用yaml文件格式，比起json来说要更好看一些
package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	JWTContextKey = "user"
)

var (
	// C 全局配置文件，在Init调用前为nil
	C *Config
)

const (
	DurationCodeExpire = time.Minute
)

// Config 配置
type Config struct {
	App     app     `yaml:"app"`
	DB      db      `yaml:"db"`
	Redis   redis   `yaml:"redis"`
	JWT     jwt     `yaml:"jwt"`
	Qiniu   qiniu   `yaml:"qiniu"`
	Wechat  wechat  `yaml:"wechat"`
	LogConf logConf `yaml:"logConf"`
	Debug   bool    `yaml:"debug"`
}

type app struct {
	Addr   string `yaml:"addr"`
	Prefix string `yaml:"prefix"`
}

type db struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
}

type redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type jwt struct {
	Secret string   `yaml:"secret"`
	Skip   []string `yaml:"skip"`
}

type qiniu struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"` // 空间
	URL       string `yaml:"url"`
}
type wechat struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

type logConf struct {
	LogPath     string `yaml:"log_path"`
	LogFileName string `yaml:"log_file_name"`
}

func init() {
	configFile := "default.yml"

	// 如果有设置 ENV ，则使用ENV中的环境
	if v, ok := os.LookupEnv("ENV"); ok {
		configFile = v + ".yml"
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(fmt.Sprintf("configs/%s", configFile))

	if err != nil {
		log.Println("Read configs error!")
		log.Panic(err)
		return
	}

	config := &Config{}

	err = yaml.Unmarshal(data, config)

	if err != nil {
		log.Println("Unmarshal configs error!")
		log.Panic(err)
		return
	}

	C = config

	log.Println("Config " + configFile + " loaded.")
	if C.Debug {
		log.Printf("%+v\n", C)
	}

}
