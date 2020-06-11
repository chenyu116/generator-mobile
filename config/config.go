package config

import (
	"github.com/spf13/viper"
	"log"
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

type ServeConfig struct {
	HostPort string `mapstructure:"hostPort"`
}

type CommonConfig struct {
	StaticPath string `mapstructure:"staticPath"`
}

type DBConfig struct {
	HostPort string `mapstructure:"hostPort"`
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
type OssConfig struct {
	EndPoint        string `mapstructure:"endPoint"`
	AccessKeyID     string `mapstructure:"accessKeyID"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	BucketName      string `mapstructure:"bucketName"`
}

type Config struct {
	Serve    ServeConfig    `mapstructure:"serve"`
	DbServer DBConfig       `mapstructure:"dbServer"`
	Common   CommonConfig   `mapstructure:"common"`
	Rabbitmq RabbitmqConfig `mapstructure:"rabbitmq"`
	Oss      OssConfig      `mapstructure:"oss"`
}

func SetConfigPath(path string) {
	viper.SetConfigFile(path)
}

func ReadConfig() {
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err.Error())
	}

	err = viper.Unmarshal(&c)

	if err != nil {
		log.Fatal(err)
	}
}
