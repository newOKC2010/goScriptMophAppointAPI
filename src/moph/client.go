package moph

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

var baseURL, clientKey, secretKey string

func Init(url, ck, sk string) {
	baseURL, clientKey, secretKey = url, ck, sk
}

func post(path string, body any) (map[string]any, error) {
	if body == nil {
		body = map[string]any{}
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", baseURL+path, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("client-key", clientKey)
	req.Header.Set("secret-key", secretKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var result map[string]any
	json.Unmarshal(raw, &result)
	return result, nil
}
