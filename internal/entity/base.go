package enetity

import "time"

type BaseModel struct {
	CreateTime *int64     `gorm:"column:create_time;index" json:"create_time,omitempty"`
	CreateDate *time.Time `gorm:"column:create_date;index" json:"create_date,omitempty"`
	UpdateTime *int64     `gorm:"column:update_time;index" json:"update_time,omitempty"`
	UpdateDate *time.Time `gorm:"column:update_date;index" json:"update_date,omitempty"`
}
