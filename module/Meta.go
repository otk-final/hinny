package module

import "time"

type ApiTag struct {
	Name        string `json:"serviceName"`
	Description string `json:"description"`
	PathCount   int    `json:"pathCount"`
}

type ApiPath struct {
	PrimaryId   string                 `json:"primaryId"`
	Tag         *ApiTag                `json:"service"`
	TagName     string                 `json:"tagName"`
	Path        string                 `json:"path"`
	Description string                 `json:"description"`
	Method      string                 `json:"method"`
	Parameters  []interface{}          `json:"parameters"`
	Responses   map[string]interface{} `json:"responses"`
	Deprecated  bool                   `json:"deprecated"`
}

type ApiDefinition struct {
	Properties  map[string]interface{} `json:"properties"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
}

/**
	元数据-请求
 */
type MetaRequest struct {
	Version string        `json:"version"`
	Header  []interface{} `json:"header"`
	Uri     []interface{} `json:"uri"`
	Query   []interface{} `json:"query"`
	Body    string        `json:"body"`
	TimeOut time.Duration `json:"timeOut"`
}

/**
	元数据-响应
 */
type MetaResponse struct {
	Version string        `json:"version"`
	Header  []interface{} `json:"header"`
	Body    string        `json:"body"`
	Code    int           `json:"code"`
}

/**
	元数据-验证脚本
 */
type MetaValid struct {
	Version    string `json:"version"`
	Script     string `json:"script"`
	ScriptType string `json:"scriptType"`
}

/**
	元数据-验证结果
 */
type MetaResult struct {
	Version string `json:"version"`
	Rule    string `json:"rule"`
	Msg     string `json:"msg"`
	Ok      bool   `json:"ok"`
}

type MetaOut struct {
	Path     *ApiPath      `json:"path"`
	Request  *MetaRequest  `json:"request"`
	Response *MetaResponse `json:"response"`
	Curl     string        `json:"curl"`
	Valid    *MetaValid    `json:"valid"`
	Result   []*MetaResult `json:"result"`
}

//获取值
func (that *MetaRequest) GetTypeValues(funcType string) []interface{} {
	if funcType == "uri" {
		return that.Uri
	} else if funcType == "header" {
		return that.Header
	} else if funcType == "query" {
		return that.Query
	}
	return nil
}

//设置值
func (that *MetaRequest) TypeValues(funcType string, values []interface{}) {
	if funcType == "uri" {
		that.Uri = values
	} else if funcType == "header" {
		that.Header = values
	} else if funcType == "query" {
		that.Query = values
	}
}
