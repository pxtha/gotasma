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
	"github.com/sirupsen/logrus"
)

type (
	repository interface {
		FindByName(ctx context.Context, name string, createrID string) (*types.Project, error)
		FindByProjectID(ctx context.Context, projectID string) (*types.Project, error)
		FindAllByUserID(ctx context.Context, id string, role types.Role) ([]*types.Project, error)
		Create(context.Context, *types.Project) error
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, req *types.ProjectInfo) error
		//Update action: addToSet or Pull
		UpdateDevsID(ctx context.Context, devsID []string, projectID string, addToSet bool) error
	}

	PolicyService interface {
		Validate(ctx context.Context, obj string, act string) error
	}

	UserService interface {
		//User can be DEV or PM
		CheckUsersExist(ctx context.Context, userID string) (string, error)
		GetDevsInfo(ctx context.Context, userIDs []string) ([]*types.UserInfo, error)
	}

	Service struct {
		repo   repository
		policy PolicyService
		user   UserService
	}
)

func New(repo repository, policy PolicyService, updateUser UserService) *Service {
	return &Service{
		repo:   repo,
		policy: policy,
		user:   updateUser,
	}
}

func (s *Service) Update(ctx context.Context, id string, req *types.ProjectInfo) (*types.Project, error) {

	//only PM can create Project
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}
	//validate incoming data
	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to update project due to invalid req, %w", err)
		return nil, err
	}

	project, err := s.repo.FindByProjectID(ctx, id)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by id: %v", err)
		return nil, fmt.Errorf("failed to check existing project by id: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return nil, status.Project().NotFoundProject
	}

	//validate each task
	for _, task := range req.Tasks {
		if err := validator.Validate(task); err != nil {
			logrus.Errorf("Fail to update project due to invalid req, %w", err)
			return nil, err
		}
		for _, devInTask := range task.DevsID {
			devIsInProject := false
			for _, dev := range project.DevsID {
				if devInTask == dev {
					devIsInProject = true
				}
			}
			if !devIsInProject {
				logrus.Errorf("Devs %v in tasks %v not found ", devInTask, task.TaskID)
				return nil, status.Project().NotFoundDev
			}
		}
	}
	//Just update tasks field of project
	//VueJS will do all the validation for data comes in
	//So that user can edit tasks as long as they want, only save tasks when user use update task API

	if err = s.repo.Update(ctx, id, req); err != nil {
		logrus.Errorf("failed to update project: %v", err)
		return nil, fmt.Errorf("failed to update project, %w", err)
	}

	//TODO add new history elasticsearch
	//TODO calculate workload

	return project, nil
}

func (s *Service) Create(ctx context.Context, req *types.CreateProjectRequest) (*types.Project, error) {

	pm := auth.FromContext(ctx)

	//only PM can create Project
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to create project due to invalid req, %w", err)
		return nil, err
	}

	existingProject, err := s.repo.FindByName(ctx, req.Name, pm.UserID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by name: %v", err)
		return nil, fmt.Errorf("failed to check existing project by name: %w", err)
	}

	if existingProject != nil {
		logrus.Errorf("Project already exsit")
		return nil, status.Project().DuplicateProject
	}

	project := &types.Project{
		CreaterID: pm.UserID,
		Highlight: true,
		Name:      req.Name,
		ProjectID: uuid.New(),
		Desc:      req.Desc,
	}

	if err = s.repo.Create(ctx, project); err != nil {
		logrus.Errorf("failed to create project: %v", err)
		return nil, fmt.Errorf("failed to create project, %w", err)
	}

	return project, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		logrus.Errorf("Fail to delete project due to %v", err)
		return status.Project().NotFoundProject
	}
	return nil
}

func (s *Service) FindByID(ctx context.Context, id string) (*types.Project, error) {

	project, err := s.repo.FindByProjectID(ctx, id)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by ID: %v", err)
		return nil, fmt.Errorf("failed to check existing project by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return nil, status.Project().NotFoundProject
	}

	//Remove devIds from tasks if dev no longer in project
	taskInfo := make([]types.Task, 0)

	for _, task := range project.Tasks {

		for i, inTask := range task.DevsID {
			check := false
			for _, inProject := range project.DevsID {
				if inTask == inProject {
					check = true
				}
			}
			if check == false {
				task.DevsID[i] = task.DevsID[len(task.DevsID)-1]
				task.DevsID = task.DevsID[:len(task.DevsID)-1]
			}
		}
		taskInfo = append(taskInfo, types.Task{
			Label:            task.Label,
			AllChildren:      task.AllChildren,
			Children:         task.Children,
			DevsID:           task.DevsID,
			Duration:         task.Duration,
			Effort:           task.Effort,
			End:              task.End,
			EstimateDuration: task.EstimateDuration,
			Parent:           task.Parent,
			Parents:          task.Parents,
			Start:            task.Start,
			TaskID:           task.TaskID,
			Type:             task.Type,
		})
	}

	var info *types.Project
	info = &types.Project{
		Name:      project.Name,
		CreatedAt: project.CreatedAt,
		CreaterID: project.CreaterID,
		DevsID:    project.DevsID,
		Highlight: project.Highlight,
		ProjectID: project.ProjectID,
		Tasks:     taskInfo,
		UpdateAt:  project.UpdateAt,
		Desc:      project.Desc,
	}

	return info, nil
}

func (s *Service) FindAllProjects(ctx context.Context) ([]*types.Project, error) {

	user := auth.FromContext(ctx)
	logrus.Infof("Devs id: %v", user.UserID)

	projects, err := s.repo.FindAllByUserID(ctx, user.UserID, user.Role)

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

func (s *Service) AddDevs(ctx context.Context, userIDs []string, projectID string) ([]string, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}

	// Check project exist
	project, err := s.repo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by ID: %v", err)
		return nil, fmt.Errorf("failed to check existing project by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return nil, status.Project().NotFoundProject
	}

	for _, userID := range userIDs {
		_, err = s.user.CheckUsersExist(ctx, userID)
		if err != nil {
			logrus.Errorf("Add Dev: User %v", err)
			return nil, status.User().NotFoundUser
		}
	}

	//check if user already in project
	for _, userID := range userIDs {

		for _, devID := range project.DevsID {
			if userID == devID {
				return nil, status.Project().AlreadyInProject
			}
		}
		// cannot add PM to project
		if userID == project.CreaterID {
			return nil, status.Sercurity().InvalidAction
		}
	}

	if err := s.repo.UpdateDevsID(ctx, userIDs, projectID, true); err != nil {
		logrus.Errorf("Failed to update devs ID in project info %v", err)
		return nil, err
	}

	return userIDs, nil
}

func (s *Service) RemoveDev(ctx context.Context, userID string, projectID string) error {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return err
	}
	// Check project exist
	project, err := s.repo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		return fmt.Errorf("failed to check existing project by ID: %w", err)
	}
	if db.IsErrNotFound(err) {
		return status.Project().NotFoundProject
	}

	_, err = s.user.CheckUsersExist(ctx, userID)
	if err != nil {
		logrus.Errorf("Remove Dev: User %v", err)
		return status.User().NotFoundUser
	}
	//check if user have not in project yet
	NotInProject := true
	for _, devID := range project.DevsID {
		if userID == devID {
			NotInProject = false
		}
	}
	if NotInProject {
		return status.Project().NotInProject
	}

	if err := s.repo.UpdateDevsID(ctx, []string{userID}, projectID, false); err != nil {
		logrus.Errorf("Failed to update devs ID in project info %v", err)
		return err
	}

	return nil
}

func (s *Service) FindAllDevs(ctx context.Context, projectID string) ([]*types.UserInfo, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}
	// Check project exist
	project, err := s.repo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by ID: %v", err)
		return nil, fmt.Errorf("failed to check existing project by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return nil, status.Project().NotFoundProject
	}

	//TODO remove devsid
	for _, user := range project.DevsID {
		devNotFound, err := s.user.CheckUsersExist(ctx, user)
		if err != nil {
			logrus.Infof("Removing Dev IDs from project INFO %v", devNotFound)
			if err := s.repo.UpdateDevsID(ctx, []string{devNotFound}, projectID, false); err != nil {
				logrus.Errorf("Failed to update devs ID in project info %v", err)
				return nil, err
			}
		}
	}

	var info []*types.UserInfo

	info, err = s.user.GetDevsInfo(ctx, project.DevsID)

	return info, nil
}
