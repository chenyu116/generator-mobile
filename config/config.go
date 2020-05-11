package config

import (
	"github.com/chenyu116/generator-mobile/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var c Config

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/mnt")
	viper.AddConfigPath(".")
}

func GetConfig() Config {
	return c
}

func GetRabbitmqConfig() RabbitmqConfig {
	return c.Rabbitmq
}

func GetCommonConfig() CommonConfig {
	return c.Common
}

func Get(key string) interface{} {
	return viper.Get(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func IsSet(key string) bool {
	return viper.IsSet(key)
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
}

type WebsocketConfig struct {
	HostPort string       `mapstructure:"hostPort"`
	Tls      websocketTls `mapstructure:"tls"`
}
type websocketTls struct {
	Enable bool   `mapstructure:"enable"`
	Cert   string `mapstructure:"cert"`
	Key    string `mapstructure:"key"`
}

type ApiConfig struct {
	HostPort string `mapstructure:"hostPort"`
}

type CommonConfig struct {
	ConnectTimeout uint32   `mapstructure:"connectTimeout"`
	Version        string   `mapstructure:"version"`
	AllowHosts     []string `mapstructure:"allowHosts"`
	AllowOrigin    []string `mapstructure:"allowOrigin"`
}

type RabbitmqConfig struct {
	HostPort    string              `mapstructure:"hostPort"`
	Username    string              `mapstructure:"username"`
	Password    string              `mapstructure:"password"`
	Prefetch    int                 `mapstructure:"prefetch"`
	VHost       string              `mapstructure:"vHost"`
	Exchanges   []map[string]string `mapstructure:"exchanges"`
	QueuePrefix string              `mapstructure:"queuePrefix"`
}
type Config struct {
	Websocket WebsocketConfig `mapstructure:"websocket"`
	ApiServer ApiConfig       `mapstructure:"apiServer"`
	Common    CommonConfig    `mapstructure:"common"`
	Rabbitmq  RabbitmqConfig  `mapstructure:"rabbitmq"`
}

func SetConfigPath(path string) {
	viper.SetConfigFile(path)
}

func ReadConfig() {
	err := viper.ReadInConfig()

	if err != nil {
		logger.ZapLogger.Fatal("viper.ReadInConfig", zap.Error(err))
	}
	err = viper.Unmarshal(&c)

	if err != nil {
		logger.ZapLogger.Fatal("viper.Unmarshal", zap.Error(err))
	}
}
