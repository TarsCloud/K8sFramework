package watch

import (
	"base"
	"encoding/json"
	"fmt"
	k8sApiV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"runtime"
	"strings"
	"sync"
)

type tServiceRecord struct {
	services map[string]base.ServerServant
	mutex    sync.RWMutex
}

func loadServerServantFromService(service *k8sApiV1.Service) base.ServerServant {
	servantAnnotations, ok := service.Annotations[base.TafServantLabel]
	if !ok {
		return nil
	}
	var servant base.ServerServant
	err := json.Unmarshal([]byte(servantAnnotations), &servant)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, fmt.Sprintf("Unexpected ServerServant Fromat")))
		return nil
	}
	return servant
}

func (c *tServiceRecord) onAddService(service *k8sApiV1.Service) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	servant := loadServerServantFromService(service)
	if servant != nil {
		c.services[service.Name] = servant
	}
}

func (c *tServiceRecord) onUpdateService(oldService, newService *k8sApiV1.Service) {
	c.mutex.Lock()
	c.mutex.Unlock()
	servant := loadServerServantFromService(newService)
	if servant != nil {
		c.services[newService.Name] = servant
	}
}

func (c *tServiceRecord) onDeleteService(service *k8sApiV1.Service) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.services, service.Name)
}

func (c *tServiceRecord) GetServerServant(app string, name string) base.ServerServant {
	serverName := strings.ToLower(app + "-" + name)
	return c.services[serverName]
}

var serviceRecord tServiceRecord

func prepareServiceWatch() {
	svcInformer := informerFactory.Core().V1().Services().Informer()
	svcInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			service := obj.(*k8sApiV1.Service)
			serviceRecord.onAddService(service)
		},

		DeleteFunc: func(obj interface{}) {
			service := obj.(*k8sApiV1.Service)
			serviceRecord.onDeleteService(service)
		},

		UpdateFunc: func(oldObj, newObj interface{}) {
			oldService := oldObj.(*k8sApiV1.Service)

			newService := newObj.(*k8sApiV1.Service)

			serviceRecord.onUpdateService(oldService, newService)
		},
	})
}

func init() {
	serviceRecord.services = make(map[string]base.ServerServant, 30)
	serviceRecord.mutex = sync.RWMutex{}
	registryK8SWatchPrepare(prepareServiceWatch)
}
