package controllers

import (
	"fmt"
	"strings"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getJob(name, namespace string) (*batchv1.Job, error) {
	if err := checkMetadata(name, namespace); err != nil {
		return nil, err
	}

	into := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if err := sdk.Get(into); err != nil {
		return nil, err
	}
	return into, nil
}

func deleteJob(name, namespace string) error {
	if err := checkMetadata(name, namespace); err != nil {
		return err
	}

	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if err := sdk.Delete(job); !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func getPVC(name, namespace string) (*corev1.PersistentVolumeClaim, error) {
	if err := checkMetadata(name, namespace); err != nil {
		return nil, err
	}

	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	if err := sdk.Get(pvc); err != nil {
		return nil, err
	}

	return pvc, nil
}

func createPVC(name, namespace string, spec *corev1.PersistentVolumeClaimSpec) (*corev1.PersistentVolumeClaim, error) {
	if err := checkMetadata(name, namespace); err != nil {
		return nil, err
	}
	if spec == nil {
		return nil, fmt.Errorf("spec should be provided")
	}

	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: *spec,
	}
	if err := sdk.Create(pvc); err != nil {
		return nil, err
	}

	return pvc, nil
}

func checkMetadata(name, namespace string) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(namespace) == "" {
		return fmt.Errorf("name and namespace should not be empty")
	}
	return nil
}
