package db

import (
	"github.com/go-xorm/xorm"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"time"
)

/**
	工作空间表
 */
type Workspace struct {
	Id          int64  `xorm:"bigint(20) notnull 	'id'"`
	Application string `xorm:"varchar(64) notnull 	'application'"`
	WsName      string `xorm:"varchar(64) notnull 	'ws_name'"`
	WsKey       string `xorm:"varchar(64) notnull 	'ws_key'"`
	ApiUrl      string `xorm:"varchar(256) notnull 	'api_url'"`
}

/**
	案例模板
 */
type CaseTemplate struct {
	Id          int64     `xorm:"pk bigint(20) notnull 		'id'"`
	Application string    `xorm:"varchar(255) notnull 		'application'"`
	Case        string    `xorm:"varchar(255) 				'case'"`
	Module      string    `xorm:"varchar(255) 				'module'"`
	Group       string    `xorm:"varchar(255) 				'group'"`
	Description string    `xorm:"text 						'description'"`
	ServiceKey  string    `xorm:"varchar(64) notnull 		'service_key'"`
	PathKey     string    `xorm:"varchar(128) notnull 		'path_key'"`
	MetaRequest string    `xorm:"text 						'request'"`
	ScriptType  string    `xorm:"varchar(32) 				'script_type'"`
	Script      string    `xorm:"text 						'script'"`
	CreateTime  time.Time `xorm:"datetime 					'create_time'"`
}

/**
	案例日志
 */
type CaseLog struct {
	Id           int64     `xorm:"bigint(20) 				'id'"`
	WsId         int64     `xorm:"bigint(20) 				'ws_id'"`
	CaseId       int64     `xorm:"bigint(20) 				'case_id'"`
	MetaRequest  string    `xorm:"text 						'request'"`
	MetaResponse string    `xorm:"text 						'response'"`
	MetaResult   string    `xorm:"text 						'result'"`
	Status       int       `xorm:"tinyint(3) notnull 		'status'"`
	Curl         string    `xorm:"varchar(128) notnull 		'curl'"`
	CreateTime   time.Time `xorm:"datetime 					'create_time'"`
}

/**
	元数据-请求
 */
type MetaRequest struct {
	Version string                 `json:"version"`
	Header  map[string]string      `json:"header"`
	Uri     map[string]string      `json:"uri"`
	Query   map[string]interface{} `json:"query"`
	Body    interface{}            `json:"body"`
}

/**
	元数据-响应
 */
type MetaResponse struct {
	Version string            `json:"version"`
	Header  map[string]string `json:"header"`
	Body    interface{}       `json:"body"`
	Code    int8              `json:"code"`
}

/**
	元数据-验证结果
 */
type MetaResult struct {
	Version string `json:"version"`
	Rule    string `json:"rule"`
	Msg     string `json:"msg"`
	Ok      bool   `json:"ok"`
}

var Session *xorm.Engine
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
	//统一去除前缀前缀
	mapper := core.NewPrefixMapper(core.SnakeMapper{}, "hinny_")
	engine.ShowSQL(true)
	engine.SetTableMapper(mapper)
	//暴露
	Session = engine
}
