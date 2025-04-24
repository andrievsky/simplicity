package images

import (
	"errors"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

func Transcode(in *Format, out *Format, r io.Reader, w io.Writer) error {
	if in == nil || out == nil {
		return errors.New("format cannot be nil")
	}

	if in == out {
		return errors.New("input and output formats are the same")
	}

	if out == Source {
		return errors.New("output format cannot be source")
	}

	img, err := decode(in, r)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}
	if out.SizeDefined() {
		img = resizeImage(img, out.Width, out.Height)
	}
	if err := encode(out, img, w); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}
	return nil
}

func resizeImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

func decode(format *Format, r io.Reader) (image.Image, error) {
	ext, err := resolveExt(format.Ext)
	if err != nil {
		return nil, err
	}

	switch ext {
	case JpegExt:
		return decodeJpeg(r)
	case PngExt:
		return decodePng(r)
	default:
		return nil, errors.New("unsupported image format")
	}
}

func decodeJpeg(r io.Reader) (image.Image, error) {
	img, err := jpeg.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func decodePng(r io.Reader) (image.Image, error) {
	img, err := png.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func encode(format *Format, img image.Image, w io.Writer) error {
	ext, err := resolveExt(format.Ext)
	if err != nil {
		return err
	}

	switch ext {
	case JpegExt:
		return encodeJpeg(img, w)
	case PngExt:
		return encodePng(img, w)
	default:
		return errors.New("unsupported image format")
	}
}

func encodeJpeg(img image.Image, w io.Writer) error {
	opts := jpeg.Options{
		Quality: 85,
	}
	return jpeg.Encode(w, img, &opts)
}

func encodePng(img image.Image, w io.Writer) error {
	return png.Encode(w, img)
}
