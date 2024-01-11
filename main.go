package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Api struct {
		Gemini string `yaml:"gemini"`
	} `yaml:"api"`
}

func main() {
	// 读取YAML文件
	yamlFile, err := os.ReadFile("cfg/config.yaml")
	if err != nil {
		fmt.Printf("无法读取YAML文件：%v", err)
	}

	// 解析YAML文件
	var cfg Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		fmt.Printf("无法解析YAML文件：%v", err)
	}

	apiKey := cfg.Api.Gemini
	fmt.Println("apiKey is :", apiKey)
}
