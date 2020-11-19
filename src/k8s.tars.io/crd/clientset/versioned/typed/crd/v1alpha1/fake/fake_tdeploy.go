/*
Copyright The Kubernetes Authors.

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

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "k8s.tars.io/crd/v1alpha1"
)

// FakeTDeploys implements TDeployInterface
type FakeTDeploys struct {
	Fake *FakeCrdV1alpha1
	ns   string
}

var tdeploysResource = schema.GroupVersionResource{Group: "crd", Version: "v1alpha1", Resource: "tdeploys"}

var tdeploysKind = schema.GroupVersionKind{Group: "crd", Version: "v1alpha1", Kind: "TDeploy"}

// Get takes name of the tDeploy, and returns the corresponding tDeploy object, and an error if there is any.
func (c *FakeTDeploys) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.TDeploy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(tdeploysResource, c.ns, name), &v1alpha1.TDeploy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TDeploy), err
}

// List takes label and field selectors, and returns the list of TDeploys that match those selectors.
func (c *FakeTDeploys) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.TDeployList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(tdeploysResource, tdeploysKind, c.ns, opts), &v1alpha1.TDeployList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.TDeployList{ListMeta: obj.(*v1alpha1.TDeployList).ListMeta}
	for _, item := range obj.(*v1alpha1.TDeployList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested tDeploys.
func (c *FakeTDeploys) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(tdeploysResource, c.ns, opts))

}

// Create takes the representation of a tDeploy and creates it.  Returns the server's representation of the tDeploy, and an error, if there is any.
func (c *FakeTDeploys) Create(ctx context.Context, tDeploy *v1alpha1.TDeploy, opts v1.CreateOptions) (result *v1alpha1.TDeploy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(tdeploysResource, c.ns, tDeploy), &v1alpha1.TDeploy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TDeploy), err
}

// Update takes the representation of a tDeploy and updates it. Returns the server's representation of the tDeploy, and an error, if there is any.
func (c *FakeTDeploys) Update(ctx context.Context, tDeploy *v1alpha1.TDeploy, opts v1.UpdateOptions) (result *v1alpha1.TDeploy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(tdeploysResource, c.ns, tDeploy), &v1alpha1.TDeploy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TDeploy), err
}

// Delete takes name of the tDeploy and deletes it. Returns an error if one occurs.
func (c *FakeTDeploys) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(tdeploysResource, c.ns, name), &v1alpha1.TDeploy{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTDeploys) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(tdeploysResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.TDeployList{})
	return err
}

// Patch applies the patch and returns the patched tDeploy.
func (c *FakeTDeploys) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.TDeploy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(tdeploysResource, c.ns, name, pt, data, subresources...), &v1alpha1.TDeploy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TDeploy), err
}