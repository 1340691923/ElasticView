package navicat_service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg"
	proto2 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/tidwall/gjson"
	"log"
	"sort"
)

type NavicatService struct{}

func NewNavicatService() *NavicatService {
	return &NavicatService{}
}

func (this *NavicatService) CrudGetList(ctx context.Context, esClient pkg.EsI, crudFilter *dto.CrudFilter) (res json.RawMessage, count int64, err error) {
	q, err := GetWhereSql(crudFilter.Relation)
	if err != nil {
		return
	}

	queryBody, err := q.Source()
	if err != nil {
		return
	}
	req := proto2.SearchRequest{
		Index: []string{crudFilter.IndexName},
	}

	sortArr := []map[string]interface{}{}
	for _, tmp := range crudFilter.SortList {
		sortArr = append(sortArr, map[string]interface{}{tmp.Col: tmp.SortRule})
	}

	searchBody := Search{
		Query: queryBody,
		From:  int(sqlstore.CreatePage(crudFilter.Page, crudFilter.Limit)),
		Size:  crudFilter.Limit,
		Sort:  sortArr,
	}

	resp, err := esClient.Search(
		ctx, req, searchBody,
	)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}
	res = resp.JsonRawMessage()

	switch esClient.Version() {
	case 6:
		count = gjson.GetBytes(resp.ResByte(), "hits.total").Int()
	default:
		count = gjson.GetBytes(resp.ResByte(), "hits.total.value").Int()
	}

	return
}

func (this *NavicatService) CrudGetDSL(ctx context.Context, esClient pkg.EsI, crudFilter *dto.CrudFilter) (res Search, err error) {
	q, err := GetWhereSql(crudFilter.Relation)
	if err != nil {
		return
	}

	queryBody, err := q.Source()
	if err != nil {
		return
	}
	sortArr := []map[string]interface{}{}
	for _, tmp := range crudFilter.SortList {
		sortArr = append(sortArr, map[string]interface{}{tmp.Col: tmp.SortRule})
	}

	res = Search{
		Query: queryBody,
		From:  int(sqlstore.CreatePage(crudFilter.Page, crudFilter.Limit)),
		Size:  crudFilter.Limit,
		Sort:  sortArr,
	}

	return
}

func (this *NavicatService) CrudDownload(ctx context.Context, esClient pkg.EsI, filter *dto.CrudFilter) (downloadFileName string, titleList []string, searchData [][]string, err error) {

	mappingRes, err := esClient.GetMapping(ctx, []string{filter.IndexName})
	if err != nil {
		log.Println("err", err)
		return
	}
	if mappingRes.StatusErr() != nil {
		err = mappingRes.StatusErr()
		return
	}

	fields := map[string]interface{}{}

	err = json.Unmarshal(mappingRes.ResByte(), &fields)

	if err != nil {
		log.Println("err", err)
		return
	}

	fieldsArr := []string{"_index", "_type", "_id"}
	data, ok := fields[filter.IndexName].(map[string]interface{})

	if !ok {
		err = errors.New("该索引没有映射结构")
		return
	}
	propertiesArr := []string{}
	properties := map[string]interface{}{}
	mappings, ok := data["mappings"].(map[string]interface{})
	if !ok {
		err = errors.New("该索引没有映射结构")
		return
	}
	switch esClient.Version() {
	case 6:

		typeName := ""

		for key := range mappings {
			typeName = key
		}

		typeObj := mappings[typeName].(map[string]interface{})

		properties, ok = typeObj["properties"].(map[string]interface{})
		if !ok {
			err = errors.New("该索引没有映射结构")
			return
		}

	default:

		properties, ok = mappings["properties"].(map[string]interface{})
		if !ok {
			err = errors.New("该索引没有映射结构3")
			return
		}

	}

	for key := range properties {
		propertiesArr = append(propertiesArr, key)
	}

	sort.Strings(propertiesArr)
	fieldsArr = append(fieldsArr, propertiesArr...)
	q, err := GetWhereSql(filter.Relation)
	if err != nil {
		log.Println("err", err)
		return
	}

	querySource, err := q.Source()

	if err != nil {
		log.Println("err", err)
		return
	}

	searchBody := Search{
		Query: querySource,
		Size:  8000,
		Sort: []map[string]interface{}{
			{
				"_id": "desc",
			},
		},
	}

	resp, err := esClient.Search(
		ctx, proto2.SearchRequest{Index: []string{filter.IndexName}}, searchBody,
	)
	if err != nil {
		return
	}
	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}

	histsArr := gjson.GetBytes(resp.ResByte(), "hits.hits").Array()

	lastIdArr := histsArr[len(histsArr)-1].Get("sort").Array()

	llist := [][]string{}

	flushHitsDataFn := func(hits []gjson.Result) {
		for _, data := range hits {
			list := []string{data.Get("_index").String(), data.Get("_type").String(), data.Get("_id").String()}

			m := data.Get("_source").Map()

			for _, field := range fieldsArr {
				if field == "_index" || field == "_type" || field == "_id" {
					continue
				}
				if value, ok := m[field]; ok {
					list = append(list, util.ToExcelData(value.Value()))
				} else {
					list = append(list, "")
				}
			}

			llist = append(llist, list)
		}
	}

	flushHitsDataFn(histsArr)
	haveData := true

	for haveData {
		searchAfter := []interface{}{}
		for _, v := range lastIdArr {
			searchAfter = append(searchAfter, v.Value())
		}

		searchBody := Search{
			Query: querySource,
			Size:  8000,
			Sort: []map[string]interface{}{
				{
					"_id": "desc",
				},
			},
			SearchAfter: &searchAfter,
		}
		var searchResp *proto2.Response
		searchResp, err = esClient.Search(ctx, proto2.SearchRequest{Index: []string{filter.IndexName}}, searchBody)

		if err != nil {
			log.Println("err", err)
			return
		}

		if searchResp.StatusErr() != nil {
			err = searchResp.StatusErr()
			return
		}

		hitsArr := gjson.GetBytes(searchResp.ResByte(), "hits.hits").Array()

		if len(hitsArr) == 0 {
			break
		}

		lastIdArr = hitsArr[len(hitsArr)-1].Get("sort").Array()
		flushHitsDataFn(hitsArr)
	}
	downloadFileName = "test"
	titleList = fieldsArr
	searchData = llist

	return

}
