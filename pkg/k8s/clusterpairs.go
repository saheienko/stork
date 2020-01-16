package k8s

import (
	"github.com/libopenstorage/stork/pkg/apis/stork/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetClusterPair(name, namespace string) (*v1alpha1.ClusterPair, error) {
	if err := validate(name, namespace); err != nil {
		return nil, err
	}
	cp := &v1alpha1.ClusterPair{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if err := sdk.Get(cp); err != nil {
		return nil, err
	}
	return cp, nil
}
