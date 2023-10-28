package binding

import "github.com/gin-gonic/gin/binding"

func validate(obj any) error {
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}
