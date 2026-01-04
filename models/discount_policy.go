package models

import (
	"time"

	"gorm.io/gorm"
)

// DiscountPolicy 其他优惠策略
type DiscountPolicy struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	PolicyName string  `gorm:"size:100;not null" json:"policy_name"`        // 策略名称，如"月累计折扣"
	PolicyType string  `gorm:"size:50;not null" json:"policy_type"`         // 策略类型：monthly_accumulate, student, elder, etc.
	Threshold  float64 `gorm:"type:decimal(10,2);default:0" json:"threshold"` // 阈值（如月累计200元）
	DiscountRate float64 `gorm:"type:decimal(5,4);default:0" json:"discount_rate"` // 折扣比例（0-1之间）
	DiscountAmount float64 `gorm:"type:decimal(10,2);default:0" json:"discount_amount"` // 固定优惠金额
	CardTypeFilter string `gorm:"size:50" json:"card_type_filter"`          // 适用的卡类型（空表示所有）
	Status    string  `gorm:"size:20;default:'active'" json:"status"`       // 状态：active, inactive
}

// TableName 指定表名
func (DiscountPolicy) TableName() string {
	return "discount_policies"
}
