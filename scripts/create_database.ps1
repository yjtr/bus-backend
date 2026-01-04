# PowerShell 脚本：创建数据库
# 使用方法：在PowerShell中运行: .\scripts\create_database.ps1

# PostgreSQL路径（根据你的安装版本调整，可能是15或16）
$pgPath = "C:\Program Files\PostgreSQL\16\bin\psql.exe"
if (-not (Test-Path $pgPath)) {
    $pgPath = "C:\Program Files\PostgreSQL\15\bin\psql.exe"
}

if (-not (Test-Path $pgPath)) {
    Write-Host "错误：找不到PostgreSQL安装路径" -ForegroundColor Red
    Write-Host "请手动修改脚本中的路径，或使用pgAdmin创建数据库" -ForegroundColor Yellow
    exit 1
}

Write-Host "正在连接到PostgreSQL..." -ForegroundColor Green
Write-Host "请输入密码: yl685306" -ForegroundColor Yellow

# 使用环境变量设置密码（避免交互式输入）
$env:PGPASSWORD = "yl685306"

# 执行SQL命令创建数据库
$sqlCommand = "CREATE DATABASE bus_fare_system;"
& $pgPath -U postgres -d postgres -c $sqlCommand

if ($LASTEXITCODE -eq 0) {
    Write-Host "数据库 'bus_fare_system' 创建成功！" -ForegroundColor Green
} else {
    Write-Host "创建数据库失败，错误代码: $LASTEXITCODE" -ForegroundColor Red
    Write-Host "可能的原因：" -ForegroundColor Yellow
    Write-Host "1. 数据库已存在" -ForegroundColor Yellow
    Write-Host "2. 密码不正确" -ForegroundColor Yellow
    Write-Host "3. PostgreSQL服务未启动" -ForegroundColor Yellow
}

# 清除密码环境变量
Remove-Item Env:\PGPASSWORD
