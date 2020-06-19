package controller

import (
	"crypto/sha1"
	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

type Cache interface {
	Load(key []byte);
	Save(key []byte, value string);
}

type configMapCache struct {
}

func generateKey(template *wfv1.Template) []byte {
	h := sha1.New()
	h.Write([]byte(template))
	return h.Sum(nil)
}

func (c *configMapCache) Load(key []byte) {

}

func (c *configMapCache) Save(key []byte, value string) {

}
