package valid

type ValidatorType string

const (
	SQL  ValidatorType = "sql"
	JSON ValidatorType = "json"
	LINE ValidatorType = "line"
)

type ValidatorWrapper struct {
	ValidatorType *ValidatorType `json:"validator_type"` //验证类型
	MetaData      string         `json:"meta_data"`      //元数据
	ExpectData    string         `json:"expect_data"`    //预期值
}

/*
	结果验证策略
 */
type ResultValidatorPolicy interface {
	/*
		@src 	原始值
		@target	目标值
	*/
	valid(schemeField string, target string) error
}
