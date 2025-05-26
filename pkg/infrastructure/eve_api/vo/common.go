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
