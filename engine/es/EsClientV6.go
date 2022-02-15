package es

import (
	"context"
	"errors"

	elasticV6 "github.com/olivere/elastic"
)

type EsClientV6 struct {
	Client          *elasticV6.Client
	esConnectConfig EsConnect
}

func NewEsClientV6(esConnectConfig *EsConnect) (esClient *elasticV6.Client, err error) {

	optList := []elasticV6.ClientOptionFunc{
		elasticV6.SetSniff(false),
		elasticV6.SetHealthcheck(false),
	}

	optList = append(optList, elasticV6.SetURL(esConnectConfig.Ip))

	if esConnectConfig.User != "" || esConnectConfig.Pwd != "" {
		optList = append(optList, elasticV6.SetBasicAuth(esConnectConfig.User, esConnectConfig.Pwd))
	}

	esClient, err = elasticV6.NewSimpleClient(optList...)

	if err != nil {
		return
	}

	return
}

func (this *EsClientV6) Ping() (interface{}, int, error) {
	return this.Client.Ping(this.esConnectConfig.Ip).Do(context.Background())
}

func (this *EsClientV6) CreateIndex(indexName string, body interface{}) (interface{}, error) {
	return this.Client.CreateIndex(indexName).BodyJson(body).Do(context.Background())
}

func (this *EsClientV6) CatIndices() (interface{}, error) {
	return this.Client.CatIndices().Pretty(true).Do(context.Background())
}

func (this *EsClientV6) CatAliases() (interface{}, error) {
	return this.Client.CatAliases().Pretty(true).Do(context.Background())
}

func (this *EsClientV6) CatAllocation() (interface{}, error) {
	return this.Client.CatAllocation().Pretty(true).Do(context.Background())
}

func (this *EsClientV6) CatCount() (interface{}, error) {
	return this.Client.CatCount().Pretty(true).Do(context.Background())
}

func (this *EsClientV6) CatHealth() (interface{}, error) {
	return this.Client.CatHealth().Pretty(true).Do(context.Background())
}

func (this *EsClientV6) CatShards() (interface{}, error) {
	return this.Client.CatShards().Pretty(true).Do(context.Background())
}

func (this *EsClientV6) DeleteIndex(indexNameList []string) (interface{}, error) {
	return this.Client.DeleteIndex(indexNameList...).Do(context.Background())
}

func (this *EsClientV6) CloseIndex(indexNameList []string) error {
	var err error
	for _, indexName := range indexNameList {
		_, err = this.Client.CloseIndex(indexName).Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *EsClientV6) OpenIndex(indexNameList []string) error {
	var err error
	for _, indexName := range indexNameList {
		_, err = this.Client.OpenIndex(indexName).Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *EsClientV6) FreezeIndex(indexNameList []string) error {
	return nil
}

func (this *EsClientV6) UnfreezeIndex(indexNameList []string) error {
	return nil
}

func (this *EsClientV6) CreateMapping(indexName string, body Json) (interface{}, error) {
	return this.Client.PutMapping().
		Index(indexName).
		BodyJson(body).
		IncludeTypeName(true).
		Do(context.Background())
}

func (this *EsClientV6) IndexPutSettings(indexName string, body Json) (interface{}, error) {
	return this.Client.IndexPutSettings().Index(indexName).BodyJson(body).Do(context.TODO())
}

func (this *EsClientV6) Refresh(indexName ...string) (interface{}, error) {
	return this.Client.Refresh(indexName...).Do(context.TODO())
}

func (this *EsClientV6) Flush(indexName ...string) (interface{}, error) {
	return this.Client.Flush(indexName...).Do(context.TODO())
}

func (this *EsClientV6) Fsync(indexName ...string) (interface{}, error) {
	return this.Client.SyncedFlush(indexName...).Do(context.TODO())
}
func (this *EsClientV6) Rollover(alias, indexName string) (interface{}, error) {
	return this.Client.RolloverIndex(alias).NewIndex(indexName).Do(context.TODO())
}

func (this *EsClientV6) IndexSegments(indexName string) (interface{}, error) {
	return this.Client.IndexSegments(indexName).Pretty(true).Do(context.TODO())
}

func (this *EsClientV6) Alias(alias, indexName string) (interface{}, error) {

	return this.Client.Alias().
		Action(elasticV6.NewAliasAddAction(alias).Index(indexName).IsWriteIndex(false)).
		Do(context.TODO())
}

func (this *EsClientV6) Reindex(sourceIndex, destinationIndex string) (interface{}, error) {
	return this.Client.Reindex().SourceIndex(sourceIndex).DestinationIndex(destinationIndex).Do(context.Background())
}

func (this *EsClientV6) TasksList() (interface{}, error) {
	return this.Client.TasksList().
		Pretty(true).
		Human(true).
		Detailed(true).
		Do(context.TODO())
}

func (this *EsClientV6) IndexStats(indexName []string, metrics []string) (interface{}, error) {
	return this.Client.IndexStats().Index(indexName...).Metric(metrics...).Do(context.Background())
}

func (this *EsClientV6) GetMapping(indexName string, body Json, typeName ...string) (interface{}, error) {
	return this.Client.GetMapping().Index(indexName).Type(typeName...).Do(context.Background())
}

func (this *EsClientV6) GetMappings() (interface{}, error) {
	return this.Client.GetMapping().Do(context.Background())
}

func (this *EsClientV6) PutData(indexName string, body Json, typeName ...string) (interface{}, error) {
	if len(typeName) == 0 {
		return nil, errors.New("Type 不能为空!")
	}
	return this.Client.Index().Index(indexName).Type(typeName[0]).BodyJson(body).Do(context.Background())
}

func (this *EsClientV6) DeleteById(indexName, id string, typeName ...string) (interface{}, error) {
	if len(typeName) == 0 {
		return nil, errors.New("Type 不能为空!")
	}
	return this.Client.Delete().Index(indexName).Type(typeName[0]).Id(id).Do(context.Background())
}

func (this *EsClientV6) Search(
	indexName string,
	query elasticV6.Query,
	sort *Sort,
	page *Page,
	fields []string,
	isInclude bool,
	typeName ...string,
) (*elasticV6.SearchResult, error) {

	search := this.Client.Search(indexName).Query(query)
	if sort != nil {
		search = search.Sort(sort.Field, sort.Ascending)
	}
	if page != nil {
		search = search.From(page.PageNum).Size(page.PageSize)
	}

	if len(fields) > 0 {
		fsc := elasticV6.NewFetchSourceContext(true)
		if isInclude {
			search.FetchSourceContext(fsc.Include(fields...))
		} else {
			search.FetchSourceContext(fsc.Exclude(fields...))
		}
	}

	var T interface{}

	T, err := search.Do(context.Background())
	return T.(*elasticV6.SearchResult), err
}

func (this *EsClientV6) Count(indexName string, query elasticV6.Query, typeName ...string) (int64, error) {
	if len(typeName) == 0 {
		return 0, errors.New("Type 不能为空!")
	}
	return this.Client.Count(indexName).Type(typeName[0]).Query(query).Do(context.Background())
}

func (this *EsClientV6) DeleteByQuery(indexName string, query elasticV6.Query, typeName ...string) (interface{}, error) {
	if len(typeName) == 0 {
		return nil, errors.New("Type 不能为空!")
	}
	return this.Client.DeleteByQuery(indexName).Type(typeName[0]).Query(query).Pretty(true).Do(context.Background())
}

func (this *EsClientV6) UpdateByID(indexName string, id string, query elasticV6.Query, typeName ...string) (interface{}, error) {
	if len(typeName) == 0 {
		return nil, errors.New("Type 不能为空!")
	}
	return this.Client.Update().Index(indexName).Type(typeName[0]).Id(id).Doc(query).Do(context.Background())
}
