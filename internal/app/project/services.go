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
		IndexNewHistory(ctx context.Context, project *types.History) error
	}

	PolicyService interface {
		Validate(ctx context.Context, obj string, act string) error
	}

	UserService interface {
		// Projects
		AddProject(ctx context.Context, userID string, projectID string) error
		RemoveProject(ctx context.Context, userID string, projectID string) error
		FindByProjectID(ctx context.Context, projectID string) ([]*types.UserInfo, error)
		FindByID(ctx context.Context, id string) (*types.User, error)

		// Tasks
		AssignTask(ctx context.Context, projectID string, req *types.AssignDev) error
		UnAssignTask(ctx context.Context, projectID string, req *types.UnAssignDev) error
	}

	HolidayService interface {
		RemoveProject(ctx context.Context, holidayID string, projectID string) error
		AssignProject(ctx context.Context, holidayID string, projectID string) error
		FindByProjectID(ctx context.Context, projectID string) ([]*types.HolidayInfo, error)
	}

	TaskService interface {
		FindByProjectID(ctx context.Context, projectID string) ([]*types.Task, error)
		FindDetailByProjectID(ctx context.Context, projectID string) ([]*types.TaskDetailInfo, error)
		FindByID(ctx context.Context, id string) (*types.Task, error)
		FindByIDs(ctx context.Context, ids []string) ([]*types.TaskInfo, error)

		Update(ctx context.Context, projectID string, req *types.Task) error
		Create(ctx context.Context, projectID string, req *types.Task) (*types.Task, error)
		Delete(ctx context.Context, id string) error
	}

	Workload interface {
		FindByID(ctx context.Context, projectID string, userID string) (*types.WorkLoad, error)
		Update(ctx context.Context, projectID string, userID string, overload []*types.Interval) error
		Create(ctx context.Context, projectID string, userID string) error
		Delete(ctx context.Context, projectID string, userID string) error
	}

	Service struct {
		mongo    mongoRepository
		policy   PolicyService
		holiday  HolidayService
		elastic  elasticRepository
		user     UserService
		task     TaskService
		workload Workload
	}
)

func New(mongo mongoRepository, policy PolicyService, elastic elasticRepository, holiday HolidayService, user UserService, task TaskService, workload Workload) *Service {
	return &Service{
		mongo:    mongo,
		policy:   policy,
		elastic:  elastic,
		holiday:  holiday,
		user:     user,
		task:     task,
		workload: workload,
	}
}

// Manage project
func (s *Service) Save(ctx context.Context, id string, req *types.SaveProject) ([]*types.Task, error) {
	// only PM can create Project
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}
	// validate incoming data
	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to update project due to invalid req, %w", err)
		return nil, fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
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

	// Get tasks from db: oldtask := s.mongo.getbyprojectid()
	dbTasks, err := s.task.FindByProjectID(ctx, project.ProjectID)
	if err != nil {
		logrus.Errorf("Database error: failed to get old tasks of project by project_id: %v", err)
		return nil, fmt.Errorf("Fail to save tasks of project: %w", err)
	}

	for _, newtask := range req.Tasks {
		exist := false
		if err := validator.Validate(newtask); err != nil {
			logrus.Errorf("Fail to update project due to invalid req, %w", err)
			return nil, fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
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
		// This task not exist in db -> create new task
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
			//Remove task from dev
			req := &types.UnAssignDev{
				TaskID: oldtask.TaskID,
				UserID: "_all_devs_",
			}
			if err := s.user.UnAssignTask(ctx, id, req); err != nil {
				logrus.Errorf("Failed to update tasks_id in user info %v", err)
				return nil, err
			}
			//History on elastic search
			//TODO: add context: user info
			if err := s.task.Delete(ctx, oldtask.TaskID); err != nil {
				logrus.Errorf("Fail to delete tasks due to, %w", err)
				return nil, fmt.Errorf("Fail to delete  tasks due to, %w", err)
			}
		}
	}
	// History on elastic search
	// TODO: add context: user info
	//Gets for history
	tasks, err := s.task.FindByProjectID(ctx, id)
	if err != nil {
		logrus.Errorf("Database error: failed to get all tasks of project by project_id: %v", err)
		return nil, fmt.Errorf("Fail to save history to es: %w", err)
	}
	users, err := s.user.FindByProjectID(ctx, id)
	if err != nil {
		logrus.Errorf("Database error: failed to get users of project by project_id: %v", err)
		return nil, fmt.Errorf("Fail to save history to es: %w", err)
	}
	holidays, err := s.holiday.FindByProjectID(ctx, id)
	if err != nil {
		logrus.Errorf("Database error: failed to get holidays of project by project_id: %v", err)
		return nil, fmt.Errorf("Fail to save history to es: %w", err)
	}

	history := &types.History{
		// old
		ProjectID: project.ProjectID,
		Desc:      project.Desc,
		Tasks:     tasks,
		Devs:      users,
		Holiday:   holidays,
		CreatedAt: project.CreatedAt,
		CreaterID: project.ProjectID,
		Highlight: project.Highlight,
		Name:      project.Name,
		Action:    "Save project",
	}

	// TODO calculate workload
	err = s.elastic.IndexNewHistory(ctx, history)
	if err != nil {
		//TODO retry if fail
		//OR remove Updated info
		return nil, err
	}

	return req.Tasks, nil
}

func (s *Service) Update(ctx context.Context, id string, req *types.UpdateProject) (*types.UpdateProject, error) {

	//only PM can update Project
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}
	//validate incoming data
	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to update project due to invalid req, %w", err)
		return nil, fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
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

	//Gets for history
	tasks, err := s.task.FindByProjectID(ctx, id)
	if err != nil {
		logrus.Errorf("Database error: failed to get all tasks of project by project_id: %v", err)
		return nil, fmt.Errorf("Fail to save history to es: %w", err)
	}
	users, err := s.user.FindByProjectID(ctx, id)
	if err != nil {
		logrus.Errorf("Database error: failed to get users of project by project_id: %v", err)
		return nil, fmt.Errorf("Fail to save history to es: %w", err)
	}
	holidays, err := s.holiday.FindByProjectID(ctx, id)
	if err != nil {
		logrus.Errorf("Database error: failed to get holidays of project by project_id: %v", err)
		return nil, fmt.Errorf("Fail to save history to es: %w", err)
	}
	//History elastic search info
	history := &types.History{
		//news info
		Name:      req.Name,
		Desc:      req.Desc,
		Highlight: req.Highlight,
		Tasks:     tasks,
		Devs:      users,
		Holiday:   holidays,
		//old
		ProjectID: project.ProjectID,
		CreatedAt: project.CreatedAt,
		Action:    "Update project",
	}

	//TODO calculate workload
	err = s.elastic.IndexNewHistory(ctx, history)

	if err != nil {
		//TODO retry if fail
		//OR remove Created project
		logrus.Error(err)
	}

	return req, nil
}

func (s *Service) Create(ctx context.Context, req *types.CreateProjectRequest) (*types.Project, error) {

	pm := auth.FromContext(ctx)

	//only PM can create Project
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to create project due to invalid req, %w", err)
		return nil, fmt.Errorf(err.Error()+" err: %w", status.Gen().BadRequest)
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
	history := &types.History{
		//news info
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
		logrus.Error(err)
	}
	return project, nil
}

func (s *Service) Delete(ctx context.Context, projectID string) error {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return err
	}

	project, err := s.mongo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by id: %v", err)
		return fmt.Errorf("failed to check existing project by id: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("Project doesn't exist")
		return status.Project().NotFoundProject
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
	//Delete all tasks of project from db
	//Remove all tasks of this project from devs
	tasks, err := s.task.FindByProjectID(ctx, projectID)
	if err != nil {
		logrus.Errorf("Fail get all task of project to delete project, due to %v", err)
		return fmt.Errorf("Fail to delete of all tasks of project %v", err)
	}
	for _, task := range tasks {
		req := &types.UnAssignDev{
			TaskID: task.TaskID,
			UserID: "_all_devs_",
		}
		if err := s.user.UnAssignTask(ctx, projectID, req); err != nil {
			logrus.Errorf("Failed to update tasks_id in user info %v", err)
			return err
		}
		if err := s.task.Delete(ctx, task.TaskID); err != nil {
			logrus.Errorf("Fail to delete tasks due to, %w", err)
			return fmt.Errorf("Fail to delete  tasks due to, %w", err)
		}
	}
	//Delete all workload by project id
	if err := s.workload.Delete(ctx, projectID, "_all_devs_"); err != nil {
		logrus.Errorf("Fail to delete project due to %v", err)
		return fmt.Errorf("Fail to delete all workload of this project from all user due to: %w", err)
	}
	//Remove project, return err if project not exist
	if err := s.mongo.Delete(ctx, projectID); err != nil {
		logrus.Errorf("Fail to delete project due to %v", err)
		return status.Project().NotFoundProject
	}

	//

	history := &types.History{
		// news info
		// old
		ProjectID: project.ProjectID,
		Desc:      project.Desc,
		CreatedAt: project.CreatedAt,
		Highlight: project.Highlight,
		Name:      project.Name,
		Action:    "Delete project",
	}
	err = s.elastic.IndexNewHistory(ctx, history)
	if err != nil {
		//TODO retry if fail
		//OR remove Created project
		logrus.Error(err)
	}
	return nil
}

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

// Mange holiday of project
func (s *Service) AddHoliday(ctx context.Context, req *types.AddHolidayRequest, projectID string) (string, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return "", err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to add holiday to project due to invalid req, %w", err)
		return "", fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
	}

	// Check project exist
	project, err := s.mongo.FindByProjectID(ctx, projectID)

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
	//Get for history
	holidays, err := s.holiday.FindByProjectID(ctx, projectID)
	if err != nil {
		logrus.Errorf("Database error: failed to get holidays of project by project_id: %v", err)
		return "", fmt.Errorf("Fail to save history to es: %w", err)
	}
	history := &types.History{
		// news info
		// old
		Holiday:   holidays,
		ProjectID: project.ProjectID,
		Desc:      project.Desc,
		CreatedAt: project.CreatedAt,
		Highlight: project.Highlight,
		Name:      project.Name,
		Action:    "Add holiday",
	}
	err = s.elastic.IndexNewHistory(ctx, history)
	if err != nil {
		//TODO retry if fail
		//OR remove Created project
		logrus.Error(err)
	}

	return req.HolidayID, nil
}

func (s *Service) RemoveHoliday(ctx context.Context, req *types.RemoveHolidayRequest, projectID string) (string, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return "", err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to remove holiday from project due to invalid req, %w", err)
		return "", fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
	}

	// Check project exist
	project, err := s.mongo.FindByProjectID(ctx, projectID)

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

	//get for history
	holidays, err := s.holiday.FindByProjectID(ctx, projectID)
	if err != nil {
		logrus.Errorf("Database error: failed to get holidays of project by project_id: %v", err)
		return "", fmt.Errorf("Fail to save history to es: %w", err)
	}
	history := &types.History{
		// news info
		// old
		Holiday:   holidays,
		ProjectID: project.ProjectID,
		Desc:      project.Desc,
		CreatedAt: project.CreatedAt,
		Highlight: project.Highlight,
		Name:      project.Name,
		Action:    "Remove holiday",
	}
	err = s.elastic.IndexNewHistory(ctx, history)
	if err != nil {
		//TODO retry if fail
		//OR remove Created project
		logrus.Error(err)
	}

	return req.HolidayID, nil
}

func (s *Service) FindAllHolidays(ctx context.Context, projectID string) ([]*types.HolidayInfo, error) {

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

// Mange devs of project
func (s *Service) AddDev(ctx context.Context, req *types.AddUsersRequest, projectID string) (string, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return "", err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to add user to project due to invalid req, %w", err)
		return "", fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
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

	if err := s.user.AddProject(ctx, req.UserID, projectID); err != nil {
		logrus.Errorf("Failed to update projects_ID in user info %v", err)
		return "", err
	}
	//Create new workload record
	if err = s.workload.Create(ctx, projectID, req.UserID); err != nil {
		logrus.Errorf("Failed to create new workload for this user in this project %v", err)
		return "", err
	}

	users, err := s.user.FindByProjectID(ctx, projectID)
	if err != nil {
		logrus.Errorf("Database error: failed to get users of project by project_id: %v", err)
		return "", fmt.Errorf("Fail to save history to es: %w", err)
	}

	history := &types.History{
		// news info
		// old
		Devs:      users,
		ProjectID: project.ProjectID,
		Desc:      project.Desc,
		CreatedAt: project.CreatedAt,
		Highlight: project.Highlight,
		Name:      project.Name,
		Action:    "Add dev",
	}

	err = s.elastic.IndexNewHistory(ctx, history)
	if err != nil {
		//TODO retry if fail
		//OR remove Created project
		logrus.Error(err)
	}

	return req.UserID, nil
}

func (s *Service) RemoveDev(ctx context.Context, req *types.RemoveUserRequest, projectID string) (string, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return "", err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to remove user from project due to invalid req, %w", err)
		return "", fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
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
		logrus.Error("Cannot remove, this is creater of this project")
		return "", status.Project().ProjectCreater
	}

	// Remove project ==> remove all tasks of this project in this dev
	// Remove all tasks of this project from devs

	//TODO: get only task assigned to this user
	//Doing: get all task of project - skipp task that not in project - skip err - WARNING
	userInfo, err := s.user.FindByID(ctx, req.UserID)
	if err != nil {
		logrus.Errorf("Fail get user info, due to %v", err)
		return "", fmt.Errorf("Fail get all tasks in user info, due to%w", err)
	}

	// Get only task assigned to user
	tasks, err := s.task.FindByIDs(ctx, userInfo.TasksID)
	if err != nil {
		logrus.Errorf("Fail get all task of user in this project, due to %v", err)
		return "", fmt.Errorf("Fail get all task of user in this project,%w", err)
	}

	for _, task := range tasks {
		logrus.Info(task.TaskID)
		req := &types.UnAssignDev{
			TaskID: task.TaskID,
			UserID: req.UserID,
		}
		if err := s.user.UnAssignTask(ctx, projectID, req); err != nil {
			logrus.Errorf("Failed to update tasks_id in user info %v", err)
			return "", fmt.Errorf("Failed to update tasks_id in user info %w", err)
		}
	}

	//remove dev from project
	if err := s.user.RemoveProject(ctx, req.UserID, projectID); err != nil {
		logrus.Errorf("Fail to update projects_ID in user info due to %v", err)
		return "", err
	}
	// Delete workload of this project in this user
	if err = s.workload.Delete(ctx, projectID, req.UserID); err != nil {
		logrus.Errorf("Failed to delete workload of this project in this user %v", err)
		return "", err
	}

	users, err := s.user.FindByProjectID(ctx, projectID)
	if err != nil {
		logrus.Errorf("Database error: failed to get users of project by project_id: %v", err)
		return "", fmt.Errorf("Fail to save history to es: %w", err)
	}

	history := &types.History{
		// news info
		// old
		Devs:      users,
		ProjectID: project.ProjectID,
		Desc:      project.Desc,
		CreatedAt: project.CreatedAt,
		Highlight: project.Highlight,
		Name:      project.Name,
		Action:    "Remove dev",
	}

	err = s.elastic.IndexNewHistory(ctx, history)
	if err != nil {
		//TODO retry if fail
		//OR remove Created project
		logrus.Error(err)
	}

	return req.UserID, nil
}

func (s *Service) AssignDev(ctx context.Context, projectID string, req *types.AssignDev) (*types.WorkLoadInfo, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to assign task to user due to invalid req, %w", err)
		return nil, fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
	}

	// Check project exist
	project, err := s.mongo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by ID: %v", err)
		return nil, fmt.Errorf("failed to check existing project by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Error("Project doesn't exist")
		return nil, status.Project().NotFoundProject
	}

	// Validate task
	task, err := s.task.FindByID(ctx, req.TaskID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing task by ID: %v", err)
		return nil, fmt.Errorf("failed to check existing task by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Error("Task doesn't exist")
		return nil, status.Task().NotFoundTask
	}

	if task.ProjectID != project.ProjectID {
		logrus.Error("Task not in project")
		return nil, status.Task().NotInProject
	}

	// Validate dev
	if err := s.user.AssignTask(ctx, projectID, req); err != nil {
		logrus.Errorf("Failed to update tasks_id in user info %v", err)
		return nil, err
	}

	// calculate workload of this user in this project
	userInfo, err := s.user.FindByID(ctx, req.UserID)
	if err != nil {
		logrus.Errorf("Fail get user info, due to %v", err)
		return nil, fmt.Errorf("Fail get all tasks in user info, due to%v", err)
	}

	// Get only task assigned to user
	tasks, err := s.task.FindByIDs(ctx, userInfo.TasksID)
	if err != nil {
		logrus.Errorf("Fail get all task of user in this project, due to %v", err)
		return nil, fmt.Errorf("Fail get all task of user in this project,%v", err)
	}
	overload := make(map[int]int)
	for i := 0; i < len(tasks)-1; i++ {
		logrus.Info("loop main", tasks[i])
		for j := i + 1; j < len(tasks); j++ {
			logrus.Info("loop child", tasks[j])
			if tasks[i].End > tasks[j].Start && tasks[i].End < tasks[j].End {
				overload[tasks[j].Start] = tasks[i].End
			}
			if tasks[i].End > tasks[j].Start && tasks[i].End >= tasks[j].End {
				overload[tasks[j].Start] = tasks[j].End
			}
		}
	}
	logrus.Info(overload)

	overloadInfo := make([]*types.Interval, 0)
	for from, to := range overload {
		overloadInfo = append(overloadInfo, &types.Interval{
			From: time.Unix(0, int64(from)*int64(time.Millisecond)),
			To:   time.Unix(0, int64(to)*int64(time.Millisecond)),
		})
	}

	if err := s.workload.Update(ctx, projectID, req.UserID, overloadInfo); err != nil {
		logrus.Errorf("Fail calculate new workload of this user in this project, due to %v", err)
		return nil, fmt.Errorf("Fail calculate new workload: %v", err)
	}

	workloadInfo := &types.WorkLoadInfo{
		ProjectID: projectID,
		UserID:    req.UserID,
		Overload:  overloadInfo,
	}

	users, err := s.user.FindByProjectID(ctx, projectID)
	if err != nil {
		logrus.Errorf("Database error: failed to get users of project by project_id: %v", err)
		return nil, fmt.Errorf("Fail to save history to es: %w", err)
	}

	history := &types.History{
		// news info
		// old
		WorkLoad:  workloadInfo,
		Devs:      users,
		ProjectID: project.ProjectID,
		Desc:      project.Desc,
		CreatedAt: project.CreatedAt,
		Highlight: project.Highlight,
		Name:      project.Name,
		Action:    "Assign dev to task",
	}

	err = s.elastic.IndexNewHistory(ctx, history)
	if err != nil {
		//TODO retry if fail
		//OR remove Created project
		logrus.Error(err)
	}

	return workloadInfo, nil
}

func (s *Service) UnAssignDev(ctx context.Context, projectID string, req *types.UnAssignDev) (*types.WorkLoadInfo, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to unassign task from user due to invalid req, %w", err)
		return nil, fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
	}

	// Check project exist
	project, err := s.mongo.FindByProjectID(ctx, projectID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing project by ID: %v", err)
		return nil, fmt.Errorf("failed to check existing project by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Error("Project doesn't exist")
		return nil, status.Project().NotFoundProject
	}

	// Validate task
	task, err := s.task.FindByID(ctx, req.TaskID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing task by ID: %v", err)
		return nil, fmt.Errorf("failed to check existing task by ID: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Error("Task doesn't exist")
		return nil, status.Task().NotFoundTask
	}

	if task.ProjectID != project.ProjectID {
		logrus.Error("Task not in project")
		return nil, status.Task().NotInProject
	}

	// Validate dev
	if err := s.user.UnAssignTask(ctx, projectID, req); err != nil {
		logrus.Errorf("Failed to update tasks_id in user info %v", err)
		return nil, err
	}

	//calculate workload
	userInfo, err := s.user.FindByID(ctx, req.UserID)
	if err != nil {
		logrus.Errorf("Fail get user info, due to %v", err)
		return nil, fmt.Errorf("Fail get all tasks in user info, due to%v", err)
	}

	// Get only task assigned to user
	tasks, err := s.task.FindByIDs(ctx, userInfo.TasksID)
	if err != nil {
		logrus.Errorf("Fail get all task of user in this project, due to %v", err)
		return nil, fmt.Errorf("Fail get all task of user in this project,%v", err)
	}

	overload := make(map[int]int)
	for i := 0; i < len(tasks)-1; i++ {
		for j := i + 1; j < len(tasks); j++ {
			if tasks[i].End > tasks[j].Start && tasks[i].End < tasks[j].End {
				overload[tasks[j].Start] = tasks[i].End
			}
			if tasks[i].End > tasks[j].Start && tasks[i].End >= tasks[j].End {
				overload[tasks[j].Start] = tasks[j].End
			}
		}
	}

	overloadInfo := make([]*types.Interval, 0)
	for from, to := range overload {
		overloadInfo = append(overloadInfo, &types.Interval{
			From: time.Unix(0, int64(from)*int64(time.Millisecond)),
			To:   time.Unix(0, int64(to)*int64(time.Millisecond)),
		})
	}

	if err := s.workload.Update(ctx, projectID, req.UserID, overloadInfo); err != nil {
		logrus.Errorf("Fail calculate new workload of this user in this project, due to %v", err)
		return nil, fmt.Errorf("Fail calculate new workload: %v", err)
	}

	workloadInfo := &types.WorkLoadInfo{
		ProjectID: projectID,
		UserID:    req.UserID,
		Overload:  overloadInfo,
	}

	users, err := s.user.FindByProjectID(ctx, projectID)
	if err != nil {
		logrus.Errorf("Database error: failed to get users of project by project_id: %v", err)
		return nil, fmt.Errorf("Fail to save history to es: %w", err)
	}

	history := &types.History{
		// news info
		// old
		WorkLoad:  workloadInfo,
		Devs:      users,
		ProjectID: project.ProjectID,
		Desc:      project.Desc,
		CreatedAt: project.CreatedAt,
		Highlight: project.Highlight,
		Name:      project.Name,
		Action:    "Assign dev to task",
	}

	err = s.elastic.IndexNewHistory(ctx, history)
	if err != nil {
		//TODO retry if fail
		//OR remove Created project
		logrus.Error(err)
	}

	return workloadInfo, nil
}

func (s *Service) FindAllDevs(ctx context.Context, projectID string) ([]*types.UserInfo, error) {

	// if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
	// 	return nil, err
	// }
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

// Manage tasks
func (s *Service) FindAllTasks(ctx context.Context, projectID string) ([]*types.TaskDetailInfo, error) {

	// if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
	// 	return nil, err
	// }
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

	var info []*types.TaskDetailInfo

	info, err = s.task.FindDetailByProjectID(ctx, projectID)
	if err != nil {
		logrus.Errorf("Fail to get tasks of project: %v", err)
		return nil, fmt.Errorf("failed to get tasks of project: %w", err)
	}
	return info, nil
}
