package queue

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// Task type names
const (
	TypeSyncMetrica     = "sync:metrica"
	TypeSyncDirect      = "sync:direct"
	TypeSyncProject     = "sync:project"
	TypeAnalyzeMetrics  = "analyze:metrics"
	TypeGenerateReport  = "generate:report"
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

// AnalyzeMetricsPayload is the payload for metrics analysis task
type AnalyzeMetricsPayload struct {
	ProjectID uint     `json:"project_id"`
	Periods   []string `json:"periods"`
}

// GenerateReportPayload is the payload for report generation task
type GenerateReportPayload struct {
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

// NewAnalyzeMetricsTask creates a new metrics analysis task
func NewAnalyzeMetricsTask(projectID uint, periods []string) *asynq.Task {
	payload := AnalyzeMetricsPayload{
		ProjectID: projectID,
		Periods:   periods,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal payload: %v", err))
	}
	return asynq.NewTask(TypeAnalyzeMetrics, payloadBytes)
}

// ParseAnalyzeMetricsPayload parses metrics analysis task payload
func ParseAnalyzeMetricsPayload(task *asynq.Task) (*AnalyzeMetricsPayload, error) {
	var payload AnalyzeMetricsPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return &payload, nil
}

// NewGenerateReportTask creates a new report generation task
func NewGenerateReportTask(projectID uint) *asynq.Task {
	payload := GenerateReportPayload{
		ProjectID: projectID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal payload: %v", err))
	}
	return asynq.NewTask(TypeGenerateReport, payloadBytes)
}

// ParseGenerateReportPayload parses report generation task payload
func ParseGenerateReportPayload(task *asynq.Task) (*GenerateReportPayload, error) {
	var payload GenerateReportPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return &payload, nil
}

