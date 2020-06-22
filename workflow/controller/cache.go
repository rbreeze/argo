package controller

import (
	"crypto/sha1"
	"encoding/json"
	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Cache interface {
	Load(key []byte);
	Save(key []byte, value string);
}

type configMapCache struct {
	kubeClient kubernetes.Interface
	namespace string
}

func NewConfigMapCache(ns string, ki kubernetes.Interface) *configMapCache {

	return &configMapCache{namespace: ns, kubeClient: ki}
}

func generateKey(template *wfv1.Template) []byte {
	h := sha1.New()
	h.Write([]byte(template.Name))
	return h.Sum(nil)
}

func (c *configMapCache) Load(key string) (*wfv1.Outputs, bool) {
	// TODO: return value stored in ConfigMap cache under key, or nil if none exists
	log.Infof("REM Loading from cache %s...\n", key)
	cm, err := c.kubeClient.CoreV1().ConfigMaps(c.namespace).Get(key, metav1.GetOptions{})
	if err != nil || cm == nil {
		return nil, false
	}
	entry := make(map[string]interface{})
	err = json.Unmarshal([]byte(cm.Data[key]), &entry)
	if err != nil {
		panic(err)
	}
	for k, v := range entry {
        log.Infof("REM k:", k, "v:", v)
    }
    //result := wfv1.Outputs{}
	return nil, false
}

func (c *configMapCache) Save(key string, value *wfv1.Outputs) bool {
	// TODO: store value to ConfigMap cache
	wfname := "whalesay"
	log.Infof("Saving to cache %s...\n", key)
	outputsJSON, err := json.Marshal(value)
	if err != nil {
		return false
	}
	_, err = c.kubeClient.CoreV1().ConfigMaps(c.namespace).Create(
		&apiv1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: wfname + "-cache",
			},
			Data: map[string]string{
				wfname + "." + "key": string(outputsJSON),
			},
		},
	)
	return true
}
