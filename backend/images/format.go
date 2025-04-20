package images

import "fmt"

type AspectRatio struct {
	Width  int
	Height int
}

type Format struct {
	Name        string
	Ext         string
	Width       int
	AspectRatio AspectRatio
}

func (f Format) FileName() string {
	return fmt.Sprintf("%s.%s", f.Name, f.Ext)
}

var Source = &Format{"source", "data", 0, AspectRatio{}}
var WebLange = &Format{"web-std", "jpeg", 1280, AspectRatio{3, 2}}
var WebThumbSquare = &Format{"web-thumb-sq", "jpeg", 400, AspectRatio{1, 1}}

func resolveFormat(candidate string) (*Format, error) {
	switch candidate {
	case Source.Name:
		return Source, nil
	case WebLange.Name:
		return WebLange, nil
	case WebThumbSquare.Name:
		return WebThumbSquare, nil
	default:
		return nil, fmt.Errorf("unknown format: %s", candidate)
	}
}

func resolveExt(ext string) (string, error) {
	switch ext {
	case "jpeg", "jpg":
		return "jpeg", nil
	case "png":
		return "png", nil
	default:
		return "", fmt.Errorf("unknown extension: %s", ext)
	}
}

func resolveMime(ext string) string {
	switch ext {
	case "jpeg", "jpg":
		return "image/jpeg"
	case "png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}
