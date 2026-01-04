package main

// 这是一个Go程序，用于在数据库迁移后初始化示例数据
// 使用方法：go run scripts/seed_data.go

import (
	"awesomeProject/config"
	"awesomeProject/models"
	"awesomeProject/utils"
	"fmt"
	"log"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	_, err = utils.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	db := utils.DB

	// 1. 创建示例线路
	routes := []models.Route{
		{RouteID: "A", Name: "1路", Area: "市区", Status: "active", Direction: "up"},
		{RouteID: "B", Name: "2路", Area: "市区", Status: "active", Direction: "up"},
		{RouteID: "K1", Name: "快速1路", Area: "市区", Status: "active", Direction: "up"},
	}

	for _, route := range routes {
		if err := db.FirstOrCreate(&route, models.Route{RouteID: route.RouteID}).Error; err != nil {
			log.Printf("创建线路失败: %v", err)
		}
	}

	// 2. 创建示例站点
	stations := []models.Station{
		{StationID: "ST001", Name: "市中心站", Latitude: 39.9042, Longitude: 116.4074, Address: "市中心广场", IsTransfer: true},
		{StationID: "ST002", Name: "火车站", Latitude: 39.9019, Longitude: 116.4250, Address: "火车站广场", IsTransfer: true},
		{StationID: "ST003", Name: "大学城站", Latitude: 39.9200, Longitude: 116.4000, Address: "大学城入口", IsTransfer: false},
		{StationID: "ST004", Name: "商业区站", Latitude: 39.9100, Longitude: 116.4200, Address: "商业区中心", IsTransfer: false},
		{StationID: "ST005", Name: "体育场站", Latitude: 39.9150, Longitude: 116.4100, Address: "体育场门口", IsTransfer: false},
		{StationID: "ST006", Name: "医院站", Latitude: 39.9050, Longitude: 116.4150, Address: "市人民医院", IsTransfer: false},
	}

	for _, station := range stations {
		if err := db.FirstOrCreate(&station, models.Station{StationID: station.StationID}).Error; err != nil {
			log.Printf("创建站点失败: %v", err)
		}
	}

	// 3. 创建线路-站点关联（1路线路）
	var route1 models.Route
	if err := db.Where("route_id = ?", "A").First(&route1).Error; err == nil {
		var stations1 []models.Station
		db.Where("station_id IN ?", []string{"ST001", "ST002", "ST003", "ST004"}).Find(&stations1)

		for i, station := range stations1 {
			routeStation := models.RouteStation{
				RouteID:   route1.ID,
				StationID: station.ID,
				Sequence:  i + 1,
				Direction: "up",
			}
			db.FirstOrCreate(&routeStation, models.RouteStation{
				RouteID:   route1.ID,
				StationID: station.ID,
				Direction: "up",
			})
		}
	}

	// 4. 创建票价规则
	var route1ForFare models.Route
	if err := db.Where("route_id = ?", "A").First(&route1ForFare).Error; err == nil {
		fare := models.Fare{
			RouteID:   route1ForFare.ID,
			BasePrice: 2.00,
			FareType:  "uniform",
			Status:    "active",
		}
		db.FirstOrCreate(&fare, models.Fare{RouteID: route1ForFare.ID, StartStation: 0, EndStation: 0})
	}

	var route2 models.Route
	if err := db.Where("route_id = ?", "B").First(&route2).Error; err == nil {
		fare := models.Fare{
			RouteID:      route2.ID,
			BasePrice:    2.00,
			FareType:     "segment",
			SegmentCount: 1,
			ExtraPrice:   1.00,
			Status:       "active",
		}
		db.FirstOrCreate(&fare, models.Fare{RouteID: route2.ID, StartStation: 0, EndStation: 0})
	}

	// 5. 创建换乘优惠规则
	var routeA, routeB models.Route
	var stationCenter models.Station
	if err := db.Where("route_id = ?", "A").First(&routeA).Error; err == nil {
		if err := db.Where("route_id = ?", "B").First(&routeB).Error; err == nil {
			if err := db.Where("station_id = ?", "ST001").First(&stationCenter).Error; err == nil {
				transfer := models.Transfer{
					FromRouteID:    routeA.ID,
					FromStationID:  stationCenter.ID,
					ToRouteID:      routeB.ID,
					ToStationID:    stationCenter.ID,
					DiscountAmount: 2.00,
					TimeWindow:     60,
					Status:         "active",
				}
				db.FirstOrCreate(&transfer, models.Transfer{
					FromRouteID:   routeA.ID,
					FromStationID: stationCenter.ID,
					ToRouteID:     routeB.ID,
					ToStationID:   stationCenter.ID,
				})
			}
		}
	}

	// 6. 创建折扣策略
	discountPolicies := []models.DiscountPolicy{
		{PolicyName: "月度累计折扣-8折", PolicyType: "monthly_accumulate", Threshold: 200.00, DiscountRate: 0.80, Status: "active"},
		{PolicyName: "月度累计折扣-5折", PolicyType: "monthly_accumulate", Threshold: 500.00, DiscountRate: 0.50, Status: "active"},
		{PolicyName: "学生卡折扣", PolicyType: "student", Threshold: 0.00, DiscountRate: 0.50, CardTypeFilter: "student", Status: "active"},
		{PolicyName: "老人卡免费", PolicyType: "elder", Threshold: 0.00, DiscountRate: 1.00, CardTypeFilter: "elder", Status: "active"},
	}

	for _, policy := range discountPolicies {
		db.FirstOrCreate(&policy, models.DiscountPolicy{
			PolicyType: policy.PolicyType,
			Threshold:  policy.Threshold,
		})
	}

	fmt.Println("数据库初始化数据已成功创建！")
}
