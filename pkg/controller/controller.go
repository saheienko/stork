package controller

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type controllerNotInitError struct {
}

func (c *controllerNotInitError) Error() string {
	return "Contoller has not been initialized"
}

// Controller to track updates for objects
type controller struct {
	sync.Mutex
	started  bool
	handlers map[string][]sdk.Handler
}

var controllerInst *controller

// Init the controller
func Init() error {
	if controllerInst != nil {
		return nil
	}
	controllerInst = &controller{}
	controllerInst.handlers = make(map[string][]sdk.Handler)
	sdk.Handle(controllerInst)
	return nil
}

// Run the controller
func Run() error {
	if controllerInst == nil {
		return &controllerNotInitError{}
	}
	controllerInst.Lock()
	defer controllerInst.Unlock()
	go sdk.Run(context.TODO())
	controllerInst.started = true
	return nil
}

// Register registers a controller with empty label selector.
func Register(
	gkv *schema.GroupVersionKind,
	namespace string,
	resyncPeriod time.Duration,
	handler sdk.Handler,
) error {
	return RegisterWithLabelSelector(gkv, namespace, resyncPeriod, handler, "")
}

// RegisterWithLabelSelector registers a controller that will handle updates on resources which
// satisfy a label selector.
func RegisterWithLabelSelector(
	gkv *schema.GroupVersionKind,
	namespace string,
	resyncPeriod time.Duration,
	handler sdk.Handler,
	labelSelector string,
) error {
	logrus.Debugf("Registering controller for %v", gkv)
	if controllerInst == nil {
		return &controllerNotInitError{}
	}
	controllerInst.Lock()
	defer controllerInst.Unlock()
	if controllerInst.started {
		return fmt.Errorf("can't register new handlers after starting controller")
	}

	objectType := gkv.String()
	// Only add Watch if we aren't already watching for the object already.
	// resyncPeriod will be ignored for second call if different
	if controllerInst.handlers[objectType] == nil {
		controllerInst.handlers[objectType] = make([]sdk.Handler, 0)
		sdk.Watch(gkv.GroupVersion().String(), gkv.Kind, namespace, resyncPeriod, sdk.WithNumWorkers(10), sdk.WithLabelSelector(labelSelector))
	}
	logrus.Debugf("Registered controller for %v", gkv)
	controllerInst.handlers[objectType] = append(controllerInst.handlers[objectType], handler)
	return nil
}

// Handle handles updates for registered types
func (c *controller) Handle(ctx context.Context, event sdk.Event) error {
	gkv := event.Object.GetObjectKind().GroupVersionKind().String()
	var firstErr error
	if handlers, ok := c.handlers[gkv]; ok {
		for _, handler := range handlers {
			err := handler.Handle(ctx, event)
			// Keep track of the first error
			if err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}
