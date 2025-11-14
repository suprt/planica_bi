package services

import (
	"context"
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
func (s *GoalService) CreateGoal(ctx context.Context, goal *models.Goal) error {
	// Validate required fields
	if goal.CounterID == 0 {
		return errors.New("counter_id is required")
	}
	if goal.GoalID == 0 {
		return errors.New("goal_id is required")
	}

	// Validate that counter exists
	_, err := s.counterRepo.GetByID(ctx, goal.CounterID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("counter not found")
		}
		return err
	}

	// Check if goal with this GoalID already exists for this counter
	existing, err := s.goalRepo.GetByGoalID(ctx, goal.CounterID, goal.GoalID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existing != nil {
		return errors.New("goal with this GoalID already exists for this counter")
	}

	return s.goalRepo.Create(ctx, goal)
}

// GetGoal retrieves a goal by ID
func (s *GoalService) GetGoal(ctx context.Context, id uint) (*models.Goal, error) {
	return s.goalRepo.GetByID(ctx, id)
}

// GetGoalsByCounter retrieves all goals for a counter
func (s *GoalService) GetGoalsByCounter(ctx context.Context, counterID uint) ([]*models.Goal, error) {
	// Validate that counter exists
	_, err := s.counterRepo.GetByID(ctx, counterID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("counter not found")
		}
		return nil, err
	}

	return s.goalRepo.GetByCounterID(ctx, counterID)
}

// GetGoalsByProject retrieves all goals for all counters of a project
func (s *GoalService) GetGoalsByProject(ctx context.Context, projectID uint) ([]*models.Goal, error) {
	// Get all counters for the project
	counters, err := s.counterRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if len(counters) == 0 {
		return []*models.Goal{}, nil
	}

	// Get all counter IDs
	var counterIDs []uint
	for _, counter := range counters {
		counterIDs = append(counterIDs, counter.ID)
	}

	// Get all goals for these counters
	return s.goalRepo.GetByCounterIDs(ctx, counterIDs)
}

// GetGoalByGoalID retrieves a goal by Yandex goal ID and counter ID
func (s *GoalService) GetGoalByGoalID(ctx context.Context, counterID uint, goalID int64) (*models.Goal, error) {
	return s.goalRepo.GetByGoalID(ctx, counterID, goalID)
}

// UpdateGoal updates a goal
func (s *GoalService) UpdateGoal(ctx context.Context, goal *models.Goal) error {
	// Validate that counter exists if CounterID changed
	if goal.CounterID != 0 {
		_, err := s.counterRepo.GetByID(ctx, goal.CounterID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("counter not found")
			}
			return err
		}
	}

	return s.goalRepo.Update(ctx, goal)
}

// DeleteGoal deletes a goal
func (s *GoalService) DeleteGoal(ctx context.Context, id uint) error {
	return s.goalRepo.Delete(ctx, id)
}

// DeleteGoalsByCounter deletes all goals for a counter
func (s *GoalService) DeleteGoalsByCounter(ctx context.Context, counterID uint) error {
	// Validate that counter exists
	_, err := s.counterRepo.GetByID(ctx, counterID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("counter not found")
		}
		return err
	}

	return s.goalRepo.DeleteByCounterID(ctx, counterID)
}

// GetConversionGoals retrieves only conversion goals for a counter
func (s *GoalService) GetConversionGoals(ctx context.Context, counterID uint) ([]*models.Goal, error) {
	goals, err := s.GetGoalsByCounter(ctx, counterID)
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
