package models

import (
	"time"

	"gorm.io/gorm"
)

// Fare 基本票价规则
type Fare struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	RouteID      uint    `gorm:"index" json:"route_id"`                      // 线路ID（0表示通用规则）
	StartStation uint    `gorm:"index" json:"start_station"`                 // 起始站点ID（0表示通用）
	EndStation   uint    `gorm:"index" json:"end_station"`                   // 结束站点ID（0表示通用）
	BasePrice    float64 `gorm:"type:decimal(10,2);not null" json:"base_price"` // 基础票价
	FareType     string  `gorm:"size:50;default:'uniform'" json:"fare_type"`    // 计价类型：uniform(统一票价), segment(分段计价), distance(按距离)
	SegmentCount int     `gorm:"default:0" json:"segment_count"`                // 区段数（用于分段计价）
	ExtraPrice   float64 `gorm:"type:decimal(10,2);default:0" json:"extra_price"` // 续程价（分段计价用）
	Status       string  `gorm:"size:20;default:'active'" json:"status"`         // 状态：active, inactive
}

// TableName 指定表名
func (Fare) TableName() string {
	return "fares"
}
