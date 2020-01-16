package k8s

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetStorageClass(name string) (*storagev1.StorageClass, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}
	sc := &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	if err := sdk.Get(sc); err != nil {
		return nil, err
	}
	return sc, nil
}
