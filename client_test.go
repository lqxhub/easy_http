package easy_http

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	//new 一个构造器
	builder := NewClientBuilder()

	//是否跳过服务器证书校验
	builder.SkipVerify(false)

	//设置超时时间
	builder.TimeOut(time.Second * 5)

	//设置代理
	builder.ProxyUrl("http://127.0.0.1:10809")

	//设置根证书
	var certPool [1]string
	certPool[0] = "D:\\server.pem"
	builder.Cert(certPool[:])

	//设置双向校验证书
	var tlsPath [1]*TlsPath
	tlsPath[0] = &TlsPath{
		certFile: "D:\\client.pem",
		keyFile:  "D:\\client.key",
	}
	builder.Tls(tlsPath[:])

	//设置http请求header
	header := make(map[string]string)
	header["Accept-Language"] = "Accept-Language: en,zh"
	builder.Header(header)

	//设置http请求cookie
	cookie := make(map[string]string)
	cookie["name"] = "value"
	builder.Cookie(EasyCookie(cookie))

	//开启cookie jar
	builder.Jar(true, nil)

	//设置 Response 处理函数
	builder.BuildResponse(EasyBuildResponse)

	//构造client
	client, err := builder.Build()
	if err != nil {
		fmt.Println(err)
		return
	}
	response := client.Get("https://baidu.com")
	fmt.Println(response.Error())
	fmt.Println(response.StatusCode())
	fmt.Println(string(response.Content()))
}

func TestClientGetAsyn(t *testing.T) {
	builder := NewClientBuilder()
	client, _ := builder.Build()
	client.GetAsyn("http://baidu.com", call)
	time.Sleep(5 * time.Second)
}

func TestClientGetAsynWithCallback(t *testing.T) {
	builder := NewClientBuilder()
	client, err := builder.Build()
	if err != nil {
		panic(err.Error())
	}

	client.GetAsynWithCallback("http://baidu.com", &Get{})
	time.Sleep(5 * time.Second)
}

func call(response IResponse) {
	fmt.Println(response.Error())
	fmt.Println(string(response.Content()))
}

type Get struct {
}

func (g Get) call(response IResponse) {
	fmt.Println(response.Error())
	fmt.Println(string(response.Content()))
}

func TestClientProxy(t *testing.T) {
	builder := NewClientBuilder()
	builder.ProxyUrl("http://127.0.0.1:10809")
	client, err := builder.Build()
	if err != nil {
		fmt.Println(err)
		return
	}
	response := client.Get("https://google.com")
	fmt.Println(response.Error())
	fmt.Println(response.StatusCode())
	fmt.Println(string(response.Content()))
}

func TestClientStl(t *testing.T) {
	builder := NewClientBuilder()
	//builder.Tls("D:\\client.pem", "D:\\client.key")
	client, err := builder.Build()
	if err != nil {
		fmt.Println(err)
		return
	}
	response := client.Get("https://lqx.com")
	fmt.Println(response.Error())
	fmt.Println(response.StatusCode())
	fmt.Println(string(response.Content()))
}

func TestClientHeader(t *testing.T) {
	builder := NewClientBuilder()
	header := make(map[string]string)
	header["token"] = "ttttt"
	client, err := builder.Header(header).Build()

	if err != nil {
		panic(err)
	}

	header2 := make(map[string]string)
	header2["User-Agent"] = "dsdsdsdsdsd"

	client.Header(header2).Get("http://127.0.0.1:8088")

	client.Get("http://127.0.0.1:8088")
}

func TestClientPostForm(t *testing.T) {
	builder := NewClientBuilder()
	client, err := builder.Build()

	if err != nil {
		panic(err)
	}

	data := make(map[string]string)
	data["name"] = "vvvvv1"
	data["name2"] = "2222222"
	client.PostForm("http://127.0.0.1:8088", EasyPost(data))
}

func TestClientPostMultipart(t *testing.T) {
	builder := NewClientBuilder()
	client, err := builder.Build()

	if err != nil {
		panic(err)
	}

	data := make(map[string]string)
	data["name"] = "vvvvv1"
	data["name2"] = "2222222"

	multipartBuilder := NewMultipartBuilder()
	multipart, err := multipartBuilder.FromDate(data).AddFile("file", "D:\\a.txt").
		AddFile("file2", "D:\\b.log").Builder()

	if err != nil {
		panic(err)
	}
	client.PostMultipart("http://127.0.0.1:8088", multipart)
}

func TestClientSendWithMethod(t *testing.T) {
	builder := NewClientBuilder()
	client, err := builder.Build()
	if err != nil {
		panic(err)
	}
	response := client.SendWithMethod("http://127.0.0.1:8088", http.MethodGet, nil, func(request *http.Request) {

	})
	fmt.Println(response)

}

func TestClientHttpRequest(t *testing.T) {
	builder := NewClientBuilder()
	client, err := builder.Build()
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(http.MethodGet, "http://baidu.com", nil)
	if err != nil {
		panic(err)
	}
	response, err := client.DoRequest(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
	}
}
