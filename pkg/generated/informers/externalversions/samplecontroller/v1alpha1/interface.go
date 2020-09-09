/*
this is a header file
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/schwarzeni/k8s-sample-controller-doc-demo/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Foos returns a FooInformer.
	Foos() FooInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Foos returns a FooInformer.
func (v *version) Foos() FooInformer {
	return &fooInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
