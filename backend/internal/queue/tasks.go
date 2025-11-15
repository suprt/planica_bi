package queue

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// Task type names
const (
	TypeSyncMetrica = "sync:metrica"
	TypeSyncDirect  = "sync:direct"
	TypeSyncProject = "sync:project"
)

// SyncMetricaPayload is the payload for Metrica sync task
type SyncMetricaPayload struct {
	ProjectID uint `json:"project_id"`
	Year      int  `json:"year"`
	Month     int  `json:"month"`
}

// SyncDirectPayload is the payload for Direct sync task
type SyncDirectPayload struct {
	ProjectID uint `json:"project_id"`
	Year      int  `json:"year"`
	Month     int  `json:"month"`
}

// SyncProjectPayload is the payload for project sync task
type SyncProjectPayload struct {
	ProjectID uint `json:"project_id"`
}

// NewSyncMetricaTask creates a new Metrica sync task
func NewSyncMetricaTask(projectID uint, year, month int) *asynq.Task {
	payload := SyncMetricaPayload{
		ProjectID: projectID,
		Year:      year,
		Month:     month,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal payload: %v", err))
	}
	return asynq.NewTask(TypeSyncMetrica, payloadBytes)
}

// NewSyncDirectTask creates a new Direct sync task
func NewSyncDirectTask(projectID uint, year, month int) *asynq.Task {
	payload := SyncDirectPayload{
		ProjectID: projectID,
		Year:      year,
		Month:     month,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal payload: %v", err))
	}
	return asynq.NewTask(TypeSyncDirect, payloadBytes)
}

// NewSyncProjectTask creates a new project sync task
func NewSyncProjectTask(projectID uint) *asynq.Task {
	payload := SyncProjectPayload{
		ProjectID: projectID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal payload: %v", err))
	}
	return asynq.NewTask(TypeSyncProject, payloadBytes)
}

// ParseSyncMetricaPayload parses Metrica sync task payload
func ParseSyncMetricaPayload(task *asynq.Task) (*SyncMetricaPayload, error) {
	var payload SyncMetricaPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return &payload, nil
}

// ParseSyncDirectPayload parses Direct sync task payload
func ParseSyncDirectPayload(task *asynq.Task) (*SyncDirectPayload, error) {
	var payload SyncDirectPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return &payload, nil
}

// ParseSyncProjectPayload parses project sync task payload
func ParseSyncProjectPayload(task *asynq.Task) (*SyncProjectPayload, error) {
	var payload SyncProjectPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return &payload, nil
}

