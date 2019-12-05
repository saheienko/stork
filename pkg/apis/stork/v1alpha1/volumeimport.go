package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	VolumeImportResourceName   = "volumeimport"
	VolumeImportResourcePlural = "volumeimports"
)

type VolumeImportType string

const (
	VolumeImportTypeRsync VolumeImportType = "rsync"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeImport
type VolumeImport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              VolumeImportSpec   `json:"spec"`
	Status            VolumeImportStatus `json:"status"`
}

type VolumeImportSpec struct {
	Type            VolumeImportType        `json:"type,omitempty"`
	ClusterPairName string                  `json:"clusterPairName,omitempty"`
	Source          VolumeImportSource      `json:"source"`
	Destination     VolumeImportDestination `json:"destination"`
}

type VolumeImportSource struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type VolumeImportDestination struct {
	PVC VolumeImportPVC `json:"pvc,omitempty"`
}

type VolumeImportPVC struct {
	Metadata metav1.ObjectMeta                `json:"metadata,omitempty"`
	Spec     *corev1.PersistentVolumeClaimSpec `json:"spec,omitempty"`
}

type VolumeImportStatus struct {
	RsyncJobName string
	ConditionType string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeImportList is a list of VolumeImport resources.
type VolumeImportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []VolumeImport `json:"items"`
}
