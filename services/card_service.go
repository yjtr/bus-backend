package services

import (
	"TapTransit-backend/models"

	"gorm.io/gorm"
)

type CardService struct {
	db *gorm.DB
}

type CardFilter struct {
	CardID     string
	CardNoLike string
	HolderName string
	Status     string
}

type CardDiscount struct {
	DiscountRate   float64
	DiscountAmount float64
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

// ListCards 查询卡片列表（支持简单筛选）
func (s *CardService) ListCards(filter CardFilter) ([]models.Card, error) {
	query := s.db.Model(&models.Card{})
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.CardID != "" {
		query = query.Where("card_id = ?", filter.CardID)
	} else if filter.CardNoLike != "" {
		query = query.Where("card_id LIKE ?", "%"+filter.CardNoLike+"%")
	}
	if filter.HolderName != "" {
		query = query.Where("holder_name LIKE ?", "%"+filter.HolderName+"%")
	}
	var cards []models.Card
	if err := query.Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

func (s *CardService) GetCardDiscount(cardType string) (CardDiscount, error) {
	if cardType == "" || cardType == "normal" {
		return CardDiscount{}, nil
	}
	var policy models.DiscountPolicy
	err := s.db.Where(
		"policy_type = ? AND (card_type_filter = ? OR card_type_filter = '') AND status = 'active'",
		cardType,
		cardType,
	).First(&policy).Error
	if err != nil {
		return CardDiscount{}, err
	}
	return CardDiscount{
		DiscountRate:   policy.DiscountRate,
		DiscountAmount: policy.DiscountAmount,
	}, nil
}
