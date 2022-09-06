package views

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var StatisFs embed.FS

func GetFileSystem() http.FileSystem {
	fsys, err := fs.Sub(StatisFs, "dist")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}
