package valid

import (
	"encoding/json"
	"strings"
	"strconv"
	"errors"
)

type HttpBodyValidator struct {
	Body string
}

func (that *HttpBodyValidator) Get(schemeField string, body string) (*HttpBodyValidator, string, error) {

}

func GetVal(schemeField string, body string) (*HttpBodyValidator, string, error) {

	that := &HttpBodyValidator{
		Body: body,
	}

	//解析报文
	node := make(map[string]interface{})
	if err := json.Unmarshal([]byte(that.Body), node); err != nil {
		return that, "", err
	}

	//获取值
	out, funcName := that.jsonExpressionResolve(schemeField, node)
	if out == nil {
		return that, "", errors.New("value is nil")
	}

	if !strings.EqualFold("length", funcName) {
		return that, out.(string), nil
	}

	//解析函数



	return that, "", nil
}

func (that *HttpBodyValidator) valid(schemeField string, target string) error {

	//解析报文
	node := make(map[string]interface{})
	if err := json.Unmarshal([]byte(that.Body), node); err != nil {
		return err
	}

	//获取值
	out, funcName := that.jsonExpressionResolve(schemeField, node)
	if out == nil {
		return errors.New("value is nil")
	}

	//解析函数
	if strings.EqualFold("length", funcName) {
		//长度
		targetLength, _ := strconv.Atoi(target)
		srcLength := 0
		switch out.(type) {
		case string:
			srcLength = len(out.(string))
			break
		case map[string]interface{}:
			srcLength = len(out.(map[string]interface{}))
			break
		case []interface{}:
			srcLength = len(out.([]interface{}))
			break
		}
		//判断长度
		if targetLength != srcLength {
			return errors.New("length is not match")
		}
	} else if strings.EqualFold("scheme", funcName) {
		//格式

	} else {
		//值匹配
		if out.(string) != target {
			return errors.New("context is not match")
		}
	}
	return nil
}

/*
	example:
	json格式schemeField 定位到指定字段
	一级对象：field
	二级对象：field1.field2
	二级对象数组：field1.field2[idx] ,idx下标(从0开始)
	其他：field1.field2[idx].field3
	条目数:field(length)			如果field是对象 则返回子属性个数，如果是数组，则返回数组length
*/
func (that *HttpBodyValidator) jsonExpressionResolve(schemeField string, root interface{}) (interface{}, string) {
	/**
		获取下个节点数据类型
	 */
	typeClassicQuery := func(current interface{}, idx int, fieldName string) interface{} {
		//判断类型
		switch current.(type) {
		case []interface{}:
			return current.([]interface{})[idx]
		case map[string]interface{}:
			return current.(map[string]interface{})[fieldName]
		case string:
			return current.(string)
		case int64, int32, int16, int8:
			return current.(int)
		case bool:
			return current.(bool)
		default:
			return nil
		}
	}

	//分割截取
	splits := strings.Split(schemeField, ".")
	funcName := "" //默认函数名
	for i := 0; i < len(splits); i++ {

		fieldExp := splits[i] //表达式
		idx := 0              //默认下标
		fieldName := fieldExp //默认属性名

		if strings.LastIndexAny(fieldExp, "]") != -1 {
			//数组取值
			idx, _ = strconv.Atoi(getBetweenCtx(fieldExp, "[", "]"))
			fieldName = fieldExp[0:strings.LastIndex(fieldExp, "]")]
		} else if strings.LastIndexAny(fieldExp, ")") != -1 {
			//函数取值
			funcName = getBetweenCtx(fieldExp, "(", ")")
			fieldName = fieldExp[0:strings.LastIndex(fieldExp, ")")]
		}

		//递归查找
		root = typeClassicQuery(root, idx, fieldName)
		if root == nil {
			return nil, funcName
		}

		if i+1 == len(splits) {
			return root, funcName
		}
	}
	return nil, funcName
}

func getBetweenCtx(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}
