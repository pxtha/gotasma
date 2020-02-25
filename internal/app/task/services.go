package task

import (
	"context"
	"fmt"

	"github.com/gotasma/internal/app/auth"
	"github.com/gotasma/internal/app/status"
	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/uuid"
	"github.com/gotasma/internal/pkg/validator"

	"github.com/sirupsen/logrus"
)

type (
	mongoRepository interface {
		Create(context.Context, *types.Task) error
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, projectID string, req *types.Task) error

		FindByProjectID(ctx context.Context, projectID string) ([]*types.Task, error)
		FindByID(ctx context.Context, id string) (*types.Task, error)
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
func (s *Service) FindByIDs(ctx context.Context, ids []string) ([]*types.TaskInfo, error) {

	tasks := make([]*types.TaskInfo, 0)
	for _, id := range ids {
		task, err := s.repo.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &types.TaskInfo{
			TaskID: task.TaskID,
			Label:  task.Label,
			End:    task.End,
			Start:  task.Start,
		})
	}
	return tasks, nil
}

func (s *Service) FindByID(ctx context.Context, id string) (*types.Task, error) {

	task, err := s.repo.FindByID(ctx, id)

	return task, err
}

func (s *Service) Create(ctx context.Context, projectID string, req *types.Task) (*types.Task, error) {

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to create tasks due to invalid req, %w", err)
		return nil, fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
	}
	user := auth.FromContext(ctx)
	// Create new task
	tasks := &types.Task{
		ProjectID:        projectID,
		CreaterID:        user.UserID,
		AllChildren:      req.AllChildren,
		Children:         req.Children,
		Duration:         req.Duration,
		Effort:           req.Effort,
		End:              req.End,
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

	return tasks, nil
}

func (s *Service) Update(ctx context.Context, projectID string, req *types.Task) error {
	// Validate [input]
	if err := validator.Validate(req); err != nil {
		logrus.Error("Failed to validate input update task %v", err)
		return fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
	}

	return s.repo.Update(ctx, projectID, req)
}

func (s *Service) Delete(ctx context.Context, id string) error {

	if err := s.repo.Delete(ctx, id); err != nil {
		logrus.Errorf("Fail to delete task due to %v", err)
		return err
	}

	return nil
}
