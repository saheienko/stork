package k8s

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func GetPersistentVolumeClaim(name, namespace string) (*corev1.PersistentVolumeClaim, error) {
	if err := validate(name, namespace); err != nil {
		return nil, err
	}
	pvc := &corev1.PersistentVolumeClaim{
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

func ListPersistentVolumeClaims(namespace string, selector map[string]string) (*corev1.PersistentVolumeClaimList, error) {
	if err := validateNamespace(namespace); err != nil {
		return nil, err
	}
	pvcList := &corev1.PersistentVolumeClaimList{}
	opts := sdk.WithListOptions(&metav1.ListOptions{
		LabelSelector: labels.FormatLabels(selector),
	})
	if err := sdk.List(namespace, pvcList, opts); err != nil {
		return nil, err
	}
	return pvcList, nil
}
