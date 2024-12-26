package main

//go:generate go install github.com/1340691923/ElasticView/cmd/ev_builder/

//go:generate ev_builder.exe

//go:generate swag init -g cmd/ev/main.go -o resources/docs -exclude resources,logs,config

//go:generate gowatch
