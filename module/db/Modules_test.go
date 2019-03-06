package db

import (
	"testing"
	"fmt"
	"encoding/json"
	"time"
	"otk-final/hinny/module"
)

func init() {
	Install("mysql",
		"dev62:dev62.123456@(192.168.30.62:3306)/platform_behavior?charset=utf8")
}

func TestWorkspace(t *testing.T) {
	id, err := Conn.Insert(&Workspace{
		Application: "lovelorn",
		WsName:      "个人自测平台",
		WsKey:       "huangxy-local",
		ApiUrl:      "http://192.168.30.61:18080",
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}

func TestCaseTemplate(t *testing.T) {

	head := make(map[string]string)
	head["tenantId"] = "lovelorn"
	head["token"] = "2668caf4-a889-4a40-8f69-0939c1235e85"
	head["userId"] = "277204731355136"

	uri := make(map[string]string)
	uri["version"] = "v1"

	query := make(map[string]interface{})
	query["pageNo"] = "1"
	query["pageSize"] = "20"

	req := &module.MetaRequest{
		Header: nil,
		Uri:    nil,
		Query:  nil,
		Body:   "",
	}

	reqCtx, _ := json.Marshal(req)
	fmt.Println(req)

	script := "function(a,b){" +
		"}"

	temp := &CaseTemplate{
		Application: "lovelorn",
		Case:        "测试案例名称",
		Module:      "模块",
		Group:       "分组",
		Description: "备注",
		ServiceKey:  "rest-controller",
		PathKey:     "/{version}/pv/wish-case-items/action/list-admin",
		MetaRequest: string(reqCtx),
		ScriptType:  "javascript",
		Script:      script,
		CreateTime:  time.Now(),
	}

	id, err := Conn.Insert(temp)

	if err != nil {
		panic(err)
	}

	fmt.Println(id)
}

func TestCaseLog(t *testing.T) {

	r := &module.MetaResult{
		Rule: "规则说明",
		Ok:   true,
		Msg:  "验证通过",
	}
	resultCtx, _ := json.Marshal([]*module.MetaResult{r})

	temp := &CaseTemplate{}
	ok, err := Conn.ID(2).Get(temp)

	if !ok && err != nil {
		panic(err)
	}

	log := &CaseLog{
		WsId:        1,
		CaseId:      temp.Id,
		MetaRequest: temp.MetaRequest,
		MetaResult:  string(resultCtx),
		Status:      1,
		Curl:        "curl++++",
		CreateTime:  time.Now(),
	}

	id, err := Conn.Insert(log)
	if err != nil {
		panic(err)
	}
	fmt.Println(id)

}
