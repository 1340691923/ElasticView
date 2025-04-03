package dto

type LiveBroadcast struct {
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type BatchLiveBroadcast struct {
	Channel string        `json:"channel"`
	List    []interface{} `json:"list"`
}
