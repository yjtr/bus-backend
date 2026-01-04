package models

import (
	"time"

	"gorm.io/gorm"
)

// Route 公交线路信息
type Route struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	RouteID   string `gorm:"uniqueIndex;not null;size:50" json:"route_id"` // 线路编号，如"A"、"1"、"K1"
	Name      string `gorm:"size:100;not null" json:"name"`                // 线路名称
	Status    string `gorm:"size:20;default:'active'" json:"status"`       // 状态：active, inactive
}

// TableName 指定表名
func (Route) TableName() string {
	return "routes"
}
