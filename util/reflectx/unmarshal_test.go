package reflectx

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	m := map[string]any{
		"a": 1,
		"b": "1",
	}
	type testMapStruct struct {
		A int `json:"a"`
		B string `json:"b"`
	}
	var mStruct testMapStruct

	err := unmarshalMap(m, reflect.ValueOf(&mStruct).Elem(), "json")
	require.Nil(t, err)
	require.Nil(t, err)
	//fmt.Printf("%+v\n",mStruct)
}
