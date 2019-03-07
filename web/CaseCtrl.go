package web

import (
	"net/http"
	"otk-final/hinny/module/db"
	"strings"
	"otk-final/hinny/service"
	"io/ioutil"
	"encoding/json"
)

type CaseModuleGroup struct {
	Module string   `json:"module"`
	Groups []string `json:"groups"`
}

func GetCaseModuleCroups(response http.ResponseWriter, request *http.Request) {
	db.Conn.Cols("").GroupBy("module,group").Find(db.CaseTemplate{})

	rows, err := db.Conn.Query("select module,group_concat(group_name) groups from hinny_case_template group by module")
	if err != nil {
		panic(err)
	}

	out := make([]CaseModuleGroup, 0)
	for _, row := range rows {
		item := &CaseModuleGroup{
			Module: string(row["module"]),
		}
		groups, ok := row["groups"]
		if !ok {
			item.Groups = []string{}
		} else {
			item.Groups = strings.Split(string(groups), ",")
		}
		out = append(out, *item)
	}
	view.JSON(response, 200, out)
}

/**
	案例执行
 */
func CaseExecute(response http.ResponseWriter, request *http.Request) {
	key := request.Header.Get("workspace")

	//获取
	body, _ := ioutil.ReadAll(request.Body)
	input := &service.CaseInput{}
	err := json.Unmarshal(body, input)
	if err != nil {
		panic(err)
	}

	//查询唯一请求
	path := service.GetPathPrimary(key, input.PrimaryId)
	if path == nil {
		view.JSON(response, 500, "未查询到指定接口")
		return
	}

	//查询工作空间
	ws := &db.Workspace{}
	ok, err := db.Conn.Where("ws_key=?", key).Get(ws)
	if !ok || err != nil {
		panic(err)
	}

	//格式化metaRequest内容
	//reqCtx, _ := json.Marshal(input.Request)
	/**
		组件参数
		记录模板
		异步发送请求，通过chan获取响应
	 */
	//temp := &db.CaseTemplate{
	//	Application: "",
	//	Case:        "",
	//	Module:      "",
	//	Group:       "",
	//	Description: "",
	//	ServiceKey:  path.Tag.Name,
	//	PathKey:     path.Path,
	//	MetaRequest: string(reqCtx),
	//	ScriptType:  input.Valid.ScriptType,
	//	Script:      input.Valid.Script,
	//	CreateTime:  time.Now(),
	//}
	////新增数据库
	//db.Conn.Insert(temp)
	//执行
	out, _ := service.Execute(ws.ApiUrl, path, input)
	//响应
	view.JSON(response, 200, out)
}
