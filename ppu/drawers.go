package ppu

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"github.com/popsul/gones/bus"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"image"
	"image/color"
	"image/png"
	"os"
)

const dst = "./screen"
const width = 256
const height = 224

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
	renderer   sdl.Renderer
	window     *sdl.Window
	surface    *sdl.Surface
	controller *sdl.GameController
	frame      int64
	keypad     *bus.Keypad
}

func NewPngDrawer() *PngDrawer {
	return &PngDrawer{
		0,
	}
}

func (D *PngDrawer) Draw(buffer []uint8) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
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

func NewSDLDrawer(keypad *bus.Keypad) *SDLDrawer {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	ttf.Init()
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width*2, height*2, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	//surface.FillRect(nil, 0)
	rect := sdl.Rect{W: width * 2, H: height * 2}
	surface.FillRect(&rect, 0xff000000)
	surface.SetClipRect(&rect)
	window.UpdateSurface()

	rdr, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	var controller *sdl.GameController = nil
	for i := 0; i < sdl.NumJoysticks(); i++ {
		if sdl.IsGameController(i) {
			fmt.Printf("Found controller %d\n", i)
			controller = sdl.GameControllerOpen(i)
			if controller != nil {
				fmt.Printf("Opened %s\n", controller.Name())
			}
		}
	}

	return &SDLDrawer{
		*rdr,
		window,
		surface,
		controller,
		0,
		keypad,
	}
}

func (D *SDLDrawer) Draw(buffer []byte) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := uint((x + (y * 0x100)) * 4)
			D.surface.Set(x*2, y*2, color.RGBA{
				R: buffer[index],
				G: buffer[index+1],
				B: buffer[index+2],
				A: 0xff,
			})
			D.surface.Set(x*2+1, y*2, color.RGBA{
				R: buffer[index],
				G: buffer[index+1],
				B: buffer[index+2],
				A: 0xff,
			})
			D.surface.Set(x*2, y*2+1, color.RGBA{
				R: buffer[index],
				G: buffer[index+1],
				B: buffer[index+2],
				A: 0xff,
			})
			D.surface.Set(x*2+1, y*2+1, color.RGBA{
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
		switch event.(type) {
		case *sdl.KeyboardEvent:
			ev := event.(*sdl.KeyboardEvent)
			if ev.Type == sdl.KEYDOWN {
				//fmt.Printf("keydown 0x%02x\n", D.matchKey(ev.Keysym.Sym))
				D.keypad.KeyDown(D.matchKey(ev.Keysym.Sym))
			} else if ev.Type == sdl.KEYUP {
				//println("keyup")
				D.keypad.KeyUp(D.matchKey(ev.Keysym.Sym))
			} else {
				println(ev.Type)
				println(ev.Keysym.Scancode)
			}
		case *sdl.ControllerButtonEvent:
			ev := event.(*sdl.ControllerButtonEvent)
			if ev.Type == sdl.CONTROLLERBUTTONDOWN {
				//fmt.Printf("keydown 0x%02x\n", D.matchControllerButton(ev.Button))
				D.keypad.KeyDown(D.matchControllerButton(ev.Button))
			} else if ev.Type == sdl.CONTROLLERBUTTONUP {
				//println("keyup")
				D.keypad.KeyUp(D.matchControllerButton(ev.Button))
			} else {
				println(ev.Type)
			}
			//fmt.Println(event.GetType())
		}
	}
}

func (D *SDLDrawer) matchKey(key sdl.Keycode) uint {
	//Maps a keyboard key to a nes key.
	// A, B, SELECT, START, ↑, ↓, ←, →
	switch key {
	case sdl.K_SEMICOLON:
		return 0
	case sdl.K_COMMA:
		return 1
	case sdl.K_TAB:
		return 2
	case sdl.K_SPACE:
		return 3
	case sdl.K_w:
		return 4
	case sdl.K_s:
		return 5
	case sdl.K_a:
		return 6
	case sdl.K_d:
		return 7
	}
	return 8
}

func (D *SDLDrawer) matchControllerButton(key uint8) uint {
	switch key {
	case sdl.CONTROLLER_BUTTON_A:
		return 0
	case sdl.CONTROLLER_BUTTON_B:
		return 1
	case sdl.CONTROLLER_BUTTON_BACK:
		return 2
	case sdl.CONTROLLER_BUTTON_START:
		return 3
	case sdl.CONTROLLER_BUTTON_DPAD_UP:
		return 4
	case sdl.CONTROLLER_BUTTON_DPAD_DOWN:
		return 5
	case sdl.CONTROLLER_BUTTON_DPAD_LEFT:
		return 6
	case sdl.CONTROLLER_BUTTON_DPAD_RIGHT:
		return 7

	}
	return 8
}
