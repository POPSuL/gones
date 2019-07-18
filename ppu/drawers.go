package ppu

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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

type FyneDrawer struct {
	root   fyne.Window
	raster *canvas.Raster
}

type SDLDrawer struct {
	renderer sdl.Renderer
	window   *sdl.Window
	surface  *sdl.Surface
	frame    int64
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

func NewFyneDrawer() *FyneDrawer {
	//parentDrawer := NewPngDrawer()
	app := app.New()
	root := app.NewWindow("Gones")
	r := canvas.NewRaster(func(w, h int) image.Image {

		println("sdfsdf")
		return image.NewRGBA(image.Rect(0, 0, 256, 224))
	})

	root.SetContent(r)

	go root.ShowAndRun()

	return &FyneDrawer{
		root,
		r,
	}
}

func (D *FyneDrawer) Draw(buffer []uint8) {

	//println("saaaaaa")
	//D.root.Canvas().Refresh(nil)
}

func NewSDLDrawer() *SDLDrawer {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	ttf.Init()
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		256, 224, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	//surface.FillRect(nil, 0)
	rect := sdl.Rect{0, 0, 256, 244}
	surface.FillRect(&rect, 0xff000000)
	surface.SetClipRect(&rect)
	window.UpdateSurface()

	rdr, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	return &SDLDrawer{
		*rdr,
		window,
		surface,
		0,
	}
}

func (D *SDLDrawer) Draw(buffer []byte) {
	for y := 0; y < 224; y++ {
		for x := 0; x < 256; x++ {
			index := uint((x + (y * 0x100)) * 4)
			D.surface.Set(x, y, color.RGBA{
				R: buffer[index],
				G: buffer[index+1],
				B: buffer[index+2],
				A: 0xff,
			})
		}
	}

	D.frame++

	//println(D.frame)

	D.window.UpdateSurface()

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		if event.GetType() == sdl.QUIT {
			os.Exit(0)
		}
	}
}
