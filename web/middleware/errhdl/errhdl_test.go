package errhdl

import (
	"GinLearning/web"
	"bytes"
	"html/template"
	"net/http"
	"testing"
)

func TestNewMiddlewareBuilder_Builder(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.AddCode(http.StatusNotFound, []byte(`
<html>
	<body>
		<h1>gugugu 走丢了</h1>
	</body>
</html>
`)).AddCode(http.StatusBadRequest, []byte(`
<html>
	<body>
		<h1>请求有误</h1>
	</body>
</html>
`))
	s := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	s.Start(":8081")
}

func TestNewMiddlewareBuilder_Builder_withTemplate(t *testing.T) {
	s := web.NewHTTPServer()
	s.Get("/user", func(ctx *web.Context) {
		ctx.RespData = []byte("hello, world")
	})
	page := `
<html>
	<h1>gugugu 走丢了</h1>
</html>
`
	tpl, err := template.New("404").Parse(page)
	if err != nil {
		t.Fatal(err)
	}
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, nil)
	if err != nil {
		t.Fatal(err)
	}
	s.Use(NewMiddlewareBuilder().
		RegisterError(404, buffer.Bytes()).Build())

	s.Start(":8081")
}
