package server

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/registry"
	"github.com/1340691923/ElasticView/pkg/services/es_link_service"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"reflect"
	"sync"
)

type Server struct {
	cfg                *config.Config
	sqlx               *sqlstore.SqlStore
	logger             *logger.AppLogger
	shutdownOnce       sync.Once
	context            context.Context
	shutdownFn         context.CancelFunc
	childRoutines      *errgroup.Group
	isInitialized      bool
	mtx                sync.Mutex
	backgroundServices []registry.BackgroundService
	esLinkService      *es_link_service.EsLinkService
}

func NewServer(
	cfg *config.Config,
	logger *logger.AppLogger,
	sqlx *sqlstore.SqlStore,
	backgroundServiceProvider registry.BackgroundServiceRegistry,
	esLinkService *es_link_service.EsLinkService,
) *Server {
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)
	svr := &Server{
		cfg:                cfg,
		logger:             logger,
		shutdownFn:         shutdownFn,
		context:            childCtx,
		childRoutines:      childRoutines,
		sqlx:               sqlx,
		backgroundServices: backgroundServiceProvider.GetServices(),
		esLinkService:      esLinkService,
	}
	return svr
}

func (this *Server) Init() (err error) {
	log.Println("init")
	return this.esLinkService.FlushEsLinkList()
	//return nil
}

func (this *Server) Run(exitFn ...func(svr *Server) error) (err error) {

	services := this.backgroundServices

	for _, svc := range services {
		if registry.IsDisabled(svc) {
			continue
		}

		service := svc
		serviceName := reflect.TypeOf(service).String()
		this.childRoutines.Go(func() error {
			select {
			case <-this.context.Done():
				return this.context.Err()
			default:
			}
			this.logger.Info("开启后台服务", zap.String("服务名", serviceName))
			err := service.Run(this.context)

			if err != nil && !errors.Is(err, context.Canceled) {
				this.logger.Error("停止后台服务异常", zap.String("服务名", serviceName), zap.Error(err))
				return fmt.Errorf("%s run error: %w", serviceName, err)
			}

			this.logger.Info("后台服务已停止", zap.String("服务名", serviceName), zap.Error(err))
			return nil
		})
	}

	this.logger.Info("等待后台服务启动中...")
	err = this.childRoutines.Wait()
	if err != nil {
		return errors.Wrap(err, "服务启动异常")
	}
	for _, fn := range exitFn {
		err = fn(this)
		if err != nil {
			return err
		}
	}

	return err
}

func (this *Server) Shutdown(ctx context.Context) (err error) {
	this.shutdownOnce.Do(func() {
		this.logger.Info("开始停止进程")
		this.shutdownFn()
		select {

		case <-ctx.Done():
			this.logger.Warn("关闭服务超时")
			err = fmt.Errorf("关闭服务超时")
		}

	})
	return
}

func (this *Server) GetLogger() *logger.AppLogger {
	return this.logger
}

func (this *Server) CloseLog() error {
	this.logger.Sync()
	return nil
}

func (this *Server) CloseSqlx() error {
	err := this.sqlx.Close()
	if err != nil {
		this.logger.Error("sqlStore 连接关闭失败", zap.Error(err))
		return err
	}
	this.logger.Info("sqlStore 连接关闭成功")
	return nil
}
