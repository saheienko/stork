package k8s

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfigMap(name, namespace string) (*corev1.ConfigMap, error) {
	if err := validate(name, namespace); err != nil {
		return nil, err
	}
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if err := sdk.Get(cm); err != nil {
		return nil, err
	}
	return cm, nil
}

func CreateConfigMap(cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	if err := sdk.Create(cm); err != nil {
		return nil, err
	}
	return cm, nil
}
