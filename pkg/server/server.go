package server

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm/migrator"
	"github.com/1340691923/ElasticView/pkg/registry"
	"github.com/1340691923/ElasticView/pkg/services/es_link_service"
	"github.com/1340691923/ElasticView/pkg/services/eve_service"
	"github.com/1340691923/ElasticView/pkg/services/gm_user"
	"github.com/1340691923/ElasticView/pkg/services/plugin_install_service"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"reflect"
	"runtime"
	"sync"
)

type Server struct {
	cfg                *config.Config
	orm                *orm.Gorm
	migrator           *migrator.Migrator
	logger             *logger.AppLogger
	shutdownOnce       sync.Once
	context            context.Context
	shutdownFn         context.CancelFunc
	childRoutines      *errgroup.Group
	isInitialized      bool
	mtx                sync.Mutex
	backgroundServices []registry.BackgroundService
	esLinkService      *es_link_service.EsLinkService
	eveService         *eve_service.EvEService
	gmUserService      *gm_user.GmUserService
	pluginInstaller    *plugin_install_service.PluginInstaller
}

func NewServer(
	cfg *config.Config,
	migrator *migrator.Migrator,
	logger *logger.AppLogger,
	orm *orm.Gorm,
	backgroundServiceProvider registry.BackgroundServiceRegistry,
	esLinkService *es_link_service.EsLinkService,
	eveService *eve_service.EvEService,
	gmUserService *gm_user.GmUserService,
	pluginInstaller *plugin_install_service.PluginInstaller,
) *Server {
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)
	svr := &Server{
		migrator:           migrator,
		cfg:                cfg,
		logger:             logger,
		shutdownFn:         shutdownFn,
		context:            childCtx,
		childRoutines:      childRoutines,
		orm:                orm,
		backgroundServices: backgroundServiceProvider.GetServices(),
		esLinkService:      esLinkService,
		eveService:         eveService,
		gmUserService:      gmUserService,
		pluginInstaller:    pluginInstaller,
	}
	return svr
}

func (this *Server) Init() (err error) {
	return nil
}

func (this *Server) RunMigrator() (err error) {

	err = this.migrator.Start()

	if err != nil {
		err = errors.WithMessage(err, "migrator初始化失败")
		return
	}

	this.logger.Info("migrator初始化成功")
	return
}

func (this *Server) Run(exitFn ...func(svr *Server) error) (err error) {
	if !this.cfg.DeBug && runtime.GOOS == "windows" {
		this.runSystray()
	}
	services := this.backgroundServices

	this.eveService.FlushAccessToken(this.context)

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
			err = service.Run(this.context)

			if err != nil && !errors.Is(err, context.Canceled) {
				this.logger.Sugar().Errorf("停止后台服务(%s)异常:%+v", serviceName, err)
				return errors.WithStack(fmt.Errorf("%s run error: %w", serviceName, err))
			}

			this.logger.Info("后台服务已停止", zap.String("服务名", serviceName), zap.Error(err))
			return nil
		})
	}

	this.logger.Info("等待后台服务启动中...")
	err = this.childRoutines.Wait()
	if err != nil {
		return errors.Wrap(err, "后台服务停止异常")
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
	this.logger.Info("退出进程")
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

func (this *Server) ResetAdminPwd(pwd string) error {
	if pwd == "" {
		return errors.New("密码不能为空")
	}
	return this.gmUserService.UpdatePassById(context.Background(), 1, pwd)
}
