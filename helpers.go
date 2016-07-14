package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"path/filepath"
	"strings"
)

func pathToObject(path string) string {
	pathbase := strings.TrimPrefix(path, params.source)
	return filepath.Join(params.prefix, pathbase)
}

func injectSize(name string) string {
	ext := filepath.Ext(name)
	origname := strings.TrimSuffix(name, ext)
	return fmt.Sprintf("%s_%d.%s", origname, params.size, ext)
}

func getImagingFormat(format string) imaging.Format {

	formats := map[string]imaging.Format{
		"jpg":  imaging.JPEG,
		"jpeg": imaging.JPEG,
		"png":  imaging.PNG,
		"tif":  imaging.TIFF,
		"tiff": imaging.TIFF,
		"bmp":  imaging.BMP,
		"gif":  imaging.GIF,
	}

	return formats[strings.ToLower(format)]
}

func verbose(format string, args ...interface{}) {
	if params.verbose {
		fmt.Printf(format, args...)
	}
}
