package es

import (
	"github.com/1340691923/ElasticView/engine/es"
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/platform-basic-libs/my_error"
	"github.com/1340691923/ElasticView/platform-basic-libs/service/es_settings"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

type EsBackAbstruct struct {
}

type SnapshotRepositoryRes struct {
	Name                   string `json:"name"`
	Type                   string `json:"type"`
	Location               string `json:"location"`
	Compress               string `json:"compress"`
	MaxRestoreBytesPerSec  string `json:"max_restore_bytes_per_sec"`
	MaxSnapshotBytesPerSec string `json:"max_snapshot_bytes_per_sec"`
	ChunkSize              string `json:"chunk_size"`
	Readonly               string `json:"readonly"`
}

func (this *EsBackAbstruct) SnapshotRepositoryList(ctx *fiber.Ctx) (ctxRes map[string]interface{}, err error) {
	esSnapshotInfo := es.EsSnapshotInfo{}
	err = ctx.BodyParser(&esSnapshotInfo)
	if err != nil {
		return
	}
	esClinet, err := es.GetEsClientV6ByID(esSnapshotInfo.EsConnect)
	if err != nil {
		return
	}

	clusterSettings, err := es_settings.NewSettings(esClinet.(*es.EsClientV6).Client)
	if err != nil {
		return
	}
	pathRepo := clusterSettings.GetPathRepo()

	if len(pathRepo) == 0 {
		err = my_error.NewError(`path.repo没有设置`, 199999)
		return
	}

	res, err := esClinet.(*es.EsClientV6).Client.SnapshotGetRepository(esSnapshotInfo.SnapshotInfoList...).Do(ctx.Context())
	if err != nil {
		return
	}

	list := []SnapshotRepositoryRes{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for name, settings := range res {
		var t SnapshotRepositoryRes
		t.Type = settings.Type
		t.Name = name
		b, err := json.Marshal(settings.Settings)
		if err != nil {
			logs.Logger.Sugar().Errorf("err", err)
			continue
		}
		err = json.Unmarshal(b, &t)
		if err != nil {
			logs.Logger.Sugar().Errorf("err", err)
			continue
		}
		list = append(list, t)
	}

	ctxRes = map[string]interface{}{
		"list":     list,
		"res":      res,
		"pathRepo": pathRepo,
	}

	return
}
