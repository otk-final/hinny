package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	"github.com/sony/sonyflake"
	"log"
	"time"
)

/**
	工作空间表
 */
type Workspace struct {
	Kid         uint64 `json:"kid"         xorm:"pk bigint(20) notnull 'kid'"`
	Application string `json:"application" xorm:"varchar(64)   notnull 'application'"`
	WsName      string `json:"wsName"      xorm:"varchar(64)   notnull 'ws_name'"`
	WsKey       string `json:"wsKey"       xorm:"varchar(64)   notnull 'ws_key'"`
	ApiUrl      string `json:"apiUrl"      xorm:"varchar(256)  notnull 'api_url'"`
}

/**
	案例模板
 */
type CaseTemplate struct {
	Kid         uint64    `json:"kid"           xorm:"pk bigint(20) 			'kid'"`
	Application string    `json:"application"   xorm:"varchar(255) notnull 		'application'"`
	Module      string    `json:"module"        xorm:"varchar(255) 				'module'"`
	CaseType    string    `json:"caseType"      xorm:"varchar(255) 				'case_type'"`
	CaseName    string    `json:"caseName"      xorm:"varchar(255) 				'case_name'"`
	Description string    `json:"description"   xorm:"text 						'description'"`
	Path        string    `json:"path"          xorm:"varchar(256) notnull 		'path'"`
	MetaRequest string    `json:"metaRequest"   xorm:"text 						'request'"`
	ScriptType  string    `json:"scriptType"    xorm:"varchar(32) 				'script_type'"`
	Script      string    `json:"script"        xorm:"text 						'script'"`
	CreateTime  time.Time `json:"createTime"    xorm:"datetime 					'create_time'"`
}

/**
	案例日志
 */
type CaseLog struct {
	Kid          uint64    `json:"kid"          xorm:"pk bigint(20) 		    'kid'"`
	WsKId        uint64    `json:"wsKid"        xorm:"bigint(20) 				'ws_kid'"`
	CaseKid      uint64    `json:"caseKid"      xorm:"bigint(20) 				'case_kid'"`
	PathIdentity string    `json:"pathIdentity" xorm:"varchar(64) 				'path_identity'"`
	Path         string    `json:"path"         xorm:"bigint(20) 				'path'"`
	MetaRequest  string    `json:"metaRequest"  xorm:"text 						'request'"`
	MetaResponse string    `json:"metaResponse" xorm:"text 						'response'"`
	Script       string    `json:"script"       xorm:"text 						'script'"`
	ScriptType   string    `json:"scriptType"   xorm:"text 						'script_type'"`
	MetaResult   string    `json:"metaResult"   xorm:"text 						'result'"`
	Status       int       `json:"status"       xorm:"tinyint(3) notnull 		'status'"`
	Curl         string    `json:"curl"         xorm:"varchar(128) notnull 		'curl'"`
	CreateTime   time.Time `json:"createTime"   xorm:"datetime 					'create_time'"`
}

var Conn *xorm.Engine
var idGeneral *sonyflake.Sonyflake
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
	Conn = engine
}

func InstallIDGeneral(startTime time.Time, machineID uint16) {
	st := &sonyflake.Settings{
		StartTime: startTime,
		MachineID: func() (uint16, error) {
			return machineID, nil
		},
	}
	idGeneral = sonyflake.NewSonyflake(*st)
}

func GetNextKid() (uint64, error) {
	return idGeneral.NextID()
}
