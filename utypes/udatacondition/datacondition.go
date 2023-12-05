package udatacondition

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/uerrors"
	"git.umu.work/be/goframework/logger"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"time"
)

type Condition interface {
	Query(ctx context.Context, table string, dao gen.DO) (gen.DO, error)
}

type TimeRange struct {
	field string
	start *time.Time
	end   *time.Time
}

func NewTimeRange(field string, start, end *time.Time) *TimeRange {
	return &TimeRange{
		field: field,
		start: start,
		end:   end,
	}
}

func (t *TimeRange) Query(ctx context.Context, table string, dao gen.DO) (gen.DO, error) {
	timeRangeField := field.NewTime(table, t.field)
	if t.start != nil {
		dao = *dao.Where(timeRangeField.Gte(*t.start)).(*gen.DO)
	}
	if t.end != nil {
		dao = *dao.Where(timeRangeField.Lte(*t.start)).(*gen.DO)
	}

	return dao, nil
}

type Pagination interface {
	GetTotal(ctx context.Context) int64
	Query(ctx context.Context, table string, dao gen.DO) (gen.DO, error)
}

type Order struct {
	field  string
	isDesc bool
}

func NewOrder(field string, isDesc bool) *Order {
	return &Order{
		field:  field,
		isDesc: isDesc,
	}
}

func (o *Order) Query(ctx context.Context, table string, dao gen.DO) (gen.DO, error) {
	orderField := field.NewField(table, o.field)
	if o.isDesc {
		dao = *dao.Order(orderField.Desc()).(*gen.DO)
	} else {
		dao = *dao.Order(orderField).(*gen.DO)
	}

	return dao, nil
}

func (o *Order) IsDesc(ctx context.Context) bool {
	return o.isDesc
}

type PagePagination struct {
	pageNum  uint32
	pageSize uint32
	total    int64
	order    *Order
}

func NewPagePagination(pageSize, pageNum uint32, order *Order) Pagination {
	return &PagePagination{
		pageNum:  pageNum,
		pageSize: pageSize,
		order:    order,
	}
}

func (p *PagePagination) Query(ctx context.Context, table string, dao gen.DO) (gen.DO, error) {
	total, err := dao.Count()
	if err != nil {
		return gen.DO{}, err
	}
	p.total = total
	dao = *dao.Offset(p.getOffset()).Limit(p.getLimit()).(*gen.DO)
	if p.order != nil {
		dao, err = p.order.Query(ctx, table, dao)
		if err != nil {
			return gen.DO{}, err
		}
	}

	return dao, nil
}

func (p *PagePagination) getOffset() int {
	return (p.getPage() - 1) * p.getLimit()
}

func (p *PagePagination) getLimit() int {
	if p.pageSize <= 0 {
		p.pageSize = 10
	}
	return int(p.pageSize)
}

func (p *PagePagination) getPage() int {
	if p.pageNum <= 0 {
		p.pageNum = 1
	}
	return int(p.pageNum)
}

func (p *PagePagination) GetTotal(ctx context.Context) int64 {
	return p.total
}

func (p *PagePagination) GetResult() uint64 {
	return uint64(p.total)
}

type CursorPagination struct {
	cursor uint64
	next   int
	order  *Order
	total  int64
}

func NewCursorPagination(cursor uint64, next int, order *Order) Pagination {
	return &CursorPagination{
		cursor: cursor,
		next:   next,
		order:  order,
	}
}

func (p *CursorPagination) GetTotal(ctx context.Context) int64 {
	return p.total
}

func (p *CursorPagination) Query(ctx context.Context, table string, dao gen.DO) (gen.DO, error) {
	logger.GetLogger(ctx).Debug(fmt.Sprintf("nowPos:%+v, next:%+v\n", p.cursor, p.next))
	total, err := dao.Count()
	if err != nil {
		return gen.DO{}, err
	}
	p.total = total
	idField := field.NewUint64(table, "id")
	if p.order != nil {
		if p.order.IsDesc(ctx) {
			if p.cursor <= 0 {
				return gen.DO{}, uerrors.UErrorRepoRecordConstraint
			}
			dao = *dao.Where(idField.Lte(p.cursor)).Limit(p.next + 1).(*gen.DO)
		} else {
			dao = *dao.Where(idField.Gte(p.cursor)).Limit(p.next + 1).(*gen.DO)
		}
		dao, err = p.order.Query(ctx, table, dao)
		if err != nil {
			return gen.DO{}, err
		}
	} else {
		return gen.DO{}, uerrors.UErrorPaginationNotFoundOrder
	}

	return dao, nil
}
