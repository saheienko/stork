package k8s

import (
	"time"

	snap_v1 "github.com/kubernetes-incubator/external-storage/snapshot/pkg/apis/crd/v1"
	"github.com/libopenstorage/stork/pkg/apis/stork/v1alpha1"
	"github.com/portworx/sched-ops/k8s/stork"
)

// SnapshotScheduleOps is an interface to perform k8s VolumeSnapshotSchedule operations
type SnapshotScheduleOps interface {
	// GetSnapshotSchedule gets the SnapshotSchedule
	GetSnapshotSchedule(string, string) (*v1alpha1.VolumeSnapshotSchedule, error)
	// CreateSnapshotSchedule creates a SnapshotSchedule
	CreateSnapshotSchedule(*v1alpha1.VolumeSnapshotSchedule) (*v1alpha1.VolumeSnapshotSchedule, error)
	// UpdateSnapshotSchedule updates the SnapshotSchedule
	UpdateSnapshotSchedule(*v1alpha1.VolumeSnapshotSchedule) (*v1alpha1.VolumeSnapshotSchedule, error)
	// ListSnapshotSchedules lists all the SnapshotSchedules
	ListSnapshotSchedules(string) (*v1alpha1.VolumeSnapshotScheduleList, error)
	// DeleteSnapshotSchedule deletes the SnapshotSchedule
	DeleteSnapshotSchedule(string, string) error
	// ValidateSnapshotSchedule validates the given SnapshotSchedule. It checks the status of each of
	// the snapshots triggered for this schedule and returns a map of successfull snapshots. The key of the
	// map will be the schedule type and value will be list of snapshots for that schedule type.
	// The caller is expected to validate if the returned map has all snapshots expected at that point of time
	ValidateSnapshotSchedule(string, string, time.Duration, time.Duration) (
		map[v1alpha1.SchedulePolicyType][]*v1alpha1.ScheduledVolumeSnapshotStatus, error)
}

// GroupSnapshotOps is an interface to perform k8s GroupVolumeSnapshot operations
type GroupSnapshotOps interface {
	// GetGroupSnapshot returns the group snapshot for the given name and namespace
	GetGroupSnapshot(name, namespace string) (*v1alpha1.GroupVolumeSnapshot, error)
	// ListGroupSnapshots lists all group snapshots for the given namespace
	ListGroupSnapshots(namespace string) (*v1alpha1.GroupVolumeSnapshotList, error)
	// CreateGroupSnapshot creates the given group snapshot
	CreateGroupSnapshot(*v1alpha1.GroupVolumeSnapshot) (*v1alpha1.GroupVolumeSnapshot, error)
	// UpdateGroupSnapshot updates the given group snapshot
	UpdateGroupSnapshot(*v1alpha1.GroupVolumeSnapshot) (*v1alpha1.GroupVolumeSnapshot, error)
	// DeleteGroupSnapshot deletes the group snapshot with the given name and namespace
	DeleteGroupSnapshot(name, namespace string) error
	// ValidateGroupSnapshot checks if the group snapshot with given name and namespace is in ready state
	//  If retry is true, the validation will be retried with given timeout and retry internal
	ValidateGroupSnapshot(name, namespace string, retry bool, timeout, retryInterval time.Duration) error
	// GetSnapshotsForGroupSnapshot returns all child snapshots for the group snapshot
	GetSnapshotsForGroupSnapshot(name, namespace string) ([]*snap_v1.VolumeSnapshot, error)
}

// VolumeSnapshotRestoreOps is interface to perform isnapshot restore using CRD
type VolumeSnapshotRestoreOps interface {
	// CreateVolumeSnapshotRestore restore snapshot to pvc specifed in CRD, if no pvcs defined we restore to
	// parent volumes
	CreateVolumeSnapshotRestore(snap *v1alpha1.VolumeSnapshotRestore) (*v1alpha1.VolumeSnapshotRestore, error)
	// UpdateVolumeSnapshotRestore updates given volumesnapshorestore CRD
	UpdateVolumeSnapshotRestore(snap *v1alpha1.VolumeSnapshotRestore) (*v1alpha1.VolumeSnapshotRestore, error)
	// GetVolumeSnapshotRestore returns details of given restore crd status
	GetVolumeSnapshotRestore(name, namespace string) (*v1alpha1.VolumeSnapshotRestore, error)
	// ListVolumeSnapshotRestore return list of volumesnapshotrestores in given namespaces
	ListVolumeSnapshotRestore(namespace string) (*v1alpha1.VolumeSnapshotRestoreList, error)
	// DeleteVolumeSnapshotRestore delete given volumesnapshotrestore CRD
	DeleteVolumeSnapshotRestore(name, namespace string) error
	// ValidateVolumeSnapshotRestore validates given volumesnapshotrestore CRD
	ValidateVolumeSnapshotRestore(name, namespace string, timeout, retry time.Duration) error
}

// RuleOps is an interface to perform operations for k8s stork rule
type RuleOps interface {
	// GetRule fetches the given stork rule
	GetRule(name, namespace string) (*v1alpha1.Rule, error)
	// CreateRule creates the given stork rule
	CreateRule(rule *v1alpha1.Rule) (*v1alpha1.Rule, error)
	// DeleteRule deletes the given stork rule
	DeleteRule(name, namespace string) error
}

// ClusterPairOps is an interface to perfrom k8s ClusterPair operations
type ClusterPairOps interface {
	// CreateClusterPair creates the ClusterPair
	CreateClusterPair(*v1alpha1.ClusterPair) (*v1alpha1.ClusterPair, error)
	// GetClusterPair gets the ClusterPair
	GetClusterPair(string, string) (*v1alpha1.ClusterPair, error)
	// ListClusterPairs gets all the ClusterPairs
	ListClusterPairs(string) (*v1alpha1.ClusterPairList, error)
	// UpdateClusterPair updates the ClusterPair
	UpdateClusterPair(*v1alpha1.ClusterPair) (*v1alpha1.ClusterPair, error)
	// DeleteClusterPair deletes the ClusterPair
	DeleteClusterPair(string, string) error
	// ValidateClusterPair validates clusterpair status
	ValidateClusterPair(string, string, time.Duration, time.Duration) error
}

// ClusterDomainsOps is an interface to perform k8s ClusterDomains operations
type ClusterDomainsOps interface {
	// CreateClusterDomainsStatus creates the ClusterDomainStatus
	CreateClusterDomainsStatus(*v1alpha1.ClusterDomainsStatus) (*v1alpha1.ClusterDomainsStatus, error)
	// GetClusterDomainsStatus gets the ClusterDomainsStatus
	GetClusterDomainsStatus(string) (*v1alpha1.ClusterDomainsStatus, error)
	// UpdateClusterDomainsStatus updates the ClusterDomainsStatus
	UpdateClusterDomainsStatus(*v1alpha1.ClusterDomainsStatus) (*v1alpha1.ClusterDomainsStatus, error)
	// DeleteClusterDomainsStatus deletes the ClusterDomainsStatus
	DeleteClusterDomainsStatus(string) error
	// ListClusterDomainStatuses lists ClusterDomainsStatus
	ListClusterDomainStatuses() (*v1alpha1.ClusterDomainsStatusList, error)
	// ValidateClusterDomainsStatus validates the ClusterDomainsStatus
	ValidateClusterDomainsStatus(string, map[string]bool, time.Duration, time.Duration) error
	// CreateClusterDomainUpdate creates the ClusterDomainUpdate
	CreateClusterDomainUpdate(*v1alpha1.ClusterDomainUpdate) (*v1alpha1.ClusterDomainUpdate, error)
	// GetClusterDomainUpdate gets the ClusterDomainUpdate
	GetClusterDomainUpdate(string) (*v1alpha1.ClusterDomainUpdate, error)
	// UpdateClusterDomainUpdate updates the ClusterDomainUpdate
	UpdateClusterDomainUpdate(*v1alpha1.ClusterDomainUpdate) (*v1alpha1.ClusterDomainUpdate, error)
	// DeleteClusterDomainUpdate deletes the ClusterDomainUpdate
	DeleteClusterDomainUpdate(string) error
	// ValidateClusterDomainUpdate validates ClusterDomainUpdate
	ValidateClusterDomainUpdate(string, time.Duration, time.Duration) error
	// ListClusterDomainUpdates lists ClusterDomainUpdates
	ListClusterDomainUpdates() (*v1alpha1.ClusterDomainUpdateList, error)
}

// MigrationOps is an interface to perfrom k8s Migration operations
type MigrationOps interface {
	// CreateMigration creates the Migration
	CreateMigration(*v1alpha1.Migration) (*v1alpha1.Migration, error)
	// GetMigration gets the Migration
	GetMigration(string, string) (*v1alpha1.Migration, error)
	// ListMigrations lists all the Migrations
	ListMigrations(string) (*v1alpha1.MigrationList, error)
	// UpdateMigration updates the Migration
	UpdateMigration(*v1alpha1.Migration) (*v1alpha1.Migration, error)
	// DeleteMigration deletes the Migration
	DeleteMigration(string, string) error
	// ValidateMigration validate the Migration status
	ValidateMigration(string, string, time.Duration, time.Duration) error
	// GetMigrationSchedule gets the MigrationSchedule
	GetMigrationSchedule(string, string) (*v1alpha1.MigrationSchedule, error)
	// CreateMigrationSchedule creates a MigrationSchedule
	CreateMigrationSchedule(*v1alpha1.MigrationSchedule) (*v1alpha1.MigrationSchedule, error)
	// UpdateMigrationSchedule updates the MigrationSchedule
	UpdateMigrationSchedule(*v1alpha1.MigrationSchedule) (*v1alpha1.MigrationSchedule, error)
	// ListMigrationSchedules lists all the MigrationSchedules
	ListMigrationSchedules(string) (*v1alpha1.MigrationScheduleList, error)
	// DeleteMigrationSchedule deletes the MigrationSchedule
	DeleteMigrationSchedule(string, string) error
	// ValidateMigrationSchedule validates the given MigrationSchedule. It checks the status of each of
	// the migrations triggered for this schedule and returns a map of successfull migrations. The key of the
	// map will be the schedule type and value will be list of migrations for that schedule type.
	// The caller is expected to validate if the returned map has all migrations expected at that point of time
	ValidateMigrationSchedule(string, string, time.Duration, time.Duration) (
		map[v1alpha1.SchedulePolicyType][]*v1alpha1.ScheduledMigrationStatus, error)
}

// SchedulePolicyOps is an interface to manage SchedulePolicy Object
type SchedulePolicyOps interface {
	// CreateSchedulePolicy creates a SchedulePolicy
	CreateSchedulePolicy(*v1alpha1.SchedulePolicy) (*v1alpha1.SchedulePolicy, error)
	// GetSchedulePolicy gets the SchedulePolicy
	GetSchedulePolicy(string) (*v1alpha1.SchedulePolicy, error)
	// ListSchedulePolicies lists all the SchedulePolicies
	ListSchedulePolicies() (*v1alpha1.SchedulePolicyList, error)
	// UpdateSchedulePolicy updates the SchedulePolicy
	UpdateSchedulePolicy(*v1alpha1.SchedulePolicy) (*v1alpha1.SchedulePolicy, error)
	// DeleteSchedulePolicy deletes the SchedulePolicy
	DeleteSchedulePolicy(string) error
}

// BackupLocationOps is an interface to perfrom k8s BackupLocation operations
type BackupLocationOps interface {
	// CreateBackupLocation creates the BackupLocation
	CreateBackupLocation(*v1alpha1.BackupLocation) (*v1alpha1.BackupLocation, error)
	// GetBackupLocation gets the BackupLocation
	GetBackupLocation(string, string) (*v1alpha1.BackupLocation, error)
	// ListBackupLocations lists all the BackupLocations
	ListBackupLocations(string) (*v1alpha1.BackupLocationList, error)
	// UpdateBackupLocation updates the BackupLocation
	UpdateBackupLocation(*v1alpha1.BackupLocation) (*v1alpha1.BackupLocation, error)
	// DeleteBackupLocation deletes the BackupLocation
	DeleteBackupLocation(string, string) error
	// ValidateBackupLocation validates the BackupLocation
	ValidateBackupLocation(string, string, time.Duration, time.Duration) error
}

// ApplicationBackupRestoreOps is an interface to perfrom k8s Application Backup
// and Restore operations
type ApplicationBackupRestoreOps interface {
	// CreateApplicationBackup creates the ApplicationBackup
	CreateApplicationBackup(*v1alpha1.ApplicationBackup) (*v1alpha1.ApplicationBackup, error)
	// GetApplicationBackup gets the ApplicationBackup
	GetApplicationBackup(string, string) (*v1alpha1.ApplicationBackup, error)
	// ListApplicationBackups lists all the ApplicationBackups
	ListApplicationBackups(string) (*v1alpha1.ApplicationBackupList, error)
	// UpdateApplicationBackup updates the ApplicationBackup
	UpdateApplicationBackup(*v1alpha1.ApplicationBackup) (*v1alpha1.ApplicationBackup, error)
	// DeleteApplicationBackup deletes the ApplicationBackup
	DeleteApplicationBackup(string, string) error
	// ValidateApplicationBackup validates the ApplicationBackup
	ValidateApplicationBackup(string, string, time.Duration, time.Duration) error
	// CreateApplicationRestore creates the ApplicationRestore
	CreateApplicationRestore(*v1alpha1.ApplicationRestore) (*v1alpha1.ApplicationRestore, error)
	// GetApplicationRestore gets the ApplicationRestore
	GetApplicationRestore(string, string) (*v1alpha1.ApplicationRestore, error)
	// ListApplicationRestores lists all the ApplicationRestores
	ListApplicationRestores(string) (*v1alpha1.ApplicationRestoreList, error)
	// UpdateApplicationRestore updates the ApplicationRestore
	UpdateApplicationRestore(*v1alpha1.ApplicationRestore) (*v1alpha1.ApplicationRestore, error)
	// DeleteApplicationRestore deletes the ApplicationRestore
	DeleteApplicationRestore(string, string) error
	// ValidateApplicationRestore validates the ApplicationRestore
	ValidateApplicationRestore(string, string, time.Duration, time.Duration) error
	// GetApplicationBackupSchedule gets the ApplicationBackupSchedule
	GetApplicationBackupSchedule(string, string) (*v1alpha1.ApplicationBackupSchedule, error)
	// CreateApplicationBackupSchedule creates an ApplicationBackupSchedule
	CreateApplicationBackupSchedule(*v1alpha1.ApplicationBackupSchedule) (*v1alpha1.ApplicationBackupSchedule, error)
	// UpdateApplicationBackupSchedule updates the ApplicationBackupSchedule
	UpdateApplicationBackupSchedule(*v1alpha1.ApplicationBackupSchedule) (*v1alpha1.ApplicationBackupSchedule, error)
	// ListApplicationBackupSchedules lists all the ApplicationBackupSchedules
	ListApplicationBackupSchedules(string) (*v1alpha1.ApplicationBackupScheduleList, error)
	// DeleteApplicationBackupSchedule deletes the ApplicationBackupSchedule
	DeleteApplicationBackupSchedule(string, string) error
	// ValidateApplicationBackupSchedule validates the given ApplicationBackupSchedule. It checks the status of each of
	// the backups triggered for this schedule and returns a map of successfull backups. The key of the
	// map will be the schedule type and value will be list of backups for that schedule type.
	// The caller is expected to validate if the returned map has all backups expected at that point of time
	ValidateApplicationBackupSchedule(string, string, time.Duration, time.Duration) (
		map[v1alpha1.SchedulePolicyType][]*v1alpha1.ScheduledApplicationBackupStatus, error)
}

// ApplicationCloneOps is an interface to perform k8s Application Clone operations
type ApplicationCloneOps interface {
	// CreateApplicationClone creates the ApplicationClone
	CreateApplicationClone(*v1alpha1.ApplicationClone) (*v1alpha1.ApplicationClone, error)
	// GetApplicationClone gets the ApplicationClone
	GetApplicationClone(string, string) (*v1alpha1.ApplicationClone, error)
	// ListApplicationClones lists all the ApplicationClones
	ListApplicationClones(string) (*v1alpha1.ApplicationCloneList, error)
	// UpdateApplicationClone updates the ApplicationClone
	UpdateApplicationClone(*v1alpha1.ApplicationClone) (*v1alpha1.ApplicationClone, error)
	// DeleteApplicationClone deletes the ApplicationClone
	DeleteApplicationClone(string, string) error
	// ValidateApplicationClone validates the ApplicationClone
	ValidateApplicationClone(string, string, time.Duration, time.Duration) error
}

// VolumeSnapshotSchedule APIs - BEGIN

func (k *k8sOps) GetSnapshotSchedule(name string, namespace string) (*v1alpha1.VolumeSnapshotSchedule, error) {
	return stork.Instance().GetSnapshotSchedule(name, namespace)
}

func (k *k8sOps) ListSnapshotSchedules(namespace string) (*v1alpha1.VolumeSnapshotScheduleList, error) {
	return stork.Instance().ListSnapshotSchedules(namespace)
}

func (k *k8sOps) CreateSnapshotSchedule(snapshotSchedule *v1alpha1.VolumeSnapshotSchedule) (*v1alpha1.VolumeSnapshotSchedule, error) {
	return stork.Instance().CreateSnapshotSchedule(snapshotSchedule)
}

func (k *k8sOps) UpdateSnapshotSchedule(snapshotSchedule *v1alpha1.VolumeSnapshotSchedule) (*v1alpha1.VolumeSnapshotSchedule, error) {
	return stork.Instance().UpdateSnapshotSchedule(snapshotSchedule)
}
func (k *k8sOps) DeleteSnapshotSchedule(name string, namespace string) error {
	return stork.Instance().DeleteSnapshotSchedule(name, namespace)
}

func (k *k8sOps) ValidateSnapshotSchedule(name string, namespace string, timeout, retryInterval time.Duration) (
	map[v1alpha1.SchedulePolicyType][]*v1alpha1.ScheduledVolumeSnapshotStatus, error) {
	return stork.Instance().ValidateSnapshotSchedule(name, namespace, timeout, retryInterval)
}

// VolumeSnapshotSchedule APIs - END

// GroupSnapshot APIs - BEGIN

func (k *k8sOps) GetGroupSnapshot(name, namespace string) (*v1alpha1.GroupVolumeSnapshot, error) {
	return stork.Instance().GetGroupSnapshot(name, namespace)
}

func (k *k8sOps) ListGroupSnapshots(namespace string) (*v1alpha1.GroupVolumeSnapshotList, error) {
	return stork.Instance().ListGroupSnapshots(namespace)
}

func (k *k8sOps) CreateGroupSnapshot(snap *v1alpha1.GroupVolumeSnapshot) (*v1alpha1.GroupVolumeSnapshot, error) {
	return stork.Instance().CreateGroupSnapshot(snap)
}

func (k *k8sOps) UpdateGroupSnapshot(snap *v1alpha1.GroupVolumeSnapshot) (*v1alpha1.GroupVolumeSnapshot, error) {
	return stork.Instance().UpdateGroupSnapshot(snap)
}

func (k *k8sOps) DeleteGroupSnapshot(name, namespace string) error {
	return stork.Instance().DeleteGroupSnapshot(name, namespace)
}

func (k *k8sOps) ValidateGroupSnapshot(name, namespace string, retry bool, timeout, retryInterval time.Duration) error {
	return stork.Instance().ValidateGroupSnapshot(name, namespace, retry, timeout, retryInterval)
}

func (k *k8sOps) GetSnapshotsForGroupSnapshot(name, namespace string) ([]*snap_v1.VolumeSnapshot, error) {
	return stork.Instance().GetSnapshotsForGroupSnapshot(name, namespace)
}

// GroupSnapshot APIs - END

// Restore Snapshot APIs - BEGIN

func (k *k8sOps) CreateVolumeSnapshotRestore(snapRestore *v1alpha1.VolumeSnapshotRestore) (*v1alpha1.VolumeSnapshotRestore, error) {
	return stork.Instance().CreateVolumeSnapshotRestore(snapRestore)
}

func (k *k8sOps) UpdateVolumeSnapshotRestore(snapRestore *v1alpha1.VolumeSnapshotRestore) (*v1alpha1.VolumeSnapshotRestore, error) {
	return stork.Instance().UpdateVolumeSnapshotRestore(snapRestore)
}

func (k *k8sOps) GetVolumeSnapshotRestore(name, namespace string) (*v1alpha1.VolumeSnapshotRestore, error) {
	return stork.Instance().GetVolumeSnapshotRestore(name, namespace)
}

func (k *k8sOps) ListVolumeSnapshotRestore(namespace string) (*v1alpha1.VolumeSnapshotRestoreList, error) {
	return stork.Instance().ListVolumeSnapshotRestore(namespace)
}

func (k *k8sOps) DeleteVolumeSnapshotRestore(name, namespace string) error {
	return stork.Instance().DeleteVolumeSnapshotRestore(name, namespace)
}

func (k *k8sOps) ValidateVolumeSnapshotRestore(name, namespace string, timeout, retryInterval time.Duration) error {
	return stork.Instance().ValidateVolumeSnapshotRestore(name, namespace, timeout, retryInterval)
}

// Restore Snapshot APIs - END

// Rule APIs - BEGIN

func (k *k8sOps) GetRule(name, namespace string) (*v1alpha1.Rule, error) {
	return stork.Instance().GetRule(name, namespace)
}

func (k *k8sOps) CreateRule(rule *v1alpha1.Rule) (*v1alpha1.Rule, error) {
	return stork.Instance().CreateRule(rule)
}

func (k *k8sOps) DeleteRule(name, namespace string) error {
	return stork.Instance().DeleteRule(name, namespace)
}

// Rule APIs - END

// ClusterPair APIs - BEGIN

func (k *k8sOps) GetClusterPair(name string, namespace string) (*v1alpha1.ClusterPair, error) {
	return stork.Instance().GetClusterPair(name, namespace)
}

func (k *k8sOps) ListClusterPairs(namespace string) (*v1alpha1.ClusterPairList, error) {
	return stork.Instance().ListClusterPairs(namespace)
}

func (k *k8sOps) CreateClusterPair(pair *v1alpha1.ClusterPair) (*v1alpha1.ClusterPair, error) {
	return stork.Instance().CreateClusterPair(pair)
}

func (k *k8sOps) UpdateClusterPair(pair *v1alpha1.ClusterPair) (*v1alpha1.ClusterPair, error) {
	return stork.Instance().UpdateClusterPair(pair)
}

func (k *k8sOps) DeleteClusterPair(name string, namespace string) error {
	return stork.Instance().DeleteClusterPair(name, namespace)
}

func (k *k8sOps) ValidateClusterPair(name string, namespace string, timeout, retryInterval time.Duration) error {
	return stork.Instance().ValidateClusterPair(name, namespace, timeout, retryInterval)
}

// ClusterPair APIs - END

// Migration APIs - BEGIN

func (k *k8sOps) GetMigration(name string, namespace string) (*v1alpha1.Migration, error) {
	return stork.Instance().GetMigration(name, namespace)
}

func (k *k8sOps) ListMigrations(namespace string) (*v1alpha1.MigrationList, error) {
	return stork.Instance().ListMigrations(namespace)
}

func (k *k8sOps) CreateMigration(migration *v1alpha1.Migration) (*v1alpha1.Migration, error) {
	return stork.Instance().CreateMigration(migration)
}

func (k *k8sOps) DeleteMigration(name string, namespace string) error {
	return stork.Instance().DeleteMigration(name, namespace)
}

func (k *k8sOps) UpdateMigration(migration *v1alpha1.Migration) (*v1alpha1.Migration, error) {
	return stork.Instance().UpdateMigration(migration)
}

func (k *k8sOps) ValidateMigration(name string, namespace string, timeout, retryInterval time.Duration) error {
	return stork.Instance().ValidateMigration(name, namespace, timeout, retryInterval)
}

func (k *k8sOps) GetMigrationSchedule(name string, namespace string) (*v1alpha1.MigrationSchedule, error) {
	return stork.Instance().GetMigrationSchedule(name, namespace)
}

func (k *k8sOps) ListMigrationSchedules(namespace string) (*v1alpha1.MigrationScheduleList, error) {
	return stork.Instance().ListMigrationSchedules(namespace)
}

func (k *k8sOps) CreateMigrationSchedule(migrationSchedule *v1alpha1.MigrationSchedule) (*v1alpha1.MigrationSchedule, error) {
	return stork.Instance().CreateMigrationSchedule(migrationSchedule)
}

func (k *k8sOps) UpdateMigrationSchedule(migrationSchedule *v1alpha1.MigrationSchedule) (*v1alpha1.MigrationSchedule, error) {
	return stork.Instance().UpdateMigrationSchedule(migrationSchedule)
}
func (k *k8sOps) DeleteMigrationSchedule(name string, namespace string) error {
	return stork.Instance().DeleteMigrationSchedule(name, namespace)
}

func (k *k8sOps) ValidateMigrationSchedule(name string, namespace string, timeout, retryInterval time.Duration) (
	map[v1alpha1.SchedulePolicyType][]*v1alpha1.ScheduledMigrationStatus, error) {
	return stork.Instance().ValidateMigrationSchedule(name, namespace, timeout, retryInterval)
}

// Migration APIs - END

// SchedulePolicy APIs - BEGIN

func (k *k8sOps) GetSchedulePolicy(name string) (*v1alpha1.SchedulePolicy, error) {
	return stork.Instance().GetSchedulePolicy(name)
}

func (k *k8sOps) ListSchedulePolicies() (*v1alpha1.SchedulePolicyList, error) {
	return stork.Instance().ListSchedulePolicies()
}

func (k *k8sOps) CreateSchedulePolicy(schedulePolicy *v1alpha1.SchedulePolicy) (*v1alpha1.SchedulePolicy, error) {
	return stork.Instance().CreateSchedulePolicy(schedulePolicy)
}

func (k *k8sOps) DeleteSchedulePolicy(name string) error {
	return stork.Instance().DeleteSchedulePolicy(name)
}

func (k *k8sOps) UpdateSchedulePolicy(schedulePolicy *v1alpha1.SchedulePolicy) (*v1alpha1.SchedulePolicy, error) {
	return stork.Instance().UpdateSchedulePolicy(schedulePolicy)
}

// SchedulePolicy APIs - END

// ClusterDomain CRD - BEGIN

// CreateClusterDomainsStatus creates the ClusterDomainStatus
func (k *k8sOps) CreateClusterDomainsStatus(clusterDomainsStatus *v1alpha1.ClusterDomainsStatus) (*v1alpha1.ClusterDomainsStatus, error) {
	return stork.Instance().CreateClusterDomainsStatus(clusterDomainsStatus)
}

// GetClusterDomainsStatus gets the ClusterDomainsStatus
func (k *k8sOps) GetClusterDomainsStatus(name string) (*v1alpha1.ClusterDomainsStatus, error) {
	return stork.Instance().GetClusterDomainsStatus(name)
}

// UpdateClusterDomainsStatus updates the ClusterDomainsStatus
func (k *k8sOps) UpdateClusterDomainsStatus(clusterDomainsStatus *v1alpha1.ClusterDomainsStatus) (*v1alpha1.ClusterDomainsStatus, error) {
	return stork.Instance().UpdateClusterDomainsStatus(clusterDomainsStatus)
}

// DeleteClusterDomainsStatus deletes the ClusterDomainsStatus
func (k *k8sOps) DeleteClusterDomainsStatus(name string) error {
	return stork.Instance().DeleteClusterDomainsStatus(name)
}

func (k *k8sOps) ValidateClusterDomainsStatus(name string, domainMap map[string]bool, timeout, retryInterval time.Duration) error {
	return stork.Instance().ValidateClusterDomainsStatus(name, domainMap, timeout, retryInterval)
}

// ListClusterDomainStatuses lists ClusterDomainsStatus
func (k *k8sOps) ListClusterDomainStatuses() (*v1alpha1.ClusterDomainsStatusList, error) {
	return stork.Instance().ListClusterDomainStatuses()
}

// CreateClusterDomainUpdate creates the ClusterDomainUpdate
func (k *k8sOps) CreateClusterDomainUpdate(clusterDomainUpdate *v1alpha1.ClusterDomainUpdate) (*v1alpha1.ClusterDomainUpdate, error) {
	return stork.Instance().CreateClusterDomainUpdate(clusterDomainUpdate)
}

// GetClusterDomainUpdate gets the ClusterDomainUpdate
func (k *k8sOps) GetClusterDomainUpdate(name string) (*v1alpha1.ClusterDomainUpdate, error) {
	return stork.Instance().GetClusterDomainUpdate(name)
}

// UpdateClusterDomainUpdate updates the ClusterDomainUpdate
func (k *k8sOps) UpdateClusterDomainUpdate(clusterDomainUpdate *v1alpha1.ClusterDomainUpdate) (*v1alpha1.ClusterDomainUpdate, error) {
	return stork.Instance().UpdateClusterDomainUpdate(clusterDomainUpdate)
}

// DeleteClusterDomainUpdate deletes the ClusterDomainUpdate
func (k *k8sOps) DeleteClusterDomainUpdate(name string) error {
	return stork.Instance().DeleteClusterDomainUpdate(name)
}

// ValidateClusterDomainUpdate validates ClusterDomainUpdate
func (k *k8sOps) ValidateClusterDomainUpdate(name string, timeout, retryInterval time.Duration) error {
	return stork.Instance().ValidateClusterDomainUpdate(name, timeout, retryInterval)
}

// ListClusterDomainUpdates lists ClusterDomainUpdates
func (k *k8sOps) ListClusterDomainUpdates() (*v1alpha1.ClusterDomainUpdateList, error) {
	return stork.Instance().ListClusterDomainUpdates()
}

// ClusterDomain CRD - END

// BackupLocation APIs - BEGIN

func (k *k8sOps) GetBackupLocation(name string, namespace string) (*v1alpha1.BackupLocation, error) {
	return stork.Instance().GetBackupLocation(name, namespace)
}

func (k *k8sOps) ListBackupLocations(namespace string) (*v1alpha1.BackupLocationList, error) {
	return stork.Instance().ListBackupLocations(namespace)
}

func (k *k8sOps) CreateBackupLocation(backupLocation *v1alpha1.BackupLocation) (*v1alpha1.BackupLocation, error) {
	return stork.Instance().CreateBackupLocation(backupLocation)
}

func (k *k8sOps) DeleteBackupLocation(name string, namespace string) error {
	return stork.Instance().DeleteBackupLocation(name, namespace)
}

func (k *k8sOps) UpdateBackupLocation(backupLocation *v1alpha1.BackupLocation) (*v1alpha1.BackupLocation, error) {
	return stork.Instance().UpdateBackupLocation(backupLocation)
}

func (k *k8sOps) ValidateBackupLocation(name, namespace string, timeout, retryInterval time.Duration) error {
	return stork.Instance().ValidateBackupLocation(name, namespace, timeout, retryInterval)
}

// BackupLocation APIs - END

// ApplicationBackupRestore APIs - BEGIN

func (k *k8sOps) GetApplicationBackup(name string, namespace string) (*v1alpha1.ApplicationBackup, error) {
	return stork.Instance().GetApplicationBackup(name, namespace)
}

func (k *k8sOps) ListApplicationBackups(namespace string) (*v1alpha1.ApplicationBackupList, error) {
	return stork.Instance().ListApplicationBackups(namespace)
}

func (k *k8sOps) CreateApplicationBackup(backup *v1alpha1.ApplicationBackup) (*v1alpha1.ApplicationBackup, error) {
	return stork.Instance().CreateApplicationBackup(backup)
}

func (k *k8sOps) DeleteApplicationBackup(name string, namespace string) error {
	return stork.Instance().DeleteApplicationBackup(name, namespace)
}

func (k *k8sOps) UpdateApplicationBackup(backup *v1alpha1.ApplicationBackup) (*v1alpha1.ApplicationBackup, error) {
	return stork.Instance().UpdateApplicationBackup(backup)
}

func (k *k8sOps) ValidateApplicationBackup(name, namespace string, timeout, retryInterval time.Duration) error {
	return stork.Instance().ValidateApplicationBackup(name, namespace, timeout, retryInterval)
}

func (k *k8sOps) GetApplicationRestore(name string, namespace string) (*v1alpha1.ApplicationRestore, error) {
	return stork.Instance().GetApplicationRestore(name, namespace)
}

func (k *k8sOps) ListApplicationRestores(namespace string) (*v1alpha1.ApplicationRestoreList, error) {
	return stork.Instance().ListApplicationRestores(namespace)
}

func (k *k8sOps) CreateApplicationRestore(restore *v1alpha1.ApplicationRestore) (*v1alpha1.ApplicationRestore, error) {
	return stork.Instance().CreateApplicationRestore(restore)
}

func (k *k8sOps) DeleteApplicationRestore(name string, namespace string) error {
	return stork.Instance().DeleteApplicationRestore(name, namespace)
}

func (k *k8sOps) ValidateApplicationRestore(name, namespace string, timeout, retryInterval time.Duration) error {
	return stork.Instance().ValidateApplicationRestore(name, namespace, timeout, retryInterval)
}

func (k *k8sOps) UpdateApplicationRestore(restore *v1alpha1.ApplicationRestore) (*v1alpha1.ApplicationRestore, error) {
	return stork.Instance().UpdateApplicationRestore(restore)
}

func (k *k8sOps) GetApplicationBackupSchedule(name string, namespace string) (*v1alpha1.ApplicationBackupSchedule, error) {
	return stork.Instance().GetApplicationBackupSchedule(name, namespace)
}

func (k *k8sOps) ListApplicationBackupSchedules(namespace string) (*v1alpha1.ApplicationBackupScheduleList, error) {
	return stork.Instance().ListApplicationBackupSchedules(namespace)
}

func (k *k8sOps) CreateApplicationBackupSchedule(applicationBackupSchedule *v1alpha1.ApplicationBackupSchedule) (*v1alpha1.ApplicationBackupSchedule, error) {
	return stork.Instance().CreateApplicationBackupSchedule(applicationBackupSchedule)
}

func (k *k8sOps) UpdateApplicationBackupSchedule(applicationBackupSchedule *v1alpha1.ApplicationBackupSchedule) (*v1alpha1.ApplicationBackupSchedule, error) {
	return stork.Instance().UpdateApplicationBackupSchedule(applicationBackupSchedule)
}

func (k *k8sOps) DeleteApplicationBackupSchedule(name string, namespace string) error {
	return stork.Instance().DeleteApplicationBackupSchedule(name, namespace)
}

func (k *k8sOps) ValidateApplicationBackupSchedule(name string, namespace string, timeout, retryInterval time.Duration) (
	map[v1alpha1.SchedulePolicyType][]*v1alpha1.ScheduledApplicationBackupStatus, error) {
	return stork.Instance().ValidateApplicationBackupSchedule(name, namespace, timeout, retryInterval)
}

// ApplicationBackupRestore APIs - END

// ApplicationClone APIs - BEGIN

func (k *k8sOps) GetApplicationClone(name string, namespace string) (*v1alpha1.ApplicationClone, error) {
	return stork.Instance().GetApplicationClone(name, namespace)
}

func (k *k8sOps) ListApplicationClones(namespace string) (*v1alpha1.ApplicationCloneList, error) {
	return stork.Instance().ListApplicationClones(namespace)
}

func (k *k8sOps) CreateApplicationClone(clone *v1alpha1.ApplicationClone) (*v1alpha1.ApplicationClone, error) {
	return stork.Instance().CreateApplicationClone(clone)
}

func (k *k8sOps) DeleteApplicationClone(name string, namespace string) error {
	return stork.Instance().DeleteApplicationClone(name, namespace)
}

func (k *k8sOps) UpdateApplicationClone(clone *v1alpha1.ApplicationClone) (*v1alpha1.ApplicationClone, error) {
	return stork.Instance().UpdateApplicationClone(clone)
}

func (k *k8sOps) ValidateApplicationClone(name, namespace string, timeout, retryInterval time.Duration) error {
	return stork.Instance().ValidateApplicationClone(name, namespace, timeout, retryInterval)
}

// ApplicationClone APIs - END
