package script

import (
	"github.com/spf13/viper"
)

// SetFileConfig 从指定的配置文件中读取配置
func SetFileConfig(fileName string, filePaths []string) (map[string]interface{}, error) {
	var fileConfig = map[string]interface{}{}
	FileConfigViper := viper.New()
	FileConfigViper.SetConfigName(fileName)
	for _, filePath := range filePaths {
		FileConfigViper.AddConfigPath(filePath)
	}
	err := FileConfigViper.ReadInConfig()
	if err != nil {
		return fileConfig, err
	}

	err = FileConfigViper.Unmarshal(&fileConfig)
	return fileConfig, err
}

// InitFileConfig 从默认的配置文件位置读取配置
func InitFileConfig() (map[string]interface{}, error) {
	fileName := "config"
	filePaths := []string{"/etc/jwtcenter/", "$HOME/.jwtcenter", "."}
	fileConfig, err := SetFileConfig(fileName, filePaths)
	return fileConfig, err
}
