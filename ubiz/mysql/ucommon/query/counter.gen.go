// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"git.umu.work/AI/uglib/ubiz/mysql/ucommon/model"
)

func newCounter(db *gorm.DB, opts ...gen.DOOption) counter {
	_counter := counter{}

	_counter.counterDo.UseDB(db, opts...)
	_counter.counterDo.UseModel(&model.Counter{})

	tableName := _counter.counterDo.TableName()
	_counter.ALL = field.NewAsterisk(tableName)
	_counter.ID = field.NewUint64(tableName, "id")
	_counter.MsgID = field.NewUint64(tableName, "msg_id")
	_counter.CounterKey = field.NewString(tableName, "counter_key")
	_counter.CounterNum = field.NewUint32(tableName, "counter_num")

	_counter.fillFieldMap()

	return _counter
}

type counter struct {
	counterDo counterDo

	ALL        field.Asterisk
	ID         field.Uint64 // id
	MsgID      field.Uint64 // 业务计数场景的消息(业务)id
	CounterKey field.String // 业务计数场景的唯一key
	CounterNum field.Uint32 // 计数器数量

	fieldMap map[string]field.Expr
}

func (c counter) Table(newTableName string) *counter {
	c.counterDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c counter) As(alias string) *counter {
	c.counterDo.DO = *(c.counterDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *counter) updateTableName(table string) *counter {
	c.ALL = field.NewAsterisk(table)
	c.ID = field.NewUint64(table, "id")
	c.MsgID = field.NewUint64(table, "msg_id")
	c.CounterKey = field.NewString(table, "counter_key")
	c.CounterNum = field.NewUint32(table, "counter_num")

	c.fillFieldMap()

	return c
}

func (c *counter) WithContext(ctx context.Context) *counterDo { return c.counterDo.WithContext(ctx) }

func (c counter) TableName() string { return c.counterDo.TableName() }

func (c counter) Alias() string { return c.counterDo.Alias() }

func (c *counter) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *counter) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 4)
	c.fieldMap["id"] = c.ID
	c.fieldMap["msg_id"] = c.MsgID
	c.fieldMap["counter_key"] = c.CounterKey
	c.fieldMap["counter_num"] = c.CounterNum
}

func (c counter) clone(db *gorm.DB) counter {
	c.counterDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c counter) replaceDB(db *gorm.DB) counter {
	c.counterDo.ReplaceDB(db)
	return c
}

type counterDo struct{ gen.DO }

//SelectForUpdateByCounterUniqueKey
//
//sql(select * from counter where counter_key=@key and msg_id=@msgID for update)
func (c counterDo) SelectForUpdateByCounterUniqueKey(key string, msgID uint64) (result model.Counter, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, key)
	params = append(params, msgID)
	generateSQL.WriteString("select * from counter where counter_key=? and msg_id=? for update ")

	var executeSQL *gorm.DB

	executeSQL = c.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result)
	err = executeSQL.Error
	return
}

func (c counterDo) Debug() *counterDo {
	return c.withDO(c.DO.Debug())
}

func (c counterDo) WithContext(ctx context.Context) *counterDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c counterDo) ReadDB() *counterDo {
	return c.Clauses(dbresolver.Read)
}

func (c counterDo) WriteDB() *counterDo {
	return c.Clauses(dbresolver.Write)
}

func (c counterDo) Session(config *gorm.Session) *counterDo {
	return c.withDO(c.DO.Session(config))
}

func (c counterDo) Clauses(conds ...clause.Expression) *counterDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c counterDo) Returning(value interface{}, columns ...string) *counterDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c counterDo) Not(conds ...gen.Condition) *counterDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c counterDo) Or(conds ...gen.Condition) *counterDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c counterDo) Select(conds ...field.Expr) *counterDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c counterDo) Where(conds ...gen.Condition) *counterDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c counterDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *counterDo {
	return c.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (c counterDo) Order(conds ...field.Expr) *counterDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c counterDo) Distinct(cols ...field.Expr) *counterDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c counterDo) Omit(cols ...field.Expr) *counterDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c counterDo) Join(table schema.Tabler, on ...field.Expr) *counterDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c counterDo) LeftJoin(table schema.Tabler, on ...field.Expr) *counterDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c counterDo) RightJoin(table schema.Tabler, on ...field.Expr) *counterDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c counterDo) Group(cols ...field.Expr) *counterDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c counterDo) Having(conds ...gen.Condition) *counterDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c counterDo) Limit(limit int) *counterDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c counterDo) Offset(offset int) *counterDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c counterDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *counterDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c counterDo) Unscoped() *counterDo {
	return c.withDO(c.DO.Unscoped())
}

func (c counterDo) Create(values ...*model.Counter) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c counterDo) CreateInBatches(values []*model.Counter, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c counterDo) Save(values ...*model.Counter) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c counterDo) First() (*model.Counter, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Counter), nil
	}
}

func (c counterDo) Take() (*model.Counter, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Counter), nil
	}
}

func (c counterDo) Last() (*model.Counter, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Counter), nil
	}
}

func (c counterDo) Find() ([]*model.Counter, error) {
	result, err := c.DO.Find()
	return result.([]*model.Counter), err
}

func (c counterDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Counter, err error) {
	buf := make([]*model.Counter, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c counterDo) FindInBatches(result *[]*model.Counter, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c counterDo) Attrs(attrs ...field.AssignExpr) *counterDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c counterDo) Assign(attrs ...field.AssignExpr) *counterDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c counterDo) Joins(fields ...field.RelationField) *counterDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c counterDo) Preload(fields ...field.RelationField) *counterDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c counterDo) FirstOrInit() (*model.Counter, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Counter), nil
	}
}

func (c counterDo) FirstOrCreate() (*model.Counter, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Counter), nil
	}
}

func (c counterDo) FindByPage(offset int, limit int) (result []*model.Counter, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c counterDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c counterDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c counterDo) Delete(models ...*model.Counter) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *counterDo) withDO(do gen.Dao) *counterDo {
	c.DO = *do.(*gen.DO)
	return c
}