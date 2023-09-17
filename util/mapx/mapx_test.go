package mapx

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSetStruct(t *testing.T) {
	type a struct {
		A int
	}

	var test *a

	test2 := a{
		A: 1,
	}

	v := reflect.ValueOf(test).Elem()
	Set(v, test2)

	fmt.Printf("%+v\n", test)
}
