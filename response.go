package easy_http

import (
	"io/ioutil"
	"net/http"
)

type BuildResponse func(resp *http.Response, err error) IResponse

//使用client发去请求后,返回一个实现了这个接口的对象
//只要实现这个接口,就能作为返回值
//在 `BuildResponse` 函数中构造出返回的对象
//默认提供了 `HttpResponse` 实现了这个接口
//可以根据自己的需求自己重新实现这个接口
type IResponse interface {
	//返回这个请求的错误
	Error() error

	//返回这个请求的http状态码
	StatusCode() int

	//返回HTTP请求的header信息
	Header() http.Header

	//返回HTTP内容长度
	ContentLength() int64

	//返回HTTP的内容
	Content() []byte

	//返回HTTP包中的 response信息
	Resp() *http.Response

	//返回这次请求的request信息
	Request() *http.Request

	//根据name 返回response的cookie
	Cookie(name string) *http.Cookie
}

type HttpResponse struct {
	err             error
	ResponseContent []byte
	httpResp        *http.Response
}

func (h *HttpResponse) Error() error {
	return h.err
}

func (h *HttpResponse) StatusCode() int {
	if h.httpResp == nil {
		return 0
	}
	return h.httpResp.StatusCode
}

func (h *HttpResponse) Header() http.Header {
	if h.httpResp == nil {
		return nil
	}
	return h.httpResp.Header
}

func (h *HttpResponse) ContentLength() int64 {
	if h.httpResp == nil {
		return 0
	}
	return h.httpResp.ContentLength
}

func (h *HttpResponse) Content() []byte {
	return h.ResponseContent
}

func (h *HttpResponse) Resp() *http.Response {
	return h.httpResp
}

func (h *HttpResponse) Request() *http.Request {
	if h.httpResp == nil {
		return nil
	}
	return h.httpResp.Request
}

func (h *HttpResponse) Cookie(name string) *http.Cookie {
	for _, cookie := range h.httpResp.Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

//默认构造HTTP response的函数
func EasyBuildResponse(resp *http.Response, err error) IResponse {
	response := new(HttpResponse)
	if err != nil {
		response.err = err
		return response
	}
	response.httpResp = resp

	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.err = err
		return response
	}
	response.ResponseContent = all
	_ = resp.Body.Close()
	return response
}
