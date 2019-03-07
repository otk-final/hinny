package service

import "otk-final/hinny/module"

type ValidDefine struct {
}

type ValidCode = int

const (
	FAIL    ValidCode = -1 //异常
	PART    ValidCode = 0  //部分成功
	SUCCESS ValidCode = 1  //成功
)

func NewValid(scriptType string, script string) *ValidDefine {

	return &ValidDefine{}
}

func (that ValidDefine) valid(resp module.MetaResponse) ([]*module.MetaResult, ValidCode) {

	return nil, FAIL
}
