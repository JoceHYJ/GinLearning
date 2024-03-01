package v2

import "net/http"

type Context struct {
	Resp http.ResponseWriter
	Req  *http.Request
}
