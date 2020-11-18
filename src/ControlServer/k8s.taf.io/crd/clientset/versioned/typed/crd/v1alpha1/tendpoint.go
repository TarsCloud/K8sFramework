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

package v1alpha1

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	scheme "k8s.taf.io/crd/clientset/versioned/scheme"
	v1alpha1 "k8s.taf.io/crd/v1alpha1"
)

// TEndpointsGetter has a method to return a TEndpointInterface.
// A group's client should implement this interface.
type TEndpointsGetter interface {
	TEndpoints(namespace string) TEndpointInterface
}

// TEndpointInterface has methods to work with TEndpoint resources.
type TEndpointInterface interface {
	Create(ctx context.Context, tEndpoint *v1alpha1.TEndpoint, opts v1.CreateOptions) (*v1alpha1.TEndpoint, error)
	Update(ctx context.Context, tEndpoint *v1alpha1.TEndpoint, opts v1.UpdateOptions) (*v1alpha1.TEndpoint, error)
	UpdateStatus(ctx context.Context, tEndpoint *v1alpha1.TEndpoint, opts v1.UpdateOptions) (*v1alpha1.TEndpoint, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.TEndpoint, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.TEndpointList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.TEndpoint, err error)
	TEndpointExpansion
}

// tEndpoints implements TEndpointInterface
type tEndpoints struct {
	client rest.Interface
	ns     string
}

// newTEndpoints returns a TEndpoints
func newTEndpoints(c *CrdV1alpha1Client, namespace string) *tEndpoints {
	return &tEndpoints{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the tEndpoint, and returns the corresponding tEndpoint object, and an error if there is any.
func (c *tEndpoints) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.TEndpoint, err error) {
	result = &v1alpha1.TEndpoint{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tendpoints").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of TEndpoints that match those selectors.
func (c *tEndpoints) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.TEndpointList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.TEndpointList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tendpoints").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested tEndpoints.
func (c *tEndpoints) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("tendpoints").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a tEndpoint and creates it.  Returns the server's representation of the tEndpoint, and an error, if there is any.
func (c *tEndpoints) Create(ctx context.Context, tEndpoint *v1alpha1.TEndpoint, opts v1.CreateOptions) (result *v1alpha1.TEndpoint, err error) {
	result = &v1alpha1.TEndpoint{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("tendpoints").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(tEndpoint).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a tEndpoint and updates it. Returns the server's representation of the tEndpoint, and an error, if there is any.
func (c *tEndpoints) Update(ctx context.Context, tEndpoint *v1alpha1.TEndpoint, opts v1.UpdateOptions) (result *v1alpha1.TEndpoint, err error) {
	result = &v1alpha1.TEndpoint{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("tendpoints").
		Name(tEndpoint.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(tEndpoint).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *tEndpoints) UpdateStatus(ctx context.Context, tEndpoint *v1alpha1.TEndpoint, opts v1.UpdateOptions) (result *v1alpha1.TEndpoint, err error) {
	result = &v1alpha1.TEndpoint{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("tendpoints").
		Name(tEndpoint.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(tEndpoint).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the tEndpoint and deletes it. Returns an error if one occurs.
func (c *tEndpoints) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tendpoints").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *tEndpoints) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tendpoints").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched tEndpoint.
func (c *tEndpoints) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.TEndpoint, err error) {
	result = &v1alpha1.TEndpoint{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("tendpoints").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
