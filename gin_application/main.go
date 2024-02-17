package main

import (
	"GinLearning/gin_application/common"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
)

func main() {
	InitConfig()
	common.InitDB()
	r := gin.Default()
	port := viper.GetString("server.port")
	r.Run(":" + port)
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
	// 这里面，可以从 yaml 当中读取并写入到 Cfg 里面

}
