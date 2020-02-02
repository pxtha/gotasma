package api

import (
	"praslar.com/gotasma/internal/app/holiday"
)

func newHolidayService(policy holiday.PolicyServices) (*holiday.Services, error) {
	s, err := dialDefaultMongoDB()
	if err != nil {
		return nil, err
	}

	repo := holiday.NewMongoRepository(s)
	return holiday.New(repo, policy), nil
}

func newHolidayHandler(srv *holiday.Services) *holiday.Handler {
	return holiday.NewHanlder(srv)
}
