package moph

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// ตัวแปร global สำหรับเก็บ config ของ MOPH API
var (
	baseURL   string // URL ของ MOPH API
	clientKey string // Client key สำหรับ authentication
	secretKey string // Secret key สำหรับ authentication
)

// Init ตั้งค่า MOPH API client ด้วย URL และ authentication keys
func Init(url, ck, sk string) {
	baseURL, clientKey, secretKey = url, ck, sk
}

// post ส่ง HTTP POST request ไปยัง MOPH API พร้อม authentication headers
func post(path string, body any) (map[string]any, error) {
	// ถ้า body เป็น nil ให้ใช้ empty object แทน
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
