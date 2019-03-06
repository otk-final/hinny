package swagger

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"otk-final/hinny/module"
)

/**
Swagger 抓取接口实现
*/
type SwaggerHandler struct {
}

func (handler *SwaggerHandler) DocFetch(fetchUrl string) ([]module.ApiTag, []module.ApiPath, []module.ApiDefinition) {

	nullFunc := new(struct {
		item      map[string]interface{}
		stringCvt func(name string) string
		objCvt    func(name string) map[string]interface{}
	})

	//空字符串
	nullFunc.stringCvt = func(name string) string {
		out := nullFunc.item[name]
		if out == nil {
			return ""
		}
		return out.(string)
	}
	//空对象
	nullFunc.objCvt = func(name string) map[string]interface{} {
		out := nullFunc.item[name]
		if out == nil {
			return make(map[string]interface{})
		}
		return out.(map[string]interface{})
	}

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
	apiTags := make([]module.ApiTag, 0)
	for _, tag := range tags {
		tagMap := tag.(map[string]interface{})

		tagName := tagMap["name"].(string)
		apiTag := &module.ApiTag{
			Name:        tagName,
			PathCount:   0,
			Description: tagMap["description"].(string),
		}
		apiTags = append(apiTags, *apiTag)
	}

	//请求路径
	paths := apiOut["paths"].(map[string]interface{})
	apiPaths := make([]module.ApiPath, 0)
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
	apiDefinitions := make([]module.ApiDefinition, 0)
	for _, v := range definitions {
		nullFunc.item = v.(map[string]interface{})
		define := &module.ApiDefinition{
			Title:       nullFunc.stringCvt("title"),
			Description: nullFunc.stringCvt("description"),
			Properties:  nullFunc.objCvt("properties"),
		}
		apiDefinitions = append(apiDefinitions, *define)
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

func getTag(array []module.ApiTag, name string) *module.ApiTag {
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
func (handler *SwaggerHandler) convertPath(key string, val interface{}) (string, []module.ApiPath) {
	pathArray := make([]module.ApiPath, 0)
	methods := val.(map[string]interface{})
	var tagName string
	for mk, mv := range methods {
		item := mv.(map[string]interface{})
		//唯一标识
		tagName = item["tags"].([]interface{})[0].(string)
		//参数
		parameters := item["parameters"].([]interface{})
		//独立方法
		path := module.ApiPath{
			PrimaryId:   item["operationId"].(string),
			TagName:     tagName,
			Path:        key,
			Method:      mk,
			Parameters:  parameters,
			Description: item["summary"].(string),
			Deprecated:  item["deprecated"].(bool),
			Responses:   item["responses"].(map[string]interface{}),
		}
		pathArray = append(pathArray, path)
	}
	return tagName, pathArray
}