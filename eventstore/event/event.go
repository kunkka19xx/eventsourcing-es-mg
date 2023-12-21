package event

import (
	"encoding/json"
	"eventstore-intro/eventstore/aggregate"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/google/uuid"
	"time"
)

type EventModel struct {
	ID            string
	Timestamp     time.Time
	AggregateID   string
	AggregateType aggregate.AggregateType
	Version       int64
	EventType     string
	Data          []byte
	Metadata      []byte
}

func NewBaseEvent(aggregate aggregate.Aggregate, eventType string) EventModel {
	return EventModel{
		ID:            uuid.New().String(),
		Timestamp:     time.Now().Local(),
		AggregateID:   aggregate.GetID(),
		AggregateType: aggregate.GetType(),
		Version:       aggregate.GetVersion(),
		EventType:     eventType,
		//Metadata:      event.Metadata(),
	}
}

func NewRecordedEvent(event *esdb.RecordedEvent) EventModel {
	return EventModel{
		ID:          event.EventID.String(),
		EventType:   event.EventType,
		Data:        event.Data,
		Timestamp:   event.CreatedDate,
		AggregateID: event.StreamID,
		Version:     int64(event.EventNumber),
		Metadata:    event.UserMetadata,
	}
}

// GetEventID get EventID of the Event.
func (e *EventModel) GetEventID() string {
	return e.ID
}

// GetTimeStamp get timestamp of the Event.
func (e *EventModel) GetTimeStamp() time.Time {
	return e.Timestamp
}

// GetData The data attached to the Event serialized to bytes.
func (e *EventModel) GetData() []byte {
	return e.Data
}

// SetData add the data attached to the Event serialized to bytes.
func (e *EventModel) SetData(data []byte) *EventModel {
	e.Data = data
	return e
}

// GetJsonData json unmarshal data attached to the Event.
func (e *EventModel) GetJsonData(data interface{}) error {
	return json.Unmarshal(e.GetData(), data)
}

// SetJsonData serialize to json and set data attached to the Event.
func (e *EventModel) SetJsonData(data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	e.Data = dataBytes
	return nil
}

// GetEventType returns the EventType of the event.
func (e *EventModel) GetEventType() string {
	return e.EventType
}

// GetAggregateType is the AggregateType that the Event can be applied to.
func (e *EventModel) GetAggregateType() aggregate.AggregateType {
	return e.AggregateType
}

// SetAggregateType set the AggregateType that the Event can be applied to.
func (e *EventModel) SetAggregateType(aggregateType aggregate.AggregateType) {
	e.AggregateType = aggregateType
}

// GetAggregateId is the Id of the Aggregate that the Event belongs to
func (e *EventModel) GetAggregateId() string {
	return e.AggregateID
}

// GetVersion is the version of the Aggregate after the Event has been applied.
func (e *EventModel) GetVersion() int64 {
	return e.Version
}

// SetVersion set the version of the Aggregate.
func (e *EventModel) SetVersion(aggregateVersion int64) {
	e.Version = aggregateVersion
}

// GetMetadata is app-specific metadata such as request Id, originating user etc.
func (e *EventModel) GetMetadata() []byte {
	return e.Metadata
}

// SetMetadata add app-specific metadata serialized as json for the Event.
func (e *EventModel) SetMetadata(metaData interface{}) error {

	metaDataBytes, err := json.Marshal(metaData)
	if err != nil {
		return err
	}

	e.Metadata = metaDataBytes
	return nil
}

// GetJsonMetadata unmarshal app-specific metadata serialized as json for the Event.
func (e *EventModel) GetJsonMetadata(metaData interface{}) error {
	return json.Unmarshal(e.GetMetadata(), metaData)
}
