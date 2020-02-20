package task

import (
	"context"
	"fmt"

	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/uuid"
	"github.com/gotasma/internal/pkg/validator"
	"github.com/sirupsen/logrus"
)

type (
	mongoRepository interface {
		Create(context.Context, *types.Task) error
		// Delete(ctx context.Context, id string) error
		Update(ctx context.Context, projectID string, req *types.Task) error

		FindByProjectID(ctx context.Context, projectID string) ([]*types.Task, error)
		// FindByID(ctx context.Context, id string) (*types.Task, error)
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

func (s *Service) FindByProjectID(ctx context.Context, projectID string) ([]*types.Task, error) {

	tasks, err := s.repo.FindByProjectID(ctx, projectID)

	return tasks, err
}

func (s *Service) Create(ctx context.Context, projectID string, req *types.Task) (*types.Task, error) {

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to create tasks due to invalid req, %w", err)
		return nil, err
	}

	//Create new project
	tasks := &types.Task{
		ProjectID:        projectID,
		AllChildren:      req.AllChildren,
		Children:         req.Children,
		Duration:         req.Duration,
		Effort:           req.Duration,
		End:              req.Duration,
		EstimateDuration: req.EstimateDuration,
		Label:            req.Label,
		Parent:           req.Parent,
		Parents:          req.Parents,
		Start:            req.Start,
		TaskID:           uuid.New(),
		Type:             req.Type,
	}

	if err := s.repo.Create(ctx, tasks); err != nil {
		logrus.Errorf("failed to create project: %v", err)
		return nil, fmt.Errorf("failed to create project, %w", err)
	}

	return nil, nil
}

func (s *Service) Update(ctx context.Context, projectID string, req *types.Task) error {
	//Validate [input]
	if err := validator.Validate(req); err != nil {
		logrus.Error("Failed to validate input update task %v", err)
		return err
	}

	return s.repo.Update(ctx, projectID, req)
}
