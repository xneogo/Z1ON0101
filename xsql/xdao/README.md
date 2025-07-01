## xdao

xdao 是 /seele/xsalIface/repo.go 中 DaoModel 的标准化实现方案

### 快速开始

#### 支持以下接口功能

```go

// DaoModel
// interactive between DataObj and db
// reflect db instance rows to DataObj
type DaoModel[EntityObj any, DObj DataObj[EntityObj]] interface {
	Init(cons SqlConstructor, tableName func() string, omits func() []string, b Bind)
	TableName() string
	Omits() []string
	GetScanner() Scanner
	GetBuilder() Builder

	SelectOne(ctx context.Context, db XDB, where map[string]interface{}) (DObj, error)
	SelectMulti(ctx context.Context, db XDB, where map[string]interface{}) ([]DObj, error)
	Insert(ctx context.Context, db XDB, data []map[string]interface{}) (int64, error)
	Update(ctx context.Context, db XDB, where, data map[string]interface{}) (int64, error)
	Delete(ctx context.Context, db XDB, where map[string]interface{}) (int64, error)
	CountOf(ctx context.Context, db XDB, where map[string]interface{}) (count int, err error)
	ToEntity(ctx context.Context, t DObj) *EntityObj
	MultiToEntity(ctx context.Context, ts []DObj) []*EntityObj
}
```

其中 EntityObj 为项目 service 层级定义的数据结构；DataObj 为 repo 层级定义的数据结构。

按照泛型定义的写法，强制关联当前 DaoModel 接口实例与目标数据结构的关系。

#### 复杂sql

```go

// ComplexQuery
// you can use this default logic or
// you can build your own query logic with or without tableName or columns
// depends on your ToSql func
func ComplexQuery[ans any](tableName string, columns ...string) xsqlIface.ComplexQueryMod[ans] {
	return func(
		ctx context.Context,
		db xsqlIface.XDB,
		scanner xsqlIface.Scanner,
		f xsqlIface.ToSql,
		bind xsqlIface.Bind,
	) (res []ans, err error) {
		if nil == db {
			return nil, errors.New("manager.XDB object couldn't be nil")
		}
		cond, vals, err := f(tableName, columns...)
		if nil != err {
			return nil, err
		}
		xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
		xlog.Infof(ctx, "build cond: %s vals: %v", cond, vals)
		row, err := db.QueryContext(ctx, cond, vals...)
		if nil != err || nil == row {
			return nil, err
		}
		defer row.Close()
		err = scanner.Scan(row, &res, bind)
		return res, err
	}
}

func ComplexExec(tableName string) xsqlIface.ComplexExecMod {
	return func(
		ctx context.Context,
		db xsqlIface.XDB,
		f xsqlIface.ToSql,
	) (int64, error) {
		if nil == db {
			return 0, errors.New("manager.XDB object couldn't be nil")
		}
		cond, vals, err := f(tableName)
		if nil != err {
			return 0, err
		}
		xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
		result, err := db.ExecContext(ctx, cond, vals...)
		if nil != err {
			return 0, err
		}
		return result.RowsAffected()
	}
}
```

默认支持 Query 和 Exec 两种操作，查询默认返回列表结构体，具体取用根据实际逻辑

常规理解下，复杂sql操作涉及的数据结构与库表关联的数据结构无关，因此这两个定义中需要单独定义取值结果对应的数据结构

ToSql 定义为
```go
type ToSql func(tName string, columns ...string) (string, []interface{}, error)
```

在执行数据库相关指令流程中，基本上都会生成 sql， args， error 部分来确定最终的执行。



