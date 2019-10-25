package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/compute/metadata"
	snapv1 "github.com/kubernetes-incubator/external-storage/snapshot/pkg/apis/crd/v1"
	snapshotVolume "github.com/kubernetes-incubator/external-storage/snapshot/pkg/volume"
	storkvolume "github.com/libopenstorage/stork/drivers/volume"
	storkapi "github.com/libopenstorage/stork/pkg/apis/stork/v1alpha1"
	"github.com/libopenstorage/stork/pkg/errors"
	"github.com/libopenstorage/stork/pkg/log"
	"github.com/portworx/sched-ops/k8s"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
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
	pvNamePrefix              = "pvc-"
)

type gcp struct {
	projectID string
	zone      string
	service   *compute.Service
	storkvolume.ClusterPairNotSupported
	storkvolume.MigrationNotSupported
	storkvolume.GroupSnapshotNotSupported
	storkvolume.ClusterDomainsNotSupported
	storkvolume.CloneNotSupported
	storkvolume.SnapshotRestoreNotSupported
}

func (g *gcp) Init(_ interface{}) error {
	var err error
	g.zone, err = metadata.Zone()
	if err != nil {
		return fmt.Errorf("error getting zone for gce: %v", err)
	}

	g.projectID, err = metadata.ProjectID()
	if err != nil {
		return fmt.Errorf("error getting projectID for gce: %v", err)
	}

	creds, err := google.FindDefaultCredentials(oauth2.NoContext, compute.ComputeScope)
	if err != nil {
		return err
	}

	//	client := oauth2.NewClient(oauth2.NoContext, creds.TokenSource)

	g.service, err = compute.NewService(context.TODO(), option.WithTokenSource(creds.TokenSource))
	if err != nil {
		return err
	}

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

func (g *gcp) StartBackup(backup *storkapi.ApplicationBackup,
	pvcs []*v1.PersistentVolumeClaim,
) ([]*storkapi.ApplicationBackupVolumeInfo, error) {
	volumeInfos := make([]*storkapi.ApplicationBackupVolumeInfo, 0)

	for _, pvc := range pvcs {
		if pvc.DeletionTimestamp != nil {
			log.ApplicationBackupLog(backup).Warnf("Ignoring PVC %v which is being deleted", pvc.Name)
			continue
		}
		volumeInfo := &storkapi.ApplicationBackupVolumeInfo{}
		volumeInfo.PersistentVolumeClaim = pvc.Name
		volumeInfo.Namespace = pvc.Namespace
		volumeInfo.DriverName = driverName
		volumeInfos = append(volumeInfos, volumeInfo)

		pvName, err := k8s.Instance().GetVolumeForPersistentVolumeClaim(pvc)
		if err != nil {
			return nil, fmt.Errorf("Error getting PV name for PVC (%v/%v): %v", pvc.Namespace, pvc.Name, err)
		}
		pv, err := k8s.Instance().GetPersistentVolume(pvName)
		if err != nil {
			return nil, fmt.Errorf("Error getting pv %v: %v", pvName, err)
		}
		volume := pv.Spec.GCEPersistentDisk.PDName
		volumeInfo.Volume = volume
		snapshot := &compute.Snapshot{
			Name: "stork-snapshot-" + string(uuid.NewUUID()),
			Labels: map[string]string{
				"created-by":           "stork",
				"backup-uid":           string(backup.UID),
				"source-pvc-name":      pvc.Name,
				"source-pvc-namespace": pvc.Namespace,
			},
		}
		snapshotCall := g.service.Disks.CreateSnapshot(g.projectID, g.zone, volume, snapshot)

		// Set the Request ID to the unique name to get idempotency
		//snapshotCall = snapshotCall.RequestId(snapshot.Name)

		_, err = snapshotCall.Do()
		if err != nil {
			return nil, fmt.Errorf("Error triggering backup for volume: %v (PVC: %v, Namespace: %v): %v", volume, pvc.Name, pvc.Namespace, err)
		}
		volumeInfo.BackupID = snapshot.Name
	}
	return volumeInfos, nil
}

func (g *gcp) GetBackupStatus(backup *storkapi.ApplicationBackup) ([]*storkapi.ApplicationBackupVolumeInfo, error) {
	for _, vInfo := range backup.Status.Volumes {
		snapshot, err := g.service.Snapshots.Get(g.projectID, vInfo.BackupID).Do()
		/*
			filter := fmt.Sprintf("(labels.created-by=\"stork\") AND (labels.backup-uid=\"%v\") AND (labels.source-pvc-name=\"%v\") AND (labels.source-pvc-namespace=\"%v\")",
				backup.UID, vInfo.PersistentVolumeClaim, vInfo.Namespace)
			snapshots, err := g.service.Snapshots.List(g.projectID).Filter(filter).Do()
		*/
		if err != nil {
			return nil, err
		}
		//snapshot := snapshots.Items[0]
		switch snapshot.Status {
		case "CREATING", "UPLOADING":
			vInfo.Status = storkapi.ApplicationBackupStatusInProgress
			vInfo.Reason = fmt.Sprintf("Volume backup in progress: %v", snapshot.Status)
		case "DELETING", "FAILED":
			vInfo.Status = storkapi.ApplicationBackupStatusFailed
			vInfo.Reason = fmt.Sprintf("Backup failed for volume: %v", snapshot.Status)
		case "READY":
			vInfo.BackupID = snapshot.SelfLink
			vInfo.Status = storkapi.ApplicationBackupStatusSuccessful
			vInfo.Reason = fmt.Sprintf("Backup successful for volume")
		}
	}

	return backup.Status.Volumes, nil

}

// CancelBackup returns ErrNotSupported
func (g *gcp) CancelBackup(*storkapi.ApplicationBackup) error {
	//return &errors.ErrNotSupported{}
	return nil
}

// DeleteBackup returns ErrNotSupported
func (g *gcp) DeleteBackup(*storkapi.ApplicationBackup) error {
	return nil
	//return &errors.ErrNotSupported{}
}

func (g *gcp) UpdateMigratedPersistentVolumeSpec(
	object runtime.Unstructured,
) (runtime.Unstructured, error) {
	metadata, err := meta.Accessor(object)
	if err != nil {
		return nil, err
	}

	// Get access to the csi section of the PV
	_, found, err := unstructured.NestedString(object.UnstructuredContent(), "spec", "csi", "driver")
	if err != nil {
		return nil, err
	}

	// Determine if CSI is used
	if found {
		if err := unstructured.SetNestedField(object.UnstructuredContent(), metadata.GetName(), "spec", "csi", "volumeHandle"); err != nil {
			return nil, err
		}
		return object, nil
	}

	// Fallback to in-tree driver in case CSI isn't found
	err = unstructured.SetNestedField(object.UnstructuredContent(), metadata.GetName(), "spec", "gcePersistentDisk", "pdName")
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (g *gcp) generatePVName() string {
	return pvNamePrefix + string(uuid.NewUUID())
}

// StartRestore returns ErrNotSupported
func (g *gcp) StartRestore(
	restore *storkapi.ApplicationRestore,
	volumeBackupInfos []*storkapi.ApplicationBackupVolumeInfo,
) ([]*storkapi.ApplicationRestoreVolumeInfo, error) {

	volumeInfos := make([]*storkapi.ApplicationRestoreVolumeInfo, 0)
	for _, backupVolumeInfo := range volumeBackupInfos {
		volumeInfo := &storkapi.ApplicationRestoreVolumeInfo{}
		volumeInfo.PersistentVolumeClaim = backupVolumeInfo.PersistentVolumeClaim
		volumeInfo.SourceNamespace = backupVolumeInfo.Namespace
		volumeInfo.SourceVolume = backupVolumeInfo.Volume
		volumeInfo.RestoreVolume = g.generatePVName()
		volumeInfos = append(volumeInfos, volumeInfo)
		disk := &compute.Disk{

			Name:           volumeInfo.RestoreVolume,
			SourceSnapshot: backupVolumeInfo.BackupID,
			Labels: map[string]string{
				"created-by":           "stork",
				"restore-uid":          string(restore.UID),
				"source-pvc-name":      volumeInfo.PersistentVolumeClaim,
				"source-pvc-namespace": volumeInfo.SourceNamespace,
			},
		}
		_, err := g.service.Disks.Insert(g.projectID, g.zone, disk).Do()
		if err != nil {
			return nil, err
		}
	}
	return volumeInfos, nil
}

func (g *gcp) CancelRestore(*storkapi.ApplicationRestore) error {
	//return &errors.ErrNotSupported{}
	return nil
}

func (g *gcp) GetRestoreStatus(restore *storkapi.ApplicationRestore) ([]*storkapi.ApplicationRestoreVolumeInfo, error) {
	for _, vInfo := range restore.Status.Volumes {
		disk, err := g.service.Disks.Get(g.projectID, g.zone, vInfo.RestoreVolume).Do()
		if err != nil {
			return nil, err
		}
		switch disk.Status {
		case "CREATING", "RESTORING":
			vInfo.Status = storkapi.ApplicationRestoreStatusInProgress
			vInfo.Reason = fmt.Sprintf("Volume restore in progress: %v", disk.Status)
		case "DELETING", "FAILED":
			vInfo.Status = storkapi.ApplicationRestoreStatusFailed
			vInfo.Reason = fmt.Sprintf("Restore failed for volume: %v", disk.Status)
		case "READY":
			vInfo.Status = storkapi.ApplicationRestoreStatusSuccessful
			vInfo.Reason = fmt.Sprintf("Restore successful for volume")
		}
	}

	return restore.Status.Volumes, nil
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
	g := &gcp{}
	err := g.Init(nil)
	if err != nil {
		logrus.Errorf("Error init'ing gcp driver")
		return
	}
	if err := storkvolume.Register(driverName, g); err != nil {
		logrus.Panicf("Error registering gcp volume driver: %v", err)
	}
}
