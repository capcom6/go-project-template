package agent

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrAgentStopped = errors.New("agent stopped")
)
