package models

import (
	"database/sql/driver"

	"github.com/bytedance/sonic"
)

type JSONMap map[string]any

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	b, err := sonic.Marshal(j)
	return string(b), err
}

func (j *JSONMap) Scan(src any) error {
	if src == nil {
		*j = JSONMap{}
		return nil
	}
	switch v := src.(type) {
	case []byte:
		return sonic.Unmarshal(v, j)
	case string:
		return sonic.Unmarshal([]byte(v), j)
	default:
		*j = JSONMap{}
		return nil
	}
}

type ActivityLog struct {
	BaseModelWithoutUpdatedAt
	LogName     string  `gorm:"size:255;column:log_name;index:idx_log_name" json:"log_name,omitzero"`
	Desc        string  `gorm:"type:text;not null;column:desc" json:"desc,omitzero"`
	SubjectType string  `gorm:"size:255;column:subject_type;index:idx_subject" json:"subject_type,omitzero"`
	Event       string  `gorm:"size:255;column:event" json:"event,omitzero"`
	SubjectID   *string `gorm:"column:subject_id;index:idx_subject" json:"subject_id,omitzero"`
	CauserType  string  `gorm:"size:255;column:causer_type;index:idx_causer" json:"causer_type,omitzero"`
	CauserID    *string `gorm:"column:causer_id;index:idx_causer" json:"causer_id,omitzero"`
	Properties  JSONMap `gorm:"type:json;column:properties" json:"properties,omitzero"`
	TraceID     *string `gorm:"size:36;column:trace_id" json:"trace_id,omitzero"`
	IP          *string `gorm:"size:45;column:ip" json:"ip,omitzero"`
	UserAgent   *string `gorm:"size:255;column:user_agent" json:"user_agent,omitzero"`
	Subject     any     `gorm:"-" json:"subject,omitzero"`
	Causer      any     `gorm:"-" json:"causer,omitzero"`
}
