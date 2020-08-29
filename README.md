# easy_http

#### 介绍
easy_http是对golang的http client的简单封装,不依赖第三方库

1. 提供了链式调用和方便使用的接口

2. easy_http 支持https双向证书校验

3. 支持代理,可以方便的设置代理

4. easy_http 具有很强的拓展性,可以在easy_http的基础上非常方便的定制自己http请求

5. easy_http 封装了multipart 协议, 可以方便的上传文件

6. http Response 支持异步回调

7. 可以方便的自定义Response


#### 安装教程
1. 运行下面的命令

   ``` $ go get -u github.com/bruce12397/easy_http```
2. 把代码引入你的代码

    ```import "github.com/bruce12397/easy_http"```
    
    
######  创建一个客户端
```
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
```
######  post请求
```
builder := NewClientBuilder()
client, err := builder.Build()

if err != nil {
	panic(err)
}

data := make(map[string]string)
data["name"] = "vvvvv1"
data["name2"] = "2222222"
client.PostForm("http://127.0.0.1:8088", EasyPost(data))
```

######  文件上传

```
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
```

######  异步请求

```
//函数回调
func call(response IResponse) {
	fmt.Println(response.Error())
	fmt.Println(string(response.Content()))
}
builder := NewClientBuilder()
client, _ := builder.Build()
client.GetAsyn("http://baidu.com", call)
time.Sleep(5 * time.Second)

//接口回调方式
type Get struct {
}

func (g Get) call(response IResponse) {
	fmt.Println(response.Error())
	fmt.Println(string(response.Content()))
}

builder := NewClientBuilder()
client, err := builder.Build()
if err != nil {
	panic(err.Error())
}

client.GetAsynWithCallback("http://baidu.com", &Get{})
```

###### 为这一次请求 设置 cookie和 header

```
builder := NewClientBuilder()
client, _ := builder.Build()
cookie := make(map[string]string)
cookie["name"] = "value"

header := make(map[string]string)
header["Accept-Language"] = "Accept-Language: en,zh"

client.Cookies(EasyCookie(cookie)).Header(header).Get("http://127.0.0.1:8088/")
```

###### 自定义请求

```
builder := NewClientBuilder()
client, err := builder.Build()
if err != nil {
	panic(err)
}
response := client.SendWithMethod("http://127.0.0.1:8088", http.MethodGet, nil, func(request *http.Request) {

})
fmt.Println(response)
```

###### 使用http 原生的 请求

```
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
```

###### easy函数 使用

```
//构造 post form 的 value
values := make(map[string]string)
values["a"] = "aaaaaa"
values["b"] = "bbbbbb"
values["c"] = "中文"
EasyPost(values).Encode()

//构造 cookie 
cookie := make(map[string]string)
cookie["name"] = "value"
EasyCookie(cookie)

//构造 post 上传文件的 multipart
multipartBuilder := NewMultipartBuilder()
multipartBuilder.AddFile("file1", "d:\\a.txt")
multipartBuilder.AddFromDate("name", "value")
builder, err := multipartBuilder.Builder()

```

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request