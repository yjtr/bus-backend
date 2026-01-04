# PowerShell 脚本：重置数据库
# 使用方法：在PowerShell中运行: .\scripts\reset_database.ps1

# PostgreSQL路径
$pgPath = "C:\Program Files\PostgreSQL\16\bin\psql.exe"
if (-not (Test-Path $pgPath)) {
    $pgPath = "C:\Program Files\PostgreSQL\15\bin\psql.exe"
}

if (-not (Test-Path $pgPath)) {
    Write-Host "错误：找不到PostgreSQL安装路径" -ForegroundColor Red
    exit 1
}

Write-Host "正在重置数据库..." -ForegroundColor Green

# 使用环境变量设置密码
$env:PGPASSWORD = "yl685306"

# 断开所有连接
Write-Host "步骤1: 断开所有连接到bus_fare_system的会话..." -ForegroundColor Yellow
& $pgPath -U postgres -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'bus_fare_system' AND pid <> pg_backend_pid();" 2>$null

# 删除数据库
Write-Host "步骤2: 删除数据库 bus_fare_system..." -ForegroundColor Yellow
& $pgPath -U postgres -d postgres -c "DROP DATABASE IF EXISTS bus_fare_system;"

if ($LASTEXITCODE -ne 0) {
    Write-Host "删除数据库时出现警告（可能数据库不存在）" -ForegroundColor Yellow
}

# 创建新数据库
Write-Host "步骤3: 创建新数据库 bus_fare_system..." -ForegroundColor Yellow
& $pgPath -U postgres -d postgres -c "CREATE DATABASE bus_fare_system;"

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ 数据库重置成功！" -ForegroundColor Green
    Write-Host "现在可以重新运行程序: go run main.go" -ForegroundColor Cyan
} else {
    Write-Host "❌ 创建数据库失败" -ForegroundColor Red
}

# 清除密码环境变量
Remove-Item Env:\PGPASSWORD
