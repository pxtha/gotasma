package types

import "time"

type (
	Interval struct {
		From time.Time
		To   time.Time
	}

	WorkLoad struct {
		WorkLoadID string `json:"workload_id" bson:"workload_id"`
		UserID     string `json:"user_id" bson:"user_id"`
		ProjectID  string `json:"project_id,omitempty" bson:"project_id"`
		//Overload: From ... to ....
		Overload  map[int]int `json:"overload,omitempty" bson:"overload"`
		CreatedAt time.Time   `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdateAt  time.Time   `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	}

	WorkLoadInfo struct {
		UserID    string      `json:"user_id" bson:"user_id"`
		ProjectID string      `json:"project_id,omitempty" bson:"project_id"`
		Overload  []*Interval `json:"overload,omitempty" bson:"overload"`
	}
)
