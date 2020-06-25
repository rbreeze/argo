package controller

import (
	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/test"
	"k8s.io/client-go/kubernetes"
	"testing"
	mock "github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

var MockParamValue string = "Hello world"

var MockParam = wfv1.Parameter{
	Name: "hello",
	Value: &MockParamValue,
}

func (_m *MockCache) NewConfigMapCache(cm string, ns string, ki kubernetes.Interface) *configMapCache {
	ret := _m.Called(cm, ns, ki)
	res := ret.Get(0)
	return res.(*configMapCache)
}

func (_m *MockCache) Load(key []byte) (*wfv1.Outputs, bool) {
	outputs := wfv1.Outputs{}
	outputs = append(outputs, MockParam)
	return &outputs, true
}

func (_m *MockCache) Save(key []byte, value string) bool {
	return true
}

func TestCacheLoad(t *testing.T) {
	mockCache := mocks.ConfigMapCache{}
}

func TestCacheSave(t *testing.T) {

}
