package parameter

type ParameterType string

const (
	SQL  ParameterType = "sql"
	JSON ParameterType = "json"
	HEAD ParameterType = "head"
	PATH ParameterType = "path"
	URI  ParameterType = "uri"
)

type ParametersWrapper struct {
	ParameterType *ParameterType `json:"parameter_type"` //参数类型
	MetaData      string         `json:"meta_data"`      //元数据
	PresetData    string         `json:"preset_data"`    //预设置
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
