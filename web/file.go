package web

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type FileUploader struct {
	FileField   string                                        // 文件在表单中的字段名
	DstPathFunc func(fileHeader *multipart.FileHeader) string // 目标路径
}

// Handle 处理文件上传请求 返回一个 HandleFunc
func (u FileUploader) Handle() HandleFunc {

	// 这里可以额外做一些校验
	// if u.FileField == "" {
	// 	u.FileField = "file"
	// }
	//if u.DstPathFunc == nil {
	//	// 默认值
	//}

	return func(ctx *Context) {
		// 处理文件上传逻辑
		// 第一步：解析请求中的文件内容
		// 第二步：计算目标路径
		// 第三步：保存文件到目标路径
		// 第四步：返回响应
		file, fileHeader, err := ctx.Req.FormFile(u.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		defer file.Close()
		// 计算目标路径
		// 将目标计算的逻辑交给用户
		dst := u.DstPathFunc(fileHeader)

		// 创建目录结构
		err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("目录创建失败" + err.Error())
			return
		}

		// O_TRUNC表示如果目标路径存在同名文件，则清空该文件的内容
		// 保存文件到目标路径
		dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		defer dstFile.Close()
		// buffer 会影响性能 考虑复用
		_, err = io.CopyBuffer(dstFile, file, nil)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("上传成功")
	}
}

//// FileUploader 的第二种设计模式 Option + HandleFunc
//
//type FileUploaderOption func(uploader *FileUploader
//
//func NewFileUploader(opts...FileUploaderOption) *FileUploader  {
//	res := &FileUploader{
//		FileField: "file",
//		      DstPathFunc: func(fileHeader *multipart.FileHeader) string {
//				  return filepath.Join("testdata", "upload", uuid.New().String())
//			  },
//	}
//	return res
//}
//
//func (u FileUploader) HandleFunc(ctx *Context) {
//	// 文件上传逻辑
//}

type FileDownloader struct {
	Dir string
}

func (f *FileDownloader) Handle() HandleFunc {
	// 文件下载逻辑
	return func(ctx *Context) {
		req, _ := ctx.QueryValue("file").String()
		path := filepath.Join(f.Dir, filepath.Clean(req))
		fn := filepath.Base(path)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		http.ServeFile(ctx.Resp, ctx.Req, path)
	}
}

// 静态资源处理采用 Option Handle 模式

type StaticResourceHandler struct {
	dir                     string
	extensionContentTypeMap map[string]string
}

func NewStaticResourceHandler(dir string) *StaticResourceHandler {
	res := &StaticResourceHandler{
		dir: dir,
		extensionContentTypeMap: map[string]string{
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
			"doc":  "application/msword",
			"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
	}
	return res
}

func (s *StaticResourceHandler) Handle(ctx *Context) {
	// 静态资源处理逻辑
	file, err := ctx.PathValue("file").String()
	if err != nil {
		ctx.RespStatusCode = http.StatusBadRequest
		ctx.RespData = []byte("请求路径错误")
		return
	}
	dst := filepath.Join(s.dir, file)
	data, err := os.ReadFile(dst)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("服务器内部错误") // 避免攻击者通过调用 API 查看文件是否存在
		return
	}
	ext := filepath.Ext(dst)[1:] // 获取文件扩展名
	header := ctx.Resp.Header()
	// 可能的 Content-Type 文本 图片 多媒体
	header.Set("Content-Type", s.extensionContentTypeMap[ext])
	header.Set("Content-Length", strconv.Itoa(len(data)))
	ctx.RespData = data
	ctx.RespStatusCode = http.StatusOK
}
