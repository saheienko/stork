package k8s

import (
	"github.com/portworx/sched-ops/k8s/talisman"
	talisman_v1beta2 "github.com/portworx/talisman/pkg/apis/portworx/v1beta2"
)

// VolumePlacementStrategyOps is an interface to perform CRUD volume placememt strategy ops
type VolumePlacementStrategyOps interface {
	// CreateVolumePlacementStrategy creates a new volume placement strategy
	CreateVolumePlacementStrategy(spec *talisman_v1beta2.VolumePlacementStrategy) (*talisman_v1beta2.VolumePlacementStrategy, error)
	// UpdateVolumePlacementStrategy updates an existing volume placement strategy
	UpdateVolumePlacementStrategy(spec *talisman_v1beta2.VolumePlacementStrategy) (*talisman_v1beta2.VolumePlacementStrategy, error)
	// ListVolumePlacementStrategies lists all volume placement strategies
	ListVolumePlacementStrategies() (*talisman_v1beta2.VolumePlacementStrategyList, error)
	// DeleteVolumePlacementStrategy deletes the volume placement strategy with given name
	DeleteVolumePlacementStrategy(name string) error
	// GetVolumePlacementStrategy returns the volume placememt strategy with given name
	GetVolumePlacementStrategy(name string) (*talisman_v1beta2.VolumePlacementStrategy, error)
}

// VolumePlacementStrategy APIs - BEGIN

func (k *k8sOps) CreateVolumePlacementStrategy(spec *talisman_v1beta2.VolumePlacementStrategy) (*talisman_v1beta2.VolumePlacementStrategy, error) {
	return talisman.Instance().CreateVolumePlacementStrategy(spec)
}

func (k *k8sOps) UpdateVolumePlacementStrategy(spec *talisman_v1beta2.VolumePlacementStrategy) (*talisman_v1beta2.VolumePlacementStrategy, error) {
	return talisman.Instance().UpdateVolumePlacementStrategy(spec)
}

func (k *k8sOps) ListVolumePlacementStrategies() (*talisman_v1beta2.VolumePlacementStrategyList, error) {
	return talisman.Instance().ListVolumePlacementStrategies()
}

func (k *k8sOps) DeleteVolumePlacementStrategy(name string) error {
	return talisman.Instance().DeleteVolumePlacementStrategy(name)
}

func (k *k8sOps) GetVolumePlacementStrategy(name string) (*talisman_v1beta2.VolumePlacementStrategy, error) {
	return talisman.Instance().GetVolumePlacementStrategy(name)
}

// VolumePlacementStrategy APIs - END
