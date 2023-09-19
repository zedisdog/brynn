package mapx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshal1(t *testing.T) {
	type test struct {
		A int `input:"a"`
	}

	m := map[string]any{
		"a": 1,
	}
	var a test

	Unmarshal(m, &a, false)
	require.Equal(t, 1, a.A)
}

func TestUnmarshal2(t *testing.T) {
	s := make([]int, 0, 2)
	s = append(s, 1)
	type test struct {
		A []int `input:"a"`
	}

	m := map[string]any{
		"a": s,
	}
	var a test

	Unmarshal(m, &a, false)
	require.Equal(t, s, a.A)
}

func TestUnmarshal3(t *testing.T) {
	type test2 struct {
		A int `input:"a"`
	}
	type test struct {
		test2
	}

	m := map[string]any{
		"a": 1,
	}
	var a test

	err := Unmarshal(m, &a, false)
	require.Nil(t, err)
	require.Equal(t, 1, a.A)
}

func TestUnmarshal4(t *testing.T) {
	type test struct {
		A int `input:"a"`
	}

	m := map[string]any{}
	var a test

	err := Unmarshal(m, &a, false)
	require.NotNil(t, err)
}

func TestUnmarshal5(t *testing.T) {
	type test struct {
		A int `input:"a,optional"`
	}

	m := map[string]any{}
	var a test

	err := Unmarshal(m, &a, false)
	require.Nil(t, err)
	require.Equal(t, 0, a.A)
}

func TestUnmarshal6(t *testing.T) {
	type test struct {
		A int `input:"a"`
	}

	m := map[string]any{
		"a": "1",
	}
	var a test

	err := Unmarshal(m, &a, true)
	require.Nil(t, err)
	require.Equal(t, 1, a.A)
}
