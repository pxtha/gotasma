package types

import (
	"time"
)

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
					"workloads"  :{ "type":"nested" },
					"devs"		 :{ "type":"nested" },
					"tasks"		 :{ "type":"nested"	 },
					"holidays"	 :{	"type":"keyword" },
					"created_at" :{ "type":"date" 	 },
					"updated_at" :{ "type":"date"	 },
					"highlight"	 :{ "type":"boolean" }
				}
			}
	}`
)

type (
	History struct {
		Name      string         `json:"name,omitempty" `
		Desc      string         `json:"desc,omitempty"`
		ProjectID string         `json:"project_id,omitempty"`
		CreaterID string         `json:"creater_id,omitempty"`
		WorkLoad  *WorkLoadInfo  `json:"workloads,omitempty"`
		Devs      []*UserInfo    `json:"devs,omitempty"`
		Tasks     []*Task        `json:"tasks,omitempty"`
		Holiday   []*HolidayInfo `json:"holiday,omitempty"`
		Action    string         `json:"action,omitempty"`
		CreatedAt time.Time      `json:"created_at,omitempty"`
		Highlight bool           `json:"highlight"`
	}
)
