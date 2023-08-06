package config

import (
	"bytes"
	"strings"

	"github.com/spf13/viper"
)

func NewWithBytesContent(configType string, content []byte) *viper.Viper {
	c := viper.New()
	c.SetConfigType(configType)
	err := c.ReadConfig(bytes.NewBuffer(content))
	if err != nil {
		panic(err)
	}
	c.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	c.AutomaticEnv()

	return c
}
