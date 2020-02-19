package project

import (
	"context"
	"fmt"
	"time"

	"github.com/gotasma/internal/app/auth"
	"github.com/gotasma/internal/app/status"
	"github.com/gotasma/internal/app/types"
	"github.com/gotasma/internal/pkg/db"
	"github.com/gotasma/internal/pkg/uuid"
	"github.com/gotasma/internal/pkg/validator"

	"github.com/sirupsen/logrus"
)

type (
	mongoRepository interface {
		//Check project exists
		FindByName(ctx context.Context, name string, createrID string) (*types.Project, error)
		//View project details (inclued tasks)
		FindByProjectID(ctx context.Context, projectID string) (*types.Project, error)
		//Search all project belong to user (Just projectinfo not include tasks but has number of tasks: int)
		//If Pm searchs: return all project has creater_id == pmID
		//If Devs search: return all project has devs_id == devid
		FindAllByUserID(ctx context.Context, id string, role types.Role) ([]*types.Project, error)
		//Find all project has holiday id
		FindAllByHolidaysID(ctx context.Context, holidayID string) ([]*types.Project, error)
		//Find all tasks has user id
		Create(context.Context, *types.Project) error
		Delete(ctx context.Context, id string) error
		//Trusted data get from client
		//ID must be uuid
		//Update project API inclued (add, delete, update tasks)
		Save(ctx context.Context, id string, req *types.SaveProject) error
		Update(ctx context.Context, id string, req *types.UpdateProject) error
		//UpdateDevsID action: addToSet or Pull
		UpdateDevsID(ctx context.Context, devsID []string, projectID string, addToSet bool) error
		UpdateHolidaysID(ctx context.Context, holiday string, projectID string) error
		//FindholidayByID
	}

	elasticRepository interface {
		IndexNewHistory(ctx context.Context, project *types.ProjectHistory) error
	}

	//PolicyService check permission of client
	PolicyService interface {
		Validate(ctx context.Context, obj string, act string) error
	}

	HolidayService interface {
		RemoveProject(ctx context.Context, holidayID string, projectID string) error
		AssignProject(ctx context.Context, holidayID string, projectID string) error
	}

	Service struct {
		mongo   mongoRepository
		policy  PolicyService
		holiday HolidayService
		elastic elasticRepository
	}
)

func New(mongo mongoRepository, policy PolicyService, elastic elasticRepository, holiday HolidayService) *Service {
	return &Service{
		mongo:   mongo,
		policy:  policy,
		elastic: elastic,
		holiday: holiday,
	}
}

//Save all tasks, only update tasks has new update time
func (s *Service) Save(ctx context.Context, id string, req *types.SaveProject) (*types.ProjectHistory, error) {

	//only PM can create Project
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}
	//validate incoming data
	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to update project due to invalid req, %w", err)
		return nil, err
	}

	//Get lastest project
	//Get lastest edit user
	project, err := s.mongo.FindByProjectID(ctx, id)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by id: %v", err)
		return nil, fmt.Errorf("failed to check existing project by id: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return nil, status.Project().NotFoundProject
	}

	//validate each task
	for i, task := range req.Tasks {

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
				return nil, fmt.Errorf("Dev: "+devInTask+" %w", status.Project().NotFoundDev)
			}
		}
		//For sync data
		//Update only changed task
		for _, dbTask := range project.Tasks {
			if dbTask.TaskID == task.TaskID {
				if dbTask.UpdateAt == task.UpdateAt {
					logrus.Info("Updating task INFO" + task.Label + "This task UpdateTime not change => not updated" + "Code: project/update/line 113 ")
					req.Tasks[i] = dbTask
				}
				break
			}
		}
	}

	//Just update tasks field of project
	//VueJS will do all the validation for data comes in
	//So that user can edit tasks as long as they want, only save tasks when user use update task API

	history := &types.ProjectHistory{
		//news info
		Tasks:    req.Tasks,
		UpdateAt: time.Now(),

		//old
		ProjectID: project.ProjectID,
		Desc:      project.Desc,
		DevsID:    project.DevsID,
		CreatedAt: project.CreatedAt,
		CreaterID: project.ProjectID,
		Highlight: project.Highlight,
		Name:      project.Name,
		Action:    "Save project",
	}

	if err = s.mongo.Save(ctx, id, req); err != nil {
		logrus.Errorf("failed to update project: %v", err)
		return nil, fmt.Errorf("failed to update project, %w", err)
	}

	//TODO calculate workload
	err = s.elastic.IndexNewHistory(ctx, history)

	if err != nil {
		//TODO retry if fail
		//OR remove Updated info
		return nil, err
	}

	return history, nil
}

//Note: update data must inclue all the field,
func (s *Service) Update(ctx context.Context, id string, req *types.UpdateProject) (*types.ProjectHistory, error) {

	//only PM can update Project
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}
	//validate incoming data
	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to update project due to invalid req, %w", err)
		return nil, err
	}

	project, err := s.mongo.FindByProjectID(ctx, id)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by id: %v", err)
		return nil, fmt.Errorf("failed to check existing project by id: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return nil, status.Project().NotFoundProject
	}

	history := &types.ProjectHistory{
		//news info
		UpdateAt:  time.Now(),
		Name:      req.Name,
		Desc:      req.Desc,
		Highlight: req.Highlight,

		//old
		Tasks:     project.Tasks,
		ProjectID: project.ProjectID,
		DevsID:    project.DevsID,
		CreatedAt: project.CreatedAt,
		CreaterID: project.ProjectID,
		Action:    "Update project",
	}

	if err = s.mongo.Update(ctx, id, req); err != nil {
		logrus.Errorf("failed to update project: %v", err)
		return nil, fmt.Errorf("failed to update project, %w", err)
	}

	//TODO calculate workload
	err = s.elastic.IndexNewHistory(ctx, history)

	if err != nil {
		//TODO retry if fail
		//OR remove Updated info
		return nil, err
	}

	return history, nil
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

	existingProject, err := s.mongo.FindByName(ctx, req.Name, pm.UserID)

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

	if err = s.mongo.Create(ctx, project); err != nil {
		logrus.Errorf("failed to create project: %v", err)
		return nil, fmt.Errorf("failed to create project, %w", err)
	}

	history := &types.ProjectHistory{
		//news info
		UpdateAt: time.Now(),
		//old
		Tasks:     project.Tasks,
		ProjectID: project.ProjectID,
		Desc:      project.Desc,
		DevsID:    project.DevsID,
		CreatedAt: project.CreatedAt,
		CreaterID: project.ProjectID,
		Highlight: project.Highlight,
		Name:      project.Name,
		Action:    "Create project",
	}
	err = s.elastic.IndexNewHistory(ctx, history)
	if err != nil {
		//TODO retry if fail
		//OR remove Created project
		return nil, err
	}
	return project, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return err
	}

	//Remove project from all holiday
	if err := s.holiday.RemoveProject(ctx, "_all_holiday_", id); err != nil {
		logrus.Errorf("Fail to delete project due to %v", err)
		return fmt.Errorf("failed to remove project %v", err)
	}

	if err := s.mongo.Delete(ctx, id); err != nil {
		logrus.Errorf("Fail to delete project due to %v", err)
		return status.Project().NotFoundProject
	}

	return nil
}

func (s *Service) FindByID(ctx context.Context, id string) (*types.Project, error) {

	//TODO only devs of this project
	project, err := s.mongo.FindByProjectID(ctx, id)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by ID: %v", err)
		return nil, fmt.Errorf("failed to check existing project by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return nil, status.Project().NotFoundProject
	}
	return project, nil
}

func (s *Service) FindAllProjects(ctx context.Context) ([]*types.ProjectInfo, error) {

	user := auth.FromContext(ctx)
	logrus.Infof("Devs id: %v", user.UserID, " Find all project Code: projects/findall/ line: 259")

	projects, err := s.mongo.FindAllByUserID(ctx, user.UserID, user.Role)

	info := make([]*types.ProjectInfo, 0)
	for _, project := range projects {
		info = append(info, &types.ProjectInfo{
			ProjectID: project.ProjectID,
			Name:      project.Name,
			CreaterID: project.CreaterID,
			DevsID:    project.DevsID,
			Highlight: project.Highlight,
			UpdateAt:  project.UpdateAt,
			Tasks:     len(project.Tasks),
			Desc:      project.Desc,
			CreatedAt: project.CreatedAt,
		})
	}

	return info, err
}

//Removing devs_id not existing anymore (dev deleted)
//TODO: history
func (s *Service) RemoveDevs(ctx context.Context, userID string) error {
	projects, err := s.mongo.FindAllByUserID(ctx, userID, types.DEV)
	for _, project := range projects {
		//Remove dev from project
		if err := s.mongo.UpdateDevsID(ctx, []string{userID}, project.ProjectID, false); err != nil {
			logrus.Errorf("Failed to update devs ID in project info %v", err)
			return err
		}
		//Update all project's tasks
		taskInfo := make([]*types.Task, 0)
		for _, task := range project.Tasks {

			//remove devid from tasks
			for i, dev := range task.DevsID {
				if dev == userID {
					task.DevsID[i] = task.DevsID[len(task.DevsID)-1]
					task.DevsID = task.DevsID[:len(task.DevsID)-1]
				}
			}
			taskInfo = append(taskInfo, &types.Task{
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
				UpdateAt:         task.UpdateAt,
			})
		}

		newTasks := &types.SaveProject{
			Tasks: taskInfo,
		}

		if err = s.mongo.Save(ctx, project.ProjectID, newTasks); err != nil {
			logrus.Errorf("failed to update project: %v", err)
			return fmt.Errorf("failed to update project, %w", err)
		}

		//New history record
		history := &types.ProjectHistory{
			//news info
			UpdateAt: time.Now(),
			Tasks:    taskInfo,
			//old
			Name:      project.Name,
			Desc:      project.Desc,
			Highlight: project.Highlight,
			ProjectID: project.ProjectID,
			DevsID:    project.DevsID,
			CreatedAt: project.CreatedAt,
			CreaterID: project.ProjectID,
			Action:    "Remove devs",
		}

		return s.elastic.IndexNewHistory(ctx, history)

	}
	return err
}

func (s *Service) AddHoliday(ctx context.Context, holidayID string, projectID string) (string, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return "", err
	}

	// Check project exist
	_, err := s.mongo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by ID: %v", err)
		return "", fmt.Errorf("failed to check existing project by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return "", status.Project().NotFoundProject
	}

	if err := s.holiday.AssignProject(ctx, holidayID, projectID); err != nil {
		logrus.Errorf("Failed to update projects_id in holiday info %v", err)
		return "", err
	}

	return holidayID, nil
}

func (s *Service) RemoveHoliday(ctx context.Context, holidayID string, projectID string) (string, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return "", err
	}

	// Check project exist
	_, err := s.mongo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by ID: %v", err)
		return "", fmt.Errorf("failed to check existing project by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return "", status.Project().NotFoundProject
	}

	if err := s.holiday.RemoveProject(ctx, holidayID, projectID); err != nil {
		logrus.Errorf("Fail to delete project due to %v", err)
		return "", fmt.Errorf("failed to remove holiday from project %v", err)
	}

	return holidayID, nil
}

func (s *Service) AddDevs(ctx context.Context, userIDs []string, projectID string) ([]string, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}

	// Check project exist
	project, err := s.mongo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by ID: %v", err)
		return nil, fmt.Errorf("failed to check existing project by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return nil, status.Project().NotFoundProject
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

	if err := s.mongo.UpdateDevsID(ctx, userIDs, projectID, true); err != nil {
		logrus.Errorf("Failed to update devs ID in project info %v", err)
		return nil, err
	}

	return userIDs, nil
}

//Remove dev of this project
func (s *Service) RemoveDev(ctx context.Context, userID string, projectID string) error {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return err
	}
	// Check project exist
	project, err := s.mongo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		return fmt.Errorf("failed to check existing project by ID: %w", err)
	}
	if db.IsErrNotFound(err) {
		return status.Project().NotFoundProject
	}

	NotInProject := true
	for _, devID := range project.DevsID {
		if userID == devID {
			NotInProject = false
		}
	}
	if NotInProject {
		return status.Project().NotInProject
	}

	if err := s.mongo.UpdateDevsID(ctx, []string{userID}, projectID, false); err != nil {
		logrus.Errorf("Failed to update devs ID in project info %v", err)
		return err
	}

	return nil
}
