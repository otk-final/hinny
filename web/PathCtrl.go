package web

import (
	"net/http"
	"otk-final/hinny/service"
	"strings"
	"otk-final/hinny/service/swagger"
	"otk-final/hinny/module"
	"encoding/json"
	"github.com/kataras/iris/core/errors"
)

func init() {

}

func GetServices(response http.ResponseWriter, request *http.Request) {
	//当前工作空间
	ws := GetWorkspaceFromHeader(request)
	//查询
	out, err := service.ApiTagList(&swagger.SwaggerHandler{}, ws)
	if err != nil {
		view.JSON(response, 200, err)
		return
	}

	view.JSON(response, 200, out)
}

func GetPaths(response http.ResponseWriter, request *http.Request) {

	//当前工作空间
	ws := GetWorkspaceFromHeader(request)
	findType := request.URL.Query().Get("type")
	findValue := request.URL.Query().Get("typeValue")

	//查询
	out, err := service.ApiPathList(&swagger.SwaggerHandler{}, ws)
	if err != nil {
		view.JSON(response, 500, err)
		return
	}

	multipleContains := func(target string, srcArray ... string) bool {
		for _, src := range srcArray {
			if strings.Contains(src, target) {
				return true
			}
		}
		return false
	}

	outs := make([]module.ApiPath, 0)

	for _, path := range out {
		match := matchPath(findType, findValue, &path, multipleContains)
		if match {
			outs = append(outs, path)
		}
	}
	view.JSON(response, 200, outs)
}

func GetPath(key uint64, primaryId string) (*module.MetaOut, error) {

	//查询唯一
	path := service.GetPathPrimary(key, primaryId)
	if path == nil {
		return nil, errors.New("未查询到指定接口")
	}

	parameters := path.Parameters
	//过滤出，uri,query,header,body参数

	kindOf := func(parameters []interface{}, in string) []interface{} {
		kindArray := make([]interface{}, 0)
		for _, p := range parameters {
			item := p.(map[string]interface{})
			kind, ok := item["in"]
			if ok && kind.(string) == in {
				kindArray = append(kindArray, item)
			}
		}
		return kindArray
	}

	/**
		1，生成scheme文件
		2，补齐相关参数
	 */
	req := &module.MetaRequest{
		Header: kindOf(parameters, "header"),
		Uri:    kindOf(parameters, "path"),
		Query:  kindOf(parameters, "query"),
	}

	//获取请求body
	bodyParams := kindOf(parameters, "body")
	if bodyParams != nil && len(bodyParams) > 0 {
		req.Body = getReqBodyDefineJson(key, bodyParams[0])
	} else {
		req.Body = "{}"
	}

	//只取值200的返回信息
	resp := &module.MetaResponse{
		Body: getRespBodyDefineJson(key, *path, "200"),
	}

	path.Parameters = nil
	path.Responses = nil

	//获取标准示例验证脚本
	return &module.MetaOut{
		Path:     path,
		Request:  req,
		Response: resp,
	}, nil
}

func GetPrimaryPath(response http.ResponseWriter, request *http.Request) {

	//当前工作空间
	ws := GetWorkspaceFromHeader(request)

	values := request.URL.Query()
	primaryId := values["primaryId"][0]

	out, err := GetPath(ws.Kid, primaryId)
	if err != nil {
		view.JSON(response, 500, err.Error())
		return
	}

	//获取当前工作空间中的默认脚本
	out.Valid = &module.MetaValid{
		Script:     ws.Script,
		ScriptType: "javascript",
	}

	view.JSON(response, 200, out)
}

func matchPath(findType string, findValue string, path *module.ApiPath,
	matchFunc func(target string, srcArray ... string) bool) bool {
	if findType == "service" {
		return matchFunc(findValue, path.Tag.Name, path.Tag.Description, path.Description)
	}
	if findType == "path" {
		return matchFunc(findValue, path.Path)
	}
	return false
}

func getReqBodyDefineJson(key uint64, bodyParam interface{}) string {
	bodyMap := bodyParam.(map[string]interface{})
	schema, ok := bodyMap["schema"]
	if !ok {
		return "{}"
	}

	schemaMap := schema.(map[string]interface{})
	var ref interface{}
	var bodyOut interface{}
	//判断类型
	schemaType := schemaMap["type"]
	if "array" == schemaType {
		ref, ok = schemaMap["items"].(map[string]interface{})["$ref"]
		if !ok {
			return "[]"
		}
		bodyOut = service.GetDefinitionArray(key, ref.(string))
	} else {
		ref, ok = schemaMap["$ref"]
		if !ok {
			return "{}"
		}
		bodyOut = service.GetDefinitionMap(key, ref.(string))
	}

	json, err := json.Marshal(bodyOut)
	if err != nil {
		panic(err)
	}
	return string(json)
}

func getRespBodyDefineJson(key uint64, path module.ApiPath, code string) string {
	codeMap := path.Responses[code].(map[string]interface{})
	return getReqBodyDefineJson(key, codeMap)
}
