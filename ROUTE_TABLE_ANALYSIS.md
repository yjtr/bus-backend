# Routes 表问题分析和简化方案

## 问题分析

### 为什么 routes 表不存在？

1. **数据库状态不一致**
   - 之前的迁移失败导致数据库状态异常
   - 需要重置数据库

2. **模型关联复杂度**
   - Route 模型中的 `Stations []RouteStation` 关联可能在某些情况下影响迁移
   - 虽然这是标准 GORM 写法，但在迁移顺序错误时可能导致问题

## 简化方案

### 简化后的 Route 模型（只保留扣费必需字段）

```go
type Route struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

    RouteID   string `gorm:"uniqueIndex;not null;size:50" json:"route_id"` // 线路编号
    Name      string `gorm:"size:100;not null" json:"name"`                // 线路名称
    Status    string `gorm:"size:20;default:'active'" json:"status"`       // 状态：active, inactive
}
```

### 移除的字段

1. ❌ `Area` - 区域字段（扣费不需要）
2. ❌ `Direction` - 方向字段（扣费不需要）
3. ❌ `Stations []RouteStation` - 关联定义（表结构不需要，只用于查询时 Preload）

### 扣费功能所需的最小字段

从 `services/fare_service.go` 分析，扣费只需要：
- `routeID` (uint) - 用于查询票价规则（fares 表）
- 线路基本信息（RouteID, Name, Status）

其他关联数据通过其他表查询：
- 票价规则 → `fares` 表（通过 route_id 查询）
- 站点信息 → `route_stations` 表（通过 route_id 查询）
- 换乘规则 → `transfers` 表（通过 from_route_id/to_route_id 查询）

## 解决方案

1. ✅ **简化 Route 模型** - 移除不必要的字段和关联
2. ✅ **修复 Controller** - 移除对 Stations 关联的 Preload
3. ✅ **重置数据库** - 清空状态，重新迁移

## 下一步

1. 重置数据库：
```powershell
.\scripts\reset_database.ps1
```

2. 运行程序：
```bash
go run main.go
```

简化后的 routes 表应该能正常创建。
