# GORM
> GORM 支持的数据库类型: MySQL, PostgreSQL, SQlite, SQL Server.  
> 这里使用 MySQL.  
> [详细讲解](https://www.topgoer.com/数据库操作/gorm/)

## 1. 连接数据库
```golang
import (
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func main() {
  // 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name
  dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```

## 2. 使用连接池
GORM 使用 `database/sql` 维护连接池  
```golang
sqlDB, err := db.DB()

// SetMaxIdleConns 设置空闲连接池中连接的最大数量
sqlDB.SetMaxIdleConns(10)

// SetMaxOpenConns 设置打开数据库连接的最大数量. 
sqlDB.SetMaxOpenConns(100)

// SetConnMaxLifetime 设置了连接可复用的最大时间. 
sqlDB.SetConnMaxLifetime(time.Hour)
```

## 3. CRUD
定义数据库数据映射结构体:  
```golang
type User struct {
	Id          string
	Name        string
	StudentName string
	Password    string
	Email       string
	Salt        string
	RoleIds     string
	Enabled 	bool
}

func (User) TableName() string {
	return "t_user"
}
```
**注意**:  
- 结构体名称首字母必须大写.
- 结构体字段首字母必须大写, 字段按照驼峰命名法, 对应数据库中的 `_`.
- 如果结构体名称与数据库表名称不同, 则需要编写 `TableName()` 方法, 返回对应的数据库表的名称.

### 3.1 查询

使用 `gorm.DB.Find()` 来进行全量查询.  
使用 `gorm.DB.Where()` 来指定查询条件.

查询 API:  
```golang
// 获取第一条记录,按主键排序
db.First(&user)
//// SELECT * FROM users ORDER BY id LIMIT 1;

// 获取一条记录,不指定排序
db.Take(&user)
//// SELECT * FROM users LIMIT 1;

// 获取最后一条记录,按主键排序
db.Last(&user)
//// SELECT * FROM users ORDER BY id DESC LIMIT 1;

// 获取所有的记录
db.Find(&users)
//// SELECT * FROM users;

// 通过主键进行查询 (仅适用于主键是数字类型)
db.First(&user, 10)
//// SELECT * FROM users WHERE id = 10;
```

使用 `Where()`:  
```golang
// 获取第一条匹配的记录
db.Where("name = ?", "riicarus").First(&user)
//// SELECT * FROM users WHERE name = 'riicarus' limit 1;

// 获取所有匹配的记录
db.Where("name = ?", "riicarus").Find(&users)
//// SELECT * FROM users WHERE name = 'riicarus';

// <>
db.Where("name <> ?", "riicarus").Find(&users)

// IN
db.Where("name in (?)", []string{"riicarus", "riicarus 2"}).Find(&users)

// LIKE
db.Where("name LIKE ?", "%jin%").Find(&users)

// AND
db.Where("name = ? AND age >= ?", "riicarus", "22").Find(&users)

// Time
db.Where("updated_at > ?", lastWeek).Find(&users)

// BETWEEN
db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)
```

使用 struct/map:  
将对应的 key-value 加入 WHERE 条件中.  
```golang
// Struct
db.Where(&User{Name: "riicarus", Age: 20}).First(&user)
//// SELECT * FROM users WHERE name = "riicarus" AND age = 20 LIMIT 1;

// Map
db.Where(map[string]interface{}{"name": "riicarus", "age": 20}).Find(&users)
//// SELECT * FROM users WHERE name = "riicarus" AND age = 20;

// 多主键 slice 查询
db.Where([]int64{20, 21, 22}).Find(&users)
//// SELECT * FROM users WHERE id IN (20, 21, 22);
```

NOT/OR 用法相似.  

### 3.2 创建
使用 `gorm.DB.Create()` 来进行创建. 可以使用 `gorm.DB.NewRecord()` 来判断是否记录存在.  
也可以通过设置标签来指定字段的默认值: `gorm:"default:'riicarus'"`.  
在创建过程中, gorm 会排除所有值为 `nil`, `false` 或 `0` 的字段, 但是如果有默认值, 会使用其默认值.  
如果想要使零值生效, 可以使用指针类型或者其他类型的值, 或者设置默认值为零值.  

### 3.3 更新
`gorm.DB.Save()` 方法会更新所有的字段, 无论这些字段有没有被修改过.  
如果只想更新修改过的字段, 那么需要使用 `gorm.DB.Update()` 或 `gorm.DB.Updates()` 方法.  
```golang
// 如果单个属性被更改了,更新它
db.Model(&user).Update("name", "hello")
//// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;

// 使用组合条件更新单个属性
db.Model(&user).Where("active = ?", true).Update("name", "hello")
//// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;

// 使用 `map` 更新多个属性,只会更新那些被更改了的字段
db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
//// UPDATE users SET name='hello', age=18, actived=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

// 使用 `struct` 更新多个属性,只会更新那些被修改了的和非空的字段
db.Model(&user).Updates(User{Name: "hello", Age: 18})
//// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;

// 警告: 当使用结构体更新的时候, GORM 只会更新那些非空的字段
// 例如下面的更新,没有东西会被更新,因为像 "", 0, false 是这些字段类型的空值
db.Model(&user).Updates(User{Name: "", Age: 0, Actived: false})
```
可以使用表达式:  
```golang
DB.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))
//// UPDATE "products" SET "price" = price * '2' + '100', "updated_at" = '2013-11-17 21:34:10' WHERE "id" = '2';

DB.Model(&product).Updates(map[string]interface{}{"price": gorm.Expr("price * ? + ?", 2, 100)})
//// UPDATE "products" SET "price" = price * '2' + '100', "updated_at" = '2013-11-17 21:34:10' WHERE "id" = '2';

DB.Model(&product).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
//// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = '2';

DB.Model(&product).Where("quantity > 1").UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
//// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = '2' AND quantity > 1;
```

### 3.4 删除
使用 `gorm.DB.Delete()` 方法使用主键对对应元素进行删除.  
**注意**: 当删除一条记录的时候, 你需要确定这条记录的主键有值, GORM会使用主键来删除这条记录. 如果主键字段为空,GORM会删除模型中所有的记录.
