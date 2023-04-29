package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func request(r *gin.Engine, method string, url string, payload map[string]any, headers map[string]string) []byte {
	w := httptest.NewRecorder()
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	r.ServeHTTP(w, req)

	return w.Body.Bytes()
}

func decodeJSON[T any](body []byte) T {
	var decoded T
	json.Unmarshal(body, &decoded)
	return decoded
}
