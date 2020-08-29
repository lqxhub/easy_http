package easy_http

//User-agent
const (
	//chrome pc
	HTTP_USER_AGENT_CHROME_PC = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11"
	//chrome mobile
	HTTP_USER_AGENT_CHROME_MOBILE = "Mozilla/5.0 (Linux; U; Android 2.2.1; zh-cn; HTC_Wildfire_A3333 Build/FRG83D) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"

	//Firefox pc
	HTTP_USER_AGENT_FIREFOX_PC = "Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1"
	//Firefox Mobil
	HTTP_USER_AGENT_FIREFOX_MOBLIE = "Mozilla/5.0 (Androdi; Linux armv7l; rv:5.0) Gecko/ Firefox/5.0 fennec/5.0"

	//IE11
	HTTP_USER_AGENT_IE11 = "Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv 11.0) like Gecko"
)

//Content-Type
const (
	HTTP_CONTENT_TYPE_FROM_DATA = "application/x-www-form-urlencoded"
	HTTP_CONTENT_TYPE_TEXT      = "text/plain"
	HTTP_CONTENT_TYPE_JSON      = "application/json"
	HTTP_CONTENT_TYPE_XML       = "application/xml"
)
