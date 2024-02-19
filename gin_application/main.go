package main

import (
	"GinLearning/gin_application/common"
	"GinLearning/gin_application/route"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
)

func main() {
	InitConfig()
	common.InitDB()
	r := gin.Default()
	r = route.CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		r.Run(":" + port)
	} else {
		r.Run() // 默认端口 8080
	}
}

func InitConfig() {
	workDir, _ := os.Getwd()
	//fmt.Printf("workDir: %v\n", workDir)
	//v := viper.New()
	//v.SetConfigName("application")
	//v.SetConfigType("yaml")
	//v.AddConfigPath(workDir + "/gin_application/config")
	//err := v.ReadInConfig()
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/gin_application/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
