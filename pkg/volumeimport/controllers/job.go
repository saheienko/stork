package controllers

import (
	"context"
	"reflect"
	"time"

	"github.com/libopenstorage/stork/pkg/controller"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/apis/core"
)

type JobController struct {
}

func NewJobController() (*JobController, error) {
	return &JobController{}, nil
}

func (c *JobController) Init() error {
	return controller.Register(
		&schema.GroupVersionKind{
			Group:   batchv1.GroupName,
			Version: batchv1.SchemeGroupVersion.Version,
			Kind:    reflect.TypeOf(batchv1.Job{}).Name(),
		},
		core.NamespaceAll,
		5*time.Minute,
		c,
	)
}

func (c *JobController) Handle(ctx context.Context, event sdk.Event) error {
	return ReconcileProtected(ctx, event)
}
