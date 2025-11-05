package curl

import (
	"akali/libs/logs"
	"akali/libs/tools"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// Protocol
	ProtocolHttp  = "http"
	ProtocolHttps = "https"

	// headers Key
	HeadersConnection  = "Connection"
	HeadersUserAgent   = "User-Agent"
	HeadersContentType = "Content-Type"

	// ContentType Value
	ContentTypeJson          = "application/json"
	ContentTypeJsonUTF8      = "application/json;charset=utf-8"
	ContentTypeFormUrlEncode = "application/x-www-form-urlencoded"
)

type Curl struct {
	protocol string
	client   *http.Client
	request  *http.Request
	method   string
	path     string
	url      string
	port     string
	timeout  time.Duration
	headers  map[string]string
	cookies  map[string]string
	queries  map[string]any
	body     map[string]any
	traceLog *logs.TraceLog
	tools    *tools.Tools
}

func Initialize(tools *tools.Tools) *Curl {
	// 預設超時
	timeOut := 15 * time.Second

	// 初始化 traceLog
	traceLog := logs.TraceLogInit()

	return &Curl{
		timeout:  timeOut,
		headers:  make(map[string]string),
		cookies:  make(map[string]string),
		queries:  make(map[string]any),
		body:     make(map[string]any),
		traceLog: traceLog,
		tools:    tools,
	}
}

// 新增請求
func (c *Curl) NewRequest(domain, port, path string) *Curl {
	c.client = http.DefaultClient
	c.port = port
	c.path = path
	c.url = fmt.Sprintf("%s://%s:%s%s", c.protocol, domain, port, path)

	c.traceLog.SetDomain(domain)
	return c
}

// 設定請求超時時間
func (c *Curl) SetTimeOut(timeout time.Duration) *Curl {
	if timeout > 0 && timeout < 30*time.Second {
		c.timeout = timeout
	}
	return c
}

// 設定 protocol HTTP 協議
func (c *Curl) SetHttp() *Curl {
	c.protocol = ProtocolHttp
	return c
}

// 設定 protocol HTTPs 協議
func (c *Curl) SetHttps() *Curl {
	c.protocol = ProtocolHttps
	return c
}

// 設定 headers
func (c *Curl) SetHeaders(headers map[string]string) *Curl {
	c.headers = headers
	return c
}

// 設定 cookies
func (c *Curl) SetCookies(cookies map[string]string) *Curl {
	c.cookies = cookies
	return c
}

// 設定 url 查詢參數
func (c *Curl) SetQueries(queries map[string]any) *Curl {
	c.queries = queries
	return c
}

// 設定請求 body 資訊
func (c *Curl) SetBody(bodyData map[string]any) *Curl {
	c.body = bodyData
	return c
}

func (c *Curl) SetTraceID(tid string) *Curl {
	c.traceLog.SetTraceID(tid)
	return c
}

func (c *Curl) Get() (string, error) {
	return c.setMethod(http.MethodGet).send()
}

func (c *Curl) Post() (string, error) {
	return c.setMethod(http.MethodPost).send()
}

func (c *Curl) Put() (string, error) {
	return c.setMethod(http.MethodPut).send()
}

func (c *Curl) Delete() (string, error) {
	return c.setMethod(http.MethodDelete).send()
}

func (c *Curl) send() (string, error) {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("發送 Curl Request 異常, Err: %s", c.tools.PanicParser(err)))
		}
	}()
	var body io.Reader

	if c.tools.InStrArray(c.method, []string{http.MethodPost, http.MethodPut}) && len(c.body) > 0 {
		if contentType, exist := c.headers[HeadersContentType]; exist {
			switch strings.ToLower(contentType) {
			case ContentTypeJson, ContentTypeJsonUTF8:
				if bts, err := json.Marshal(c.body); err != nil {
					return "", err
				} else {
					body = bytes.NewReader(bts)
				}
			case ContentTypeFormUrlEncode:
				formData := url.Values{}

				for k, v := range c.body {
					formData.Add(k, fmt.Sprintf("%v", v))
				}
				body = strings.NewReader(formData.Encode())
			}
		}
	}

	if req, err := http.NewRequest(c.method, c.url, body); err != nil {
		return "", err
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
		c.request = req.WithContext(ctx)
	}
	c.setHeaders().setCookies().setQueries()

	// 寫入 trace log 相關資訊
	c.traceLog.SetTopic("api")
	c.traceLog.SetUrl(c.path)

	if c.tools.InStrArray(c.method, []string{http.MethodGet, http.MethodDelete}) {
		c.traceLog.SetUrl(c.path + "?" + c.request.URL.RawQuery)
	}
	c.traceLog.SetMethod(c.method)
	c.traceLog.SetArgs(c.body)
	c.traceLog.SetHeaders(c.headers)

	// 發送
	var result []byte
	start := time.Now()
	resp, err := c.client.Do(c.request)
	elapsed := time.Since(start)

	if err != nil {
		// 檢查是否為超時
		if errors.Is(err, context.DeadlineExceeded) {
			c.traceLog.PrintCurl(fmt.Sprintf("Timeout after %s", elapsed), err)
		} else if resp == nil {
			c.traceLog.PrintCurl("Response nil", err)
		} else {
			c.traceLog.PrintCurl(resp.Status, err)
		}
		return string(result), err
	}
	// 解析成 json 字串
	result, err = io.ReadAll(resp.Body)

	if err != nil {
		c.traceLog.PrintCurl("Parser response error", err)
	} else {
		c.traceLog.SetResponse(string(result))
		c.traceLog.PrintCurl("Success", nil)
	}

	return string(result), err
}

func (c *Curl) setMethod(method string) *Curl {
	c.method = method
	return c
}

func (c *Curl) setHeaders() *Curl {
	var foundConnection, foundUserAgent bool

	for k, v := range c.headers {
		c.request.Header.Set(k, v)

		switch k {
		case HeadersConnection:
			foundConnection = true
		case HeadersUserAgent:
			foundUserAgent = true
		}
	}
	// 預設 Connection
	if !foundConnection {
		c.request.Header.Set(HeadersConnection, "close")
		c.headers[HeadersConnection] = "close"
	}

	// 預設 User-Agent
	if !foundUserAgent {
		userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36"
		c.request.Header.Set(HeadersUserAgent, userAgent)
		c.headers[HeadersUserAgent] = userAgent
	}
	return c
}

func (c *Curl) setCookies() *Curl {
	for k, v := range c.cookies {
		c.request.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	return c
}

func (c *Curl) setQueries() *Curl {
	q := c.request.URL.Query()

	for k, v := range c.queries {
		q.Add(k, fmt.Sprintf("%v", v))
	}
	c.request.URL.RawQuery = q.Encode()
	return c
}
