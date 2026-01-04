-- 创建数据库脚本
-- 使用方法: psql -U postgres -f scripts/create_database.sql

-- 创建数据库（如果不存在）
SELECT 'CREATE DATABASE bus_fare_system'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'bus_fare_system')\gexec

-- 连接到新数据库并设置权限
\c bus_fare_system

-- 给postgres用户所有权限（如果需要）
GRANT ALL PRIVILEGES ON DATABASE bus_fare_system TO postgres;
