package holiday

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

const (
	MilisecondInDay = 86400000
)

type (
	Repository interface {
		Create(ctx context.Context, holiday *types.Holiday) error
		Update(ctx context.Context, holiday *types.HolidayInfo, id string) error
		Delete(ctx context.Context, id string) error

		FindByTitle(ctx context.Context, title string, createrID string) (*types.Holiday, error)
		FindByID(ctx context.Context, id string) (*types.Holiday, error)
		FindByProjectID(ctx context.Context, id string) ([]*types.Holiday, error)
		FindAll(ctx context.Context, createrID string) ([]*types.Holiday, error)

		//Assign or remove project from holiday
		UpdateProjectsID(ctx context.Context, holidayID string, projectID string, addToSet bool) error
	}
	PolicyServices interface {
		Validate(ctx context.Context, obj string, act string) error
	}

	Service struct {
		repo   Repository
		policy PolicyServices
	}
)

func New(repo Repository, policy PolicyServices) *Service {
	return &Service{
		repo:   repo,
		policy: policy,
	}
}

func (s *Service) Create(ctx context.Context, req *types.HolidayRequest) (*types.Holiday, error) {
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}

	if err := validator.Validate(req); err != nil {
		logrus.Error("Failed to validate input create holiday %v", err)
		validateErr := err.Error()
		return nil, fmt.Errorf(validateErr+"err: %w", status.Gen().BadRequest)
	}

	//Duration >= 1 day
	holidayDuration := ((req.End - req.Start) / MilisecondInDay)
	if holidayDuration < 1 {
		logrus.Error("Failed to validate input create holiday ")
		return nil, status.Holiday().InvalidHoliday
	}

	//check existing Holiday of this PM
	pm := auth.FromContext(ctx)
	existingHoliday, err := s.repo.FindByTitle(ctx, req.Title, pm.UserID)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.Error("Failed to check existing holiday by title %w", err)
		return nil, fmt.Errorf("Failed to check existing holiday by title: %w", err)
	}
	if existingHoliday != nil {
		logrus.Error("Holiday all ready exist")
		return nil, status.Holiday().DuplicatedHoliday
	}

	holiday := &types.Holiday{
		Title:     req.Title,
		Start:     req.Start,
		End:       req.End,
		HolidayID: uuid.New(),
		Duration:  holidayDuration,
		CreaterID: pm.UserID,
	}
	if err := s.repo.Create(ctx, holiday); err != nil {
		return nil, fmt.Errorf("Faild to insert Holiday, %w", err)
	}
	return holiday, nil
}

//Update holiday info
func (s *Service) Update(ctx context.Context, id string, req *types.HolidayRequest) (*types.HolidayInfo, error) {
	//Check policy
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}

	//Validate [input]
	if err := validator.Validate(req); err != nil {
		logrus.Error("Failed to validate input update holiday %v", err)
		validateErr := err.Error()
		return nil, fmt.Errorf(validateErr+"err: %w", status.Gen().BadRequest)
	}

	holidayDuration := ((req.End - req.Start) / MilisecondInDay)
	if holidayDuration < 1 {
		logrus.Error("Failed to validate input update holiday ")
		return nil, status.Holiday().InvalidHoliday
	}

	//check existing holiday of this PM
	_, err := s.repo.FindByID(ctx, id)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.Error("Failed to check existing holiday by title %w", err)
		return nil, fmt.Errorf("Failed to check existing holiday by title: %w", err)
	}
	if db.IsErrNotFound(err) {
		logrus.Error("Not found holiday")
		return nil, status.Holiday().NotFoundHoliday
	}

	info := &types.HolidayInfo{
		Title:    req.Title,
		Start:    req.Start,
		End:      req.End,
		Duration: ((req.End - req.Start) / MilisecondInDay),
		UpdateAt: time.Now(),
	}

	if err := s.repo.Update(ctx, info, id); err != nil {
		return nil, fmt.Errorf("Faild to update Holiday, %w", err)
	}

	return info, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	//check policy
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		logrus.Errorf("Fail to delete holiday due to %v", err)
		return status.Holiday().NotFoundHoliday
	}
	return nil
}

func (s *Service) FindAll(ctx context.Context) ([]*types.Holiday, error) {

	user := auth.FromContext(ctx)
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}

	var holidays []*types.Holiday

	holidays, err := s.repo.FindAll(ctx, user.UserID)
	if err != nil {
		logrus.Errorf("Database err: failed to find holidays, err: %v", err)
		return nil, fmt.Errorf("Failed to find holidays of this user: %w", err)
	}

	info := make([]*types.Holiday, 0)
	for _, holiday := range holidays {
		info = append(info, &types.Holiday{
			Title:      holiday.Title,
			Start:      holiday.Start,
			End:        holiday.End,
			Duration:   holiday.Duration,
			HolidayID:  holiday.HolidayID,
			CreatedAt:  holiday.CreatedAt,
			CreaterID:  holiday.CreaterID,
			UpdateAt:   holiday.UpdateAt,
			ProjectsID: holiday.ProjectsID,
		})
	}
	return info, err
}

func (s *Service) FindByProjectID(ctx context.Context, projectID string) ([]*types.HolidayInfo, error) {

	holidays, err := s.repo.FindByProjectID(ctx, projectID)

	info := make([]*types.HolidayInfo, 0)
	for _, holiday := range holidays {
		info = append(info, &types.HolidayInfo{
			HolidayID: holiday.HolidayID,
			Title:     holiday.Title,
			Start:     holiday.Start,
			End:       holiday.End,
			Duration:  holiday.Duration,
			UpdateAt:  holiday.UpdateAt,
		})
	}
	return info, err
}

func (s *Service) FindByID(ctx context.Context, id string) (string, error) {

	_, err := s.repo.FindByID(ctx, id)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.Error("Failed to check existing holiday by id %w", err)
		return "", fmt.Errorf("Failed to check existing holiday by id: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Error("Holiday all ready exist")
		return "", status.Holiday().NotFoundHoliday
	}

	return id, nil
}

func (s *Service) AssignProject(ctx context.Context, holidayID string, projectID string) error {

	holiday, err := s.repo.FindByID(ctx, holidayID)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.Error("Failed to check existing holiday by id %w", err)
		return fmt.Errorf("Failed to check existing holiday by id: %w", err)
	}

	if db.IsErrNotFound(err) {
		logrus.Error("Holiday not found")
		return status.Holiday().NotFoundHoliday
	}

	hasProject := false
	for _, project := range holiday.ProjectsID {
		if project == projectID {
			hasProject = true
		}
	}
	if hasProject {
		logrus.Errorf("Project already had this holiday")
		return status.Holiday().AlreadyInProject
	}
	return s.repo.UpdateProjectsID(ctx, holidayID, projectID, true)
}

func (s *Service) RemoveProject(ctx context.Context, holidayID string, projectID string) error {

	holidays, err := s.repo.FindByProjectID(ctx, projectID)
	if holidayID == "_all_holiday_" {
		//remove project from all holidays
		for _, holiday := range holidays {
			if err := s.repo.UpdateProjectsID(ctx, holiday.HolidayID, projectID, false); err != nil {
				logrus.Error("Database error, Failed to remove project from holiday %w", err)
				return fmt.Errorf("Failed to remove project from holiday: %w", err)
			}
		}
	} else {
		//remove project from one holiday
		//check holiday exist, is projectID in this holiday?
		holiday, err := s.repo.FindByID(ctx, holidayID)
		if err != nil && !db.IsErrNotFound(err) {
			logrus.Error("Failed to check existing holiday by id %w", err)
			return fmt.Errorf("Failed to check existing holiday by id: %w", err)
		}

		if db.IsErrNotFound(err) {
			logrus.Error("Holiday not found")
			return status.Holiday().NotFoundHoliday
		}

		hasProject := false
		for _, project := range holiday.ProjectsID {
			if project == projectID {
				hasProject = true
			}
		}
		if !hasProject {
			logrus.Errorf("Project didn't have this holiday")
			return status.Holiday().NotFoundProject
		}
		if err := s.repo.UpdateProjectsID(ctx, holidayID, projectID, false); err != nil {
			logrus.Error("Database error, Failed to remove project from holiday %w", err)
			return fmt.Errorf("Failed to remove project from holiday: %w", err)
		}
	}
	return err
}
