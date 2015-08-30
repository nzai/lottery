package config

import (
	"os"
	"path/filepath"
	"github.com/Unknwon/goconfig"
)

const(
	configFileName = "lottery.ini"
)

var configFile = &goconfig.ConfigFile{}

//	指定配置文件
func SetRootDir(root string) error {
	
	path := filepath.Join(root, configFileName)
	_, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
			return err
	}
	
	configFile, err = goconfig.LoadConfigFile(path)

	return err
}

//	获取配置项
func String(section, key, defaultValue string) string {	
	return configFile.MustValue(section,key,defaultValue)
}

//	获取配置项
func Int(section, key string, defaultValue int) int {
	return configFile.MustInt(section, key, defaultValue)
}

