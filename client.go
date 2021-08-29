package easy_http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {

	//超时时间
	timeOut time.Duration

	//HTTP请求中的header信息
	header map[string]string

	//HTTP请求中,携带的cookies
	cookies []*http.Cookie

	//发起请求的client(go 自带的client)
	client *http.Client

	//临时header,每次请求后会重置
	_header map[string]string

	//临时cookie,每次请求后会重置
	_cookies []*http.Cookie

	//处理HTTP返回的response
	buildResponse BuildResponse
}

//为这次请求设置header,只有这次会生效
func (c *Client) Header(header map[string]string) *Client {
	c._header = header
	return c
}

//为这次请求设置cookie,只有这次会生效
func (c *Client) Cookies(cookies []*http.Cookie) *Client {
	c._cookies = cookies
	return c
}

//为client添加header,原来的会保留,整个生命周期有效
func (c *Client) AddHeader(header map[string]string) {
	for k, v := range header {
		c.header[k] = v
	}
}

//为client设置新的header,原来的会删除,整个生命周期有效
func (c *Client) SetHeader(header map[string]string) {
	c.header = header
}

//为client添加cookie,原来的会保留,整个生命周期有效
func (c *Client) AddCookies(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		c.cookies = append(c.cookies, cookie)
	}
}

//为client设置新的cookie,原来的会删除,整个生命周期有效
func (c *Client) SetCookies(cookies []*http.Cookie) {
	c.cookies = cookies
}

//发起HTTP请求
func (c *Client) DoRequest(r *http.Request) (*http.Response, error) {
	return c.client.Do(r)
}

//指定请求的方法,发送请求
//`req` 参数 可以处理这次请求的request
func (c *Client) SendWithMethod(url, method string, body io.Reader, req func(request *http.Request)) IResponse {
	request, err := c.getRequest(method, url, body)
	if err != nil {
		return c.buildResponse(nil, err)
	}
	if req != nil {
		req(request)
	}
	return c.buildResponse(c.DoRequest(request))
}

//使用异步回调的方式,指定请求的方法,发送请求
//`req` 参数 可以处理这次请求的request
//`call` 参数,请求成功后的回调函数
func (c *Client) SendWithMethodCallBack(url, method string, body io.Reader, req func(request *http.Request), call func(response IResponse)) error {
	request, err := c.getRequest(method, url, body)
	if err != nil {
		return err
	}
	if req != nil {
		req(request)
	}
	go func() {
		call(c.buildResponse(c.DoRequest(request)))
	}()
	return nil
}

//发起GET 请求
func (c *Client) Get(url string) IResponse {
	return c.SendWithMethod(url, http.MethodGet, nil, nil)
}

//GET 异步请求,使用回调函数
func (c *Client) GetAsyn(url string, call func(response IResponse)) error {
	return c.SendWithMethodCallBack(url, http.MethodGet, nil, nil, call)
}

//GET 异步请求,使用接口回调
func (c *Client) GetAsynWithCallback(url string, call ICallBack) error {
	return c.GetAsyn(url, call.EasyResponseCallback)
}

//post 的form请求
func (c *Client) PostForm(url string, values url.Values) IResponse {
	var reader io.Reader
	if values != nil {
		reader = strings.NewReader(values.Encode())
	}
	return c.SendWithMethod(url, http.MethodPost, reader, EasyPostFromRequest)
}

//Post form 异步请求,使用回调函数
func (c *Client) PostFormAsyn(url string, values url.Values, call func(response IResponse)) error {
	if call == nil {
		return errors.New("callback function is nil")
	}
	if values == nil {
		return errors.New("values is nil")
	}
	reader := strings.NewReader(values.Encode())
	return c.SendWithMethodCallBack(url, http.MethodPost, reader, EasyPostFromRequest, call)
}

//Post form 异步请求,使用接口回调
func (c *Client) PostFormAsynWithCallback(url string, values url.Values, call ICallBack) error {
	return c.PostFormAsyn(url, values, call.EasyResponseCallback)
}

//post 的bytes请求
func (c *Client) PostBytes(url string, value []byte, req func(request *http.Request)) IResponse {
	if value == nil {
		return c.buildResponse(nil, errors.New("PostBytes value is nil"))
	}
	reader := bytes.NewReader(value)
	return c.SendWithMethod(url, http.MethodPost, reader, req)
}

//post 的bytes请求
func (c *Client) PostBytesAsyn(url string, value []byte, req func(request *http.Request), call func(response IResponse)) error {
	if call == nil {
		return errors.New("callback function is nil")
	}
	if value == nil {
		return errors.New("value is nil")
	}
	reader := bytes.NewReader(value)
	return c.SendWithMethodCallBack(url, http.MethodPost, reader, req, call)
}

//post 的json请求
func (c *Client) PostJson(url string, value interface{}) IResponse {
	if value == nil {
		return c.buildResponse(nil, errors.New("PostJson value is nil"))
	}
	by, err := json.Marshal(value)
	if err != nil {
		return c.buildResponse(nil, err)
	}
	return c.PostBytes(url, by, EasyPostJsonRequest)
}

//Post json 异步请求,使用回调函数
func (c *Client) PostJsonAsyn(url string, value interface{}, call func(response IResponse)) error {
	if call == nil {
		return errors.New("callback function is nil")
	}
	if value == nil {
		return errors.New("value is nil")
	}
	by, err := json.Marshal(value)
	if err != nil {
		return errors.New("value json encode error: " + err.Error())
	}
	return c.PostBytesAsyn(url, by, EasyPostJsonRequest, call)
}

//Post json 异步请求,使用接口回调
func (c *Client) PostJsonAsynWithCallback(url string, values interface{}, call ICallBack) error {
	return c.PostJsonAsyn(url, values, call.EasyResponseCallback)
}

//post 的xml请求
func (c *Client) PostXml(url string, value interface{}) IResponse {
	if value == nil {
		return c.buildResponse(nil, errors.New("PostJson value is nil"))
	}
	by, err := xml.Marshal(value)
	if err != nil {
		return c.buildResponse(nil, err)
	}
	return c.PostBytes(url, by, EasyPostXmlRequest)
}

//Post xml 异步请求,使用回调函数
func (c *Client) PostXmlAsyn(url string, value interface{}, call func(response IResponse)) error {
	if call == nil {
		return errors.New("callback function is nil")
	}
	if value == nil {
		return errors.New("value is nil")
	}
	by, err := json.Marshal(value)
	if err != nil {
		return errors.New("value json encode error: " + err.Error())
	}
	return c.PostBytesAsyn(url, by, EasyPostXmlRequest, call)
}

//Post xml 异步请求,使用接口回调
func (c *Client) PostXmlAsynWithCallback(url string, values interface{}, call ICallBack) error {
	return c.PostXmlAsyn(url, values, call.EasyResponseCallback)
}

//post 的multipart请求
func (c *Client) PostMultipart(url string, body IMultipart) IResponse {
	return c.SendWithMethod(url, http.MethodPost, body, func(request *http.Request) {
		request.Header.Set("Content-Type", body.ContentType())
	})
}

//post 的multipart请求,使用回调函数
func (c *Client) PostMultipartAsyn(url string, body IMultipart, call func(response IResponse)) error {
	if call == nil {
		return errors.New("callback function is nil")
	}
	return c.SendWithMethodCallBack(url, http.MethodPost, body, func(request *http.Request) {
		request.Header.Set("Content-Type", body.ContentType())
	}, call)
}

//post 的multipart请求,使用接口回调
func (c *Client) PostMultipartAsynWithCallback(url string, body IMultipart, call ICallBack) error {
	return c.PostMultipartAsyn(url, body, call.EasyResponseCallback)
}

//初始化一个 http.Request, 并填充属性
func (c *Client) getRequest(method, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range c.header {
		request.Header.Set(k, v)
	}
	for k, v := range c._header {
		request.Header.Set(k, v)
	}
	c._header = nil

	if _, e := request.Header["User-Agent"]; !e {
		request.Header.Set("User-Agent", HTTP_USER_AGENT_CHROME_PC)
	}

	for _, v := range c.cookies {
		request.AddCookie(v)
	}
	for _, v := range c._cookies {
		request.AddCookie(v)
	}
	c._cookies = nil

	return request, nil
}
