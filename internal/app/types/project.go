package types

import "time"

type (
	Task struct {
		Label            string   `json:"label,omitempty" bson:"label,omitempty"`
		TaskID           string   `json:"task_id,omitempty" bson:"task_id,omitempty"`
		Start            int      `json:"start,omitempty" bson:"start,omitempty"`
		Duration         int      `json:"duration,omitempty" bson:"duration,omitempty"`
		EstimateDuration int      `json:"estimate_duration,omitempty" bson:"estimate_duration,omitempty"`
		End              int      `json:"end,omitempty" bson:"end,omitempty"`
		Parent           string   `json:"parent,omitempty" bson:"parent,omitempty"`
		AllParent        []string `json:"all_parent,omitempty" bson:"all_parent,omitempty"`
		Children         []string `json:"children,omitempty" bson:"children,omitempty"`
		AllChildren      []string `json:"all_children,omitempty" bson:"all_children,omitempty"`
		Progress         int      `json:"progress,omitempty" bson:"progress,omitempty"`
		Effort           int      `json:"effort,omitempty" bson:"effort,omitempty"`
		DevsID           []string `json:"devs_id,omitempty" bson:"dev_id,omitempty"`
	}

	Project struct {
		Name      string    `json:"name,omitempty" bson:"name,omitempty"`
		ProjectID string    `json:"project_id,omitempty" bson:"project_id,omitempty"`
		CreaterID string    `json:"creater_id,omitempty" bson:"creater_id,omitempty"`
		DevsID    []string  `json:"devs_id,omitempty" bson:"devs_id,omitempty"`
		Tasks     []Task    `json:"tasks,omitempty" bson:"tasks,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdateAt  time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
		Highlight bool      `json:"highlight,omitempty" bson:"highlight,omitempty"`
	}

	CreateProjectRequest struct {
		Name string `json:"name,omitempty"  validate:"required,gt=3"`
	}

	RemoveUserRequest struct {
		UserID string `json:"user_id"`
	}

	AddUsersRequest struct {
		UserIDs []string `json:"user_ids"`
	}
)
