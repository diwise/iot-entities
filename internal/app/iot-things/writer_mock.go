// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package iotthings

import (
	"context"
	"github.com/diwise/iot-things/internal/app/iot-things/things"
	"sync"
)

// Ensure, that ThingsWriterMock does implement ThingsWriter.
// If this is not the case, regenerate this file with moq.
var _ ThingsWriter = &ThingsWriterMock{}

// ThingsWriterMock is a mock implementation of ThingsWriter.
//
//	func TestSomethingThatUsesThingsWriter(t *testing.T) {
//
//		// make and configure a mocked ThingsWriter
//		mockedThingsWriter := &ThingsWriterMock{
//			AddThingFunc: func(ctx context.Context, t things.Thing) error {
//				panic("mock out the AddThing method")
//			},
//			AddValueFunc: func(ctx context.Context, t things.Thing, m things.Value) error {
//				panic("mock out the AddValue method")
//			},
//			DeleteThingFunc: func(ctx context.Context, thingID string) error {
//				panic("mock out the DeleteThing method")
//			},
//			UpdateThingFunc: func(ctx context.Context, t things.Thing) error {
//				panic("mock out the UpdateThing method")
//			},
//		}
//
//		// use mockedThingsWriter in code that requires ThingsWriter
//		// and then make assertions.
//
//	}
type ThingsWriterMock struct {
	// AddThingFunc mocks the AddThing method.
	AddThingFunc func(ctx context.Context, t things.Thing) error

	// AddValueFunc mocks the AddValue method.
	AddValueFunc func(ctx context.Context, t things.Thing, m things.Value) error

	// DeleteThingFunc mocks the DeleteThing method.
	DeleteThingFunc func(ctx context.Context, thingID string) error

	// UpdateThingFunc mocks the UpdateThing method.
	UpdateThingFunc func(ctx context.Context, t things.Thing) error

	// calls tracks calls to the methods.
	calls struct {
		// AddThing holds details about calls to the AddThing method.
		AddThing []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// T is the t argument value.
			T things.Thing
		}
		// AddValue holds details about calls to the AddValue method.
		AddValue []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// T is the t argument value.
			T things.Thing
			// M is the m argument value.
			M things.Value
		}
		// DeleteThing holds details about calls to the DeleteThing method.
		DeleteThing []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ThingID is the thingID argument value.
			ThingID string
		}
		// UpdateThing holds details about calls to the UpdateThing method.
		UpdateThing []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// T is the t argument value.
			T things.Thing
		}
	}
	lockAddThing    sync.RWMutex
	lockAddValue    sync.RWMutex
	lockDeleteThing sync.RWMutex
	lockUpdateThing sync.RWMutex
}

// AddThing calls AddThingFunc.
func (mock *ThingsWriterMock) AddThing(ctx context.Context, t things.Thing) error {
	if mock.AddThingFunc == nil {
		panic("ThingsWriterMock.AddThingFunc: method is nil but ThingsWriter.AddThing was just called")
	}
	callInfo := struct {
		Ctx context.Context
		T   things.Thing
	}{
		Ctx: ctx,
		T:   t,
	}
	mock.lockAddThing.Lock()
	mock.calls.AddThing = append(mock.calls.AddThing, callInfo)
	mock.lockAddThing.Unlock()
	return mock.AddThingFunc(ctx, t)
}

// AddThingCalls gets all the calls that were made to AddThing.
// Check the length with:
//
//	len(mockedThingsWriter.AddThingCalls())
func (mock *ThingsWriterMock) AddThingCalls() []struct {
	Ctx context.Context
	T   things.Thing
} {
	var calls []struct {
		Ctx context.Context
		T   things.Thing
	}
	mock.lockAddThing.RLock()
	calls = mock.calls.AddThing
	mock.lockAddThing.RUnlock()
	return calls
}

// AddValue calls AddValueFunc.
func (mock *ThingsWriterMock) AddValue(ctx context.Context, t things.Thing, m things.Value) error {
	if mock.AddValueFunc == nil {
		panic("ThingsWriterMock.AddValueFunc: method is nil but ThingsWriter.AddValue was just called")
	}
	callInfo := struct {
		Ctx context.Context
		T   things.Thing
		M   things.Value
	}{
		Ctx: ctx,
		T:   t,
		M:   m,
	}
	mock.lockAddValue.Lock()
	mock.calls.AddValue = append(mock.calls.AddValue, callInfo)
	mock.lockAddValue.Unlock()
	return mock.AddValueFunc(ctx, t, m)
}

// AddValueCalls gets all the calls that were made to AddValue.
// Check the length with:
//
//	len(mockedThingsWriter.AddValueCalls())
func (mock *ThingsWriterMock) AddValueCalls() []struct {
	Ctx context.Context
	T   things.Thing
	M   things.Value
} {
	var calls []struct {
		Ctx context.Context
		T   things.Thing
		M   things.Value
	}
	mock.lockAddValue.RLock()
	calls = mock.calls.AddValue
	mock.lockAddValue.RUnlock()
	return calls
}

// DeleteThing calls DeleteThingFunc.
func (mock *ThingsWriterMock) DeleteThing(ctx context.Context, thingID string) error {
	if mock.DeleteThingFunc == nil {
		panic("ThingsWriterMock.DeleteThingFunc: method is nil but ThingsWriter.DeleteThing was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		ThingID string
	}{
		Ctx:     ctx,
		ThingID: thingID,
	}
	mock.lockDeleteThing.Lock()
	mock.calls.DeleteThing = append(mock.calls.DeleteThing, callInfo)
	mock.lockDeleteThing.Unlock()
	return mock.DeleteThingFunc(ctx, thingID)
}

// DeleteThingCalls gets all the calls that were made to DeleteThing.
// Check the length with:
//
//	len(mockedThingsWriter.DeleteThingCalls())
func (mock *ThingsWriterMock) DeleteThingCalls() []struct {
	Ctx     context.Context
	ThingID string
} {
	var calls []struct {
		Ctx     context.Context
		ThingID string
	}
	mock.lockDeleteThing.RLock()
	calls = mock.calls.DeleteThing
	mock.lockDeleteThing.RUnlock()
	return calls
}

// UpdateThing calls UpdateThingFunc.
func (mock *ThingsWriterMock) UpdateThing(ctx context.Context, t things.Thing) error {
	if mock.UpdateThingFunc == nil {
		panic("ThingsWriterMock.UpdateThingFunc: method is nil but ThingsWriter.UpdateThing was just called")
	}
	callInfo := struct {
		Ctx context.Context
		T   things.Thing
	}{
		Ctx: ctx,
		T:   t,
	}
	mock.lockUpdateThing.Lock()
	mock.calls.UpdateThing = append(mock.calls.UpdateThing, callInfo)
	mock.lockUpdateThing.Unlock()
	return mock.UpdateThingFunc(ctx, t)
}

// UpdateThingCalls gets all the calls that were made to UpdateThing.
// Check the length with:
//
//	len(mockedThingsWriter.UpdateThingCalls())
func (mock *ThingsWriterMock) UpdateThingCalls() []struct {
	Ctx context.Context
	T   things.Thing
} {
	var calls []struct {
		Ctx context.Context
		T   things.Thing
	}
	mock.lockUpdateThing.RLock()
	calls = mock.calls.UpdateThing
	mock.lockUpdateThing.RUnlock()
	return calls
}
