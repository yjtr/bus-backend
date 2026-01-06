-- 创建具有超级用户权限的新用户
-- 使用方法: psql -U postgres -f scripts/create_superuser.sql

-- 创建用户（如果不存在）
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'bus_admin') THEN
        CREATE USER bus_admin WITH PASSWORD 'bus_admin_2024';
        RAISE NOTICE '用户 bus_admin 已创建';
    ELSE
        RAISE NOTICE '用户 bus_admin 已存在';
    END IF;
END
$$;

-- 授予超级用户权限
ALTER USER bus_admin WITH SUPERUSER;

-- 授予创建数据库权限
ALTER USER bus_admin CREATEDB;

-- 授予所有数据库的连接权限（可选）
-- GRANT ALL PRIVILEGES ON DATABASE bus_fare_system TO bus_admin;

-- 显示用户信息
SELECT usename, usesuper, usecreatedb 
FROM pg_user 
WHERE usename = 'bus_admin';
