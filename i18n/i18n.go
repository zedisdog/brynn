package i18n

import (
	"github.com/zedisdog/brynn/util/conv"
	"strings"
)

type P map[string]any

type Language string

var DefaultLang Language = EnUs

const (
	EnUs Language = "en_US"
	ZhCn Language = "zh_cn"
)

var data = map[Language]map[string]string{}

func Set(data map[Language]map[string]string) {
	data = data
}

func Register(lang Language, key string, value string) {
	if data[lang] == nil {
		data[lang] = make(map[string]string)
	}

	data[lang][key] = value
}

func TransByLang(lang Language, key string) (value string) {
	l, ok := data[lang]
	if !ok {
		return key
	}

	value, ok = l[key]
	if !ok {
		return key
	}

	return
}

func Trans(key string) string {
	return TransByLang(DefaultLang, key)
}

func Transf(tempKey string, params map[string]any) (result string) {
	result = Trans(tempKey)

	for key, param := range params {
		s, err := conv.ConvertTo[string](param)
		if err != nil {
			panic(err)
		}
		result = strings.Replace(result, ":"+key, s, -1)
	}

	return
}
