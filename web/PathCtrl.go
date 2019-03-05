package web

import (
	"net/http"
	"otk-final/hinny/service"
	"strings"
	"otk-final/hinny/service/swagger"
	"otk-final/hinny/module"
	"encoding/json"
)

func init() {

}

func GetServices(response http.ResponseWriter, request *http.Request) {
	//当前工作空间
	key := request.Header.Get("workspace")

	//查询
	out, err := service.ApiTagList(&swagger.SwaggerHandler{}, key)
	if err != nil {
		view.JSON(response, 200, err)
		return
	}

	view.JSON(response, 200, out)
}

func GetPaths(response http.ResponseWriter, request *http.Request) {

	//当前工作空间
	key := request.Header.Get("workspace")
	findType := request.URL.Query()["type"][0]
	findValue := request.URL.Query()["typeValue"][0]

	//查询
	out, err := service.ApiPathList(&swagger.SwaggerHandler{}, key)
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

func GetPrimaryPath(response http.ResponseWriter, request *http.Request) {

	//当前工作空间
	key := request.Header.Get("workspace")

	values := request.URL.Query()
	primaryId := values["primaryId"][0]

	//查询唯一
	path := service.GetPathPrimary(key, primaryId)
	if path == nil {
		view.JSON(response, 500, "未查询到指定接口")
		return
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

	out := &module.MetaOut{
		Path:     path,
		Request:  req,
		Response: resp,
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

func getReqBodyDefineJson(key string, bodyParam interface{}) string {

	bodyMap := bodyParam.(map[string]interface{})
	schemaMap := bodyMap["schema"].(map[string]interface{})

	bodyMapper := service.GetDefinitionMap(key, schemaMap["$ref"].(string))
	json, err := json.Marshal(bodyMapper)
	if err != nil {
		panic(err)
	}
	return string(json)
}

func getRespBodyDefineJson(key string, path module.ApiPath, code string) string {
	codeMap := path.Responses[code].(map[string]interface{})
	schemaMap := codeMap["schema"].(map[string]interface{})
	bodyMapper := service.GetDefinitionMap(key, schemaMap["$ref"].(string))
	json, err := json.Marshal(bodyMapper)
	if err != nil {
		panic(err)
	}
	return string(json)
}
