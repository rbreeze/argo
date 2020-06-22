package controller

import (
	"crypto/sha1"
	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

type Cache interface {
	Load(key []byte);
	Save(key []byte, value string);
}

type configMapCache struct {
}

func NewConfigMapCache() *configMapCache {
	return &configMapCache{}
}

func generateKey(template *wfv1.Template) []byte {
	h := sha1.New()
	h.Write([]byte(template.Name))
	return h.Sum(nil)
}

func (c *configMapCache) Load(key string) (*wfv1.Outputs, bool) {
	// TODO: return value stored in ConfigMap cache under key, or nil if none exists
	return nil, false
}

func (c *configMapCache) Save(key string, value *wfv1.Outputs) bool {
	// TODO: store value to ConfigMap cache
	var cm *v1.ConfigMap
	cm.Name = key
	return true
}
