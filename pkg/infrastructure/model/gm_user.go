package model

import (
	"github.com/1340691923/ElasticView/pkg/util"
	"time"
)

// GmUserModel BI用户
type GmUserModel struct {
	Id       int    `gorm:"column:id;primary_key;NOT NULL" json:"id"`
	Username string `gorm:"column:username;default:NULL;index:gm_user_username,unique" json:"username"`
	Password string `gorm:"column:password;default:NULL" json:"password"`

	Avatar string `gorm:"column:avatar;default:NULL" json:"avatar"`

	Realname      string `gorm:"column:realname;default:" json:"realname"`
	Email         string `gorm:"column:email;default:" json:"email"`
	WorkWechatUid string `gorm:"column:work_wechat_uid;default:" json:"work_wechat_uid"`
	DingtalkId    string `gorm:"column:dingtalk_id;default:" json:"dingtalk_id"`
	FeishuOpenId  string `gorm:"column:feishu_open_id;default:" json:"feishu_open_id"`

	UpdateTime    time.Time `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"`
	CreateTime    time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"`
	LastLoginTime time.Time `gorm:"column:last_login_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"last_login_time"`

	IsBan int32 `gorm:"column:is_ban;default:0;NOT NULL" json:"is_ban"`
}

func (g *GmUserModel) TableName() string {
	return "gm_user"
}

// 密码进行md5混淆
func (this *GmUserModel) GetPassword() string {
	return util.MD5HexHash(util.Str2bytes(this.Password))
}
