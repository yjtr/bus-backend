-- 修改用户密码
-- 使用方法: psql -U postgres -f scripts/change_user_password.sql

-- 修改bus_admin用户密码
ALTER USER bus_admin WITH PASSWORD 'your_new_password_here';

-- 查看用户信息
SELECT usename, usesuper, usecreatedb 
FROM pg_user 
WHERE usename = 'bus_admin';
