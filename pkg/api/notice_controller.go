package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/infrastructure/response"
	"github.com/1340691923/ElasticView/pkg/services/gm_user"
	"github.com/1340691923/ElasticView/pkg/services/notice_service"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// BI用户控制器
type NoticeController struct {
	*BaseController
	log           *logger.AppLogger
	cfg           *config.Config
	orm           *orm.Gorm
	jwtSvr        *jwt_svr.Jwt
	gmUserService *gm_user.GmUserService
	noticeService *notice_service.NoticeService
}

func NewNoticeController(baseController *BaseController, log *logger.AppLogger, cfg *config.Config, orm *orm.Gorm, jwtSvr *jwt_svr.Jwt, gmUserService *gm_user.GmUserService, noticeService *notice_service.NoticeService) *NoticeController {
	return &NoticeController{BaseController: baseController, log: log, cfg: cfg, orm: orm, jwtSvr: jwtSvr, gmUserService: gmUserService, noticeService: noticeService}
}

func (this *NoticeController) GetList(ctx *gin.Context) {

	var req dto.NoticeReq

	err := ctx.ShouldBind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	cliams, err := this.jwtSvr.ParseToken(this.GetToken(ctx))
	if err != nil {
		this.Error(ctx, err)
		return
	}
	roleIds, err := this.gmUserService.GetRolesByUserID(cliams.UserID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	list, count, err := this.noticeService.GetList(cliams.UserID, roleIds, req.ReadType, req.Title, req.Page, req.PageSize)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, map[string]interface{}{
		"count": count,
		"list":  list,
	})

}

func (this *NoticeController) MarkReadNotice(ctx *gin.Context) {

	var req dto.MarkReadNoticeReq

	err := ctx.ShouldBind(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	cliams, err := this.jwtSvr.ParseToken(this.GetToken(ctx))
	if err != nil {
		this.Error(ctx, err)
		return
	}
	eg := errgroup.Group{}
	eg.SetLimit(50)
	for _, id := range req.Ids {
		id := id
		eg.Go(func() error {
			err = this.noticeService.MarkReadMsg(ctx, cliams.UserID, id)
			if err != nil {
				this.log.Sugar().Errorf("MarkReadMsg err", err)
			}
			return nil
		})
	}
	eg.Wait()

	this.Success(ctx, response.OperateSuccess, nil)
}

func (this *NoticeController) Truncate(ctx *gin.Context) {
	this.noticeService.Truncate()
	this.Success(ctx, response.OperateSuccess, nil)
}
