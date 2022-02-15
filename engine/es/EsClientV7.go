package es

import (
	"context"

	elasticV6 "github.com/olivere/elastic"
	elasticV7 "github.com/olivere/elastic/v7"
)

type EsClientV7 struct {
	Client          *elasticV7.Client
	esConnectConfig EsConnect
}

func NewEsClientV7(esConnectConfig *EsConnect) (esClient *elasticV7.Client, err error) {

	optList := []elasticV7.ClientOptionFunc{
		elasticV7.SetSniff(false),
		elasticV7.SetHealthcheck(false),
	}

	optList = append(optList, elasticV7.SetURL(esConnectConfig.Ip))

	if esConnectConfig.User != "" || esConnectConfig.Pwd != "" {
		optList = append(optList, elasticV7.SetBasicAuth(esConnectConfig.User, esConnectConfig.Pwd))
	}

	esClient, err = elasticV7.NewSimpleClient(optList...)

	if err != nil {
		return
	}

	return
}

func (this *EsClientV7) Ping() (interface{}, int, error) {
	return this.Client.Ping(this.esConnectConfig.Ip).Do(context.Background())
}

func (this *EsClientV7) CreateIndex(indexName string, body interface{}) (interface{}, error) {
	return this.Client.CreateIndex(indexName).BodyJson(body).Do(context.Background())
}

func (this *EsClientV7) CatIndices() (interface{}, error) {
	return this.Client.CatIndices().Pretty(true).Do(context.Background())
}

func (this *EsClientV7) CatAliases() (interface{}, error) {
	return this.Client.CatAliases().Pretty(true).Do(context.Background())
}

func (this *EsClientV7) CatAllocation() (interface{}, error) {
	return this.Client.CatAllocation().Pretty(true).Do(context.Background())
}

func (this *EsClientV7) CatCount() (interface{}, error) {
	return this.Client.CatCount().Pretty(true).Do(context.Background())
}

func (this *EsClientV7) CatHealth() (interface{}, error) {
	return this.Client.CatHealth().Pretty(true).Do(context.Background())
}

func (this *EsClientV7) CatShards() (interface{}, error) {
	return this.Client.CatShards().Pretty(true).Do(context.Background())
}

func (this *EsClientV7) DeleteIndex(indexNameList []string) (interface{}, error) {
	return this.Client.DeleteIndex(indexNameList...).Do(context.Background())
}

func (this *EsClientV7) CloseIndex(indexNameList []string) error {
	var err error
	for _, indexName := range indexNameList {
		_, err = this.Client.CloseIndex(indexName).Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *EsClientV7) OpenIndex(indexNameList []string) error {
	var err error
	for _, indexName := range indexNameList {
		_, err = this.Client.OpenIndex(indexName).Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *EsClientV7) FreezeIndex(indexNameList []string) error {
	var err error
	for _, indexName := range indexNameList {
		_, err = this.Client.FreezeIndex(indexName).Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *EsClientV7) UnfreezeIndex(indexNameList []string) error {
	var err error
	for _, indexName := range indexNameList {
		_, err = this.Client.UnfreezeIndex(indexName).Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *EsClientV7) CreateMapping(indexName string, body Json) (interface{}, error) {
	return this.Client.PutMapping().
		Index(indexName).
		BodyJson(body).
		Do(context.Background())
}

func (this *EsClientV7) IndexPutSettings(indexName string, body Json) (interface{}, error) {
	return this.Client.IndexPutSettings().Index(indexName).BodyJson(body).Do(context.TODO())
}

func (this *EsClientV7) Reindex(sourceIndex, destinationIndex string) (interface{}, error) {
	return this.Client.Reindex().SourceIndex(sourceIndex).DestinationIndex(destinationIndex).Do(context.Background())
}

func (this *EsClientV7) TasksList() (interface{}, error) {
	return this.Client.TasksList().Pretty(true).
		Human(true).
		//Header("X-Opaque-Id", "123456").
		Do(context.TODO())
}

func (this *EsClientV7) Refresh(indexName ...string) (interface{}, error) {
	return this.Client.Refresh(indexName...).Do(context.TODO())
}

func (this *EsClientV7) Flush(indexName ...string) (interface{}, error) {
	return this.Client.Flush(indexName...).Do(context.TODO())
}

func (this *EsClientV7) Fsync(indexName ...string) (interface{}, error) {
	return this.Client.SyncedFlush(indexName...).Do(context.TODO())
}

func (this *EsClientV7) Rollover(alias, indexName string) (interface{}, error) {
	return this.Client.RolloverIndex(alias).NewIndex(indexName).Do(context.TODO())
}

func (this *EsClientV7) IndexSegments(indexName string) (interface{}, error) {
	return this.Client.IndexSegments(indexName).Pretty(true).Do(context.TODO())
}

func (this *EsClientV7) Alias(alias, indexName string) (interface{}, error) {
	return this.Client.Alias().
		Action(elasticV7.NewAliasAddAction(alias).Index(indexName).IsWriteIndex(false)).
		Do(context.TODO())
}

func (this *EsClientV7) IndexStats(indexName []string, metrics []string) (interface{}, error) {
	return this.Client.IndexStats().Index(indexName...).Metric(metrics...).Do(context.Background())
}

func (this *EsClientV7) GetMapping(indexName string, body Json, typeName ...string) (interface{}, error) {
	return this.Client.GetMapping().Index(indexName).Type("_doc").Do(context.Background())
}

func (this *EsClientV7) PutData(indexName string, body Json, typeName ...string) (interface{}, error) {
	return this.Client.Index().Index(indexName).BodyJson(body).Do(context.Background())
}

func (this *EsClientV7) DeleteById(indexName, id string, typeName ...string) (interface{}, error) {
	return this.Client.Delete().Index(indexName).Id(id).Do(context.Background())
}

func (this *EsClientV7) Search(indexName string, query elasticV6.Query, sort *Sort, page *Page, fields []string, isInclude bool, typeName ...string) (*elasticV6.SearchResult, error) {

	search := this.Client.Search(indexName).Query(query)
	if sort != nil {
		search = search.Sort(sort.Field, sort.Ascending)
	}
	if page != nil {
		search = search.From(page.PageNum).Size(page.PageSize)
	}

	if len(fields) > 0 {
		fsc := elasticV7.NewFetchSourceContext(true)
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

func (this *EsClientV7) Count(indexName string, query elasticV6.Query, typeName ...string) (int64, error) {
	return this.Client.Count(indexName).Query(query).Do(context.Background())
}

func (this *EsClientV7) DeleteByQuery(indexName string, query elasticV6.Query, typeName ...string) (interface{}, error) {
	return this.Client.DeleteByQuery(indexName).Query(query).Pretty(true).Do(context.Background())
}

func (this *EsClientV7) UpdateByID(indexName string, id string, query elasticV6.Query, typeName ...string) (interface{}, error) {
	return this.Client.Update().Index(indexName).Id(id).Doc(query).Do(context.Background())
}

func (this *EsClientV7) GetMappings() (interface{}, error) {
	return this.Client.GetMapping().Do(context.Background())
}
