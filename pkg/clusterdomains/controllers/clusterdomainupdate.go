package controllers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/libopenstorage/stork/drivers/volume"
	storkv1 "github.com/libopenstorage/stork/pkg/apis/stork/v1alpha1"
	"github.com/libopenstorage/stork/pkg/log"
	"github.com/portworx/sched-ops/k8s/apiextensions"
	"github.com/portworx/sched-ops/k8s/stork"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func NewClusterDomainUpdate(mgr manager.Manager, d volume.Driver, r record.EventRecorder) *ClusterDomainUpdateController {
	return &ClusterDomainUpdateController{
		client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		Driver:   d,
		Recorder: r,
	}
}

// ClusterDomainUpdateController clusterdomainupdate controller
type ClusterDomainUpdateController struct {
	client runtimeclient.Client
	scheme *runtime.Scheme

	Driver   volume.Driver
	Recorder record.EventRecorder
}

// Init initialize the clusterdomainupdate controller
func (c *ClusterDomainUpdateController) Init(mgr manager.Manager) error {
	err := c.createCRD()
	if err != nil {
		return err
	}

	// Create a new controller
	ctrl, err := controller.New("cluster-domain-update-controller", mgr, controller.Options{Reconciler: c})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Migration
	return ctrl.Watch(&source.Kind{Type: &storkv1.ClusterDomainUpdate{}}, &handler.EnqueueRequestForObject{})
}

func (c *ClusterDomainUpdateController) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logrus.Printf("Reconciling ClusterDomainUpdate %s/%s", request.Namespace, request.Name)

	// Fetch the ApplicationBackup instance
	clusterDomainUpdate := &storkv1.ClusterDomainUpdate{}
	err := c.client.Get(context.TODO(), request.NamespacedName, clusterDomainUpdate)
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

	return reconcile.Result{}, c.handle(context.TODO(), clusterDomainUpdate)
}

func (c *ClusterDomainUpdateController) handle(ctx context.Context, clusterDomainUpdate *storkv1.ClusterDomainUpdate) error {
	if clusterDomainUpdate.DeletionTimestamp != nil {
		return nil
	}

	switch clusterDomainUpdate.Status.Status {
	case storkv1.ClusterDomainUpdateStatusInitial:
		var (
			action string
			err    error
		)
		if clusterDomainUpdate.Spec.Active {
			action = "activate"
			err = c.Driver.ActivateClusterDomain(clusterDomainUpdate)
		} else {
			action = "deactivate"
			err = c.Driver.DeactivateClusterDomain(clusterDomainUpdate)
		}
		if err != nil {
			err = fmt.Errorf("unable to %v cluster domain: %v", action, err)
			log.ClusterDomainUpdateLog(clusterDomainUpdate).Errorf(err.Error())
			clusterDomainUpdate.Status.Status = storkv1.ClusterDomainUpdateStatusFailed
			clusterDomainUpdate.Status.Reason = err.Error()
			c.Recorder.Event(
				clusterDomainUpdate,
				v1.EventTypeWarning,
				string(storkv1.ClusterDomainUpdateStatusFailed),
				err.Error(),
			)

		} else {
			clusterDomainUpdate.Status.Status = storkv1.ClusterDomainUpdateStatusSuccessful
		}

		err = c.client.Update(context.TODO(), clusterDomainUpdate)
		if err != nil {
			log.ClusterDomainUpdateLog(clusterDomainUpdate).Errorf("Error updating ClusterDomainUpdate: %v", err)
			return err
		}
		// Do a dummy update on the cluster domain status so that it queries
		// the storage driver and gets updated too
		if clusterDomainUpdate.Status.Status == storkv1.ClusterDomainUpdateStatusSuccessful {
			cdsList, err := stork.Instance().ListClusterDomainStatuses()
			if err != nil {
				return err
			}
			for _, cds := range cdsList.Items {
				_, err := stork.Instance().UpdateClusterDomainsStatus(&cds)
				if err != nil {
					return err
				}
			}

		}
		return nil
	case storkv1.ClusterDomainUpdateStatusFailed, storkv1.ClusterDomainUpdateStatusSuccessful:
		return nil
	default:
		log.ClusterDomainUpdateLog(clusterDomainUpdate).Errorf("Invalid status for cluster domain update: %v", clusterDomainUpdate.Status.Status)
	}

	return nil
}

// createCRD creates the CRD for ClusterDomainsStatus object
func (c *ClusterDomainUpdateController) createCRD() error {
	resource := apiextensions.CustomResource{
		Name:       storkv1.ClusterDomainUpdateResourceName,
		Plural:     storkv1.ClusterDomainUpdatePlural,
		Group:      storkv1.SchemeGroupVersion.Group,
		Version:    storkv1.SchemeGroupVersion.Version,
		Scope:      apiextensionsv1beta1.ClusterScoped,
		Kind:       reflect.TypeOf(storkv1.ClusterDomainUpdate{}).Name(),
		ShortNames: []string{storkv1.ClusterDomainUpdateShortName},
	}
	err := apiextensions.Instance().CreateCRD(resource)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return apiextensions.Instance().ValidateCRD(resource, validateCRDTimeout, validateCRDInterval)
}
