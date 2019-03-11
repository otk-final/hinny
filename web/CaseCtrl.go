package web

import (
	"net/http"
	"otk-final/hinny/module/db"
	"strings"
	"otk-final/hinny/service"
	"io/ioutil"
	"encoding/json"
	"time"
	"fmt"
	"net/url"
	"github.com/gorilla/mux"
	"otk-final/hinny/module"
)

type CaseModuleGroup struct {
	Module    string   `json:"module"`
	CaseTypes []string `json:"caseTypes"`
}

type CaseArchiveInput struct {
	LogKid      uint64 `json:"logKid"`
	Module      string `json:"module"`
	CaseType    string `json:"caseType"`
	CaseName    string `json:"caseName"`
	Description string `json:"description"`
}

func GetCaseModules(response http.ResponseWriter, request *http.Request) {

	out := make([]CaseModuleGroup, 0)

	application := request.Header.Get("application")
	if application == "" {
		view.JSON(response, 200, out)
		return
	}

	rows, err := db.Conn.Select("module,group_concat(case_type) caseTypes").
		Table("hinny_case_template").
		Where("application = ?", application).
		GroupBy("module").Query()

	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		item := &CaseModuleGroup{
			Module: string(row["module"]),
		}
		types, ok := row["caseTypes"]
		if !ok {
			item.CaseTypes = []string{}
		} else {
			item.CaseTypes = strings.Split(string(types), ",")
		}
		out = append(out, *item)
	}
	view.JSON(response, 200, out)
}

/**
	案例执行
 */
func CaseExecute(response http.ResponseWriter, request *http.Request) {
	ws := GetWorkspaceFromHeader(request)

	//获取
	body, _ := ioutil.ReadAll(request.Body)
	input := &service.CaseInput{}
	err := json.Unmarshal(body, input)
	if err != nil {
		panic(err)
	}

	//查询唯一请求
	path := service.GetPathPrimary(ws.Kid, input.PrimaryId)
	if path == nil {
		view.JSON(response, 500, "未查询到指定接口")
		return
	}

	//执行
	out, _ := service.Execute(ws, path, input)
	//响应
	view.JSON(response, 200, out)
}

/**
	案例保存模板
 */
func CaseSave(response http.ResponseWriter, request *http.Request) {
	application := request.Header.Get("application")

	//获取
	body, _ := ioutil.ReadAll(request.Body)
	input := &CaseArchiveInput{}
	err := json.Unmarshal(body, input)
	if err != nil {
		panic(err)
	}

	//查询原始记录，保存至模板信息表中
	log := &db.CaseLog{}
	ok, err := db.Conn.ID(input.LogKid).Get(log)
	if !ok || err != nil {
		panic(nil)
	}

	tempKid, _ := db.GetNextKid()

	temp := &db.CaseTemplate{
		Kid:         tempKid,
		Application: application,
		CaseName:    input.CaseName,
		Module:      input.Module,
		CaseType:    input.CaseType,
		Description: input.Description,
		Path:        log.Path,
		MetaRequest: log.MetaRequest,
		ScriptType:  log.ScriptType,
		Script:      log.Script,
		CreateTime:  time.Now(),
	}

	//创建事务
	session := db.Conn.NewSession()
	defer session.Close()

	//事务提交或者回滚
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			session.Rollback()
			view.JSON(response, 200, "保存异常")
		}
	}()

	session.Begin()
	//新增模板
	count, err := session.Insert(temp)
	if count != 1 && err != nil {
		panic(err)
	}

	//修改日志
	count, err = session.ID(log.Kid).Update(db.CaseLog{CaseKid: tempKid})
	if count != 1 && err != nil {
		panic(err)
	}
	session.Commit()

	view.JSON(response, 200, "ok")
}

/**
	获取模块列表
 */
func GetCases(response http.ResponseWriter, request *http.Request) {
	temps := make([]*db.CaseTemplate, 0)

	application := request.Header.Get("application")
	if application == "" {
		view.JSON(response, 200, temps)
		return
	}

	/**
		获取查询参数
	 */
	values, err := url.ParseQuery(request.URL.RawQuery)

	dynamicArgs := func() (string, []interface{}) {
		args := make([]interface{}, 0)

		//项目为必要条件
		sql := "application = ? "
		args = append(args, application)

		//模块
		if module := values.Get("module"); module != "" {
			sql += "and module = ? "
			args = append(args, module)
		}

		//类型
		if caseType := values.Get("caseType"); caseType != "" {
			sql += "and case_type = ? "
			args = append(args, caseType)
		}

		//其他
		if searchText := values.Get("searchText"); searchText != "" {
			sql += "and (case_name like ? or path like ?) "
			args = append(args, "%"+searchText+"%", "%"+searchText+"%")
		}

		return sql, args
	}

	/**
		查询
	 */
	sql, args := dynamicArgs()

	err = db.Conn.Cols("kid", "application", "case_type", "module", "case_name", "create_time").Where(sql, args...).Find(&temps)
	if err != nil {
		panic(err)
	}

	view.JSON(response, 200, temps)
}

func GetCaseLog(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	logKid := vars["kid"]

	//查询数据库记录,获取执行记录
	log := &db.CaseLog{}
	ok, err := db.Conn.ID(logKid).Get(log)
	if !ok || err != nil {
		panic(err)
	}

	//获取Path相关基本信息
	tpl, err := GetPath(log.WsKId, log.PathIdentity)
	if err != nil {
		view.JSON(response, 500, err.Error())
		return
	}

	logRequest := &module.MetaRequest{}
	err = json.Unmarshal([]byte(log.MetaRequest), logRequest)
	if err != nil {
		panic(err)
	}

	logResponse := &module.MetaResponse{}
	err = json.Unmarshal([]byte(log.MetaResponse), logResponse)
	if err != nil {
		panic(err)
	}

	logResults := make([]*module.MetaResult, 0)
	err = json.Unmarshal([]byte(log.MetaResult), &logResults)
	if err != nil {
		panic(err)
	}

	getVal := func(dbArray []interface{}, name string) (bool, interface{}) {
		for _, item := range dbArray {
			itemMap := item.(map[string]interface{})
			if itemMap["name"] == name {
				return true, itemMap["value"]
			}
		}
		return false, nil
	}

	//对request进行默认值设置
	renderValue := func(tplArray []interface{}, dbArray []interface{}) []interface{} {
		for _, tpl := range tplArray {
			tplMap := tpl.(map[string]interface{})
			exist, val := getVal(dbArray, tplMap["name"].(string))

			//不存在，将当前require改为false,前端不进行默认勾选
			tplMap["required"] = exist
			if exist {
				//存在设置值
				tplMap["value"] = val
			}
		}
		return tplArray
	}

	//重设request,response,valid
	tpl.Request = &module.MetaRequest{
		Header: renderValue(tpl.Request.Header, logRequest.Header),
		Uri:    renderValue(tpl.Request.Uri, logRequest.Uri),
		Query:  renderValue(tpl.Request.Query, logRequest.Query),
		Body:   logRequest.Body,
	}
	tpl.Response = logResponse
	tpl.Valid = &module.MetaValid{
		Script:     log.Script,
		ScriptType: log.ScriptType,
	}
	tpl.Result = logResults

	view.JSON(response, 200, tpl)
}

/**
	获取案例下，的日志记录
 */
func GetCaseLogs(response http.ResponseWriter, request *http.Request) {
	caseKid := request.URL.Query().Get("caseKid")

	/**
		查询数据库，创建时间倒叙
	 */
	out := make([]*db.CaseLog, 0)
	err := db.Conn.Cols("kid", "path", "status", "create_time").
		Where("case_kid=?", caseKid).
		Desc("create_time").Find(&out)
	if err != nil {
		panic(err)
	}

	view.JSON(response, 200, out)
}

func GetCaseTpl(response http.ResponseWriter, request *http.Request) {
	//获取Path相关基本信息
	out, err := GetPath(0, "")
	if err != nil {
		view.JSON(response, 500, err.Error())
		return
	}
	//TODO 封装相关验证信息

	//封装返回  MetaRequest/MetaValid
	view.JSON(response, 200, out)
}
