package httpx

func isOptional(arr []string) bool {
	for _, item := range arr {
		if item == "optional" {
			return true
		}
	}

	return false
}
