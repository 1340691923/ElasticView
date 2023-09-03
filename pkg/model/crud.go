package model

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
)

type Curd interface {
	ProcessSqlWhere(sqlA sqlstore.SelectBuilder) sqlstore.SelectBuilder
	TableName() string
	ProcessSqlInsert(sqlA sqlstore.InsertBuilder) sqlstore.InsertBuilder
	ProcessSqlUpdate(id int, sqlA sqlstore.UpdateBuilder) sqlstore.UpdateBuilder
	GetId() int
}

func SearchList(curd Curd, page, limit int, columns string, list interface{}, sqlx *sqlstore.SqlStore, log *logger.AppLogger) (err error) {
	sqlA := sqlstore.SqlBuilder.
		Select(columns).
		From(curd.TableName())
	sqlA = curd.ProcessSqlWhere(sqlA)
	sql, args, err := sqlA.
		Limit(uint64(limit)).
		Offset(sqlstore.CreatePage(int(page), int(limit))).
		OrderBy("id desc").
		ToSql()
	log.Sugar().Info(sql, args)
	err = sqlx.Select(list, sql, args...)
	return
}

func SearchAll(curd Curd, columns string, list interface{}, sqlx *sqlstore.SqlStore, log *logger.AppLogger) (err error) {
	sqlA := sqlstore.SqlBuilder.
		Select(columns).
		From(curd.TableName())
	sqlA = curd.ProcessSqlWhere(sqlA)
	sql, args, err := sqlA.
		OrderBy("id desc").
		ToSql()
	log.Sugar().Info(sql, args)
	err = sqlx.Select(list, sql, args...)
	return
}

func Count(curd Curd, sqlx *sqlstore.SqlStore, log *logger.AppLogger) (count int, err error) {
	sqlA := sqlstore.SqlBuilder.
		Select("count(*)").
		From(curd.TableName())
	sqlA = curd.ProcessSqlWhere(sqlA)
	sql, args, err := sqlA.ToSql()
	err = sqlx.Get(&count, sql, args...)
	return
}

func Delete(service Curd, sqlx *sqlstore.SqlStore, log *logger.AppLogger) (err error) {
	sql, args, err := sqlstore.SqlBuilder.Delete(service.TableName()).Where(sqlstore.Eq{"id": service.GetId()}).ToSql()
	if err != nil {
		return
	}
	log.Sugar().Info(sql, args)
	_, err = sqlx.Exec(sql, args...)
	return
}

func Insert(service Curd, sqlx *sqlstore.SqlStore, log *logger.AppLogger) (lastInsertId int, err error) {
	sqlA := sqlstore.SqlBuilder.Insert(service.TableName())

	sqlA = service.ProcessSqlInsert(sqlA)

	sql, args, err := sqlA.ToSql()
	log.Sugar().Info(sql, args)
	res, err := sqlx.Exec(sql, args...)
	if err != nil {
		return
	}
	id, err := res.LastInsertId()

	lastInsertId = int(id)
	return
}

func Update(service Curd, sqlx *sqlstore.SqlStore, log *logger.AppLogger) (err error) {
	sqlA := sqlstore.SqlBuilder.Update(service.TableName())
	sqlA = service.ProcessSqlUpdate(service.GetId(), sqlA)
	sql, args, err := sqlA.ToSql()
	log.Sugar().Info(sql, args)
	_, err = sqlx.Exec(sql, args...)
	return
}
