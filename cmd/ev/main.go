package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/server"
	"github.com/1340691923/ElasticView/pkg/util"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var args *config.CommandLineArgs

func init() {
	args = &config.CommandLineArgs{}
	flag.StringVar(&args.HomePath, "homePath", util.GetCurrentDirectory(), "ev程序所在文件夹")
	flag.StringVar(&args.CmdName, "cmdName", "ev", "二进制名称")
	flag.StringVar(&args.ConfigFile, "configFile", "config/config.yml", "配置文件路径")
	flag.Parse()
}

func main() {
	svr, err := server.Initialize(args)
	if err != nil {
		log.Println("初始化ev失败")
		panic(err)
	}
	err = svr.Init()
	if err != nil {
		log.Println("初始化ev失败")
		panic(err)
	}
	go listenToSystemSignals(context.Background(), svr)

	if err = svr.Run(func(svr *server.Server) error {
		log.Println("服务退出成功")
		return nil
	}, func(svr *server.Server) error {
		return svr.CloseSqlx()
	}, func(svr *server.Server) error {
		return svr.CloseLog()
	}); err != nil {
		log.Println("启动EV失败")
		panic(err)
	}
}

func listenToSystemSignals(ctx context.Context, s *server.Server) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-signalChan:
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			if err := s.Shutdown(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "ev服务关闭超时\n")
			}
			return
		}
	}
}
