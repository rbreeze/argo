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
	configMapName string
	kubeClient kubernetes.Interface
	namespace string
}

func NewConfigMapCache(cm string, ns string, ki kubernetes.Interface) *configMapCache {
	return &configMapCache{
		configMapName: cm,
		namespace: ns,
		kubeClient: ki,
	}
}

func generateKey(template *wfv1.Template) []byte {
	h := sha1.New()
	h.Write([]byte(template.Name))
	return h.Sum(nil)
}

func (c *configMapCache) Load(key string) (*wfv1.Outputs, bool) {
	// TODO: return value stored in ConfigMap cache under key, or nil if none exists
	log.Infof("REM Loading key %s from cache %s...", key, c.configMapName)
	cm, err := c.kubeClient.CoreV1().ConfigMaps(c.namespace).Get(c.configMapName, metav1.GetOptions{})
	if err != nil {
		log.Infof("Error loading ConfigMap %s: %s", c.configMapName, err)
		return nil, false
	}
	if cm == nil {
		log.Infof("Cache miss: ConfigMap does not exist")
		return nil, false
	}
	log.Infof("ConfigMap %s loaded", c.configMapName)
	entry, ok := cm.Data[key];
	if !ok {
		log.Infof("Cache miss: Entry for %s doesn't exist", key)
		return nil, false
	}
	rawEntry := make(map[string]interface{})
	err = json.Unmarshal([]byte(entry), &rawEntry)
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
	log.Infof("Saving to cache %s...\n", key)
	outputsJSON, err := json.Marshal(value)
	if err != nil {
		return false
	}
	cm, err := c.kubeClient.CoreV1().ConfigMaps(c.namespace).Get(c.configMapName, metav1.GetOptions{})
	if err != nil {
		log.Infof("Error saving to cache: %s", err)
		return false
	}
	opts := apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: c.configMapName,
		},
		Data: map[string]string{
			key: string(outputsJSON),
		},
	}
	if cm == nil {
		_, err = c.kubeClient.CoreV1().ConfigMaps(c.namespace).Create(&opts)
	} else {
		_, err = c.kubeClient.CoreV1().ConfigMaps(c.namespace).Update(&opts)
	}

	if err != nil {
		log.Infof("Error creating new cache entry for %s: %s", key, err)
		return false
	}
	return true
}
