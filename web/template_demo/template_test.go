package template_demo

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html"
	"html/template"
	"testing"
)

// Template 基本语法

func TestHelloWorld(t *testing.T) {
	type User struct {
		Name string
	}
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`Hello, {{ .Name}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, User{Name: "World"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, World`, buffer.String())
}

func TestHelloWorld_Map(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`Hello, {{ .Name}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, map[string]string{"Name": "World"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, World`, buffer.String())
}

func TestHelloWorld_Slice(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`Hello, {{index . 0}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, []string{"World"})
	require.NoError(t, err)
	assert.Equal(t, `Hello, World`, buffer.String())
}

func TestHelloWorld_Basic(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`Hello, {{.}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, "World")
	require.NoError(t, err)
	assert.Equal(t, `Hello, World`, buffer.String())
}

// Template 方法调用

func TestFuncCall(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
Slice Length: {{len .Slice}}
{{.Hello "Tomato" "Sprite"}}
Print Number: {{printf "%.2f" 1.234}}`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, FuncCall{
		Slice: []string{"a", "b"},
	})
	require.NoError(t, err)
	assert.Equal(t, `
Slice Length: 2
Hello, Tomato·Sprite
Print Number: 1.23`, buffer.String())
}

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(first string, last string) string {
	return fmt.Sprintf("Hello, %s·%s", first, last)
}

// Template 循环

func TestLoop(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- range $idx, $elem := .Slice}}
{{- .}}
{{$idx}}-{{$elem}}
{{end -}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, FuncCall{
		Slice: []string{"a", "b"},
	})
	require.NoError(t, err)
	assert.Equal(t, `a
0-a
b
1-b
`, buffer.String())
}

// TestForLoop 测试间接的 for loop 按下标索引
func TestForLoop(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- range $idx, $elem := .}}
{{- $idx}},
{{- end -}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, make([]int, 10))
	require.NoError(t, err)
	assert.Equal(t, `0,1,2,3,4,5,6,7,8,9,`, buffer.String())
}

// Template 条件判断 if-else
// 比较操作符 前缀表达式

func TestIfElseBlock(t *testing.T) {
	type User struct {
		Age int
	}
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- if and (gt .Age 0) (le .Age 6) -}}
children: (0 ,6]
{{- else if and (gt .Age 6) (lt .Age 18) -}}
teens: (6 ,18)
{{- else -}}
adults: [18 ,]
{{- end -}}
`)
	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, User{23})
	require.NoError(t, err)
	assert.Equal(t, `adults: [18 ,]`, buffer.String())
}

// Template Pipeline 不建议使用
func TestPipeline(t *testing.T) {
	testCases := []struct {
		name string

		tpl  string
		data any

		want string
	}{
		// 这些例子来自官方文档
		// https://pkg.go.dev/text/template#hdr-Pipelines
		{
			name: "string constant",
			tpl:  `{{"\"output\""}}`,
			want: `"output"`,
		},
		{
			name: "raw string constant",
			tpl:  "{{`\"output\"`}}",
			want: `"output"`,
		},
		{
			name: "function call",
			tpl:  `{{printf "%q" "output"}}`,
			want: `"output"`,
		},
		{
			name: "take argument from pipeline",
			tpl:  `{{"output" | printf "%q"}}`,
			want: `"output"`,
		},
		{
			name: "parenthesized argument",
			tpl:  `{{printf "%q" (print "out" "put")}}`,
			want: `"output"`,
		},
		{
			name: "elaborate call",
			// printf "%s%s" "out" "put"
			tpl:  `{{"put" | printf "%s%s" "out" | printf "%q"}}`,
			want: `"output"`,
		},
		{
			name: "longer chain",
			tpl:  `{{"output" | printf "%s" | printf "%q"}}`,
			want: `"output"`,
		},
		{
			name: "with action using dot",
			tpl:  `{{with "output"}}{{printf "%q" .}}{{end}}`,
			want: `"output"`,
		},
		{
			name: "with action that creates and uses a variable",
			tpl:  `{{with $x := "output" | printf "%q"}}{{$x}}{{end}}`,
			want: `"output"`,
		},
		{
			name: "with action that uses the variable in another action",
			tpl:  `{{with $x := "output"}}{{printf "%q" $x}}{{end}}`,
			want: `"output"`,
		},
		{
			name: "pipeline with action that uses the variable in another action",
			tpl:  `{{with $x := "output"}}{{$x | printf "%q"}}{{end}}`,
			want: `"output"`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tpl := template.New(tc.name)
			tpl, err := tpl.Parse(tc.tpl)
			if err != nil {
				t.Fatal(err)
			}
			bs := &bytes.Buffer{}
			err = tpl.Execute(bs, tc.data)
			if err != nil {
				t.Fatal(err)
			}
			//assert.Equal(t, tc.want, bs.String())
			// 在断言中使用html.UnescapeString处理实际值
			assert.Equal(t, tc.want, html.UnescapeString(bs.String()))
		})
	}
}
