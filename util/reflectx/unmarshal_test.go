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
		"c": 1,
	}
	type testMapStruct struct {
		A int    `json:"a"`
		B string `json:"b"`
		C *int   `json:"c"`
	}
	var mStruct testMapStruct

	err := unmarshalMap(m, reflect.ValueOf(&mStruct).Elem(), "json")
	require.Nil(t, err)
	require.Equal(t, "1", mStruct.B)
	require.Equal(t, 1, mStruct.A)
	require.Equal(t, 1, *mStruct.C)
}

func TestMapStrust(t *testing.T) {
	m := map[string]any{
		"a": 1,
		"b": "1",
		"c": 1,
		"d": 1,
		"e": map[string]any{
			"a": 1,
		},
	}
	type testMapStruct2 struct {
		D int `json:"d"`
	}
	type testMapStruct struct {
		A int    `json:"a"`
		B string `json:"b"`
		C *int   `json:"c"`
		testMapStruct2
		E struct {
			A int `json:"a"`
		} `json:"e"`
	}
	var mStruct testMapStruct

	err := unmarshalMap(m, reflect.ValueOf(&mStruct).Elem(), "json")
	require.Nil(t, err)
	require.Equal(t, "1", mStruct.B)
	require.Equal(t, 1, mStruct.A)
	require.Equal(t, 1, *mStruct.C)
	require.Equal(t, 1, mStruct.D)
	require.Equal(t, 1, mStruct.E.A)
}
