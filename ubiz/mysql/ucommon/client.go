package mysql

import (
	"context"
	"git.umu.work/AI/uglib/ubiz/mysql/ucommon/model"
	"git.umu.work/AI/uglib/ubiz/mysql/ucommon/query"
	"git.umu.work/be/goframework/store/gorm"
)

type UcommonDBClient struct {
	DB    *gorm.DB
	Query *query.Query
}

func NewUcommonDB(ctx context.Context) *UcommonDBClient {
	db, err := gorm.GetDB(ctx, model.DBCOMMON)
	if err != nil {
		panic(err)
	}

	return &UcommonDBClient{
		DB:    db,
		Query: query.Use(db),
	}
}
