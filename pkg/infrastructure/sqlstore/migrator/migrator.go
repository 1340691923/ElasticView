package migrator

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/access_control"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	util2 "github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/eve-plugin-sdk-go/util"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"time"
)

var rbacInstance *access_control.Rbac

var migrators = []*gormigrate.Migration{
	{
		ID: "0.0.1",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(
				&model.EsLinkV2{},
				&model.GmUserModel{},
				&model.EslinkRoleCfgReletion{},
				&model.EslinkCfgV2{},
				&model.GmOperaterLog{},
				&model.GmRole{},
				&model.GmRoleEslinkCfgV2{},
				&model.GmGuid{},
			)
			if err != nil {
				return err
			}

			if err := tx.Exec("INSERT INTO gm_user (id, username, password, realname) VALUES(1, 'admin', '21232f297a57a5a743894a0e4a801fc3',  '肖文龙');").Error; err != nil {
				return err
			}
			if err := tx.Exec(`INSERT INTO gm_role (id,role_name,description,role_list) VALUES (1,'admin','超级管理员','[{"path":"/permission","component":"layout","redirect":"/permission/role","alwaysShow":true,"meta":{"title":"权限","icon":"system"},"children":[{"path":"role","name":"role","component":"views/permission/role","meta":{"title":"权限组管理","icon":"role"},"children":[]},{"path":"user","name":"user","component":"views/permission/user","meta":{"title":"用户管理","icon":"el-icon-user"},"children":[]},{"path":"operater_log","name":"operater_log","component":"views/permission/operater_log","meta":{"title":"操作日志列表","icon":"dict"},"children":[]}]},{"path":"/connect-tree","component":"layout","redirect":"/connect-tree/link","alwaysShow":false,"meta":{"title":"数据源管理","icon":"el-icon-link"},"children":[{"path":"link","name":"link","component":"views/connect-tree/link","meta":{"title":"数据源管理","icon":"el-icon-link"},"children":[]},{"path":"auth","name":"auth","component":"views/connect-tree/auth","meta":{"title":"鉴权管理","icon":"el-icon-user"},"children":[]}]},{"path":"/plugins","component":"layout","redirect":"/plugins/market","alwaysShow":false,"meta":{"title":"插件","icon":"el-icon-link"},"children":[{"path":"/plugins/market","name":"market","component":"views/plugins/market","meta":{"title":"插件市场","icon":"el-icon-link"},"children":[]}]}]');`).Error; err != nil {
				return err
			}

			now := time.Now().Format(time.DateTime)

			if err := tx.Exec("INSERT INTO es_link_v2 (id,ip,created,updated,remark,version,create_by) VALUES (1,'http://127.0.0.1:9200',?,?,'默认连接','elasticsearch6.x',1);", now, now).Error; err != nil {
				return err
			}

			if err := tx.Exec("INSERT INTO eslink_cfg_v2(id, `user`, pwd, rootpem, certpem, keypem, created, updated, create_by, remark) VALUES(1, '', '', '', '', '', ?, ?, 1, '空鉴权');", now, now).Error; err != nil {
				return err
			}

			if err := tx.Exec("INSERT INTO eslink_role_cfg_reletion(id, es_link_id, role_cfg_id, created, updated) VALUES(1, 1, 1, ?, ?);", now, now).Error; err != nil {
				return err
			}

			if err := tx.Exec("INSERT INTO gm_role_eslink_cfg_v2 (id,role_id,es_link_cfg_id,es_link_id,created,updated) VALUES (1,1,1,1,?,?);", now, now).Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}, //0.0.1
	{
		ID: "0.0.4",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(&model.GmUserModel{}, &model.UserRoleRelationModel{})
			if err != nil {
				return errors.WithStack(err)
			}

			err = tx.Create(&model.UserRoleRelationModel{
				UserId:     1,
				RoleId:     1,
				UpdateTime: time.Now(),
				CreateTime: time.Now(),
			}).Error
			if err != nil {
				return errors.WithStack(err)
			}

			baseRole := &model.GmRole{
				RoleName:    "基础权限组",
				Description: "可以进行数据源设置操作",
				RoleList:    util2.StringPtr(`[{"path":"/connect-tree","component":"layout","redirect":"/connect-tree/link","alwaysShow":false,"meta":{"title":"数据源管理","icon":"el-icon-link"},"children":[{"path":"link","name":"link","component":"views/connect-tree/link","meta":{"title":"数据源管理","icon":"el-icon-link"},"children":[]},{"path":"auth","name":"auth","component":"views/connect-tree/auth","meta":{"title":"鉴权管理","icon":"el-icon-user"},"children":[]}]}]`),
			}

			err = tx.Create(baseRole).Error
			if err != nil {
				return errors.WithStack(err)
			}

			rbacInstance.RemoveFilteredPolicy(0, strconv.Itoa(baseRole.Id)) //先全清掉
			eg := errgroup.Group{}

			apis := []string{"/api/es_link/InsertAction", "/api/es_link/DeleteAction", "/api/es_link/UpdateAction", "/api/es_link/InsertEsCfgAction", "/api/es_link/UpdateEsCfgAction", "/api/es_link/DeleteEsCfgAction"}

			for _, api := range apis {
				api := api
				roleId := baseRole.Id
				eg.Go(func() error {
					_, err = rbacInstance.AddPolicy(strconv.Itoa(roleId), api, "*")
					if err != nil {
						return errors.WithStack(err)
					}
					return nil
				})
			}

			err = eg.Wait()
			if err != nil {
				return errors.WithStack(err)
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}, //0.0.1
}

type Migrator struct {
	orm  *gorm.DB
	cfg  *config.Config
	log  *logger.AppLogger
	rbac *access_control.Rbac
}

func NewMigrator(orm *sqlstore.SqlStore, cfg *config.Config, log *logger.AppLogger, rbac *access_control.Rbac) *Migrator {
	rbacInstance = rbac
	return &Migrator{orm: orm.DB, cfg: cfg, log: log, rbac: rbac}
}

func (this *Migrator) Start() error {

	storeMaxVersion := "0.0.0"
	currentVersion := config.Version
	type VersionInfo struct {
		Id           int `gorm:"primaryKey"`
		LocalVersion string
	}
	var versionInfo VersionInfo

	hasVersionInfo := this.orm.Migrator().HasTable(&versionInfo)

	if hasVersionInfo {
		err := this.orm.First(&versionInfo).Error
		if err != nil {
			this.log.Error("err", zap.Error(err))
		} else {
			storeMaxVersion = versionInfo.LocalVersion
		}
	}

	if !hasVersionInfo {
		err := this.orm.AutoMigrate(&versionInfo)
		if err != nil {
			this.log.Error("err", zap.Error(err))
		}
	} else {
		err := this.orm.First(&versionInfo).Error
		if err != nil {
			this.log.Error("err", zap.Error(err))
		} else {
			storeMaxVersion = versionInfo.LocalVersion
		}
	}

	if hasVersionInfo && (storeMaxVersion == currentVersion) {
		return nil
	}

	defer func() {
		versionInfo.Id = 1
		versionInfo.LocalVersion = currentVersion
		this.orm.Save(&versionInfo)
	}()

	m := gormigrate.New(this.orm, gormigrate.DefaultOptions, migrators)

	isRollback := util.LessThan(currentVersion, storeMaxVersion)

	type IDs struct {
		ID string
	}

	var ids []IDs
	for _, v := range migrators {
		ids = append(ids, IDs{ID: v.ID})
	}

	if isRollback {
		sort.Slice(ids, func(i, j int) bool {
			return util.LessThan(ids[j].ID, ids[i].ID)
		})

		for _, v := range ids {
			if util.LessThan(storeMaxVersion, v.ID) {
				continue
			}
			if util.LessThan(v.ID, currentVersion) {
				break
			}

			err := m.RollbackTo(v.ID)
			if err != nil {
				this.log.Error("err", zap.Error(err))
			}
		}
		return nil
	}
	sort.Slice(ids, func(i, j int) bool {
		return util.LessThan(ids[i].ID, ids[j].ID)
	})

	for _, v := range ids {
		if util.LessThan(v.ID, storeMaxVersion) {
			continue
		}
		if util.LessThan(currentVersion, v.ID) {
			break
		}

		err := m.MigrateTo(v.ID)
		if err != nil {
			this.log.Error("err", zap.Error(err))
		}
	}

	return nil
}

func (this *Migrator) GetMigratorsVersions() []string {
	versions := []string{}
	for _, v := range migrators {
		versions = append(versions, v.ID)
	}
	return versions
}
