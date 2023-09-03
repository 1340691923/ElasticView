package model

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/util"
)

type GmOperaterLog struct {
	Id             int      `db:"id" json:"id"`
	OperaterName   string   `db:"operater_name" json:"operater_name"`     //操作者id
	OperaterId     int      `db:"operater_id" json:"operater_id"`         //操作者id
	OperaterAction string   `db:"operater_action" json:"operater_action"` //请求路由
	Created        string   `db:"created" json:"created"`
	Method         string   `db:"method" json:"method"` //请求方法
	Body           []byte   `db:"body" json:"-"`        //请求body
	OperaterRoleId int      `db:"operater_role_id" json:"operater_role_id"`
	FilterDate     []string `db:"-" json:"date"`
	BodyStr        string   `db:"-" json:"body_str"` //请求body
	Sqlx           *sqlstore.SqlStore
}

func (this *GmOperaterLog) ProcessSqlInsert(sqlA sqlstore.InsertBuilder) sqlstore.InsertBuilder {
	return sqlA
}

func (this *GmOperaterLog) ProcessSqlUpdate(id int, sqlA sqlstore.UpdateBuilder) sqlstore.UpdateBuilder {
	return sqlA
}

func (this *GmOperaterLog) TableName() string {
	return "gm_operater_log"
}

func (this *GmOperaterLog) ProcessSqlWhere(sqlA sqlstore.SelectBuilder) sqlstore.SelectBuilder {
	if this.OperaterId != 0 {
		sqlA = sqlA.Where(sqlstore.Eq{"operater_id": this.OperaterId})
	}
	if this.OperaterRoleId != 0 {
		sqlA = sqlA.Where(sqlstore.Eq{"operater_role_id": this.OperaterRoleId})
	}
	if this.OperaterAction != "" {
		sqlA = sqlA.Where(sqlstore.Eq{"operater_action": this.OperaterAction})
	}
	if len(this.FilterDate) == 2 {
		sqlA = sqlA.Where(sqlstore.GtOrEq{"created": this.FilterDate[0]})
		sqlA = sqlA.Where(sqlstore.LtOrEq{"created": this.FilterDate[1]})
	}
	return sqlA
}

func (this *GmOperaterLog) Insert() (err error) {
	body, err := util.GzipCompress(util.Bytes2str(this.Body))
	if err != nil {
		return
	}
	_, err = sqlstore.SqlBuilder.
		Insert(this.TableName()).
		SetMap(map[string]interface{}{
			"operater_name":    this.OperaterName,
			"operater_id":      this.OperaterId,
			"operater_action":  this.OperaterAction,
			"method":           this.Method,
			"body":             body,
			"operater_role_id": this.OperaterRoleId,
		}).RunWith(this.Sqlx).Exec()
	if err != nil {
		return
	}
	return
}

func (this *GmOperaterLog) GetId() int {
	return this.Id
}
