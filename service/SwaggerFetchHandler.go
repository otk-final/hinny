package service

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
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
	Name        string `json:"serviceName"`
	Description string `json:"description"`
	PathCount   int    `json:"pathCount"`
}

type ApiPath struct {
	PrimaryId   string         `json:"primary_id"`
	Tag         *ApiTag        `json:"service"`
	TagName     string         `json:"tag_name"`
	Path        string         `json:"path"`
	Description string         `json:"description"`
	Method      string         `json:"method"`
	Parameters  []interface{}  `json:"parameters"`
	Definition  *ApiDefinition `json:"responses"`
	Deprecated  bool           `json:"deprecated"`
}

type ApiDefinition struct {
	Properties  interface{} `json:"properties"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
}

/**
Swagger 抓取接口实现
*/
type SwaggerHandler struct {
}

func (handler *SwaggerHandler) DocFetch(fetchUrl string) ([]ApiTag, []ApiPath, []ApiDefinition) {

	/**
	通过地址获取swagger暴露的接口信息，进行归档
	*/
	resp, err := http.Get(fetchUrl)
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

	//标签/服务
	tags := apiOut["tags"].([]interface{})
	apiTags := make([]ApiTag, 0)
	for _, tag := range tags {
		tagMap := tag.(map[string]interface{})

		tagName := tagMap["name"].(string)
		apiTag := &ApiTag{
			Name:        tagName,
			PathCount:   0,
			Description: tagMap["description"].(string),
		}
		apiTags = append(apiTags, *apiTag)
	}

	//请求路径
	paths := apiOut["paths"].(map[string]interface{})
	apiPaths := make([]ApiPath, 0)
	tagPathCounts := make(map[string]int)
	for k, v := range paths {
		//转换
		tagName, slice := handler.convertPath(k, v)
		//存在继续最近数量
		tagPathCounts[tagName] = tagPathCounts[tagName] + len(slice)

		apiPaths = append(apiPaths, slice...)
	}

	//元数据
	definitions := apiOut["definitions"].(map[string]interface{})
	apiDefinitions := make([]ApiDefinition, 0)
	for k, v := range definitions {
		fmt.Println(k, v)
	}

	//回填Tag中的数量
	for i := 0; i < len(apiTags); i++ {
		apiTags[i].PathCount = tagPathCounts[apiTags[i].Name]
	}

	//回填路径中的所属Tag,和元数据definition
	for i := 0; i < len(apiPaths); i++ {
		apiPaths[i].Tag = getTag(apiTags, apiPaths[i].TagName)
	}

	return apiTags, apiPaths, apiDefinitions
}

func getTag(array []ApiTag, name string) *ApiTag {
	for _, item := range array {
		if item.Name == name {
			return &item
		}
	}
	return nil
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
		//参数
		parameters := item["parameters"].([]interface{})
		//独立方法
		path := ApiPath{
			PrimaryId:   item["operationId"].(string),
			TagName:     tagName,
			Path:        key,
			Method:      mk,
			Parameters:  parameters,
			Description: item["summary"].(string),
			Deprecated:  item["deprecated"].(bool),
		}
		pathArray = append(pathArray, path)
	}
	return tagName, pathArray
}
