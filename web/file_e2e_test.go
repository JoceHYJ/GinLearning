//go:build e2e

package web

import (
	"github.com/stretchr/testify/require"
	"html/template"
	"log"
	"mime/multipart"
	"path/filepath"
	"testing"
)

func TestFileUploader(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	require.NoError(t, err)
	engine := &GoTemplateEngine{
		T: tpl,
	}

	s := NewHTTPServer(ServerWithTemplateEngine(engine))

	s.Get("/upload", func(ctx *Context) {
		err := ctx.Render("upload.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
	})

	// 上传文件
	uploader := FileUploader{
		FileField: "myfile",
		DstPathFunc: func(fileHeader *multipart.FileHeader) string {
			return filepath.Join("testdata", "uploads", fileHeader.Filename)
		},
	}
	s.Post("/upload", uploader.Handle())

	s.Start(":8081")
}

func TestFileDownloader(t *testing.T) {

	s := NewHTTPServer()

	// 下载文件
	downloader := FileDownloader{
		Dir: filepath.Join("testdata", "downloads"),
	}
	s.Get("/download", downloader.Handle())
	// http://localhost:8081/download?file=myfile.txt
	s.Start(":8081")
}

func TestStaticResourceHandler(t *testing.T) {
	s := NewHTTPServer()
	sr, err := NewStaticResourceHandler(filepath.Join("testdata", "static"))
	require.NoError(t, err)
	// http://localhost:8081/static/xxx.jpg
	s.Get("/static/:file", sr.Handle)
	s.Start(":8081")
}
