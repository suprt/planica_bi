package services

import (
	"errors"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"gorm.io/gorm"
)

// GoalService handles business logic for goals
type GoalService struct {
	goalRepo    *repositories.GoalRepository
	counterRepo *repositories.CounterRepository
}

// NewGoalService creates a new goal service
func NewGoalService(goalRepo *repositories.GoalRepository, counterRepo *repositories.CounterRepository) *GoalService {
	return &GoalService{
		goalRepo:    goalRepo,
		counterRepo: counterRepo,
	}
}

// CreateGoal creates a new goal
func (s *GoalService) CreateGoal(goal *models.Goal) error {
	// Validate that counter exists
	_, err := s.counterRepo.GetByID(goal.CounterID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("counter not found")
		}
		return err
	}

	// Check if goal with this GoalID already exists for this counter
	existing, err := s.goalRepo.GetByGoalID(goal.CounterID, goal.GoalID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existing != nil {
		return errors.New("goal with this GoalID already exists for this counter")
	}

	return s.goalRepo.Create(goal)
}

// GetGoal retrieves a goal by ID
func (s *GoalService) GetGoal(id uint) (*models.Goal, error) {
	return s.goalRepo.GetByID(id)
}

// GetGoalsByCounter retrieves all goals for a counter
func (s *GoalService) GetGoalsByCounter(counterID uint) ([]*models.Goal, error) {
	// Validate that counter exists
	_, err := s.counterRepo.GetByID(counterID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("counter not found")
		}
		return nil, err
	}

	return s.goalRepo.GetByCounterID(counterID)
}

// GetGoalByGoalID retrieves a goal by Yandex goal ID and counter ID
func (s *GoalService) GetGoalByGoalID(counterID uint, goalID int64) (*models.Goal, error) {
	return s.goalRepo.GetByGoalID(counterID, goalID)
}

// UpdateGoal updates a goal
func (s *GoalService) UpdateGoal(goal *models.Goal) error {
	// Validate that counter exists if CounterID changed
	if goal.CounterID != 0 {
		_, err := s.counterRepo.GetByID(goal.CounterID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("counter not found")
			}
			return err
		}
	}

	return s.goalRepo.Update(goal)
}

// DeleteGoal deletes a goal
func (s *GoalService) DeleteGoal(id uint) error {
	return s.goalRepo.Delete(id)
}

// DeleteGoalsByCounter deletes all goals for a counter
func (s *GoalService) DeleteGoalsByCounter(counterID uint) error {
	// Validate that counter exists
	_, err := s.counterRepo.GetByID(counterID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("counter not found")
		}
		return err
	}

	return s.goalRepo.DeleteByCounterID(counterID)
}

// GetConversionGoals retrieves only conversion goals for a counter
func (s *GoalService) GetConversionGoals(counterID uint) ([]*models.Goal, error) {
	goals, err := s.GetGoalsByCounter(counterID)
	if err != nil {
		return nil, err
	}

	var conversionGoals []*models.Goal
	for _, goal := range goals {
		if goal.IsConversion {
			conversionGoals = append(conversionGoals, goal)
		}
	}

	return conversionGoals, nil
}

