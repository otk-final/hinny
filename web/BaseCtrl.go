package web

import (
	"github.com/otk-final/respond"
	"net/http"
)

type BaseCtrl interface {
}

var view *respond.Engine

func init() {
	view = respond.New()
	view.LoadHTMLGlob("./view/*")
}

func Index(response http.ResponseWriter, request *http.Request) {
	view.HTML(response, 200, "index.html", nil)
}
