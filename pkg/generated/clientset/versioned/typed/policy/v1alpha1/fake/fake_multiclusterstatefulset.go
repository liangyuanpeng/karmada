/*
Copyright The Karmada Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/karmada-io/karmada/pkg/apis/policy/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMultiClusterStatefulSets implements MultiClusterStatefulSetInterface
type FakeMultiClusterStatefulSets struct {
	Fake *FakePolicyV1alpha1
	ns   string
}

var multiclusterstatefulsetsResource = v1alpha1.SchemeGroupVersion.WithResource("multiclusterstatefulsets")

var multiclusterstatefulsetsKind = v1alpha1.SchemeGroupVersion.WithKind("MultiClusterStatefulSet")

// Get takes name of the multiClusterStatefulSet, and returns the corresponding multiClusterStatefulSet object, and an error if there is any.
func (c *FakeMultiClusterStatefulSets) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.MultiClusterStatefulSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(multiclusterstatefulsetsResource, c.ns, name), &v1alpha1.MultiClusterStatefulSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MultiClusterStatefulSet), err
}

// List takes label and field selectors, and returns the list of MultiClusterStatefulSets that match those selectors.
func (c *FakeMultiClusterStatefulSets) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.MultiClusterStatefulSetList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(multiclusterstatefulsetsResource, multiclusterstatefulsetsKind, c.ns, opts), &v1alpha1.MultiClusterStatefulSetList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.MultiClusterStatefulSetList{ListMeta: obj.(*v1alpha1.MultiClusterStatefulSetList).ListMeta}
	for _, item := range obj.(*v1alpha1.MultiClusterStatefulSetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested multiClusterStatefulSets.
func (c *FakeMultiClusterStatefulSets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(multiclusterstatefulsetsResource, c.ns, opts))

}

// Create takes the representation of a multiClusterStatefulSet and creates it.  Returns the server's representation of the multiClusterStatefulSet, and an error, if there is any.
func (c *FakeMultiClusterStatefulSets) Create(ctx context.Context, multiClusterStatefulSet *v1alpha1.MultiClusterStatefulSet, opts v1.CreateOptions) (result *v1alpha1.MultiClusterStatefulSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(multiclusterstatefulsetsResource, c.ns, multiClusterStatefulSet), &v1alpha1.MultiClusterStatefulSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MultiClusterStatefulSet), err
}

// Update takes the representation of a multiClusterStatefulSet and updates it. Returns the server's representation of the multiClusterStatefulSet, and an error, if there is any.
func (c *FakeMultiClusterStatefulSets) Update(ctx context.Context, multiClusterStatefulSet *v1alpha1.MultiClusterStatefulSet, opts v1.UpdateOptions) (result *v1alpha1.MultiClusterStatefulSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(multiclusterstatefulsetsResource, c.ns, multiClusterStatefulSet), &v1alpha1.MultiClusterStatefulSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MultiClusterStatefulSet), err
}

// Delete takes name of the multiClusterStatefulSet and deletes it. Returns an error if one occurs.
func (c *FakeMultiClusterStatefulSets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(multiclusterstatefulsetsResource, c.ns, name, opts), &v1alpha1.MultiClusterStatefulSet{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMultiClusterStatefulSets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(multiclusterstatefulsetsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.MultiClusterStatefulSetList{})
	return err
}

// Patch applies the patch and returns the patched multiClusterStatefulSet.
func (c *FakeMultiClusterStatefulSets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.MultiClusterStatefulSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(multiclusterstatefulsetsResource, c.ns, name, pt, data, subresources...), &v1alpha1.MultiClusterStatefulSet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MultiClusterStatefulSet), err
}