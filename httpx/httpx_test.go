package httpx

import (
	"fmt"
	"mime"
	"testing"
)

func TestXxx(t *testing.T) {
	types, _ := mime.ExtensionsByType("application/x-www-form-urlencoded")
	fmt.Printf("%+v\n", types)
}
