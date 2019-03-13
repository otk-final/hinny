package global

import (
	"github.com/spf13/viper"
	"fmt"
)

var Conf *viper.Viper

func init()  {
	initConfiguration()
}


func initConfiguration() {
	Conf = viper.New()
	Conf.AddConfigPath("./config/")
	Conf.SetConfigName("conf")
	Conf.SetConfigType("toml")

	err := Conf.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("加载配置文件异常:%s", err))
	}
}