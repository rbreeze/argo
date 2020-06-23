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

func (c *configMapCache) Load(key string) (*wfv1.Outputs, bool) {
	log.Infof("got here")
	log.Infof("REM Got here. Key: %s", key)
	log.Infof("REM ConfigMapName: %s", c.configMapName)
	log.Infof("REM CACHE LOAD Loading key %s from cache %s...", key, c.configMapName)
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
	rawEntry, ok := cm.Data[key];
	if !ok || rawEntry == "" {
		log.Infof("Cache miss: Entry for %s doesn't exist", key)
		return nil, false
	}
	// What do i want to do here?
	// the cache entry is a plain JSON string
	// I want to Unmarshal the JSON string into a Go var
	// the json contains a map of strings to interfaces
	// one of them is an Output
	// so I want to take the interface{} and cast it to a wfv1.Outputs
	// so that I can access it like: outputs.parameters
	var entry CacheEntry
	err = json.Unmarshal([]byte(rawEntry), &entry)
	if err != nil {
		panic(err)
	}
	outputs := entry.Outputs
	log.Infof("REM CACHE LOAD Cache: %s", outputs)
	return &outputs, true
}

func (c *configMapCache) Save(key string, value *wfv1.Outputs) bool {
	log.Infof("Saving to cache %s...\n", key)
	//outputsJSON, err := json.Marshal(value)
	//if err != nil {
	//	return false
	//}
	cm, err := c.kubeClient.CoreV1().ConfigMaps(c.namespace).Get(c.configMapName, metav1.GetOptions{})
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
