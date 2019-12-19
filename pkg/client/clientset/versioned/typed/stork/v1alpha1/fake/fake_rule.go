/*
Copyright 2018 Openstorage.org

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
	v1alpha1 "github.com/libopenstorage/stork/pkg/apis/stork/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRules implements RuleInterface
type FakeRules struct {
	Fake *FakeStorkV1alpha1
	ns   string
}

var rulesResource = schema.GroupVersionResource{Group: "stork", Version: "v1alpha1", Resource: "rules"}

var rulesKind = schema.GroupVersionKind{Group: "stork", Version: "v1alpha1", Kind: "Rule"}

// Get takes name of the rule, and returns the corresponding rule object, and an error if there is any.
func (c *FakeRules) Get(name string, options v1.GetOptions) (result *v1alpha1.Rule, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(rulesResource, c.ns, name), &v1alpha1.Rule{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Rule), err
}

// List takes label and field selectors, and returns the list of Rules that match those selectors.
func (c *FakeRules) List(opts v1.ListOptions) (result *v1alpha1.RuleList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(rulesResource, rulesKind, c.ns, opts), &v1alpha1.RuleList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.RuleList{ListMeta: obj.(*v1alpha1.RuleList).ListMeta}
	for _, item := range obj.(*v1alpha1.RuleList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested rules.
func (c *FakeRules) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(rulesResource, c.ns, opts))

}

// Create takes the representation of a rule and creates it.  Returns the server's representation of the rule, and an error, if there is any.
func (c *FakeRules) Create(rule *v1alpha1.Rule) (result *v1alpha1.Rule, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(rulesResource, c.ns, rule), &v1alpha1.Rule{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Rule), err
}

// Update takes the representation of a rule and updates it. Returns the server's representation of the rule, and an error, if there is any.
func (c *FakeRules) Update(rule *v1alpha1.Rule) (result *v1alpha1.Rule, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(rulesResource, c.ns, rule), &v1alpha1.Rule{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Rule), err
}

// Delete takes name of the rule and deletes it. Returns an error if one occurs.
func (c *FakeRules) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(rulesResource, c.ns, name), &v1alpha1.Rule{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRules) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(rulesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.RuleList{})
	return err
}

// Patch applies the patch and returns the patched rule.
func (c *FakeRules) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Rule, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(rulesResource, c.ns, name, data, subresources...), &v1alpha1.Rule{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Rule), err
}
