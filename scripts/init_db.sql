-- 公交刷卡收费系统 - 数据库初始化脚本
-- 此脚本用于创建示例数据，供开发和测试使用

-- 注意：表结构会通过GORM AutoMigrate自动创建，此脚本主要用于初始化数据

-- 1. 插入示例线路数据
INSERT INTO routes (route_id, name, area, status, direction, created_at, updated_at) VALUES
('A', '1路', '市区', 'active', 'up', NOW(), NOW()),
('B', '2路', '市区', 'active', 'up', NOW(), NOW()),
('K1', '快速1路', '市区', 'active', 'up', NOW(), NOW())
ON CONFLICT (route_id) DO NOTHING;

-- 2. 插入示例站点数据
INSERT INTO stations (station_id, name, latitude, longitude, address, is_transfer, created_at, updated_at) VALUES
('ST001', '市中心站', 39.9042, 116.4074, '市中心广场', true, NOW(), NOW()),
('ST002', '火车站', 39.9019, 116.4250, '火车站广场', true, NOW(), NOW()),
('ST003', '大学城站', 39.9200, 116.4000, '大学城入口', false, NOW(), NOW()),
('ST004', '商业区站', 39.9100, 116.4200, '商业区中心', false, NOW(), NOW()),
('ST005', '体育场站', 39.9150, 116.4100, '体育场门口', false, NOW(), NOW()),
('ST006', '医院站', 39.9050, 116.4150, '市人民医院', false, NOW(), NOW())
ON CONFLICT (station_id) DO NOTHING;

-- 3. 插入线路-站点关联（需要先获取route和station的ID，这里使用示例）
-- 注意：实际使用时需要通过子查询获取ID
-- 示例：1路线路（假设route ID为1）
-- INSERT INTO route_stations (route_id, station_id, sequence, direction, created_at, updated_at)
-- SELECT r.id, s.id, row_number() OVER (ORDER BY s.station_id), 'up', NOW(), NOW()
-- FROM routes r, stations s
-- WHERE r.route_id = 'A' AND s.station_id IN ('ST001', 'ST002', 'ST003', 'ST004')
-- ON CONFLICT DO NOTHING;

-- 4. 插入票价规则
-- 统一票价示例：1路线路统一2元
INSERT INTO fares (route_id, start_station, end_station, base_price, fare_type, status, created_at, updated_at)
SELECT id, 0, 0, 2.00, 'uniform', 'active', NOW(), NOW()
FROM routes WHERE route_id = 'A'
ON CONFLICT DO NOTHING;

-- 分段计价示例：2路线路起步2元，每段加1元
INSERT INTO fares (route_id, start_station, end_station, base_price, fare_type, segment_count, extra_price, status, created_at, updated_at)
SELECT id, 0, 0, 2.00, 'segment', 1, 1.00, 'active', NOW(), NOW()
FROM routes WHERE route_id = 'B'
ON CONFLICT DO NOTHING;

-- 5. 插入换乘优惠规则
-- 示例：从1路的市中心站换乘到2路，60分钟内优惠2元
INSERT INTO transfers (from_route_id, from_station_id, to_route_id, to_station_id, discount_amount, time_window, status, created_at, updated_at)
SELECT 
    r1.id, s1.id, r2.id, s2.id, 2.00, 60, 'active', NOW(), NOW()
FROM routes r1, routes r2, stations s1, stations s2
WHERE r1.route_id = 'A' AND r2.route_id = 'B' 
  AND s1.station_id = 'ST001' AND s2.station_id = 'ST001'
ON CONFLICT DO NOTHING;

-- 6. 插入折扣策略
-- 月度累计折扣：满200元享8折，满500元享5折
INSERT INTO discount_policies (policy_name, policy_type, threshold, discount_rate, status, created_at, updated_at) VALUES
('月度累计折扣-8折', 'monthly_accumulate', 200.00, 0.80, 'active', NOW(), NOW()),
('月度累计折扣-5折', 'monthly_accumulate', 500.00, 0.50, 'active', NOW(), NOW()),
('学生卡折扣', 'student', 0.00, 0.50, 'active', NOW(), NOW()),
('老人卡折扣', 'elder', 0.00, 0.00, 'active', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 7. 插入示例卡片（可选，用于测试）
-- INSERT INTO cards (card_id, holder_name, card_type, status, balance, created_at, updated_at) VALUES
-- ('12345678', '测试用户1', 'normal', 'active', 0.00, NOW(), NOW()),
-- ('87654321', '测试用户2', 'student', 'active', 0.00, NOW(), NOW()),
-- ('11111111', '测试用户3', 'elder', 'active', 0.00, NOW(), NOW())
-- ON CONFLICT (card_id) DO NOTHING;
