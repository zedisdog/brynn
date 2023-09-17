package httpx

import (
	"encoding/json"
	"fmt"
	"mime"
	"testing"
)

func TestXxx(t *testing.T) {
	types, _ := mime.ExtensionsByType("application/x-www-form-urlencoded")
	fmt.Printf("%+v\n", types)
}

func TestXxx2(t *testing.T) {
	j := `{"a": 1, "b": 2, "c": [1,2,3]}`
	var v any
	json.Unmarshal([]byte(j), &v)
	fmt.Printf("%#v\n", v.(map[string]any)["c"])
}
