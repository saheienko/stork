/*
Copyright the Velero contributors.

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

package v1

import (
	v1 "github.com/heptio/velero/pkg/apis/velero/v1"
	scheme "github.com/heptio/velero/pkg/generated/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// DeleteBackupRequestsGetter has a method to return a DeleteBackupRequestInterface.
// A group's client should implement this interface.
type DeleteBackupRequestsGetter interface {
	DeleteBackupRequests(namespace string) DeleteBackupRequestInterface
}

// DeleteBackupRequestInterface has methods to work with DeleteBackupRequest resources.
type DeleteBackupRequestInterface interface {
	Create(*v1.DeleteBackupRequest) (*v1.DeleteBackupRequest, error)
	Update(*v1.DeleteBackupRequest) (*v1.DeleteBackupRequest, error)
	UpdateStatus(*v1.DeleteBackupRequest) (*v1.DeleteBackupRequest, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.DeleteBackupRequest, error)
	List(opts metav1.ListOptions) (*v1.DeleteBackupRequestList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.DeleteBackupRequest, err error)
	DeleteBackupRequestExpansion
}

// deleteBackupRequests implements DeleteBackupRequestInterface
type deleteBackupRequests struct {
	client rest.Interface
	ns     string
}

// newDeleteBackupRequests returns a DeleteBackupRequests
func newDeleteBackupRequests(c *VeleroV1Client, namespace string) *deleteBackupRequests {
	return &deleteBackupRequests{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the deleteBackupRequest, and returns the corresponding deleteBackupRequest object, and an error if there is any.
func (c *deleteBackupRequests) Get(name string, options metav1.GetOptions) (result *v1.DeleteBackupRequest, err error) {
	result = &v1.DeleteBackupRequest{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("deletebackuprequests").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of DeleteBackupRequests that match those selectors.
func (c *deleteBackupRequests) List(opts metav1.ListOptions) (result *v1.DeleteBackupRequestList, err error) {
	result = &v1.DeleteBackupRequestList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("deletebackuprequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested deleteBackupRequests.
func (c *deleteBackupRequests) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("deletebackuprequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a deleteBackupRequest and creates it.  Returns the server's representation of the deleteBackupRequest, and an error, if there is any.
func (c *deleteBackupRequests) Create(deleteBackupRequest *v1.DeleteBackupRequest) (result *v1.DeleteBackupRequest, err error) {
	result = &v1.DeleteBackupRequest{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("deletebackuprequests").
		Body(deleteBackupRequest).
		Do().
		Into(result)
	return
}

// Update takes the representation of a deleteBackupRequest and updates it. Returns the server's representation of the deleteBackupRequest, and an error, if there is any.
func (c *deleteBackupRequests) Update(deleteBackupRequest *v1.DeleteBackupRequest) (result *v1.DeleteBackupRequest, err error) {
	result = &v1.DeleteBackupRequest{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("deletebackuprequests").
		Name(deleteBackupRequest.Name).
		Body(deleteBackupRequest).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *deleteBackupRequests) UpdateStatus(deleteBackupRequest *v1.DeleteBackupRequest) (result *v1.DeleteBackupRequest, err error) {
	result = &v1.DeleteBackupRequest{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("deletebackuprequests").
		Name(deleteBackupRequest.Name).
		SubResource("status").
		Body(deleteBackupRequest).
		Do().
		Into(result)
	return
}

// Delete takes name of the deleteBackupRequest and deletes it. Returns an error if one occurs.
func (c *deleteBackupRequests) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("deletebackuprequests").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *deleteBackupRequests) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("deletebackuprequests").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched deleteBackupRequest.
func (c *deleteBackupRequests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.DeleteBackupRequest, err error) {
	result = &v1.DeleteBackupRequest{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("deletebackuprequests").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
