package config

import (
	"github.com/spf13/viper"
)

// InitDefaultConfig 配置配置项的默认值
func InitDefaultConfig(ConfigViper *viper.Viper) {
	ConfigViper.SetDefault("Address", "0.0.0.0:5000")
	ConfigViper.SetDefault("PrivateKeyPath", "autogen_rsa.pem")
	ConfigViper.SetDefault("PublicKeyPath", "autogen_rsa_pub.pem")
	ConfigViper.SetDefault("SignMethod", "RS256")
	ConfigViper.SetDefault("Iss", "jwt-signer")
}
