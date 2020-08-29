package easy_http

import (
	"fmt"
	"testing"
)

func TestEasyGet(t *testing.T) {
	strUrl := "http://www.baidu.com"

	values := make(map[string]string)
	values["a"] = "aaaaaa"
	values["b"] = "bbbbbb"
	values["c"] = "中文"

	fmt.Println(EasyGet(strUrl, values))
}

func TestEasyPost(t *testing.T) {

	values := make(map[string]string)
	values["a"] = "aaaaaa"
	values["b"] = "bbbbbb"
	values["c"] = "中文"
	fmt.Println(EasyPost(values).Encode())
}

func TestEasyCookie(t *testing.T) {
	builder := NewClientBuilder()
	cookie := make(map[string]string)
	cookie["name"] = "value"
	client, err := builder.Cookie(EasyCookie(cookie)).Build()
	if err != nil {
		panic(err)
	}
	cookie2 := make(map[string]string)
	cookie2["name22"] = "value222"
	client.Cookies(EasyCookie(cookie2)).Get("http://127.0.0.1:8088/")
	client.Get("http://127.0.0.1:8088/")
}

func TestEasyMultipart(t *testing.T) {
	multipartBuilder := NewMultipartBuilder()
	multipartBuilder.AddFile("file1", "d:\\a.txt")
	multipartBuilder.AddFromDate("name", "value")
	builder, err := multipartBuilder.Builder()
	fmt.Println(builder, err)
}
