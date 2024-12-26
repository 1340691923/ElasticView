package testing

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/server"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/davecgh/go-spew/spew"
	"log"
	"testing"
)

var args *config.CommandLineArgs

func init() {
	args = &config.CommandLineArgs{
		ConfigFile: "D:\\eve\\ev2\\config_dev\\config.yml",
		CmdName:    "ev",
		HomePath:   util.GetCurrentDirectory(),
	}
}

func TestPluginInstallService(t *testing.T) {
	svr, err := server.InitializeProvideInstaller(args)
	if err != nil {
		t.Fatal(err)
	}
	err = svr.Remove(context.Background(), "ev-tools", "0.0.3")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("succ")
}

func TestEvBackDaoGetCheckEvKeyStatus(t *testing.T) {
	dao, err := server.InitializeEvApiDao(args)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(dao.GetWxArticleList(context.Background()))
}

func TestInitializeGmRoleEslinkCfgV2Dao(t *testing.T) {
	svr, err := server.InitializeGmRoleEslinkCfgV2Dao(args)
	if err != nil {
		t.Fatal(err)
	}
	l, err := svr.QueryByRoleID(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(l)
	log.Println("succ")
}
