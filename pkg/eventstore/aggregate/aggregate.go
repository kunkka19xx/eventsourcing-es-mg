package aggregate

import (
	"eventstore-intro/pkg/eventstore/config"
	"eventstore-intro/pkg/eventstore/event"
	"fmt"
	"github.com/pkg/errors"
)

type AggregateRoot interface {
	GetUncommittedEvents() []event.EventModel
	ClearUncommittedEvents()
	GetID() string
	SetID(id string) *AggregateBase
	GetVersion() int64
	SetType(aggregateType AggregateType)
	GetType() AggregateType
	SetAppliedEvents(events []event.EventModel)
	GetAppliedEvents() []event.EventModel
	Apply(event event.EventModel) error
	RaiseEvent(event event.EventModel) error
	When(event event.EventModel) error
	Load(events []event.EventModel) error
}

type AggregateBase struct {
	ID                string
	Version           int64
	AppliedEvents     []event.EventModel
	UncommittedEvents []event.EventModel
	Type              AggregateType
	onEvent           onEvent
	appliedEvents     bool
}

type onEvent func(event event.EventModel) error

type Aggregate interface {
	AggregateRoot
}

type AggregateType string

func NewAggregateBase(onEvent onEvent) *AggregateBase {
	if onEvent == nil {
		return nil
	}

	return &AggregateBase{
		Version:           config.AGGREGATE_START_VERSION,
		AppliedEvents:     make([]event.EventModel, 0, config.CHANGES_EVENT_CAPAPACITY),
		UncommittedEvents: make([]event.EventModel, 0, config.CHANGES_EVENT_CAPAPACITY),
		onEvent:           onEvent,
		appliedEvents:     false,
	}
}

// GetUncommittedEvents get AggregateBase uncommitted Event's
func (a *AggregateBase) GetUncommittedEvents() []event.EventModel {
	return a.UncommittedEvents
}

// SetID set AggregateBase ID
func (a *AggregateBase) SetID(id string) *AggregateBase {
	a.ID = fmt.Sprintf("%s-%s", a.GetType(), id)
	return a
}

// GetID get AggregateBase ID
func (a *AggregateBase) GetID() string {
	return a.ID
}

// SetType set AggregateBase AggregateType
func (a *AggregateBase) SetType(aggregateType AggregateType) {
	a.Type = aggregateType
}

// GetVersion get AggregateBase version
func (a *AggregateBase) GetVersion() int64 {
	return a.Version
}

// ClearUncommittedEvents clear AggregateBase uncommitted Event's
func (a *AggregateBase) ClearUncommittedEvents() {
	a.UncommittedEvents = make([]event.EventModel, 0, config.CHANGES_EVENT_CAPAPACITY)
}

// GetAppliedEvents get AggregateBase applied Event's
func (a *AggregateBase) GetAppliedEvents() []event.EventModel {
	return a.AppliedEvents
}

// SetAppliedEvents set AggregateBase applied Event's
func (a *AggregateBase) SetAppliedEvents(events []event.EventModel) {
	a.AppliedEvents = events
}

// Load add existing events from event store to aggregate using When interface method
func (a *AggregateBase) Load(events []event.EventModel) error {

	for _, evt := range events {
		if evt.GetAggregateID() != a.GetID() {
			return errors.New("ErrInvalidAggregate")
		}

		if err := a.onEvent(evt); err != nil {
			return err
		}

		if a.appliedEvents {
			a.AppliedEvents = append(a.AppliedEvents, evt)
		}
		a.Version++
	}

	return nil
}

// Apply push event to aggregate uncommitted events using When method
func (a *AggregateBase) Apply(event event.EventModel) error {
	if event.GetAggregateID() != a.GetID() {
		return errors.New("ErrInvalidAggregate")
	}
	event.SetAggregateType(a.GetType())
	if err := a.onEvent(event); err != nil {
		return err
	}
	a.Version++
	event.SetVersion(a.GetVersion())
	a.UncommittedEvents = append(a.UncommittedEvents, event)
	return nil
}

// RaiseEvent push event to aggregate applied events using When method, used for load directly from eventstore
func (a *AggregateBase) RaiseEvent(event event.EventModel) error {
	if event.GetAggregateID() != a.GetID() {
		return errors.New("Invalid ID")
	}
	event.SetAggregateType(a.GetType())
	if err := a.onEvent(event); err != nil {
		return err
	}
	if a.appliedEvents {
		a.AppliedEvents = append(a.AppliedEvents, event)
	}
	a.Version = event.GetVersion()
	return nil
}

func (a *AggregateBase) GetType() AggregateType {
	return a.Type
}
