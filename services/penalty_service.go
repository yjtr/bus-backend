package services

import (
	"TapTransit-backend/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PenaltyService struct {
	db          *gorm.DB
	fareService *FareService
}

func NewPenaltyService(db *gorm.DB, fareService *FareService) *PenaltyService {
	return &PenaltyService{
		db:          db,
		fareService: fareService,
	}
}

// ProcessPenaltyFares 处理缺少下车刷卡的pending交易（按罚款计费）
// 查找超过指定时间（默认30分钟）的pending交易，按罚款计费完成
func (s *PenaltyService) ProcessPenaltyFares(timeoutMinutes int) (int, error) {
	if timeoutMinutes <= 0 {
		timeoutMinutes = 30 // 默认30分钟
	}

	// 计算超时时间点
	timeoutTime := time.Now().Add(-time.Duration(timeoutMinutes) * time.Minute)

	// 查找所有超过超时时间的pending交易
	var pendingTransactions []models.Transaction
	err := s.db.Where("status = ? AND board_time < ?", "pending", timeoutTime).
		Find(&pendingTransactions).Error

	if err != nil {
		return 0, fmt.Errorf("查询pending交易失败: %w", err)
	}

	successCount := 0
	for _, transaction := range pendingTransactions {
		if err := s.processPenaltyFare(&transaction); err != nil {
			fmt.Printf("处理罚款计费失败 (Transaction ID: %d): %v\n", transaction.ID, err)
			continue
		}
		successCount++
	}

	return successCount, nil
}

// processPenaltyFare 处理单个pending交易的罚款计费
func (s *PenaltyService) processPenaltyFare(transaction *models.Transaction) error {
	// 获取线路信息
	var route models.Route
	if err := s.db.First(&route, transaction.RouteID).Error; err != nil {
		return fmt.Errorf("线路不存在: %w", err)
	}

	// 检查线路是否需要tap_in_out模式
	if route.TapMode != "tap_in_out" {
		// 如果不是tap_in_out模式，不应该有pending交易，跳过
		return nil
	}

	// 使用罚款计费逻辑计算费用
	fareResult, err := s.fareService.CalculateFareV2(
		transaction.CardID,
		transaction.RouteID,
		transaction.StartStation,
		nil, // 没有下车站点
		transaction.BoardTime,
		true, // 是罚款计费
	)
	if err != nil {
		return fmt.Errorf("计算罚款费用失败: %w", err)
	}

	// 更新交易记录
	transaction.Fare = fareResult.BaseFare
	transaction.ActualFare = fareResult.ActualFare
	transaction.DiscountType = fareResult.DiscountType
	transaction.DiscountAmount = fareResult.DiscountAmount
	transaction.PenaltyFare = fareResult.PenaltyFare
	transaction.Status = "completed"
	// EndStation保持为nil，AlightTime保持为nil（表示未下车）

	// 注意：罚款计费不计入月度累计（设计文档要求）
	// 所以这里不更新月度累计

	// 保存更新后的交易
	if err := s.db.Save(transaction).Error; err != nil {
		return fmt.Errorf("更新交易记录失败: %w", err)
	}

	return nil
}

// StartPenaltyProcessor 启动定时任务，定期处理罚款计费
// intervalMinutes: 检查间隔（分钟）
// timeoutMinutes: pending交易超时时间（分钟）
func (s *PenaltyService) StartPenaltyProcessor(intervalMinutes, timeoutMinutes int) {
	if intervalMinutes <= 0 {
		intervalMinutes = 5 // 默认每5分钟检查一次
	}
	if timeoutMinutes <= 0 {
		timeoutMinutes = 120 // 默认2小时（120分钟）超时
	}

	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	go func() {
		for range ticker.C {
			count, err := s.ProcessPenaltyFares(timeoutMinutes)
			if err != nil {
				fmt.Printf("处理罚款计费任务失败: %v\n", err)
			} else if count > 0 {
				fmt.Printf("处理了 %d 笔罚款计费交易\n", count)
			}
		}
	}()
}
