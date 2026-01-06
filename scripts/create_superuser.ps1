# 创建PostgreSQL超级用户脚本
# 使用方法: .\scripts\create_superuser.ps1

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
    Write-Host "Please confirm PostgreSQL is installed, or specify path manually" -ForegroundColor Yellow
    exit 1
}

Write-Host "Using PostgreSQL: $pgPath" -ForegroundColor Cyan

# 设置密码环境变量（使用postgres用户的密码）
$env:PGPASSWORD = "yl685306"

# 用户名和密码
$newUser = "bus_admin"
$newPassword = "bus_admin_2024"

Write-Host ""
Write-Host "Creating user $newUser ..." -ForegroundColor Yellow

# 创建用户并授予权限（使用SQL文件）
$sqlFile = Join-Path $PSScriptRoot "create_superuser_temp.sql"

# 创建临时SQL文件
$sqlContent = @"
DO `$`$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = '$newUser') THEN
        CREATE USER $newUser WITH PASSWORD '$newPassword';
        RAISE NOTICE 'User $newUser created';
    ELSE
        RAISE NOTICE 'User $newUser already exists, updating privileges';
    END IF;
END
`$`$;

ALTER USER $newUser WITH SUPERUSER CREATEDB;

SELECT usename, usesuper, usecreatedb 
FROM pg_user 
WHERE usename = '$newUser';
"@

# 使用UTF8NoBOM编码避免BOM问题
$utf8NoBom = New-Object System.Text.UTF8Encoding $false
[System.IO.File]::WriteAllText($sqlFile, $sqlContent, $utf8NoBom)

# 执行SQL文件
& $pgPath -U postgres -d postgres -f $sqlFile

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "User created successfully!" -ForegroundColor Green
    Write-Host "Username: $newUser" -ForegroundColor Cyan
    Write-Host "Password: $newPassword" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Please keep the password safe, consider using a stronger password in production" -ForegroundColor Yellow
} else {
    Write-Host ""
    Write-Host "Failed to create user" -ForegroundColor Red
}

# 清理临时文件
if (Test-Path $sqlFile) {
    Remove-Item $sqlFile -Force
}

# 清除密码环境变量
Remove-Item Env:\PGPASSWORD -ErrorAction SilentlyContinue
