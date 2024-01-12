package api

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestGemini(t *testing.T) {
	type args struct {
		prompts string
		text    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "translation",
			args: args{
				prompts: "请你把下面的话翻译成英文：",
				text:    "放鸽子了",
			},
		},
		{
			name: "coding",
			args: args{
				prompts: "请帮我编写代码：",
				text:    "使用Go输出 “大明放鸽子了” ",
			},
		},
		{
			name: "coding",
			args: args{
				prompts: "请告诉我正月可以理发嘛：",
				text:    "正月可以理发嘛？答：",
			},
		},
		{
			name: "coding",
			args: args{
				prompts: "请告诉我正月为什么不可以理发嘛：",
				text:    "正月为什么不可以理发嘛？答：",
			},
		},
	}

	// 读取YAML文件
	yamlFile, err := os.ReadFile("../cfg/config.yaml")
	if err != nil {
		fmt.Printf("无法读取YAML文件：%v", err)
	}
	type Config struct {
		Api struct {
			Gemini string `yaml:"gemini"`
		} `yaml:"api"`
	}

	// 解析YAML文件
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("无法解析YAML文件：%v", err)
	}

	apiKey := config.Api.Gemini
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := Gemini(tt.args.prompts, tt.args.text, apiKey)
			if err != nil {
				t.Errorf("Gemini() error = %v", err)
			}

			fmt.Println("result: ", text)
		})
	}
}
