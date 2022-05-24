# Gin+Gorm开发在线备忘录（To do list）

Gin：

- [gin框架](https://www.topgoer.com/gin%E6%A1%86%E6%9E%B6/)
- [github gin README](https://github.com/gin-gonic/gin)
- [Go Gin 简明教程](https://geektutu.com/post/quick-go-gin.html)

Gorm：

- [GORM 指南](https://gorm.io/zh_CN/docs/)
- **[GORM入门指南](https://www.topgoer.com/%E6%95%B0%E6%8D%AE%E5%BA%93%E6%93%8D%E4%BD%9C/gorm/%E5%85%A5%E9%97%A8%E6%8C%87%E5%8D%97/)**

## 记录ShouldBind使用遇到的坑（GET or POST？）

这里是我自己马虎了：

![](https://img-qingbo.oss-cn-beijing.aliyuncs.com/img/20220523213051.png)

因为API里使用的是`GET`，所以在PostMan中就应该这么写URL：`localhost:8080/api/task/getall?page_size=2&page_num=1`

![](https://img-qingbo.oss-cn-beijing.aliyuncs.com/img/20220523213234.png)

## gorm分页查询报错 sql: no rows in result set

参考文章：[gorm分页查询count报错 sql: no rows in result set](https://www.jianshu.com/p/fa267de4f5d0)

原本gorm的查询是这么写的：`DB.Where("uid = ?", uid).Limit(s.PageSize).Offset(s.PageSize * (s.PageNum - 1)).Find(&tasks).Count(&count)`。

![](https://img-qingbo.oss-cn-beijing.aliyuncs.com/img/20220524145038.png)



==注意==：`Count()`查询必须在`where`条件之后，`limit`,`offset` 分页之前

- 如果写在`limit`,`offset` 分页 之前之后，在第二页开始就会报错 `sql: no rows in result set`
- `count`也不能太前，否则查询出来的总数将是所有数据总数，非条件过滤后的条数

> ```go
> // 错误db.Order("id desc").Limit(10).Offset(10).Find(&List).Count(&totalRows)
> //count应该条件where之后，分页条件之前，如以下，但是结果报错:incorrect table name
> db.Order("id desc").Count(&totalRows).Limit(10).Offset(10).Find(&List)
> 
> //正确写法
> db.Table("TableName").Order("id desc").Count(&totalRows).Limit(10).Offset(10).Find(&List)
> ```

修改后应该如下：`DB.Table("task").Order("created_at desc").Where("uid = ?", uid).Count(&count).Limit(s.PageSize).Offset(s.PageSize * (s.PageNum - 1)).Find(&tasks)`

![](https://img-qingbo.oss-cn-beijing.aliyuncs.com/img/20220524144936.png)



## 关于分页查询中添加keyword进行模糊查询的gorm写法

代码中我是这样写的：

```go
err := mysql.DB.Table("task").Where("uid = ? and tittle like ?", uid, "%"+s.KeyWord+"%").Or("uid = ? and context like ?", uid, "%"+s.KeyWord+"%").
	Count(&count).Limit(s.PageSize).Offset(s.PageSize * (s.PageNum - 1)).Find(&tasks).Error
```

但也可以这样写，连续带上两个`Where`：

```go
err := mysql.DB.Table("task").Where("uid = ?", uid).Where("tittle like ? or context like ?", "%"+s.KeyWord+"%", "%"+s.KeyWord+"%").
   Count(&count).Limit(s.PageSize).Offset(s.PageSize * (s.PageNum - 1)).Find(&tasks).Error
```













