package conf

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"learnerai/golib/app/logger"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type GlobalAppConfig struct {
	App    App    `yaml:"app"`
	Server Server `yaml:"server"`
}

type App struct {
	RuntimeRootPath  string `yaml:"RuntimeRootPath"`
	LogSavePath      string `yaml:"LogSavePath"`
	LogSaveName      string `yaml:"LogSaveName"`
	LogFileExt       string `yaml:"LogFileExt"`
	LogLevel         string `yaml:"LogLevel"`
	TimeFormat       string `yaml:"TimeFormat"`
	ExpireTime       int64  `yaml:"ExpireTime"`
	AbsoluteAppPath  string `yaml:"AbsoluteAppPath"`
	InterfaceTimeOut int64  `yaml:"InterfaceTimeOut"`
}

type Server struct {
	RunMode         string        `yaml:"RunMode"`
	HttpPort        int32         `yaml:"HttpPort"`
	ConfigPath      string        `yaml:"ConfigPath"`
	Env             string        `yaml:"Env"`
	ReadTimeout     time.Duration `yaml:"ReadTimeout"`
	WriteTimeout    time.Duration `yaml:"WriteTimeout"`
	ShutDownTimeout time.Duration `yaml:"ShutDownTimeout"`
}

var AppConfig = &GlobalAppConfig{}

var AppSetting = &App{}
var ServerSetting = &Server{}

var ViperAppConfig *viper.Viper

var MapDeviceId2Flag = make(map[string]bool, 0)
var MapReidDeviceId2Flag = make(map[string]bool, 0)
var MapHour2Flag = make(map[string]bool, 0)

func LoadApiYaml() {
	ViperAppConfig.AddConfigPath("./conf/")
	ViperAppConfig.SetConfigName("api")
	ViperAppConfig.SetConfigType("yaml")
	err := ViperAppConfig.ReadInConfig()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	err = ViperAppConfig.Unmarshal(AppConfig)
	if err != nil {
		fmt.Println(" Unmarshal error +v%", err)
		os.Exit(1)
	}

	AppSetting = &AppConfig.App
	ServerSetting = &AppConfig.Server

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
}

func InitViper() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("initViper down [stack=%+v] [recover=%+v]", string(debug.Stack()), r)
		}
	}()
	ViperAppConfig = viper.New()
	LoadApiYaml()
}

var GlobalEdgeSetting = &EdgeServerLocalConfig{}

func GetGlobalConf(ctx *gin.Context) {
	err := GetConfig(ServerSetting.ConfigPath, GlobalEdgeSetting)
	if err != nil {
		logger.Errorf(ctx, "GetConfig error err=[%+v]", err)
		os.Exit(1)
	}
}
