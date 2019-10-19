package gcp

import (
	snapv1 "github.com/kubernetes-incubator/external-storage/snapshot/pkg/apis/crd/v1"
	snapshotVolume "github.com/kubernetes-incubator/external-storage/snapshot/pkg/volume"
	storkvolume "github.com/libopenstorage/stork/drivers/volume"
	"github.com/libopenstorage/stork/pkg/errors"
	"github.com/portworx/sched-ops/k8s"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	k8shelper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
)

const (
	// driverName is the name of the gcp driver implementation
	driverName = "gce"
	// provisioner names for portworx volumes
	provisionerName = "kubernetes.io/gce-pd"
	// pvcProvisionerAnnotation is the annotation on PVC which has the
	// provisioner name
	pvcProvisionerAnnotation = "volume.beta.kubernetes.io/storage-provisioner"
	// pvProvisionedByAnnotation is the annotation on PV which has the
	// provisioner name
	pvProvisionedByAnnotation = "pv.kubernetes.io/provisioned-by"
)

type gcp struct {
	storkvolume.ClusterPairNotSupported
	storkvolume.MigrationNotSupported
	storkvolume.GroupSnapshotNotSupported
	storkvolume.ClusterDomainsNotSupported
	storkvolume.BackupRestoreNotSupported
	storkvolume.CloneNotSupported
	storkvolume.SnapshotRestoreNotSupported
}

func (g *gcp) Init(_ interface{}) error {
	return nil
}

func (g *gcp) String() string {
	return driverName
}

func (g *gcp) Stop() error {
	return nil
}

func (g *gcp) OwnsPVC(pvc *v1.PersistentVolumeClaim) bool {

	provisioner := ""
	// Check for the provisioner in the PVC annotation. If not populated
	// try getting the provisioner from the Storage class.
	if val, ok := pvc.Annotations[pvcProvisionerAnnotation]; ok {
		provisioner = val
	} else {
		storageClassName := k8shelper.GetPersistentVolumeClaimClass(pvc)
		if storageClassName != "" {
			storageClass, err := k8s.Instance().GetStorageClass(storageClassName)
			if err == nil {
				provisioner = storageClass.Provisioner
			} else {
				logrus.Warnf("Error getting storageclass %v for pvc %v: %v", storageClassName, pvc.Name, err)
			}
		}
	}

	if provisioner == "" {
		// Try to get info from the PV since storage class could be deleted
		pv, err := k8s.Instance().GetPersistentVolume(pvc.Spec.VolumeName)
		if err != nil {
			logrus.Warnf("Error getting pv %v for pvc %v: %v", pvc.Spec.VolumeName, pvc.Name, err)
			return false
		}
		// Check the annotation in the PV for the provisioner
		if val, ok := pv.Annotations[pvProvisionedByAnnotation]; ok {
			provisioner = val
		} else {
			// Finally check the volume reference in the spec
			if pv.Spec.GCEPersistentDisk != nil {
				return true
			}
		}
	}

	if provisioner != provisionerName &&
		!isCsiProvisioner(provisioner) {
		logrus.Debugf("Provisioner in Storageclass not GCE: %v", provisioner)
		return false
	}
	return true
}

func isCsiProvisioner(provisioner string) bool {
	return false
}

func (g *gcp) StartBackup(backup *stork_crd.ApplicationBackup,
	pvcs []*v1.PersistentVolumeClaim,
) ([]*stork_crd.ApplicationBackupVolumeInfo, error) {
	for _, pvc := range pvcs {
		if pvc.DeletionTimestamp != nil {
			log.ApplicationBackupLog(backup).Warnf("Ignoring PVC %v which is being deleted", pvc.Name)
			continue
		}
		volumeInfo := &stork_crd.ApplicationBackupVolumeInfo{}
		volumeInfo.PersistentVolumeClaim = pvc.Name
		volumeInfo.Namespace = pvc.Namespace
		volumeInfo.DriverName = driverName
		volumeInfos = append(volumeInfos, volumeInfo)

		volume, err := k8s.Instance().GetVolumeForPersistentVolumeClaim(pvc)
		if err != nil {
			return nil, fmt.Errorf("Error getting volume for PVC: %v", err)
		}
		volumeInfo.Volume = volume
		taskID := p.getBackupRestoreTaskID(backup.UID, volumeInfo.Namespace, volumeInfo.PersistentVolumeClaim)
		credID := p.getCredID(backup.Spec.BackupLocation, backup.Namespace)
		request := &api.CloudBackupCreateRequest{
			VolumeID:       volume,
			CredentialUUID: credID,
			Name:           taskID,
		}
		request.Labels = make(map[string]string)
		request.Labels[cloudBackupOwnerLabel] = "stork"
		p.addApplicationBackupCloudsnapInfo(request, backup)

		_, err = volDriver.CloudBackupCreate(request)
		if err != nil {
			if _, ok := err.(*ost_errors.ErrExists); !ok {
				return nil, err
			}
		}
	}
}

func (g *gcp) InspectVolume(volumeID string) (*storkvolume.Info, error) {
	return nil, &errors.ErrNotSupported{}
}

func (g *gcp) GetClusterID() (string, error) {
	return "", &errors.ErrNotSupported{}
}

func (g *gcp) GetNodes() ([]*storkvolume.NodeInfo, error) {
	return nil, &errors.ErrNotSupported{}
}

func (g *gcp) GetPodVolumes(podSpec *v1.PodSpec, namespace string) ([]*storkvolume.Info, error) {
	return nil, &errors.ErrNotSupported{}
}

func (g *gcp) GetSnapshotPlugin() snapshotVolume.Plugin {
	return nil
}

func (g *gcp) GetSnapshotType(snap *snapv1.VolumeSnapshot) (string, error) {
	return "", &errors.ErrNotSupported{}
}

func (g *gcp) GetVolumeClaimTemplates([]v1.PersistentVolumeClaim) (
	[]v1.PersistentVolumeClaim, error) {
	return nil, &errors.ErrNotSupported{}
}

func init() {
	if err := storkvolume.Register(driverName, &gcp{}); err != nil {
		logrus.Panicf("Error registering gcp volume driver: %v", err)
	}
}
