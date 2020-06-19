package controller

import (
	"crypto/sha1"
)

type Cache interface {
	GenerateKey(template string, input string);
	Load(key string);
	Save(key string, value string);
}

type cache struct {
}

func (c* cache) GenerateKey(template string, input string) []byte {
	h := sha1.New()
	h.Write([]byte(template+input))
	return h.Sum(nil)
}

func Load(key []byte) {

}

func Save(key []byte, value string) {

}

