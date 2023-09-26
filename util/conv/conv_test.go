package conv

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

// func TestConvertStringTo(t *testing.T) {
// 	s := "true"
// 	res, err := ConvertStringTo[bool](s)
// 	require.Nil(t, err)
// 	require.Equal(t, true, res)

// 	s2 := "123"
// 	res2, err := ConvertStringTo[int](s2)
// 	require.Nil(t, err)
// 	require.Equal(t, int(123), res2)
// }

func TestBoolTo(t *testing.T) {
	fmt.Printf("%+v\n", BoolTo[int8](true))
}

func TestConvert(t *testing.T) {
	a := 1

	b, err := ConvertTo[string](a)
	require.Nil(t, err)
	require.Equal(t, "1", b)

	type Int int
	var c Int = 1
	d, err := ConvertTo[string](c)
	require.Nil(t, err)
	require.Equal(t, "1", d)
}
