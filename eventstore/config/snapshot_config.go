package config

type Config struct {
	SnapshotFrequency int64 `json:"snapshotFrequency" validate:"required,gte=0"`
}
