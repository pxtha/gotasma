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
		Create(context.Context, *types.Project) error
		Delete(ctx context.Context, id string) error
		Save(ctx context.Context, id string, req *types.SaveProject) error
		Update(ctx context.Context, id string, req *types.UpdateProject) error

		FindByName(ctx context.Context, name string, createrID string) (*types.Project, error)
		FindByProjectID(ctx context.Context, projectID string) (*types.Project, error)
		FindAllByUserID(ctx context.Context, id string, role types.Role) ([]*types.Project, error)
		FindAllByHolidaysID(ctx context.Context, holidayID string) ([]*types.Project, error)
	}

	elasticRepository interface {
		IndexNewHistory(ctx context.Context, project *types.ProjectHistory) error
	}

	PolicyService interface {
		Validate(ctx context.Context, obj string, act string) error
	}

	UserService interface {
		RemoveProject(ctx context.Context, userID string, projectID string) error
		AssignProject(ctx context.Context, userID string, projectID string) error
		FindByProjectID(ctx context.Context, projectID string) ([]*types.UserInfo, error)
	}

	HolidayService interface {
		RemoveProject(ctx context.Context, holidayID string, projectID string) error
		AssignProject(ctx context.Context, holidayID string, projectID string) error
		FindByProjectID(ctx context.Context, projectID string) ([]*types.HolidayInfo, error)
	}

	TaskService interface {
		FindByProjectID(ctx context.Context, projectID string) ([]*types.Task, error)
		Update(ctx context.Context, projectID string, req *types.Task) error
		Create(ctx context.Context, projectID string, req *types.Task) (*types.Task, error)
		Delete(ctx context.Context, id string) error
	}

	Service struct {
		mongo   mongoRepository
		policy  PolicyService
		holiday HolidayService
		elastic elasticRepository
		user    UserService
		task    TaskService
	}
)

func New(mongo mongoRepository, policy PolicyService, elastic elasticRepository, holiday HolidayService, user UserService, task TaskService) *Service {
	return &Service{
		mongo:   mongo,
		policy:  policy,
		elastic: elastic,
		holiday: holiday,
		user:    user,
		task:    task,
	}
}

//Save all tasks, only update tasks has new update time, TODO: elasticsearch HISTORY
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

	project, err := s.mongo.FindByProjectID(ctx, id)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by id: %v", err)
		return nil, fmt.Errorf("failed to check existing project by id: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return nil, status.Project().NotFoundProject
	}

	//Get tasks from db: oldtask := s.mongo.getbyprojectid()
	dbTasks, err := s.task.FindByProjectID(ctx, project.ProjectID)
	if err != nil {
		logrus.Errorf("Database error: failed to get old tasks of project by project_id: %v", err)
		return nil, fmt.Errorf("Fail to save tasks of project: %w", err)
	}

	for _, newtask := range req.Tasks {
		exist := false
		if err := validator.Validate(newtask); err != nil {
			logrus.Errorf("Fail to update project due to invalid req, %w", err)
			return nil, err
		}

		// For sync data
		// Update only changed task base on UPDATETIME
		for _, oldtask := range dbTasks {
			if newtask.TaskID == oldtask.TaskID {
				if newtask.UpdateAt != oldtask.UpdateAt {
					//update this task
					//History on elastic search
					//TODO: add context: user info
					if err := s.task.Update(ctx, id, newtask); err != nil {
						logrus.Errorf("Fail to update tasks due to, %w", err)
						return nil, fmt.Errorf("Fail to update  tasks due to, %w", err)
					}
				}
				exist = true
				break
			}
		}
		//this task not exist in db -> create new task
		if !exist {
			//History on elastic search
			//TODO: add context: user info
			_, err := s.task.Create(ctx, id, newtask)
			if err != nil {
				return nil, fmt.Errorf("Fail to create new tasks due to, %w", err)
			}
		}
	}

	for _, oldtask := range dbTasks {
		deleted := true
		for _, newtask := range req.Tasks {
			if newtask.TaskID == oldtask.TaskID {
				deleted = false
				break
			}
		}
		//this task not exist in in newtask -> delete task
		if deleted {
			//History on elastic search
			//TODO: add context: user info
			if err := s.task.Delete(ctx, oldtask.TaskID); err != nil {
				logrus.Errorf("Fail to delete tasks due to, %w", err)
				return nil, fmt.Errorf("Fail to delete  tasks due to, %w", err)
			}
		}
	}
	//History on elastic search
	//TODO: add context: user info
	history := &types.ProjectHistory{
		//news info
		UpdateAt: time.Now(),
		//old
		ProjectID: project.ProjectID,
		Desc:      project.Desc,
		CreatedAt: project.CreatedAt,
		CreaterID: project.ProjectID,
		Highlight: project.Highlight,
		Name:      project.Name,
		Action:    "Save project",
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

//Update Note: update data must inclue all the field,  TODO: elasticsearch HISTORY
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

	//Check project existing
	project, err := s.mongo.FindByProjectID(ctx, id)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by id: %v", err)
		return nil, fmt.Errorf("failed to check existing project by id: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return nil, status.Project().NotFoundProject
	}

	//Update project info
	if err = s.mongo.Update(ctx, id, req); err != nil {
		logrus.Errorf("failed to update project: %v", err)
		return nil, fmt.Errorf("failed to update project, %w", err)
	}

	//History elastic search info
	history := &types.ProjectHistory{
		//news info
		UpdateAt:  time.Now(),
		Name:      req.Name,
		Desc:      req.Desc,
		Highlight: req.Highlight,

		//old
		ProjectID: project.ProjectID,
		CreatedAt: project.CreatedAt,
		CreaterID: project.ProjectID,
		Action:    "Update project",
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

//Create TODO: elasticsearch HISTORY
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

	//Check duplicate project's name of this PM
	existingProject, err := s.mongo.FindByName(ctx, req.Name, pm.UserID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by name: %v", err)
		return nil, fmt.Errorf("failed to check existing project by name: %w", err)
	}

	if existingProject != nil {
		logrus.Errorf("Project already exsit")
		return nil, status.Project().DuplicateProject
	}

	//Create new project
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

	//History elastic search info
	history := &types.ProjectHistory{
		//news info
		UpdateAt: time.Now(),
		//old
		ProjectID: project.ProjectID,
		Desc:      project.Desc,
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

//Delete TODO: elasticsearch HISTORY
func (s *Service) Delete(ctx context.Context, projectID string) error {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return err
	}

	//Remove this project from all holiday
	if err := s.holiday.RemoveProject(ctx, "_all_holiday_", projectID); err != nil {
		logrus.Errorf("Fail to delete project due to %v", err)
		return fmt.Errorf("failed to remove project %v", err)
	}

	//Remove this project from all devs
	if err := s.user.RemoveProject(ctx, "_all_devs_", projectID); err != nil {
		logrus.Errorf("Fail to delete project due to %v", err)
		return fmt.Errorf("failed to remove project %v", err)
	}

	//Remove project, return err if project not exist
	if err := s.mongo.Delete(ctx, projectID); err != nil {
		logrus.Errorf("Fail to delete project due to %v", err)
		return status.Project().NotFoundProject
	}

	return nil
}

//FindByID TODO: only devs of this project
func (s *Service) FindByID(ctx context.Context, id string) (*types.Project, error) {

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

//FindAllProjects  by user_id: PM and DEVs, TODO: policy
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
			Highlight: project.Highlight,
			UpdateAt:  project.UpdateAt,
			Desc:      project.Desc,
			CreatedAt: project.CreatedAt,
		})
	}

	return info, err
}

//AddHoliday Mange holiday of project
func (s *Service) AddHoliday(ctx context.Context, req *types.AddHolidayRequest, projectID string) (string, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return "", err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to add holiday to project due to invalid req, %w", err)
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

	if err := s.holiday.AssignProject(ctx, req.HolidayID, projectID); err != nil {
		logrus.Errorf("Failed to update projects_id in holiday info %v", err)
		return "", err
	}

	return req.HolidayID, nil
}

//RemoveHoliday Mange holiday of project
func (s *Service) RemoveHoliday(ctx context.Context, req *types.RemoveHolidayRequest, projectID string) (string, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return "", err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to remove holiday from project due to invalid req, %w", err)
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

	if err := s.holiday.RemoveProject(ctx, req.HolidayID, projectID); err != nil {
		logrus.Errorf("Fail to delete project due to %v", err)
		return "", err
	}

	return req.HolidayID, nil
}

//FindAllHolidays Mange holiday of project
func (s *Service) FindAllHolidays(ctx context.Context, projectID string) ([]*types.HolidayInfo, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}
	// Check project exist
	_, err := s.mongo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by ID: %v", err)
		return nil, fmt.Errorf("failed to check existing project by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist", err)
		return nil, status.Project().NotFoundProject
	}

	var info []*types.HolidayInfo

	info, err = s.holiday.FindByProjectID(ctx, projectID)
	if err != nil {
		logrus.Errorf("Fail to get holidays of project: %v", err)
		return nil, fmt.Errorf("failed to get holidays of project: %w", err)
	}
	return info, nil
}

//AddDev Mange devs of project
func (s *Service) AddDev(ctx context.Context, req *types.AddUsersRequest, projectID string) (string, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return "", err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to add user to project due to invalid req, %w", err)
		return "", err
	}

	// Check project exist
	project, err := s.mongo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by ID: %v", err)
		return "", fmt.Errorf("failed to check existing project by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist", err)
		return "", status.Project().NotFoundProject
	}

	if req.UserID == project.CreaterID {
		logrus.Error("Cannot at, this is creater of this project")
		return "", status.Project().AlreadyInProject
	}

	if err := s.user.AssignProject(ctx, req.UserID, projectID); err != nil {
		logrus.Errorf("Failed to update devs ID in project info %v", err)
		return "", err
	}

	return req.UserID, nil
}

//RemoveDev  Mange devs of project
func (s *Service) RemoveDev(ctx context.Context, req *types.RemoveUserRequest, projectID string) (string, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return "", err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to remove user from project due to invalid req, %w", err)
		return "", err
	}

	// Check project exist
	project, err := s.mongo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		return "", fmt.Errorf("failed to check existing project by ID: %w", err)
	}
	if db.IsErrNotFound(err) {
		return "", status.Project().NotFoundProject
	}

	if req.UserID == project.CreaterID {
		logrus.Error("Cannot at, this is creater of this project")
		return "", status.Project().AlreadyInProject
	}

	if err := s.user.RemoveProject(ctx, req.UserID, projectID); err != nil {
		logrus.Errorf("Fail to delete project due to %v", err)
		return "", err
	}

	return req.UserID, nil
}

//FindAllDevs  Mange devs of project
func (s *Service) FindAllDevs(ctx context.Context, projectID string) ([]*types.UserInfo, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}
	// Check project exist
	_, err := s.mongo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by ID: %v", err)
		return nil, fmt.Errorf("failed to check existing project by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return nil, status.Project().NotFoundProject
	}

	var info []*types.UserInfo

	info, err = s.user.FindByProjectID(ctx, projectID)
	if err != nil {
		logrus.Errorf("Fail to get devs of project: %v", err)
		return nil, fmt.Errorf("failed to get devs of project: %w", err)
	}
	return info, nil
}
