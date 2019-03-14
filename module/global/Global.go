package global

import (
	"github.com/spf13/viper"
	"fmt"
	"github.com/spf13/pflag"
	"flag"
)

var Conf *viper.Viper

func init() {
	initConfiguration()
}

func initConfiguration() {
	Conf = viper.New()

	//环境变量
	Conf.AutomaticEnv()

	//命令行 **启动参数需要提前定义？？？这特么是坑么**
	//pflag.String("application", "inn", "")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	Conf.BindPFlags(pflag.CommandLine)

	//配置文件
	Conf.AddConfigPath("./config/")
	Conf.SetConfigName("conf")
	Conf.SetConfigType("toml")



	err := Conf.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("加载配置文件异常:%s", err))
	}
}
