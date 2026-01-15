package model

type Job struct {
	ID        string
	Status    string
	Type      string
	RawURL    string
	ResultURL string
	Error     string

	Params TransformParams

	QueuedAt    int64
	StartedAt   int64
	CompletedAt int64
}
