package utils

import (
	"TapTransit-backend/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// GetMonthlyAggregate 获取卡片月度累计金额（从数据库）
func GetMonthlyAggregate(db *gorm.DB, cardID string, month string) (float64, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	var aggregate models.MonthlyAggregate
	err := db.Where("card_id = ? AND month = ?", cardID, month).First(&aggregate).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil // 如果没有记录，返回0
	}
	if err != nil {
		return 0, fmt.Errorf("查询月度累计失败: %w", err)
	}

	return aggregate.TotalAmount, nil
}

// IncrementMonthlyAggregate 增加卡片月度累计金额（使用数据库）
func IncrementMonthlyAggregate(db *gorm.DB, cardID string, amount float64) error {
	month := time.Now().Format("2006-01")

	// 使用ON CONFLICT UPDATE或先查询后更新
	var aggregate models.MonthlyAggregate
	err := db.Where("card_id = ? AND month = ?", cardID, month).First(&aggregate).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		aggregate = models.MonthlyAggregate{
			CardID:      cardID,
			Month:       month,
			TotalAmount: amount,
			UpdatedAt:   time.Now(),
		}
		return db.Create(&aggregate).Error
	} else if err != nil {
		return fmt.Errorf("查询月度累计失败: %w", err)
	}

	// 更新现有记录
	aggregate.TotalAmount += amount
	aggregate.UpdatedAt = time.Now()
	return db.Save(&aggregate).Error
}

// GetCurrentMonthAggregate 获取当前月份累计金额
func GetCurrentMonthAggregate(db *gorm.DB, cardID string) (float64, error) {
	return GetMonthlyAggregate(db, cardID, "")
}
