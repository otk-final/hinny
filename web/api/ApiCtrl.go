package api

import (
	"otk-final/hinny/web"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

/**
	接口抓取配置
 */
type DocFetchParameter struct {
	fetchUrl   string
	serverName string
}


func DocFetch(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	vars := mux.Vars(request)

	/**
		执行类
	 */
	var handler ApiCtrl
	if vars["fetchType"] == "swagger" {
		handler = &SwaggerHandler{}
	}

	if handler == nil {
		http.Error(response, "fetchType is illegal", 200)
		return
	}

	/**
		DocFetchParameter 参数
	 */
	var parameter DocFetchParameter
	if err := json.Unmarshal(body, &parameter); err == nil {
		handler.docFetch(&parameter)
		//响应抓取的结果集

	} else {
		http.Error(response, err.Error(), 200)
	}
}

type ApiCtrl interface {
	web.BaseCtrl
	docFetch(args *DocFetchParameter)
}

/**
 	Swagger 抓取接口实现
 */
type SwaggerHandler struct {
	metaType interface{}
}

func (handler *SwaggerHandler) docFetch(args *DocFetchParameter) {

	/**
		通过地址获取swagger暴露的接口信息，进行归档
	 */
	resp, err := http.Get(args.fetchUrl)
	if err != nil {
		panic(err)
	}

	//转换为结构数据进行存储
	body, _ := ioutil.ReadAll(resp.Body)

	//栏目

	//单元请求路径

	//元数据
	fmt.Println(body)
}

