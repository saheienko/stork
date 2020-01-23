package k8s

import (
	prometheusclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	snap_client "github.com/kubernetes-incubator/external-storage/snapshot/pkg/client"
	autopilotclientset "github.com/libopenstorage/autopilot-api/pkg/client/clientset/versioned"
	ostclientset "github.com/libopenstorage/operator/pkg/client/clientset/versioned"
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
	talismanclientset "github.com/portworx/talisman/pkg/client/clientset/versioned"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	dynamicclient "k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ClientSetter is an interface to allow setting different clients on the Ops object
type ClientSetter interface {
	// SetConfig sets the config and resets the client
	SetConfig(config *rest.Config)
	// SetConfigFromPath sets the config from a kubeconfig file
	SetConfigFromPath(configPath string) error
	// SetClient set the k8s clients
	SetClient(
		kubernetes.Interface,
		rest.Interface,
		storkclientset.Interface,
		apiextensionsclient.Interface,
		dynamicclient.Interface,
		ocp_clientset.Interface,
		ocp_security_clientset.Interface,
		autopilotclientset.Interface,
	)
	// SetBaseClient sets the kubernetes clientset
	SetBaseClient(kubernetes.Interface)
	// SetSnapshotClient sets the snapshot clientset
	SetSnapshotClient(rest.Interface)
	// SetStorkClient sets the stork clientset
	SetStorkClient(storkclientset.Interface)
	// SetOpenstorageOperatorClient sets the openstorage operator clientset
	SetOpenstorageOperatorClient(ostclientset.Interface)
	// SetAPIExtensionsClient sets the api extensions clientset
	SetAPIExtensionsClient(apiextensionsclient.Interface)
	// SetDynamicClient sets the dynamic clientset
	SetDynamicClient(dynamicclient.Interface)
	// SetOpenshiftAppsClient sets the openshift apps clientset
	SetOpenshiftAppsClient(ocp_clientset.Interface)
	// SetOpenshiftSecurityClient sets the openshift security clientset
	SetOpenshiftSecurityClient(ocp_security_clientset.Interface)
	// SetTalismanClient sets the talisman clientset
	SetTalismanClient(talismanclientset.Interface)
	// SetAutopilotClient sets the autopilot clientset
	SetAutopilotClient(autopilotclientset.Interface)
	// SetPrometheusClient sets the prometheus clientset
	SetPrometheusClient(prometheusclient.Interface)
}

// SetConfig sets the config and resets the client
func (k *k8sOps) SetConfig(config *rest.Config) {
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
}

// SetConfigFromPath takes the path to a kubeconfig file
// and then internally calls SetConfig to set it
func (k *k8sOps) SetConfigFromPath(configPath string) error {
	if configPath == "" {
		k.SetConfig(nil)
		return nil
	}
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return err
	}

	k.SetConfig(config)
	return nil
}

// SetClient set the k8s clients
func (k *k8sOps) SetClient(
	client kubernetes.Interface,
	snapClient rest.Interface,
	storkClient storkclientset.Interface,
	apiExtensionClient apiextensionsclient.Interface,
	dynamicInterface dynamicclient.Interface,
	ocpClient ocp_clientset.Interface,
	ocpSecurityClient ocp_security_clientset.Interface,
	autopilotClient autopilotclientset.Interface,
) {
	admissionregistration.SetInstance(admissionregistration.New(client.AdmissionregistrationV1beta1()))
	apiextensions.SetInstance(apiextensions.New(apiExtensionClient))
	apps.SetInstance(apps.New(client.AppsV1(), client.CoreV1()))
	autopilot.SetInstance(autopilot.New(autopilotClient))
	batch.SetInstance(batch.New(client.BatchV1()))
	core.SetInstance(core.New(client.CoreV1(), client.StorageV1()))
	discovery.SetInstance(discovery.New(client.Discovery()))
	dynamic.SetInstance(dynamic.New(dynamicInterface))
	externalstorage.SetInstance(externalstorage.New(snapClient))
	openshift.SetInstance(openshift.New(client, ocpClient, ocpSecurityClient))
	rbac.SetInstance(rbac.New(client.RbacV1()))
	storage.SetInstance(storage.New(client.StorageV1(), client.StorageV1beta1()))
	stork.SetInstance(stork.New(client, storkClient, snapClient))
}

// SetBaseClient sets the kubernetes clientset
func (k *k8sOps) SetBaseClient(client kubernetes.Interface) {
	admissionregistration.SetInstance(admissionregistration.New(client.AdmissionregistrationV1beta1()))
	apps.SetInstance(apps.New(client.AppsV1(), client.CoreV1()))
	batch.SetInstance(batch.New(client.BatchV1()))
	core.SetInstance(core.New(client.CoreV1(), client.StorageV1()))
	discovery.SetInstance(discovery.New(client.Discovery()))
	rbac.SetInstance(rbac.New(client.RbacV1()))
	storage.SetInstance(storage.New(client.StorageV1(), client.StorageV1beta1()))
}

// SetSnapshotClient sets the snapshot clientset
func (k *k8sOps) SetSnapshotClient(snapClient rest.Interface) {
	externalstorage.SetInstance(externalstorage.New(snapClient))
}

// SetStorkClient sets the stork clientset
func (k *k8sOps) SetStorkClient(storkClient storkclientset.Interface) {
}

// SetOpenstorageOperatorClient sets the openstorage operator clientset
func (k *k8sOps) SetOpenstorageOperatorClient(ostClient ostclientset.Interface) {
	operator.SetInstance(operator.New(ostClient))
}

// SetAPIExtensionsClient sets the api extensions clientset
func (k *k8sOps) SetAPIExtensionsClient(apiExtensionsClient apiextensionsclient.Interface) {
	apiextensions.SetInstance(apiextensions.New(apiExtensionsClient))
}

// SetDynamicClient sets the dynamic clientset
func (k *k8sOps) SetDynamicClient(dynamicClient dynamicclient.Interface) {
	dynamic.SetInstance(dynamic.New(dynamicClient))
}

// SetOpenshiftAppsClient sets the openshift apps clientset
func (k *k8sOps) SetOpenshiftAppsClient(ocpAppsClient ocp_clientset.Interface) {
}

// SetOpenshiftSecurityClient sets the openshift security clientset
func (k *k8sOps) SetOpenshiftSecurityClient(ocpSecurityClient ocp_security_clientset.Interface) {
}

// SetAutopilotClient sets the autopilot clientset
func (k *k8sOps) SetAutopilotClient(autopilotClient autopilotclientset.Interface) {
	autopilot.SetInstance(autopilot.New(autopilotClient))
}

// SetTalismanClient sets the talisman clientset
func (k *k8sOps) SetTalismanClient(talismanClient talismanclientset.Interface) {
	talisman.SetInstance(talisman.New(talismanClient))
}

// SetPrometheusClient sets the prometheus clientset
func (k *k8sOps) SetPrometheusClient(prometheusClient prometheusclient.Interface) {
	prometheus.SetInstance(prometheus.New(prometheusClient))
}

func (k *k8sOps) loadClientFromKubeconfig(kubeconfig string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	return k.loadClientFor(config)
}

func (k *k8sOps) loadClientFromConfigBytes(kubeconfig []byte) error {
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		return err
	}

	return k.loadClientFor(config)
}

func (k *k8sOps) loadClientFor(config *rest.Config) error {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	snapClient, _, err := snap_client.NewClient(config)
	if err != nil {
		return err
	}

	storkClient, err := storkclientset.NewForConfig(config)
	if err != nil {
		return err
	}

	ostClient, err := ostclientset.NewForConfig(config)
	if err != nil {
		return err
	}

	talismanClient, err := talismanclientset.NewForConfig(config)
	if err != nil {
		return err
	}

	apiExtensionClient, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return err
	}

	dynamicInterface, err := dynamicclient.NewForConfig(config)
	if err != nil {
		return err
	}

	ocpClient, err := ocp_clientset.NewForConfig(config)
	if err != nil {
		return err
	}

	ocpSecurityClient, err := ocp_security_clientset.NewForConfig(config)
	if err != nil {
		return err
	}

	autopilotClient, err := autopilotclientset.NewForConfig(config)
	if err != nil {
		return err
	}

	prometheusClient, err := prometheusclient.NewForConfig(config)
	if err != nil {
		return err
	}

	admissionregistration.SetInstance(admissionregistration.New(client.AdmissionregistrationV1beta1()))
	apiextensions.SetInstance(apiextensions.New(apiExtensionClient))
	apps.SetInstance(apps.New(client.AppsV1(), client.CoreV1()))
	autopilot.SetInstance(autopilot.New(autopilotClient))
	batch.SetInstance(batch.New(client.BatchV1()))
	core.SetInstance(core.New(client.CoreV1(), client.StorageV1()))
	discovery.SetInstance(discovery.New(client.Discovery()))
	dynamic.SetInstance(dynamic.New(dynamicInterface))
	externalstorage.SetInstance(externalstorage.New(snapClient))
	openshift.SetInstance(openshift.New(client, ocpClient, ocpSecurityClient))
	operator.SetInstance(operator.New(ostClient))
	prometheus.SetInstance(prometheus.New(prometheusClient))
	rbac.SetInstance(rbac.New(client.RbacV1()))
	storage.SetInstance(storage.New(client.StorageV1(), client.StorageV1beta1()))
	stork.SetInstance(stork.New(client, storkClient, snapClient))
	talisman.SetInstance(talisman.New(talismanClient))

	return nil
}
