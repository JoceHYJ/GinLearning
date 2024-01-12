package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Payload struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type Candidate struct {
	Content struct {
		Parts []Part `json:"parts"`
		Role  string `json:"role"`
	} `json:"content"`
	FinishReason  string         `json:"finishReason"`
	Index         int            `json:"index"`
	SafetyRatings []SafetyRating `json:"safetyRatings"`
}

type SafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

type Result struct {
	Candidates     []Candidate `json:"candidates"`
	PromptFeedback struct {
		SafetyRatings []SafetyRating `json:"safetyRatings"`
	} `json:"promptFeedback"`
}

func Gemini(prompts, text, apiKey string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", apiKey)
	method := "POST"

	payload := Payload{Contents: []Content{
		{Parts: []Part{
			{Text: fmt.Sprintf("%s%s", prompts, text)},
		}},
	}}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(payloadJSON))

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = res.Body.Close() }()

	decoder := json.NewDecoder(res.Body)
	decoder.DisallowUnknownFields()
	result := new(Result)
	err = decoder.Decode(result)
	if err != nil {
		return "", err
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
