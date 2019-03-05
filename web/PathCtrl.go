package web

import (
	"net/http"
	"otk-final/hinny/service"
	"strings"
	"otk-final/hinny/service/swagger"
	"otk-final/hinny/module"
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
	key := request.Header.Get("key")
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

	/**
		1，生成scheme文件
		2，补齐相关参数
	 */


	view.JSON(response, 200, path)
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
