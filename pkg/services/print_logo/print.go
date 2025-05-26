package print_logo

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"go.uber.org/zap"
)

type PrintLogo struct {
	log *zap.Logger
	cfg *config.Config
}

func ProvidePrintLogo(log *logger.AppLogger, cfg *config.Config) (*PrintLogo, error) {
	log = log.Named("PrintLogo")
	return &PrintLogo{
		log: logger.ZapLog2AppLog(log),
		cfg: cfg,
	}, nil
}

func (this *PrintLogo) Run(ctx context.Context) (err error) {
	this.printStartLogo()
	<-ctx.Done()
	this.printByeLogo()
	return nil
}

func (this *PrintLogo) printStartLogo() {
	fmt.Println(`

___________.__                   __  .__         ____   ____.__               
\_   _____/|  | _____    _______/  |_|__| ____   \   \ /   /|__| ______  _  __
 |    __)_ |  | \__  \  /  ___/\   __\  |/ ___\   \   Y   / |  |/ __ \ \/ \/ /
 |        \|  |__/ __ \_\___ \  |  | |  \  \___    \     /  |  \  ___/\     / 
/_________/|____(______/______> |__| |__|\_____>    \___/   |__|\_____>\/\_/  
                                                                              


               您的 Ev 插件小管家，已为您就位! 
	`)
}

func (this *PrintLogo) printByeLogo() {
	fmt.Println(`

				 __.                 
				 | |__ ___.__. ____  
				 | __ <   |  |/ __ \ 
				 | \_\ \___  \  ___/ 
				 |___  / ____|\____>
					 

               Bye,您的 Ev 插件小管家期待下次为您服务! 
	`)
}
