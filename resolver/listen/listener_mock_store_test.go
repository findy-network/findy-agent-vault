// Code generated by MockGen. DO NOT EDIT.
// Source: db/store/db.go

// Package mock is a generated GoMock package.
package listen

import (
	model "github.com/findy-network/findy-agent-vault/db/model"
	model0 "github.com/findy-network/findy-agent-vault/graph/model"
	paginator "github.com/findy-network/findy-agent-vault/paginator"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDB is a mock of DB interface
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// GetListenerAgents mocks base method
func (m *MockDB) GetListenerAgents(info *paginator.BatchInfo) (*model.Agents, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListenerAgents", info)
	ret0, _ := ret[0].(*model.Agents)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListenerAgents indicates an expected call of GetListenerAgents
func (mr *MockDBMockRecorder) GetListenerAgents(info interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListenerAgents", reflect.TypeOf((*MockDB)(nil).GetListenerAgents), info)
}

// Close mocks base method
func (m *MockDB) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDB)(nil).Close))
}

// AddAgent mocks base method
func (m *MockDB) AddAgent(a *model.Agent) (*model.Agent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAgent", a)
	ret0, _ := ret[0].(*model.Agent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAgent indicates an expected call of AddAgent
func (mr *MockDBMockRecorder) AddAgent(a interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAgent", reflect.TypeOf((*MockDB)(nil).AddAgent), a)
}

// GetAgent mocks base method
func (m *MockDB) GetAgent(id, agentID *string) (*model.Agent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAgent", id, agentID)
	ret0, _ := ret[0].(*model.Agent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgent indicates an expected call of GetAgent
func (mr *MockDBMockRecorder) GetAgent(id, agentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgent", reflect.TypeOf((*MockDB)(nil).GetAgent), id, agentID)
}

// AddConnection mocks base method
func (m *MockDB) AddConnection(c *model.Connection) (*model.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddConnection", c)
	ret0, _ := ret[0].(*model.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddConnection indicates an expected call of AddConnection
func (mr *MockDBMockRecorder) AddConnection(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddConnection", reflect.TypeOf((*MockDB)(nil).AddConnection), c)
}

// GetConnection mocks base method
func (m *MockDB) GetConnection(id, tenantID string) (*model.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnection", id, tenantID)
	ret0, _ := ret[0].(*model.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnection indicates an expected call of GetConnection
func (mr *MockDBMockRecorder) GetConnection(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnection", reflect.TypeOf((*MockDB)(nil).GetConnection), id, tenantID)
}

// GetConnections mocks base method
func (m *MockDB) GetConnections(info *paginator.BatchInfo, tenantID string) (*model.Connections, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnections", info, tenantID)
	ret0, _ := ret[0].(*model.Connections)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnections indicates an expected call of GetConnections
func (mr *MockDBMockRecorder) GetConnections(info, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnections", reflect.TypeOf((*MockDB)(nil).GetConnections), info, tenantID)
}

// GetConnectionCount mocks base method
func (m *MockDB) GetConnectionCount(tenantID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectionCount", tenantID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectionCount indicates an expected call of GetConnectionCount
func (mr *MockDBMockRecorder) GetConnectionCount(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionCount", reflect.TypeOf((*MockDB)(nil).GetConnectionCount), tenantID)
}

// AddCredential mocks base method
func (m *MockDB) AddCredential(c *model.Credential) (*model.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCredential", c)
	ret0, _ := ret[0].(*model.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCredential indicates an expected call of AddCredential
func (mr *MockDBMockRecorder) AddCredential(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCredential", reflect.TypeOf((*MockDB)(nil).AddCredential), c)
}

// UpdateCredential mocks base method
func (m *MockDB) UpdateCredential(c *model.Credential) (*model.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCredential", c)
	ret0, _ := ret[0].(*model.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCredential indicates an expected call of UpdateCredential
func (mr *MockDBMockRecorder) UpdateCredential(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCredential", reflect.TypeOf((*MockDB)(nil).UpdateCredential), c)
}

// GetCredential mocks base method
func (m *MockDB) GetCredential(id, tenantID string) (*model.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCredential", id, tenantID)
	ret0, _ := ret[0].(*model.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCredential indicates an expected call of GetCredential
func (mr *MockDBMockRecorder) GetCredential(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCredential", reflect.TypeOf((*MockDB)(nil).GetCredential), id, tenantID)
}

// GetCredentials mocks base method
func (m *MockDB) GetCredentials(info *paginator.BatchInfo, tenantID string, connectionID *string) (*model.Credentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCredentials", info, tenantID, connectionID)
	ret0, _ := ret[0].(*model.Credentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCredentials indicates an expected call of GetCredentials
func (mr *MockDBMockRecorder) GetCredentials(info, tenantID, connectionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCredentials", reflect.TypeOf((*MockDB)(nil).GetCredentials), info, tenantID, connectionID)
}

// GetCredentialCount mocks base method
func (m *MockDB) GetCredentialCount(tenantID string, connectionID *string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCredentialCount", tenantID, connectionID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCredentialCount indicates an expected call of GetCredentialCount
func (mr *MockDBMockRecorder) GetCredentialCount(tenantID, connectionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCredentialCount", reflect.TypeOf((*MockDB)(nil).GetCredentialCount), tenantID, connectionID)
}

// GetConnectionForCredential mocks base method
func (m *MockDB) GetConnectionForCredential(id, tenantID string) (*model.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectionForCredential", id, tenantID)
	ret0, _ := ret[0].(*model.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectionForCredential indicates an expected call of GetConnectionForCredential
func (mr *MockDBMockRecorder) GetConnectionForCredential(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionForCredential", reflect.TypeOf((*MockDB)(nil).GetConnectionForCredential), id, tenantID)
}

// AddProof mocks base method
func (m *MockDB) AddProof(p *model.Proof) (*model.Proof, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProof", p)
	ret0, _ := ret[0].(*model.Proof)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProof indicates an expected call of AddProof
func (mr *MockDBMockRecorder) AddProof(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProof", reflect.TypeOf((*MockDB)(nil).AddProof), p)
}

// UpdateProof mocks base method
func (m *MockDB) UpdateProof(p *model.Proof) (*model.Proof, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProof", p)
	ret0, _ := ret[0].(*model.Proof)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProof indicates an expected call of UpdateProof
func (mr *MockDBMockRecorder) UpdateProof(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProof", reflect.TypeOf((*MockDB)(nil).UpdateProof), p)
}

// GetProof mocks base method
func (m *MockDB) GetProof(id, tenantID string) (*model.Proof, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProof", id, tenantID)
	ret0, _ := ret[0].(*model.Proof)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProof indicates an expected call of GetProof
func (mr *MockDBMockRecorder) GetProof(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProof", reflect.TypeOf((*MockDB)(nil).GetProof), id, tenantID)
}

// GetProofs mocks base method
func (m *MockDB) GetProofs(info *paginator.BatchInfo, tenantID string, connectionID *string) (*model.Proofs, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProofs", info, tenantID, connectionID)
	ret0, _ := ret[0].(*model.Proofs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProofs indicates an expected call of GetProofs
func (mr *MockDBMockRecorder) GetProofs(info, tenantID, connectionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProofs", reflect.TypeOf((*MockDB)(nil).GetProofs), info, tenantID, connectionID)
}

// GetProofCount mocks base method
func (m *MockDB) GetProofCount(tenantID string, connectionID *string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProofCount", tenantID, connectionID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProofCount indicates an expected call of GetProofCount
func (mr *MockDBMockRecorder) GetProofCount(tenantID, connectionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProofCount", reflect.TypeOf((*MockDB)(nil).GetProofCount), tenantID, connectionID)
}

// GetConnectionForProof mocks base method
func (m *MockDB) GetConnectionForProof(id, tenantID string) (*model.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectionForProof", id, tenantID)
	ret0, _ := ret[0].(*model.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectionForProof indicates an expected call of GetConnectionForProof
func (mr *MockDBMockRecorder) GetConnectionForProof(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionForProof", reflect.TypeOf((*MockDB)(nil).GetConnectionForProof), id, tenantID)
}

// AddMessage mocks base method
func (m_2 *MockDB) AddMessage(m *model.Message) (*model.Message, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "AddMessage", m)
	ret0, _ := ret[0].(*model.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMessage indicates an expected call of AddMessage
func (mr *MockDBMockRecorder) AddMessage(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMessage", reflect.TypeOf((*MockDB)(nil).AddMessage), m)
}

// UpdateMessage mocks base method
func (m_2 *MockDB) UpdateMessage(m *model.Message) (*model.Message, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "UpdateMessage", m)
	ret0, _ := ret[0].(*model.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMessage indicates an expected call of UpdateMessage
func (mr *MockDBMockRecorder) UpdateMessage(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessage", reflect.TypeOf((*MockDB)(nil).UpdateMessage), m)
}

// GetMessage mocks base method
func (m *MockDB) GetMessage(id, tenantID string) (*model.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessage", id, tenantID)
	ret0, _ := ret[0].(*model.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessage indicates an expected call of GetMessage
func (mr *MockDBMockRecorder) GetMessage(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessage", reflect.TypeOf((*MockDB)(nil).GetMessage), id, tenantID)
}

// GetMessages mocks base method
func (m *MockDB) GetMessages(info *paginator.BatchInfo, tenantID string, connectionID *string) (*model.Messages, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessages", info, tenantID, connectionID)
	ret0, _ := ret[0].(*model.Messages)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessages indicates an expected call of GetMessages
func (mr *MockDBMockRecorder) GetMessages(info, tenantID, connectionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessages", reflect.TypeOf((*MockDB)(nil).GetMessages), info, tenantID, connectionID)
}

// GetMessageCount mocks base method
func (m *MockDB) GetMessageCount(tenantID string, connectionID *string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageCount", tenantID, connectionID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageCount indicates an expected call of GetMessageCount
func (mr *MockDBMockRecorder) GetMessageCount(tenantID, connectionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageCount", reflect.TypeOf((*MockDB)(nil).GetMessageCount), tenantID, connectionID)
}

// GetConnectionForMessage mocks base method
func (m *MockDB) GetConnectionForMessage(id, tenantID string) (*model.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectionForMessage", id, tenantID)
	ret0, _ := ret[0].(*model.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectionForMessage indicates an expected call of GetConnectionForMessage
func (mr *MockDBMockRecorder) GetConnectionForMessage(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionForMessage", reflect.TypeOf((*MockDB)(nil).GetConnectionForMessage), id, tenantID)
}

// AddEvent mocks base method
func (m *MockDB) AddEvent(e *model.Event) (*model.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEvent", e)
	ret0, _ := ret[0].(*model.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddEvent indicates an expected call of AddEvent
func (mr *MockDBMockRecorder) AddEvent(e interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEvent", reflect.TypeOf((*MockDB)(nil).AddEvent), e)
}

// MarkEventRead mocks base method
func (m *MockDB) MarkEventRead(id, tenantID string) (*model.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkEventRead", id, tenantID)
	ret0, _ := ret[0].(*model.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkEventRead indicates an expected call of MarkEventRead
func (mr *MockDBMockRecorder) MarkEventRead(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkEventRead", reflect.TypeOf((*MockDB)(nil).MarkEventRead), id, tenantID)
}

// GetEvent mocks base method
func (m *MockDB) GetEvent(id, tenantID string) (*model.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvent", id, tenantID)
	ret0, _ := ret[0].(*model.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvent indicates an expected call of GetEvent
func (mr *MockDBMockRecorder) GetEvent(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockDB)(nil).GetEvent), id, tenantID)
}

// GetEvents mocks base method
func (m *MockDB) GetEvents(info *paginator.BatchInfo, tenantID string, connectionID *string) (*model.Events, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvents", info, tenantID, connectionID)
	ret0, _ := ret[0].(*model.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvents indicates an expected call of GetEvents
func (mr *MockDBMockRecorder) GetEvents(info, tenantID, connectionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvents", reflect.TypeOf((*MockDB)(nil).GetEvents), info, tenantID, connectionID)
}

// GetEventCount mocks base method
func (m *MockDB) GetEventCount(tenantID string, connectionID *string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventCount", tenantID, connectionID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventCount indicates an expected call of GetEventCount
func (mr *MockDBMockRecorder) GetEventCount(tenantID, connectionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventCount", reflect.TypeOf((*MockDB)(nil).GetEventCount), tenantID, connectionID)
}

// GetConnectionForEvent mocks base method
func (m *MockDB) GetConnectionForEvent(id, tenantID string) (*model.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectionForEvent", id, tenantID)
	ret0, _ := ret[0].(*model.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectionForEvent indicates an expected call of GetConnectionForEvent
func (mr *MockDBMockRecorder) GetConnectionForEvent(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionForEvent", reflect.TypeOf((*MockDB)(nil).GetConnectionForEvent), id, tenantID)
}

// GetJobForEvent mocks base method
func (m *MockDB) GetJobForEvent(id, tenantID string) (*model.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobForEvent", id, tenantID)
	ret0, _ := ret[0].(*model.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobForEvent indicates an expected call of GetJobForEvent
func (mr *MockDBMockRecorder) GetJobForEvent(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobForEvent", reflect.TypeOf((*MockDB)(nil).GetJobForEvent), id, tenantID)
}

// GetJobOutput mocks base method
func (m *MockDB) GetJobOutput(id, tenantID string, protocolType model0.ProtocolType) (*model.JobOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobOutput", id, tenantID, protocolType)
	ret0, _ := ret[0].(*model.JobOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobOutput indicates an expected call of GetJobOutput
func (mr *MockDBMockRecorder) GetJobOutput(id, tenantID, protocolType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobOutput", reflect.TypeOf((*MockDB)(nil).GetJobOutput), id, tenantID, protocolType)
}

// AddJob mocks base method
func (m *MockDB) AddJob(j *model.Job) (*model.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddJob", j)
	ret0, _ := ret[0].(*model.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddJob indicates an expected call of AddJob
func (mr *MockDBMockRecorder) AddJob(j interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddJob", reflect.TypeOf((*MockDB)(nil).AddJob), j)
}

// UpdateJob mocks base method
func (m *MockDB) UpdateJob(j *model.Job) (*model.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateJob", j)
	ret0, _ := ret[0].(*model.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateJob indicates an expected call of UpdateJob
func (mr *MockDBMockRecorder) UpdateJob(j interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateJob", reflect.TypeOf((*MockDB)(nil).UpdateJob), j)
}

// GetJob mocks base method
func (m *MockDB) GetJob(id, tenantID string) (*model.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJob", id, tenantID)
	ret0, _ := ret[0].(*model.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJob indicates an expected call of GetJob
func (mr *MockDBMockRecorder) GetJob(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJob", reflect.TypeOf((*MockDB)(nil).GetJob), id, tenantID)
}

// GetJobs mocks base method
func (m *MockDB) GetJobs(info *paginator.BatchInfo, tenantID string, connectionID *string, completed *bool) (*model.Jobs, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobs", info, tenantID, connectionID, completed)
	ret0, _ := ret[0].(*model.Jobs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobs indicates an expected call of GetJobs
func (mr *MockDBMockRecorder) GetJobs(info, tenantID, connectionID, completed interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobs", reflect.TypeOf((*MockDB)(nil).GetJobs), info, tenantID, connectionID, completed)
}

// GetJobCount mocks base method
func (m *MockDB) GetJobCount(tenantID string, connectionID *string, completed *bool) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobCount", tenantID, connectionID, completed)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobCount indicates an expected call of GetJobCount
func (mr *MockDBMockRecorder) GetJobCount(tenantID, connectionID, completed interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobCount", reflect.TypeOf((*MockDB)(nil).GetJobCount), tenantID, connectionID, completed)
}

// GetConnectionForJob mocks base method
func (m *MockDB) GetConnectionForJob(id, tenantID string) (*model.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectionForJob", id, tenantID)
	ret0, _ := ret[0].(*model.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectionForJob indicates an expected call of GetConnectionForJob
func (mr *MockDBMockRecorder) GetConnectionForJob(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionForJob", reflect.TypeOf((*MockDB)(nil).GetConnectionForJob), id, tenantID)
}