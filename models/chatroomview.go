package models

import "time"

type ChatroomView struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	ProjectID  string    `json:"project_id"`
	CaseID     string    `json:"case_id"`
	UUID       string    `json:"uuid"`
	CreateDate time.Time `json:"create_date"`
	CountRead  int64       `json:"count_read"`
}
