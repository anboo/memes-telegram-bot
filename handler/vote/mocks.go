// Code generated by MockGen. DO NOT EDIT.
// Source: handler/vote/interfaces.go

// Package vote is a generated GoMock package.
package vote

import (
	context "context"
	mem "memes-bot/storage/mem"
	vote "memes-bot/storage/vote"
	reflect "reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gomock "github.com/golang/mock/gomock"
)

// MockMemRepository is a mock of MemRepository interface.
type MockMemRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMemRepositoryMockRecorder
}

// MockMemRepositoryMockRecorder is the mock recorder for MockMemRepository.
type MockMemRepositoryMockRecorder struct {
	mock *MockMemRepository
}

// NewMockMemRepository creates a new mock instance.
func NewMockMemRepository(ctrl *gomock.Controller) *MockMemRepository {
	mock := &MockMemRepository{ctrl: ctrl}
	mock.recorder = &MockMemRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMemRepository) EXPECT() *MockMemRepositoryMockRecorder {
	return m.recorder
}

// Find mocks base method.
func (m *MockMemRepository) Find(ctx context.Context, id string) (mem.Mem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, id)
	ret0, _ := ret[0].(mem.Mem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockMemRepositoryMockRecorder) Find(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockMemRepository)(nil).Find), ctx, id)
}

// UpdateRating mocks base method.
func (m *MockMemRepository) UpdateRating(ctx context.Context, memId string, diff int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRating", ctx, memId, diff)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRating indicates an expected call of UpdateRating.
func (mr *MockMemRepositoryMockRecorder) UpdateRating(ctx, memId, diff interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRating", reflect.TypeOf((*MockMemRepository)(nil).UpdateRating), ctx, memId, diff)
}

// MockVoteRepository is a mock of VoteRepository interface.
type MockVoteRepository struct {
	ctrl     *gomock.Controller
	recorder *MockVoteRepositoryMockRecorder
}

// MockVoteRepositoryMockRecorder is the mock recorder for MockVoteRepository.
type MockVoteRepositoryMockRecorder struct {
	mock *MockVoteRepository
}

// NewMockVoteRepository creates a new mock instance.
func NewMockVoteRepository(ctrl *gomock.Controller) *MockVoteRepository {
	mock := &MockVoteRepository{ctrl: ctrl}
	mock.recorder = &MockVoteRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVoteRepository) EXPECT() *MockVoteRepositoryMockRecorder {
	return m.recorder
}

// Save mocks base method.
func (m *MockVoteRepository) Save(ctx context.Context, v *vote.Vote) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockVoteRepositoryMockRecorder) Save(ctx, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockVoteRepository)(nil).Save), ctx, v)
}

// MockTelegramAPI is a mock of TelegramAPI interface.
type MockTelegramAPI struct {
	ctrl     *gomock.Controller
	recorder *MockTelegramAPIMockRecorder
}

// MockTelegramAPIMockRecorder is the mock recorder for MockTelegramAPI.
type MockTelegramAPIMockRecorder struct {
	mock *MockTelegramAPI
}

// NewMockTelegramAPI creates a new mock instance.
func NewMockTelegramAPI(ctrl *gomock.Controller) *MockTelegramAPI {
	mock := &MockTelegramAPI{ctrl: ctrl}
	mock.recorder = &MockTelegramAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTelegramAPI) EXPECT() *MockTelegramAPIMockRecorder {
	return m.recorder
}

// Request mocks base method.
func (m *MockTelegramAPI) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Request", c)
	ret0, _ := ret[0].(*tgbotapi.APIResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Request indicates an expected call of Request.
func (mr *MockTelegramAPIMockRecorder) Request(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Request", reflect.TypeOf((*MockTelegramAPI)(nil).Request), c)
}