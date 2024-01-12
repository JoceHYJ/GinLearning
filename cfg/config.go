package cfg

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

func NewConfig() (*Config, error) {
	var config Config
	// 读取YAML文件
	yamlFile, err := os.ReadFile("../../cfg/config.yaml")
	if err != nil {
		return &config, fmt.Errorf("无法读取YAML文件：%v", err)
	}

	// 解析YAML文件
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return &config, fmt.Errorf("无法解析YAML文件：%v", err)
	}

	return &config, nil
}
