package services

import (
	"context"
	"errors"

	"github.com/suprt/planica_bi/backend/internal/models"
)

// CounterService handles business logic for Yandex counters
type CounterService struct {
	counterRepo CounterRepositoryInterface
}

// NewCounterService creates a new counter service
func NewCounterService(counterRepo CounterRepositoryInterface) *CounterService {
	return &CounterService{
		counterRepo: counterRepo,
	}
}

// CreateCounter creates a new counter
func (s *CounterService) CreateCounter(ctx context.Context, counter *models.YandexCounter) error {
	// Validate required fields
	if counter.ProjectID == 0 {
		return errors.New("project_id is required")
	}
	if counter.CounterID == 0 {
		return errors.New("counter_id is required")
	}

	// Check if counter with this CounterID already exists
	existing, err := s.counterRepo.GetByCounterID(ctx, counter.CounterID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("counter with this CounterID already exists")
	}

	// If this is set as primary, unset other primary counters for this project
	if counter.IsPrimary {
		if err := s.unsetPrimaryCounters(ctx, counter.ProjectID); err != nil {
			return err
		}
	}

	return s.counterRepo.Create(ctx, counter)
}

// GetCountersByProject retrieves all counters for a project
func (s *CounterService) GetCountersByProject(ctx context.Context, projectID uint) ([]*models.YandexCounter, error) {
	return s.counterRepo.GetByProjectID(ctx, projectID)
}

// unsetPrimaryCounters unsets all primary flags for counters in a project
func (s *CounterService) unsetPrimaryCounters(ctx context.Context, projectID uint) error {
	counters, err := s.counterRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return err
	}

	for _, counter := range counters {
		if counter.IsPrimary {
			counter.IsPrimary = false
			if err := s.counterRepo.Update(ctx, counter); err != nil {
				return err
			}
		}
	}

	return nil
}
