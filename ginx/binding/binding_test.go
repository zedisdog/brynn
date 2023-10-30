package binding

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderBinding(t *testing.T) {
	h := Header
	assert.Equal(t, "header", h.Name())

	type tHeader struct {
		Limit int `header:"limit"`
	}

	var theader tHeader
	req := requestWithBody("GET", "/", "")
	req.Header.Add("limit", "1000")
	assert.NoError(t, h.Bind(req, &theader))
	assert.Equal(t, 1000, theader.Limit)

	req = requestWithBody("GET", "/", "")
	assert.Error(t, h.Bind(req, &theader))

	req = requestWithBody("GET", "/", "")
	req.Header.Add("fail", `{fail:fail}`)

	type failStruct struct {
		Fail map[string]any `header:"fail"`
	}

	err := h.Bind(req, &failStruct{})
	assert.Error(t, err)
}

func TestBindingJSONNilBody(t *testing.T) {
	var obj FooStruct
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	err := JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func TestBindingJSON(t *testing.T) {
	testBodyBinding(t,
		JSON, "json",
		"/", "/",
		`{"foo": "bar"}`, `{"bar": "foo"}`)
}

func TestBindingJSONSlice(t *testing.T) {
	EnableDecoderDisallowUnknownFields = true
	defer func() {
		EnableDecoderDisallowUnknownFields = false
	}()

	testBodyBindingSlice(t, JSON, "json", "/", "/", `[]`, ``)
	testBodyBindingSlice(t, JSON, "json", "/", "/", `[{"foo": "123"}]`, `[{}]`)
	testBodyBindingSlice(t, JSON, "json", "/", "/", `[{"foo": "123"}]`, `[{"foo": ""}]`)
	testBodyBindingSlice(t, JSON, "json", "/", "/", `[{"foo": "123"}]`, `[{"foo": 123}]`)
	testBodyBindingSlice(t, JSON, "json", "/", "/", `[{"foo": "123"}]`, `[{"bar": 123}]`)
	testBodyBindingSlice(t, JSON, "json", "/", "/", `[{"foo": "123"}]`, `[{"foo": "123456789012345678901234567890123"}]`)
}

func TestBindingJSONUseNumber(t *testing.T) {
	testBodyBindingUseNumber(t,
		JSON, "json",
		"/", "/",
		`{"foo": 123}`, `{"bar": "foo"}`)
}

func TestBindingJSONUseNumber2(t *testing.T) {
	testBodyBindingUseNumber2(t,
		JSON, "json",
		"/", "/",
		`{"foo": 123}`, `{"bar": "foo"}`)
}

func TestBindingJSONDisallowUnknownFields(t *testing.T) {
	t.Skip("not implement")
	testBodyBindingDisallowUnknownFields(t, JSON,
		"/", "/",
		`{"foo": "bar"}`, `{"foo": "bar", "what": "this"}`)
}

func TestBindingJSONStringMap(t *testing.T) {
	testBodyBindingStringMap(t, JSON,
		"/", "/",
		`{"foo": "bar", "hello": "world"}`, `{"num": 2}`)
}

type FooStruct struct {
	Foo string `msgpack:"foo" json:"foo" form:"foo" xml:"foo" binding:"required,max=32"`
}

func testBodyBinding(t *testing.T, b binding.Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)

	obj = FooStruct{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

func testBodyBindingSlice(t *testing.T, b binding.Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	var obj1 []FooStruct
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj1)
	require.NoError(t, err)
	require.NotNil(t, obj1)

	var obj2 []FooStruct
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj2)
	require.Error(t, err)
}

type FooStructUseNumber struct {
	Foo any `json:"foo" binding:"required"`
}

func testBodyBindingUseNumber(t *testing.T, b binding.Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStructUseNumber{}
	req := requestWithBody("POST", path, body)
	EnableDecoderUseNumber = true
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	// we hope it is int64(123)
	v, e := obj.Foo.(json.Number).Int64()
	assert.NoError(t, e)
	assert.Equal(t, int64(123), v)

	obj = FooStructUseNumber{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}
func testBodyBindingUseNumber2(t *testing.T, b binding.Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStructUseNumber{}
	req := requestWithBody("POST", path, body)
	EnableDecoderUseNumber = false
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	// it will return float64(123) if not use EnableDecoderUseNumber
	// maybe it is not hoped
	assert.Equal(t, float64(123), obj.Foo)

	obj = FooStructUseNumber{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

type FooStructDisallowUnknownFields struct {
	Foo any `json:"foo" binding:"required"`
}

func testBodyBindingDisallowUnknownFields(t *testing.T, b binding.Binding, path, badPath, body, badBody string) {
	EnableDecoderDisallowUnknownFields = true
	defer func() {
		EnableDecoderDisallowUnknownFields = false
	}()

	obj := FooStructDisallowUnknownFields{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)

	obj = FooStructDisallowUnknownFields{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "what")
}
func testBodyBindingStringMap(t *testing.T, b binding.Binding, path, badPath, body, badBody string) {
	obj := make(map[string]string)
	req := requestWithBody("POST", path, body)
	if b.Name() == "form" {
		req.Header.Add("Content-Type", binding.MIMEPOSTForm)
	}
	err := b.Bind(req, &obj)
	require.NoError(t, err)
	require.NotNil(t, obj)
	require.Len(t, obj, 2)
	require.Equal(t, "bar", obj["foo"])
	require.Equal(t, "world", obj["hello"])

	if badPath != "" && badBody != "" {
		obj = make(map[string]string)
		req = requestWithBody("POST", badPath, badBody)
		err = b.Bind(req, &obj)
		assert.Error(t, err)
	}

	objInt := make(map[string]int)
	req = requestWithBody("POST", path, body)
	err = b.Bind(req, &objInt)
	require.Error(t, err)
}

var (
	JSON = jsonBinding{}
	//XML           = xmlBinding{}
	//Form          = formBinding{}
	//Query         = queryBinding{}
	//FormPost      = formPostBinding{}
	//FormMultipart = formMultipartBinding{}
	//ProtoBuf      = protobufBinding{}
	//MsgPack       = msgpackBinding{}
	//YAML          = yamlBinding{}
	//Uri           = uriBinding{}
	Header = headerBinding{}
	//TOML          = tomlBinding{}
)

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
