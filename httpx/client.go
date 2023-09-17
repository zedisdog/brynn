package httpx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type IClient interface {
	SetHeader(header http.Header) IClient
	AddHeader(key string, value string) IClient
	Debug() IClient
	SetTimeout(d time.Duration) IClient
	SetQuery(values url.Values) IClient
	AddQuery(key string, value string) IClient
	SetCookies(cookies []*http.Cookie) IClient
	AddCookie(cookie *http.Cookie) IClient
	Get(url string) getResult
	Delete(url string) getResult
	PutJson(url string, data any) getResult
	PutXml(url string, data any) getResult
	PostJson(url string, data any) getResult
	PostXml(url string, data any) getResult
	PostForm(u string, data any) getResult
	Request(req *http.Request) (resp *http.Response, err error)
}

type getResult interface {
	Response() (res *http.Response, err error)
	ScanJsonBody(data any) (err error)
}

type OptFunc func(*client)

func WithBaseUrl(baseUrl string) OptFunc {
	return func(client *client) {
		client.BaseUrl = baseUrl
	}
}

func WithTimeout(d time.Duration) OptFunc {
	return func(c *client) {
		c.timeout = d
	}
}

func NewClient(opts ...OptFunc) IClient {
	c := &client{
		header:  make(http.Header),
		query:   make(url.Values),
		cookies: make([]*http.Cookie, 0),
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

type client struct {
	BaseUrl string
	debug   bool
	header  http.Header
	timeout time.Duration
	query   url.Values
	res     *http.Response
	Err     error
	cookies []*http.Cookie
}

func (c client) SetHeader(header http.Header) IClient {
	c.header = header
	return &c
}

func (c client) AddHeader(key string, value string) IClient {
	c.header.Add(key, value)
	return &c
}

func (c client) Debug() IClient {
	c.debug = true
	return &c
}

func (c client) SetTimeout(d time.Duration) IClient {
	c.timeout = d
	return &c
}

func (c client) SetQuery(values url.Values) IClient {
	c.query = values
	return &c
}

func (c client) AddQuery(key string, value string) IClient {
	c.query.Add(key, value)
	return &c
}

func (c client) SetCookies(cookies []*http.Cookie) IClient {
	c.cookies = cookies
	return &c
}

func (c client) AddCookie(cookie *http.Cookie) IClient {
	c.cookies = append(c.cookies, cookie)
	return &c
}

func (c client) Get(url string) getResult {
	c.get(url)
	return &c
}

func (c *client) get(url string) {
	var req *http.Request
	req, c.Err = c.makeRequest(http.MethodGet, url, nil)
	if c.Err != nil {
		return
	}

	c.res, c.Err = c.Request(req)
	return
}

func (c client) Delete(url string) getResult {
	c.delete(url)
	return &c
}

func (c *client) delete(url string) {
	var req *http.Request
	req, c.Err = c.makeRequest(http.MethodDelete, url, nil)
	if c.Err != nil {
		return
	}

	c.res, c.Err = c.Request(req)
	return
}

func (c client) PutJson(url string, data any) getResult {
	var (
		content []byte
	)
	content, c.Err = json.Marshal(data)
	if c.Err != nil {
		return &c
	}
	c.header.Add("Content-Type", "application/json")
	c.put(url, bytes.NewReader(content))
	return &c
}

func (c client) PutXml(url string, data any) getResult {
	var content []byte
	content, c.Err = xml.Marshal(data)
	if c.Err != nil {
		return &c
	}
	c.header.Add("Content-Type", "text/xml")
	c.put(url, bytes.NewReader(content))
	return &c
}

func (c *client) put(url string, body io.Reader) {
	var req *http.Request
	req, c.Err = c.makeRequest(http.MethodPut, url, body)
	if c.Err != nil {
		return
	}

	c.res, c.Err = c.Request(req)
	return
}

func (c client) PostJson(url string, data any) getResult {
	var content []byte
	content, c.Err = json.Marshal(data)
	if c.Err != nil {
		return &c
	}

	c.header.Add("Content-Type", "application/json")
	c.post(url, bytes.NewReader(content))
	return &c
}

func (c client) PostXml(url string, data any) getResult {
	var content []byte
	content, c.Err = xml.Marshal(data)
	if c.Err != nil {
		return &c
	}
	c.header.Add("Content-Type", "text/xml")
	c.post(url, bytes.NewReader(content))
	return &c
}

func (c client) PostForm(u string, data any) getResult {
	c.header.Add("Content-Type", "application/x-www-form-urlencoded")
	vData := reflect.ValueOf(data)
	if vData.Kind() == reflect.Pointer {
		vData = vData.Elem()
	}
	if vData.Kind() == reflect.String {
		c.post(u, strings.NewReader(vData.Interface().(string)))
		return &c
	}

	if vData.Kind() != reflect.Struct {
		c.Err = errors.New("unsupported type of data")
		return &c
	}

	tData := reflect.TypeOf(data)
	if tData.Kind() == reflect.Pointer {
		tData = tData.Elem()
	}

	values := url.Values{}
	for i := 0; i < vData.NumField(); i++ {
		fieldValue := vData.Field(i)
		if fieldValue.IsValid() {
			continue
		}
		fieldType := tData.Field(i)

		var (
			key   string
			value string
		)
		parts := strings.Split(fieldType.Tag.Get("form"), ",")
		if len(parts) > 0 {
			key = parts[0]
		} else {
			key = fieldType.Name
		}

		value = gconv.String(fieldValue.Interface())

		values.Add(key, value)
	}

	c.post(u, strings.NewReader(values.Encode()))
	return &c
}

func (c *client) post(url string, body io.Reader) {
	req, err := c.makeRequest(http.MethodPost, url, body)
	if err != nil {
		return
	}

	c.res, c.Err = c.Request(req)
	return
}

func (c client) buildUrl(u string) string {
	if !strings.HasPrefix(u, "http") {
		u = fmt.Sprintf(
			"%s/%s",
			strings.TrimRight(c.BaseUrl, "/"),
			strings.TrimLeft(u, "/"),
		)
	}

	var parsedUrl *url.URL
	parsedUrl, c.Err = url.Parse(u)
	if c.Err != nil {
		return ""
	}

	for k, values := range c.query {
		for _, v := range values {
			parsedUrl.Query().Add(k, v)
		}
	}

	return parsedUrl.String()
}

func (c client) makeRequest(method string, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, c.buildUrl(url), body)
	if err != nil {
		return
	}

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	for k, values := range c.header {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	return
}

func (c client) Request(req *http.Request) (resp *http.Response, err error) {
	client := http.Client{
		Timeout: c.timeout,
	}

	var start time.Time
	if c.debug {
		var content []byte
		content, err = io.ReadAll(req.Body)
		if err != nil {
			return
		}
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(content))
		fmt.Printf(
			"[request: method=%s, url=%s, body=%s, header=%v]\n",
			req.Method,
			req.URL.String(),
			string(content),
			req.Header,
		)
		start = time.Now()
	}

	resp, err = client.Do(req)
	if err != nil {
		return
	}

	if c.debug {
		fmt.Printf("[request time: %fs]\n", time.Since(start).Seconds())

		var content []byte
		content, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return
		}
		resp.Body = io.NopCloser(bytes.NewReader(content))
		fmt.Printf(
			"[response: body=%s]\n",
			string(content),
		)
	}

	return
}

func (c client) Response() (res *http.Response, err error) {
	return c.res, c.Err
}

func (c client) ScanJsonBody(data any) (err error) {
	if c.Err != nil {
		return c.Err
	}
	defer c.res.Body.Close()
	content, err := io.ReadAll(c.res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, data)
	return
}
