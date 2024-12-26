package views

import "embed"

//go:embed dist/*.html
var IndexFileTemplate embed.FS

//go:embed dist/css/*
var CssFs embed.FS

//go:embed dist/js/*
var JsFs embed.FS

//go:embed dist/img/*
var ImgFs embed.FS

//go:embed dist/favicon.ico
var FaviconFs embed.FS

func GetFavicon() ([]byte, error) {
	return FaviconFs.ReadFile("dist/favicon.ico")
}
