package k8s

import (
	"time"

	ocp_appsv1_api "github.com/openshift/api/apps/v1"
	ocp_securityv1_api "github.com/openshift/api/security/v1"
	"github.com/portworx/sched-ops/k8s/openshift"
	v1 "k8s.io/api/core/v1"
)

// SecurityContextConstraintsOps is an interface to list, get and update security context constraints
type SecurityContextConstraintsOps interface {
	// ListSecurityContextConstraints returns the list of all SecurityContextConstraints, and an error if there is any.
	ListSecurityContextConstraints() (*ocp_securityv1_api.SecurityContextConstraintsList, error)
	// GetSecurityContextConstraints takes name of the securityContextConstraints and returns the corresponding securityContextConstraints object, and an error if there is any.
	GetSecurityContextConstraints(string) (*ocp_securityv1_api.SecurityContextConstraints, error)
	// UpdateSecurityContextConstraints takes the representation of a securityContextConstraints and updates it. Returns the server's representation of the securityContextConstraints, and an error, if there is any.
	UpdateSecurityContextConstraints(*ocp_securityv1_api.SecurityContextConstraints) (*ocp_securityv1_api.SecurityContextConstraints, error)
}

// DeploymentConfigOps is an interface to perform ocp deployment config operations
type DeploymentConfigOps interface {
	// ListDeploymentConfigs lists all deployments for the given namespace
	ListDeploymentConfigs(namespace string) (*ocp_appsv1_api.DeploymentConfigList, error)
	// GetDeploymentConfig returns a deployment for the give name and namespace
	GetDeploymentConfig(name, namespace string) (*ocp_appsv1_api.DeploymentConfig, error)
	// CreateDeploymentConfig creates the given deployment
	CreateDeploymentConfig(*ocp_appsv1_api.DeploymentConfig) (*ocp_appsv1_api.DeploymentConfig, error)
	// UpdateDeploymentConfig updates the given deployment
	UpdateDeploymentConfig(*ocp_appsv1_api.DeploymentConfig) (*ocp_appsv1_api.DeploymentConfig, error)
	// DeleteDeploymentConfig deletes the given deployment
	DeleteDeploymentConfig(name, namespace string) error
	// ValidateDeploymentConfig validates the given deployment if it's running and healthy
	ValidateDeploymentConfig(deployment *ocp_appsv1_api.DeploymentConfig, timeout, retryInterval time.Duration) error
	// ValidateTerminatedDeploymentConfig validates if given deployment is terminated
	ValidateTerminatedDeploymentConfig(*ocp_appsv1_api.DeploymentConfig) error
	// GetDeploymentConfigPods returns pods for the given deployment
	GetDeploymentConfigPods(*ocp_appsv1_api.DeploymentConfig) ([]v1.Pod, error)
	// DescribeDeploymentConfig gets the deployment status
	DescribeDeploymentConfig(name, namespace string) (*ocp_appsv1_api.DeploymentConfigStatus, error)
	// GetDeploymentConfigsUsingStorageClass returns all deployments using the given storage class
	GetDeploymentConfigsUsingStorageClass(scName string) ([]ocp_appsv1_api.DeploymentConfig, error)
}

// Security Context Constraints APIs - BEGIN

func (k *k8sOps) ListSecurityContextConstraints() (result *ocp_securityv1_api.SecurityContextConstraintsList, err error) {
	return openshift.Instance().ListSecurityContextConstraints()
}

func (k *k8sOps) GetSecurityContextConstraints(name string) (result *ocp_securityv1_api.SecurityContextConstraints, err error) {
	return openshift.Instance().GetSecurityContextConstraints(name)
}

func (k *k8sOps) UpdateSecurityContextConstraints(securityContextConstraints *ocp_securityv1_api.SecurityContextConstraints) (result *ocp_securityv1_api.SecurityContextConstraints, err error) {
	return openshift.Instance().UpdateSecurityContextConstraints(securityContextConstraints)
}

// Security Context Constraints APIs - END

// DeploymentConfig APIs - BEGIN

func (k *k8sOps) ListDeploymentConfigs(namespace string) (*ocp_appsv1_api.DeploymentConfigList, error) {
	return openshift.Instance().ListDeploymentConfigs(namespace)
}

func (k *k8sOps) GetDeploymentConfig(name, namespace string) (*ocp_appsv1_api.DeploymentConfig, error) {
	return openshift.Instance().GetDeploymentConfig(name, namespace)
}

func (k *k8sOps) CreateDeploymentConfig(deployment *ocp_appsv1_api.DeploymentConfig) (*ocp_appsv1_api.DeploymentConfig, error) {
	return openshift.Instance().CreateDeploymentConfig(deployment)
}

func (k *k8sOps) DeleteDeploymentConfig(name, namespace string) error {
	return openshift.Instance().DeleteDeploymentConfig(name, namespace)
}

func (k *k8sOps) DescribeDeploymentConfig(depName, depNamespace string) (*ocp_appsv1_api.DeploymentConfigStatus, error) {
	return openshift.Instance().DescribeDeploymentConfig(depName, depNamespace)
}

func (k *k8sOps) UpdateDeploymentConfig(deployment *ocp_appsv1_api.DeploymentConfig) (*ocp_appsv1_api.DeploymentConfig, error) {
	return openshift.Instance().UpdateDeploymentConfig(deployment)
}

func (k *k8sOps) ValidateDeploymentConfig(deployment *ocp_appsv1_api.DeploymentConfig, timeout, retryInterval time.Duration) error {
	return openshift.Instance().ValidateDeploymentConfig(deployment, timeout, retryInterval)
}

func (k *k8sOps) ValidateTerminatedDeploymentConfig(deployment *ocp_appsv1_api.DeploymentConfig) error {
	return openshift.Instance().ValidateTerminatedDeploymentConfig(deployment)
}

func (k *k8sOps) GetDeploymentConfigPods(deployment *ocp_appsv1_api.DeploymentConfig) ([]v1.Pod, error) {
	return openshift.Instance().GetDeploymentConfigPods(deployment)
}

func (k *k8sOps) GetDeploymentConfigsUsingStorageClass(scName string) ([]ocp_appsv1_api.DeploymentConfig, error) {
	return openshift.Instance().GetDeploymentConfigsUsingStorageClass(scName)
}

// DeploymentConfig APIs - END
