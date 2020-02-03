package project

import (
	"context"
	"fmt"

	"praslar.com/gotasma/internal/app/auth"
	"praslar.com/gotasma/internal/app/status"
	"praslar.com/gotasma/internal/app/types"
	"praslar.com/gotasma/internal/pkg/db"
	"praslar.com/gotasma/internal/pkg/uuid"
	"praslar.com/gotasma/internal/pkg/validator"
)

type (
	repository interface {
		FindByName(ctx context.Context, name string, createrID string) (*types.Project, error)
		Create(context.Context, *types.Project) error
	}

	PolicyService interface {
		Validate(ctx context.Context, obj string, act string) error
	}

	Service struct {
		repo   repository
		policy PolicyService
	}
)

func New(repo repository, policy PolicyService) *Service {
	return &Service{
		repo:   repo,
		policy: policy,
	}
}

func (s *Service) Create(ctx context.Context, req *types.CreateProjectRequest) (*types.Project, error) {

	pm := auth.FromContext(ctx)

	//only PM can create Project
	if err := s.policy.Validate(ctx, types.ObjectProject, types.ActionProject); err != nil {
		return nil, err
	}

	if err := validator.Validate(req); err != nil {
		return nil, err
	}

	existingProject, err := s.repo.FindByName(ctx, req.Name, pm.UserID)

	if err != nil && !db.IsErrNotFound(err) {
		return nil, fmt.Errorf("failed to check existing project by name: %w", err)
	}

	if existingProject != nil {
		return nil, status.Project().DuplicateProject
	}

	project := &types.Project{
		CreaterID: pm.UserID,
		Highlight: true,
		Name:      req.Name,
		ProjectID: uuid.New(),
		DevsID:    []string{},
		Tasks:     []types.Task{},
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project, %w", err)
	}

	return project, nil
}

func (s *Service) FindAll(context.Context) []types.Project {

}
