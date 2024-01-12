package api

import (
	"fmt"
	"github.com/JoceHYJ/GinLearning/cfg"
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
			name: "ask",
			args: args{prompts: "我遇到了一些问题", text: "外甥女非要在正月理发，身为舅舅的我应该怎么做"},
		},
	}

	config, err := cfg.NewConfig()
	if err != nil {
		t.Errorf("failed to load config, error = %v", err)
		return
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
