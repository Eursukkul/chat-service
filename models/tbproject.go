package models

import "time"

type Project struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	ProjectName string    `json:"project_name"`
	ProjectID   string    `json:"project_id"`
	SecretKey   string    `json:"secret_key"`
	CreateBy    string    `json:"create_by"`
	CreateDate  time.Time `json:"create_date"`
	ModifyBy    string    `json:"modify_by"`
	ModifyDate  time.Time `json:"modify_date"`
}
