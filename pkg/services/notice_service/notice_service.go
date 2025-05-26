package notice_service

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/consts"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager"
	"github.com/1340691923/ElasticView/pkg/services/live_svr"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"
	"time"
)

type NoticeService struct {
	live           *live_svr.Live
	log            *logger.AppLogger
	cfg            *config.Config
	dao            *dao.NoticeDao
	userDao        *dao.GmUserDao
	pluginRegistry manager.Service
}

func NewNoticeService(live *live_svr.Live, log *logger.AppLogger, cfg *config.Config, dao *dao.NoticeDao, userDao *dao.GmUserDao, pluginRegistry manager.Service) *NoticeService {
	return &NoticeService{live: live, log: log, cfg: cfg, dao: dao, userDao: userDao, pluginRegistry: pluginRegistry}
}

// 全部用户都接收
func (this *NoticeService) LiveBroadcastEvMsg2All(data *dto.NoticeData) (err error) {

	noticeData, err := this.sendLiveEvMsg(consts.EvAllMsgChannel, nil, data)

	return this.live.LiveBroadcast(consts.EvAllMsgChannel, noticeData)
}

// 部分权限组接收
func (this *NoticeService) LiveBroadcastEvMsg2Roles(roles []int, data *dto.NoticeData) (err error) {

	noticeData, err := this.sendLiveEvMsg(consts.EvRoleMsgChannel, roles, data)

	for _, roleId := range roles {
		err = this.live.LiveBroadcast(fmt.Sprintf("%s:%d", consts.EvRoleMsgChannel, roleId), noticeData)
		if err != nil {
			return err
		}
	}

	return nil
}

// 部分用户接收
func (this *NoticeService) LiveBroadcastEvMsg2Users(users []int, data *dto.NoticeData) (err error) {

	noticeData, err := this.sendLiveEvMsg(consts.EvUserMsgChannel, users, data)

	for _, userId := range users {
		err = this.live.LiveBroadcast(fmt.Sprintf("%s:%d", consts.EvUserMsgChannel, userId), noticeData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *NoticeService) sendLiveEvMsg(typ string, ids []int, data *dto.NoticeData) (res *model.Notice, err error) {
	p, has := this.pluginRegistry.Plugin(context.Background(), data.PluginAlias)

	if has {
		pluginName := p.PluginData().PluginJsonData.PluginName
		if data.Source != "" {
			data.Source = fmt.Sprintf("%s(%s)", data.Source, pluginName)
		} else {
			data.Source = pluginName
		}
	}

	if data.FromUid > 0 {
		var gmUser model.GmUserModel
		gmUser, err = this.userDao.GetUserById(context.Background(), data.FromUid)
		if err != nil {
			return
		}
		if gmUser.Realname != "" {
			data.Source = fmt.Sprintf("%s(%s)", data.Source, gmUser.Realname)
		}
	}
	noticePojo := &model.Notice{
		Title:       data.Title,
		Content:     data.Content,
		Type:        data.Type,
		Level:       string(data.Level),
		FromUid:     data.FromUid,
		PluginAlias: data.PluginAlias,
		Source:      data.Source,
		Created:     time.Now(),
		Updated:     time.Now(),
	}
	if data.IsTask {
		noticePojo.IsTask = 1
	}
	if data.NoticeJumpBtn != nil {
		noticePojo.BtnJumpType = data.NoticeJumpBtn.JumpType
		noticePojo.BtnDesc = data.NoticeJumpBtn.Text
		noticePojo.BtnJumpUrl = data.NoticeJumpBtn.JumpUrl
	}
	switch typ {
	case consts.EvAllMsgChannel:
		noticePojo.TargetType = "all"
		res, err = this.dao.CreateNotice(noticePojo, nil)
	case consts.EvRoleMsgChannel:
		noticePojo.TargetType = "roles"
		res, err = this.dao.CreateNotice(noticePojo, ids)
	case consts.EvUserMsgChannel:
		noticePojo.TargetType = "users"
		res, err = this.dao.CreateNotice(noticePojo, ids)
	}
	return
}

func (this *NoticeService) GetList(userID int, roleIDs []int,
	readType int, title string, page, pageSize int) (notices []model.Notice,
	total int64, err error) {
	return this.dao.GetUserNoticesWithReadStatus(userID, roleIDs, readType, title, page, pageSize)
}

func (this *NoticeService) MarkReadMsg(ctx context.Context, userID int, noticeId int) (err error) {
	return this.dao.MarkUserNoticeAsRead(ctx, userID, noticeId)
}

func (this *NoticeService) Truncate() {
	this.dao.Truncate()
}
