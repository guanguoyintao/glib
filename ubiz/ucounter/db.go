package ucounter

import (
	"context"
	"git.umu.work/AI/uglib/ubiz"
	mysql "git.umu.work/AI/uglib/ubiz/mysql/ucommon"
	"git.umu.work/AI/uglib/ubiz/mysql/ucommon/model"
	"git.umu.work/AI/uglib/ubiz/mysql/ucommon/query"
	"git.umu.work/AI/uglib/uerrors"
	"gorm.io/gorm/clause"
)

type DBCounter struct {
	namespace CounterNameSpaceType
	dBClient  *mysql.UcommonDBClient
}

func NewDBCounter(ctx context.Context, namespace CounterNameSpaceType) (ubiz.UCounter, error) {
	counter := &DBCounter{
		namespace: namespace,
		dBClient:  mysql.NewUcommonDB(ctx),
	}

	return counter, nil
}

func (d *DBCounter) Get(ctx context.Context, key uint64) (uint32, error) {
	counterModel := d.dBClient.Query.Counter
	dos, err := counterModel.WithContext(ctx).Where(counterModel.MsgID.Eq(key), counterModel.CounterKey.Eq(string(d.namespace))).Find()
	if err != nil {
		return 0, err
	}
	if dos == nil || len(dos) == 0 {
		err = counterModel.WithContext(ctx).Create(&model.Counter{
			MsgID:      key,
			CounterKey: string(d.namespace),
			CounterNum: 0,
		})
		if err != nil {
			return 0, err
		}

		return 0, nil
	} else if len(dos) > 1 {
		return 0, uerrors.UErrorRepoRecordConstraint
	}

	return dos[0].CounterNum, nil
}

func (d *DBCounter) Incr(ctx context.Context, key uint64, value uint32) (uint32, error) {
	var result uint32
	err := d.dBClient.Query.Transaction(func(tx *query.Query) error {
		counterModel := tx.Counter
		dos, err := counterModel.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(counterModel.MsgID.Eq(key), counterModel.CounterKey.Eq(string(d.namespace))).Find()
		if err != nil {
			return err
		}
		if dos == nil && len(dos) == 0 {
			err = counterModel.WithContext(ctx).Create(&model.Counter{
				MsgID:      key,
				CounterKey: string(d.namespace),
				CounterNum: value,
			})
			if err != nil {
				return err
			}
			result = 0
		} else if len(dos) > 1 {
			return uerrors.UErrorRepoRecordConstraint
		} else if len(dos) == 1 {
			_, err := counterModel.WithContext(ctx).Where(counterModel.MsgID.Eq(key), counterModel.CounterKey.Eq(string(d.namespace))).
				UpdateSimple(counterModel.CounterNum.Add(value))
			if err != nil {
				return err
			}
			result = dos[0].CounterNum
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (d *DBCounter) Decr(ctx context.Context, key uint64, value uint32) (uint32, error) {
	var result uint32
	err := d.dBClient.Query.Transaction(func(tx *query.Query) error {
		counterModel := tx.Counter
		dos, err := counterModel.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(counterModel.MsgID.Eq(key), counterModel.CounterKey.Eq(string(d.namespace))).Find()
		if err != nil {
			return err
		}
		if dos == nil && len(dos) == 0 {
			err = counterModel.WithContext(ctx).Create(&model.Counter{
				MsgID:      key,
				CounterKey: string(d.namespace),
				CounterNum: value,
			})
			if err != nil {
				return err
			}
			result = 0
		} else if len(dos) > 1 {
			return uerrors.UErrorRepoRecordConstraint
		} else if len(dos) == 1 {
			_, err := counterModel.WithContext(ctx).Where(counterModel.MsgID.Eq(key), counterModel.CounterKey.Eq(string(d.namespace))).
				UpdateSimple(counterModel.CounterNum.Sub(value))
			if err != nil {
				return err
			}
			result = dos[0].CounterNum
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return result, nil
}
