package controllers

import (
	"awesomeProject/services"
	"awesomeProject/utils"

	"github.com/gin-gonic/gin"
)

type BusController struct {
	uploadService *services.UploadService
}

func NewBusController(uploadService *services.UploadService) *BusController {
	return &BusController{
		uploadService: uploadService,
	}
}

// UploadBatchRecords 批量上传乘车记录
// @Summary 批量上传乘车记录
// @Description 网关上传批量乘车记录
// @Tags 公交数据
// @Accept json
// @Produce json
// @Param records body []services.BatchRecordRequest true "批量记录"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/bus/batchRecords [post]
func (c *BusController) UploadBatchRecords(ctx *gin.Context) {
	var records []services.BatchRecordRequest
	if err := ctx.ShouldBindJSON(&records); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	successCount, err := c.uploadService.UploadBatchRecords(records)
	if err != nil {
		utils.InternalServerError(ctx, "处理记录失败: "+err.Error())
		return
	}

	utils.Success(ctx, gin.H{
		"received": successCount,
		"total":    len(records),
	})
}
