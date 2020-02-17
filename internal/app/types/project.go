package types

import "time"

const (
	Mapping = `
	{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"projects": {
				"properties":{
					"name"		 :{ "type":"keyword" },
					"description":{ "type":"text"	 },
					"project_id" :{ "type":"keyword" },
					"creater_id" :{ "type":"keyword" },
					"devs_id"	 :{ "type":"keyword" },
					"tasks"		 :{ "type":"nested"	 },
					"created_at" :{ "type":"date" 	 },
					"updated_at" :{ "type":"date"	 },
					"highlight"	 :{ "type":"boolean"}
				}
			}
		}
	}`
)

type (
	Task struct {
		Label            string    `json:"label,omitempty" bson:"label" validate:"required"`
		TaskID           string    `json:"task_id,omitempty" bson:"task_id" validate:"required"`
		Start            int       `json:"start,omitempty" bson:"start" validate:"gt=1262304000000"`
		Duration         int       `json:"duration,omitempty" bson:"duration"`
		EstimateDuration int       `json:"estimate_duration,omitempty" bson:"estimate_duration"`
		End              int       `json:"end,omitempty" bson:"end"  validate:"gt=1262304000000"`
		Parent           string    `json:"parent,omitempty" bson:"parent"`
		Parents          []string  `json:"parents,omitempty" bson:"parents"`
		Children         []string  `json:"children,omitempty" bson:"children"`
		AllChildren      []string  `json:"all_children,omitempty" bson:"all_children"`
		Effort           int       `json:"effort,omitempty" bson:"effort"`
		DevsID           []string  `json:"devs_id,omitempty" bson:"dev_id"`
		Type             string    `json:"type,omitempty" bson:"type" validate:"required"`
		UpdateAt         time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	}

	Project struct {
		Name      string    `json:"name,omitempty" bson:"name"`
		Desc      string    `json:"desc,omitempty" bson:"desc"`
		ProjectID string    `json:"project_id,omitempty" bson:"project_id"`
		CreaterID string    `json:"creater_id,omitempty" bson:"creater_id"`
		DevsID    []string  `json:"devs_id,omitempty" bson:"devs_id"`
		Tasks     []*Task   `json:"tasks,omitempty" bson:"tasks"`
		CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
		UpdateAt  time.Time `json:"updated_at,omitempty" bson:"updated_at"`
		Highlight bool      `json:"highlight,omitempty" bson:"highlight"`
	}

	ProjectInfo struct {
		Name      string    `json:"name,omitempty" `
		Desc      string    `json:"description,omitempty" `
		ProjectID string    `json:"project_id,omitempty" `
		CreaterID string    `json:"creater_id,omitempty"`
		DevsID    []string  `json:"devs_id,omitempty" `
		Tasks     int       `json:"tasks,omitempty" `
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdateAt  time.Time `json:"updated_at,omitempty" `
		Highlight bool      `json:"highlight,omitempty" `
	}

	CreateProjectRequest struct {
		Name string `json:"name,omitempty"  validate:"required,gt=3"`
		Desc string `json:"desc,omitempty"`
	}

	UpdateProject struct {
		Name  string  `json:"name,omitempty" validate:"required"`
		Tasks []*Task `json:"tasks,omitempty" validate:"required"`
	}

	RemoveUserRequest struct {
		UserID string `json:"user_id"`
	}

	AddUsersRequest struct {
		UserIDs []string `json:"user_ids"`
	}
)
