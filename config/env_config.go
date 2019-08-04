package config

import (
	"github.com/spf13/viper"
)

// InitEnvConfig 从环境变量获得的配置内容初始化配置
func InitEnvConfig(ConfigViper *viper.Viper) {
	EnvConfigViper := viper.New()
	EnvConfigViper.SetEnvPrefix("jwtsigner") // will be uppercased automatically
	EnvConfigViper.BindEnv("address")
	EnvConfigViper.BindEnv("private_key_path")
	EnvConfigViper.BindEnv("public_key_path")
	EnvConfigViper.BindEnv("sign_method")
	EnvConfigViper.BindEnv("iss")
	if EnvConfigViper.Get("address") != nil {
		ConfigViper.Set("Address", EnvConfigViper.Get("address"))
	}
	if EnvConfigViper.Get("private_key_path") != nil {
		ConfigViper.Set("PrivateKeyPath", EnvConfigViper.Get("private_key_path"))
	}
	if EnvConfigViper.Get("public_key_path") != nil {
		ConfigViper.Set("PublicKeyPath", EnvConfigViper.Get("public_key_path"))
	}
	if EnvConfigViper.Get("sign_method") != nil {
		ConfigViper.Set("SignMethod", EnvConfigViper.Get("sign_method"))
	}
	if EnvConfigViper.Get("iss") != nil {
		ConfigViper.Set("Iss", EnvConfigViper.Get("iss"))
	}
}
