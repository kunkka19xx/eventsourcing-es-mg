package command

type BaseCommand struct {
	AggregateID string `json:"aggregateId" validate:"required,min=1"`
}

type Command interface {
	GetAggregateID() string
}

func NewBaseCommand(aggregateID string) BaseCommand {
	return BaseCommand{AggregateID: aggregateID}
}

func (c *BaseCommand) GetAggregateID() string {
	return c.AggregateID
}
