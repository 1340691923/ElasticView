package es_optimize

import (
	elasticV6 "github.com/olivere/elastic"
	elasticV7 "github.com/olivere/elastic/v7"
)

type OptimizeInterface interface {
	SetIndexName(indexName string)
	DoV6(client *elasticV6.Client) (err error)
	DoV7(client *elasticV7.Client) (err error)
}

// 索引操作工厂
var optimizeMap = map[string]OptimizeInterface{
	"_refresh":     newRefresh(),
	"_cache/clear": newCacheClear(),
	"_flush":       newFlush(),
	"_forcemerge":  newForcemerge(),
	"open":         newOpen(),
	"close":        newClose(),
}

func OptimizeFactory(command string) OptimizeInterface {
	if _, ok := optimizeMap[command]; ok {
		return optimizeMap[command]
	} else {
		return nil
	}
}
