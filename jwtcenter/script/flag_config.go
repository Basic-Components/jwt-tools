package script

import (
	"path"
	"strings"

	"github.com/small-tk/pathlib"
	"github.com/spf13/pflag"
)

//InitFlagConfig 从命令行获取配置
func InitFlagConfig() (map[string]interface{}, error) {
	loglevel := pflag.StringP("loglevel", "l", "", "log的等级")
	address := pflag.StringP("address", "a", "", "要启动的服务器地址")
	privatekeypath := pflag.StringP("private_key_path", "r", "", "签名使用的私钥地址")
	publickeypath := pflag.StringP("public_key_path", "u", "", "验签使用的公钥地址")
	hashkey := pflag.StringP("hash_key", "k", "", "hash签名使用的盐")
	name := pflag.StringP("name", "n", "", "服务名称")
	confPath := pflag.StringP("config", "c", "", "配置文件位置")
	pflag.Parse()
	var flagConfig = map[string]interface{}{}

	if *confPath != "" {
		p, err := pathlib.New(*confPath).Absolute()
		if err != nil {
			return flagConfig, err
		}
		if p.Exists() && p.IsFile() {
			filenameWithSuffix := path.Base(*confPath)
			fileSuffix := path.Ext(filenameWithSuffix)
			fileName := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
			dir, err := p.Parent()
			if err != nil {
				return flagConfig, err
			}
			filePaths := []string{dir.Path}
			targetfileconf, err := SetFileConfig(fileName, filePaths)
			if err != nil {
				return flagConfig, err
			}
			for k, v := range targetfileconf {
				flagConfig[k] = v
			}
		}
	}
	if *loglevel != "" {
		flagConfig["LOG_LEVEL"] = *loglevel
	}

	if *address != "" {
		flagConfig["ADDRESS"] = *address
	}
	if *privatekeypath != "" {
		flagConfig["PRIVATE_KEY_PATH"] = *privatekeypath
	}
	if *publickeypath != "" {
		flagConfig["PUBLIC_KEY_PATH"] = *publickeypath
	}
	if *hashkey != "" {
		flagConfig["HASH_KEY"] = *hashkey
	}
	if *name != "" {
		flagConfig["COMPONENT_NAME"] = *name
	}
	return flagConfig, nil
}
