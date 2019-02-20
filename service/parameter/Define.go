package parameter

import "net/http"

type ParameterType string

const (
	SQL  ParameterType = "sql"
	BODY ParameterType = "body"
	HEAD ParameterType = "head"
	PATH ParameterType = "path"
	URI  ParameterType = "uri"
)

type ParametersPolicy struct {
	ParameterType ParameterType `json:"parameter_type"` //参数类型
	MetaData      string        `json:"meta_data"`      //元数据
	PresetData    string        `json:"preset_data"`    //预设置
}

/*
	参数打包策略
 */
type ParameterPackPolicy interface {
	/*
		@schemeField 	属性
		@defaultVal		默认值
	*/
	pack(schemeField string, defaultVal string) error
}

/*
	语法
	head.{{loginToken}}
	body.{{data.f1.f2[1].array(length))}}
 */


func FromJson(body interface{}, schemeNode string) error {


	return nil
}

func FromHeader(handler http.Header, key string) string {

	return ""
}
