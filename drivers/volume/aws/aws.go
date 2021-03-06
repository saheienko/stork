package aws

import (
	"fmt"
	"regexp"

	aws_sdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	snapv1 "github.com/kubernetes-incubator/external-storage/snapshot/pkg/apis/crd/v1"
	snapshotVolume "github.com/kubernetes-incubator/external-storage/snapshot/pkg/volume"
	"github.com/kubernetes-sigs/aws-ebs-csi-driver/pkg/cloud"
	storkvolume "github.com/libopenstorage/stork/drivers/volume"
	storkapi "github.com/libopenstorage/stork/pkg/apis/stork/v1alpha1"
	"github.com/libopenstorage/stork/pkg/errors"
	"github.com/libopenstorage/stork/pkg/log"
	"github.com/portworx/sched-ops/k8s/core"
	"github.com/portworx/sched-ops/k8s/storage"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	k8shelper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
)

const (
	// driverName is the name of the aws driver implementation
	driverName = "aws"
	// provisioner names for ebs volumes
	provisionerName = "kubernetes.io/aws-ebs"
	// pvcProvisionerAnnotation is the annotation on PVC which has the
	// provisioner name
	pvcProvisionerAnnotation = "volume.beta.kubernetes.io/storage-provisioner"
	// pvProvisionedByAnnotation is the annotation on PV which has the
	// provisioner name
	pvProvisionedByAnnotation = "pv.kubernetes.io/provisioned-by"
	pvNamePrefix              = "pvc-"

	// Tags used for snapshots and disks
	restoreUIDTag         = "restore-uid"
	nameTag               = "Name"
	createdByTag          = "created-by"
	backupUIDTag          = "backup-uid"
	sourcePVCNameTag      = "source-pvc-name"
	sourcePVCNamespaceTag = "source-pvc-namespace"
)

type aws struct {
	client *ec2.EC2
	storkvolume.ClusterPairNotSupported
	storkvolume.MigrationNotSupported
	storkvolume.GroupSnapshotNotSupported
	storkvolume.ClusterDomainsNotSupported
	storkvolume.CloneNotSupported
	storkvolume.SnapshotRestoreNotSupported
}

func (a *aws) Init(_ interface{}) error {

	s, err := session.NewSession(&aws_sdk.Config{})
	if err != nil {
		return err
	}
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(s),
			},
			&credentials.SharedCredentialsProvider{},
		})
	metadata, err := cloud.NewMetadata()
	if err != nil {
		return err
	}

	s, err = session.NewSession(&aws_sdk.Config{
		Region:      aws_sdk.String(metadata.GetRegion()),
		Credentials: creds,
	})
	if err != nil {
		return err
	}
	a.client = ec2.New(s)

	return nil
}

func (a *aws) String() string {
	return driverName
}

func (a *aws) Stop() error {
	return nil
}

func (a *aws) OwnsPVC(pvc *v1.PersistentVolumeClaim) bool {

	provisioner := ""
	// Check for the provisioner in the PVC annotation. If not populated
	// try getting the provisioner from the Storage class.
	if val, ok := pvc.Annotations[pvcProvisionerAnnotation]; ok {
		provisioner = val
	} else {
		storageClassName := k8shelper.GetPersistentVolumeClaimClass(pvc)
		if storageClassName != "" {
			storageClass, err := storage.Instance().GetStorageClass(storageClassName)
			if err == nil {
				provisioner = storageClass.Provisioner
			} else {
				logrus.Warnf("Error getting storageclass %v for pvc %v: %v", storageClassName, pvc.Name, err)
			}
		}
	}

	if provisioner == "" {
		// Try to get info from the PV since storage class could be deleted
		pv, err := core.Instance().GetPersistentVolume(pvc.Spec.VolumeName)
		if err != nil {
			logrus.Warnf("Error getting pv %v for pvc %v: %v", pvc.Spec.VolumeName, pvc.Name, err)
			return false
		}
		return a.OwnsPV(pv)
	}

	if provisioner != provisionerName &&
		!isCsiProvisioner(provisioner) {
		logrus.Debugf("Provisioner in Storageclass not AWS EBS: %v", provisioner)
		return false
	}
	return true
}

func (a *aws) OwnsPV(pv *v1.PersistentVolume) bool {
	var provisioner string
	// Check the annotation in the PV for the provisioner
	if val, ok := pv.Annotations[pvProvisionedByAnnotation]; ok {
		provisioner = val
	} else {
		// Finally check the volume reference in the spec
		if pv.Spec.AWSElasticBlockStore != nil {
			return true
		}
	}
	if provisioner != provisionerName &&
		!isCsiProvisioner(provisioner) {
		logrus.Debugf("Provisioner in Storageclass not AWS EBS: %v", provisioner)
		return false
	}
	return true
}

func isCsiProvisioner(provisioner string) bool {
	return false
}

func (a *aws) StartBackup(backup *storkapi.ApplicationBackup,
	pvcs []v1.PersistentVolumeClaim,
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

		pvName, err := core.Instance().GetVolumeForPersistentVolumeClaim(&pvc)
		if err != nil {
			return nil, fmt.Errorf("error getting PV name for PVC (%v/%v): %v", pvc.Namespace, pvc.Name, err)
		}
		pv, err := core.Instance().GetPersistentVolume(pvName)
		if err != nil {
			return nil, fmt.Errorf("error getting pv %v: %v", pvName, err)
		}
		volume := pvc.Spec.VolumeName
		ebsName := a.getEBSVolumeID(pv.Spec.AWSElasticBlockStore.VolumeID)

		ebsVolume, err := a.getEBSVolume(ebsName, nil)
		if err != nil {
			return nil, err
		}

		volumeInfo.Volume = volume
		volumeInfo.Zones = []string{*ebsVolume.AvailabilityZone}
		snapshotInput := &ec2.CreateSnapshotInput{
			VolumeId: aws_sdk.String(ebsName),
			Description: aws_sdk.String(fmt.Sprintf("Created by stork for %v for PVC %v Namespace %v Volume: %v",
				backup.Name, pvc.Name, pvc.Namespace, pv.Spec.AWSElasticBlockStore.VolumeID)),
			TagSpecifications: []*ec2.TagSpecification{
				{
					ResourceType: aws_sdk.String(ec2.ResourceTypeSnapshot),
					Tags: []*ec2.Tag{
						{
							Key:   aws_sdk.String(createdByTag),
							Value: aws_sdk.String("stork"),
						},
						{
							Key:   aws_sdk.String(backupUIDTag),
							Value: aws_sdk.String(string(backup.UID)),
						},
						{
							Key:   aws_sdk.String(sourcePVCNameTag),
							Value: aws_sdk.String(pvc.Name),
						},
						{
							Key:   aws_sdk.String(sourcePVCNamespaceTag),
							Value: aws_sdk.String(pvc.Namespace),
						},
						{
							Key:   aws_sdk.String(nameTag),
							Value: aws_sdk.String("stork-snapshot-" + volume),
						},
					},
				},
			},
		}
		sourceTags := make([]*ec2.Tag, 0)
		for _, tag := range ebsVolume.Tags {
			if *tag.Key == nameTag ||
				*tag.Key == createdByTag ||
				*tag.Key == backupUIDTag ||
				*tag.Key == sourcePVCNameTag ||
				*tag.Key == sourcePVCNamespaceTag {
				continue
			}
			sourceTags = append(sourceTags, tag)
		}
		snapshotInput.TagSpecifications[0].Tags = append(snapshotInput.TagSpecifications[0].Tags, sourceTags...)
		snapshot, err := a.client.CreateSnapshot(snapshotInput)
		if err != nil {
			return nil, err
		}

		volumeInfo.BackupID = *snapshot.SnapshotId
	}
	return volumeInfos, nil
}

func (a *aws) getEBSVolumeID(volumeID string) string {
	return regexp.MustCompile("vol-.*").FindString(volumeID)
}

func (a *aws) getEBSVolume(volumeID string, filters map[string]string) (*ec2.Volume, error) {
	input := &ec2.DescribeVolumesInput{}
	if volumeID != "" {
		input.VolumeIds = []*string{&volumeID}
	}
	if len(filters) > 0 {
		input.Filters = make([]*ec2.Filter, 0)
		for k, v := range filters {
			input.Filters = append(input.Filters, &ec2.Filter{
				Name:   aws_sdk.String(k),
				Values: []*string{aws_sdk.String(v)},
			})
		}
	}

	output, err := a.client.DescribeVolumes(input)
	if err != nil {
		return nil, err
	}

	if len(output.Volumes) != 1 {
		return nil, fmt.Errorf("received %v volumes for %v", len(output.Volumes), volumeID)
	}
	return output.Volumes[0], nil
}

func (a *aws) getEBSSnapshot(snapshotID string) (*ec2.Snapshot, error) {
	input := &ec2.DescribeSnapshotsInput{
		SnapshotIds: []*string{&snapshotID},
	}

	output, err := a.client.DescribeSnapshots(input)
	if err != nil {
		return nil, err
	}

	if len(output.Snapshots) != 1 {
		return nil, fmt.Errorf("received %v snapshots for %v", len(output.Snapshots), snapshotID)
	}
	return output.Snapshots[0], nil
}

func (a *aws) GetBackupStatus(backup *storkapi.ApplicationBackup) ([]*storkapi.ApplicationBackupVolumeInfo, error) {
	volumeInfos := make([]*storkapi.ApplicationBackupVolumeInfo, 0)

	for _, vInfo := range backup.Status.Volumes {
		if vInfo.DriverName != driverName {
			continue
		}
		snapshot, err := a.getEBSSnapshot(vInfo.BackupID)
		if err != nil {
			return nil, err
		}
		switch *snapshot.State {
		case "pending":
			vInfo.Status = storkapi.ApplicationBackupStatusInProgress
			vInfo.Reason = fmt.Sprintf("Volume backup in progress: %v (%v)", *snapshot.State, *snapshot.Progress)
		case "error":
			vInfo.Status = storkapi.ApplicationBackupStatusFailed
			vInfo.Reason = fmt.Sprintf("Backup failed for volume: %v", *snapshot.State)
		case "completed":
			vInfo.Status = storkapi.ApplicationBackupStatusSuccessful
			vInfo.Reason = "Backup successful for volume"
		}
		volumeInfos = append(volumeInfos, vInfo)
	}
	return volumeInfos, nil

}

func (a *aws) CancelBackup(backup *storkapi.ApplicationBackup) error {
	return a.DeleteBackup(backup)
}

func (a *aws) DeleteBackup(backup *storkapi.ApplicationBackup) error {
	for _, vInfo := range backup.Status.Volumes {
		if vInfo.DriverName != driverName {
			continue
		}
		input := &ec2.DeleteSnapshotInput{
			SnapshotId: aws_sdk.String(vInfo.BackupID),
		}

		_, err := a.client.DeleteSnapshot(input)
		if err != nil {
			// Do nothing if snapshot isn't found
			if awsErr, ok := err.(awserr.Error); ok {
				if awsErr.Code() == "InvalidSnapshot.NotFound" {
					continue
				}
			}
			return err
		}
	}
	return nil
}

func (a *aws) UpdateMigratedPersistentVolumeSpec(
	pv *v1.PersistentVolume,
) (*v1.PersistentVolume, error) {
	if pv.Spec.CSI != nil {
		pv.Spec.CSI.VolumeHandle = pv.Name
		return pv, nil
	}

	pv.Spec.AWSElasticBlockStore.VolumeID = pv.Name
	return pv, nil
}

func (a *aws) generatePVName() string {
	return pvNamePrefix + string(uuid.NewUUID())
}

func (a *aws) StartRestore(
	restore *storkapi.ApplicationRestore,
	volumeBackupInfos []*storkapi.ApplicationBackupVolumeInfo,
) ([]*storkapi.ApplicationRestoreVolumeInfo, error) {

	volumeInfos := make([]*storkapi.ApplicationRestoreVolumeInfo, 0)
	for _, backupVolumeInfo := range volumeBackupInfos {
		volumeInfo := &storkapi.ApplicationRestoreVolumeInfo{}
		volumeInfo.PersistentVolumeClaim = backupVolumeInfo.PersistentVolumeClaim
		volumeInfo.SourceNamespace = backupVolumeInfo.Namespace
		volumeInfo.SourceVolume = backupVolumeInfo.Volume
		volumeInfo.DriverName = driverName
		volumeInfo.RestoreVolume = a.generatePVName()
		volumeInfos = append(volumeInfos, volumeInfo)

		if len(backupVolumeInfo.Zones) == 0 {
			return nil, fmt.Errorf("zone missing in backup for volume (%v) %v", backupVolumeInfo.Namespace, backupVolumeInfo.PersistentVolumeClaim)
		}
		ebsSnapshot, err := a.getEBSSnapshot(backupVolumeInfo.BackupID)
		if err != nil {
			return nil, err
		}
		input := &ec2.CreateVolumeInput{
			SnapshotId:       aws_sdk.String(backupVolumeInfo.BackupID),
			AvailabilityZone: aws_sdk.String(backupVolumeInfo.Zones[0]),
			TagSpecifications: []*ec2.TagSpecification{
				{
					ResourceType: aws_sdk.String(ec2.ResourceTypeVolume),
					Tags: []*ec2.Tag{

						{
							Key:   aws_sdk.String(createdByTag),
							Value: aws_sdk.String("stork"),
						},
						{
							Key:   aws_sdk.String(restoreUIDTag),
							Value: aws_sdk.String(string(restore.UID)),
						},
						{
							Key:   aws_sdk.String(sourcePVCNameTag),
							Value: aws_sdk.String(volumeInfo.PersistentVolumeClaim),
						},
						{
							Key:   aws_sdk.String(sourcePVCNamespaceTag),
							Value: aws_sdk.String(volumeInfo.SourceNamespace),
						},
						{
							Key:   aws_sdk.String(nameTag),
							Value: aws_sdk.String(volumeInfo.RestoreVolume),
						},
					},
				},
			},
		}

		sourceTags := make([]*ec2.Tag, 0)
		for _, tag := range ebsSnapshot.Tags {
			if *tag.Key == nameTag ||
				*tag.Key == createdByTag ||
				*tag.Key == restoreUIDTag ||
				*tag.Key == sourcePVCNameTag ||
				*tag.Key == sourcePVCNamespaceTag {
				continue
			}
			sourceTags = append(sourceTags, tag)
		}
		input.TagSpecifications[0].Tags = append(input.TagSpecifications[0].Tags, sourceTags...)
		output, err := a.client.CreateVolume(input)
		if err != nil {
			return nil, err
		}
		volumeInfo.RestoreVolume = *output.VolumeId
	}
	return volumeInfos, nil
}

func (a *aws) CancelRestore(*storkapi.ApplicationRestore) error {
	return nil
}

func (a *aws) GetRestoreStatus(restore *storkapi.ApplicationRestore) ([]*storkapi.ApplicationRestoreVolumeInfo, error) {
	volumeInfos := make([]*storkapi.ApplicationRestoreVolumeInfo, 0)
	for _, vInfo := range restore.Status.Volumes {
		if vInfo.DriverName != driverName {
			continue
		}
		ebsVolume, err := a.getEBSVolume(vInfo.RestoreVolume, nil)
		if err != nil {
			return nil, err
		}
		switch *ebsVolume.State {
		default:
			fallthrough
		case "error", "deleting", "deleted":
			vInfo.Status = storkapi.ApplicationRestoreStatusFailed
			vInfo.Reason = fmt.Sprintf("Restore failed for volume: %v", *ebsVolume.State)
		case "creating":
			vInfo.Status = storkapi.ApplicationRestoreStatusInProgress
			vInfo.Reason = fmt.Sprintf("Volume restore in progress: %v", *ebsVolume.State)
		case "available", "in-use":
			vInfo.Status = storkapi.ApplicationRestoreStatusSuccessful
			vInfo.Reason = "Restore successful for volume"
		}
		volumeInfos = append(volumeInfos, vInfo)
	}

	return volumeInfos, nil
}

func (a *aws) InspectVolume(volumeID string) (*storkvolume.Info, error) {
	return nil, &errors.ErrNotSupported{}
}

func (a *aws) GetClusterID() (string, error) {
	return "", &errors.ErrNotSupported{}
}

func (a *aws) GetNodes() ([]*storkvolume.NodeInfo, error) {
	return nil, &errors.ErrNotSupported{}
}

func (a *aws) GetPodVolumes(podSpec *v1.PodSpec, namespace string) ([]*storkvolume.Info, error) {
	return nil, &errors.ErrNotSupported{}
}

func (a *aws) GetSnapshotPlugin() snapshotVolume.Plugin {
	return nil
}

func (a *aws) GetSnapshotType(snap *snapv1.VolumeSnapshot) (string, error) {
	return "", &errors.ErrNotSupported{}
}

func (a *aws) GetVolumeClaimTemplates([]v1.PersistentVolumeClaim) (
	[]v1.PersistentVolumeClaim, error) {
	return nil, &errors.ErrNotSupported{}
}

func init() {
	a := &aws{}
	err := a.Init(nil)
	if err != nil {
		logrus.Debugf("Error init'ing aws driver: %v", err)
		return
	}
	if err := storkvolume.Register(driverName, a); err != nil {
		logrus.Panicf("Error registering aws volume driver: %v", err)
	}
}
