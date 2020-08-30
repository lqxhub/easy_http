# easy_http

#### Description
easy_http is a simple encapsulation of golang's http client and does not rely on third-party libraries

1. provides chain call and convenient interface

2. easy_http supports https two-way certificate verification

3. support proxy you can set proxy easily

4. easy_http With strong scalability, you can customize your own http request very conveniently on the basis of easy_http

5. easy_http encapsulates the multipart protocol, files can be uploaded easily

6. http Response support asynchronous callback

7. response can be easily customized


#### Installation
1.  you can use the below Go command to install

    ``` $ go get -u github.com/bruce12397/easy_http```
    
2. Import it in your code
    
   ```import "github.com/bruce12397/easy_http"```
        
        
#### Instructions

###### create a client

```
//new A constructor
builder := NewClientBuilder()

//Whether to skip server certificate verification
builder.SkipVerify(false)

//set timeout
builder.TimeOut(time.Second * 5)

//set http proxy
builder.ProxyUrl("http://127.0.0.1:10809")

//set the root certificate
var certPool [1]string
certPool[0] = "D:\\server.pem"
builder.Cert(certPool[:])

//set up two way verification certificate
var tlsPath [1]*TlsPath
tlsPath[0] = &TlsPath{
    CertFile: "D:\\client.pem",
    KeyFile:  "D:\\client.key",
}

builder.Tls(tlsPath[:])

//set http request header
header := make(map[string]string)
header["Accept-Language"] = "Accept-Language: en,zh"
builder.Header(header)

//set http request cookie
cookie := make(map[string]string)
cookie["name"] = "value"
builder.Cookie(EasyCookie(cookie))

//open cookie jar
builder.Jar(true, nil)

//set response processing function
builder.BuildResponse(EasyBuildResponse)

//build client
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

######  post request

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

######  file upload

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

######  asynchronous request

```
//function callback
func call(response IResponse) {
	fmt.Println(response.Error())
	fmt.Println(string(response.Content()))
}
builder := NewClientBuilder()
client, _ := builder.Build()
client.GetAsyn("http://baidu.com", call)
time.Sleep(5 * time.Second)

//interface callback method
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


###### Set cookie and header for this request

```
builder := NewClientBuilder()
client, _ := builder.Build()
cookie := make(map[string]string)
cookie["name"] = "value"

header := make(map[string]string)
header["Accept-Language"] = "Accept-Language: en,zh"

client.Cookies(EasyCookie(cookie)).Header(header).Get("http://127.0.0.1:8088/")
```

###### Custom request

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

###### Use http native request

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

###### easy function use

```
//Construct the value of the post form
values := make(map[string]string)
values["a"] = "aaaaaa"
values["b"] = "bbbbbb"
values["c"] = "中文"
EasyPost(values).Encode()

//Construct cookie 
cookie := make(map[string]string)
cookie["name"] = "value"
EasyCookie(cookie)

//Construct a multipart of post uploaded files
multipartBuilder := NewMultipartBuilder()
multipartBuilder.AddFile("file1", "d:\\a.txt")
multipartBuilder.AddFromDate("name", "value")
builder, err := multipartBuilder.Builder()

```

#### Contribution

1.  Fork the repository
2.  Create Feat_xxx branch
3.  Commit your code
4.  Create Pull Request