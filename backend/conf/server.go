package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

// 存储于数据库的配置
type MongoConfig struct {
	User     string `yaml:"User" json:"user" bson:"user"`
	Password string `yaml:"Password" json:"password" bson:"password"`
	Host     string `yaml:"Host" json:"host" bson:"host"`
}

type EdgeServerLocalConfig struct {
	Mongo  MongoConfig  `yaml:"mongo" json:"mongo" bson:"mongo"`
	Redis  RedisConfig  `yaml:"redis" json:"redis" bson:"redis"`
	Milvus MilvusConfig `yaml:"milvus" json:"milvus" bson:"milvus"`
}

type RedisConfig struct {
	Password string `yaml:"Password"`
	Host     string `yaml:"Host"`
}

type MilvusConfig struct {
	User     string `yaml:"User" json:"user" bson:"user"`
	Password string `yaml:"Password" json:"password" bson:"password"`
	Host     string `yaml:"Host" json:"host" bson:"host"`
}

// Deprecated: use New() instead if you want runtime auto reload.
func GetConfig(path string, globalConfig *EdgeServerLocalConfig) (err error) {
	ViperServerConfig := viper.New()
	ViperServerConfig.AddConfigPath(path)
	ViperServerConfig.SetConfigName("server")
	ViperServerConfig.SetConfigType("yaml")
	err = ViperServerConfig.ReadInConfig()

	fmt.Printf("GetConfig request path=[%s]", path)
	if err != nil {
		fmt.Printf("GetConfig error err=[%+v]", err)
		return
	}

	err = ViperServerConfig.Unmarshal(globalConfig)
	if err != nil {
		fmt.Printf("GetConfig error err=[%+v]", err)
		return
	}
	return
}

func (c *EdgeServerLocalConfig) UnmarshalFrom(v *viper.Viper) error {
	return v.Unmarshal(c)
}

func (c *EdgeServerLocalConfig) OnReloadError(err error) {
	fmt.Println("reload server config error:", err)
}

type Config = SafeConfig[EdgeServerLocalConfig]

// FIXME: precedence depends on undocumented behavior:
// config paths are searched in order of added.
// but it is unlikely to change.
func New(path string, additionalPaths ...string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(path)
	for _, path := range additionalPaths {
		v.AddConfigPath(path)
	}
	v.SetConfigName("server")
	v.SetConfigType("yaml")

	// can set default value to v here: v.SetDefault(k, v)

	return Watch[EdgeServerLocalConfig](v)
}
