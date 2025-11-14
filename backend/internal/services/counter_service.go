package services

import (
	"errors"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"gorm.io/gorm"
)

// CounterService handles business logic for Yandex counters
type CounterService struct {
	counterRepo *repositories.CounterRepository
}

// NewCounterService creates a new counter service
func NewCounterService(counterRepo *repositories.CounterRepository) *CounterService {
	return &CounterService{
		counterRepo: counterRepo,
	}
}

// CreateCounter creates a new counter
func (s *CounterService) CreateCounter(counter *models.YandexCounter) error {
	// Check if counter with this CounterID already exists
	existing, err := s.counterRepo.GetByCounterID(counter.CounterID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existing != nil {
		return errors.New("counter with this CounterID already exists")
	}

	// If this is set as primary, unset other primary counters for this project
	if counter.IsPrimary {
		if err := s.unsetPrimaryCounters(counter.ProjectID); err != nil {
			return err
		}
	}

	return s.counterRepo.Create(counter)
}

// GetCounter retrieves a counter by ID
func (s *CounterService) GetCounter(id uint) (*models.YandexCounter, error) {
	return s.counterRepo.GetByID(id)
}

// GetCountersByProject retrieves all counters for a project
func (s *CounterService) GetCountersByProject(projectID uint) ([]*models.YandexCounter, error) {
	return s.counterRepo.GetByProjectID(projectID)
}

// GetCounterByCounterID retrieves a counter by Yandex counter ID
func (s *CounterService) GetCounterByCounterID(counterID int64) (*models.YandexCounter, error) {
	return s.counterRepo.GetByCounterID(counterID)
}

// GetPrimaryCounter retrieves the primary counter for a project
func (s *CounterService) GetPrimaryCounter(projectID uint) (*models.YandexCounter, error) {
	return s.counterRepo.GetPrimaryByProjectID(projectID)
}

// UpdateCounter updates a counter
func (s *CounterService) UpdateCounter(counter *models.YandexCounter) error {
	// If this is set as primary, unset other primary counters for this project
	if counter.IsPrimary {
		if err := s.unsetPrimaryCounters(counter.ProjectID); err != nil {
			return err
		}
	}

	return s.counterRepo.Update(counter)
}

// SetAsPrimary sets a counter as primary and unsets others
func (s *CounterService) SetAsPrimary(id uint) error {
	counter, err := s.counterRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Unset other primary counters for this project
	if err := s.unsetPrimaryCounters(counter.ProjectID); err != nil {
		return err
	}

	// Set this counter as primary
	counter.IsPrimary = true
	return s.counterRepo.Update(counter)
}

// DeleteCounter deletes a counter
func (s *CounterService) DeleteCounter(id uint) error {
	return s.counterRepo.Delete(id)
}

// DeleteCountersByProject deletes all counters for a project
func (s *CounterService) DeleteCountersByProject(projectID uint) error {
	return s.counterRepo.DeleteByProjectID(projectID)
}

// unsetPrimaryCounters unsets all primary flags for counters in a project
func (s *CounterService) unsetPrimaryCounters(projectID uint) error {
	counters, err := s.counterRepo.GetByProjectID(projectID)
	if err != nil {
		return err
	}

	for _, counter := range counters {
		if counter.IsPrimary {
			counter.IsPrimary = false
			if err := s.counterRepo.Update(counter); err != nil {
				return err
			}
		}
	}

	return nil
}

