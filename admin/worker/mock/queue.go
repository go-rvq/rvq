// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/go-rvq/rvq/admin/worker"
	"sync"
)

// Ensure, that QueueMock does implement worker.Queue.
// If this is not the case, regenerate this file with moq.
var _ worker.Queue = &QueueMock{}

// QueueMock is a mock implementation of worker.Queue.
//
//	func TestSomethingThatUsesQueue(t *testing.T) {
//
//		// make and configure a mocked worker.Queue
//		mockedQueue := &QueueMock{
//			AddFunc: func(ctx context.Context, job worker.QueJobInterface) error {
//				panic("mock out the Add method")
//			},
//			KillFunc: func(ctx context.Context, job worker.QueJobInterface) error {
//				panic("mock out the Kill method")
//			},
//			ListenFunc: func(jobDefs []*worker.QorJobDefinition, getJob func(qorJobID uint) (worker.QueJobInterface, error)) error {
//				panic("mock out the Listen method")
//			},
//			RemoveFunc: func(ctx context.Context, job worker.QueJobInterface) error {
//				panic("mock out the Remove method")
//			},
//			ShutdownFunc: func(ctx context.Context) error {
//				panic("mock out the Shutdown method")
//			},
//		}
//
//		// use mockedQueue in code that requires worker.Queue
//		// and then make assertions.
//
//	}
type QueueMock struct {
	// AddFunc mocks the Add method.
	AddFunc func(ctx context.Context, job worker.QueJobInterface) error

	// KillFunc mocks the Kill method.
	KillFunc func(ctx context.Context, job worker.QueJobInterface) error

	// ListenFunc mocks the Listen method.
	ListenFunc func(jobDefs []*worker.QorJobDefinition, getJob func(qorJobID uint) (worker.QueJobInterface, error)) error

	// RemoveFunc mocks the Remove method.
	RemoveFunc func(ctx context.Context, job worker.QueJobInterface) error

	// ShutdownFunc mocks the Shutdown method.
	ShutdownFunc func(ctx context.Context) error

	// calls tracks calls to the methods.
	calls struct {
		// Add holds details about calls to the Add method.
		Add []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Job is the job argument value.
			Job worker.QueJobInterface
		}
		// Kill holds details about calls to the Kill method.
		Kill []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Job is the job argument value.
			Job worker.QueJobInterface
		}
		// Listen holds details about calls to the Listen method.
		Listen []struct {
			// JobDefs is the jobDefs argument value.
			JobDefs []*worker.QorJobDefinition
			// GetJob is the getJob argument value.
			GetJob func(qorJobID uint) (worker.QueJobInterface, error)
		}
		// Remove holds details about calls to the Remove method.
		Remove []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Job is the job argument value.
			Job worker.QueJobInterface
		}
		// Shutdown holds details about calls to the Shutdown method.
		Shutdown []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
	}
	lockAdd      sync.RWMutex
	lockKill     sync.RWMutex
	lockListen   sync.RWMutex
	lockRemove   sync.RWMutex
	lockShutdown sync.RWMutex
}

// Add calls AddFunc.
func (mock *QueueMock) Add(ctx context.Context, job worker.QueJobInterface) error {
	if mock.AddFunc == nil {
		panic("QueueMock.AddFunc: method is nil but Queue.Add was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Job worker.QueJobInterface
	}{
		Ctx: ctx,
		Job: job,
	}
	mock.lockAdd.Lock()
	mock.calls.Add = append(mock.calls.Add, callInfo)
	mock.lockAdd.Unlock()
	return mock.AddFunc(ctx, job)
}

// AddCalls gets all the calls that were made to Add.
// Check the length with:
//
//	len(mockedQueue.AddCalls())
func (mock *QueueMock) AddCalls() []struct {
	Ctx context.Context
	Job worker.QueJobInterface
} {
	var calls []struct {
		Ctx context.Context
		Job worker.QueJobInterface
	}
	mock.lockAdd.RLock()
	calls = mock.calls.Add
	mock.lockAdd.RUnlock()
	return calls
}

// Kill calls KillFunc.
func (mock *QueueMock) Kill(ctx context.Context, job worker.QueJobInterface) error {
	if mock.KillFunc == nil {
		panic("QueueMock.KillFunc: method is nil but Queue.Kill was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Job worker.QueJobInterface
	}{
		Ctx: ctx,
		Job: job,
	}
	mock.lockKill.Lock()
	mock.calls.Kill = append(mock.calls.Kill, callInfo)
	mock.lockKill.Unlock()
	return mock.KillFunc(ctx, job)
}

// KillCalls gets all the calls that were made to Kill.
// Check the length with:
//
//	len(mockedQueue.KillCalls())
func (mock *QueueMock) KillCalls() []struct {
	Ctx context.Context
	Job worker.QueJobInterface
} {
	var calls []struct {
		Ctx context.Context
		Job worker.QueJobInterface
	}
	mock.lockKill.RLock()
	calls = mock.calls.Kill
	mock.lockKill.RUnlock()
	return calls
}

// Listen calls ListenFunc.
func (mock *QueueMock) Listen(jobDefs []*worker.QorJobDefinition, getJob func(qorJobID uint) (worker.QueJobInterface, error)) error {
	if mock.ListenFunc == nil {
		panic("QueueMock.ListenFunc: method is nil but Queue.Listen was just called")
	}
	callInfo := struct {
		JobDefs []*worker.QorJobDefinition
		GetJob  func(qorJobID uint) (worker.QueJobInterface, error)
	}{
		JobDefs: jobDefs,
		GetJob:  getJob,
	}
	mock.lockListen.Lock()
	mock.calls.Listen = append(mock.calls.Listen, callInfo)
	mock.lockListen.Unlock()
	return mock.ListenFunc(jobDefs, getJob)
}

// ListenCalls gets all the calls that were made to Listen.
// Check the length with:
//
//	len(mockedQueue.ListenCalls())
func (mock *QueueMock) ListenCalls() []struct {
	JobDefs []*worker.QorJobDefinition
	GetJob  func(qorJobID uint) (worker.QueJobInterface, error)
} {
	var calls []struct {
		JobDefs []*worker.QorJobDefinition
		GetJob  func(qorJobID uint) (worker.QueJobInterface, error)
	}
	mock.lockListen.RLock()
	calls = mock.calls.Listen
	mock.lockListen.RUnlock()
	return calls
}

// Remove calls RemoveFunc.
func (mock *QueueMock) Remove(ctx context.Context, job worker.QueJobInterface) error {
	if mock.RemoveFunc == nil {
		panic("QueueMock.RemoveFunc: method is nil but Queue.Remove was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Job worker.QueJobInterface
	}{
		Ctx: ctx,
		Job: job,
	}
	mock.lockRemove.Lock()
	mock.calls.Remove = append(mock.calls.Remove, callInfo)
	mock.lockRemove.Unlock()
	return mock.RemoveFunc(ctx, job)
}

// RemoveCalls gets all the calls that were made to Remove.
// Check the length with:
//
//	len(mockedQueue.RemoveCalls())
func (mock *QueueMock) RemoveCalls() []struct {
	Ctx context.Context
	Job worker.QueJobInterface
} {
	var calls []struct {
		Ctx context.Context
		Job worker.QueJobInterface
	}
	mock.lockRemove.RLock()
	calls = mock.calls.Remove
	mock.lockRemove.RUnlock()
	return calls
}

// Shutdown calls ShutdownFunc.
func (mock *QueueMock) Shutdown(ctx context.Context) error {
	if mock.ShutdownFunc == nil {
		panic("QueueMock.ShutdownFunc: method is nil but Queue.Shutdown was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockShutdown.Lock()
	mock.calls.Shutdown = append(mock.calls.Shutdown, callInfo)
	mock.lockShutdown.Unlock()
	return mock.ShutdownFunc(ctx)
}

// ShutdownCalls gets all the calls that were made to Shutdown.
// Check the length with:
//
//	len(mockedQueue.ShutdownCalls())
func (mock *QueueMock) ShutdownCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockShutdown.RLock()
	calls = mock.calls.Shutdown
	mock.lockShutdown.RUnlock()
	return calls
}
