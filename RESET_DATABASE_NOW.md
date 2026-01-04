# 🔄 立即重置数据库

## 当前问题
数据库状态不一致，迁移失败。需要重置数据库。

## 快速重置（3步）

### 步骤1：运行重置脚本

在 PowerShell 中运行：

```powershell
.\scripts\reset_database.ps1
```

### 步骤2：如果脚本不工作，使用命令

```powershell
$env:PGPASSWORD="yl685306"; & "C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -d postgres -c "DROP DATABASE IF EXISTS bus_fare_system; CREATE DATABASE bus_fare_system;"
```

### 步骤3：重新运行程序

```bash
go run main.go
```

---

## 或者使用 pgAdmin（图形界面）

1. 打开 pgAdmin
2. 连接到 PostgreSQL 16
3. 展开 "Databases"
4. 右键 `bus_fare_system` → "Delete/Drop"
5. **勾选 "Cascade" 选项**（重要！）
6. 确认删除
7. 重新创建：右键 "Databases" → "Create" → "Database..." → 输入 `bus_fare_system` → "Save"

---

重置后，迁移应该能正常完成！
