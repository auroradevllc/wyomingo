package wyoming

const (
	TimerStartedEvent   EventType = "timer-started"
	TimerUpdatedEvent   EventType = "timer-updated"
	TimerCancelledEvent EventType = "timer-cancelled"
	TimerFinishedEvent  EventType = "timer-finished"
)

type TimerStarted struct {
	ID           string `json:"id"`
	TotalSeconds int    `json:"total_seconds"`
	Name         string `json:"name,omitempty"`
	StartHours   int    `json:"start_hours,omitempty"`
	StartMinutes int    `json:"start_minutes,omitempty"`
	StartSeconds int    `json:"start_seconds,omitempty"`
}

type TimerUpdated struct {
	ID           string `json:"id"`
	Active       bool   `json:"is_active"`
	TotalSeconds int    `json:"total_seconds"`
}

type TimerCancelled struct {
	ID string `json:"id"`
}

type TimerFinished struct {
	ID string `json:"id"`
}
