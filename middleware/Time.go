package middleware

import (
	"time"

	"github.com/1340691923/ElasticView/engine/logs"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Timer(ctx *fiber.Ctx) error {

	// start timer
	start := time.Now()
	// next routes
	err := ctx.Next()
	// stop timer
	stop := time.Now()

	logs.Logger.Info("时间拦截器", zap.String("消耗时间：", stop.Sub(start).String()))
	return err

}
