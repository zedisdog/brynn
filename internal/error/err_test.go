package error

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	err := NewWithSkip("123", 0)
	err2 := WrapWithSkip(err, "321", 0)
	fmt.Printf("%#v", err2)
}
