package web

import (
	"net/http"
	"otk-final/hinny/service"
	"fmt"
	"strings"
)

var ApiTags = make([]service.ApiTag, 0)
var ApiPaths = make([]service.ApiPath, 0)
var ApiDefinition = make([]service.ApiDefinition, 0)

func init() {

}

func GetServices(response http.ResponseWriter, request *http.Request) {
	//当前工作空间
	workspaceKey := request.Header.Get("workspace")
	workspaceKey = "dev"
	ws := FindWorkspace(workspaceKey)
	if ws == nil {
		view.JSON(response, 200, make([]service.ApiTag, 0))
		return
	}

	//抓取解析
	handler := &service.SwaggerHandler{}
	ApiTags, ApiPaths, ApiDefinition = handler.DocFetch(ws.Host + "/v2/api-docs")

	view.JSON(response, 200, ApiTags)
}

func GetPaths(response http.ResponseWriter, request *http.Request) {
	//args := mux.Vars(request)

	//当前工作空间
	//workspaceKey := request.Header.Get("workspaceKey")
	findType := request.URL.Query()["type"][0]
	findValue := request.URL.Query()["typeValue"][0]

	multipleContains := func(target string, srcArray ... string) bool {
		for _, src := range srcArray {

			if strings.Contains(src, target) {
				return true
			}
		}
		return false
	}

	outs := make([]service.ApiPath, 0)

	for _, path := range ApiPaths {
		match := matchPath(findType, findValue, &path, multipleContains)
		if match {
			outs = append(outs, path)
		}
	}
	fmt.Println(findType, findValue)
	view.JSON(response, 200, outs)
}

func GetPrimaryPath(response http.ResponseWriter, request *http.Request) {
	values := request.URL.Query()
	primaryId := values["primaryId"][0]
	var out service.ApiPath
	for _, path := range ApiPaths {
		if path.PrimaryId == primaryId {
			out = path
			break
		}
	}
	view.JSON(response, 200, out)
}

func matchPath(findType string, findValue string, path *service.ApiPath,
	matchFunc func(target string, srcArray ... string) bool) bool {
	if findType == "service" {
		return matchFunc(findValue, path.Tag.Name, path.Tag.Description, path.Description)
	}
	if findType == "path" {
		return matchFunc(findValue, path.Path)
	}
	return false
}
