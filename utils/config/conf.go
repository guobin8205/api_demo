package conf

import (
	"path/filepath"

	"fmt"

	"github.com/go-ini/ini"
)

type Config struct {
	conf *ini.File
}

var (
	c = new(Config)

	//项目根路径
	AppPath string
	//配置文件路径
	ConfPath string
	//当前选择的配置文件
	RunMode  string
	HttpPort string
	RPCServer string
)

func init() {
	c.conf = parseConfig()
	//只读操作增加性能
	c.conf.BlockMode = false
	selConfFile := c.conf.Section("").Key("loadfile").String()
	loadConfFile := filepath.Join(ConfPath, selConfFile)
	if selConfFile == "" {
		panic(fmt.Sprintf("no config file to load:%s", loadConfFile))
	}
	if err := c.conf.Append(loadConfFile); err != nil {
		panic(err)
	}
	RunMode = c.conf.Section("").Key("runmode").String()
	HttpPort = c.conf.Section("").Key("httpport").String()
	RPCServer = c.conf.Section("").Key("rpcserver").String()
}

func String(key string) string {
	return c.conf.Section("").Key(key).String()
}

func MustString(key, def string) string {
	return c.conf.Section("").Key(key).MustString(def)
}

func Int(key string) int {
	r, err := c.conf.Section("").Key(key).Int()
	if err != nil {
		return 0
	}
	return r
}

func MustInt(key string, def int) int {
	return c.conf.Section("").Key(key).MustInt(def)
}

func parseConfig() *ini.File {
	AppPath, _ = filepath.Abs("./")
	ConfPath = filepath.Join(AppPath, "conf")

	//载入入口配置文件
	indexConfPath := filepath.Join(ConfPath, "app.conf")
	conf, err := ini.Load(indexConfPath)
	if err != nil {
		panic(err)
	}
	return conf
}
