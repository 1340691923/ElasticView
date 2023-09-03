package alias_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"

	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/my_error"
	"github.com/1340691923/ElasticView/pkg/vo"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/errgroup"
)

type AliasService struct{}

func NewAliasService() *AliasService {
	return &AliasService{}
}

func (this *AliasService) EsIndexGetAlias(ctx context.Context, esClient pkg.EsI, esAliasInfo *dto.EsAliasInfo) (res []vo.AliasInfo, err error) {
	if esAliasInfo.IndexName == "" {
		err = my_error.NewBusiness(my_error.ParmasNullError, my_error.IndexNameNullError)
		return
	}

	aliasRes, err := esClient.GetAliases(ctx, []string{esAliasInfo.IndexName})
	if err != nil {
		return
	}

	if aliasRes.StatusErr() != nil {
		err = aliasRes.StatusErr()
		return
	}
	gjson.GetBytes(
		aliasRes.ResByte(),
		fmt.Sprintf("%s.aliases", esAliasInfo.IndexName),
	).ForEach(func(key, value gjson.Result) bool {
		res = append(res, vo.AliasInfo{AliasName: cast.ToString(key)})
		return true
	})
	return
}

func (this *AliasService) MoveAliasToIndex(ctx context.Context, esClient pkg.EsI, esAliasInfo *dto.EsAliasInfo) (err error) {

	_, err = esClient.MoveToAnotherIndexAliases(
		ctx,
		proto.AliasAction{Actions: []proto.AliasAddAction{
			{
				Add: proto.AliasAdd{
					Indices: esAliasInfo.NewIndexList,
					Alias:   esAliasInfo.AliasName,
				},
			},
		}})

	if err != nil {
		return
	}

	return
}

func (this *AliasService) AddAliasToIndex(ctx context.Context, esClient pkg.EsI, esAliasInfo *dto.EsAliasInfo) (err error) {
	if esAliasInfo.IndexName == "" {
		err = my_error.NewBusiness(my_error.ParmasNullError, my_error.IndexNameNullError)
		return
	}
	resp, err := esClient.AddAliases(ctx, []string{esAliasInfo.IndexName}, esAliasInfo.AliasName)

	if err != nil {
		return
	}

	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}

	return
}

func (this *AliasService) BatchAddAliasToIndex(ctx context.Context, esClient pkg.EsI, esAliasInfo *dto.EsAliasInfo) (err error) {
	if esAliasInfo.IndexName == "" {
		err = my_error.NewBusiness(my_error.ParmasNullError, my_error.IndexNameNullError)
		return
	}
	eg := errgroup.Group{}
	if len(esAliasInfo.NewAliasNameList) > 10 {
		err = errors.New("别名列表数量不能大于10")
		return
	}
	for _, aliasName := range esAliasInfo.NewAliasNameList {
		aliasName := aliasName
		eg.Go(func() error {
			resp, err := esClient.AddAliases(ctx, []string{esAliasInfo.IndexName}, aliasName)
			if err != nil {
				return err
			}
			if resp.StatusErr() != nil {
				return resp.StatusErr()
			}
			return nil
		})

	}
	err = eg.Wait()
	if err != nil {
		return
	}

	return
}

func (this *AliasService) RemoveAlias(ctx context.Context, esClient pkg.EsI, esAliasInfo *dto.EsAliasInfo) (err error) {
	if esAliasInfo.IndexName == "" {
		err = my_error.NewBusiness(my_error.ParmasNullError, my_error.IndexNameNullError)
		return
	}
	resp, err := esClient.RemoveAliases(
		ctx,
		[]string{esAliasInfo.IndexName},
		[]string{esAliasInfo.AliasName},
	)

	if err != nil {
		return
	}

	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}

	return
}
