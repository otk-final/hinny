package module

import (
	"github.com/go-xorm/xorm"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

var DB *xorm.Engine

/**

	初始化
	driverName:mysql
	dataSourceName:mysql name:password@(ip:port)/xxx?charset=utf8
 */

func Install(driverName string, dataSourceName string) {
	/*数据库支持*/

	engine, err := xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		log.Fatal("initializing db :", err.Error())
	}
	//最大连接数
	engine.SetMaxIdleConns(10)

	err = engine.Ping()
	if err != nil {
		log.Fatal("ping db ping:", err.Error())
	}

	//暴露
	DB = engine
}


