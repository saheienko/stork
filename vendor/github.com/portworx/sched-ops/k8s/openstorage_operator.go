package k8s

import (
	corev1alpha1 "github.com/libopenstorage/operator/pkg/apis/core/v1alpha1"
	"github.com/portworx/sched-ops/k8s/operator"
)

// StorageClusterOps is an interface to perfrom k8s StorageCluster operations
type StorageClusterOps interface {
	// CreateStorageCluster creates the given StorageCluster
	CreateStorageCluster(*corev1alpha1.StorageCluster) (*corev1alpha1.StorageCluster, error)
	// UpdateStorageCluster updates the given StorageCluster
	UpdateStorageCluster(*corev1alpha1.StorageCluster) (*corev1alpha1.StorageCluster, error)
	// GetStorageCluster gets the StorageCluster with given name and namespace
	GetStorageCluster(string, string) (*corev1alpha1.StorageCluster, error)
	// ListStorageClusters lists all the StorageClusters
	ListStorageClusters(string) (*corev1alpha1.StorageClusterList, error)
	// DeleteStorageCluster deletes the given StorageCluster
	DeleteStorageCluster(string, string) error
	// UpdateStorageClusterStatus update the status of given StorageCluster
	UpdateStorageClusterStatus(*corev1alpha1.StorageCluster) (*corev1alpha1.StorageCluster, error)
}

// StorageCluster APIs - BEGIN

func (k *k8sOps) CreateStorageCluster(cluster *corev1alpha1.StorageCluster) (*corev1alpha1.StorageCluster, error) {
	return operator.Instance().CreateStorageCluster(cluster)
}

func (k *k8sOps) UpdateStorageCluster(cluster *corev1alpha1.StorageCluster) (*corev1alpha1.StorageCluster, error) {
	return operator.Instance().UpdateStorageCluster(cluster)
}

func (k *k8sOps) GetStorageCluster(name, namespace string) (*corev1alpha1.StorageCluster, error) {
	return operator.Instance().GetStorageCluster(name, namespace)
}

func (k *k8sOps) ListStorageClusters(namespace string) (*corev1alpha1.StorageClusterList, error) {
	return operator.Instance().ListStorageClusters(namespace)
}

func (k *k8sOps) DeleteStorageCluster(name, namespace string) error {
	return operator.Instance().DeleteStorageCluster(name, namespace)
}

func (k *k8sOps) UpdateStorageClusterStatus(cluster *corev1alpha1.StorageCluster) (*corev1alpha1.StorageCluster, error) {
	return operator.Instance().UpdateStorageClusterStatus(cluster)
}

// StorageCluster APIs - END
