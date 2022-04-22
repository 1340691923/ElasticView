package crontab

import (
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/1340691923/ElasticView/platform-basic-libs/service/data_conversion"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron"
	"log"
)

var Crontab *cron.Cron

func CrontabFn(reqData request.TransferReq, taskId int) {
	log.Println("ji", reqData, taskId)
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

/*scheduler := timing.GetTaskSchedulerInstance()
go scheduler.Start()
gmTimedList := model.GmTimedList{}
timedList, err := gmTimedList.GetTaskList()
if err != nil {
	logs.Logger.Sugar().Errorf("err", err)
	return
}
for _, timed := range timedList {
	timing.AddTask(timed.Action, timed.Data, timed.TaskId, timed.ExecTime)
}*/
