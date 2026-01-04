-- 重置数据库脚本
-- 如果数据库迁移出现问题，可以先删除数据库重新创建

-- 注意：这个脚本需要在连接到postgres数据库时执行，而不是bus_fare_system数据库

-- 1. 断开所有连接到bus_fare_system的会话
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = 'bus_fare_system' AND pid <> pg_backend_pid();

-- 2. 删除数据库
DROP DATABASE IF EXISTS bus_fare_system;

-- 3. 重新创建数据库
CREATE DATABASE bus_fare_system;

-- 执行完成后，重新运行程序即可
