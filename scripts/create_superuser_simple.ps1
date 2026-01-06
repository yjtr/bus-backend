# 创建PostgreSQL超级用户脚本（简化版）
# 使用方法: .\scripts\create_superuser_simple.ps1

$ErrorActionPreference = "Stop"

# PostgreSQL路径（支持版本16和15）
$pgPaths = @(
    "C:\Program Files\PostgreSQL\16\bin\psql.exe",
    "C:\Program Files\PostgreSQL\15\bin\psql.exe"
)

$pgPath = $null
foreach ($path in $pgPaths) {
    if (Test-Path $path) {
        $pgPath = $path
        break
    }
}

if ($null -eq $pgPath) {
    Write-Host "Error: Cannot find PostgreSQL psql.exe" -ForegroundColor Red
    exit 1
}

Write-Host "Using PostgreSQL: $pgPath" -ForegroundColor Cyan

# 设置密码环境变量
$env:PGPASSWORD = "yl685306"

# 用户名和密码
$newUser = "bus_admin"
$newPassword = "bus_admin_2024"

Write-Host ""
Write-Host "Creating user $newUser ..." -ForegroundColor Yellow

# 直接执行SQL命令（不使用DO块，更简单）
$sql = "CREATE USER $newUser WITH PASSWORD '$newPassword' SUPERUSER CREATEDB;"

& $pgPath -U postgres -d postgres -c $sql 2>&1 | Out-Null

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "User created successfully!" -ForegroundColor Green
    Write-Host "Username: $newUser" -ForegroundColor Cyan
    Write-Host "Password: $newPassword" -ForegroundColor Cyan
} else {
    # 如果用户已存在，尝试只更新权限
    Write-Host "User may already exist, updating privileges..." -ForegroundColor Yellow
    $sqlUpdate = "ALTER USER $newUser WITH SUPERUSER CREATEDB;"
    & $pgPath -U postgres -d postgres -c $sqlUpdate 2>&1 | Out-Null
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "User privileges updated!" -ForegroundColor Green
        Write-Host "Username: $newUser" -ForegroundColor Cyan
        Write-Host "Password: $newPassword" -ForegroundColor Cyan
    } else {
        Write-Host ""
        Write-Host "Failed to create/update user" -ForegroundColor Red
    }
}

# 清除密码环境变量
Remove-Item Env:\PGPASSWORD -ErrorAction SilentlyContinue
