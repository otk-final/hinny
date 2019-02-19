package testcase

import (
	"otk-final/hinny/service/valid"
	"otk-final/hinny/web/api"
	"otk-final/hinny/service/parameter"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"otk-final/hinny/service"
)

//流程类型
type StepType int8

const (
	//顺序
	Order StepType = 1
	//分支
	Branch StepType = 2
)

type CaseStepConfig struct {
	Name              string                         `json:"name"`
	Description       string                         `json:"description"`
	Tag               *api.ApiTag                    `json:"tag"`
	Path              *api.ApiPath                   `json:"path"`
	ParametersWrapper *[]parameter.ParametersWrapper `json:"parameters_wrapper"` //参数策略
	ValidatorWrapper  *[]valid.ValidatorWrapper      `json:"validator_wrapper"`  //验证策略
	NextCases         *[]CaseStepConfig              `json:"next_cases"`
}

type CaseFlowConfig struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Steps       *[]CaseStepConfig `json:"steps"`
}

/**
	案例执行上下文
*/
type CaseProcessContext struct {
	Identity string
	ArgsMap  map[string]interface{}
}

/**
	提交案例
 */
func PostCase(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	flowConfig := &CaseFlowConfig{
		Name:        "案例",
		Description: "标准测试案例",
	}

	//反序列化
	if err := json.Unmarshal(body, flowConfig); err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	//入库

	//通知调度器
	service.ExecCase(flowConfig)

	//响应
}
