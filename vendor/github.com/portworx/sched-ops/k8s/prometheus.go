package k8s

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/portworx/sched-ops/k8s/prometheus"
)

// PrometheusOps is an interface to perform Prometheus object operations
type PrometheusOps interface {
	ServiceMonitorOps
	PrometheusPodOps
	PrometheusRuleOps
	AlertManagerOps
}

// PrometheusPodOps is an interface to perform Prometheus operations
type PrometheusPodOps interface {
	// ListPrometheuses lists all prometheus instances in a given namespace
	ListPrometheuses(namespace string) (*monitoringv1.PrometheusList, error)
	// GetPrometheus gets the prometheus instance that matches the given name
	GetPrometheus(name, namespace string) (*monitoringv1.Prometheus, error)
	// CreatePrometheus creates the given prometheus
	CreatePrometheus(*monitoringv1.Prometheus) (*monitoringv1.Prometheus, error)
	// UpdatePrometheus updates the given prometheus
	UpdatePrometheus(*monitoringv1.Prometheus) (*monitoringv1.Prometheus, error)
	// DeletePrometheus deletes the given prometheus
	DeletePrometheus(name, namespace string) error
}

// ServiceMonitorOps is an interface to perform ServiceMonitor operations
type ServiceMonitorOps interface {
	// ListServiceMonitors lists all servicemonitors in a given namespace
	ListServiceMonitors(namespace string) (*monitoringv1.ServiceMonitorList, error)
	// GetServiceMonitor gets the service monitor instance that matches the given name
	GetServiceMonitor(name, namespace string) (*monitoringv1.ServiceMonitor, error)
	// CreateServiceMonitor creates the given service monitor
	CreateServiceMonitor(*monitoringv1.ServiceMonitor) (*monitoringv1.ServiceMonitor, error)
	// UpdateServiceMonitor updates the given service monitor
	UpdateServiceMonitor(*monitoringv1.ServiceMonitor) (*monitoringv1.ServiceMonitor, error)
	// DeleteServiceMonitor deletes the given service monitor
	DeleteServiceMonitor(name, namespace string) error
}

// PrometheusRuleOps is an interface to perform PrometheusRule operations
type PrometheusRuleOps interface {
	// ListPrometheusRule creates the given prometheus rule
	ListPrometheusRules(namespace string) (*monitoringv1.PrometheusRuleList, error)
	// GetPrometheusRule gets the prometheus rule that matches the given name
	GetPrometheusRule(name, namespace string) (*monitoringv1.PrometheusRule, error)
	// CreatePrometheusRule creates the given prometheus rule
	CreatePrometheusRule(*monitoringv1.PrometheusRule) (*monitoringv1.PrometheusRule, error)
	// UpdatePrometheusRule updates the given prometheus rule
	UpdatePrometheusRule(*monitoringv1.PrometheusRule) (*monitoringv1.PrometheusRule, error)
	// DeletePrometheusRule deletes the given prometheus rule
	DeletePrometheusRule(name, namespace string) error
}

// AlertManagerOps is an interface to perform AlertManager operations
type AlertManagerOps interface {
	// ListAlertManagerss lists all alertmanager instances in a given namespace
	ListAlertManagers(namespace string) (*monitoringv1.AlertmanagerList, error)
	// GetAlertManager gets the alert manager that matches the given name
	GetAlertManager(name, namespace string) (*monitoringv1.Alertmanager, error)
	// CreateAlertManager creates the given alert manager
	CreateAlertManager(*monitoringv1.Alertmanager) (*monitoringv1.Alertmanager, error)
	// UpdateAlertManager updates the given alert manager
	UpdateAlertManager(*monitoringv1.Alertmanager) (*monitoringv1.Alertmanager, error)
	// DeleteAlertManager deletes the given alert manager
	DeleteAlertManager(name, namespace string) error
}

func (k *k8sOps) ListPrometheuses(namespace string) (*monitoringv1.PrometheusList, error) {
	return prometheus.Instance().ListPrometheuses(namespace)
}

func (k *k8sOps) GetPrometheus(name string, namespace string) (*monitoringv1.Prometheus, error) {
	return prometheus.Instance().GetPrometheus(name, namespace)
}

func (k *k8sOps) CreatePrometheus(p *monitoringv1.Prometheus) (*monitoringv1.Prometheus, error) {
	return prometheus.Instance().CreatePrometheus(p)
}

func (k *k8sOps) UpdatePrometheus(p *monitoringv1.Prometheus) (*monitoringv1.Prometheus, error) {
	return prometheus.Instance().UpdatePrometheus(p)
}

func (k *k8sOps) DeletePrometheus(name, namespace string) error {
	return prometheus.Instance().DeletePrometheus(name, namespace)
}

func (k *k8sOps) ListServiceMonitors(namespace string) (*monitoringv1.ServiceMonitorList, error) {
	return prometheus.Instance().ListServiceMonitors(namespace)
}

func (k *k8sOps) GetServiceMonitor(name string, namespace string) (*monitoringv1.ServiceMonitor, error) {
	return prometheus.Instance().GetServiceMonitor(name, namespace)
}

func (k *k8sOps) CreateServiceMonitor(serviceMonitor *monitoringv1.ServiceMonitor) (*monitoringv1.ServiceMonitor, error) {
	return prometheus.Instance().CreateServiceMonitor(serviceMonitor)
}

func (k *k8sOps) UpdateServiceMonitor(serviceMonitor *monitoringv1.ServiceMonitor) (*monitoringv1.ServiceMonitor, error) {
	return prometheus.Instance().UpdateServiceMonitor(serviceMonitor)
}

func (k *k8sOps) DeleteServiceMonitor(name, namespace string) error {
	return prometheus.Instance().DeleteServiceMonitor(name, namespace)
}

func (k *k8sOps) ListPrometheusRules(namespace string) (*monitoringv1.PrometheusRuleList, error) {
	return prometheus.Instance().ListPrometheusRules(namespace)
}

func (k *k8sOps) GetPrometheusRule(name string, namespace string) (*monitoringv1.PrometheusRule, error) {
	return prometheus.Instance().GetPrometheusRule(name, namespace)
}

func (k *k8sOps) CreatePrometheusRule(rule *monitoringv1.PrometheusRule) (*monitoringv1.PrometheusRule, error) {
	return prometheus.Instance().CreatePrometheusRule(rule)
}

func (k *k8sOps) UpdatePrometheusRule(rule *monitoringv1.PrometheusRule) (*monitoringv1.PrometheusRule, error) {
	return prometheus.Instance().UpdatePrometheusRule(rule)
}

func (k *k8sOps) DeletePrometheusRule(name, namespace string) error {
	return prometheus.Instance().DeletePrometheusRule(name, namespace)
}

func (k *k8sOps) ListAlertManagers(namespace string) (*monitoringv1.AlertmanagerList, error) {
	return prometheus.Instance().ListAlertManagers(namespace)
}

func (k *k8sOps) GetAlertManager(name string, namespace string) (*monitoringv1.Alertmanager, error) {
	return prometheus.Instance().GetAlertManager(name, namespace)
}

func (k *k8sOps) CreateAlertManager(alertmanager *monitoringv1.Alertmanager) (*monitoringv1.Alertmanager, error) {
	return prometheus.Instance().CreateAlertManager(alertmanager)
}

func (k *k8sOps) UpdateAlertManager(alertmanager *monitoringv1.Alertmanager) (*monitoringv1.Alertmanager, error) {
	return prometheus.Instance().UpdateAlertManager(alertmanager)
}

func (k *k8sOps) DeleteAlertManager(name, namespace string) error {
	return prometheus.Instance().DeleteAlertManager(name, namespace)
}
