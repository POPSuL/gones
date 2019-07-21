package ppu

import (
	"github.com/popsul/gones/bus"
	. "github.com/popsul/gones/common"
)

var COLORS = [64][3]uint8{
	{0x80, 0x80, 0x80}, {0x00, 0x3D, 0xA6}, {0x00, 0x12, 0xB0}, {0x44, 0x00, 0x96},
	{0xA1, 0x00, 0x5E}, {0xC7, 0x00, 0x28}, {0xBA, 0x06, 0x00}, {0x8C, 0x17, 0x00},
	{0x5C, 0x2F, 0x00}, {0x10, 0x45, 0x00}, {0x05, 0x4A, 0x00}, {0x00, 0x47, 0x2E},
	{0x00, 0x41, 0x66}, {0x00, 0x00, 0x00}, {0x05, 0x05, 0x05}, {0x05, 0x05, 0x05},
	{0xC7, 0xC7, 0xC7}, {0x00, 0x77, 0xFF}, {0x21, 0x55, 0xFF}, {0x82, 0x37, 0xFA},
	{0xEB, 0x2F, 0xB5}, {0xFF, 0x29, 0x50}, {0xFF, 0x22, 0x00}, {0xD6, 0x32, 0x00},
	{0xC4, 0x62, 0x00}, {0x35, 0x80, 0x00}, {0x05, 0x8F, 0x00}, {0x00, 0x8A, 0x55},
	{0x00, 0x99, 0xCC}, {0x21, 0x21, 0x21}, {0x09, 0x09, 0x09}, {0x09, 0x09, 0x09},
	{0xFF, 0xFF, 0xFF}, {0x0F, 0xD7, 0xFF}, {0x69, 0xA2, 0xFF}, {0xD4, 0x80, 0xFF},
	{0xFF, 0x45, 0xF3}, {0xFF, 0x61, 0x8B}, {0xFF, 0x88, 0x33}, {0xFF, 0x9C, 0x12},
	{0xFA, 0xBC, 0x20}, {0x9F, 0xE3, 0x0E}, {0x2B, 0xF0, 0x35}, {0x0C, 0xF0, 0xA4},
	{0x05, 0xFB, 0xFF}, {0x5E, 0x5E, 0x5E}, {0x0D, 0x0D, 0x0D}, {0x0D, 0x0D, 0x0D},
	{0xFF, 0xFF, 0xFF}, {0xA6, 0xFC, 0xFF}, {0xB3, 0xEC, 0xFF}, {0xDA, 0xAB, 0xEB},
	{0xFF, 0xA8, 0xF9}, {0xFF, 0xAB, 0xB3}, {0xFF, 0xD2, 0xB0}, {0xFF, 0xEF, 0xA6},
	{0xFF, 0xF7, 0x9C}, {0xD7, 0xE8, 0x95}, {0xA6, 0xED, 0xAF}, {0xA2, 0xF2, 0xDA},
	{0x99, 0xFF, 0xFC}, {0xDD, 0xDD, 0xDD}, {0x11, 0x11, 0x11}, {0x11, 0x11, 0x11},
}

type Renderer struct {
	frameBuffer []uint8
	background  []Tile
	serial      uint
	drawer      Drawer
}

func NewRenderer(keypad *bus.Keypad) *Renderer {
	R := new(Renderer)
	R.drawer = NewSDLDrawer(keypad)
	R.serial = 0
	R.frameBuffer = make([]uint8, 256*256*4)
	return R
}

func (R *Renderer) shouldPixelHide(x uint, y uint) bool {
	// TODO: WTF??
	//return false
	tileX := ^^(x / 8)
	tileY := ^^(y / 8)
	backgroundIndex := tileY*33 + tileX

	//println("xx", x, y)

	sprite := uint(len(R.background)) >= backgroundIndex && R.background[backgroundIndex].Pattern != nil
	if sprite {
		return true
	}

	return false
	// NOTE: If background pixel is not transparent, we need to hide sprite.
	//return !((sprite[y % 8] && sprite[y % 8][x % 8] % 4) == 0)
}

func (R *Renderer) Render(data *RenderingData) {
	//return
	if data.background != nil && len(data.background) > 0 {
		R.renderBackground(data.background, data.palette)
	}

	if data.sprites != nil && len(data.sprites) > 0 {
		R.renderSprites(data.sprites, data.palette)
	}

	R.drawer.Draw(R.frameBuffer)
}

func (R *Renderer) renderBackground(background []Tile, palette []byte) {
	R.background = background
	for i := uint(0); i < uint(len(background)); i++ {
		x := (i % 33) * 8
		y := ^^(i / 33) * 8
		R.renderTile(background[i], x, y, palette)
	}
}

func (R *Renderer) renderSprites(sprites []SpriteWithAttribute, palette []byte) {

	for _, sprite := range sprites {
		R.renderSprite(sprite, palette)
	}
}

func (R *Renderer) renderTile(tile Tile, tileX uint, tileY uint, palette []byte) {
	//{ sprite, paletteId, scrollX, scrollY }: Tile
	offsetX := tile.scrollX % 8
	offsetY := tile.scrollY % 8
	for i := uint(0); i < 8; i++ {
		for j := uint(0); j < 8; j++ {
			paletteIndex := tile.paletteId*4 + tile.Pattern[i][j]
			colorId := palette[paletteIndex]
			color := COLORS[colorId]
			x := tileX + j - offsetX
			y := tileY + i - offsetY
			if x >= 0 && 0xFF >= x && y >= 0 && y < 224 {
				index := (x + (y * 0x100)) * 4
				R.frameBuffer[index] = color[0]
				R.frameBuffer[index+1] = color[1]
				R.frameBuffer[index+2] = color[2]
				R.frameBuffer[index+3] = 0xFF
			}
		}
	}
}

func (R *Renderer) renderSprite(sprite SpriteWithAttribute, palette []byte) {
	isVerticalReverse := I2b(sprite.attribute & 0x80)
	isHorizontalReverse := I2b(sprite.attribute & 0x40)
	isLowPriority := I2b(sprite.attribute & 0x20)
	paletteId := sprite.attribute & 0x03
	for i := uint(0); i < 8; i++ {
		for j := uint(0); j < 8; j++ {
			x := sprite.x + B2ix(isHorizontalReverse, 7-j, j)
			y := sprite.y + B2ix(isVerticalReverse, 7-i, i)
			if isLowPriority && R.shouldPixelHide(x, y) {
				continue
			}
			if sprite.sprite != nil && sprite.sprite[i][j] > 0 {
				colorId := palette[paletteId*4+sprite.sprite[i][j]+0x10]
				color := COLORS[colorId]
				index := (x + y*0x100) * 4
				R.frameBuffer[index] = color[0]
				R.frameBuffer[index+1] = color[1]
				R.frameBuffer[index+2] = color[2]
				//R.frameBuffer[index+3] = 0xff
				// data[index + 3] = 0xFF;
			}
		}
	}
}
