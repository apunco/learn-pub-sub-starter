package enums

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)
