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

	err := UnmarshalMap(m, reflect.ValueOf(&mStruct).Elem(), "json")
	require.Nil(t, err)
	require.Equal(t, "1", mStruct.B)
	require.Equal(t, 1, mStruct.A)
	require.Equal(t, 1, *mStruct.C)
}

func TestMapStrust(t *testing.T) {
	m := map[string]any{
		"a": float64(1),
		"b": "1",
		"c": 1,
		"d": 1,
		"e": map[string]any{
			"a": 1,
		},
		"f": []any{1},
		"g": map[string]any{
			"d": 1,
		},
		"h": []any{
			map[string]any{
				"d": 1,
			},
		},
	}
	type testMapStruct2 struct {
		D int `json:"d"`
	}
	type testMapStruct struct {
		A *int   `json:"a"`
		B string `json:"b"`
		C *int   `json:"c"`
		testMapStruct2
		E struct {
			A int `json:"a"`
		} `json:"e"`
		F []int             `json:"f"`
		G *testMapStruct2   `json:"g"`
		H []*testMapStruct2 `json:"h"`
	}
	var mStruct testMapStruct

	err := UnmarshalMap(m, reflect.ValueOf(&mStruct).Elem(), "json")
	require.Nil(t, err)
	require.Equal(t, "1", mStruct.B)
	require.Equal(t, int(1), *mStruct.A)
	require.Equal(t, 1, *mStruct.C)
	require.Equal(t, 1, mStruct.D)
	require.Equal(t, 1, mStruct.E.A)
	require.Equal(t, []int{1}, mStruct.F)
	require.Equal(t, 1, mStruct.G.D)
	require.Equal(t, 1, mStruct.H[0].D)
}

func TestSlice(t *testing.T) {
	a := []any{1}
	var b []int

	res, err := UnmarshalSlice(a, reflect.ValueOf(b))
	require.Nil(t, err)
	require.Equal(t, []int{1}, res.Interface())
}

func TestSliceStruct(t *testing.T) {
	type s struct {
		A int `json:"a"`
	}
	a := []any{
		map[string]any{
			"a": 1,
		},
	}
	b := make([]s, 0, len(a))

	res, err := UnmarshalSlice(a, reflect.ValueOf(b), "json")
	require.Nil(t, err)
	require.Equal(t, []s{
		{A: 1},
	}, res.Interface())
}

func TestSliceStructPtr(t *testing.T) {
	type s struct {
		A int `json:"a"`
	}
	a := []any{
		map[string]any{
			"a": 1,
		},
	}
	b := make([]*s, 0, len(a))

	res, err := UnmarshalSlice(a, reflect.ValueOf(b), "json")
	require.Nil(t, err)
	require.Equal(t, []*s{
		{A: 1},
	}, res.Interface())
}
