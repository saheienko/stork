package controllers

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"time"

	"github.com/libopenstorage/stork/pkg/apis/stork"
	storkapi "github.com/libopenstorage/stork/pkg/apis/stork/v1alpha1"
	"github.com/libopenstorage/stork/pkg/controller"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/portworx/sched-ops/k8s"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type VolumeImportController struct {
}

func NewVolumeImportController() (*VolumeImportController, error) {
	return &VolumeImportController{}, nil
}

func (c *VolumeImportController) Init() error {
	if err := c.createCRD(); err != nil {
		return err
	}

	return controller.Register(
		&schema.GroupVersionKind{
			Group:   stork.GroupName,
			Version: storkapi.SchemeGroupVersion.Version,
			Kind:    reflect.TypeOf(storkapi.VolumeImport{}).Name(),
		},
		corev1.NamespaceAll,
		5*time.Minute,
		c,
	)
}

func (c *VolumeImportController) Handle(ctx context.Context, event sdk.Event) error {
	return ReconcileProtected(ctx, event)
}

func (c *VolumeImportController) createCRD() error {
	resource := k8s.CustomResource{
		Name:    storkapi.VolumeImportResourceName,
		Plural:  storkapi.VolumeImportResourcePlural,
		Group:   stork.GroupName,
		Version: storkapi.SchemeGroupVersion.Version,
		Scope:   apiextensionsv1beta1.NamespaceScoped,
		Kind:    reflect.TypeOf(storkapi.VolumeImport{}).Name(),
	}
	err := k8s.Instance().CreateCRD(resource)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return k8s.Instance().ValidateCRD(resource, 10*time.Second, 2*time.Minute)
}
