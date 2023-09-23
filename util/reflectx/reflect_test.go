package reflectx

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestXxx(t *testing.T) {
// 	type a struct {
// 		A int
// 	}

// 	var test *a

// 	test2 := a{
// 		A: 1,
// 	}

// 	err := CopyValue(reflect.ValueOf(test2), reflect.TypeOf(test2), reflect.ValueOf(test), reflect.TypeOf(test))
// 	require.Nil(t, err)
// }

func TestNewToValue(t *testing.T) {
	type a struct {
		A int
	}

	var test *a

	fmt.Printf("%+v\n", test)

	err := NewTypeToValue(reflect.TypeOf(test), reflect.ValueOf(test))
	require.Nil(t, err)

	fmt.Printf("%+v\n", test)
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
