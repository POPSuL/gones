package ppu

import (
	"fmt"
	"github.com/popsul/gones/bus"
	"github.com/veandco/go-sdl2/sdl"
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

type SDLDrawer struct {
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
			index := (x + (y * 0x100)) * 4
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

func NewSDLDrawer(keypad *bus.Keypad) *SDLDrawer {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width*2, height*2, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	surface.Format.Format = sdl.PIXELFORMAT_RGB888

	rect := sdl.Rect{W: width * 2, H: height * 2}
	surface.FillRect(&rect, 0xff000000)
	surface.SetClipRect(&rect)
	window.UpdateSurface()

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
		window,
		surface,
		controller,
		0,
		keypad,
	}
}

func (D *SDLDrawer) Draw(buffer []byte) {
	buff := D.surface.Pixels()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := (x + (y * 0x100)) * 4
			for xShift := 0; xShift < 2; xShift++ {
				for yShift := 0; yShift < 2; yShift++ {
					i := int32(y*2+yShift)*D.surface.Pitch + int32(x*2+xShift)*int32(D.surface.Format.BytesPerPixel)
					buff[i+2] = buffer[index+0]
					buff[i+1] = buffer[index+1]
					buff[i+0] = buffer[index+2]
				}
			}
		}
	}

	D.frame++

	err := D.window.UpdateSurface()
	if err != nil {
		panic(err)
	}

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
