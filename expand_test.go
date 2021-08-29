package easy_http

import (
	"fmt"
	"testing"
	"time"
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

	response := client.PostJson("http://127.0.0.1:8088", tt)
	fmt.Println(response.Error())
}

func callback(resp IResponse) {
	if resp.Error() != nil {
		panic(resp.Error())
	}
	fmt.Println(string(resp.Content()))
}

func TestClient_PostJsonAsyn(t *testing.T) {
	builder := NewClientBuilder()
	client, err := builder.Build()
	if err != nil {
		panic(err)
	}
	type T struct {
		Name string `json:"name"`
	}
	tt := &T{Name: "aaaaaaaa"}

	err = client.PostJsonAsyn("http://127.0.0.1:8088", tt, callback)
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
}

type PostJsonCallback struct{}

func (p PostJsonCallback) EasyResponseCallback(resp IResponse) {
	if resp.Error() != nil {
		panic(resp.Error())
	}
	fmt.Println(string(resp.Content()))
}

func TestClient_PostJsonAsynWithCallback(t *testing.T) {
	builder := NewClientBuilder()
	client, err := builder.Build()
	if err != nil {
		panic(err)
	}
	type T struct {
		Name string `json:"name"`
	}
	tt := &T{Name: "aaaaaaaa"}

	err = client.PostJsonAsynWithCallback("http://127.0.0.1:8088", tt, &PostJsonCallback{})
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
}
