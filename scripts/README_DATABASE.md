# 数据库创建指南

## 方式1：使用 psql 命令行（推荐）

1. 打开命令提示符（CMD）或 PowerShell

2. 连接到 PostgreSQL（根据你的安装路径调整）：
   ```bash
   psql -U postgres
   ```
   或者如果 psql 不在 PATH 中：
   ```bash
   "C:\Program Files\PostgreSQL\<版本>\bin\psql.exe" -U postgres
   ```

3. 输入 PostgreSQL 的密码（你配置文件中使用的密码）

4. 创建数据库：
   ```sql
   CREATE DATABASE bus_fare_system;
   ```

5. 验证数据库是否创建成功：
   ```sql
   \l
   ```
   你应该能看到 `bus_fare_system` 在列表中

6. 退出 psql：
   ```sql
   \q
   ```

## 方式2：使用 pgAdmin（图形界面）

1. 打开 pgAdmin

2. 连接到你的 PostgreSQL 服务器

3. 右键点击 "Databases" -> "Create" -> "Database..."

4. 在 "Database" 字段中输入: `bus_fare_system`

5. 点击 "Save"

## 方式3：使用 SQL 脚本

如果你已经连接到 postgres 数据库，可以执行：

```sql
CREATE DATABASE bus_fare_system;
```

## 验证

创建完成后，你可以运行程序，它应该能够连接到数据库并自动创建表结构。

如果仍有问题，请检查：
- PostgreSQL 服务是否正在运行
- 配置文件 `config/config.yaml` 中的密码是否正确
- 数据库名称是否正确（大小写敏感）
