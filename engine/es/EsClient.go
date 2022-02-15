package es

import elasticV6 "github.com/olivere/elastic"

type EsClient interface {
	Ping() (interface{}, int, error)                                     //Ping
	CreateIndex(indexName string, body interface{}) (interface{}, error) //创建索引
	CatIndices() (interface{}, error)                                    //查看索引信息
	CatAliases() (interface{}, error)                                    //显示别名,过滤器,路由信息
	CatAllocation() (interface{}, error)                                 //显示每个节点分片数量、占用空间
	CatCount() (interface{}, error)                                      //显示索引文档的数量
	CatHealth() (interface{}, error)                                     //查看集群健康状况
	CatShards() (interface{}, error)                                     //显示索引分片信息
	DeleteIndex(indexNameList []string) (interface{}, error)             //删除索引
	CloseIndex(indexNameList []string) error                             //关闭索引
	OpenIndex(indexNameList []string) error                              //打开索引
	FreezeIndex(indexNameList []string) error                            //冻结索引
	UnfreezeIndex(indexNameList []string) error                          //解冻索引
	CreateMapping(indexName string, body Json) (interface{}, error)      //mapping 新增字段
	GetMapping(indexName string, body Json, typeName ...string) (interface{}, error)
	PutData(indexName string, body Json, typeName ...string) (interface{}, error)
	DeleteById(indexName, id string, typeName ...string) (interface{}, error)
	Search(indexName string, query elasticV6.Query, sort *Sort, page *Page, fields []string, isInclude bool, typeName ...string) (*elasticV6.SearchResult, error)
	Count(indexName string, query elasticV6.Query, typeName ...string) (int64, error)
	Refresh(indexName ...string) (interface{}, error)
	Flush(indexName ...string) (interface{}, error)
	Fsync(indexName ...string) (interface{}, error)
	Rollover(alias, indexName string) (interface{}, error)
	IndexStats(indexName []string, metrics []string) (interface{}, error)
	Alias(alias, indexName string) (interface{}, error)
	IndexSegments(indexName string) (interface{}, error)
	UpdateByID(indexName string, id string, query elasticV6.Query, typeName ...string) (interface{}, error)
	DeleteByQuery(indexName string, query elasticV6.Query, typeName ...string) (interface{}, error)
	TasksList() (interface{}, error)
	Reindex(sourceIndex, destinationIndex string) (interface{}, error)
	GetMappings() (interface{}, error)
	IndexPutSettings(indexName string, body Json) (interface{}, error)
}
