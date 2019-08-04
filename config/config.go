package config

import (
	errs "github.com/Basic-Components/jwtrpc/errs"
	logger "github.com/Basic-Components/jwtrpc/logger"
	signals "github.com/Basic-Components/jwtrpc/signals"

	"github.com/spf13/viper"
)

// ConfigType 配置类型
type ConfigType struct {
	Address        string
	PrivateKeyPath string
	PublicKeyPath  string
	SignMethod     string
	Iss            string
}

// Init 根据不同的途径构造配置
func Init() (ConfigType, error) {
	var Config ConfigType
	ConfigViper := viper.New()
	InitDefaultConfig(ConfigViper)

	InitFileConfig(ConfigViper)
	InitEnvConfig(ConfigViper)
	ok := InitFlagConfig(ConfigViper)
	if ok {
		err := ConfigViper.Unmarshal(&Config)
		if err != nil {
			logger.Logger.Error("unable to decode into struct, %v", err)
			return Config, errs.ConfigDecodeError

		}
		if VerifyConfig(Config) {
			return Config, nil
		}
		return Config, errs.ConfigVerifyError
	}
	return Config, signals.GenkeySignal
}
