package main

import (
	"log"
	"net/url"
)

func main() {
	p := "_cat/nodes?h=ip,name,heap.percent,heap.current,heap.max,ram.percent,ram.current,ram.max,node.role,master,cpu,load_1m,load_5m,load_15m,disk.used_percent,disk.used,disk.total"

	u,err:=url.Parse(p)

	if err!=nil{
		panic(err)
	}

	log.Println(u.Query())
}
