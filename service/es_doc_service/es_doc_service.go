package es_doc_service

import (
	"context"
	"encoding/json"
	"github.com/1340691923/ElasticView/es_sdk/pkg"
	proto2 "github.com/1340691923/ElasticView/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/pkg/escache"
)

type EsDocService struct {
	esClient pkg.EsI
}

func NewEsDocService(esClient pkg.EsI) *EsDocService {
	return &EsDocService{esClient: esClient}
}

func (this *EsDocService) DeleteRowByIDAction(ctx context.Context, esDocDeleteRowByID *escache.EsDocDeleteRowByID) (err error) {
	res, err := this.esClient.Delete(ctx, proto2.DeleteRequest{
		Index:        esDocDeleteRowByID.IndexName,
		DocumentType: esDocDeleteRowByID.Type,
		DocumentID:   esDocDeleteRowByID.ID,
	})
	if err != nil {
		return
	}
	if res.StatusErr() != nil {
		return res.StatusErr()
	}
	return
}

func (this *EsDocService) EsDocUpdateByID(ctx context.Context, esDocUpdateByID *escache.EsDocUpdateByID) (err error) {

	res, err := this.esClient.Update(ctx, proto2.UpdateRequest{
		Index:        esDocUpdateByID.Index,
		DocumentType: esDocUpdateByID.Type,
		DocumentID:   esDocUpdateByID.ID,
	}, esDocUpdateByID.JSON)
	if err != nil {
		return
	}
	if res.StatusErr() != nil {
		return res.StatusErr()
	}
	return
}

func (this *EsDocService) EsDocInsert(ctx context.Context, esDocUpdateByID *escache.EsDocUpdateByID) (res json.RawMessage, err error) {

	resp, err := this.esClient.Create(ctx, proto2.CreateRequest{
		Index:        esDocUpdateByID.Index,
		DocumentType: esDocUpdateByID.Type,
		DocumentID:   esDocUpdateByID.ID,
	}, esDocUpdateByID.JSON)
	if err != nil {
		return nil, err
	}
	if resp.StatusErr() != nil {
		return nil, resp.StatusErr()
	}

	return resp.JsonRawMessage(), nil
}
