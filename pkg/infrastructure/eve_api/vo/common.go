package vo

import (
	"errors"
)

type ApiCommonRes struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (this *ApiCommonRes) Error() error {
	if this.Code != 0 {
		return errors.New(this.Msg)
	}
	return nil
}

type OpenApiPlugin struct {
	ID              int    `db:"id" json:"id"`
	PluginAlias     string `db:"plugin_alias" json:"plugin_alias"`
	PluginName      string `db:"plugin_name" json:"plugin_name"`
	Realname        string `db:"realname" json:"realname"`
	Describe        string `db:"p_desc" json:"describe"`
	Readme          string `db:"readme" json:"readme"`
	CreateTime      string `db:"create_time" json:"create_time"`
	UpdateTime      string `db:"update_time" json:"update_time"`
	Logo            string `db:"logo" json:"logo"`
	DownloadCnt     int    `db:"download_cnt" json:"download_cnt"`
	StarCnt         int    `db:"star_cnt" json:"star_cnt"`
	PublishTime     string `db:"publish_time" json:"publish_time"`
	DownloadUserCnt int    `db:"download_user_cnt" json:"download_user_cnt"`
	StarState       *int   `db:"star_state" json:"star_state"`
	HasDownload     bool   `json:"has_download" `
}

type PluginListRes struct {
	List  []OpenApiPlugin `json:"list"`
	Count int             `json:"count"`
}

type WxArticleModel struct {
	ID         int    `json:"id" db:"id"`
	Title      string `json:"title" db:"title"`
	Link       string `json:"link" db:"link"`
	CreateTime string `json:"create_time" db:"create_time"`
	Typ        string `json:"typ" db:"typ"` // info/warning/success/danger/primary
	TagName    string `json:"tag_name" db:"tag_name"`
}

type WxArticleList struct {
	Data []WxArticleModel
}

type PublishData struct {
	PluginName string `json:"plugin_name"`
	PluginPublish
}

type GetLocalPlugin struct {
	PluginID      string `json:"plugin_id"`
	PluginName    string `json:"plugin_name"`
	Version       string `json:"version"`
	HasUpdate     bool   `json:"has_update"`
	UpdateVersion string `json:"update_version"`
}

type PublishRes struct {
	List  []PublishData `json:"list"`
	Count int           `json:"count"`
}

type PluginPublish struct {
	ID                 int64  `json:"id" db:"id"`
	PluginId           int64  `json:"plugin_id" db:"plugin_id"`
	Version            string `json:"version" db:"version"`
	SourceCodeUrl      string `json:"source_code_url" db:"source_code_url"`
	Changelog          string `json:"changelog" db:"changelog"`
	CdnSourceCodeURL   string `json:"cdn_source_code_url" db:"cdn_source_code_url"`
	CreateTime         string `json:"create_time" db:"create_time"`
	UpdateTime         string `json:"update_time" db:"update_time"`
	CreateBy           int64  `json:"create_by" db:"create_by"`
	GteEvDependencyVer string `json:"gte_ev_dependency_ver" db:"gte_ev_dependency_ver"`
	LteEvDependencyVer string `json:"lte_ev_dependency_ver" db:"lte_ev_dependency_ver"`
	LinuxAmd64Crc      string `json:"linux_amd64_crc" db:"linux_amd64_crc"`
	WindowsCrc         string `json:"windows_crc" db:"windows_crc"`
	LinuxArm64Crc      string `json:"linux_arm64_crc" db:"linux_arm64_crc"`
	DarwinCrc          string `json:"darwin_crc" db:"darwin_crc"`
	State              int8   `json:"state" db:"state"`
	Msg                string `json:"msg" db:"msg"`
	DarwinURL          string `json:"darwin_url" db:"darwin_url"`
	WindowsURL         string `json:"windows_url" db:"windows_url"`
	LinuxAmd64URL      string `json:"linux_amd64_url" db:"linux_amd64_url"`
	LinuxArm64URL      string `json:"linux_arm64_url" db:"linux_arm64_url"`
	DownloadCount      int64  `json:"download_count" db:"download_count"`
}

type GetPluginDownloadUrlRes struct {
	DownloadUrl string `json:"downloadUrl"`
	DonwloadCrc string `json:"downloadCrc"`
}

type GetEvMaxVersionRes struct {
	DownloadUrl string `json:"downloadUrl"`
	Version     string `json:"version"`
}
