package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

// UserAPI 调用第三方接口请求的数据
type UserAPI struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// TempData 调用第三方接口返回的数据
type TempData struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

// ClientRequest 客户端提交的数据
type ClientRequest struct {
	UserName string      `json:"user_name"`
	Password string      `json:"password"`
	Other    interface{} `json:"other"`
}

// ClientResponse 返回客户端的数据
type ClientResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 开发 Server 程序，通过 Gin 调用其他( Restful )接口
func main() {
	//testAPI()
	r := gin.Default()
	r.POST("/getOtherAPI", getOtherAPI)
	r.Run(":8081")
}

func getOtherAPI(context *gin.Context) {
	var requestData ClientRequest
	var response ClientResponse
	err := context.Bind(requestData)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Msg = "请求的参数错误"
		response.Data = err
		context.JSON(http.StatusBadRequest, response)
		return
	}
	// 请求第三方 API 接口数据
	url := "http://127.0.0.1:8080/login"
	user := UserAPI{requestData.UserName, requestData.Password}
	data, err := getRestfulAPI(url, user, "application/json")
	fmt.Println(data, err)
	// json 的反序列化
	var temp TempData
	json.Unmarshal(data, &temp)
	fmt.Println(temp.Msg, temp.Data)
	response.Code = http.StatusOK
	response.Msg = "请求数据成功"
	response.Data = temp
	context.JSON(http.StatusOK, response)
}

func testAPI() {
	url := "http://127.0.0.1:8080/login"
	user := UserAPI{"user", "123456"}
	data, err := getRestfulAPI(url, user, "application/json")
	fmt.Println(data, err)
	// json 的反序列化
	var temp TempData
	json.Unmarshal(data, &temp)
	fmt.Println(temp.Msg, temp.Data)
}

// 发送 POST 请求
func getRestfulAPI(url string, data interface{}, contentType string) ([]byte, error) {
	// 创建调用 API 接口的 client
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("调用 API 接口出现了错误！")
		return nil, err
	}
	res, err := io.ReadAll(resp.Body)
	return res, err
}
