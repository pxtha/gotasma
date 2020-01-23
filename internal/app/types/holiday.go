package types

import "time"

type (
	Holiday struct {
		Title     string    `json:"title,omitempty" bson:"title,omitempty"`
		HolidayID string    `json:"holiday_id,omitempty" bson:"holiday_id,omitempty"`
		Start     int       `json:"start,omitempty" bson:"start,omitempty"`
		End       int       `json:"end,omitempty" bson:"end,omitempty"`
		Duration  int       `json:"end,omitempty" bson:"end,omitempty"`
		ProjectID []string  `json:"project_id" bson:"project_id"`
		CreaterID string    `json:"creater_id,omitempty" bson:"creater_id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdateAt  time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	}
	HolidayRequest struct {
		Title string `json:"title,omitempty"  validate:"required, gt=3" bson:"title,omitempty"`
		Start int    `json:"start,omitempty" validate:"required, gt=1262304000000" bson:"start,omitempty"`
		End   int    `json:"end,omitempty" validate:"required, gt=1262304000000" bson:"end,omitempty"`
	}
)
