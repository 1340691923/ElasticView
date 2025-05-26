package api

type API = string

const (
	GetEvAccessToken       API = "GetEvAccessToken"
	GetPluginList          API = "GetPluginList"
	GetPluginInfo          API = "GetPluginInfo"
	GetPluginDownloadUrl   API = "GetPluginDownloadUrl"
	GetEvPluginMaxVersion  API = "GetEvPluginMaxVersion"
	GetEvMaxVersion        API = "GetEvMaxVersion"
	GetEvPluginsMaxVersion API = "GetEvPluginsMaxVersion"
	StarPlugin             API = "StarPlugin"
	GetWxArticleList       API = "GetWxArticleList"
	//
	AddComment   API = "AddComment"
	LikeComment  API = "LikeComment"
	ListComments API = "ListComments"
)
