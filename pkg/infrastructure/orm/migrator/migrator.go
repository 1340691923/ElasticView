package migrator

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/access_control"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/migrator_cfg"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/eve-plugin-sdk-go/util"
	"github.com/go-gormigrate/gormigrate/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sort"
)

type Migrator struct {
	orm  *gorm.DB
	cfg  *config.Config
	log  *logger.AppLogger
	rbac *access_control.Rbac
}

func NewMigrator(orm *orm.Gorm, cfg *config.Config, log *logger.AppLogger, rbac *access_control.Rbac) *Migrator {
	migrator_cfg.RbacInstance = rbac
	return &Migrator{orm: orm.DB, cfg: cfg, log: log, rbac: rbac}
}

func (this *Migrator) Start() error {

	storeMaxVersion := "0.0.0"
	currentVersion := config.GetVersion()
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

	m := gormigrate.New(this.orm, gormigrate.DefaultOptions, migrator_cfg.Migrators)

	isRollback := util.LessThan(currentVersion, storeMaxVersion)

	type IDs struct {
		ID string
	}

	var ids []IDs
	for _, v := range migrator_cfg.Migrators {
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
	for _, v := range migrator_cfg.Migrators {
		versions = append(versions, v.ID)
	}
	return versions
}
