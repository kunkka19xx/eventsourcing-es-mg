package aggregate

import "encoding/json"

type Snapshot struct {
	ID      string
	Type    AggregateType
	State   []byte
	Version int64
}

// NewSnapshotFromAggregate create snapshot from the Aggregate to hold the current state
func NewSnapshotFromAggregate(aggregate Aggregate) (*Snapshot, error) {
	aggregateData, err := json.Marshal(aggregate)
	if err != nil {
		return nil, err
	}
	return &Snapshot{
		ID:      aggregate.GetID(),
		Type:    aggregate.GetType(),
		State:   aggregateData,
		Version: int64(aggregate.GetVersion()),
	}, nil
}
