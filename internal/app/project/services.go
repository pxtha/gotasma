package project

import (
	"context"
	"fmt"

	"github.com/gotasma/internal/app/auth"
	"github.com/gotasma/internal/app/status"
	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/db"
	"github.com/gotasma/internal/pkg/uuid"
	"github.com/gotasma/internal/pkg/validator"
)

type (
	repository interface {
		FindByName(ctx context.Context, name string, createrID string) (*types.Project, error)
		FindByPm(ctx context.Context, pmID string) ([]*types.Project, error)
		FindByDev(ctx context.Context, devID string) ([]*types.Project, error)
		Create(context.Context, *types.Project) (string, error)
	}

	PolicyService interface {
		Validate(ctx context.Context, obj string, act string) error
	}

	UpdateUserInfo interface {
		UpdateUserProjectsID(ctx context.Context, userID string, projectID string) error
	}

	Service struct {
		repo       repository
		policy     PolicyService
		updateUser UpdateUserInfo
	}
)

func New(repo repository, policy PolicyService, updateUser UpdateUserInfo) *Service {
	return &Service{
		repo:       repo,
		policy:     policy,
		updateUser: updateUser,
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

	projectID, err := s.repo.Create(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to create project, %w", err)
	}
	//update ProjectIDs of PM info
	if err := s.updateUser.UpdateUserProjectsID(ctx, pm.UserID, projectID); err != nil {
		return nil, fmt.Errorf("failed to update projectIDs of user, %w", err)
	}
	return project, nil
}

func (s *Service) FindAll(ctx context.Context) ([]*types.Project, error) {
	user := auth.FromContext(ctx)

	projects, err := s.repo.FindByPm(ctx, user.UserID)
	if user.Role != types.PM {
		projects, err = s.repo.FindByDev(ctx, user.UserID)
	}

	info := make([]*types.Project, 0)
	for _, project := range projects {
		info = append(info, &types.Project{
			ProjectID: project.ProjectID,
			Name:      project.Name,
			CreaterID: project.CreaterID,
			DevsID:    project.DevsID,
			Highlight: project.Highlight,
			UpdateAt:  project.UpdateAt,
		})
	}

	return info, err
}