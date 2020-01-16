package k8s

import (
	"fmt"
	snapv1 "github.com/kubernetes-incubator/external-storage/snapshot/pkg/apis/crd/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/portworx/sched-ops/task"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func GetSnapshotData(name string) (*snapv1.VolumeSnapshotData, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}
	snapData := &snapv1.VolumeSnapshotData{
		Metadata: metav1.ObjectMeta{
			Name: name,
		},
	}
	if err := sdk.Get(snapData); err != nil {
		return nil, err
	}
	return snapData, nil
}

func GetSnapshot(name, namespace string) (*snapv1.VolumeSnapshot, error) {
	if err := validate(name, namespace); err != nil {
		return nil, err
	}
	snap := &snapv1.VolumeSnapshot{
		Metadata: metav1.ObjectMeta{
			Name: name,
		},
	}
	if err := sdk.Get(snap); err != nil {
		return nil, err
	}
	return snap, nil
}

func DeleteSnapshot(name, namespace string) error {
	if err := validate(name, namespace); err != nil {
		return err
	}
	snap := &snapv1.VolumeSnapshot{
		Metadata: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return sdk.Delete(snap)
}

func ValidateSnapshotData(name string, retry bool, timeout, retryInterval time.Duration) error {
	t := func() (interface{}, bool, error) {
		snapData, err := GetSnapshotData(name)
		if err != nil {
			return "", true, err
		}

		for _, condition := range snapData.Status.Conditions {
			if condition.Status == v1.ConditionTrue {
				if condition.Type == snapv1.VolumeSnapshotDataConditionReady {
					return "", false, nil
				} else if condition.Type == snapv1.VolumeSnapshotDataConditionError {
					return "", true, fmt.Errorf("SnapshotData is failed, ID=%s Status %v", name, snapData.Status)
				}
			}
		}

		return "", true, fmt.Errorf("SnapshotData is not ready, ID=%s Status %v", name, snapData.Status)
	}

	if retry {
		if _, err := task.DoRetryWithTimeout(t, timeout, retryInterval); err != nil {
			return err
		}
	} else {
		if _, _, err := t(); err != nil {
			return err
		}
	}

	return nil
}
