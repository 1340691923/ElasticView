package crontab

import (
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/pkg/core"
	"github.com/1340691923/ElasticView/pkg/engine/db"
	"github.com/1340691923/ElasticView/pkg/request"
	"github.com/1340691923/ElasticView/service/data_conversion"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron"
	"log"
)

var Crontab *cron.Cron

func init() {
	core.Register(core.LastLevel, "計劃任務", InitTask)
}

// 初始化项目启动任务
func InitTask() (fn func(), err error) {
	fn = func() {}

	esLinkModel := model.EsLinkModel{}
	if err = esLinkModel.FlushEsLinkList(); err != nil {
		return fn, err
	}

	Crontab, err = InitCrontab()
	if err != nil {
		return fn, err
	}
	fn = func() {
		Crontab.Stop()
	}
	return fn, err
}

func CrontabFn(reqData request.TransferReq, taskId int) {
	selectType, err := reqData.ParseSelectType()
	if err != nil {
		log.Println("计划任务运行异常：", err)
		return
	}

	var obj model.DataxLinkInfoModel
	err = db.Sqlx.Get(&obj, "select * from datax_link_info where id = ?", selectType.ID)
	if err != nil {
		log.Println("计划任务运行异常：", err)
	}
	dataSource, err := data_conversion.NewDataSource(request.DataxInfoTestLinkReq{
		IP:       obj.Ip,
		Port:     obj.Port,
		DbName:   obj.DbName,
		Username: obj.Username,
		Pwd:      obj.Pwd,
		Remark:   obj.Remark,
		Typ:      obj.Typ,
	})

	err = dataSource.Transfer(taskId, &reqData)
	if err != nil {
		log.Println("计划任务运行异常：", err)
		return
	}
}

// 初始化项目启动任务
func InitCrontab() (crontab *cron.Cron, err error) {

	crontab = cron.New()
	crontab.Start()

	sql, args, err := db.SqlBuilder.
		Select("*").
		From("datax_transfer_list").
		ToSql()
	if err != nil {
		return crontab, err
	}
	var list []model.DataxListModel
	err = db.Sqlx.Select(&list, sql, args...)
	if err != nil {
		return crontab, err
	}

	for _, v := range list {
		tempV := v
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		var reqData request.TransferReq
		err := json.UnmarshalFromString(tempV.FormData, &reqData)
		if err != nil {
			return crontab, err
		}

		if tempV.CrontabSpec != "" {
			crontab.AddFunc(tempV.CrontabSpec, func() {
				CrontabFn(reqData, tempV.Id)
			})
		}
	}

	return crontab, err
}
