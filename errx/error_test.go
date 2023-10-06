package errx

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	_, file, line, _ := runtime.Caller(0)

	err := NewMsg("test").(*Error)

	require.Equal(t, file, err.Details["file"])
	require.Equal(t, line+2, err.Details["line"])
}

func TestWrap(t *testing.T) {
	_, file, line, _ := runtime.Caller(0)

	err := NewMsg("test1")
	err2 := WrapMsg(err, "test2").(*Error)

	require.Equal(t, file, err2.Details["file"])
	require.Equal(t, line+3, err2.Details["line"])
}

func TestFormat(t *testing.T) {
	_, file, line, _ := runtime.Caller(0)

	err1 := NewMsg("test1")
	err2 := WrapMsg(err1, "test2")
	err3 := WrapMsg(err2, "test3")

	require.Equal(t, "test3]<=[test2]<=[test1", fmt.Sprintf("%v", err3))
	require.Equal(t, fmt.Sprintf("%s:%d:test3\n", file, line+4), fmt.Sprintf("%+v", err3))

	except := []string{
		fmt.Sprintf("%s:%d:%s\n", file, line+4, err3.Error()),
		fmt.Sprintf("%s:%d:%s\n", file, line+3, err2.Error()),
		fmt.Sprintf("%s:%d:%s\n", file, line+2, err1.Error()),
	}
	require.Equal(t, strings.Join(except, ""), fmt.Sprintf("%#v", err3))
}
