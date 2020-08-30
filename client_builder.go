package easy_http

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

//初始化 clientBuilder
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		skipVerify:    true,
		buildResponse: EasyBuildResponse,
	}
}

//client 构造器 用来初始化一个client
//不建议直接new client使用
type ClientBuilder struct {

	//超时时间
	timeOut time.Duration

	//代理http url
	proxy string

	//tls私钥证书
	tlsPath []*TlsPath

	//cert root 证书
	certPool []string
	//是否跳过HTTPS证书校验(默认跳过)
	skipVerify bool

	//client发起HTTP请求时,header信息
	//默认会携带User-agent信息
	header map[string]string

	//client发起 HTTP请求时,自动携带的cookie
	cookie []*http.Cookie

	//重定向函数
	checkRedirect CheckRedirect

	//client是否开启cookieJar功能
	//默认不开启
	openJar bool

	//jarOptions 配置
	jarOptions *cookiejar.Options

	//处理HTTP的 response的回调函数
	//默认使用 `response.go`中的 `BuildResponse` 函数
	buildResponse BuildResponse
}

func (builder *ClientBuilder) TimeOut(t time.Duration) *ClientBuilder {
	builder.timeOut = t
	return builder
}

func (builder *ClientBuilder) ProxyUrl(u string) *ClientBuilder {
	builder.proxy = u
	return builder
}

func (builder *ClientBuilder) Tls(tlsPath []*TlsPath) *ClientBuilder {
	builder.tlsPath = tlsPath
	return builder
}

func (builder *ClientBuilder) Cert(cert []string) *ClientBuilder {
	builder.certPool = cert
	return builder
}

func (builder *ClientBuilder) SkipVerify(skip bool) *ClientBuilder {
	builder.skipVerify = skip
	return builder
}

func (builder *ClientBuilder) Cookie(cookie []*http.Cookie) *ClientBuilder {
	builder.cookie = cookie
	return builder
}

func (builder *ClientBuilder) CheckRedirect(checkRedirect CheckRedirect) *ClientBuilder {
	builder.checkRedirect = checkRedirect
	return builder
}

func (builder *ClientBuilder) Header(header map[string]string) *ClientBuilder {
	builder.header = header
	return builder
}

func (builder *ClientBuilder) Jar(open bool, options *cookiejar.Options) *ClientBuilder {
	builder.openJar = open
	builder.jarOptions = options
	return builder
}

func (builder *ClientBuilder) BuildResponse(build BuildResponse) *ClientBuilder {
	builder.buildResponse = build
	return builder
}

//构造 client
func (builder *ClientBuilder) Build() (*Client, error) {
	var (
		err         error
		proxy       *url.URL
		x509KeyPair tls.Certificate
	)

	if builder.buildResponse == nil {
		return nil, errors.New("clint not set BuildResponse")
	}

	if builder.proxy != "" {
		proxy, err = url.Parse(builder.proxy)
		if err != nil {
			return nil, err
		}
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: builder.skipVerify,
	}
	if builder.tlsPath != nil {
		certificates := make([]tls.Certificate, len(builder.tlsPath))
		for i, path := range builder.tlsPath {
			x509KeyPair, err = tls.LoadX509KeyPair(path.CertFile, path.KeyFile)
			if err != nil {
				return nil, err
			}
			certificates[i] = x509KeyPair
		}
		tlsConfig.Certificates = certificates
	}

	if builder.certPool != nil {
		tlsConfig.RootCAs = x509.NewCertPool()
		for _, certFile := range builder.certPool {
			if ca, err := ioutil.ReadFile(certFile); err != nil {
				return nil, err
			} else {
				if ok := tlsConfig.RootCAs.AppendCertsFromPEM(ca); !ok {
					return nil, fmt.Errorf("load:%s cert fail", certFile)
				}
			}
		}
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	if proxy != nil {
		transport.Proxy = http.ProxyURL(proxy)
	}
	c := &Client{
		client: &http.Client{
			Transport:     transport,
			Timeout:       builder.timeOut,
			CheckRedirect: builder.checkRedirect,
		},
		header:        builder.header,
		cookies:       builder.cookie,
		buildResponse: builder.buildResponse,
	}

	if builder.openJar {
		jar, err := cookiejar.New(builder.jarOptions)
		if err != nil {
			return nil, err
		}
		c.client.Jar = jar
	}
	return c, nil
}
