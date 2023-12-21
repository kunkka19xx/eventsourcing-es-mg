package aggregate

import (
	"eventstore-intro/eventstore/config"
	"eventstore-intro/eventstore/event"
)

type AggregateRoot interface {
	GetUncommittedEvents() []event.EventModel
	GetID() string
	SetID(id string) *AggregateBase
	GetVersion() int64
	//ClearUncommittedEvents()
	ToSnapshot()
	SetType(aggregateType AggregateType)
	GetType() AggregateType
	SetAppliedEvents(events []event.EventModel)
	GetAppliedEvents() []event.EventModel
	RaiseEvent(event event.EventModel) error
	String() string
	Load
	Apply
}

type AggregateBase struct {
	ID      string
	Version int64
	Changes []event.EventModel
	Type    AggregateType
	when    when
}

type When interface {
	When(event event.EventModel) error
}

type when func(event event.EventModel) error

type RaiseEvent interface {
	RaiseEvent(event event.EventModel) error
}

type Load interface {
	Load(events []event.EventModel) error
}

type Aggregate interface {
	When
	AggregateRoot
}

type Apply interface {
	Apply(event event.EventModel) error
}

type AggregateType string

func NewAggregateBase(when when) *AggregateBase {
	if when == nil {
		return nil
	}

	return &AggregateBase{
		Version: config.AGGREGATE_START_VERSION,
		Changes: make([]event.EventModel, 0, config.CHANGES_EVENT_CAPAPACITY),
		when:    when,
	}
}

// GetUncommittedEvents get AggregateBase uncommitted Event's
func (a *AggregateBase) GetUncommittedEvents() []event.EventModel {
	return a.Changes
}

// Apply push event to aggregate uncommitted events using When method
func (a *AggregateBase) Apply(event event.EventModel) error {

	if err := a.when(event); err != nil {
		return err
	}

	a.Version++
	a.Changes = append(a.Changes, event)
	return nil
}

// Load add existing events from event store to aggregate using When interface method
func (a *AggregateBase) Load(events []event.EventModel) error {

	for _, evt := range events {
		if err := a.when(evt); err != nil {
			return err
		}
		a.Version++
	}

	return nil
}

func (a *AggregateBase) RaiseEvent(event event.EventModel) error {

	if err := a.when(event); err != nil {
		return err
	}
	
	a.Version++
	return nil
}

// ToSnapshot prepare AggregateBase for saving Snapshot.
func (a *AggregateBase) ToSnapshot() {
	a.ClearChanges()
}

// GetVersion get AggregateBase version
func (a *AggregateBase) GetVersion() int64 {
	return a.Version
}

// ClearChanges clear AggregateBase uncommitted Event's
func (a *AggregateBase) ClearChanges() {
	a.Changes = make([]event.EventModel, 0, config.CHANGES_EVENT_CAPAPACITY)
}
func (a *AggregateBase) GetType() AggregateType {
	return a.Type
}
