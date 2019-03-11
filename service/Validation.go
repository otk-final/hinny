package service

import (
	"otk-final/hinny/module"
	"github.com/robertkrimen/otto"
	"fmt"
)

type ValidDefine struct {
	vm     *otto.Otto
	result []*module.MetaResult
	script string
}

type ValidCode = int

const (
	FAIL    ValidCode = -1 //异常
	PART    ValidCode = 0  //部分成功
	SUCCESS ValidCode = 1  //成功
)

const ValidTemplateScript = `
/**
 *  公共函数,无参返回全部
 *  请求路径：$GetUri('version')
 *  请求头：  $GetHeader('tenantId')
 *  请求参数：$GetQuery('pageNo')
 *  请求报文：$GetBody()
 * 
 *  验证函数 testing(响应码,响应头,响应体)
 * 
 * 
 *  结果集：$AppendResult('规则说明','验证结果','是否通过[true/false]')
 **/
 
function testing(code,header,body){

  //TODO
  
}
`

func NewValid(scriptType string, script string) *ValidDefine {

	ctx := &ValidDefine{
		script: script,
		result: []*module.MetaResult{},
	}

	//javascript 虚拟执行环境

	vm := otto.New()
	//设置相关变量

	//设置返回值
	vm.Set("$AppendResult", ctx.vmAppendResult())

	ctx.vm = vm
	return ctx
}

func (that *ValidDefine) Valid(request *module.MetaRequest, resp *module.MetaResponse) ([]*module.MetaResult, ValidCode) {
	/**
		如果脚本为空，直接返回
	 */
	if that.script == "" {
		return []*module.MetaResult{}, SUCCESS
	}

	//设置方法
	that.vm.Set("$GetHeader", vmGetValue(request.Header))
	that.vm.Set("$GetUri", vmGetValue(request.Uri))
	that.vm.Set("$GetQuery", vmGetValue(request.Query))
	that.vm.Set("$GetBody", vmGetBody(request.Body))

	//转换报文体为json格式，方便调用
	bodyEval, err := that.vm.Eval("(" + resp.Body + ")")

	//运行初始化脚本
	that.vm.Run(that.script)

	//默认参数为 返回码，响应头，报文体
	values, err := that.vm.Call("testing", nil, resp.Code, resp.Header, bodyEval)
	if err != nil {
		out := []*module.MetaResult{{
			Rule: "脚本编译错误",
			Msg:  err.Error(),
			Ok:   false,}}
		return out, FAIL
	}

	//返回值，暂不做处理
	fmt.Println(values.ToString())

	previous := true
	size := len(that.result)
	for i := 0; i < size; i++ {
		item := that.result[i]

		//部分成功
		if i > 0 && previous != item.Ok {
			return that.result, PART
		}

		previous = item.Ok
		if i < (size - 1) {
			continue
		}

		if previous {
			//全部成功
			return that.result, SUCCESS
		} else {
			//全部失败
			return that.result, FAIL
		}
	}
	return that.result, SUCCESS
}

func (that *ValidDefine) vmAppendResult() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		/**
			获取参数
		 */
		rule := call.Argument(0).String()     //规则说明
		msg := call.Argument(1).String()      //验证结果
		ok, _ := call.Argument(2).ToBoolean() //状态

		/**
			添加到当前上下文中
		 */
		item := &module.MetaResult{
			Rule: rule,
			Msg:  msg,
			Ok:   ok,
		}
		that.result = append(that.result, item)

		return otto.Value{}
	}
}

func vmGetBody(body string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		body, err := call.Otto.Eval("(" + body + ")")
		if err != nil {
			bodyCtx, err := otto.ToValue(body)
			if err != nil {
				return otto.NaNValue()
			}
			return bodyCtx
		}
		return body
	}
}

func vmGetValue(typeValues []interface{}) func(call otto.FunctionCall) otto.Value {

	return func(call otto.FunctionCall) otto.Value {

		var out otto.Value
		var err error

		/**
			获取参数,无参返回所有
		 */
		if len(call.ArgumentList) == 0 {
			out, err = otto.ToValue(typeValues)
			if err != nil {
				out = otto.UndefinedValue()
			}
			return out
		}

		//有参数返回第一个
		argName := call.Argument(0).String()
		for _, item := range typeValues {
			itemMap := item.(map[string]interface{})
			if itemMap["name"].(string) != argName {
				continue
			}

			//数组类型，转换为json数组
			if itemMap["type"] == "array" {
				/*
					数组借助临时变量获取
				 */
				call.Otto.Set("$tempArray", itemMap["value"])
				out, err = call.Otto.Get("$tempArray")
			} else {
				out, err = otto.ToValue(itemMap["value"])
			}

			if err != nil {
				fmt.Println(err)
				return otto.UndefinedValue()
			}
			return out
		}
		return otto.NullValue()
	}
}
