package statuschecker

import (
	"github.com/atlassian/smith"
	"github.com/atlassian/smith/pkg/resources"

	"github.com/pkg/errors"
	apiext_v1b1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ObjectStatusResult struct {
}

// ObjectStatusChecker checks object's status.
// Function is responsible for handling different versions of objects by itself.
type ObjectStatusChecker func(runtime.Object) (ObjectStatusResult, error)

// CRDStore gets a CRD definition for a Group and Kind of the resource (CRD instance).
// Returns nil if CRD definition was not found.
type CRDStore interface {
	Get(resource schema.GroupKind) (*apiext_v1b1.CustomResourceDefinition, error)
}

type Interface interface {
	CheckStatus(*unstructured.Unstructured) (isReady, retriableError bool, e error)
}

type Checker struct {
	Store      CRDStore
	KnownTypes map[schema.GroupKind]ObjectStatusChecker
}

func New(store CRDStore, kts ...map[schema.GroupKind]ObjectStatusChecker) (*Checker, error) {
	kt := make(map[schema.GroupKind]ObjectStatusChecker)
	for _, knownTypes := range kts {
		for knownGK, f := range knownTypes {
			if kt[knownGK] != nil {
				return nil, errors.Errorf("GroupKind specified more than once: %s", knownGK)
			}
			kt[knownGK] = f
		}
	}
	return &Checker{
		Store:      store,
		KnownTypes: kt,
	}, nil
}

func (c *Checker) CheckStatus(obj *unstructured.Unstructured) (isReady, retriableError bool, e error) {
	gvk := obj.GroupVersionKind()
	gk := gvk.GroupKind()

	if gk.Kind == "" || gvk.Version == "" { // Group can be empty e.g. built-in objects like ConfigMap
		return false, false, errors.Errorf("object has empty kind/version: %s", gvk)
	}

	// 1. Check if it is a known built-in resource
	if isObjectReady, ok := c.KnownTypes[gk]; ok {
		return isObjectReady(obj)
	}

	// 2. Check if it is a CRD with path/value annotation
	ready, retriable, err := c.checkPathValue(gk, obj)
	if err != nil || ready {
		return ready, retriable, err
	}

	// 3. Check if it is a CRD with Kind/GroupVersion annotation
	return c.checkForInstance(gk, obj)
}

func (c *Checker) checkForInstance(gk schema.GroupKind, obj *unstructured.Unstructured) (isReady, retriableError bool, e error) {
	// TODO Check if it is a CRD with Kind/GroupVersion annotation
	return false, false, nil
}

func (c *Checker) checkPathValue(gk schema.GroupKind, obj *unstructured.Unstructured) (isReady, retriableError bool, e error) {
	crd, err := c.Store.Get(gk)
	if err != nil {
		return false, true, err
	}
	if crd == nil {
		return false, false, nil
	}
	path := crd.Annotations[smith.CrFieldPathAnnotation]
	value := crd.Annotations[smith.CrFieldValueAnnotation]
	if len(path) == 0 || len(value) == 0 {
		return false, false, nil
	}
	actualValue, err := resources.GetJSONPathString(obj.Object, path)
	if err != nil {
		return false, false, err
	}
	if actualValue != value {
		return false, false, nil
	}
	return true, false, nil
}
