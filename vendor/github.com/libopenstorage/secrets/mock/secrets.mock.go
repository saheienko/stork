// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/libopenstorage/secrets (interfaces: Secrets)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSecrets is a mock of Secrets interface
type MockSecrets struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsMockRecorder
}

// MockSecretsMockRecorder is the mock recorder for MockSecrets
type MockSecretsMockRecorder struct {
	mock *MockSecrets
}

// NewMockSecrets creates a new mock instance
func NewMockSecrets(ctrl *gomock.Controller) *MockSecrets {
	mock := &MockSecrets{ctrl: ctrl}
	mock.recorder = &MockSecretsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSecrets) EXPECT() *MockSecretsMockRecorder {
	return m.recorder
}

// Decrypt mocks base method
func (m *MockSecrets) Decrypt(arg0, arg1 string, arg2 map[string]string) (string, error) {
	ret := m.ctrl.Call(m, "Decrypt", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decrypt indicates an expected call of Decrypt
func (mr *MockSecretsMockRecorder) Decrypt(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decrypt", reflect.TypeOf((*MockSecrets)(nil).Decrypt), arg0, arg1, arg2)
}

// DeleteSecret mocks base method
func (m *MockSecrets) DeleteSecret(arg0 string, arg1 map[string]string) error {
	ret := m.ctrl.Call(m, "DeleteSecret", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSecret indicates an expected call of DeleteSecret
func (mr *MockSecretsMockRecorder) DeleteSecret(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecret", reflect.TypeOf((*MockSecrets)(nil).DeleteSecret), arg0, arg1)
}

// Encrypt mocks base method
func (m *MockSecrets) Encrypt(arg0, arg1 string, arg2 map[string]string) (string, error) {
	ret := m.ctrl.Call(m, "Encrypt", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Encrypt indicates an expected call of Encrypt
func (mr *MockSecretsMockRecorder) Encrypt(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encrypt", reflect.TypeOf((*MockSecrets)(nil).Encrypt), arg0, arg1, arg2)
}

// GetSecret mocks base method
func (m *MockSecrets) GetSecret(arg0 string, arg1 map[string]string) (map[string]interface{}, error) {
	ret := m.ctrl.Call(m, "GetSecret", arg0, arg1)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret
func (mr *MockSecretsMockRecorder) GetSecret(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockSecrets)(nil).GetSecret), arg0, arg1)
}

// ListSecrets mocks base method
func (m *MockSecrets) ListSecrets() ([]string, error) {
	ret := m.ctrl.Call(m, "ListSecrets")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecrets indicates an expected call of ListSecrets
func (mr *MockSecretsMockRecorder) ListSecrets() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecrets", reflect.TypeOf((*MockSecrets)(nil).ListSecrets))
}

// PutSecret mocks base method
func (m *MockSecrets) PutSecret(arg0 string, arg1 map[string]interface{}, arg2 map[string]string) error {
	ret := m.ctrl.Call(m, "PutSecret", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutSecret indicates an expected call of PutSecret
func (mr *MockSecretsMockRecorder) PutSecret(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutSecret", reflect.TypeOf((*MockSecrets)(nil).PutSecret), arg0, arg1, arg2)
}

// Rencrypt mocks base method
func (m *MockSecrets) Rencrypt(arg0, arg1 string, arg2, arg3 map[string]string, arg4 string) (string, error) {
	ret := m.ctrl.Call(m, "Rencrypt", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Rencrypt indicates an expected call of Rencrypt
func (mr *MockSecretsMockRecorder) Rencrypt(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rencrypt", reflect.TypeOf((*MockSecrets)(nil).Rencrypt), arg0, arg1, arg2, arg3, arg4)
}

// String mocks base method
func (m *MockSecrets) String() string {
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String
func (mr *MockSecretsMockRecorder) String() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockSecrets)(nil).String))
}
