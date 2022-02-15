package es_optimize

import "github.com/olivere/elastic"

type OptimizeInterface interface {
	SetIndexName(indexName string)
	Do(client *elastic.Client) (err error)
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
