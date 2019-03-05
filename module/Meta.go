package module

type ApiTag struct {
	Name        string `json:"serviceName"`
	Description string `json:"description"`
	PathCount   int    `json:"pathCount"`
}

type ApiPath struct {
	PrimaryId   string                 `json:"primary_id"`
	Tag         *ApiTag                `json:"service"`
	TagName     string                 `json:"tag_name"`
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
}

/**
	元数据-响应
 */
type MetaResponse struct {
	Version string        `json:"version"`
	Header  []interface{} `json:"header"`
	Body    string        `json:"body"`
	Code    int8          `json:"code"`
}

/**
	元数据-验证脚本
 */
type MetaValid struct {
	Version    string `json:"version"`
	Script     string `json:"script"`
	ScriptType string `json:"script_type"`
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
	Valid    *MetaValid    `json:"valid"`
	Result   []*MetaResult `json:"result"`
}
