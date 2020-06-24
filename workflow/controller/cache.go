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

var sampleVal string = "Hello world"

var sampleParam = wfv1.Parameter{
	Name: "hello",
	Value: &sampleVal,
}

var sampleEntry = CacheEntry{
	ExpiresAt: "2020-06-18T17:11:05Z",
	NodeID: "memoize-abx4124-123129321123",
	Outputs: wfv1.Outputs{},
}

type Cache interface {
	Load(key []byte);
	Save(key []byte, value string);
}

type CacheEntry struct {
	ExpiresAt string `json"expiresAt"`
	NodeID string `json"nodeID"`
	Outputs wfv1.Outputs `json"outputs"`
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

func (c *configMapCache) Clear() bool {
	err := c.kubeClient.CoreV1().ConfigMaps(c.namespace).Delete(c.configMapName, &metav1.DeleteOptions{})
	if err != nil {
		log.Infof("Error deleting ConfigMap cache %s: %s", c.configMapName, err)
		return false
	}
	return true
}

func (c *configMapCache) Load(key string) (*wfv1.Outputs, bool) {
	cm, err := c.kubeClient.CoreV1().ConfigMaps(c.namespace).Get(c.configMapName, metav1.GetOptions{})
	if err != nil {
		log.Infof("Error loading ConfigMap cache %s: %s", c.configMapName, err)
		return nil, false
	}
	if cm == nil {
		log.Infof("Cache miss: ConfigMap does not exist")
		return nil, false
	}
	log.Infof("ConfigMap cache %s loaded", c.configMapName)
	rawEntry, ok := cm.Data[key];
	if !ok || rawEntry == "" {
		log.Infof("Cache miss: Entry for %s doesn't exist", key)
		return nil, false
	}
	var entry CacheEntry
	err = json.Unmarshal([]byte(rawEntry), &entry)
	if err != nil {
		panic(err)
	}
	outputs := entry.Outputs
	log.Infof("ConfigMap cache %s hit for %s", c.configMapName, key)
	return &outputs, true
}

func (c *configMapCache) Save(key string, value *wfv1.Outputs) bool {
	log.Infof("Saving to cache %s...", key)
	cm, err := c.kubeClient.CoreV1().ConfigMaps(c.namespace).Get(c.configMapName, metav1.GetOptions{})
	if len(cm.Data) == 0 {
		_, err = c.kubeClient.CoreV1().ConfigMaps(c.namespace).Create(&apiv1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: c.configMapName,
				},
			},
		)
	}
	if err != nil {
		log.Infof("Error saving to cache: %s", err)
		return false
	}
	sampleEntry.Outputs.Parameters = append(sampleEntry.Outputs.Parameters, sampleParam)
	entryJSON, err := json.Marshal(sampleEntry)
	log.Infof("REM CACHE SAVE New Entry: %s", entryJSON)
	opts := apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: c.configMapName,
		},
		Data: map[string]string{
			key: string(entryJSON),
		},
	}

	_, err = c.kubeClient.CoreV1().ConfigMaps(c.namespace).Update(&opts)

	if err != nil {
		log.Infof("Error creating new cache entry for %s: %s", key, err)
		return false
	}
	return true
}
