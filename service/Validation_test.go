package service

import (
	"testing"
	"otk-final/hinny/module"
	"fmt"
	"encoding/json"
)

func TestNewValid(t *testing.T) {

	script := `$AppendResult("默认验证","结果",true)
function testing(code,header,body){
	
console.info(code)
console.info(body)
  //$AppendResult("验证报文头是否返回token","c98aed7c-d5e1-413f-9638-a279bf2e8d8c",false)
console.info(body.rule)
  //$AppendResult("验证响应头是合法","application/json;charset=UTF-8不合法",false)
  console.info(body.msg)
  //$AppendResult("验证响应报文头是否有效","无效",false)
}`

	byte, _ := json.Marshal(&module.MetaResult{Rule: "规则"})
	body := string(byte)

	ctx := NewValid("", script)

	resp := &module.MetaResponse{
		Body:   body,
		Code:   200,
		Header: []interface{}{},
	}
	results, code := ctx.Valid(nil,resp)

	byte, _ = json.Marshal(results)

	fmt.Println(string(byte))
	fmt.Println(code)
}
