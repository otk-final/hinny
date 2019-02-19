package service

import "otk-final/hinny/web/testcase"

func ExecCase(define *testcase.CaseFlowConfig) error {

	steps := define.Steps
	for _, step := range *steps {

		/*
			解析参数
			执行
			添加到当前流程上下文
			验证解析
			通知总控当前节点执行信息（状态/请求报文/响应报文/耗时）


		 */



	}
	return nil
}
