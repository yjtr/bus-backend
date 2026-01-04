package controllers

import (
	"awesomeProject/models"
	"awesomeProject/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConfigController struct {
}

func NewConfigController() *ConfigController {
	return &ConfigController{}
}

// GetRouteConfig 获取线路配置信息
// @Summary 获取线路配置
// @Description 获取线路的站点列表、票价表、换乘优惠等信息
// @Tags 配置
// @Produce json
// @Param route_id query int true "线路ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/bus/config [get]
func (c *ConfigController) GetRouteConfig(ctx *gin.Context) {
	routeIDStr := ctx.Query("route_id")
	if routeIDStr == "" {
		utils.BadRequest(ctx, "缺少route_id参数")
		return
	}

	routeID, err := strconv.ParseUint(routeIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "route_id参数格式错误")
		return
	}

	var route models.Route
	if err := utils.DB.Where("id = ?", routeID).First(&route).Error; err != nil {
		utils.NotFound(ctx, "线路不存在")
		return
	}

	// 获取票价规则
	var fares []models.Fare
	utils.DB.Where("route_id = ? AND status = 'active'", routeID).Find(&fares)

	// 获取换乘优惠规则
	var transfers []models.Transfer
	utils.DB.Where("(from_route_id = ? OR to_route_id = ?) AND status = 'active'", routeID, routeID).Find(&transfers)

	// 获取站点列表（通过RouteStation关联）
	var routeStations []models.RouteStation
	utils.DB.Where("route_id = ?", routeID).
		Preload("Station").
		Order("sequence ASC").
		Find(&routeStations)

	stations := make([]map[string]interface{}, 0)
	for _, rs := range routeStations {
		stations = append(stations, map[string]interface{}{
			"id":          rs.Station.ID,
			"station_id":  rs.Station.StationID,
			"name":        rs.Station.Name,
			"sequence":    rs.Sequence,
			"is_transfer": rs.Station.IsTransfer,
		})
	}

	utils.Success(ctx, gin.H{
		"route_id":   route.ID,
		"route_name": route.Name,
		"stations":   stations,
		"fares":      fares,
		"transfers":  transfers,
	})
}
