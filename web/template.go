package web

import (
	"bytes"
	"context"
	"html/template"
	"io/fs"
)

type TemplateEngine interface {
	// Render 渲染页面

	Render(ctx context.Context, tplName string, data any) ([]byte, error)

	//Render(ctx Context)
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context,
	tplName string, data any) ([]byte, error) {
	res := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(res, tplName, data)
	return res.Bytes(), err
}

func (g *GoTemplateEngine) LoadFromGlob(pattern string) error {
	var err error
	g.T, err = template.ParseGlob(pattern)
	return err
}

func (g *GoTemplateEngine) LoadFromFiles(filenames ...string) error {
	var err error
	g.T, err = template.ParseFiles(filenames...)
	return err
}

func (g *GoTemplateEngine) LoadFromFS(fs fs.FS, patterns ...string) error {
	var err error
	g.T, err = template.ParseFS(fs, patterns...)
	return err
}
