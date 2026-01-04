package services

import (
	"awesomeProject/models"

	"gorm.io/gorm"
)

type CardService struct {
	db *gorm.DB
}

func NewCardService(db *gorm.DB) *CardService {
	return &CardService{db: db}
}

// GetCardByID 根据卡ID查询卡片信息
func (s *CardService) GetCardByID(cardID string) (*models.Card, error) {
	var card models.Card
	err := s.db.Where("card_id = ?", cardID).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// GetCardStatus 获取卡片状态
func (s *CardService) GetCardStatus(cardID string) (string, error) {
	card, err := s.GetCardByID(cardID)
	if err != nil {
		return "", err
	}
	return card.Status, nil
}

// CreateCard 创建新卡片
func (s *CardService) CreateCard(card *models.Card) error {
	return s.db.Create(card).Error
}

// UpdateCard 更新卡片信息
func (s *CardService) UpdateCard(cardID string, updates map[string]interface{}) error {
	return s.db.Model(&models.Card{}).Where("card_id = ?", cardID).Updates(updates).Error
}

// BlockCard 封禁卡片
func (s *CardService) BlockCard(cardID string) error {
	return s.UpdateCard(cardID, map[string]interface{}{"status": "blocked"})
}

// UnblockCard 解封卡片
func (s *CardService) UnblockCard(cardID string) error {
	return s.UpdateCard(cardID, map[string]interface{}{"status": "active"})
}
