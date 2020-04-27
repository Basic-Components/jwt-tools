package script

import (
	"github.com/xeipuuv/gojsonschema"
)

// ConfigType 配置类型
type ConfigType struct {
	Address        string `json:"ADDRESS"`
	PrivateKeyPath string `json:"PRIVATE_KEY_PATH"`
	PublicKeyPath  string `json:"PUBLIC_KEY_PATH"`
	Hashkey        string `json:"HASH_KEY"`
	ComponentName  string `json:"COMPONENT_NAME"`
	LogLevel       string `json:"LOG_LEVEL"`
	RegistEtcdURLS string `json:"REGIST_ETCD_URLS"`
	RegistVersion  string `json:"REGIST_VERSION"`
	RegistAddress  string `json:"REGIST_ADDRESS"`
}

//DefaultConfig 默认配置
var DefaultConfig = map[string]interface{}{
	"ADDRESS":          "0.0.0.0:5000",
	"PRIVATE_KEY_PATH": "autogen_rsa.pem",
	"PUBLIC_KEY_PATH":  "autogen_rsa_pub.pem",
	"HASH_KEY":         "secret can not guess",
	"COMPONENT_NAME":   "jwt-center",
	"LOG_LEVEL":        "DEBUG",
	"REGIST_ETCD_URLS": "",
	"REGIST_VERSION":   "",
	"REGIST_ADDRESS":   "",
}

//Schema 默认的配置样式
const Schema = `{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "title": "jwt-center",
    "description": "jwt签发中心",
    "type": "object",
    "required": [
        "ADDRESS",
        "PRIVATE_KEY_PATH",
        "PUBLIC_KEY_PATH",
        "HASH_KEY",
        "COMPONENT_NAME",
		"LOG_LEVEL",
		"REGIST_ETCD_URLS",
    ],
    "properties": {
        "ADDRESS": {"type": "string"},
        "PRIVATE_KEY_PATH": { "type": "string" },
        "PUBLIC_KEY_PATH": { "type": "string" },
        "LOG_LEVEL": { "type": "string", "enum": ["DEBUG", "INFO", "WARN", "ERROR"] },
        "HASH_KEY": { "type": "string" },
		"COMPONENT_NAME": { "type": "string" },
		"REGIST_ETCD_URLS": { "type": "string" },
		"REGIST_VERSION": { "type": "string" },
		"REGIST_ADDRESS": { "type": "string" }
    }
}`

//VerifyConfig 验证config是否符合要求
func VerifyConfig(conf ConfigType) (bool, *gojsonschema.Result) {
	configLoader := gojsonschema.NewGoLoader(conf)
	schemaLoader := gojsonschema.NewStringLoader(Schema)
	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		return false, result
	}
	if result.Valid() {
		return true, result
	}
	return false, result

}
