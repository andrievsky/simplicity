package images

import (
	"fmt"
	"strings"
)

type Format struct {
	Name   string
	Ext    string
	Width  int
	Height int
}

const JpegExt = "jpeg"
const PngExt = "png"

func (f *Format) FileName() string {
	return fmt.Sprintf("%s.%s", f.Name, f.Ext)
}

func (f *Format) SizeDefined() bool {
	return f.Width > 0
}

var Source = &Format{"source", "data", 0, 0}
var Canonical = &Format{"canonical", PngExt, 0, 0}
var WebLange = &Format{"web-std", JpegExt, 1280, 853}           //3:2
var WebThumbSquare = &Format{"web-thumb-sq", JpegExt, 100, 100} //1:1

func resolveFormat(candidate string) (*Format, error) {
	switch candidate {
	case Source.Name:
		return Source, nil
	case Canonical.Name:
		return Canonical, nil
	case WebLange.Name:
		return WebLange, nil
	case WebThumbSquare.Name:
		return WebThumbSquare, nil
	default:
		return nil, fmt.Errorf("unknown format: %s", candidate)
	}
}

func resolveExtFromFileName(fileName string) (string, error) {
	ext := fileName[strings.LastIndex(fileName, ".")+1:]
	return resolveExt(ext)
}

func resolveExt(ext string) (string, error) {
	switch ext {
	case "jpeg", "jpg":
		return JpegExt, nil
	case "png":
		return PngExt, nil
	default:
		return "", fmt.Errorf("unknown extension: %s", ext)
	}
}

func resolveMime(ext string) string {
	switch ext {
	case JpegExt:
		return "image/jpeg"
	case PngExt:
		return "image/png"
	default:
		return "application/octet-stream"
	}
}
