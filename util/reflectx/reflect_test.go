package reflectx

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXxx(t *testing.T) {
	s1 := []int{1, 2}
	s2 := make([]int, len(s1))
	copy(s2, s1)
	fmt.Printf("%+v\n", s2)
}

func TestGetTag(t *testing.T) {
	type a struct {
		A int `json:"testa" xml:"testb"`
	}

	test := a{
		A: 1,
	}

	content := GetTag(reflect.TypeOf(test).Field(0), "json")
	require.Equal(t, "testa", content)

	content = GetTag(reflect.TypeOf(test).Field(0), "xml", "json")
	require.Equal(t, "testb", content)
}

func TestConvertMapStrAny2MapStrType(t *testing.T) {
	m := map[string]any{
		"a": 1,
		"b": "1",
	}
	result, err := ConvertMapStrAny2MapStrType(reflect.ValueOf(m), reflect.TypeOf(""))
	require.Nil(t, err)
	require.Equal(t, "1", result.Interface().(map[string]string)["a"])
	require.Equal(t, "1", result.Interface().(map[string]string)["b"])
}

func TestSet(t *testing.T) {
	a := 1
	var b **int
	SetValue(a, reflect.ValueOf(&b))
	require.Equal(t, a, **b)

	c := []string{"1"}
	var d *[]string
	SetValue(c, reflect.ValueOf(&d))
	require.Equal(t, c, *d)

	e := map[string]string{"1": "2"}
	var f *map[string]string
	SetValue(e, reflect.ValueOf(&f))
	require.Equal(t, e, *f)
}
