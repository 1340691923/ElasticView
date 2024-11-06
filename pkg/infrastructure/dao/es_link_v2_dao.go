package dao

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// EsLinkModel es连接信息表
type EsLinkV2Dao struct {
	orm *sqlstore.SqlStore
}

func NewEsLinkV2Dao(orm *sqlstore.SqlStore) *EsLinkV2Dao {
	return &EsLinkV2Dao{orm: orm}
}

func (this *EsLinkV2Dao) GetEsConnectDataByRoleId(ctx context.Context, roleId []int, id int) (esConnectTmp *model.EsConnect, err error) {
	esConnectTmp = new(model.EsConnect)
	SQL := `
select ip,user,pwd,version,rootpem,certpem,keypem,header  from
 (
select * from
(select es_link_id,role_cfg_id,id as relation_id from eslink_role_cfg_reletion )
 j1
left join
(
select * from
(select id, user,pwd,rootpem ,certpem,keypem,remark as cfg_remark,header  from eslink_cfg_v2   ) ecv
left join
(select role_id,es_link_cfg_id from gm_role_eslink_cfg_v2) grec
on ecv.id = grec.es_link_cfg_id
) j2
on j1.role_cfg_id = j2.id
)
j1 
left join
(
select id,ip,remark,version from es_link_v2 elv
) j3
on  j1.es_link_id = j3.id
where role_id in (?) and relation_id = ? and remark is not null
`

	err = this.
		orm.
		Raw(SQL, roleId, id).
		Scan(&esConnectTmp).
		WithContext(ctx).
		Error

	if util.FilterMysqlNilErr(err) {
		err = errors.WithStack(err)
		return nil, err
	}
	return
}

func (this *EsLinkV2Dao) Save(ctx context.Context, tx *gorm.DB, data *model.EsLinkV2) error {
	err := tx.Save(data).WithContext(ctx).Error
	err = errors.WithStack(err)
	return err
}

func (this *EsLinkV2Dao) Delete(ctx context.Context, tx *gorm.DB, id int) error {
	return errors.WithStack(tx.Where("id = ?", id).Delete(model.EsLinkV2{}).WithContext(ctx).Error)
}
