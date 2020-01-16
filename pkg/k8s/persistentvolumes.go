package k8s

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPersistentVolume(name string) (*corev1.PersistentVolume, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
		},
	}
	if err := sdk.Get(pv); err != nil {
		return nil, err
	}
	return pv, nil
}
