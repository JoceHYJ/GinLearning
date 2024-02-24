package v1

import "net/http"

type Context struct {
	Resp http.ResponseWriter
	Req  *http.Request
}
