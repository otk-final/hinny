package service

import (
	"sync"
	"otk-final/hinny/module"
	"otk-final/hinny/module/db"
	"strings"
	"time"
)

type ApiHandler interface {
	DocFetch(fetchUrl string) ([]module.ApiTag, []module.ApiPath, []module.ApiDefinition)
}

/**
	接口类数据进行缓存
 */
var apiTagCached = make(map[string][]module.ApiTag)
var apiPathCached = make(map[string][]module.ApiPath)
var apiDefinitionCached = make(map[string][]module.ApiDefinition)

//互斥锁，防止同一时间修改缓存
var once = &sync.Once{}

func ApiTagList(fetchHandler ApiHandler, key string) ([]module.ApiTag, error) {
	out, ok := apiTagCached[key]
	//不存在，刷新接口
	if !ok {
		err := ApiRefresh(fetchHandler, key)
		if err != nil {
			panic(err)
		}
		out = apiTagCached[key]
	}
	return out, nil
}

func ApiPathList(fetchHandler ApiHandler, key string) ([]module.ApiPath, error) {
	out, ok := apiPathCached[key]
	//不存在，刷新接口
	if !ok {
		err := ApiRefresh(fetchHandler, key)
		if err != nil {
			panic(err)
		}
		out = apiPathCached[key]
	}
	return out, nil
}

func ApiRefresh(fetchHandler ApiHandler, key string) error {
	/**
		查询工作空间
	 */
	ws := &db.Workspace{
		ApiUrl: "http://api-dev.yryz.com/gateway/lovelorn/v2/api-docs",
	}
	ok, err := db.Session.Cols("api_url").Where("ws_key=?", key).Get(ws)
	if !ok && err != nil {
		return err
	}
	/**
		只允许独立线程进行返回
	 */
	once.Do(func() {
		//查询
		apiTags, apiPaths, apiDefinitions := fetchHandler.DocFetch(ws.ApiUrl)

		//添加至缓存
		apiTagCached[key] = apiTags
		apiPathCached[key] = apiPaths
		apiDefinitionCached[key] = apiDefinitions
	})
	return nil
}
func GetPathPrimary(key string, identity string) *module.ApiPath {

	return nil
}

/**
	根据对象类型生成相关属性
 */
func GetDefinitionMap(key string, objDefine string) map[string]interface{} {
	allDefines := apiDefinitionCached[key]

	pj := &PropertyCvtMap{
		maxDeep:    10,
		allDefines: apiDefinitionCached[key],
		getPrimary: func(objDefine string) *module.ApiDefinition {
			for _, define := range allDefines {
				if strings.LastIndex(objDefine, define.Title) != -1 {
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
