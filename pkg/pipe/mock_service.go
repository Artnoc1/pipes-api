// Code generated by mockery v1.0.0. DO NOT EDIT.

package pipe

import integrations "github.com/toggl/pipes-api/pkg/integrations"
import mock "github.com/stretchr/testify/mock"
import toggl "github.com/toggl/pipes-api/pkg/toggl"

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// AvailablePipeType provides a mock function with given fields: pid
func (_m *MockService) AvailablePipeType(pid integrations.PipeID) bool {
	ret := _m.Called(pid)

	var r0 bool
	if rf, ok := ret.Get(0).(func(integrations.PipeID) bool); ok {
		r0 = rf(pid)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// AvailableServiceType provides a mock function with given fields: sid
func (_m *MockService) AvailableServiceType(sid integrations.ExternalServiceID) bool {
	ret := _m.Called(sid)

	var r0 bool
	if rf, ok := ret.Get(0).(func(integrations.ExternalServiceID) bool); ok {
		r0 = rf(sid)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ClearPipeConnections provides a mock function with given fields: workspaceID, sid, pid
func (_m *MockService) ClearPipeConnections(workspaceID int, sid integrations.ExternalServiceID, pid integrations.PipeID) error {
	ret := _m.Called(workspaceID, sid, pid)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, integrations.ExternalServiceID, integrations.PipeID) error); ok {
		r0 = rf(workspaceID, sid, pid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateAuthorization provides a mock function with given fields: workspaceID, sid, currentWorkspaceToken, oAuthRawData
func (_m *MockService) CreateAuthorization(workspaceID int, sid integrations.ExternalServiceID, currentWorkspaceToken string, oAuthRawData []byte) error {
	ret := _m.Called(workspaceID, sid, currentWorkspaceToken, oAuthRawData)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, integrations.ExternalServiceID, string, []byte) error); ok {
		r0 = rf(workspaceID, sid, currentWorkspaceToken, oAuthRawData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreatePipe provides a mock function with given fields: workspaceID, sid, pid, params
func (_m *MockService) CreatePipe(workspaceID int, sid integrations.ExternalServiceID, pid integrations.PipeID, params []byte) error {
	ret := _m.Called(workspaceID, sid, pid, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, integrations.ExternalServiceID, integrations.PipeID, []byte) error); ok {
		r0 = rf(workspaceID, sid, pid, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAuthorization provides a mock function with given fields: workspaceID, sid
func (_m *MockService) DeleteAuthorization(workspaceID int, sid integrations.ExternalServiceID) error {
	ret := _m.Called(workspaceID, sid)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, integrations.ExternalServiceID) error); ok {
		r0 = rf(workspaceID, sid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePipe provides a mock function with given fields: workspaceID, sid, pid
func (_m *MockService) DeletePipe(workspaceID int, sid integrations.ExternalServiceID, pid integrations.PipeID) error {
	ret := _m.Called(workspaceID, sid, pid)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, integrations.ExternalServiceID, integrations.PipeID) error); ok {
		r0 = rf(workspaceID, sid, pid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAuthURL provides a mock function with given fields: sid, accountName, callbackURL
func (_m *MockService) GetAuthURL(sid integrations.ExternalServiceID, accountName string, callbackURL string) (string, error) {
	ret := _m.Called(sid, accountName, callbackURL)

	var r0 string
	if rf, ok := ret.Get(0).(func(integrations.ExternalServiceID, string, string) string); ok {
		r0 = rf(sid, accountName, callbackURL)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(integrations.ExternalServiceID, string, string) error); ok {
		r1 = rf(sid, accountName, callbackURL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPipe provides a mock function with given fields: workspaceID, sid, pid
func (_m *MockService) GetPipe(workspaceID int, sid integrations.ExternalServiceID, pid integrations.PipeID) (*Pipe, error) {
	ret := _m.Called(workspaceID, sid, pid)

	var r0 *Pipe
	if rf, ok := ret.Get(0).(func(int, integrations.ExternalServiceID, integrations.PipeID) *Pipe); ok {
		r0 = rf(workspaceID, sid, pid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Pipe)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, integrations.ExternalServiceID, integrations.PipeID) error); ok {
		r1 = rf(workspaceID, sid, pid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetServiceAccounts provides a mock function with given fields: workspaceID, sid, forceImport
func (_m *MockService) GetServiceAccounts(workspaceID int, sid integrations.ExternalServiceID, forceImport bool) (*toggl.AccountsResponse, error) {
	ret := _m.Called(workspaceID, sid, forceImport)

	var r0 *toggl.AccountsResponse
	if rf, ok := ret.Get(0).(func(int, integrations.ExternalServiceID, bool) *toggl.AccountsResponse); ok {
		r0 = rf(workspaceID, sid, forceImport)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*toggl.AccountsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, integrations.ExternalServiceID, bool) error); ok {
		r1 = rf(workspaceID, sid, forceImport)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetServicePipeLog provides a mock function with given fields: workspaceID, sid, pid
func (_m *MockService) GetServicePipeLog(workspaceID int, sid integrations.ExternalServiceID, pid integrations.PipeID) (string, error) {
	ret := _m.Called(workspaceID, sid, pid)

	var r0 string
	if rf, ok := ret.Get(0).(func(int, integrations.ExternalServiceID, integrations.PipeID) string); ok {
		r0 = rf(workspaceID, sid, pid)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, integrations.ExternalServiceID, integrations.PipeID) error); ok {
		r1 = rf(workspaceID, sid, pid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetServiceUsers provides a mock function with given fields: workspaceID, sid, forceImport
func (_m *MockService) GetServiceUsers(workspaceID int, sid integrations.ExternalServiceID, forceImport bool) (*toggl.UsersResponse, error) {
	ret := _m.Called(workspaceID, sid, forceImport)

	var r0 *toggl.UsersResponse
	if rf, ok := ret.Get(0).(func(int, integrations.ExternalServiceID, bool) *toggl.UsersResponse); ok {
		r0 = rf(workspaceID, sid, forceImport)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*toggl.UsersResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, integrations.ExternalServiceID, bool) error); ok {
		r1 = rf(workspaceID, sid, forceImport)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Ready provides a mock function with given fields:
func (_m *MockService) Ready() []error {
	ret := _m.Called()

	var r0 []error
	if rf, ok := ret.Get(0).(func() []error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]error)
		}
	}

	return r0
}

// Run provides a mock function with given fields: _a0
func (_m *MockService) Run(_a0 *Pipe) {
	_m.Called(_a0)
}

// RunPipe provides a mock function with given fields: workspaceID, sid, pid, usersSelector
func (_m *MockService) RunPipe(workspaceID int, sid integrations.ExternalServiceID, pid integrations.PipeID, usersSelector []byte) error {
	ret := _m.Called(workspaceID, sid, pid, usersSelector)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, integrations.ExternalServiceID, integrations.PipeID, []byte) error); ok {
		r0 = rf(workspaceID, sid, pid, usersSelector)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePipe provides a mock function with given fields: workspaceID, sid, pid, params
func (_m *MockService) UpdatePipe(workspaceID int, sid integrations.ExternalServiceID, pid integrations.PipeID, params []byte) error {
	ret := _m.Called(workspaceID, sid, pid, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, integrations.ExternalServiceID, integrations.PipeID, []byte) error); ok {
		r0 = rf(workspaceID, sid, pid, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WorkspaceIntegrations provides a mock function with given fields: workspaceID
func (_m *MockService) WorkspaceIntegrations(workspaceID int) ([]Integration, error) {
	ret := _m.Called(workspaceID)

	var r0 []Integration
	if rf, ok := ret.Get(0).(func(int) []Integration); ok {
		r0 = rf(workspaceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Integration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(workspaceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
