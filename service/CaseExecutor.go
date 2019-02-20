package service

import (
	"net/http"
	"otk-final/hinny/service/parameter"
	"otk-final/hinny/service/valid"
)

/**
	统一执行上下文
 */
type UnifyProcessContext struct {
	Handler          *HttpProcessHandler
	Error            error
	Resp             *http.Response
	ParametersPolicy []*parameter.ParametersPolicy //参数策略
	ValidatorPolicy  []*valid.ValidatorPolicy      //验证策略
}

func D() {

}
