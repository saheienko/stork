package k8s

import (
	"sync"
	"time"

	snap_v1 "github.com/kubernetes-incubator/external-storage/snapshot/pkg/apis/crd/v1"
	autopilotclientset "github.com/libopenstorage/autopilot-api/pkg/client/clientset/versioned"
	storkclientset "github.com/libopenstorage/stork/pkg/client/clientset/versioned"
	ocp_clientset "github.com/openshift/client-go/apps/clientset/versioned"
	ocp_security_clientset "github.com/openshift/client-go/security/clientset/versioned"
	"github.com/portworx/sched-ops/k8s/admissionregistration"
	"github.com/portworx/sched-ops/k8s/apiextensions"
	"github.com/portworx/sched-ops/k8s/apps"
	"github.com/portworx/sched-ops/k8s/autopilot"
	"github.com/portworx/sched-ops/k8s/batch"
	"github.com/portworx/sched-ops/k8s/core"
	"github.com/portworx/sched-ops/k8s/discovery"
	"github.com/portworx/sched-ops/k8s/dynamic"
	"github.com/portworx/sched-ops/k8s/externalstorage"
	"github.com/portworx/sched-ops/k8s/openshift"
	"github.com/portworx/sched-ops/k8s/operator"
	"github.com/portworx/sched-ops/k8s/prometheus"
	"github.com/portworx/sched-ops/k8s/rbac"
	"github.com/portworx/sched-ops/k8s/storage"
	"github.com/portworx/sched-ops/k8s/stork"
	"github.com/portworx/sched-ops/k8s/talisman"
	"github.com/sirupsen/logrus"
	hook "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	batch_v1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/version"
	dynamicclient "k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	masterLabelKey           = "node-role.kubernetes.io/master"
	hostnameKey              = "kubernetes.io/hostname"
	pvcStorageClassKey       = "volume.beta.kubernetes.io/storage-class"
	pvcStorageProvisionerKey = "volume.beta.kubernetes.io/storage-provisioner"
	labelUpdateMaxRetries    = 5
)

var deleteForegroundPolicy = meta_v1.DeletePropagationForeground

// Ops is an interface to perform any kubernetes related operations
type Ops interface {
	NamespaceOps
	NodeOps
	ServiceOps
	StatefulSetOps
	DeploymentOps
	DeploymentConfigOps
	JobOps
	DaemonSetOps
	RBACOps
	PodOps
	StorageClassOps
	PersistentVolumeClaimOps
	VolumeAttachmentOps
	SnapshotOps
	SnapshotScheduleOps
	GroupSnapshotOps
	RuleOps
	SecretOps
	ConfigMapOps
	EventOps
	CRDOps
	ClusterPairOps
	MigrationOps
	ClusterDomainsOps
	AutopilotRuleOps
	StorageClusterOps
	ObjectOps
	SchedulePolicyOps
	VolumePlacementStrategyOps
	BackupLocationOps
	ApplicationBackupRestoreOps
	ApplicationCloneOps
	VolumeSnapshotRestoreOps
	SecurityContextConstraintsOps
	PrometheusOps
	MutatingWebhookConfigurationOps
	ClientSetter
	GetVersion() (*version.Info, error)
	// private methods for unit tests
	privateMethods
}

// EventOps is an interface to put and get k8s events
type EventOps interface {
	// CreateEvent puts an event into k8s etcd
	CreateEvent(event *v1.Event) (*v1.Event, error)
	// ListEvents retrieves all events registered with kubernetes
	ListEvents(namespace string, opts meta_v1.ListOptions) (*v1.EventList, error)
}

// NamespaceOps is an interface to perform namespace operations
type NamespaceOps interface {
	// ListNamespaces returns all the namespaces
	ListNamespaces(labelSelector map[string]string) (*v1.NamespaceList, error)
	// GetNamespace returns a namespace object for given name
	GetNamespace(name string) (*v1.Namespace, error)
	// CreateNamespace creates a namespace with given name and metadata
	CreateNamespace(name string, metadata map[string]string) (*v1.Namespace, error)
	// DeleteNamespace deletes a namespace with given name
	DeleteNamespace(name string) error
}

// NodeOps is an interface to perform k8s node operations
type NodeOps interface {
	// CreateNode creates the given node
	CreateNode(n *v1.Node) (*v1.Node, error)
	// UpdateNode updates the given node
	UpdateNode(n *v1.Node) (*v1.Node, error)
	// GetNodes talks to the k8s api server and gets the nodes in the cluster
	GetNodes() (*v1.NodeList, error)
	// GetNodeByName returns the k8s node given it's name
	GetNodeByName(string) (*v1.Node, error)
	// SearchNodeByAddresses searches corresponding k8s node match any of the given address
	SearchNodeByAddresses(addresses []string) (*v1.Node, error)
	// FindMyNode finds LOCAL Node in Kubernetes cluster
	FindMyNode() (*v1.Node, error)
	// IsNodeReady checks if node with given name is ready. Returns nil is ready.
	IsNodeReady(string) error
	// IsNodeMaster returns true if given node is a kubernetes master node
	IsNodeMaster(v1.Node) bool
	// GetLabelsOnNode gets all the labels on the given node
	GetLabelsOnNode(string) (map[string]string, error)
	// AddLabelOnNode adds a label key=value on the given node
	AddLabelOnNode(string, string, string) error
	// RemoveLabelOnNode removes the label with key on given node
	RemoveLabelOnNode(string, string) error
	// WatchNode sets up a watcher that listens for the changes on Node.
	WatchNode(node *v1.Node, fn core.WatchFunc) error
	// CordonNode cordons the given node
	CordonNode(nodeName string, timeout, retryInterval time.Duration) error
	// UnCordonNode uncordons the given node
	UnCordonNode(nodeName string, timeout, retryInterval time.Duration) error
	// DrainPodsFromNode drains given pods from given node. If timeout is set to
	// a non-zero value, it waits for timeout duration for each pod to get deleted
	DrainPodsFromNode(nodeName string, pods []v1.Pod, timeout, retryInterval time.Duration) error
}

// ServiceOps is an interface to perform k8s service operations
type ServiceOps interface {
	// GetService gets the service by the name
	GetService(string, string) (*v1.Service, error)
	// CreateService creates the given service
	CreateService(*v1.Service) (*v1.Service, error)
	// DeleteService deletes the given service
	DeleteService(name, namespace string) error
	// ValidateDeletedService validates if given service is deleted
	ValidateDeletedService(string, string) error
	// DescribeService gets the service status
	DescribeService(string, string) (*v1.ServiceStatus, error)
	// PatchService patches the current service with the given json path
	PatchService(name, namespace string, jsonPatch []byte) (*v1.Service, error)
}

// StatefulSetOps is an interface to perform k8s stateful set operations
type StatefulSetOps interface {
	// ListStatefulSets lists all the statefulsets for a given namespace
	ListStatefulSets(namespace string) (*appsv1.StatefulSetList, error)
	// GetStatefulSet returns a statefulset for given name and namespace
	GetStatefulSet(name, namespace string) (*appsv1.StatefulSet, error)
	// CreateStatefulSet creates the given statefulset
	CreateStatefulSet(ss *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	// UpdateStatefulSet creates the given statefulset
	UpdateStatefulSet(ss *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	// DeleteStatefulSet deletes the given statefulset
	DeleteStatefulSet(name, namespace string) error
	// ValidateStatefulSet validates the given statefulset if it's running and healthy within the given timeout
	ValidateStatefulSet(ss *appsv1.StatefulSet, timeout time.Duration) error
	// ValidateTerminatedStatefulSet validates if given deployment is terminated
	ValidateTerminatedStatefulSet(ss *appsv1.StatefulSet, timeout, retryInterval time.Duration) error
	// GetStatefulSetPods returns pods for the given statefulset
	GetStatefulSetPods(ss *appsv1.StatefulSet) ([]v1.Pod, error)
	// DescribeStatefulSet gets status of the statefulset
	DescribeStatefulSet(name, namespace string) (*appsv1.StatefulSetStatus, error)
	// GetStatefulSetsUsingStorageClass returns all statefulsets using given storage class
	GetStatefulSetsUsingStorageClass(scName string) ([]appsv1.StatefulSet, error)
	// GetPVCsForStatefulSet returns all the PVCs for given stateful set
	GetPVCsForStatefulSet(ss *appsv1.StatefulSet) (*v1.PersistentVolumeClaimList, error)
	// ValidatePVCsForStatefulSet validates the PVCs for the given stateful set
	ValidatePVCsForStatefulSet(ss *appsv1.StatefulSet, timeout, retryInterval time.Duration) error
}

// DeploymentOps is an interface to perform k8s deployment operations
type DeploymentOps interface {
	// ListDeployments lists all deployments for the given namespace
	ListDeployments(namespace string, options meta_v1.ListOptions) (*appsv1.DeploymentList, error)
	// GetDeployment returns a deployment for the give name and namespace
	GetDeployment(name, namespace string) (*appsv1.Deployment, error)
	// CreateDeployment creates the given deployment
	CreateDeployment(*appsv1.Deployment) (*appsv1.Deployment, error)
	// UpdateDeployment updates the given deployment
	UpdateDeployment(*appsv1.Deployment) (*appsv1.Deployment, error)
	// DeleteDeployment deletes the given deployment
	DeleteDeployment(name, namespace string) error
	// ValidateDeployment validates the given deployment if it's running and healthy
	ValidateDeployment(deployment *appsv1.Deployment, timeout, retryInterval time.Duration) error
	// ValidateTerminatedDeployment validates if given deployment is terminated
	ValidateTerminatedDeployment(*appsv1.Deployment, time.Duration, time.Duration) error
	// GetDeploymentPods returns pods for the given deployment
	GetDeploymentPods(*appsv1.Deployment) ([]v1.Pod, error)
	// DescribeDeployment gets the deployment status
	DescribeDeployment(name, namespace string) (*appsv1.DeploymentStatus, error)
	// GetDeploymentsUsingStorageClass returns all deployments using the given storage class
	GetDeploymentsUsingStorageClass(scName string) ([]appsv1.Deployment, error)
}

// DaemonSetOps is an interface to perform k8s daemon set operations
type DaemonSetOps interface {
	// CreateDaemonSet creates the given daemonset
	CreateDaemonSet(ds *appsv1.DaemonSet) (*appsv1.DaemonSet, error)
	// ListDaemonSets lists all daemonsets in given namespace
	ListDaemonSets(namespace string, listOpts meta_v1.ListOptions) ([]appsv1.DaemonSet, error)
	// GetDaemonSet gets the the daemon set with given name
	GetDaemonSet(string, string) (*appsv1.DaemonSet, error)
	// ValidateDaemonSet checks if the given daemonset is ready within given timeout
	ValidateDaemonSet(name, namespace string, timeout time.Duration) error
	// GetDaemonSetPods returns list of pods for the daemonset
	GetDaemonSetPods(*appsv1.DaemonSet) ([]v1.Pod, error)
	// UpdateDaemonSet updates the given daemon set and returns the updated ds
	UpdateDaemonSet(*appsv1.DaemonSet) (*appsv1.DaemonSet, error)
	// DeleteDaemonSet deletes the given daemonset
	DeleteDaemonSet(name, namespace string) error
}

// JobOps is an interface to perform job operations
type JobOps interface {
	// CreateJob creates the given job
	CreateJob(job *batch_v1.Job) (*batch_v1.Job, error)
	// GetJob returns the job from given namespace and name
	GetJob(name, namespace string) (*batch_v1.Job, error)
	// DeleteJob deletes the job with given namespace and name
	DeleteJob(name, namespace string) error
	// ValidateJob validates if the job with given namespace and name succeeds.
	//     It waits for timeout duration for job to succeed
	ValidateJob(name, namespace string, timeout time.Duration) error
}

// RBACOps is an interface to perform RBAC operations
type RBACOps interface {
	// CreateRole creates the given role
	CreateRole(role *rbac_v1.Role) (*rbac_v1.Role, error)
	// UpdateRole updates the given role
	UpdateRole(role *rbac_v1.Role) (*rbac_v1.Role, error)
	// GetRole gets the given role
	GetRole(name, namespace string) (*rbac_v1.Role, error)
	// CreateClusterRole creates the given cluster role
	CreateClusterRole(role *rbac_v1.ClusterRole) (*rbac_v1.ClusterRole, error)
	// GetClusterRole gets the given cluster role
	GetClusterRole(name string) (*rbac_v1.ClusterRole, error)
	// UpdateClusterRole updates the given cluster role
	UpdateClusterRole(role *rbac_v1.ClusterRole) (*rbac_v1.ClusterRole, error)
	// CreateRoleBinding creates the given role binding
	CreateRoleBinding(role *rbac_v1.RoleBinding) (*rbac_v1.RoleBinding, error)
	// UpdateRoleBinding updates the given role binding
	UpdateRoleBinding(role *rbac_v1.RoleBinding) (*rbac_v1.RoleBinding, error)
	// GetRoleBinding gets the given role binding
	GetRoleBinding(name, namespace string) (*rbac_v1.RoleBinding, error)
	// GetClusterRoleBinding gets the given cluster role binding
	GetClusterRoleBinding(name string) (*rbac_v1.ClusterRoleBinding, error)
	// ListClusterRoleBindings lists the cluster role bindings
	ListClusterRoleBindings() (*rbac_v1.ClusterRoleBindingList, error)
	// CreateClusterRoleBinding creates the given cluster role binding
	CreateClusterRoleBinding(role *rbac_v1.ClusterRoleBinding) (*rbac_v1.ClusterRoleBinding, error)
	// UpdateClusterRoleBinding updates the given cluster role binding
	UpdateClusterRoleBinding(role *rbac_v1.ClusterRoleBinding) (*rbac_v1.ClusterRoleBinding, error)
	// CreateServiceAccount creates the given service account
	CreateServiceAccount(account *v1.ServiceAccount) (*v1.ServiceAccount, error)
	// GetServiceAccount gets the given service account
	GetServiceAccount(name, namespace string) (*v1.ServiceAccount, error)
	// DeleteRole deletes the given role
	DeleteRole(name, namespace string) error
	// DeleteRoleBinding deletes the given role binding
	DeleteRoleBinding(name, namespace string) error
	// DeleteClusterRole deletes the given cluster role
	DeleteClusterRole(roleName string) error
	// DeleteClusterRoleBinding deletes the given cluster role binding
	DeleteClusterRoleBinding(roleName string) error
	// DeleteServiceAccount deletes the given service account
	DeleteServiceAccount(accountName, namespace string) error
}

// PodOps is an interface to perform k8s pod operations
type PodOps interface {
	// CreatePod creates the given pod.
	CreatePod(pod *v1.Pod) (*v1.Pod, error)
	// UpdatePod updates the given pod
	UpdatePod(pod *v1.Pod) (*v1.Pod, error)
	// GetPods returns pods for the given namespace
	GetPods(string, map[string]string) (*v1.PodList, error)
	// GetPodsByNode returns all pods in given namespace and given k8s node name.
	//  If namespace is empty, it will return pods from all namespaces
	GetPodsByNode(nodeName, namespace string) (*v1.PodList, error)
	// GetPodsByOwner returns pods for the given owner and namespace
	GetPodsByOwner(types.UID, string) ([]v1.Pod, error)
	// GetPodsUsingPV returns all pods in cluster using given pv
	GetPodsUsingPV(pvName string) ([]v1.Pod, error)
	// GetPodsUsingPVByNodeName returns all pods running on the node using the given pv
	GetPodsUsingPVByNodeName(pvName, nodeName string) ([]v1.Pod, error)
	// GetPodsUsingPVC returns all pods in cluster using given pvc
	GetPodsUsingPVC(pvcName, pvcNamespace string) ([]v1.Pod, error)
	// GetPodsUsingPVCByNodeName returns all pods running on the node using given pvc
	GetPodsUsingPVCByNodeName(pvcName, pvcNamespace, nodeName string) ([]v1.Pod, error)
	// GetPodsUsingVolumePlugin returns all pods who use PVCs provided by the given volume plugin
	GetPodsUsingVolumePlugin(plugin string) ([]v1.Pod, error)
	// GetPodsUsingVolumePluginByNodeName returns all pods who use PVCs provided by the given volume plugin on the given node
	GetPodsUsingVolumePluginByNodeName(nodeName, plugin string) ([]v1.Pod, error)
	// GetPodByName returns pod for the given pod name and namespace
	GetPodByName(string, string) (*v1.Pod, error)
	// GetPodByUID returns pod with the given UID, or error if nothing found
	GetPodByUID(types.UID, string) (*v1.Pod, error)
	// DeletePod deletes the given pod
	DeletePod(string, string, bool) error
	// DeletePods deletes the given pods
	DeletePods([]v1.Pod, bool) error
	// IsPodRunning checks if all containers in a pod are in running state
	IsPodRunning(v1.Pod) bool
	// IsPodReady checks if all containers in a pod are ready (passed readiness probe)
	IsPodReady(v1.Pod) bool
	// IsPodBeingManaged returns true if the pod is being managed by a controller
	IsPodBeingManaged(v1.Pod) bool
	// WaitForPodDeletion waits for given timeout for given pod to be deleted
	WaitForPodDeletion(uid types.UID, namespace string, timeout time.Duration) error
	// RunCommandInPod runs given command in the given pod
	RunCommandInPod(cmds []string, podName, containerName, namespace string) (string, error)
	// ValidatePod validates the given pod if it's ready
	ValidatePod(pod *v1.Pod, timeout, retryInterval time.Duration) error
	// WatchPods sets up a watcher that listens for the changes to pods in given namespace
	WatchPods(namespace string, fn core.WatchFunc, listOptions meta_v1.ListOptions) error
}

// StorageClassOps is an interface to perform k8s storage class operations
type StorageClassOps interface {
	// GetStorageClasses returns all storageClasses that match given optional label selector
	GetStorageClasses(labelSelector map[string]string) (*storagev1.StorageClassList, error)
	// GetStorageClass returns the storage class for the give namme
	GetStorageClass(name string) (*storagev1.StorageClass, error)
	// CreateStorageClass creates the given storage class
	CreateStorageClass(sc *storagev1.StorageClass) (*storagev1.StorageClass, error)
	// DeleteStorageClass deletes the given storage class
	DeleteStorageClass(name string) error
	// GetStorageClassParams returns the parameters of the given sc in the native map format
	GetStorageClassParams(sc *storagev1.StorageClass) (map[string]string, error)
	// ValidateStorageClass validates the given storage class
	// TODO: This is currently the same as GetStorageClass. If no one is using it,
	// we should remove this method
	ValidateStorageClass(name string) (*storagev1.StorageClass, error)
}

// PersistentVolumeClaimOps is an interface to perform k8s PVC operations
type PersistentVolumeClaimOps interface {
	// CreatePersistentVolumeClaim creates the given persistent volume claim
	CreatePersistentVolumeClaim(*v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error)
	// UpdatePersistentVolumeClaim updates an existing persistent volume claim
	UpdatePersistentVolumeClaim(*v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error)
	// DeletePersistentVolumeClaim deletes the given persistent volume claim
	DeletePersistentVolumeClaim(name, namespace string) error
	// ValidatePersistentVolumeClaim validates the given pvc
	ValidatePersistentVolumeClaim(vv *v1.PersistentVolumeClaim, timeout, retryInterval time.Duration) error
	// ValidatePersistentVolumeClaimSize validates the given pvc size
	ValidatePersistentVolumeClaimSize(vv *v1.PersistentVolumeClaim, expectedPVCSize int64, timeout, retryInterval time.Duration) error
	// GetPersistentVolumeClaim returns the PVC for given name and namespace
	GetPersistentVolumeClaim(pvcName string, namespace string) (*v1.PersistentVolumeClaim, error)
	// GetPersistentVolumeClaims returns all PVCs in given namespace and that match the optional labelSelector
	GetPersistentVolumeClaims(namespace string, labelSelector map[string]string) (*v1.PersistentVolumeClaimList, error)
	// CreatePersistentVolume creates the given PV
	CreatePersistentVolume(pv *v1.PersistentVolume) (*v1.PersistentVolume, error)
	// GetPersistentVolume returns the PV for given name
	GetPersistentVolume(pvName string) (*v1.PersistentVolume, error)
	// DeletePersistentVolume deletes the PV for given name
	DeletePersistentVolume(pvName string) error
	// GetPersistentVolumes returns all PVs in cluster
	GetPersistentVolumes() (*v1.PersistentVolumeList, error)
	// GetVolumeForPersistentVolumeClaim returns the volumeID for the given PVC
	GetVolumeForPersistentVolumeClaim(*v1.PersistentVolumeClaim) (string, error)
	// GetPersistentVolumeClaimParams fetches custom parameters for the given PVC
	GetPersistentVolumeClaimParams(*v1.PersistentVolumeClaim) (map[string]string, error)
	// GetPersistentVolumeClaimStatus returns the status of the given pvc
	GetPersistentVolumeClaimStatus(*v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaimStatus, error)
	// GetPVCsUsingStorageClass returns all PVCs that use the given storage class
	GetPVCsUsingStorageClass(scName string) ([]v1.PersistentVolumeClaim, error)
	// GetStorageProvisionerForPVC returns storage provisioner for given PVC if it exists
	GetStorageProvisionerForPVC(pvc *v1.PersistentVolumeClaim) (string, error)
}

// SnapshotOps is an interface to perform k8s VolumeSnapshot operations
type SnapshotOps interface {
	// GetSnapshot returns the snapshot for given name and namespace
	GetSnapshot(name string, namespace string) (*snap_v1.VolumeSnapshot, error)
	// ListSnapshots lists all snapshots in the given namespace
	ListSnapshots(namespace string) (*snap_v1.VolumeSnapshotList, error)
	// CreateSnapshot creates the given snapshot
	CreateSnapshot(*snap_v1.VolumeSnapshot) (*snap_v1.VolumeSnapshot, error)
	// UpdateSnapshot updates the given snapshot
	UpdateSnapshot(*snap_v1.VolumeSnapshot) (*snap_v1.VolumeSnapshot, error)
	// DeleteSnapshot deletes the given snapshot
	DeleteSnapshot(name string, namespace string) error
	// ValidateSnapshot validates the given snapshot.
	ValidateSnapshot(name string, namespace string, retry bool, timeout, retryInterval time.Duration) error
	// GetVolumeForSnapshot returns the volumeID for the given snapshot
	GetVolumeForSnapshot(name string, namespace string) (string, error)
	// GetSnapshotStatus returns the status of the given snapshot
	GetSnapshotStatus(name string, namespace string) (*snap_v1.VolumeSnapshotStatus, error)
	// GetSnapshotData returns the snapshot for given name
	GetSnapshotData(name string) (*snap_v1.VolumeSnapshotData, error)
	// CreateSnapshotData creates the given volume snapshot data object
	CreateSnapshotData(*snap_v1.VolumeSnapshotData) (*snap_v1.VolumeSnapshotData, error)
	// DeleteSnapshotData deletes the given snapshot
	DeleteSnapshotData(name string) error
	// ValidateSnapshotData validates the given snapshot data object
	ValidateSnapshotData(name string, retry bool, timeout, retryInterval time.Duration) error
}

// VolumeAttachmentOps is an interface to perform k8s VolumeAttachmentOps operations
type VolumeAttachmentOps interface {
	// ListVolumeAttachments lists all volume attachments
	ListVolumeAttachments() (*storagev1beta1.VolumeAttachmentList, error)
	// DeleteVolumeAttachment deletes a given Volume Attachment by name
	DeleteVolumeAttachment(name string) error
	// CreateVolumeAttachment creates a volume attachment
	CreateVolumeAttachment(*storagev1beta1.VolumeAttachment) (*storagev1beta1.VolumeAttachment, error)
	// UpdateVolumeAttachment updates a volume attachment
	UpdateVolumeAttachment(*storagev1beta1.VolumeAttachment) (*storagev1beta1.VolumeAttachment, error)
	// UpdateVolumeAttachmentStatus updates a volume attachment status
	UpdateVolumeAttachmentStatus(*storagev1beta1.VolumeAttachment) (*storagev1beta1.VolumeAttachment, error)
}

// SecretOps is an interface to perform k8s Secret operations
type SecretOps interface {
	// GetSecret gets the secrets object given its name and namespace
	GetSecret(name string, namespace string) (*v1.Secret, error)
	// CreateSecret creates the given secret
	CreateSecret(*v1.Secret) (*v1.Secret, error)
	// UpdateSecret updates the given secret
	UpdateSecret(*v1.Secret) (*v1.Secret, error)
	// UpdateSecretData updates or creates a new secret with the given data
	UpdateSecretData(string, string, map[string][]byte) (*v1.Secret, error)
	// DeleteSecret deletes the given secret
	DeleteSecret(name, namespace string) error
}

// ConfigMapOps is an interface to perform k8s ConfigMap operations
type ConfigMapOps interface {
	// GetConfigMap gets the config map object for the given name and namespace
	GetConfigMap(name string, namespace string) (*v1.ConfigMap, error)
	// CreateConfigMap creates a new config map object if it does not already exist.
	CreateConfigMap(configMap *v1.ConfigMap) (*v1.ConfigMap, error)
	// DeleteConfigMap deletes the given config map
	DeleteConfigMap(name, namespace string) error
	// UpdateConfigMap updates the given config map object
	UpdateConfigMap(configMap *v1.ConfigMap) (*v1.ConfigMap, error)
	// WatchConfigMap sets up a watcher that listens for changes on the config map
	WatchConfigMap(configMap *v1.ConfigMap, fn core.WatchFunc) error
}

// ObjectOps is an interface to perform generic Object operations
type ObjectOps interface {
	// GetObject returns the latest object given a generic Object
	GetObject(object runtime.Object) (runtime.Object, error)
	// UpdateObject updates a generic Object
	UpdateObject(object runtime.Object) (runtime.Object, error)
}

// MutatingWebhookConfigurationOps is interface to perform CRUD ops on mutatting webhook controller
type MutatingWebhookConfigurationOps interface {
	// GetMutatingWebhookConfiguration returns a given MutatingWebhookConfiguration
	GetMutatingWebhookConfiguration(name string) (*hook.MutatingWebhookConfiguration, error)
	// CreateMutatingWebhookConfiguration creates given MutatingWebhookConfiguration
	CreateMutatingWebhookConfiguration(req *hook.MutatingWebhookConfiguration) (*hook.MutatingWebhookConfiguration, error)
	// UpdateMutatingWebhookConfiguration updates given MutatingWebhookConfiguration
	UpdateMutatingWebhookConfiguration(*hook.MutatingWebhookConfiguration) (*hook.MutatingWebhookConfiguration, error)
	// DeleteMutatingWebhookConfiguration deletes given MutatingWebhookConfiguration
	DeleteMutatingWebhookConfiguration(name string) error
}

type privateMethods interface {
}

var (
	instance Ops
	once     sync.Once
)

type k8sOps struct{}

// Instance returns a singleton instance of k8sOps type
func Instance() Ops {
	once.Do(func() {
		instance = &k8sOps{}
	})
	return instance
}

// NewInstanceFromClients returns new instance of k8sOps by using given
// clients
func NewInstanceFromClients(
	kubernetesClient kubernetes.Interface,
	snapClient rest.Interface,
	storkClient storkclientset.Interface,
	apiExtensionClient apiextensionsclient.Interface,
	dynamicClient dynamicclient.Interface,
	ocpClient ocp_clientset.Interface,
	ocpSecurityClient ocp_security_clientset.Interface,
	autopilotClient autopilotclientset.Interface,
) Ops {
	admissionregistration.SetInstance(admissionregistration.New(kubernetesClient.AdmissionregistrationV1beta1()))
	apiextensions.SetInstance(apiextensions.New(apiExtensionClient))
	apps.SetInstance(apps.New(kubernetesClient.AppsV1(), kubernetesClient.CoreV1()))
	autopilot.SetInstance(autopilot.New(autopilotClient))
	batch.SetInstance(batch.New(kubernetesClient.BatchV1()))
	core.SetInstance(core.New(kubernetesClient.CoreV1(), kubernetesClient.StorageV1()))
	discovery.SetInstance(discovery.New(kubernetesClient.Discovery()))
	dynamic.SetInstance(dynamic.New(dynamicClient))
	externalstorage.SetInstance(externalstorage.New(snapClient))
	openshift.SetInstance(openshift.New(kubernetesClient, ocpClient, ocpSecurityClient))
	rbac.SetInstance(rbac.New(kubernetesClient.RbacV1()))
	storage.SetInstance(storage.New(kubernetesClient.StorageV1(), kubernetesClient.StorageV1beta1()))
	stork.SetInstance(stork.New(kubernetesClient, storkClient, snapClient))
	return &k8sOps{}
}

// NewInstanceFromConfigFile returns new instance of k8sOps by using given
// config file
func NewInstanceFromConfigFile(config string) (Ops, error) {
	newInstance := &k8sOps{}
	err := newInstance.loadClientFromKubeconfig(config)
	if err != nil {
		logrus.Errorf("Unable to set new instance: %v", err)
		return nil, err
	}
	return newInstance, nil
}

// NewInstanceFromConfigBytes returns new instance of k8sOps by using given
// config bytes
func NewInstanceFromConfigBytes(config []byte) (Ops, error) {
	newInstance := &k8sOps{}
	err := newInstance.loadClientFromConfigBytes(config)
	if err != nil {
		logrus.Errorf("Unable to set new instance: %v", err)
		return nil, err
	}
	return newInstance, nil
}

// NewInstanceFromRestConfig returns new instance of k8sOps by using given
// k8s rest client config
func NewInstanceFromRestConfig(config *rest.Config) (Ops, error) {
	// set config for k8s clients
	admissionregistration.Instance().SetConfig(config)
	apiextensions.Instance().SetConfig(config)
	apps.Instance().SetConfig(config)
	autopilot.Instance().SetConfig(config)
	batch.Instance().SetConfig(config)
	core.Instance().SetConfig(config)
	discovery.Instance().SetConfig(config)
	dynamic.Instance().SetConfig(config)
	externalstorage.Instance().SetConfig(config)
	openshift.Instance().SetConfig(config)
	operator.Instance().SetConfig(config)
	prometheus.Instance().SetConfig(config)
	rbac.Instance().SetConfig(config)
	storage.Instance().SetConfig(config)
	stork.Instance().SetConfig(config)
	talisman.Instance().SetConfig(config)

	return &k8sOps{}, nil
}

func (k *k8sOps) GetVersion() (*version.Info, error) {
	return discovery.Instance().GetVersion()
}

// Namespace APIs - BEGIN

func (k *k8sOps) ListNamespaces(labelSelector map[string]string) (*v1.NamespaceList, error) {
	return core.Instance().ListNamespaces(labelSelector)
}

func (k *k8sOps) GetNamespace(name string) (*v1.Namespace, error) {
	return core.Instance().GetNamespace(name)
}

func (k *k8sOps) CreateNamespace(name string, metadata map[string]string) (*v1.Namespace, error) {
	return core.Instance().CreateNamespace(name, metadata)
}

func (k *k8sOps) DeleteNamespace(name string) error {
	return core.Instance().DeleteNamespace(name)
}

// Namespace APIs - END
func (k *k8sOps) CreateNode(n *v1.Node) (*v1.Node, error) {
	return core.Instance().CreateNode(n)
}

func (k *k8sOps) UpdateNode(n *v1.Node) (*v1.Node, error) {
	return core.Instance().UpdateNode(n)
}

func (k *k8sOps) GetNodes() (*v1.NodeList, error) {
	return core.Instance().GetNodes()
}

func (k *k8sOps) GetNodeByName(name string) (*v1.Node, error) {
	return core.Instance().GetNodeByName(name)
}

func (k *k8sOps) IsNodeReady(name string) error {
	return core.Instance().IsNodeReady(name)
}

func (k *k8sOps) IsNodeMaster(node v1.Node) bool {
	return core.Instance().IsNodeMaster(node)
}

func (k *k8sOps) GetLabelsOnNode(name string) (map[string]string, error) {
	return core.Instance().GetLabelsOnNode(name)
}

// SearchNodeByAddresses searches the node based on the IP addresses, then it falls back to a
// search by hostname, and finally by the labels
func (k *k8sOps) SearchNodeByAddresses(addresses []string) (*v1.Node, error) {
	return core.Instance().SearchNodeByAddresses(addresses)
}

// FindMyNode finds LOCAL Node in Kubernetes cluster.
func (k *k8sOps) FindMyNode() (*v1.Node, error) {
	return core.Instance().FindMyNode()
}

func (k *k8sOps) AddLabelOnNode(name, key, value string) error {
	return core.Instance().AddLabelOnNode(name, key, value)
}

func (k *k8sOps) RemoveLabelOnNode(name, key string) error {
	return core.Instance().RemoveLabelOnNode(name, key)
}

func (k *k8sOps) WatchNode(node *v1.Node, watchNodeFn core.WatchFunc) error {
	return core.Instance().WatchNode(node, watchNodeFn)
}

func (k *k8sOps) CordonNode(nodeName string, timeout, retryInterval time.Duration) error {
	return core.Instance().CordonNode(nodeName, timeout, retryInterval)
}

func (k *k8sOps) UnCordonNode(nodeName string, timeout, retryInterval time.Duration) error {
	return core.Instance().UnCordonNode(nodeName, timeout, retryInterval)
}

func (k *k8sOps) DrainPodsFromNode(nodeName string, pods []v1.Pod, timeout time.Duration, retryInterval time.Duration) error {
	return core.Instance().DrainPodsFromNode(nodeName, pods, timeout, retryInterval)
}

func (k *k8sOps) WaitForPodDeletion(uid types.UID, namespace string, timeout time.Duration) error {
	return core.Instance().WaitForPodDeletion(uid, namespace, timeout)
}

func (k *k8sOps) RunCommandInPod(cmds []string, podName, containerName, namespace string) (string, error) {
	return core.Instance().RunCommandInPod(cmds, podName, containerName, namespace)
}

// Service APIs - BEGIN

func (k *k8sOps) CreateService(service *v1.Service) (*v1.Service, error) {
	return core.Instance().CreateService(service)
}

func (k *k8sOps) DeleteService(name, namespace string) error {
	return core.Instance().DeleteService(name, namespace)
}

func (k *k8sOps) GetService(svcName string, svcNS string) (*v1.Service, error) {
	return core.Instance().GetService(svcName, svcNS)
}

func (k *k8sOps) DescribeService(svcName string, svcNamespace string) (*v1.ServiceStatus, error) {
	return core.Instance().DescribeService(svcName, svcNamespace)
}

func (k *k8sOps) ValidateDeletedService(svcName string, svcNS string) error {
	return core.Instance().ValidateDeletedService(svcName, svcNS)
}

func (k *k8sOps) PatchService(name, namespace string, jsonPatch []byte) (*v1.Service, error) {
	return core.Instance().PatchService(name, namespace, jsonPatch)
}

// Service APIs - END

// Deployment APIs - BEGIN

func (k *k8sOps) ListDeployments(namespace string, options meta_v1.ListOptions) (*appsv1.DeploymentList, error) {
	return apps.Instance().ListDeployments(namespace, options)
}

func (k *k8sOps) GetDeployment(name, namespace string) (*appsv1.Deployment, error) {
	return apps.Instance().GetDeployment(name, namespace)
}

func (k *k8sOps) CreateDeployment(deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return apps.Instance().CreateDeployment(deployment)
}

func (k *k8sOps) DeleteDeployment(name, namespace string) error {
	return apps.Instance().DeleteDeployment(name, namespace)
}

func (k *k8sOps) DescribeDeployment(depName, depNamespace string) (*appsv1.DeploymentStatus, error) {
	return apps.Instance().DescribeDeployment(depName, depNamespace)
}

func (k *k8sOps) UpdateDeployment(deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return apps.Instance().UpdateDeployment(deployment)
}

func (k *k8sOps) ValidateDeployment(deployment *appsv1.Deployment, timeout, retryInterval time.Duration) error {
	return apps.Instance().ValidateDeployment(deployment, timeout, retryInterval)
}

func (k *k8sOps) ValidateTerminatedDeployment(deployment *appsv1.Deployment, timeout, timeBeforeRetry time.Duration) error {
	return apps.Instance().ValidateTerminatedDeployment(deployment, timeout, timeBeforeRetry)
}

func (k *k8sOps) GetDeploymentPods(deployment *appsv1.Deployment) ([]v1.Pod, error) {
	return apps.Instance().GetDeploymentPods(deployment)
}

func (k *k8sOps) GetDeploymentsUsingStorageClass(scName string) ([]appsv1.Deployment, error) {
	return apps.Instance().GetDeploymentsUsingStorageClass(scName)
}

// Deployment APIs - END

// DaemonSet APIs - BEGIN

func (k *k8sOps) CreateDaemonSet(ds *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return apps.Instance().CreateDaemonSet(ds)
}

func (k *k8sOps) ListDaemonSets(namespace string, listOpts meta_v1.ListOptions) ([]appsv1.DaemonSet, error) {
	return apps.Instance().ListDaemonSets(namespace, listOpts)
}

func (k *k8sOps) GetDaemonSet(name, namespace string) (*appsv1.DaemonSet, error) {
	return apps.Instance().GetDaemonSet(name, namespace)
}

func (k *k8sOps) GetDaemonSetPods(ds *appsv1.DaemonSet) ([]v1.Pod, error) {
	return apps.Instance().GetDaemonSetPods(ds)
}

func (k *k8sOps) ValidateDaemonSet(name, namespace string, timeout time.Duration) error {
	return apps.Instance().ValidateDaemonSet(name, namespace, timeout)
}

func (k *k8sOps) UpdateDaemonSet(ds *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return apps.Instance().UpdateDaemonSet(ds)
}

func (k *k8sOps) DeleteDaemonSet(name, namespace string) error {
	return apps.Instance().DeleteDaemonSet(name, namespace)
}

// DaemonSet APIs - END

// Job APIs - BEGIN
func (k *k8sOps) CreateJob(job *batch_v1.Job) (*batch_v1.Job, error) {
	return batch.Instance().CreateJob(job)
}

func (k *k8sOps) GetJob(name, namespace string) (*batch_v1.Job, error) {
	return batch.Instance().GetJob(name, namespace)
}

func (k *k8sOps) DeleteJob(name, namespace string) error {
	return batch.Instance().DeleteJob(name, namespace)
}

func (k *k8sOps) ValidateJob(name, namespace string, timeout time.Duration) error {
	return batch.Instance().ValidateJob(name, namespace, timeout)
}

// Job APIs - END

// StatefulSet APIs - BEGIN

func (k *k8sOps) ListStatefulSets(namespace string) (*appsv1.StatefulSetList, error) {
	return apps.Instance().ListStatefulSets(namespace)
}

func (k *k8sOps) GetStatefulSet(name, namespace string) (*appsv1.StatefulSet, error) {
	return apps.Instance().GetStatefulSet(name, namespace)
}

func (k *k8sOps) CreateStatefulSet(statefulset *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return apps.Instance().CreateStatefulSet(statefulset)
}

func (k *k8sOps) DeleteStatefulSet(name, namespace string) error {
	return apps.Instance().DeleteStatefulSet(name, namespace)
}

func (k *k8sOps) DescribeStatefulSet(ssetName string, ssetNamespace string) (*appsv1.StatefulSetStatus, error) {
	return apps.Instance().DescribeStatefulSet(ssetName, ssetNamespace)
}

func (k *k8sOps) UpdateStatefulSet(statefulset *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return apps.Instance().UpdateStatefulSet(statefulset)
}

func (k *k8sOps) ValidateStatefulSet(statefulset *appsv1.StatefulSet, timeout time.Duration) error {
	return apps.Instance().ValidateStatefulSet(statefulset, timeout)
}

func (k *k8sOps) GetStatefulSetPods(statefulset *appsv1.StatefulSet) ([]v1.Pod, error) {
	return apps.Instance().GetStatefulSetPods(statefulset)
}

func (k *k8sOps) ValidateTerminatedStatefulSet(statefulset *appsv1.StatefulSet, timeout, timeBeforeRetry time.Duration) error {
	return apps.Instance().ValidateTerminatedStatefulSet(statefulset, timeout, timeBeforeRetry)
}

func (k *k8sOps) GetStatefulSetsUsingStorageClass(scName string) ([]appsv1.StatefulSet, error) {
	return apps.Instance().GetStatefulSetsUsingStorageClass(scName)
}

func (k *k8sOps) GetPVCsForStatefulSet(ss *appsv1.StatefulSet) (*v1.PersistentVolumeClaimList, error) {
	return apps.Instance().GetPVCsForStatefulSet(ss)
}

func (k *k8sOps) ValidatePVCsForStatefulSet(ss *appsv1.StatefulSet, timeout, retryTimeout time.Duration) error {
	return apps.Instance().ValidatePVCsForStatefulSet(ss, timeout, retryTimeout)
}

// StatefulSet APIs - END

// RBAC APIs - BEGIN

func (k *k8sOps) CreateRole(role *rbac_v1.Role) (*rbac_v1.Role, error) {
	return rbac.Instance().CreateRole(role)
}

func (k *k8sOps) UpdateRole(role *rbac_v1.Role) (*rbac_v1.Role, error) {
	return rbac.Instance().UpdateRole(role)
}

func (k *k8sOps) GetRole(name, namespace string) (*rbac_v1.Role, error) {
	return rbac.Instance().GetRole(name, namespace)
}

func (k *k8sOps) CreateClusterRole(role *rbac_v1.ClusterRole) (*rbac_v1.ClusterRole, error) {
	return rbac.Instance().CreateClusterRole(role)
}

func (k *k8sOps) GetClusterRole(name string) (*rbac_v1.ClusterRole, error) {
	return rbac.Instance().GetClusterRole(name)
}

func (k *k8sOps) UpdateClusterRole(role *rbac_v1.ClusterRole) (*rbac_v1.ClusterRole, error) {
	return rbac.Instance().UpdateClusterRole(role)
}

func (k *k8sOps) CreateRoleBinding(binding *rbac_v1.RoleBinding) (*rbac_v1.RoleBinding, error) {
	return rbac.Instance().CreateRoleBinding(binding)
}

func (k *k8sOps) UpdateRoleBinding(binding *rbac_v1.RoleBinding) (*rbac_v1.RoleBinding, error) {
	return rbac.Instance().UpdateRoleBinding(binding)
}

func (k *k8sOps) GetRoleBinding(name, namespace string) (*rbac_v1.RoleBinding, error) {
	return rbac.Instance().GetRoleBinding(name, namespace)
}

func (k *k8sOps) CreateClusterRoleBinding(binding *rbac_v1.ClusterRoleBinding) (*rbac_v1.ClusterRoleBinding, error) {
	return rbac.Instance().CreateClusterRoleBinding(binding)
}

func (k *k8sOps) UpdateClusterRoleBinding(binding *rbac_v1.ClusterRoleBinding) (*rbac_v1.ClusterRoleBinding, error) {
	return rbac.Instance().UpdateClusterRoleBinding(binding)
}

func (k *k8sOps) GetClusterRoleBinding(name string) (*rbac_v1.ClusterRoleBinding, error) {
	return rbac.Instance().GetClusterRoleBinding(name)
}

func (k *k8sOps) ListClusterRoleBindings() (*rbac_v1.ClusterRoleBindingList, error) {
	return rbac.Instance().ListClusterRoleBindings()
}

func (k *k8sOps) CreateServiceAccount(account *v1.ServiceAccount) (*v1.ServiceAccount, error) {
	return core.Instance().CreateServiceAccount(account)
}

func (k *k8sOps) GetServiceAccount(name, namespace string) (*v1.ServiceAccount, error) {
	return core.Instance().GetServiceAccount(name, namespace)
}

func (k *k8sOps) DeleteRole(name, namespace string) error {
	return rbac.Instance().DeleteRole(name, namespace)
}

func (k *k8sOps) DeleteClusterRole(roleName string) error {
	return rbac.Instance().DeleteClusterRole(roleName)
}

func (k *k8sOps) DeleteRoleBinding(name, namespace string) error {
	return rbac.Instance().DeleteRoleBinding(name, namespace)
}

func (k *k8sOps) DeleteClusterRoleBinding(bindingName string) error {
	return rbac.Instance().DeleteClusterRoleBinding(bindingName)
}

func (k *k8sOps) DeleteServiceAccount(accountName, namespace string) error {
	return core.Instance().DeleteServiceAccount(accountName, namespace)
}

// RBAC APIs - END

// Pod APIs - BEGIN

func (k *k8sOps) DeletePods(pods []v1.Pod, force bool) error {
	return core.Instance().DeletePods(pods, force)
}

func (k *k8sOps) DeletePod(name string, ns string, force bool) error {
	return core.Instance().DeletePod(name, ns, force)
}

func (k *k8sOps) CreatePod(pod *v1.Pod) (*v1.Pod, error) {
	return core.Instance().CreatePod(pod)
}

func (k *k8sOps) UpdatePod(pod *v1.Pod) (*v1.Pod, error) {
	return core.Instance().UpdatePod(pod)
}

func (k *k8sOps) GetPods(namespace string, labelSelector map[string]string) (*v1.PodList, error) {
	return core.Instance().GetPods(namespace, labelSelector)
}

func (k *k8sOps) GetPodsByNode(nodeName, namespace string) (*v1.PodList, error) {
	return core.Instance().GetPodsByNode(nodeName, namespace)
}

func (k *k8sOps) GetPodsByOwner(ownerUID types.UID, namespace string) ([]v1.Pod, error) {
	return core.Instance().GetPodsByOwner(ownerUID, namespace)
}

func (k *k8sOps) GetPodsUsingPV(pvName string) ([]v1.Pod, error) {
	return core.Instance().GetPodsUsingPV(pvName)
}

func (k *k8sOps) GetPodsUsingPVByNodeName(pvName, nodeName string) ([]v1.Pod, error) {
	return core.Instance().GetPodsUsingPVByNodeName(pvName, nodeName)
}

func (k *k8sOps) GetPodsUsingPVC(pvcName, pvcNamespace string) ([]v1.Pod, error) {
	return core.Instance().GetPodsUsingPVC(pvcName, pvcNamespace)
}

func (k *k8sOps) GetPodsUsingPVCByNodeName(pvcName, pvcNamespace, nodeName string) ([]v1.Pod, error) {
	return core.Instance().GetPodsUsingPVCByNodeName(pvcName, pvcNamespace, nodeName)
}

func (k *k8sOps) GetPodsUsingVolumePlugin(plugin string) ([]v1.Pod, error) {
	return core.Instance().GetPodsUsingVolumePlugin(plugin)
}

func (k *k8sOps) GetPodsUsingVolumePluginByNodeName(nodeName, plugin string) ([]v1.Pod, error) {
	return core.Instance().GetPodsUsingVolumePluginByNodeName(nodeName, plugin)
}

func (k *k8sOps) GetPodByName(podName string, namespace string) (*v1.Pod, error) {
	return core.Instance().GetPodByName(podName, namespace)
}

func (k *k8sOps) GetPodByUID(uid types.UID, namespace string) (*v1.Pod, error) {
	return core.Instance().GetPodByUID(uid, namespace)
}

func (k *k8sOps) IsPodRunning(pod v1.Pod) bool {
	return core.Instance().IsPodRunning(pod)
}

func (k *k8sOps) IsPodReady(pod v1.Pod) bool {
	return core.Instance().IsPodReady(pod)
}

func (k *k8sOps) IsPodBeingManaged(pod v1.Pod) bool {
	return core.Instance().IsPodBeingManaged(pod)
}

func (k *k8sOps) ValidatePod(pod *v1.Pod, timeout, retryInterval time.Duration) error {
	return core.Instance().ValidatePod(pod, timeout, retryInterval)
}

func (k *k8sOps) WatchPods(namespace string, fn core.WatchFunc, listOptions meta_v1.ListOptions) error {
	return core.Instance().WatchPods(namespace, fn, listOptions)
}

// Pod APIs - END

// StorageClass APIs - BEGIN

func (k *k8sOps) GetStorageClasses(labelSelector map[string]string) (*storagev1.StorageClassList, error) {
	return storage.Instance().GetStorageClasses(labelSelector)
}

func (k *k8sOps) GetStorageClass(name string) (*storagev1.StorageClass, error) {
	return storage.Instance().GetStorageClass(name)
}

func (k *k8sOps) CreateStorageClass(sc *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	return storage.Instance().CreateStorageClass(sc)
}

func (k *k8sOps) DeleteStorageClass(name string) error {
	return storage.Instance().DeleteStorageClass(name)
}

func (k *k8sOps) GetStorageClassParams(sc *storagev1.StorageClass) (map[string]string, error) {
	return storage.Instance().GetStorageClassParams(sc)
}

func (k *k8sOps) ValidateStorageClass(name string) (*storagev1.StorageClass, error) {
	return storage.Instance().ValidateStorageClass(name)
}

// StorageClass APIs - END

// PVC APIs - BEGIN

func (k *k8sOps) CreatePersistentVolumeClaim(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	return core.Instance().CreatePersistentVolumeClaim(pvc)
}

func (k *k8sOps) UpdatePersistentVolumeClaim(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	return core.Instance().UpdatePersistentVolumeClaim(pvc)
}

func (k *k8sOps) DeletePersistentVolumeClaim(name, namespace string) error {
	return core.Instance().DeletePersistentVolumeClaim(name, namespace)
}

func (k *k8sOps) ValidatePersistentVolumeClaim(pvc *v1.PersistentVolumeClaim, timeout, retryInterval time.Duration) error {
	return core.Instance().ValidatePersistentVolumeClaim(pvc, timeout, retryInterval)
}

func (k *k8sOps) ValidatePersistentVolumeClaimSize(pvc *v1.PersistentVolumeClaim, expectedPVCSize int64, timeout, retryInterval time.Duration) error {
	return core.Instance().ValidatePersistentVolumeClaimSize(pvc, expectedPVCSize, timeout, retryInterval)
}

func (k *k8sOps) CreatePersistentVolume(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {
	return core.Instance().CreatePersistentVolume(pv)
}

func (k *k8sOps) GetPersistentVolumeClaim(pvcName string, namespace string) (*v1.PersistentVolumeClaim, error) {
	return core.Instance().GetPersistentVolumeClaim(pvcName, namespace)
}

func (k *k8sOps) GetPersistentVolumeClaims(namespace string, labelSelector map[string]string) (*v1.PersistentVolumeClaimList, error) {
	return core.Instance().GetPersistentVolumeClaims(namespace, labelSelector)
}

func (k *k8sOps) GetPersistentVolume(pvName string) (*v1.PersistentVolume, error) {
	return core.Instance().GetPersistentVolume(pvName)
}

func (k *k8sOps) DeletePersistentVolume(pvName string) error {
	return core.Instance().DeletePersistentVolume(pvName)
}

func (k *k8sOps) GetPersistentVolumes() (*v1.PersistentVolumeList, error) {
	return core.Instance().GetPersistentVolumes()
}

func (k *k8sOps) GetVolumeForPersistentVolumeClaim(pvc *v1.PersistentVolumeClaim) (string, error) {
	return core.Instance().GetVolumeForPersistentVolumeClaim(pvc)
}

func (k *k8sOps) GetPersistentVolumeClaimStatus(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaimStatus, error) {
	return core.Instance().GetPersistentVolumeClaimStatus(pvc)
}

func (k *k8sOps) GetPersistentVolumeClaimParams(pvc *v1.PersistentVolumeClaim) (map[string]string, error) {
	return core.Instance().GetPersistentVolumeClaimParams(pvc)
}

func (k *k8sOps) GetPVCsUsingStorageClass(scName string) ([]v1.PersistentVolumeClaim, error) {
	return core.Instance().GetPVCsUsingStorageClass(scName)
}

func (k *k8sOps) GetStorageProvisionerForPVC(pvc *v1.PersistentVolumeClaim) (string, error) {
	return core.Instance().GetStorageProvisionerForPVC(pvc)
}

// PVCs APIs - END

// Snapshot APIs - BEGIN

func (k *k8sOps) CreateSnapshot(snap *snap_v1.VolumeSnapshot) (*snap_v1.VolumeSnapshot, error) {
	return externalstorage.Instance().CreateSnapshot(snap)
}

func (k *k8sOps) UpdateSnapshot(snap *snap_v1.VolumeSnapshot) (*snap_v1.VolumeSnapshot, error) {
	return externalstorage.Instance().UpdateSnapshot(snap)
}

func (k *k8sOps) DeleteSnapshot(name string, namespace string) error {
	return externalstorage.Instance().DeleteSnapshot(name, namespace)
}

func (k *k8sOps) ValidateSnapshot(name string, namespace string, retry bool, timeout, retryInterval time.Duration) error {
	return externalstorage.Instance().ValidateSnapshot(name, namespace, retry, timeout, retryInterval)
}

func (k *k8sOps) ValidateSnapshotData(name string, retry bool, timeout, retryInterval time.Duration) error {
	return externalstorage.Instance().ValidateSnapshotData(name, retry, timeout, retryInterval)
}

func (k *k8sOps) GetVolumeForSnapshot(name string, namespace string) (string, error) {
	return externalstorage.Instance().GetVolumeForSnapshot(name, namespace)
}

func (k *k8sOps) GetSnapshot(name string, namespace string) (*snap_v1.VolumeSnapshot, error) {
	return externalstorage.Instance().GetSnapshot(name, namespace)
}

func (k *k8sOps) ListSnapshots(namespace string) (*snap_v1.VolumeSnapshotList, error) {
	return externalstorage.Instance().ListSnapshots(namespace)
}

func (k *k8sOps) GetSnapshotStatus(name string, namespace string) (*snap_v1.VolumeSnapshotStatus, error) {
	return externalstorage.Instance().GetSnapshotStatus(name, namespace)
}

func (k *k8sOps) GetSnapshotData(name string) (*snap_v1.VolumeSnapshotData, error) {
	return externalstorage.Instance().GetSnapshotData(name)
}

func (k *k8sOps) CreateSnapshotData(snapData *snap_v1.VolumeSnapshotData) (*snap_v1.VolumeSnapshotData, error) {
	return externalstorage.Instance().CreateSnapshotData(snapData)
}

func (k *k8sOps) DeleteSnapshotData(name string) error {
	return externalstorage.Instance().DeleteSnapshotData(name)
}

// Snapshot APIs - END

// Secret APIs - BEGIN

func (k *k8sOps) GetSecret(name string, namespace string) (*v1.Secret, error) {
	return core.Instance().GetSecret(name, namespace)
}

func (k *k8sOps) CreateSecret(secret *v1.Secret) (*v1.Secret, error) {
	return core.Instance().CreateSecret(secret)
}

func (k *k8sOps) UpdateSecret(secret *v1.Secret) (*v1.Secret, error) {
	return core.Instance().UpdateSecret(secret)
}

func (k *k8sOps) UpdateSecretData(name string, ns string, data map[string][]byte) (*v1.Secret, error) {
	return core.Instance().UpdateSecretData(name, ns, data)
}

func (k *k8sOps) DeleteSecret(name, namespace string) error {
	return core.Instance().DeleteSecret(name, namespace)
}

// Secret APIs - END

// ConfigMap APIs - BEGIN

func (k *k8sOps) GetConfigMap(name string, namespace string) (*v1.ConfigMap, error) {
	return core.Instance().GetConfigMap(name, namespace)
}

func (k *k8sOps) CreateConfigMap(configMap *v1.ConfigMap) (*v1.ConfigMap, error) {
	return core.Instance().CreateConfigMap(configMap)
}

func (k *k8sOps) DeleteConfigMap(name, namespace string) error {
	return core.Instance().DeleteConfigMap(name, namespace)
}

func (k *k8sOps) UpdateConfigMap(configMap *v1.ConfigMap) (*v1.ConfigMap, error) {
	return core.Instance().UpdateConfigMap(configMap)
}

func (k *k8sOps) WatchConfigMap(configMap *v1.ConfigMap, fn core.WatchFunc) error {
	return core.Instance().WatchConfigMap(configMap, fn)
}

// ConfigMap APIs - END

// Event APIs - BEGIN
// CreateEvent puts an event into k8s etcd
func (k *k8sOps) CreateEvent(event *v1.Event) (*v1.Event, error) {
	return core.Instance().CreateEvent(event)
}

// ListEvents retrieves all events registered with kubernetes
func (k *k8sOps) ListEvents(namespace string, opts meta_v1.ListOptions) (*v1.EventList, error) {
	return core.Instance().ListEvents(namespace, opts)
}

// Event APIs - END

// Object APIs - BEGIN

// GetObject returns the latest object given a generic Object
func (k *k8sOps) GetObject(object runtime.Object) (runtime.Object, error) {
	return dynamic.Instance().GetObject(object)
}

// UpdateObject updates a generic Object
func (k *k8sOps) UpdateObject(object runtime.Object) (runtime.Object, error) {
	return dynamic.Instance().UpdateObject(object)
}

// Object APIs - END

// VolumeAttachment APIs - START

func (k *k8sOps) ListVolumeAttachments() (*storagev1beta1.VolumeAttachmentList, error) {
	return storage.Instance().ListVolumeAttachments()
}

func (k *k8sOps) DeleteVolumeAttachment(name string) error {
	return storage.Instance().DeleteVolumeAttachment(name)
}

func (k *k8sOps) CreateVolumeAttachment(volumeAttachment *storagev1beta1.VolumeAttachment) (*storagev1beta1.VolumeAttachment, error) {
	return storage.Instance().CreateVolumeAttachment(volumeAttachment)
}

func (k *k8sOps) UpdateVolumeAttachment(volumeAttachment *storagev1beta1.VolumeAttachment) (*storagev1beta1.VolumeAttachment, error) {
	return storage.Instance().UpdateVolumeAttachment(volumeAttachment)
}

func (k *k8sOps) UpdateVolumeAttachmentStatus(volumeAttachment *storagev1beta1.VolumeAttachment) (*storagev1beta1.VolumeAttachment, error) {
	return storage.Instance().UpdateVolumeAttachmentStatus(volumeAttachment)
}

// VolumeAttachment APIs - END

// MutatingWebhookConfig APIS - START

// GetMutatingWebhookConfiguration returns a given MutatingWebhookConfiguration
func (k *k8sOps) GetMutatingWebhookConfiguration(name string) (*hook.MutatingWebhookConfiguration, error) {
	return admissionregistration.Instance().GetMutatingWebhookConfiguration(name)
}

// CreateMutatingWebhookConfiguration creates given MutatingWebhookConfiguration
func (k *k8sOps) CreateMutatingWebhookConfiguration(cfg *hook.MutatingWebhookConfiguration) (*hook.MutatingWebhookConfiguration, error) {
	return admissionregistration.Instance().CreateMutatingWebhookConfiguration(cfg)
}

// UpdateMutatingWebhookConfiguration updates given MutatingWebhookConfiguration
func (k *k8sOps) UpdateMutatingWebhookConfiguration(cfg *hook.MutatingWebhookConfiguration) (*hook.MutatingWebhookConfiguration, error) {
	return admissionregistration.Instance().UpdateMutatingWebhookConfiguration(cfg)
}

// DeleteMutatingWebhookConfiguration deletes given MutatingWebhookConfiguration
func (k *k8sOps) DeleteMutatingWebhookConfiguration(name string) error {
	return admissionregistration.Instance().DeleteMutatingWebhookConfiguration(name)
}

// MutatingWebhookConfig APIS - END
