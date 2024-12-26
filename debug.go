package main

//go:generate go build -o build.exe  cmd/build_ev/main.go
//go:generate build.exe
//go:generate swag init -g cmd/ev/main.go -o docs -exclude resources,logs,config
//go:generate gowatch
