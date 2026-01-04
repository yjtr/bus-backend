# 快速解决数据库问题

## 当前错误
数据库 `bus_fare_system` 不存在

## 最简单的方法（3步解决）

### 方法1：使用 pgAdmin（推荐，最可靠）

1. **打开 pgAdmin**
   - 在Windows开始菜单搜索 ""
   - 或找到PostgreSQL程序组，点击 pgAdmin

2. **连接服务器**
   - 首次打开需要设置主密码（用于保护保存的密码，可以设置一个简单的）
   - 左侧展开 "Servers" → "PostgreSQL 16"（或你的版本）
   - 如果提示输入密码，输入：`yl685306`

3. **创建数据库**
   - 在左侧树形菜单中，找到 "Databases"
   - **右键点击 "Databases"**
   - 选择 **"Create"** → **"Database..."**
   - 在弹出窗口的 "Database" 字段输入：`bus_fare_system`
   - 其他选项保持默认，直接点击 **"Save"** 按钮（窗口底部）

4. **验证**
   - 左侧 "Databases" 下应该能看到 `bus_fare_system`
   - 然后重新运行程序即可

---

### 方法2：使用 PowerShell 命令（如果pgAdmin不可用）

在 PowerShell 中运行（**一行命令**）：

```powershell
$env:PGPASSWORD="yl685306"; & "C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -d postgres -c "CREATE DATABASE bus_fare_system;"
```

如果提示找不到路径，先确认PostgreSQL安装路径，然后修改命令中的路径。

---

### 方法3：使用测试脚本检查

运行测试脚本查看当前状态：

```bash
go run scripts/test_connection.go
```

这会显示：
- 是否能连接PostgreSQL
- 数据库是否存在
- 所有数据库列表

---

## 创建成功后的验证

数据库创建成功后，重新运行程序：

```bash
go run main.go
```

应该看到：
```
数据库连接成功
服务器启动在 0.0.0.0:8080
```

如果还有问题，请告诉我具体的错误信息。

---

## 如果遇到迁移错误（如表已存在、外键错误等）

如果遇到数据库迁移错误，需要重置数据库：

### 快速重置方法（PowerShell）

```powershell
.\scripts\reset_database.ps1
```

或者使用一行命令：

```powershell
$env:PGPASSWORD="yl685306"; & "C:\Program Files\PostgreSQL\16\bin\psql.exe" -U postgres -d postgres -c "DROP DATABASE IF EXISTS bus_fare_system; CREATE DATABASE bus_fare_system;"
```

这会删除并重新创建数据库，然后重新运行程序即可。
pgAdmin 4