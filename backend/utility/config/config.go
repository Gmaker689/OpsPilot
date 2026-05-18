package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	once      sync.Once
	initErr   error
)

func Init() error {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("manifest/config")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			initErr = fmt.Errorf("failed to read config: %w", err)
			return
		}
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
	})
	return initErr
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}
