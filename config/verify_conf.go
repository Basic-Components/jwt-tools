package config

import (
	logger "github.com/Basic-Components/jwtrpc/logger"

	"github.com/xeipuuv/gojsonschema"
)

const schema = `{
	"description": "jwt grpc server config",
	"type": "object",
	"required": [ "Address","SignMethod","PrivateKeyPath","PublicKeyPath","Iss"],
	"additionalProperties": false,
	"properties": {
		"Address": {
			"type": "string",
			"description": "the host and port"
		},
		"PrivateKeyPath": {
			"type": "string",
			"description": "PrivateKey's Path"
		},
		"PublicKeyPath": {
			"type": "string",
			"description": "PublicKey's Path"
		},
		"SignMethod": {
			"type": "string",
			"description": "Sign jwt's Method"
		},
		"Iss": {
			"type": "string",
			"description": "Signer's name"
		}
	}
}`

// VerifyConfig 验证config是否符合要求
func VerifyConfig(conf ConfigType) bool {
	configLoader := gojsonschema.NewGoLoader(conf)
	schemaLoader := gojsonschema.NewStringLoader(schema)
	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		logger.Log.Error("Validate error: %s", err)
		return false
	} else {
		if result.Valid() {
			logger.Log.Debug("The config is valid")
			return true
		} else {
			logger.Logger.Info("The config is not valid. see errors :\n")
			for _, err := range result.Errors() {
				logger.Log.Error("- %s", err)
			}
			return false
		}
	}
}
