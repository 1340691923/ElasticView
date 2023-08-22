package index_service

import (
	"context"
	"encoding/json"
	"github.com/1340691923/ElasticView/es_sdk/pkg"
	proto2 "github.com/1340691923/ElasticView/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/my_error"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/pkg/vo"
	"github.com/tidwall/gjson"
	"time"
)

type IndexService struct {
	esClient pkg.EsI
}

func NewIndexService(esClient pkg.EsI) *IndexService {
	return &IndexService{esClient: esClient}
}

func (this *IndexService) EsIndexCreate(ctx context.Context, esIndexInfo *dto.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		return my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError)
	}
	if esIndexInfo.Types == "update" {
		res, err := this.esClient.IndicesPutSettingsRequest(ctx,
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
		res, err := this.esClient.CreateIndex(ctx, proto2.IndicesCreateRequest{
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

func (this *IndexService) EsIndexDelete(ctx context.Context, esIndexInfo *escache.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		return my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError)
	}
	resp, err := this.esClient.DeleteIndex(ctx, proto2.IndicesDeleteRequest{
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

func (this *IndexService) EsIndexGetSettings(ctx context.Context, esIndexInfo *escache.EsIndexInfo) (settings map[string]interface{}, err error) {
	if esIndexInfo.IndexName == "" {
		err = my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError)
		return
	}

	res, err := this.esClient.IndicesGetSettingsRequest(ctx, proto2.IndicesGetSettingsRequest{
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

func (this *IndexService) EsIndexGetSettingsInfo(ctx context.Context, esIndexInfo *escache.EsIndexInfo) (settings map[string]interface{}, err error) {
	if esIndexInfo.IndexName == "" {
		err = my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError)
		return
	}

	res, err := this.esClient.IndicesGetSettingsRequest(ctx, proto2.IndicesGetSettingsRequest{
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

func (this *IndexService) EsIndexReindex(ctx context.Context, esReIndexInfo *escache.EsReIndexInfo) (res map[string]interface{}, err error) {

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

	reindexRes, err := this.esClient.Reindex(ctx, reindexRequest, esReIndexInfo.Body)
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

func (this *IndexService) EsIndexNames(ctx context.Context) (indexNames []string, err error) {
	catIndicesResponse, err := this.esClient.GetIndices(ctx, proto2.CatIndicesRequest{
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

func (this *IndexService) EsIndexCount(ctx context.Context) (indexNameLen int, err error) {
	catIndicesResponse, err := this.esClient.GetIndices(ctx, proto2.CatIndicesRequest{
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

func (this *IndexService) EsIndexStats(ctx context.Context, indexName string) (res []vo.Status, err error) {
	catIndicesResponse, err := this.esClient.GetIndices(ctx, proto2.CatIndicesRequest{
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

func (this *IndexService) EsIndexCatStatus(ctx context.Context, indexName string) (res []vo.Status, err error) {
	catIndicesResponse, err := this.esClient.GetIndices(ctx, proto2.CatIndicesRequest{
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

func (this *IndexService) Refresh(ctx context.Context, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	resp, err := this.esClient.Refresh(ctx, indexNames)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this *IndexService) CacheClear(ctx context.Context, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	resp, err := this.esClient.IndicesClearCache(ctx, indexNames)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this *IndexService) Close(ctx context.Context, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	resp, err := this.esClient.IndicesClose(ctx, indexNames)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this *IndexService) Empty(ctx context.Context, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	resp, err := this.esClient.DeleteByQuery(ctx, indexNames, nil, map[string]interface{}{
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

func (this *IndexService) Flush(ctx context.Context, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	resp, err := this.esClient.Flush(ctx, indexNames)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this *IndexService) IndicesForcemerge(ctx context.Context, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	maxNumSegments := 1
	resp, err := this.esClient.IndicesForcemerge(ctx, indexNames, &maxNumSegments)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this *IndexService) Open(ctx context.Context, indexNames []string) (err error) {
	if len(indexNames) == 0 {
		indexNames = []string{"*"}
	}
	resp, err := this.esClient.Open(ctx, indexNames)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	return
}

func (this IndexService) UpdateMapping(ctx context.Context, updateMapping *escache.UpdateMapping) (res json.RawMessage, err error) {
	resp, err := this.esClient.PutMapping(ctx, proto2.IndicesPutMappingRequest{
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

func (this IndexService) EsMappingList(ctx context.Context, esConnect *escache.EsMapGetProperties) (res json.RawMessage, err error) {
	indexNames := []string{}
	if esConnect.IndexName != "" {
		indexNames = []string{esConnect.IndexName}
	}
	resp, err := this.esClient.GetMapping(ctx, indexNames)
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
