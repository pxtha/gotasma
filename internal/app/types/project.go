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
				"properties":{
					"name"		 :{ "type":"keyword" },
					"desc"		 :{ "type":"text"	 },
					"project_id" :{ "type":"keyword" },
					"creater_id" :{ "type":"keyword" },
					"devs_id"	 :{ "type":"keyword" },
					"tasks"		 :{ "type":"nested"	 },
					"created_at" :{ "type":"date" 	 },
					"updated_at" :{ "type":"date"	 },
					"highlight"	 :{ "type":"boolean" },
					"holidays_id":{	"type":"keyword" }
				}
			}
	}`
)

type (
	Project struct {
		Name      string    `json:"name,omitempty" bson:"name"`
		Desc      string    `json:"desc,omitempty" bson:"desc"`
		ProjectID string    `json:"project_id,omitempty" bson:"project_id"`
		CreaterID string    `json:"creater_id,omitempty" bson:"creater_id"`
		CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
		UpdateAt  time.Time `json:"updated_at,omitempty" bson:"updated_at"`
		Highlight bool      `json:"highlight,omitempty" bson:"highlight"`
	}

	ProjectHistory struct {
		Name      string    `json:"name,omitempty" bson:"name"`
		Desc      string    `json:"desc,omitempty" bson:"desc"`
		ProjectID string    `json:"project_id,omitempty" bson:"project_id"`
		CreaterID string    `json:"creater_id,omitempty" bson:"creater_id"`
		Action    string    `json:"action,omitempty" bson:"action"`
		CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
		UpdateAt  time.Time `json:"updated_at,omitempty" bson:"updated_at"`
		Highlight bool      `json:"highlight" bson:"highlight"`
	}

	ProjectInfo struct {
		Name      string    `json:"name,omitempty" `
		Desc      string    `json:"description,omitempty" `
		ProjectID string    `json:"project_id,omitempty" `
		CreaterID string    `json:"creater_id,omitempty"`
		Tasks     int       `json:"tasks,omitempty" `
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdateAt  time.Time `json:"updated_at,omitempty" `
		Highlight bool      `json:"highlight,omitempty" `
	}

	CreateProjectRequest struct {
		Name string `json:"name,omitempty"  validate:"required,gt=3"`
		Desc string `json:"desc,omitempty" `
	}

	SaveProject struct {
		Tasks []*Task `json:"tasks,omitempty" validate:"required"`
	}

	UpdateProject struct {
		Name      string `json:"name" validate:"required" `
		Desc      string `json:"desc"  validate:"required"`
		Highlight bool   `json:"highlight" `
	}

	RemoveUserRequest struct {
		UserID string `json:"user_id" validate:"required"`
	}

	AddUsersRequest struct {
		UserID string `json:"user_id" validate:"required"`
	}

	RemoveHolidayRequest struct {
		HolidayID string `json:"holiday_id" validate:"required"`
	}

	AddHolidayRequest struct {
		HolidayID string `json:"holiday_id" validate:"required" `
	}

	UnAssignDev struct {
		TaskID string `json:"task_id" validate:"required"`
	}

	AssignDev struct {
		TaskID string `json:"task_id" validate:"required" `
		UserID string `json:"user_id" validate:"required" `
	}
)
