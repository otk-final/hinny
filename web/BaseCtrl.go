package web

import (
	"github.com/otk-final/lets-go/response"
	"net/http"
)

type BaseCtrl interface {
}

var view *response.Engine

func init() {
	view = response.New()
	view.LoadHTMLGlob("./view/*")
}

func Index(response http.ResponseWriter, request *http.Request) {
	view.HTML(response, 200, "index.html", nil)
}
