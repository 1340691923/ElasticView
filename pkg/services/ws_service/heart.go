package ws_service

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
)

type HeartController struct {
}

func NewHeartController() *HeartController {
	return &HeartController{}
}

func (this *HeartController) Ping(req *dto.C2S_PING, ctx *Ctx) *vo.S2C_PONG {
	return &vo.S2C_PONG{}
}
