package controllers

import (
	"awesomeProject/services"
	"awesomeProject/utils"

	"github.com/gin-gonic/gin"
)

type CardController struct {
	cardService *services.CardService
}

func NewCardController(cardService *services.CardService) *CardController {
	return &CardController{
		cardService: cardService,
	}
}

// GetCard 查询卡片信息
// @Summary 查询卡片信息
// @Description 根据卡ID查询卡片状态和信息
// @Tags 卡片管理
// @Produce json
// @Param id path string true "卡片ID"
// @Success 200 {object} models.Card
// @Router /api/v1/card/{id} [get]
func (c *CardController) GetCard(ctx *gin.Context) {
	cardID := ctx.Param("id")
	if cardID == "" {
		utils.BadRequest(ctx, "缺少卡片ID")
		return
	}

	card, err := c.cardService.GetCardByID(cardID)
	if err != nil {
		utils.NotFound(ctx, "卡片不存在")
		return
	}

	utils.Success(ctx, card)
}
