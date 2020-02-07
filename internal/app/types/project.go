package types

import "time"

type (
	Task struct {
		Label            string   `json:"label,omitempty" bson:"label" validate:"required"`
		TaskID           string   `json:"task_id,omitempty" bson:"task_id" validate:"required"`
		Start            int      `json:"start,omitempty" bson:"start" validate:"gt=1262304000000"`
		Duration         int      `json:"duration,omitempty" bson:"duration"`
		EstimateDuration int      `json:"estimate_duration,omitempty" bson:"estimate_duration"`
		End              int      `json:"end,omitempty" bson:"end"  validate:"gt=1262304000000"`
		Parent           string   `json:"parent,omitempty" bson:"parent"`
		Parents          []string `json:"parents,omitempty" bson:"parents"`
		Children         []string `json:"children,omitempty" bson:"children"`
		AllChildren      []string `json:"all_children,omitempty" bson:"all_children"`
		Effort           int      `json:"effort,omitempty" bson:"effort"`
		DevsID           []string `json:"devs_id,omitempty" bson:"dev_id"`
		Type             string   `json:"type,omitempty" bson:"type" validate:"required"`
	}

	Project struct {
		Name      string    `json:"name,omitempty" bson:"name"`
		ProjectID string    `json:"project_id,omitempty" bson:"project_id"`
		CreaterID string    `json:"creater_id,omitempty" bson:"creater_id"`
		DevsID    []string  `json:"devs_id,omitempty" bson:"devs_id"`
		Tasks     []Task    `json:"tasks,omitempty" bson:"tasks"`
		CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
		UpdateAt  time.Time `json:"updated_at,omitempty" bson:"updated_at"`
		Highlight bool      `json:"highlight,omitempty" bson:"highlight"`
	}

	CreateProjectRequest struct {
		Name string `json:"name,omitempty"  validate:"required,gt=3"`
	}

	ProjectInfo struct {
		Name  string `json:"name,omitempty" validate:"required"`
		Tasks []Task `json:"tasks,omitempty" validate:"required"`
	}

	RemoveUserRequest struct {
		UserID string `json:"user_id"`
	}

	AddUsersRequest struct {
		UserIDs []string `json:"user_ids"`
	}
)
