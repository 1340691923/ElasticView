package main

import (
	"github.com/1340691923/ElasticView/engine/db"
	"log"
)

func main() {
	log.Println(db.SqlBuilder.Select("a.test,b.test2").From("a").LeftJoin("b on a.id = b.id").ToSql())
}
