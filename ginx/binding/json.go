// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/json"
	"github.com/zedisdog/brynn/errx"
	"github.com/zedisdog/brynn/util/reflectx"
	"io"
	"net/http"
	"reflect"
)

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// any as a Number instead of as a float64.
var EnableDecoderUseNumber = false

// EnableDecoderDisallowUnknownFields is used to call the DisallowUnknownFields method
// on the JSON Decoder instance. DisallowUnknownFields causes the Decoder to
// return an error when the destination is a struct and the input contains object
// keys which do not match any non-ignored, exported fields in the destination.
var EnableDecoderDisallowUnknownFields = false

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (j jsonBinding) Bind(req *http.Request, obj any) (err error) {
	if req == nil || req.Body == nil {
		return errx.New("invalid request")
	}
	content, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}
	return j.BindBody(content, obj)
}

func (j jsonBinding) BindBody(body []byte, obj any) (err error) {
	var r any
	err = json.Unmarshal(body, &r)
	if err != nil {
		return
	}

	destValue := reflect.ValueOf(obj).Elem()
	err = reflectx.Unmarshal(r, destValue, "json")
	if err != nil {
		return
	}

	return validate(obj)
}
