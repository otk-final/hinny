package service

import (
	"io"
	"net/http"
	"net/url"
	"otk-final/hinny/service/parameter"
	"otk-final/hinny/web/api"
	"otk-final/hinny/web/testcase"
	"strings"
	"errors"
	"otk-final/hinny/service/valid"
)

func ExecCase(define *testcase.CaseFlowConfig) error {

	steps := define.Steps

	//对每个独立案例，创建统一Chan进行结果聚合
	unifyChan := make(chan *UnifyProcessContext, len(steps))

	for _, stepConfig := range steps {
		stepClient := &HttpProcessHandler{
			tag:  stepConfig.Tag,
			path: stepConfig.Path,
		}

		//执行
		err := stepClient.process(unifyChan, "", stepConfig.ParameterPolicy, stepConfig.ValidatorPolicy)
		//判断
		if err != nil {
			break
		}
	}

	//关闭
	close(unifyChan)

	return nil
}

type ArrayParameters []*parameter.ParametersPolicy

type ArrayValidators []*valid.ValidatorPolicy

type IHttpCaseCreatorDefine interface {
	setUrl(wrapper ArrayParameters, parameter ArrayParameters) string
	setHeader(wrapper ArrayParameters) *http.Header
	setBody(schemeBody string, wrapper ArrayParameters) io.Reader
}

/*
	案例执行上下文
*/
type HttpProcessHandler struct {
	tag  *api.ApiTag
	path *api.ApiPath
}

func (that *HttpProcessHandler) process(unifyChan chan *UnifyProcessContext, body string, ps ArrayParameters, vs ArrayValidators) error {

	status := &UnifyProcessContext{
		Handler:         that,
		ValidatorPolicy: vs,
	}

	/*
	解析参数
	执行
	添加到当前流程上下文
	验证解析
	通知总控当前节点执行信息（状态/请求报文/响应报文/耗时）
	*/

	pathArgs := make([]*parameter.ParametersPolicy, 0)
	headerArgs := make([]*parameter.ParametersPolicy, 0)
	urlArgs := make([]*parameter.ParametersPolicy, 0)
	bodyArgs := make([]*parameter.ParametersPolicy, 0)

	for _, p := range ps {
		if p.ParameterType == parameter.PATH {
			pathArgs = append(pathArgs, p)
		} else if p.ParameterType == parameter.URI {
			urlArgs = append(urlArgs, p)
		} else if p.ParameterType == parameter.BODY {
			bodyArgs = append(bodyArgs, p)
		} else if p.ParameterType == parameter.HEAD {
			headerArgs = append(headerArgs, p)
		} else {
			continue
		}
	}

	//路径
	urlPath := that.setUrl(pathArgs, urlArgs)

	//报文体
	bodyReader := that.setBody(body, bodyArgs)
	request, err := http.NewRequest(that.path.Method, urlPath, bodyReader)
	if err != nil {
		panic(err)
	}

	//报文头
	request.Header = *that.setHeader(headerArgs)

	//超时设置

	//创建连接，并执行
	client := &http.Client{}
	resp, err := client.Do(request)
	//关流
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	//记录响应
	status.Resp = resp

	//默认执行
	defer func() error {
		//判断当前执行是否有错误
		if err := recover(); err != nil {
			//记录错误
			status.Error = err.(error)
		}
		//对当前状态进行报送到指定channel
		unifyChan <- status
		return nil
	}()

	return nil
}

/*
	设置请求路径
*/
func (that *HttpProcessHandler) setUrl(pathArgs ArrayParameters, parameter ArrayParameters) string {

	//URI解析参数
	orgPath := that.path.Path

	//参数替换
	for _, arg := range pathArgs {
		/*
			判断预设值是否带函数 arg.PresetData
		*/
		orgPath = strings.Replace(orgPath, "{"+arg.MetaData+"}", arg.PresetData, -1)
	}

	urlPath := url.URL{
		Host:   that.path.Host,
		Scheme: "http",
		Path:   orgPath,
	}

	if len(parameter) > 0 {
		values := url.Values{}
		for _, p := range parameter {
			values.Add(p.MetaData, p.PresetData)
		}
		urlPath.RawQuery = values.Encode()
	}

	return urlPath.String()
}

func (that *HttpProcessHandler) setHeader(parameter ArrayParameters) *http.Header {
	header := &http.Header{}
	for _, p := range parameter {
		header.Set(p.MetaData, p.PresetData)
	}
	return header
}

func (that *HttpProcessHandler) setBody(schemeBody string, parameter ArrayParameters) io.Reader {

	if that.path.Method == "GET" {
		return nil
	}

	var contentType string
	var bodyContext string
	if contentType == "form-data" {
		values := url.Values{}
		for _, p := range parameter {
			values.Add(p.MetaData, p.PresetData)
		}
		bodyContext = values.Encode()
	} else if contentType == "json" {
		if len(parameter) == 0 {
			bodyContext = schemeBody
		} else {
			//替换json中的值
			for _, p := range parameter {
				bodyContext = strings.Replace(bodyContext, p.MetaData, p.PresetData, -1)
			}
		}
	} else {
		panic(errors.New("not support content-type"))
	}
	return strings.NewReader(bodyContext)
}
