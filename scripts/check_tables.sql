-- 检查数据库中已创建的表
-- 使用方法: psql -U postgres -d bus_fare_system -f scripts/check_tables.sql

SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
ORDER BY table_name;
