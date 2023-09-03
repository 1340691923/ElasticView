package index_service

import (
	"context"
	"encoding/json"
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg"
	proto2 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/pkg/my_error"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/pkg/vo"
	"github.com/tidwall/gjson"
	"time"
)

type IndexService struct{}

func NewIndexService() *IndexService {
	return &IndexService{}
}

func (this *IndexService) EsIndexCreate(ctx context.Context, esClient pkg.EsI, esIndexInfo *dto.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		return my_error.NewBusiness(my_error.ParmasNullError, my_error.IndexNameNullError)
	}
	if esIndexInfo.Types == "update" {
		res, err := esClient.IndicesPutSettingsRequest(ctx,
			proto2.IndicesPutSettingsRequest{
				Index: []string{esIndexInfo.IndexName},
			}, esIndexInfo.Settings,
		)
		if err != nil {
			return err
		}
		if res.StatusErr() != nil {
			err = res.StatusErr()
			return err
		}
	} else {
		res, err := esClient.CreateIndex(ctx, proto2.IndicesCreateRequest{
			Index: esIndexInfo.IndexName,
		},
			util.Map{
				"settings": esIndexInfo.Settings,
			})
		if err != nil {
			return err
		}
		if res.StatusErr() != nil {
			err = res.StatusErr()
			return err
		}

	}
	return nil
}

func (this *IndexService) EsIndexDelete(ctx context.Context, esClient pkg.EsI, esIndexInfo *dto.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		return my_error.NewBusiness(my_error.ParmasNullError, my_error.IndexNameNullError)
	}
	resp, err := esClient.DeleteIndex(ctx, proto2.IndicesDeleteRequest{
		Index: []string{esIndexInfo.IndexName},
	})
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return err
	}
	return
}

func (this *IndexService) EsIndexGetSettings(ctx context.Context, esClient pkg.EsI, esIndexInfo *dto.EsIndexInfo) (settings map[string]interface{}, err error) {
	if esIndexInfo.IndexName == "" {
		err = my_error.NewBusiness(my_error.ParmasNullError, my_error.IndexNameNullError)
		return
	}

	res, err := esClient.IndicesGetSettingsRequest(ctx, proto2.IndicesGetSettingsRequest{
		Index: []string{esIndexInfo.IndexName},
	})
	if err != nil {
		return
	}

	if res.StatusErr() != nil {
		err = res.StatusErr()
		return
	}

	settings = map[string]interface{}{}
	gjson.GetBytes(res.ResByte(), esIndexInfo.IndexName).Get("settings").ForEach(func(key, value gjson.Result) bool {
		settings[key.String()] = value.Value()
		return true
	})

	return
}

func (this *IndexService) EsIndexGetSettingsInfo(ctx context.Context, esClient pkg.EsI, esIndexInfo *dto.EsIndexInfo) (settings map[string]interface{}, err error) {
	if esIndexInfo.IndexName == "" {
		err = my_error.NewBusiness(my_error.ParmasNullError, my_error.IndexNameNullError)
		return
	}

	res, err := esClient.IndicesGetSettingsRequest(ctx, proto2.IndicesGetSettingsRequest{
		Index: []string{esIndexInfo.IndexName},
	})
	if err != nil {
		return
	}
	if res.StatusErr() != nil {
		err = res.StatusErr()
		return
	}

	settings = map[string]interface{}{}
	err = json.Unmarshal(res.ResByte(), &settings)
	if err != nil {
		return
	}
	return
}

func (this *IndexService) EsIndexReindex(ctx context.Context, esClient pkg.EsI, esReIndexInfo *dto.EsReIndexInfo) (res map[string]interface{}, err error) {

	reindexRequest := proto2.ReindexRequest{}

	urlValues := esReIndexInfo.UrlValues
	if urlValues.WaitForActiveShards != "" {
		reindexRequest.WaitForActiveShards = urlValues.WaitForActiveShards
	}
	if urlValues.Slices != 0 {
		reindexRequest.Slices = urlValues.Slices
	}
	if urlValues.Refresh != nil {
		reindexRequest.Refresh = urlValues.Refresh
	}
	if urlValues.Timeout != 0 {
		reindexRequest.Timeout = time.Duration(int64(urlValues.Timeout)) * time.Second
	}
	if urlValues.RequestsPerSecond != 0 {
		requestsPerSecond := urlValues.RequestsPerSecond
		reindexRequest.RequestsPerSecond = &requestsPerSecond
	}
	if urlValues.WaitForCompletion != nil {
		reindexRequest.WaitForCompletion = urlValues.WaitForCompletion
	}

	reindexRes, err := esClient.Reindex(ctx, reindexRequest, esReIndexInfo.Body)
	if err != nil {
		return
	}
	if reindexRes.StatusErr() != nil {
		err = reindexRes.StatusErr()
		return
	}

	res = map[string]interface{}{}
	err = json.Unmarshal(reindexRes.ResByte(), &res)
	if err != nil {
		return
	}
	return
}

func (this *IndexService) EsIndexNames(ctx context.Context, esClient pkg.EsI) (indexNames []string, err error) {
	catIndicesResponse, err := esClient.GetIndices(ctx, proto2.CatIndicesRequest{
		Format: "json",
	})
	if err != nil {
		return
	}
	if catIndicesResponse.StatusErr() != nil {
		err = catIndicesResponse.StatusErr()
		return
	}
	var list []proto2.CatIndex
	err = json.Unmarshal(catIndicesResponse.ResByte(), &list)
	if err != nil {
		return
	}
	for _, v := range list {
		indexNames = append(indexNames, v.Index)
	}
	return
}

func (this *IndexService) EsIndexCount(ctx context.Context, esClient pkg.EsI) (indexNameLen int, err error) {
	catIndicesResponse, err := esClient.GetIndices(ctx, proto2.CatIndicesRequest{
		Format: "json",
	})
	if err != nil {
		return
	}
	if catIndicesResponse.StatusErr() != nil {
		err = catIndicesResponse.StatusErr()
		return
	}
	var list []proto2.CatIndex
	err = json.Unmarshal(catIndicesResponse.ResByte(), &list)
	if err != nil {
		return
	}
	indexNameLen = len(list)
	return
}

func (this *IndexService) EsIndexStats(ctx context.Context, esClient pkg.EsI, indexName string) (res []vo.Status, err error) {
	catIndicesResponse, err := esClient.GetIndices(ctx, proto2.CatIndicesRequest{
		Index:  []string{indexName},
		H:      []string{"status"},
		Format: "json",
	})
	if err != nil {
		return
	}
	if catIndicesResponse.StatusErr() != nil {
		err = catIndicesResponse.StatusErr()
		return
	}
	err = json.Unmarshal(catIndicesResponse.ResByte(), &res)
	return
}

func (this *IndexService) EsIndexCatStatus(ctx context.Context, esClient pkg.EsI, indexName string) (res []vo.Status, err error) {
	catIndicesResponse, err := esClient.GetIndices(ctx, proto2.CatIndicesRequest{
		Index:  []string{indexName},
		H:      []string{"status"},
		Format: "json",
	})
	if err != nil {
		return
	}
	if catIndicesResponse.StatusErr() != nil {
		err = catIndicesResponse.StatusErr()
		return
	}

	err = json.Unmarshal(catIndicesResponse.ResByte(), &res)
	return
}

func (this *IndexService) Refresh(ctx context.Context, esClient pkg.EsI, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	resp, err := esClient.Refresh(ctx, indexNames)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this *IndexService) CacheClear(ctx context.Context, esClient pkg.EsI, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	resp, err := esClient.IndicesClearCache(ctx, indexNames)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this *IndexService) Close(ctx context.Context, esClient pkg.EsI, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	resp, err := esClient.IndicesClose(ctx, indexNames)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this *IndexService) Empty(ctx context.Context, esClient pkg.EsI, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	resp, err := esClient.DeleteByQuery(ctx, indexNames, nil, map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	})
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this *IndexService) Flush(ctx context.Context, esClient pkg.EsI, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	resp, err := esClient.Flush(ctx, indexNames)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this *IndexService) IndicesForcemerge(ctx context.Context, esClient pkg.EsI, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	maxNumSegments := 1
	resp, err := esClient.IndicesForcemerge(ctx, indexNames, &maxNumSegments)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this *IndexService) Open(ctx context.Context, esClient pkg.EsI, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	resp, err := esClient.Open(ctx, indexNames)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this *IndexService) UpdateMapping(ctx context.Context, esClient pkg.EsI, updateMapping *dto.UpdateMapping) (res json.RawMessage, err error) {
	resp, err := esClient.PutMapping(ctx, proto2.IndicesPutMappingRequest{
		Index:        []string{updateMapping.IndexName},
		DocumentType: updateMapping.TypeName,
	}, updateMapping.Properties)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	res = resp.JsonRawMessage()

	return
}

func (this *IndexService) EsMappingList(ctx context.Context, esClient pkg.EsI, esConnect *dto.EsMapGetProperties) (res json.RawMessage, err error) {
	indexNames := []string{}
	if esConnect.IndexName != "" {
		indexNames = []string{esConnect.IndexName}
	}
	resp, err := esClient.GetMapping(ctx, indexNames)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}

	res = resp.JsonRawMessage()

	return
}
