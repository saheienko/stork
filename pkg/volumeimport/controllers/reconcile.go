package controllers

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/libopenstorage/stork/pkg/apis/stork"
	storkapi "github.com/libopenstorage/stork/pkg/apis/stork/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/portworx/sched-ops/k8s"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	LabelController     = "stork.libopenstorage.org/controller"
	LabelControllerName = "controller-name"
	VolumeImportName    = "volume-import"
)

const (
	ConditionTypeCompleted = "Completed"
)

var (
	mu sync.Mutex
)

// controllerKind contains the schema.GroupVersionKind for this controller type.
var controllerKind = storkapi.SchemeGroupVersion.WithKind(reflect.TypeOf(storkapi.VolumeImport{}).Name())

func ReconcileProtected(ctx context.Context, event sdk.Event) error {
	mu.Lock()
	defer mu.Unlock()

	return Reconcile(ctx, event)
}

func Reconcile(ctx context.Context, event sdk.Event) error {
	if event.Deleted {
		return nil
	}

	var vi *storkapi.VolumeImport
	var err error
	switch o := event.Object.(type) {
	case *storkapi.VolumeImport:
		vi = o
	case *batchv1.Job:
		if c := getControllerName(o.Labels); c != VolumeImportName {
			return nil
		}

		if vi, err = getVolumeImportFor(o); err != nil {
			if errors.IsNotFound(err) {
				return sdk.Delete(o)
			}
			return err
		}
	}

	return reconcile(vi)
}

func getVolumeImportFor(job *batchv1.Job) (*storkapi.VolumeImport, error) {
	controllerRef := metav1.GetControllerOf(job)
	if controllerRef == nil {
		// TODO: should such jobs be handled? (it has proper labels)
		return nil, fmt.Errorf("job %s/%s has no controllerRef", job.Namespace, job.Name)
	}

	vi := &storkapi.VolumeImport{
		TypeMeta: metav1.TypeMeta{
			Kind:       "VolumeImport",
			APIVersion: strings.Join([]string{stork.GroupName, storkapi.SchemeGroupVersion.Version}, "/"),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      controllerRef.Name,
			Namespace: job.Namespace,
		},
	}

	if err := sdk.Get(vi); err != nil {
		return nil, err
	}

	return vi, nil
}

func reconcile(o *storkapi.VolumeImport) error {
	if o == nil {
		return nil
	}

	if o.DeletionTimestamp != nil {
		return deleteJob(o.Status.RsyncJobName, o.Namespace)
	}

	logrus.Debugf("handling %s/%s VolumeImport", o.Namespace, o.Name)

	// check if src/dst volumes is not mounted
	if err := checkClaims(o); err != nil {
		// TODO: update status
		return err
	}

	// check if a rsync job is created
	viJob, err := getJob(toJobName(o.Name), o.Namespace)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}

		// create a job if it's not exist
		viJob = jobFrom(o)
		if err = sdk.Create(viJob); err != nil {
			return err
		}
	}

	// update the Volume import status
	return updateStatus(o, viJob)
}

func updateStatus(vi *storkapi.VolumeImport, job *batchv1.Job) error {
	vi.Status.RsyncJobName = job.Name
	if isJobCompleted(job) {
		vi.Status.ConditionType = ConditionTypeCompleted
	}
	return sdk.Update(vi)
}

func checkClaims(vi *storkapi.VolumeImport) error {
	if err := ensureUnmountedPVC(vi.Spec.Source.Name, vi.Spec.Source.Namespace, vi.Name); err != nil {
		return fmt.Errorf("source pvc: %s/%s: %v", vi.Spec.Source.Namespace, vi.Spec.Source.Name, err)
	}

	dstMeta := vi.Spec.Destination.PVC.Metadata
	if err := ensureUnmountedPVC(dstMeta.Name, dstMeta.Namespace, vi.Name); err != nil {
		// return an error if pvc is not exist or spec is empty
		if !errors.IsNotFound(err) || vi.Spec.Destination.PVC.Spec == nil {
			return fmt.Errorf("destination pvc: %s/%s: %v", dstMeta.Namespace, dstMeta.Name, err)
		}

		// otherwise create pvc with provided spec
		_, err = createPVC(dstMeta.Name, dstMeta.Namespace, vi.Spec.Destination.PVC.Spec)
		return err
	}

	return nil
}

func ensureUnmountedPVC(name, namespace, viName string) error {
	pvc, err := getPVC(name, namespace)
	if err != nil {
		return err
	}
	if pvc.Status.Phase != corev1.ClaimBound {
		return fmt.Errorf("status: expected %s, got %s", corev1.ClaimBound, pvc.Status.Phase)
	}

	// check if pvc is mounted
	pods, err := getMountPods(pvc.Name, pvc.Namespace)
	if err != nil {
		return fmt.Errorf("get mounted pods: %v", err)
	}
	mounted := make([]corev1.Pod, 0)
	for _, pod := range pods {
		// pvc is mounted to pod created for this volume import
		if pod.Labels[LabelControllerName] == viName {
			continue
		}
		mounted = append(mounted, pod)
	}
	if len(mounted) > 0 {
		return fmt.Errorf("mounted to %v pods", toPodNames(pods))
	}

	return nil
}

func getMountPods(pvcName, namespace string) ([]corev1.Pod, error) {
	return k8s.Instance().GetPodsUsingPVC(pvcName, namespace)
}

func toPodNames(objs []corev1.Pod) []string {
	out := make([]string, 0, 0)
	for _, o := range objs {
		out = append(out, o.Name)
	}
	return out
}

func jobLabels(volumeImportName string) map[string]string {
	return map[string]string{
		LabelController:     VolumeImportName,
		LabelControllerName: volumeImportName,
	}
}

func toJobName(volumeImportName string) string {
	return fmt.Sprintf("job-%s", volumeImportName)
}

func getControllerName(labels map[string]string) string {
	if labels == nil {
		return ""
	}
	return labels[LabelController]
}

func isJobCompleted(j *batchv1.Job) bool {
	for _, c := range j.Status.Conditions {
		if c.Type == batchv1.JobComplete && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func jobFrom(vi *storkapi.VolumeImport) *batchv1.Job {
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            toJobName(vi.Name),
			Namespace:       vi.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(vi, controllerKind)},
			Labels:          jobLabels(vi.Name),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:     jobLabels(vi.Name),
					Finalizers: nil,
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{
						{
							Name:    "rsync",
							Image:   "eeacms/rsync",
							Command: []string{"/bin/sh", "-c", "rsync -avz /src/ /dst"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "src-vol",
									MountPath: "/src",
								},
								{
									Name:      "dst-vol",
									MountPath: "/dst",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "src-vol",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: vi.Spec.Source.Name,
								},
							},
						},
						{
							Name: "dst-vol",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: vi.Spec.Destination.PVC.Metadata.Name,
								},
							},
						},
					},
				},
			},
		},
	}
}
