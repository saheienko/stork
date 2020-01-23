package k8s

import (
	"time"

	"github.com/portworx/sched-ops/k8s/apiextensions"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

// CRDOps is an interface to perfrom k8s Customer Resource operations
type CRDOps interface {
	// CreateCRD creates the given custom resource
	// This API will be deprecated soon. Use RegisterCRD instead
	CreateCRD(resource apiextensions.CustomResource) error
	// RegisterCRD creates the given custom resource
	RegisterCRD(crd *apiextensionsv1beta1.CustomResourceDefinition) error
	// ValidateCRD checks if the given CRD is registered
	ValidateCRD(resource apiextensions.CustomResource, timeout, retryInterval time.Duration) error
	// DeleteCRD deletes the CRD for the given complete name (plural.group)
	DeleteCRD(fullName string) error
}

// CRD APIs - BEGIN

func (k *k8sOps) CreateCRD(resource apiextensions.CustomResource) error {
	return apiextensions.Instance().CreateCRD(resource)
}

func (k *k8sOps) RegisterCRD(crd *apiextensionsv1beta1.CustomResourceDefinition) error {
	return apiextensions.Instance().RegisterCRD(crd)
}

func (k *k8sOps) ValidateCRD(resource apiextensions.CustomResource, timeout, retryInterval time.Duration) error {
	return apiextensions.Instance().ValidateCRD(resource, timeout, retryInterval)
}

func (k *k8sOps) DeleteCRD(fullName string) error {
	return apiextensions.Instance().DeleteCRD(fullName)
}

// CRD APIs - END
