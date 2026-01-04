# 手动创建数据库指南

## PowerShell 命令（推荐）

在 PowerShell 中运行以下命令（注意使用 `&` 操作符）：

```powershell
# 方法1：使用 & 操作符（推荐）
& "C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -c "CREATE DATABASE bus_fare_system;"
```

如果版本是15，使用：
```powershell
& "C:\Program Files\PostgreSQL\15\bin\psql.exe" -U postgres -c "CREATE DATABASE bus_fare_system;"
```

或者使用环境变量设置密码（避免交互式输入）：
```powershell
$env:PGPASSWORD = "yl685306"
& "C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -d postgres -c "CREATE DATABASE bus_fare_system;"
Remove-Item Env:\PGPASSWORD
```

## 使用脚本（最简单）

直接运行我创建的PowerShell脚本：
```powershell
.\scripts\create_database.ps1
```

## 使用 pgAdmin（图形界面）

1. 打开 pgAdmin（通常在开始菜单中）
2. 输入主密码（安装PostgreSQL时设置的）
3. 展开 "Servers" -> "PostgreSQL 15"（或你的版本）
4. 右键点击 "Databases" -> "Create" -> "Database..."
5. 在 "Database" 字段输入：`bus_fare_system`
6. 点击 "Save"

## 使用 CMD（命令提示符）

如果你使用 CMD 而不是 PowerShell：

```cmd
"C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -c "CREATE DATABASE bus_fare_system;"
```
