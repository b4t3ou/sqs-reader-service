package domain

import "time"

const (
	DefaultAWSRegion = "eu-west-1"
)

type Event struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
