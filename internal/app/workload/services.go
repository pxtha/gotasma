package workload

import (
	"context"
	"fmt"

	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/uuid"

	"github.com/sirupsen/logrus"
)

type (
	mongoRepository interface {
		FindByID(ctx context.Context, projectID string, userID string) (*types.WorkLoad, error)

		Create(context.Context, *types.WorkLoad) error
		Delete(ctx context.Context, projectID string, userID string) error
		Update(ctx context.Context, projectID string, userID string, overload map[int]int) error
	}

	PolicyService interface {
		Validate(ctx context.Context, obj string, act string) error
	}

	Service struct {
		repo   mongoRepository
		policy PolicyService
	}
)

func New(repo mongoRepository, policy PolicyService) *Service {
	return &Service{
		repo:   repo,
		policy: policy,
	}
}

func (s *Service) FindByID(ctx context.Context, projectID string, userID string) (*types.WorkLoad, error) {

	workload, err := s.repo.FindByID(ctx, projectID, userID)

	return workload, err
}

func (s *Service) Create(ctx context.Context, projectID string, userID string) error {

	// Create new project
	workload := &types.WorkLoad{
		ProjectID:  projectID,
		UserID:     userID,
		WorkLoadID: uuid.New(),
	}

	if err := s.repo.Create(ctx, workload); err != nil {
		logrus.Errorf("failed to create workload: %v", err)
		return fmt.Errorf("failed to create workload project, %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, projectID string, userID string, overload map[int]int) error {

	if err := s.repo.Update(ctx, projectID, userID, overload); err != nil {
		logrus.Errorf("Fail to update workload due to %v", err)
		return fmt.Errorf("Fail to update workload, %w", err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, projectID string, userID string) error {

	if err := s.repo.Delete(ctx, projectID, userID); err != nil {
		logrus.Errorf("Fail to workload due to %v", err)
		return fmt.Errorf("Failed to delete workload project, %w", err)
	}

	return nil
}
