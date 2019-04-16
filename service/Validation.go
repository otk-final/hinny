package service

import (
	"otk-final/hinny/module"
	"github.com/robertkrimen/otto"
	"fmt"
	"reflect"
	"strings"
	"encoding/json"
)

type ValidDefine struct {
	vm      *otto.Otto
	result  []*module.MetaResult
	request *module.MetaRequest
	script  string
}

type ValidCode = int

const (
	FAIL    ValidCode = -1 //异常
	PART    ValidCode = 0  //部分成功
	SUCCESS ValidCode = 1  //成功
)

func NewValid(script string, request *module.MetaRequest) *ValidDefine {

	ctx := &ValidDefine{
		script:  script,
		result:  []*module.MetaResult{},
		request: request,
	}



	//运行初始化脚本

	//javascript 虚拟执行环境
	vm := otto.New()
	//设置返回值
	vm.Set("$appendResult", ctx.appendResultFunc())

	//获取值
	vm.Set("$getHeader", ctx.getValueFunc("header"))
	vm.Set("$getUri", ctx.getValueFunc("uri"))
	vm.Set("$getQuery", ctx.getValueFunc("query"))
	vm.Set("$getBody", ctx.getBodyFunc())

	//设置值
	vm.Set("$setHeader", ctx.setValueFunc("header"))
	vm.Set("$setUri", ctx.setValueFunc("uri"))
	vm.Set("$setQuery", ctx.setValueFunc("query"))
	vm.Set("$setBody", ctx.setBodyFunc())


	/**
		如果脚本为空，直接返回
 	*/
	if script == "" {
		ctx.vm = vm
		return ctx
	}


	value, err := vm.Run(script)
	if err != nil {
		ctx.result = append(ctx.result, &module.MetaResult{
			Rule: "脚本加载编译错误",
			Msg:  err.Error(),
			Ok:   false,})
	}
	fmt.Println(value)
	ctx.vm = vm
	return ctx
}

/**
	初始化相关函数
 */
func (that *ValidDefine) BeforeInit() (request *module.MetaRequest, err error) {

	/**
		方法不存在
	 */
	init, err := that.vm.Get("init")
	if !init.IsFunction() {
		return that.request, nil
	}

	values, err := that.vm.Call("init", nil)
	/**
		不存在，或者错误，返回
	 */
	if err != nil {
		that.result = append(that.result, &module.MetaResult{
			Rule: "脚本[init]编译错误",
			Msg:  err.Error(),
			Ok:   false,})
		return that.request, err
	}

	fmt.Println(values)
	return that.request, nil
}

func (that *ValidDefine) AfterValid(resp *module.MetaResponse) ([]*module.MetaResult, ValidCode) {

	/**
		方法不存在
 	*/
	testing, err := that.vm.Get("testing")
	if !testing.IsFunction() {
		return that.result, SUCCESS
	}

	//转换报文体为json格式，方便调用
	bodyEval, err := that.vm.Eval("(" + resp.Body + ")")

	//默认参数为 返回码，响应头，报文体
	values, err := that.vm.Call("testing", nil, resp.Code, resp.Header, bodyEval)
	if err != nil {
		that.result = append(that.result, &module.MetaResult{
			Rule: "脚本[testing]编译错误",
			Msg:  err.Error(),
			Ok:   false,})
	}

	//返回值，暂不做处理
	fmt.Println(values)

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

func (that *ValidDefine) appendResultFunc() func(call otto.FunctionCall) otto.Value {
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

func (that *ValidDefine) getBodyFunc() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {

		if that.request.Body == "" {
			return otto.NullValue()
		}

		body, err := call.Otto.Eval("(" + that.request.Body + ")")
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

func (that *ValidDefine) setBodyFunc() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {

		if that.request.Body == "" {
			return otto.NullValue()
		}
		obj, err := call.Argument(0).Export()
		if err != nil {
			fmt.Println(err.Error())
			return call.This
		}
		//重新序列化为json格式

		byte, err := json.Marshal(obj)
		if err != nil {
			fmt.Println(err.Error())
			return call.This
		}

		that.request.Body = string(byte)
		return call.This
	}
}

func (that *ValidDefine) setValueFunc(funcType string) func(call otto.FunctionCall) otto.Value {

	return func(call otto.FunctionCall) otto.Value {

		values := that.request.GetTypeValues(funcType)
		defer func() {
			that.request.TypeValues(funcType, values)
		}()

		//属性名
		argName := call.Argument(0).String()
		//值
		argValue, err := call.Argument(1).Export()
		if err != nil {
			return call.This
		}
		/**
			值存在则替换
		 */
		size := len(values)
		exist := false
		for idx := 0; idx < size; idx++ {
			itemMap := values[idx].(map[string]interface{})
			if itemMap["name"].(string) == argName {
				itemMap["value"] = argValue
				exist = true
			}
		}
		//不存在，则新增
		if !exist {
			typeName := reflect.TypeOf(argValue).String()
			if strings.HasPrefix(typeName, "[]") {
				typeName = "array"
			} else {
				typeName = "string"
			}
			values = append(values, map[string]interface{}{
				"name":  argName,
				"type":  typeName,
				"value": argValue,
			})
		}
		return call.This
	}
}

func (that *ValidDefine) getValueFunc(funcType string) func(call otto.FunctionCall) otto.Value {

	return func(call otto.FunctionCall) otto.Value {
		typeValues := that.request.GetTypeValues(funcType)

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
		argName, err := call.Argument(0).ToString()
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
