package config

import (
	"path"
	"strings"

	"github.com/Basic-Components/jwtrpc/logger"

	logrus "github.com/sirupsen/logrus"
	"github.com/small-tk/pathlib"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// InitFlagConfig 初始化命令行传入的参数到配置,返回值为false表示要执行创建秘钥否则为启动服务
func InitFlagConfig(ConfigViper *viper.Viper) bool {
	genkey := pflag.BoolP("genkey", "g", false, "创建rsa公私钥对")
	loglevel := pflag.StringP("loglevel", "l", "WARN", "创建rsa公私钥对")
	confPath := pflag.StringP("config", "c", "", "配置文件位置")
	address := pflag.StringP("address", "a", "", "要启动的服务器地址")
	privateKeyPath := pflag.StringP("private_key_path", "r", "", "私钥位置")
	publicKeyPath := pflag.StringP("public_key_path", "u", "", "公钥位置")
	signMethod := pflag.StringP("sign_method", "m", "", "签名方法")
	iss := pflag.StringP("iss", "i", "", "签名者")
	pflag.Parse()
	switch {
	case strings.ToUpper(*loglevel) == "TRACE":
		logger.Logger.SetLevel(logrus.TraceLevel)
	case strings.ToUpper(*loglevel) == "DEBUG":
		logger.Logger.SetLevel(logrus.DebugLevel)
	case strings.ToUpper(*loglevel) == "INFO":
		logger.Logger.SetLevel(logrus.InfoLevel)
	case strings.ToUpper(*loglevel) == "WARN":
		logger.Logger.SetLevel(logrus.WarnLevel)
	case strings.ToUpper(*loglevel) == "Error":
		logger.Logger.SetLevel(logrus.ErrorLevel)
	}
	if *genkey {
		return false
	}
	if *confPath != "" {
		p, err := pathlib.New(*confPath).Absolute()
		if err != nil {
			logger.Logger.Error("指定的配置文件获取绝对位置失败")
		} else {
			if p.Exists() && p.IsFile() {
				filenameWithSuffix := path.Base(*confPath)
				fileSuffix := path.Ext(filenameWithSuffix)
				fileName := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
				dir, err := p.Parent()
				if err != nil {
					logger.Logger.Error("指定的配置文件获取父文件夹位置失败")
				} else {
					filePaths := []string{dir.Path}
					SetFileConfig(ConfigViper, fileName, filePaths)
				}

			}
		}
	}
	if *address != "" {
		ConfigViper.Set("Address", *address)
	}
	if *privateKeyPath != "" {
		ConfigViper.Set("PrivateKeyPath", *privateKeyPath)
	}
	if *publicKeyPath != "" {
		ConfigViper.Set("PublicKeyPath", *publicKeyPath)
	}
	if *signMethod != "" {
		ConfigViper.Set("SignMethod", *signMethod)
	}
	if *iss != "" {
		ConfigViper.Set("Iss", *iss)
	}
	return true

}
