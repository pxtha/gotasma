package user

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
	"golang.org/x/crypto/bcrypt"
)

type (
	Repository interface {
		Create(context.Context, *types.User) (string, error)
		Delete(cxt context.Context, id string) error

		FindByEmail(ctx context.Context, email string) (*types.User, error)
		FindAllDev(ctx context.Context, createrID string) ([]*types.User, error)
		FindDevsByID(ctx context.Context, userIDs []string) ([]*types.User, error)
		FindByID(ctx context.Context, UserID string) (*types.User, error)
		FindByProjectID(ctx context.Context, projectID string) ([]*types.User, error)

		// Assign or remove project from user
		UpdateProjectsID(ctx context.Context, userID string, projectID string, addToSet bool) error
		// Assign or remove task from user
		UpdateTasksID(ctx context.Context, userID string, taskID string, addToSet bool) error
	}

	PolicyService interface {
		Validate(ctx context.Context, obj string, act string) error
	}

	WorkLoad interface {
		Delete(ctx context.Context, projectID string, userID string) error
		FindByUserID(ctx context.Context, userID string) ([]*types.WorkLoad, error)
	}

	Service struct {
		workload WorkLoad
		repo     Repository
		policy   PolicyService
	}
)

func New(repo Repository, policy PolicyService, workload WorkLoad) *Service {
	return &Service{
		repo:     repo,
		policy:   policy,
		workload: workload,
	}
}

func (s *Service) Register(ctx context.Context, req *types.RegisterRequest) (*types.User, error) {

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to register PM due to invalid req, %w", err)
		return nil, fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
	}

	existingUser, err := s.repo.FindByEmail(ctx, req.Email)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing PM by email, err: %v", err)
		return nil, fmt.Errorf("failed to check existing user by email: %w", err)
	}

	if existingUser != nil {
		logrus.Error("PM email already registered")
		logrus.WithContext(ctx).Debug("email already registered")
		return nil, status.User().DuplicatedEmail
	}

	rs, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("failed to check hash password, %v", err)
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	userID := uuid.New()

	user := &types.User{
		Password:  string(rs),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role:      types.PM,
		UserID:    userID,
		CreaterID: userID,
	}

	if _, err := s.repo.Create(ctx, user); err != nil {
		logrus.Errorf("failed to insert PM, %v", err)
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user.Strip(), nil
}

func (s *Service) CreateDev(ctx context.Context, req *types.RegisterRequest) (*types.User, error) {

	//Only PM can create dev
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}
	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to register DEV due to invalid req, %v", err)
		return nil, fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
	}

	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing DEV by email, err: %v", err)
		return nil, fmt.Errorf("failed to check existing user by email: %w", err)
	}

	if existingUser != nil {
		logrus.Warning("Devs email is unique in database, different PM cannot have same DEV email account")
		logrus.Infof("Devs email already registered in system")
		return nil, status.User().DuplicatedEmail
	}

	rs, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("failed to check hash password, %v", err)
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	pm := auth.FromContext(ctx)

	user := &types.User{
		Password:  string(rs),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role:      req.Role,
		UserID:    uuid.New(),
		CreaterID: pm.UserID,
	}

	if _, err := s.repo.Create(ctx, user); err != nil {
		logrus.Errorf("failed to insert DEV, %v", err)
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user.Strip(), nil
}

func (s *Service) FindByID(ctx context.Context, id string) (*types.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to find user user by ID, err: %v", err)
		return nil, err
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("User doesn't exist, err: %v", err)
		return nil, status.User().NotFoundUser
	}
	return user, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return err
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing user by ID, err: %v", err)
		return err
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("User doesn't exist, err: %v", err)
		return status.User().NotFoundUser
	}

	if user.Role == types.PM {
		logrus.Warning("This is PM_ID, cannot delete PM account")
		return status.Sercurity().InvalidAction
	}

	//Delele all devs workload
	if err := s.workload.Delete(ctx, "_all_projects_", id); err != nil {
		logrus.Error("Database error, Failed to delete all user's workload: %w", err)
		return fmt.Errorf("Failed to delete all user's workload: %w", err)
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) Auth(ctx context.Context, email, password string) (*types.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to check existing user by email, err: %v", err)
		return nil, status.Gen().Internal
	}
	if db.IsErrNotFound(err) {
		logrus.Debugf("user not found, email: %s", email)
		return nil, status.User().NotFoundUser
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.Error("invalid password")
		return nil, status.Auth().InvalidUserPassword
	}
	return user.Strip(), nil
}

func (s *Service) FindAllDev(ctx context.Context) ([]*types.UserInfo, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}

	pm := auth.FromContext(ctx)

	users, err := s.repo.FindAllDev(ctx, pm.UserID)
	if err != nil {
		logrus.Errorf("Database err: failed to find DEV, err: %v", err)
		return nil, fmt.Errorf("Failed to find devs info of this user: %w", err)
	}

	info := make([]*types.UserInfo, 0)
	for _, usr := range users {
		info = append(info, &types.UserInfo{
			Email:      usr.Email,
			FirstName:  usr.FirstName,
			LastName:   usr.LastName,
			Role:       usr.Role,
			CreaterID:  usr.CreaterID,
			UserID:     usr.UserID,
			CreatedAt:  usr.CreatedAt,
			UpdateAt:   usr.UpdateAt,
			ProjectsID: usr.ProjectsID,
			TasksID:    usr.TasksID,
		})
	}
	return info, nil
}

func (s *Service) FindByProjectID(ctx context.Context, projectID string) ([]*types.UserInfo, error) {

	users, err := s.repo.FindByProjectID(ctx, projectID)

	info := make([]*types.UserInfo, 0)
	for _, user := range users {
		info = append(info, &types.UserInfo{
			UserID:     user.UserID,
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Role:       user.Role,
			CreaterID:  user.CreaterID,
			CreatedAt:  user.CreatedAt,
			ProjectsID: user.ProjectsID,
			TasksID:    user.TasksID,
			UpdateAt:   user.UpdateAt,
		})
	}
	return info, err
}

// Mange project ids of each user
func (s *Service) RemoveProject(ctx context.Context, userID string, projectID string) error {

	if userID == "_all_devs_" {
		//remove project from all holidays
		users, err := s.repo.FindByProjectID(ctx, projectID)
		if err != nil {
			logrus.Error("Database error, Failed to remove project from user %w", err)
			return fmt.Errorf("Failed to get user by ID project: %w", err)
		}
		for _, user := range users {
			if err := s.repo.UpdateProjectsID(ctx, user.UserID, projectID, false); err != nil {
				logrus.Error("Database error, Failed to remove project from user %w", err)
				return fmt.Errorf("Failed to remove project from user: %w", err)
			}
		}
		return nil
	}
	//remove this user from this project
	//check user exist, is projectID in this user info?
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.Error("Failed to check existing user by id %w", err)
		return fmt.Errorf("Failed to check existing user by id: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Error("User not found")
		return status.User().NotFoundUser
	}

	hasProject := false
	for _, project := range user.ProjectsID {
		if project == projectID {
			hasProject = true
		}
	}

	if !hasProject {
		logrus.Errorf("Project didn't have this user")
		return status.User().NotFoundProject
	}

	if err := s.repo.UpdateProjectsID(ctx, userID, projectID, false); err != nil {
		logrus.Error("Database error, Failed to remove project from user %w", err)
		return fmt.Errorf("Failed to remove project from user: %w", err)
	}
	return err
}

func (s *Service) AddProject(ctx context.Context, userID string, projectID string) error {

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.Error("Failed to check existing user by id %w", err)
		return fmt.Errorf("Failed to check existing user by id: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Error("User not found")
		return status.User().NotFoundUser
	}

	hasProject := false
	for _, project := range user.ProjectsID {
		if project == projectID {
			hasProject = true
		}
	}

	if hasProject {
		logrus.Errorf("Project already had this user")
		return status.User().AlreadyInProject
	}

	return s.repo.UpdateProjectsID(ctx, userID, projectID, true)
}

// Manage tasks_id of each user
func (s *Service) AssignTask(ctx context.Context, projectID string, req *types.AssignDev) error {

	user, err := s.repo.FindByID(ctx, req.UserID)

	if err != nil && !db.IsErrNotFound(err) {
		logrus.Error("Failed to check existing user by id %w", err)
		return fmt.Errorf("Failed to check existing user by id: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Error("User not found")
		return status.User().NotFoundUser
	}
	//Check user in project
	hasProject := false
	for _, project := range user.ProjectsID {
		if project == projectID {
			hasProject = true
		}
	}
	if !hasProject {
		logrus.Errorf("This proecjt: %v not has this user: %v", projectID, req.UserID)
		return status.User().NotFoundProject
	}
	//Check user already has this task
	hasTask := false
	for _, task := range user.TasksID {
		if task == req.TaskID {
			hasTask = true
		}
	}
	if hasTask {
		logrus.Errorf("This task: %v already assign to this user: %v", req.UserID, req.TaskID)
		return status.User().AlreadyInTask
	}

	return s.repo.UpdateTasksID(ctx, req.UserID, req.TaskID, true)
}

func (s *Service) UnAssignTask(ctx context.Context, projectID string, req *types.UnAssignDev) error {

	if req.UserID == "_all_devs_" {
		users, err := s.repo.FindByProjectID(ctx, projectID)
		if err != nil {
			logrus.Error("Failed to get users by ID project %w", err)
			return fmt.Errorf("Failed to get users by ID project: %w", err)
		}
		for _, user := range users {
			if err := s.repo.UpdateTasksID(ctx, user.UserID, req.TaskID, false); err != nil {
				logrus.Error("Database error, Failed to remove task from user %w", err)
				return fmt.Errorf("Failed to remove task from user: %w", err)
			}
		}
		return nil
	}
	user, err := s.repo.FindByID(ctx, req.UserID)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.Error("Failed to check existing user by id %w", err)
		return fmt.Errorf("Failed to check existing user by id: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Error("User not found")
		return status.User().NotFoundUser
	}
	//Check user in project
	hasProject := false
	for _, project := range user.ProjectsID {
		if project == projectID {
			hasProject = true
		}
	}

	if !hasProject {
		logrus.Errorf("This proecjt: %v not has this user: %v", projectID, req.UserID)
		return status.User().NotFoundProject
	}
	//Check user not has this task
	hasTask := false
	for _, task := range user.TasksID {
		if task == req.TaskID {
			hasTask = true
		}
	}
	if !hasTask {
		logrus.Errorf("This task: %v have not assigned to this user yet: %v", req.UserID, req.TaskID)
		return status.User().NotFoundTask
	}

	return s.repo.UpdateTasksID(ctx, req.UserID, req.TaskID, false)
}

// Find all workloads
func (s *Service) FindWorkload(ctx context.Context, req *types.UsersWorkloadRequest) ([]*types.WorkLoad, error) {

	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Errorf("Fail to get workloads due to invalid req, %w", err)
		return nil, fmt.Errorf(err.Error()+"err: %w", status.Gen().BadRequest)
	}

	_, err := s.repo.FindByID(ctx, req.UserID)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.Errorf("failed to find user user by ID, err: %v", err)
		return nil, err
	}

	if db.IsErrNotFound(err) {
		logrus.Errorf("User doesn't exist, err: %v", err)
		return nil, status.User().NotFoundUser
	}

	workloads, err := s.workload.FindByUserID(ctx, req.UserID)
	if err != nil {
		logrus.Errorf("failed to find all workload of User by ID, err: %v", err)
		return nil, err
	}
	return workloads, nil
}
