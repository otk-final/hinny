package service

import (
	"otk-final/hinny/module"
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
	"time"
	"otk-final/hinny/module/db"
	"fmt"
	"encoding/json"
)

/**
	案例输入
 */
type CaseInput struct {
	PrimaryId string              `json:"primaryId"`
	Request   *module.MetaRequest `json:"request"`
	Valid     *module.MetaValid   `json:"valid"`
}

/**
	案例输出
 */
type CaseOutput struct {
	LogKid   uint64               `json:"logKid"`
	Time     time.Duration        `json:"time"`
	Response *module.MetaResponse `json:"response"`
	Result   []*module.MetaResult `json:"result"`
}

func Execute(ws *db.Workspace, path *module.ApiPath, input *CaseInput) (*CaseOutput, error) {

	//序列化化存储请求相关参数
	reqCtx, err := json.Marshal(input.Request)
	if err != nil {
		fmt.Print(err)
	}

	logKid, _ := db.GetNextKid()
	//记录日志
	log := &db.CaseLog{
		Kid:          logKid,
		WsKId:        ws.Kid,
		PathIdentity: path.PrimaryId,
		Path:         path.Path,
		MetaRequest:  string(reqCtx),
		CreateTime:   time.Now(),
		ScriptType:   input.Valid.ScriptType,
		Script:       input.Valid.Script,
	}

	//远程调用
	_, metaResp, err := dispatch(ws.ApiUrl, path, input)

	respCtx, err := json.Marshal(metaResp)
	if err != nil {
		fmt.Print(err)
	}

	//初始化验证对象
	valid := NewValid(input.Valid.ScriptType, input.Valid.Script)

	//转换对象进行验证
	results, validCode := valid.Valid(input.Request, metaResp)
	resultsCtx, _ := json.Marshal(results)

	//存储响应，验证
	log.MetaResult = string(resultsCtx)
	log.MetaResponse = string(respCtx)
	log.Curl = ""
	log.Status = validCode

	//保存
	db.Conn.Insert(log)

	return &CaseOutput{LogKid: logKid, Response: metaResp, Result: results}, nil
}

/**
	请求分发
 */
func dispatch(host string, path *module.ApiPath, input *CaseInput) (*http.Request, *module.MetaResponse, error) {

	//请求参数
	reqUri := input.cvtUrl(host, path.Path)
	reqValues := *input.cvtUrlValues()

	//GET请求参数拼接在URL后面
	if strings.EqualFold(path.Method, "GET") {
		reqUri.RawQuery = reqValues.Encode()
	}

	//请求对象
	request := &http.Request{
		Header: *input.cvtHeader(),           //请求头
		Method: strings.ToUpper(path.Method), //方法
		URL:    reqUri,                       //地址
	}

	//报文头提交
	if !strings.EqualFold(path.Method, "GET") {
		request.PostForm = reqValues
		request.Body = ioutil.NopCloser(strings.NewReader(input.Request.Body))
	}

	//建立连接
	client := &http.Client{Timeout: input.Request.TimeOut}

	//执行
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	//转换报文
	headerArray := make([]interface{}, 0)
	for k, v := range resp.Header {
		item := make(map[string]interface{})
		item["name"] = k
		if len(v) > 1 {
			item["type"] = "array"
		} else {
			item["type"] = "string"
		}
		item["value"] = v
		headerArray = append(headerArray, item)
	}

	//响应相关信息
	metaResp := &module.MetaResponse{
		Header: headerArray,     //头
		Code:   resp.StatusCode, //响应码
		Body:   string(body),    //响应体
	}

	fmt.Println(metaResp.Body)
	return request, metaResp, nil
}

/**
	转换Header
 */
func (that *CaseInput) cvtHeader() *http.Header {
	reqHeader := &http.Header{}
	inputHeader := that.Request.Header
	for _, item := range inputHeader {
		itemMap := item.(map[string]interface{})
		reqHeader.Add(itemMap["name"].(string), itemMap["value"].(string))
	}
	return reqHeader
}

/**
	转换路径
 */
func (that *CaseInput) cvtUrl(host string, stringPath string) *url.URL {
	uri := that.Request.Uri

	//转换路径中的占位符
	for _, item := range uri {
		itemMap := item.(map[string]interface{})
		key := itemMap["name"].(string)
		value := itemMap["value"]
		if value != nil && value != "" {
			stringPath = strings.Replace(stringPath, "{"+key+"}", value.(string), -1)
		}
	}

	//合并路径
	suffix := ""
	if !strings.HasSuffix(host, "/") && !strings.HasPrefix(stringPath, "/") {
		suffix = "/"
	}
	realPath := host + suffix + stringPath

	reqUri, err := url.Parse(realPath)
	if err != nil {
		panic(err)
	}

	return reqUri
}

/**
	转换参数
 */
func (that *CaseInput) cvtUrlValues() *url.Values {
	queries := that.Request.Query
	values := &url.Values{}
	for _, item := range queries {
		itemMap := item.(map[string]interface{})
		key := itemMap["name"].(string)
		itemValues := itemMap["value"]
		if itemValues == nil {
			continue
		}
		/**
			判断值类型
		 */
		switch itemValues.(type) {
		case []interface{}:
			array := itemValues.([]interface{})
			for _, val := range array {
				values.Add(key, val.(string))
			}
			break
		default:
			values.Add(key, itemValues.(string))
			break
		}
	}
	return values
}
