package script

import (
	"github.com/spf13/viper"
)

// InitEnvConfig 从环境变量获得的配置内容初始化配置
func InitEnvConfig() (map[string]interface{}, error) {
	var envConfig = map[string]interface{}{}
	EnvConfigViper := viper.New()
	EnvConfigViper.SetEnvPrefix("jwt-center") // will be uppercased automatically
	EnvConfigViper.BindEnv("ADDRESS")
	EnvConfigViper.BindEnv("PRIVATE_KEY_PATH")
	EnvConfigViper.BindEnv("PUBLIC_KEY_PATH")
	EnvConfigViper.BindEnv("HASH_KEY")
	EnvConfigViper.BindEnv("COMPONENT_NAME")
	EnvConfigViper.BindEnv("LOG_LEVEL")
	err := EnvConfigViper.Unmarshal(&envConfig)
	return envConfig, err
}
