package vo

type LivePrivateData struct {
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}
