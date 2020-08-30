package easy_http

import (
	"net/http"
	"net/url"
	"strings"
)

//异步回调的接口
type ICallBack interface {
	EasyResponseCallback(IResponse)
}

//request重定向的回调函数
type CheckRedirect func(req *http.Request, via []*http.Request) error

type TlsPath struct {
	//cert (pem) 路径
	CertFile string
	//key 路径
	KeyFile string
}

//构造一个简单的HTTP请求cookie
func EasyCookie(simple map[string]string) []*http.Cookie {
	if len(simple) == 0 {
		return nil
	}
	cookies := make([]*http.Cookie, 0, len(simple))
	for k, v := range simple {
		cookies = append(cookies, &http.Cookie{
			Name:  k,
			Value: v,
		})
	}
	return cookies
}

//构造一个简单的GET请求协议
func EasyGet(strUrl string, values map[string]string) string {
	if strUrl == "" || values == nil {
		return strUrl
	}
	var buf strings.Builder
	buf.WriteString(strUrl)
	buf.WriteByte('?')
	i := 0
	for k, v := range values {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(v))
		i++
	}
	return buf.String()
}

//构造一个简单的POST body,
func EasyPost(values map[string]string) url.Values {
	if values == nil {
		return nil
	}
	value := make(url.Values, len(values))
	for k, v := range values {
		value.Add(k, v)
	}
	return value
}

//POST请求中,处理request的函数
func EasyPostFromRequest(r *http.Request) {
	r.Header.Set("Content-Type", HTTP_CONTENT_TYPE_FROM_DATA)
}
