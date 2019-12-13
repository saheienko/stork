package controllers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/libopenstorage/stork/drivers/volume"
	stork_api "github.com/libopenstorage/stork/pkg/apis/stork/v1alpha1"
	"github.com/portworx/sched-ops/k8s"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	validateCRDInterval time.Duration = 5 * time.Second
	validateCRDTimeout  time.Duration = 1 * time.Minute
)

func NewClusterPair(mgr manager.Manager, d volume.Driver, r record.EventRecorder) *ClusterPairController {
	return &ClusterPairController{
		client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		Driver:   d,
		Recorder: r,
	}
}

// ClusterPairController controller to watch over ClusterPair
type ClusterPairController struct {
	client runtimeclient.Client
	scheme *runtime.Scheme

	Driver   volume.Driver
	Recorder record.EventRecorder
}

// Init initialize the cluster pair controller
func (c *ClusterPairController) Init(mgr manager.Manager) error {
	err := c.createCRD()
	if err != nil {
		return err
	}

	// Create a new controller
	ctrl, err := controller.New("cluster-pair-controller", mgr, controller.Options{Reconciler: c})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Migration
	return ctrl.Watch(&source.Kind{Type: &stork_api.ClusterPair{}}, &handler.EnqueueRequestForObject{})
}

func (a *ClusterPairController) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logrus.Printf("Reconciling ClusterPair %s/%s", request.Namespace, request.Name)

	// Fetch the ApplicationBackup instance
	backup := &stork_api.ClusterPair{}
	err := a.client.Get(context.TODO(), request.NamespacedName, backup)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, a.handle(context.TODO(), backup)
}

func (c *ClusterPairController) handle(ctx context.Context, clusterPair *stork_api.ClusterPair) error {
	if clusterPair.DeletionTimestamp != nil {
		if clusterPair.Status.RemoteStorageID != "" {
			return c.Driver.DeletePair(clusterPair)
		}
		return nil
	}

	if len(clusterPair.Spec.Options) == 0 {
		clusterPair.Status.StorageStatus = stork_api.ClusterPairStatusNotProvided
		c.Recorder.Event(clusterPair,
			v1.EventTypeNormal,
			string(clusterPair.Status.StorageStatus),
			"Skipping storage pairing since no storage options provided")
		err := c.client.Update(context.TODO(), clusterPair)
		if err != nil {
			return err
		}
	} else {
		if clusterPair.Status.StorageStatus != stork_api.ClusterPairStatusReady {
			remoteID, err := c.Driver.CreatePair(clusterPair)
			if err != nil {
				clusterPair.Status.StorageStatus = stork_api.ClusterPairStatusError
				c.Recorder.Event(clusterPair,
					v1.EventTypeWarning,
					string(clusterPair.Status.StorageStatus),
					err.Error())
			} else {
				clusterPair.Status.StorageStatus = stork_api.ClusterPairStatusReady
				c.Recorder.Event(clusterPair,
					v1.EventTypeNormal,
					string(clusterPair.Status.StorageStatus),
					"Storage successfully paired")
				clusterPair.Status.RemoteStorageID = remoteID
			}
			err = c.client.Update(context.TODO(), clusterPair)
			if err != nil {
				return err
			}
		}
	}
	if clusterPair.Status.SchedulerStatus != stork_api.ClusterPairStatusReady {
		remoteConfig, err := getClusterPairSchedulerConfig(clusterPair.Name, clusterPair.Namespace)
		if err != nil {
			return err
		}

		client, err := kubernetes.NewForConfig(remoteConfig)
		if err != nil {
			return err
		}
		if _, err = client.ServerVersion(); err != nil {
			clusterPair.Status.SchedulerStatus = stork_api.ClusterPairStatusError
			c.Recorder.Event(clusterPair,
				v1.EventTypeWarning,
				string(clusterPair.Status.SchedulerStatus),
				err.Error())
		} else {
			clusterPair.Status.SchedulerStatus = stork_api.ClusterPairStatusReady
			c.Recorder.Event(clusterPair,
				v1.EventTypeNormal,
				string(clusterPair.Status.SchedulerStatus),
				"Scheduler successfully paired")
		}
		err = c.client.Update(context.TODO(), clusterPair)
		if err != nil {
			return err
		}
	}

	return nil
}

func getClusterPairSchedulerConfig(clusterPairName string, namespace string) (*restclient.Config, error) {
	clusterPair, err := k8s.Instance().GetClusterPair(clusterPairName, namespace)
	if err != nil {
		return nil, fmt.Errorf("error getting clusterpair: %v", err)
	}
	remoteClientConfig := clientcmd.NewNonInteractiveClientConfig(
		clusterPair.Spec.Config,
		clusterPair.Spec.Config.CurrentContext,
		&clientcmd.ConfigOverrides{},
		clientcmd.NewDefaultClientConfigLoadingRules())
	return remoteClientConfig.ClientConfig()
}

func getClusterPairStorageStatus(clusterPairName string, namespace string) (stork_api.ClusterPairStatusType, error) {
	clusterPair, err := k8s.Instance().GetClusterPair(clusterPairName, namespace)
	if err != nil {
		return stork_api.ClusterPairStatusInitial, fmt.Errorf("error getting clusterpair: %v", err)
	}
	return clusterPair.Status.StorageStatus, nil
}

func getClusterPairSchedulerStatus(clusterPairName string, namespace string) (stork_api.ClusterPairStatusType, error) {
	clusterPair, err := k8s.Instance().GetClusterPair(clusterPairName, namespace)
	if err != nil {
		return stork_api.ClusterPairStatusInitial, fmt.Errorf("error getting clusterpair: %v", err)
	}
	return clusterPair.Status.SchedulerStatus, nil
}

func (c *ClusterPairController) createCRD() error {
	resource := k8s.CustomResource{
		Name:    stork_api.ClusterPairResourceName,
		Plural:  stork_api.ClusterPairResourcePlural,
		Group:   stork_api.SchemeGroupVersion.Group,
		Version: stork_api.SchemeGroupVersion.Version,
		Scope:   apiextensionsv1beta1.NamespaceScoped,
		Kind:    reflect.TypeOf(stork_api.ClusterPair{}).Name(),
	}
	err := k8s.Instance().CreateCRD(resource)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return k8s.Instance().ValidateCRD(resource, validateCRDTimeout, validateCRDInterval)
}
