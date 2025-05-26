package dao

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/api"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/vo"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type EvBackDao struct {
	log    *logger.AppLogger
	eveApi *eve_api.EvEApi
}

func NewEvBackDao(log *logger.AppLogger, eveApi *eve_api.EvEApi) *EvBackDao {
	return &EvBackDao{log: log, eveApi: eveApi}
}

func (this *EvBackDao) GetEvAccessToken(ctx context.Context, evKey string) (string, error) {
	res := vo.ApiCommonRes{
		Data: "",
	}
	err := this.eveApi.Request(ctx, api.GetEvAccessToken, &dto.EvKeyReq{EvKey: evKey}, &res)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if res.Error() != nil {
		return "", errors.WithStack(res.Error())
	}

	token := cast.ToString(res.Data)
	this.eveApi.SetAccessToken(token)

	return cast.ToString(res.Data), nil
}

func (this *EvBackDao) GetPluginList(ctx context.Context, req *dto.FromEvPluginReq) (*vo.PluginListRes, error) {
	res := vo.ApiCommonRes{Data: &vo.PluginListRes{}}
	err := this.eveApi.Request(ctx, api.GetPluginList, req, &res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return res.Data.(*vo.PluginListRes), nil
}

func (this *EvBackDao) GetWxArticleList(ctx context.Context) (*vo.WxArticleList, error) {

	res := vo.ApiCommonRes{Data: &vo.WxArticleList{}}
	err := this.eveApi.Request(ctx, api.GetWxArticleList, map[string]interface{}{}, &res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}
	if res.Data == nil {
		return nil, errors.New("暂无结果")
	}

	return res.Data.(*vo.WxArticleList), nil
}

func (this *EvBackDao) StarPlugin(ctx context.Context, req *dto.StarPlugin) error {
	res := vo.ApiCommonRes{}
	err := this.eveApi.Request(ctx, api.StarPlugin, req, &res)
	if err != nil {
		return errors.WithStack(err)
	}

	if res.Error() != nil {
		return errors.WithStack(res.Error())
	}

	return nil
}

func (this *EvBackDao) GetEvMaxVersion(ctx context.Context) (*vo.GetEvMaxVersionRes, error) {

	res := vo.ApiCommonRes{Data: &vo.GetEvMaxVersionRes{}}
	err := this.eveApi.Request(ctx, api.GetEvMaxVersion, &dto.Empty{}, &res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return res.Data.(*vo.GetEvMaxVersionRes), nil
}

func (this *EvBackDao) GetPluginDownloadUrl(ctx context.Context, req *dto.GetPluginDownloadUrlReq) (*vo.GetPluginDownloadUrlRes, error) {
	res := vo.ApiCommonRes{Data: &vo.GetPluginDownloadUrlRes{}}
	err := this.eveApi.Request(ctx, api.GetPluginDownloadUrl, req, &res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return res.Data.(*vo.GetPluginDownloadUrlRes), nil
}

func (this *EvBackDao) GetPluginInfo(ctx context.Context, req *dto.FormEvPluginInfoReq) (*vo.PublishRes, error) {
	res := vo.ApiCommonRes{Data: &vo.PublishRes{}}
	err := this.eveApi.Request(ctx, api.GetPluginInfo, req, &res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return res.Data.(*vo.PublishRes), nil
}

func (this *EvBackDao) GetEvPluginMaxVersion(ctx context.Context, pluginAlias string) (string, error) {
	res := vo.ApiCommonRes{}
	err := this.eveApi.Request(ctx, api.GetEvPluginMaxVersion, &dto.GetEvPluginMaxVersion{PluginAlias: pluginAlias}, &res)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if res.Error() != nil {
		return "", errors.WithStack(res.Error())
	}

	return cast.ToString(res.Data), nil
}

func (this *EvBackDao) GetEvPluginsMaxVersion(ctx context.Context, pluginIds []string) (map[string]interface{}, error) {
	res := vo.ApiCommonRes{
		Data: map[string]interface{}{},
	}
	err := this.eveApi.Request(ctx, api.GetEvPluginsMaxVersion, &dto.GetEvPluginsMaxVersion{PluginIds: pluginIds}, &res)
	if err != nil {
		return map[string]interface{}{}, errors.WithStack(err)
	}

	if res.Error() != nil {
		return map[string]interface{}{}, errors.WithStack(res.Error())
	}

	return res.Data.(map[string]interface{}), nil
}

func (this *EvBackDao) AddComment(ctx context.Context, req *dto.AddCommentRequest) error {
	res := vo.ApiCommonRes{}
	err := this.eveApi.Request(ctx, api.AddComment, req, &res)
	if err != nil {
		return errors.WithStack(err)
	}

	if res.Error() != nil {
		return errors.WithStack(res.Error())
	}

	return nil
}

func (this *EvBackDao) LikeComment(ctx context.Context, req *dto.LikeCommentRequest) error {
	res := vo.ApiCommonRes{}
	err := this.eveApi.Request(ctx, api.LikeComment, req, &res)
	if err != nil {
		return errors.WithStack(err)
	}

	if res.Error() != nil {
		return errors.WithStack(res.Error())
	}

	return nil
}

func (this *EvBackDao) ListComments(ctx context.Context, req *dto.ListCommentsRequest) (*[]*vo.Comment, error) {
	res := vo.ApiCommonRes{Data: &[]*vo.Comment{}}
	err := this.eveApi.Request(ctx, api.ListComments, req, &res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return res.Data.(*[]*vo.Comment), nil
}
