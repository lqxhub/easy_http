package easy_http

import (
	bytes2 "bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestClientPostJson(t *testing.T) {
	builder := NewClientBuilder()
	client, err := builder.Build()
	if err != nil {
		panic(err)
	}
	type T struct {
		Name string `json:"name"`
	}
	tt := &T{Name: "aaaaaaaa"}
	bytes, _ := json.Marshal(tt)

	buf := new(bytes2.Buffer)
	buf.Write(bytes)
	client.SendWithMethod("http://127.0.0.1:8088", http.MethodPost, buf, func(req *http.Request) {
		req.Header.Set("Content-Type", HTTP_CONTENT_TYPE_JSON)
	})
}
