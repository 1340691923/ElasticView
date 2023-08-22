package navicat_service

type Search struct {
	Query       interface{}              `json:"query"`
	From        int                      `json:"from"`
	Size        int                      `json:"size"`
	Sort        []map[string]interface{} `json:"sort"`
	SearchAfter *[]interface{}           `json:"search_after,omitempty"` // explains how the score was computed

}
