package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "otk-final/hinny/module"
	"otk-final/hinny/web"
)

/**
接口抓取配置
*/
type DocFetchParameter struct {
	Type       string `json:type`
	FetchUrl   string `json:"fetch_url"`
	ServerName string `json:"server_name"`
}

type ApiTag struct {
	Name        string    `json:name`
	Description string    `json:"description"`
	Paths       []ApiPath `json:"paths"`
}

type ApiPath struct {
	Host       string                 `json:"host"`
	Path       string                 `json:"path"`
	Summary    string                 `json:"summary"`
	Method     string                 `json:"method"`
	Parameters map[string]interface{} `json:"parameters"`
	Responses  map[string]interface{} `json:"responses"`
	Deprecated bool                   `json:"deprecated"`
}

func DocFetch(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(string(body))
	/**
	DocFetchParameter 参数
	*/
	parameter := &DocFetchParameter{}
	if err := json.Unmarshal(body, parameter); err == nil {
		/**
		执行类
		*/
		var handler ApiCtrl
		if parameter.Type == "swagger" {
			handler = &SwaggerHandler{}
		}

		if handler == nil {
			http.Error(response, "fetchType is illegal", 200)
			return
		}

		//执行
		apiTags := handler.docFetch(parameter)

		fmt.Println(apiTags)
	} else {
		http.Error(response, err.Error(), 200)
	}
}

type ApiCtrl interface {
	web.BaseCtrl
	docFetch(args *DocFetchParameter) []ApiTag
}

/**
Swagger 抓取接口实现
*/
type SwaggerHandler struct {
	metaType interface{}
}

func (handler *SwaggerHandler) docFetch(args *DocFetchParameter) []ApiTag {

	/**
	通过地址获取swagger暴露的接口信息，进行归档
	*/
	resp, err := http.Get(args.FetchUrl)
	if err != nil {
		panic(err)
	}
	//转换为结构数据进行存储
	body, _ := ioutil.ReadAll(resp.Body)

	//转换
	apiOut := make(map[string]interface{})
	if err := json.Unmarshal(body, &apiOut); err != nil {
		panic(err)
	}

	/*
		存储api元配置项
		并缓存tags/paths 基本信息，用于前端配置
	*/

	pathMap := map[string][]ApiPath{}

	paths := apiOut["paths"].(map[string]interface{})
	for k, v := range paths {
		//转换
		tagName, slice := handler.convertPath(k, v)

		//存在继续添加
		if sliceAll, ok := pathMap[tagName]; ok {
			sliceAll = append(sliceAll, slice...)
			pathMap[tagName] = sliceAll
		} else {
			pathMap[tagName] = slice
		}
	}

	tags := apiOut["tags"].([]interface{})
	apiTags := make([]ApiTag, 0)

	for _, tag := range tags {
		tagMap := tag.(map[string]interface{})

		tagName := tagMap["name"].(string)

		apiTag := &ApiTag{
			Name:        tagName,
			Paths:       pathMap[tagName],
			Description: tagMap["description"].(string),
		}
		apiTags = append(apiTags, *apiTag)
	}

	return apiTags
}

/**
返回当前tag和方法
*/
func (handler *SwaggerHandler) convertPath(key string, val interface{}) (string, []ApiPath) {

	pathArray := make([]ApiPath, 0)

	methods := val.(map[string]interface{})

	var tagName string
	for mk, mv := range methods {

		item := mv.(map[string]interface{})

		//唯一标识
		tagName = item["tags"].([]interface{})[0].(string)

		//独立方法
		path := ApiPath{
			Path:       key,
			Method:     mk,
			Summary:    item["summary"].(string),
			Deprecated: item["deprecated"].(bool),
		}
		pathArray = append(pathArray, path)
	}

	return tagName, pathArray
}
