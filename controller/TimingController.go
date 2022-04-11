package controller

import (
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	"github.com/1340691923/ElasticView/platform-basic-libs/service/timing"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	"github.com/gofiber/fiber/v2"
	"time"
)

type TimingController struct {
	BaseController
}

func (this TimingController) CancelAction(ctx *fiber.Ctx) error {
	taskId := ctx.FormValue("taskId")
	taskScheduler := timing.GetTaskSchedulerInstance()
	taskScheduler.StopOnce(taskId)
	timingsModel := model.GmTimedList{}

	_, err := db.SqlBuilder.
		Update(timingsModel.TableName()).
		Where(db.Eq{"task_id": taskId}).
		Set("status", 4).
		Set("msg", "已取消").
		Set("updated", time.Now().Format(util.TimeFormat)).
		RunWith(db.Sqlx).
		Exec()
	if err != nil {
		return this.Error(ctx, err)

	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this TimingController) ListAction(ctx *fiber.Ctx) error {
	models := request.TimingModel{}
	err := ctx.BodyParser(&models)
	if err != nil {
		return this.Error(ctx, err)
	}
	if models.Page == 0 {
		models.Page = 1
	}
	if models.Limit == 0 {
		models.Limit = 10
	}

	timingsModel := model.GmTimedList{}

	sqlBuilder := db.SqlBuilder.Select("*").From(timingsModel.TableName())

	countbuilder := db.SqlBuilder.Select("count(*)").From(timingsModel.TableName())
	if models.Action != nil {
		action := &models.Action
		sqlBuilder = sqlBuilder.Where(db.Eq{"action": action})
		countbuilder = countbuilder.Where(db.Eq{"action": action})
	}

	if models.Status != nil {
		status := &models.Status
		sqlBuilder = sqlBuilder.Where(db.Eq{"status": status})
		countbuilder = countbuilder.Where(db.Eq{"status": status})
	}
	var count int
	err = countbuilder.RunWith(db.Sqlx).QueryRow().Scan(&count)
	if err != nil {
		return this.Error(ctx, err)
	}
	sql, args, err := sqlBuilder.
		OrderBy("id desc").
		Limit(uint64(models.Limit)).
		Offset(db.CreatePage(models.Page, models.Limit)).
		RunWith(db.Sqlx).
		ToSql()

	if err != nil {
		return this.Error(ctx, err)
	}

	timingsList := []model.GmTimedList{}
	err = db.Sqlx.Select(&timingsList, sql, args...)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, util.Map{
		"data":  timingsList,
		"count": count,
	})

}
