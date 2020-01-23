package holiday

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

const (
	MilisecondInDay = 86400000
)

type (
	Repository interface {
		Create(ctx context.Context, holiday *types.Holiday) error
		FindByTitle(ctx context.Context, title string) (*types.Holiday, error)
	}
	PolicyServices interface {
		Validate(ctx context.Context, obj string, act string) error
	}
	Services struct {
		repo   Repository
		policy PolicyServices
	}
)

func New(repo Repository, policy PolicyServices) *Services {
	return &Services{
		repo:   repo,
		policy: policy,
	}
}

func (s *Services) Create(ctx context.Context, req *types.HolidayRequest) (*types.Holiday, error) {
	if err := s.policy.Validate(ctx, types.ObjectHoliday, types.ActionHoliday); err != nil {
		return nil, err
	}
	if err := validator.Validate(req); err != nil {
		return nil, err
	}

	existingHoliday, err := s.repo.FindByTitle(ctx, req.Title)
	if err != nil && !db.IsErrNotFound(err) {
		return nil, fmt.Errorf("failed to check existing holiday by title: %w", err)
	}
	if existingHoliday != nil {
		return nil, status.Hoiday().DuplicatedHoliday
	}

	pm := auth.FromContext(ctx)
	holiday := &types.Holiday{
		Title:     req.Title,
		Start:     req.Start,
		End:       req.End,
		HolidayID: uuid.New(),
		Duration:  (req.End - req.Start) / MilisecondInDay,
		ProjectID: []string{},
		CreaterID: pm.UserID,
	}
	if err := s.repo.Create(ctx, holiday); err != nil {
		return nil, fmt.Errorf("Faild to insert Holiday, %w", err)
	}
	return holiday, nil
}
