package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	// 这部分，改成从 Cfg 里面获取
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	//loc := viper.GetString("datasource.loc")
	loc := "Asia/Shanghai"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc))
	//argsV2 := "root:010729@tcp(127.0.0.1:3306)/gindb?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println("args   is :", args)
	//fmt.Println("argsV2 is: ", argsV2)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}
	DB = db
	return db
}
