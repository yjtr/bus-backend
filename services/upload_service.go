package services

import (
	"awesomeProject/models"
	"awesomeProject/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UploadService struct {
	db            *gorm.DB
	fareService   *FareService
}

func NewUploadService(db *gorm.DB, fareService *FareService) *UploadService {
	return &UploadService{
		db:          db,
		fareService: fareService,
	}
}

// BatchRecordRequest 网关上传的批量记录请求
type BatchRecordRequest struct {
	CardID        string    `json:"card_id" binding:"required"`
	BoardTime     time.Time `json:"board_time" binding:"required"`
	BoardStation  string    `json:"board_station" binding:"required"`
	AlightTime    time.Time `json:"alight_time"`
	AlightStation string    `json:"alight_station"`
	RouteID       uint      `json:"route_id"`
	GatewayID     string    `json:"gateway_id"`
}

// UploadBatchRecords 批量上传乘车记录
func (s *UploadService) UploadBatchRecords(records []BatchRecordRequest) (int, error) {
	var successCount int

	for _, record := range records {
		if err := s.processSingleRecord(record); err != nil {
			// 记录错误但继续处理下一条
			fmt.Printf("处理记录失败: %v, 错误: %v\n", record, err)
			continue
		}
		successCount++
	}

	return successCount, nil
}

// processSingleRecord 处理单条记录
func (s *UploadService) processSingleRecord(record BatchRecordRequest) error {
	// 解析站点信息（格式：线路ID-站点名称 或 站点ID）
	startStationID, startStationName, err := s.parseStation(record.BoardStation)
	if err != nil {
		return fmt.Errorf("解析上车站点失败: %w", err)
	}

	endStationID := uint(0)
	endStationName := ""
	if record.AlightStation != "" {
		endStationID, endStationName, err = s.parseStation(record.AlightStation)
		if err != nil {
			return fmt.Errorf("解析下车站点失败: %w", err)
		}
	}

	// 如果没有路由ID，尝试从站点推断
	routeID := record.RouteID
	if routeID == 0 {
		routeID, _ = s.inferRouteFromStation(startStationID)
	}

	// 检查卡片是否存在，不存在则创建
	var card models.Card
	err = s.db.Where("card_id = ?", record.CardID).First(&card).Error
	if err == gorm.ErrRecordNotFound {
		card = models.Card{
			CardID:   record.CardID,
			CardType: "normal",
			Status:   "active",
		}
		if err := s.db.Create(&card).Error; err != nil {
			return fmt.Errorf("创建卡片失败: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("查询卡片失败: %w", err)
	}

	// 检查卡片状态
	if card.Status != "active" {
		return fmt.Errorf("卡片状态异常: %s", card.Status)
	}

	// 创建交易记录
	transaction := models.Transaction{
		CardID:          record.CardID,
		RouteID:         routeID,
		StartStation:    startStationID,
		EndStation:      endStationID,
		StartStationName: startStationName,
		EndStationName:   endStationName,
		BoardTime:       record.BoardTime,
		GatewayID:       record.GatewayID,
		Status:          "pending",
	}

	// 如果有下车时间，计算费用
	if !record.AlightTime.IsZero() && endStationID > 0 {
		transaction.AlightTime = &record.AlightTime

		// 计算费用
		fareResult, err := s.fareService.CalculateFare(
			record.CardID,
			routeID,
			startStationID,
			endStationID,
			record.BoardTime,
		)
		if err != nil {
			return fmt.Errorf("计算费用失败: %w", err)
		}

		transaction.Fare = fareResult.BaseFare
		transaction.ActualFare = fareResult.ActualFare
		transaction.DiscountType = fareResult.DiscountType
		transaction.DiscountAmount = fareResult.DiscountAmount
		transaction.Status = "completed"

		// 更新Redis中的月度累计金额
		if err := utils.IncrementCardMonthlyAmount(record.CardID, fareResult.ActualFare); err != nil {
			// Redis错误不影响主流程，只记录日志
			fmt.Printf("更新Redis月度累计失败: %v\n", err)
		}

		// 删除上车信息缓存（已下车）
		utils.DeleteCardOnboardInfo(record.CardID)
	} else {
		// 只有上车记录，在Redis中缓存上车信息（用于换乘判断）
		onboardInfo := fmt.Sprintf("%d:%d:%d", routeID, startStationID, record.BoardTime.Unix())
		utils.SetCardOnboardInfo(record.CardID, onboardInfo, 2*time.Hour)
	}

	// 保存交易记录
	if err := s.db.Create(&transaction).Error; err != nil {
		return fmt.Errorf("保存交易记录失败: %w", err)
	}

	return nil
}

// parseStation 解析站点信息（简化版，实际可能需要更复杂的解析逻辑）
func (s *UploadService) parseStation(stationInfo string) (uint, string, error) {
	// 尝试查找站点（通过station_id或name匹配）
	var station models.Station
	err := s.db.Where("station_id = ? OR name = ?", stationInfo, stationInfo).First(&station).Error
	if err == nil {
		return station.ID, station.Name, nil
	}

	// 如果找不到，返回错误（调用方会使用原始字符串作为名称）
	return 0, stationInfo, fmt.Errorf("站点不存在: %s", stationInfo)
}

// inferRouteFromStation 从站点推断线路
func (s *UploadService) inferRouteFromStation(stationID uint) (uint, error) {
	var routeStation models.RouteStation
	err := s.db.Where("station_id = ?", stationID).First(&routeStation).Error
	if err != nil {
		return 0, err
	}
	return routeStation.RouteID, nil
}
