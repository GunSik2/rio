/*
Copyright 2020 Rancher Labs.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	"github.com/rancher/wrangler/pkg/generic"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	clientset "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/typed/gloo.solo.io/v1"
	informers "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/informers/externalversions/gloo.solo.io/v1"
	listers "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/listers/gloo.solo.io/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type SettingsHandler func(string, *v1.Settings) (*v1.Settings, error)

type SettingsController interface {
	generic.ControllerMeta
	SettingsClient

	OnChange(ctx context.Context, name string, sync SettingsHandler)
	OnRemove(ctx context.Context, name string, sync SettingsHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() SettingsCache
}

type SettingsClient interface {
	Create(*v1.Settings) (*v1.Settings, error)
	Update(*v1.Settings) (*v1.Settings, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.Settings, error)
	List(namespace string, opts metav1.ListOptions) (*v1.SettingsList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Settings, err error)
}

type SettingsCache interface {
	Get(namespace, name string) (*v1.Settings, error)
	List(namespace string, selector labels.Selector) ([]*v1.Settings, error)

	AddIndexer(indexName string, indexer SettingsIndexer)
	GetByIndex(indexName, key string) ([]*v1.Settings, error)
}

type SettingsIndexer func(obj *v1.Settings) ([]string, error)

type settingsController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.SettingsesGetter
	informer          informers.SettingsInformer
	gvk               schema.GroupVersionKind
}

func NewSettingsController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.SettingsesGetter, informer informers.SettingsInformer) SettingsController {
	return &settingsController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromSettingsHandlerToHandler(sync SettingsHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Settings
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Settings))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *settingsController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Settings))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateSettingsDeepCopyOnChange(client SettingsClient, obj *v1.Settings, handler func(obj *v1.Settings) (*v1.Settings, error)) (*v1.Settings, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *settingsController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *settingsController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *settingsController) OnChange(ctx context.Context, name string, sync SettingsHandler) {
	c.AddGenericHandler(ctx, name, FromSettingsHandlerToHandler(sync))
}

func (c *settingsController) OnRemove(ctx context.Context, name string, sync SettingsHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromSettingsHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *settingsController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), namespace, name)
}

func (c *settingsController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), namespace, name, duration)
}

func (c *settingsController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *settingsController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *settingsController) Cache() SettingsCache {
	return &settingsCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *settingsController) Create(obj *v1.Settings) (*v1.Settings, error) {
	return c.clientGetter.Settingses(obj.Namespace).Create(obj)
}

func (c *settingsController) Update(obj *v1.Settings) (*v1.Settings, error) {
	return c.clientGetter.Settingses(obj.Namespace).Update(obj)
}

func (c *settingsController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.Settingses(namespace).Delete(name, options)
}

func (c *settingsController) Get(namespace, name string, options metav1.GetOptions) (*v1.Settings, error) {
	return c.clientGetter.Settingses(namespace).Get(name, options)
}

func (c *settingsController) List(namespace string, opts metav1.ListOptions) (*v1.SettingsList, error) {
	return c.clientGetter.Settingses(namespace).List(opts)
}

func (c *settingsController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Settingses(namespace).Watch(opts)
}

func (c *settingsController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Settings, err error) {
	return c.clientGetter.Settingses(namespace).Patch(name, pt, data, subresources...)
}

type settingsCache struct {
	lister  listers.SettingsLister
	indexer cache.Indexer
}

func (c *settingsCache) Get(namespace, name string) (*v1.Settings, error) {
	return c.lister.Settingses(namespace).Get(name)
}

func (c *settingsCache) List(namespace string, selector labels.Selector) ([]*v1.Settings, error) {
	return c.lister.Settingses(namespace).List(selector)
}

func (c *settingsCache) AddIndexer(indexName string, indexer SettingsIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Settings))
		},
	}))
}

func (c *settingsCache) GetByIndex(indexName, key string) (result []*v1.Settings, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1.Settings))
	}
	return result, nil
}
