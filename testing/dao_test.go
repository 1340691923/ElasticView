package testing

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/server"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/davecgh/go-spew/spew"
	"log"
	"testing"

	"time"
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

func TestInitialize1(t *testing.T) {
	svr, err := server.InitializeNoticeDao(args)
	if err != nil {
		t.Fatal(err)
	}
	_, err = svr.CreateNotice(&model.Notice{
		Title:       "test2",
		Content:     "test111",
		Type:        "success",
		Level:       "success",
		IsTask:      1,
		FromUid:     1,
		PluginAlias: "",
		Source:      "测试",
		TargetType:  "roles",
		Created:     time.Now(),
		Updated:     time.Now(),
	}, []int{2})

	log.Println("succ", err)
}

func TestInitialize2(t *testing.T) {
	svr, err := server.InitializeNoticeDao(args)
	if err != nil {
		t.Fatal(err)
	}

	res, cnt, err := svr.GetUserNoticesWithReadStatus(1, []int{1, 2}, 2, "test",
		1, 10)
	spew.Dump(res)
	log.Println("succ", res, cnt, err)
}

func TestInitialize3(t *testing.T) {
	svr, err := server.InitializeNoticeDao(args)
	if err != nil {
		t.Fatal(err)
	}
	err = svr.DeleteNotice(2)

	log.Println("succ", err)
}

func TestInitialize4(t *testing.T) {
	svr, err := server.InitializeNoticeDao(args)
	if err != nil {
		t.Fatal(err)
	}
	err = svr.MarkUserNoticeAsRead(context.Background(), 1, 1)

	log.Println("succ", err)
}
