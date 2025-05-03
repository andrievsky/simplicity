package images

import (
	"errors"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"math"
)

func Transcode(in *Format, out *Format, r io.Reader, w io.Writer) error {
	slog.Info("Transcode", "in", in, "out", out)
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
		img = resizeWithBackground(img, out.Width, out.Height, color.White)
	}
	if err = encode(out, img, w); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}
	return nil
}

func resizeImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

func resizeWithCrop(img image.Image, width, height int) image.Image {
	// original dimensions
	origBounds := img.Bounds()
	origW := origBounds.Dx()
	origH := origBounds.Dy()

	// compute scale ratios
	scaleX := float64(width) / float64(origW)
	scaleY := float64(height) / float64(origH)
	// choose the larger scale → no empty bars, but some overflow to crop
	scale := math.Max(scaleX, scaleY)

	// dimensions after scaling
	newW := int(math.Ceil(float64(origW) * scale))
	newH := int(math.Ceil(float64(origH) * scale))

	// scale the entire image
	scaled := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(scaled, scaled.Bounds(), img, origBounds, draw.Over, nil)

	// compute top-left of the crop window to center-crop
	offsetX := (newW - width) / 2
	offsetY := (newH - height) / 2
	cropRect := image.Rect(offsetX, offsetY, offsetX+width, offsetY+height)

	// extract the crop
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dst, dst.Bounds(), scaled, cropRect.Min, draw.Src)

	return dst
}

func resizeWithBackground(img image.Image, width, height int, bgColor color.Color) image.Image {
	origBounds := img.Bounds()
	origW, origH := origBounds.Dx(), origBounds.Dy()

	// compute uniform scale to fit
	scaleX := float64(width) / float64(origW)
	scaleY := float64(height) / float64(origH)
	scale := math.Min(scaleX, scaleY)

	// new dimensions after scaling
	newW := int(math.Ceil(float64(origW) * scale))
	newH := int(math.Ceil(float64(origH) * scale))

	// create destination filled with bgColor
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// compute offsets to center
	offsetX := (width - newW) / 2
	offsetY := (height - newH) / 2

	// define target rect where the image will be scaled into
	targetRect := image.Rect(offsetX, offsetY, offsetX+newW, offsetY+newH)

	// scale directly into dst
	draw.CatmullRom.Scale(dst, targetRect, img, origBounds, draw.Over, nil)

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
