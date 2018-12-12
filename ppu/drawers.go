package ppu

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

const dst = "./screen"

type Drawer interface {
	Draw(buffer []uint8)
}

type PngDrawer struct {
	frame uint
}

func NewPngDrawer() *PngDrawer {
	return &PngDrawer{
		0,
	}
}

func (D *PngDrawer) Draw(buffer []uint8) {
	img := image.NewRGBA(image.Rect(0, 0, 256, 224))
	for y := 0; y < 224; y++ {
		for x := 0; x < 256; x++ {
			index := uint((x + (y * 0x100)) * 4)
			img.Set(x, y, color.RGBA{
				R: buffer[index],
				G: buffer[index+1],
				B: buffer[index+2],
				A: 0xff,
			})
		}
	}

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		_ = os.Mkdir(dst, os.ModeDir)
	}

	D.frame++
	file, err := os.Create(fmt.Sprintf("%s/img%06d.png", dst, D.frame))
	if err != nil {
		panic(err)
	}

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}
