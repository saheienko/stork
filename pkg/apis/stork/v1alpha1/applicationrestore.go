package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&ApplicationRestore{}, &ApplicationRestoreList{})
}

const (
	// ApplicationRestoreResourceName is name for "applicationrestore" resource
	ApplicationRestoreResourceName = "applicationrestore"
	// ApplicationRestoreResourcePlural is plural for "applicationrestore" resource
	ApplicationRestoreResourcePlural = "applicationrestores"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplicationRestore represents applicationrestore object
type ApplicationRestore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ApplicationRestoreSpec   `json:"spec"`
	Status            ApplicationRestoreStatus `json:"status"`
}

// ApplicationRestoreSpec is the spec used to restore applications
type ApplicationRestoreSpec struct {
	BackupName       string                              `json:"backupName"`
	BackupLocation   string                              `json:"backupLocation"`
	NamespaceMapping map[string]string                   `json:"namespaceMapping"`
	Selectors        map[string]string                   `json:"selectors"`
	EncryptionKey    *corev1.EnvVarSource                `json:"encryptionKey"`
	ReplacePolicy    ApplicationRestoreReplacePolicyType `json:"replacePolicy"`
}

// ApplicationRestoreReplacePolicyType is the replace policy for the application restore
// in case there are conflicting resources already present on the cluster
type ApplicationRestoreReplacePolicyType string

const (
	// ApplicationRestoreReplacePolicyDelete is to specify that the restore
	// should delete existing resources that conflict with resources being
	// restored
	ApplicationRestoreReplacePolicyDelete ApplicationRestoreReplacePolicyType = "Delete"
	// ApplicationRestoreReplacePolicyRetain is to specify that the restore
	// should retain existing resources that conflict with resources being
	// restored
	ApplicationRestoreReplacePolicyRetain ApplicationRestoreReplacePolicyType = "Retain"
)

// ApplicationRestoreStatus is the status of a application restore operation
type ApplicationRestoreStatus struct {
	Stage           ApplicationRestoreStageType       `json:"stage"`
	Status          ApplicationRestoreStatusType      `json:"status"`
	Resources       []*ApplicationRestoreResourceInfo `json:"resources"`
	Volumes         []*ApplicationRestoreVolumeInfo   `json:"volumes"`
	FinishTimestamp metav1.Time                       `json:"finishTimestamp"`
}

// ApplicationRestoreResourceInfo is the info for the restore of a resource
type ApplicationRestoreResourceInfo struct {
	Name                    string `json:"name"`
	Namespace               string `json:"namespace"`
	metav1.GroupVersionKind `json:",inline"`
	Status                  ApplicationRestoreStatusType `json:"status"`
	Reason                  string                       `json:"reason"`
}

// ApplicationRestoreVolumeInfo is the info for the restore of a volume
type ApplicationRestoreVolumeInfo struct {
	PersistentVolumeClaim string                       `json:"persistentVolumeClaim"`
	SourceNamespace       string                       `json:"sourceNamespace"`
	SourceVolume          string                       `json:"sourceVolume"`
	RestoreVolume         string                       `json:"restoreVolume"`
	Status                ApplicationRestoreStatusType `json:"status"`
	Reason                string                       `json:"reason"`
}

// ApplicationRestoreStatusType is the status of the application restore
type ApplicationRestoreStatusType string

const (
	// ApplicationRestoreStatusInitial is the initial state when restore is created
	ApplicationRestoreStatusInitial ApplicationRestoreStatusType = ""
	// ApplicationRestoreStatusPending for when restore is still pending
	ApplicationRestoreStatusPending ApplicationRestoreStatusType = "Pending"
	// ApplicationRestoreStatusInProgress for when restore is in progress
	ApplicationRestoreStatusInProgress ApplicationRestoreStatusType = "InProgress"
	// ApplicationRestoreStatusFailed for when restore has failed
	ApplicationRestoreStatusFailed ApplicationRestoreStatusType = "Failed"
	// ApplicationRestoreStatusPartialSuccess for when restore was partially successful
	ApplicationRestoreStatusPartialSuccess ApplicationRestoreStatusType = "PartialSuccess"
	// ApplicationRestoreStatusRetained for when restore was skipped to retain an already existing resource
	ApplicationRestoreStatusRetained ApplicationRestoreStatusType = "Retained"
	// ApplicationRestoreStatusSuccessful for when restore has completed successfully
	ApplicationRestoreStatusSuccessful ApplicationRestoreStatusType = "Successful"
)

// ApplicationRestoreStageType is the stage of the restore
type ApplicationRestoreStageType string

const (
	// ApplicationRestoreStageInitial for when restore is created
	ApplicationRestoreStageInitial ApplicationRestoreStageType = ""
	// ApplicationRestoreStageVolumes for when volumes are being restored
	ApplicationRestoreStageVolumes ApplicationRestoreStageType = "Volumes"
	// ApplicationRestoreStageApplications for when applications are being
	// restored
	ApplicationRestoreStageApplications ApplicationRestoreStageType = "Applications"
	// ApplicationRestoreStageFinal is the final stage for restore
	ApplicationRestoreStageFinal ApplicationRestoreStageType = "Final"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplicationRestoreList is a list of ApplicationRestores
type ApplicationRestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ApplicationRestore `json:"items"`
}
