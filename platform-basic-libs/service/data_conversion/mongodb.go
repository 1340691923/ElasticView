package data_conversion

import (
	"fmt"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"time"
	"github.com/globalsign/mgo"
)

type MongoDb struct {
	request.DataxInfoTestLinkReq
}

func NewMongoDb(data request.DataxInfoTestLinkReq) Datasource {
	return &MongoDb{
		data,
	}
}

func (this *MongoDb) Ping() error {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{fmt.Sprintf("%s:%v",this.IP,this.Port)},
		Timeout:  3 * time.Second,
		Database: this.DbName,
		Username: this.Username,
		Password: this.Pwd,
	}

	_, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return err
	}

	return nil
}
