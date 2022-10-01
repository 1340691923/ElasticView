package data_conversion

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/request"
	"github.com/globalsign/mgo"
	"time"
)

type MongoDb struct {
	request.DataxInfoTestLinkReq
}

func (this *MongoDb) GetTableColumns(tableName string) (interface{}, error) {
	conn, err := this.getConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	db := conn.DB(this.DbName)
	res, err := db.C(tableName).Indexes()
	return res, err
}

func (this *MongoDb) GetTables() ([]string, error) {
	conn, err := this.getConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	db := conn.DB(this.DbName)
	names, err := db.CollectionNames()
	if err != nil {
		return nil, err
	}
	return names, err
}

func NewMongoDb(data request.DataxInfoTestLinkReq) Datasource {
	return &MongoDb{
		data,
	}
}

func (this *MongoDb) getConn() (*mgo.Session, error) {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{fmt.Sprintf("%s:%v", this.IP, this.Port)},
		Timeout:  3 * time.Second,
		Database: this.DbName,
		Username: this.Username,
		Password: this.Pwd,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)

	return session, err
}

func (this *MongoDb) Ping() error {
	session, err := this.getConn()

	if err != nil {
		return err
	}

	defer session.Close()

	return nil
}
