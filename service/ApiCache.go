package service

import (
	"sync"
	"otk-final/hinny/module"
	"otk-final/hinny/module/db"
	"time"
)

type ApiHandler interface {
	DocFetch(fetchUrl string) ([]module.ApiTag, []module.ApiPath, []module.ApiDefinition)
}

/**
	接口类数据进行缓存
 */
var apiTagCached = make(map[uint64][]module.ApiTag)
var apiPathCached = make(map[uint64][]module.ApiPath)
var apiDefinitionCached = make(map[uint64][]module.ApiDefinition)

//互斥锁，防止同一时间修改缓存
var lock = &sync.Mutex{}

func ApiRemove(kid uint64) {
	delete(apiTagCached, kid)
	delete(apiPathCached, kid)
	delete(apiDefinitionCached, kid)
}

func ApiTagList(fetchHandler ApiHandler, ws *db.Workspace) ([]module.ApiTag, error) {
	out, ok := apiTagCached[ws.Kid]
	//不存在，刷新接口
	if !ok {
		err := ApiRefresh(fetchHandler, ws)
		if err != nil {
			panic(err)
		}
		out = apiTagCached[ws.Kid]
	}
	return out, nil
}

func ApiPathList(fetchHandler ApiHandler, ws *db.Workspace) ([]module.ApiPath, error) {
	out, ok := apiPathCached[ws.Kid]
	//不存在，刷新接口
	if !ok {
		err := ApiRefresh(fetchHandler, ws)
		if err != nil {
			panic(err)
		}
		out = apiPathCached[ws.Kid]
	}
	return out, nil
}

func ApiRefresh(fetchHandler ApiHandler, ws *db.Workspace) error {

	lock.Lock()
	defer lock.Unlock()

	//查询
	apiTags, apiPaths, apiDefinitions := fetchHandler.DocFetch(ws.DocUrl)

	//添加至缓存
	apiTagCached[ws.Kid] = apiTags
	apiPathCached[ws.Kid] = apiPaths
	apiDefinitionCached[ws.Kid] = apiDefinitions

	return nil
}
func GetPathPrimary(key uint64, identity string) *module.ApiPath {

	paths, ok := apiPathCached[key]
	if !ok {
		return nil
	}

	for _, path := range paths {
		if path.PrimaryId == identity {
			return &path
		}
	}
	return nil
}

func GetDefinitionArray(key uint64, objDefine string) []interface{} {
	item := GetDefinitionMap(key, objDefine)
	return []interface{}{item}
}

/**
	根据对象类型生成相关属性
 */
func GetDefinitionMap(key uint64, objDefine string) map[string]interface{} {
	allDefines := apiDefinitionCached[key]

	pj := &PropertyCvtMap{
		maxDeep:    8,
		allDefines: apiDefinitionCached[key],
		getPrimary: func(objDefine string) *module.ApiDefinition {
			for _, define := range allDefines {
				if objDefine == "#/definitions/"+define.Title {
					return &define
				}
			}
			return nil
		},
	}

	define := pj.getPrimary(objDefine)
	if define == nil {
		return nil
	}

	outMap := make(map[string]interface{})
	for k, v := range define.Properties {
		field := pj.propertyFormat(0, v.(map[string]interface{}))
		outMap[k] = field
	}
	return outMap
}

type PropertyCvtMap struct {
	maxDeep    int
	allDefines []module.ApiDefinition
	getPrimary func(objDefine string) *module.ApiDefinition
}

func (pj PropertyCvtMap) propertyFormat(deep int, property map[string]interface{}) interface{} {
	deep++
	fieldType := property["type"]
	if fieldType == "string" {
		format := property["format"]
		if format != nil && format == "date-time" {
			//日期格式化，变态逻辑 取当前时间
			return time.Now().Format("2006-01-02 15:04:05")
		}
		//字符
		return ""
	} else if fieldType == "boolean" {
		return true
	} else if fieldType == "array" {
		//数组
		array := make([]interface{}, 0)
		if items := property["items"]; items != nil {
			array = append(array, pj.propertyFormat(deep, items.(map[string]interface{})))
		}
		return array
	} else if fieldType == "integer" {
		//数字
		return 0
	} else {
		//对象
		fieldMap := make(map[string]interface{}, 0)
		if ref := property["$ref"]; ref != nil && pj.maxDeep > deep {
			nextPs := pj.getPrimary(ref.(string)).Properties
			for k, v := range nextPs {
				field := pj.propertyFormat(deep, v.(map[string]interface{}))
				fieldMap[k] = field
			}
		}
		return fieldMap
	}
}
