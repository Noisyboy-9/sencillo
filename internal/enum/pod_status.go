package enum

type PodStatus int

const (
	PodStatusPendding PodStatus = iota
	PodStatusRunning
	PodStatusSucceeded
	PodStatusFailed
	PodStatusUnkown
)
