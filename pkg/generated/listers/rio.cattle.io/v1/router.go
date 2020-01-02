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
	v1 "github.com/rancher/rio/pkg/apis/rio.cattle.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// RouterLister helps list Routers.
type RouterLister interface {
	// List lists all Routers in the indexer.
	List(selector labels.Selector) (ret []*v1.Router, err error)
	// Routers returns an object that can list and get Routers.
	Routers(namespace string) RouterNamespaceLister
	RouterListerExpansion
}

// routerLister implements the RouterLister interface.
type routerLister struct {
	indexer cache.Indexer
}

// NewRouterLister returns a new RouterLister.
func NewRouterLister(indexer cache.Indexer) RouterLister {
	return &routerLister{indexer: indexer}
}

// List lists all Routers in the indexer.
func (s *routerLister) List(selector labels.Selector) (ret []*v1.Router, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Router))
	})
	return ret, err
}

// Routers returns an object that can list and get Routers.
func (s *routerLister) Routers(namespace string) RouterNamespaceLister {
	return routerNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RouterNamespaceLister helps list and get Routers.
type RouterNamespaceLister interface {
	// List lists all Routers in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Router, err error)
	// Get retrieves the Router from the indexer for a given namespace and name.
	Get(name string) (*v1.Router, error)
	RouterNamespaceListerExpansion
}

// routerNamespaceLister implements the RouterNamespaceLister
// interface.
type routerNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Routers in the indexer for a given namespace.
func (s routerNamespaceLister) List(selector labels.Selector) (ret []*v1.Router, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Router))
	})
	return ret, err
}

// Get retrieves the Router from the indexer for a given namespace and name.
func (s routerNamespaceLister) Get(name string) (*v1.Router, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("router"), name)
	}
	return obj.(*v1.Router), nil
}
