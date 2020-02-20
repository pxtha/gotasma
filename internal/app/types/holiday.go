package types

import "time"

type (
	//TODO type of holiday, if type: public => apply to all project
	Holiday struct {
		Title      string    `json:"title,omitempty" bson:"title,omitempty"`
		HolidayID  string    `json:"holiday_id,omitempty" bson:"holiday_id,omitempty"`
		Start      int       `json:"start,omitempty" bson:"start,omitempty"`
		End        int       `json:"end,omitempty" bson:"end,omitempty"`
		ProjectsID []string  `json:"projects_id,omitempty" bson:"projects_id"`
		Duration   int       `json:"duration,omitempty" bson:"duration,omitempty"`
		CreaterID  string    `json:"creater_id,omitempty" bson:"creater_id,omitempty"`
		CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdateAt   time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	}

	HolidayInfo struct {
		HolidayID string    `json:"holiday_id,omitempty"`
		Title     string    `json:"title,omitempty" `
		Start     int       `json:"start,omitempty" `
		End       int       `json:"end,omitempty" `
		Duration  int       `json:"duration,omitempty"`
		UpdateAt  time.Time `json:"updated_at,omitempty" `
	}

	HolidayRequest struct {
		Title string `json:"title,omitempty"  validate:"required,gt=3"`
		Start int    `json:"start,omitempty" validate:"required,gt=1262304000000"`
		End   int    `json:"end,omitempty" validate:"required,gt=1262304000000"`
	}
)
