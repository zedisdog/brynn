package cryptox

import (
	"crypto"
	"encoding/hex"
)

type HashOption struct {
	Key []byte
}

// WithKey set the key to sha256
func WithKey(str []byte) WithHashOption {
	return func(option *HashOption) {
		option.Key = str
	}
}

type WithHashOption func(option *HashOption)

func Hash(str string, hashType crypto.Hash, options ...WithHashOption) (string, error) {
	var option HashOption
	for _, o := range options {
		o(&option)
	}

	c := hashType.New()
	_, err := c.Write([]byte(str))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(c.Sum(option.Key)), nil
}

func CheckHash(hashStr, plainText string, hashType crypto.Hash, options ...WithHashOption) bool {
	hash, err := Hash(plainText, hashType, options...)
	if err != nil {
		return false
	}

	return hashStr == hash
}
