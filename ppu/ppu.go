package ppu

import (
	"github.com/popsul/gones/bus"
	. "github.com/popsul/gones/common"
	"github.com/popsul/gones/interrupts"
)

const SPRITES_NUMBER = 0x100

type Ppu struct {
	// PPU power up state
	// see. https://wiki.nesdev.com/w/index.php/PPU_power_up_state
	//
	// Memory map
	/*
		| addr           |  description               |
		+----------------+----------------------------+
		| 0x0000-0x0FFF  |  Pattern table#0           |
		| 0x1000-0x1FFF  |  Pattern table#1           |
		| 0x2000-0x23BF  |  Name table                |
		| 0x23C0-0x23FF  |  Attribute table           |
		| 0x2400-0x27BF  |  Name table                |
		| 0x27C0-0x27FF  |  Attribute table           |
		| 0x2800-0x2BBF  |  Name table                |
		| 0x2BC0-0x2BFF  |  Attribute table           |
		| 0x2C00-0x2FBF  |  Name Table                |
		| 0x2FC0-0x2FFF  |  Attribute Table           |
		| 0x3000-0x3EFF  |  mirror of 0x2000-0x2EFF   |
		| 0x3F00-0x3F0F  |  background Palette        |
		| 0x3F10-0x3F1F  |  sprite Palette            |
		| 0x3F20-0x3FFF  |  mirror of 0x3F00-0x3F1F   |
	*/
	/*
		  Control Register1 0x2000
		| bit  | description                                 |
		+------+---------------------------------------------+
		|  7   | Assert NMI when VBlank 0: disable, 1:enable |
		|  6   | PPU master/slave, always 1                  |
		|  5   | Sprite size 0: 8x8, 1: 8x16                 |
		|  4   | Bg Pattern table 0:0x0000, 1:0x1000         |
		|  3   | sprite Pattern table 0:0x0000, 1:0x1000     |
		|  2   | PPU memory increment 0: +=1, 1:+=32         |
		|  1-0 | Name table 0x00: 0x2000                     |
		|      |            0x01: 0x2400                     |
		|      |            0x02: 0x2800                     |
		|      |            0x03: 0x2C00                     |
	*/
	/*
		  Control Register2 0x2001
		| bit  | description                                 |
		+------+---------------------------------------------+
		|  7-5 | Background color  0x00: Black               |
		|      |                   0x01: Green               |
		|      |                   0x02: Blue                |
		|      |                   0x04: Red                 |
		|  4   | Enable sprite                               |
		|  3   | Enable background                           |
		|  2   | Sprite mask       render left end           |
		|  1   | Background mask   render left end           |
		|  0   | Display type      0: color, 1: mono         |
	*/
	/** @var int[] */
	registers []byte
	/** @var int */
	cycle uint
	/** @var int */
	line uint
	/** @var bool */
	isValidVramAddr bool
	/** @var bool */
	isLowerVramAddr bool
	/** @var int */
	spriteRamAddr uint
	/** @var int */
	vramAddr uint
	/** @var \Nes\Bus\Ram */
	vram bus.Ram
	/** @var int */
	vramReadBuf byte
	/** @var \Nes\Bus\Ram */
	spriteRam bus.Ram
	/** @var \Nes\Bus\PpuBus */
	bus *bus.PpuBus
	/** @var \Nes\Ppu\Tile[] */
	background []Tile
	/** @var \Nes\Ppu\SpriteWithAttribute[] */
	sprites []SpriteWithAttribute
	/** @var \Nes\Ppu\Palette */
	palette Palette
	/** @var \Nes\Cpu\Interrupts */
	interrupts *interrupts.Interrupts
	/** @var bool */
	isHorizontalScroll bool
	/** @var int */
	scrollX uint
	/** @var int */
	scrollY uint
	/** @var bool */
	isHorizontalMirror bool
}

type RenderingData struct {
	palette    []byte
	background []Tile
	sprites    []SpriteWithAttribute
}

func NewPpu(ppuBus *bus.PpuBus, interrupts *interrupts.Interrupts, isHorizontalMirror bool) *Ppu {
	ppu := new(Ppu)
	ppu.registers = make([]byte, 8)
	ppu.cycle = 0
	ppu.line = 0
	ppu.isValidVramAddr = false
	ppu.isLowerVramAddr = false
	ppu.isHorizontalScroll = true
	ppu.vramAddr = 0x0000
	ppu.vram = *bus.NewRam(0x2000)
	ppu.vramReadBuf = 0
	ppu.spriteRam = *bus.NewRam(0x100)
	ppu.spriteRamAddr = 0
	ppu.background = []Tile{}
	ppu.sprites = make([]SpriteWithAttribute, SPRITES_NUMBER)
	ppu.bus = ppuBus
	ppu.interrupts = interrupts
	ppu.isHorizontalMirror = isHorizontalMirror
	ppu.scrollX = 0
	ppu.scrollY = 0
	ppu.palette = *NewPalette()

	return ppu
}

func NewRenderingData(palette []byte, background []Tile, sprites []SpriteWithAttribute) *RenderingData {
	return &RenderingData{
		palette,
		background,
		sprites,
	}
}

func (P *Ppu) ReadCharacterRAM(addr uint) byte {
	return P.bus.ReadByPpu(addr)
}
func (P *Ppu) WriteCharacterRAM(addr uint, data byte) {
	P.bus.WriteByPpu(addr, data)
}

func (P *Ppu) vramOffset() uint {
	if (P.registers[0x00] & 0x04) > 0 {
		return 32
	}
	return 1
}

func (P *Ppu) ReadVram() byte {
	buf := P.vramReadBuf
	if P.vramAddr >= 0x2000 {
		addr := P.calcVramAddr()
		P.vramAddr += P.vramOffset()
		if addr >= 0x3F00 {
			return P.vram.Read(addr)
		}
		P.vramReadBuf = P.vram.Read(addr)
	} else {
		P.vramReadBuf = P.ReadCharacterRAM(P.vramAddr)
		P.vramAddr += P.vramOffset()
	}
	return buf
}

func (P *Ppu) Read(addr uint) byte {
	/*
		| bit  | description                                 |
		+------+---------------------------------------------+
		| 7    | 1: VBlank clear by reading this register    |
		| 6    | 1: sprite hit                               |
		| 5    | 0: less than 8, 1: 9 or more                |
		| 4-0  | invalid                                     |
		|      | bit4 VRAM write flag [0: success, 1: fail]  |
	*/
	if addr == 0x0002 {
		P.isHorizontalScroll = true
		data := P.registers[0x02]
		P.clearVblank()
		// P.clearSpriteHit();
		return byte(data)
	}
	// Write OAM data here. Writes will increment OAMADDR after the write
	// reads during vertical or forced blanking return the value from OAM at that address but do not increment.
	if addr == 0x0004 {
		return P.spriteRam.Read(P.spriteRamAddr)
	}
	if addr == 0x0007 {
		return P.ReadVram()
	}
	return 0
}

func (P *Ppu) Write(addr uint, data byte) {
	if addr == 0x0003 {
		P.writeSpriteRamAddr(data)
	}
	if addr == 0x0004 {
		P.writeSpriteRamData(data)
	}
	if addr == 0x0005 {
		P.writeScrollData(data)
	}
	if addr == 0x0006 {
		P.writeVramAddr(data)
	}
	if addr == 0x0007 {
		P.writeVramData(data)
	}
	P.registers[addr] = data
}

func (P *Ppu) writeSpriteRamAddr(data byte) {
	P.spriteRamAddr = uint(data)
}

func (P *Ppu) writeSpriteRamData(data byte) {
	P.spriteRam.Write(P.spriteRamAddr, data)
	P.spriteRamAddr += 1
}

func (P *Ppu) writeScrollData(data byte) {
	if P.isHorizontalScroll {
		P.isHorizontalScroll = false
		P.scrollX = uint(data) & 0xFF
	} else {
		P.isHorizontalScroll = true
		P.scrollY = uint(data) & 0xFF
	}
}

func (P *Ppu) writeVramAddr(data byte) {
	if P.isLowerVramAddr {
		P.vramAddr += uint(data)
		P.isLowerVramAddr = false
		P.isValidVramAddr = true
	} else {
		P.vramAddr = uint(data) << 8
		P.isLowerVramAddr = true
		P.isValidVramAddr = false
	}
}

func (P *Ppu) calcVramAddr() uint {
	if P.vramAddr >= 0x3000 && P.vramAddr < 0x3f00 {
		P.vramAddr -= 0x3000
		return P.vramAddr
	} else {
		return P.vramAddr - 0x2000
	}
}

func (P *Ppu) writeVramData(data byte) {
	if P.vramAddr >= 0x2000 {
		if P.vramAddr >= 0x3f00 && P.vramAddr < 0x4000 {
			P.palette.Write(P.vramAddr-0x3f00, data)
		} else {
			P.writeVram(P.calcVramAddr(), data)
		}
	} else {
		P.WriteCharacterRAM(P.vramAddr, data)
	}
	P.vramAddr += P.vramOffset()
}

func (P *Ppu) writeVram(addr uint, data byte) {
	P.vram.Write(addr, data)
}

func (P *Ppu) nameTableId() uint {
	return uint(P.registers[0x00] & 0x03)
}

func (P *Ppu) getPalette() []byte {
	return P.palette.Read()
}

func (P *Ppu) clearSpriteHit() {
	P.registers[0x02] &= 0xbf
}

func (P *Ppu) setSpriteHit() {
	P.registers[0x02] |= 0x40
}

func (P *Ppu) hasSpriteHit() bool {
	y := uint(P.spriteRam.Read(0))
	return (y == P.line) && P.isBackgroundEnable() && P.isSpriteEnable()
}

func (P *Ppu) hasVblankIrqEnabled() bool {
	return P.registers[0]&0x80 > 0
}

func (P *Ppu) isBackgroundEnable() bool {
	return P.registers[0x01]&0x08 > 0
}

func (P *Ppu) isSpriteEnable() bool {
	return P.registers[0x01]&0x10 > 0
}

func (P *Ppu) scrollTileX() uint {
	/*
	  Name table id and address
	  +------------+------------+
	  |            |            |
	  |  0(0x2000) |  1(0x2400) |
	  |            |            |
	  +------------+------------+
	  |            |            |
	  |  2(0x2800) |  3(0x2C00) |
	  |            |            |
	  +------------+------------+
	*/
	return (P.scrollX + ((P.nameTableId() % 2) * 256)) / 8
}

func (P *Ppu) scrollTileY() uint {
	return (P.scrollY + ((P.nameTableId() / 2) * 240)) / 8
}

func (P *Ppu) tileY() uint {
	return (P.line / 8) + P.scrollTileY()
}

func (P *Ppu) backgroundTableOffset() uint {
	if (P.registers[0] & 0x10) > 0 {
		return 0x1000
	}
	return 0x0000
}

func (P *Ppu) setVblank() {
	P.registers[0x02] |= 0x80
}

func (P *Ppu) isVblank() bool {
	return (P.registers[0x02] & 0x80) > 0
}

func (P *Ppu) clearVblank() {
	P.registers[0x02] &= 0x7F
}

func (P *Ppu) getBlockId(tileX uint, tileY uint) uint {
	return ^^((tileX % 4) / 2) + (^^((tileY % 4) / 2))*2
}

func (P *Ppu) getAttribute(tileX uint, tileY uint, offset uint) uint {
	addr := ^^(tileX / 4) + (^^(tileY / 4) * 8) + 0x03C0 + offset
	return uint(P.vram.Read(P.mirrorDownSpriteAddr(addr)))
}

func (P *Ppu) getSpriteId(tileX uint, tileY uint, offset uint) uint {
	tileNumber := tileY*32 + tileX
	spriteAddr := P.mirrorDownSpriteAddr(tileNumber + offset)
	return uint(P.vram.Read(spriteAddr))
}

func (P *Ppu) mirrorDownSpriteAddr(addr uint) uint {
	if !P.isHorizontalMirror {
		return addr
	}

	if addr >= 0x0400 && addr < 0x0800 || addr >= 0x0C00 {
		return addr - 0x400
	}
	return addr
}

func (P *Ppu) buildSprites() {
	offset := I2ix(uint(P.registers[0])&0x08, 0x1000, 0x0000)
	for i := uint(0); i < SPRITES_NUMBER; i = i + 4 {
		// INFO: Offset sprite Y position, because First and last 8line is not rendered.
		y := uint(P.spriteRam.Read(i) - 8)
		// TODO: WTF
		//if (y < 0) {
		//	return
		//}
		spriteId := uint(P.spriteRam.Read(i + 1))
		attr := uint(P.spriteRam.Read(i + 2))
		x := uint(P.spriteRam.Read(i + 3))
		sprite := P.buildSprite(spriteId, offset)
		P.sprites[i/4] = *NewSpriteWithAttribute(sprite, x, y, attr, spriteId)
	}
}

func (P *Ppu) buildSprite(spriteId uint, offset uint) [][]uint {
	sprite := make([][]uint, 8)
	for index := range sprite {
		sprite[index] = make([]uint, 8)
	}

	for i := uint(0); i < 16; i++ {
		for j := uint(0); j < 8; j++ {
			addr := spriteId*16 + i + offset
			ram := P.ReadCharacterRAM(addr)
			if ram&(0x80>>j) > 0 {
				sprite[i%8][j] += 0x01 << ^^(i / 8)
			}
		}
	}
	return sprite
}

func (P *Ppu) buildTile(tileX uint, tileY uint, offset uint) Tile {
	// INFO see. http://hp.vector.co.jp/authors/VA042397/nes/ppu.html
	blockId := P.getBlockId(tileX, tileY)
	spriteId := P.getSpriteId(tileX, tileY, offset)
	attr := P.getAttribute(tileX, tileY, offset)
	paletteId := (attr >> (blockId * 2)) & 0x03
	sprite := P.buildSprite(spriteId, P.backgroundTableOffset())
	return *NewTile(
		sprite,
		paletteId,
		P.scrollX,
		P.scrollY,
	)
}

func (P *Ppu) buildBackground() {
	// INFO: Horizontal offsets range from 0 to 255. "Normal" vertical offsets range from 0 to 239,
	// while values of 240 to 255 are treated as -16 through -1 in a way, but tile data is incorrectly
	// fetched from the attribute table.
	clampedTileY := P.tileY() % 30
	tableIdOffset := I2ix((^^(P.tileY() / 30))%2, 2, 0)
	// background of a line.
	// Build viewport + 1 tile for background scroll.
	for x := uint(0); x < 32+1; x++ {
		tileX := x + P.scrollTileX()
		clampedTileX := tileX % 32
		nameTableId := ((^^(tileX / 32)) % 2) + tableIdOffset
		offsetAddrByNameTable := nameTableId * 0x400
		tile := P.buildTile(clampedTileX, clampedTileY, offsetAddrByNameTable)
		P.background = append(P.background, tile)
	}
}

func (P *Ppu) TransferSprite(index uint, data byte) {
	// The DMA transfer will begin at the current OAM write address.
	// It is common practice to initialize it to 0 with a write to PPU 0x2003 before the DMA transfer.
	// Different starting addresses can be used for a simple OAM cycling technique
	// to alleviate sprite priority conflicts by flickering. If using this technique
	// after the DMA OAMADDR should be set to 0 before the end of vblank to prevent potential OAM corruption
	// (See: Errata).
	// However, due to OAMADDR writes also having a "corruption" effect[5] this technique is not recommended.
	addr := index + P.spriteRamAddr
	P.spriteRam.Write(addr%0x100, data)
}

func (P *Ppu) Run(cycle uint) *RenderingData {
	P.cycle += cycle
	//fmt.Printf("%d %d\n", cycle, P.cycle)
	if P.line == 0 {
		P.background = []Tile{}
		P.buildSprites()
	}

	if P.cycle >= 341 {
		P.cycle -= 341
		P.line++
		if P.hasSpriteHit() {
			P.setSpriteHit()
		}
		if P.line <= 240 && P.line%8 == 0 && P.scrollY <= 240 {
			P.buildBackground()
		}
		if P.line == 241 {
			P.setVblank()
			if P.hasVblankIrqEnabled() {
				P.interrupts.AssertNmi()
			}
		}
		if P.line == 262 {
			P.clearVblank()
			P.clearSpriteHit()
			P.line = 0
			P.interrupts.ReleaseNmi()
			var bg []Tile = nil
			var sprites []SpriteWithAttribute = nil
			if P.isBackgroundEnable() {
				bg = P.background
			}
			if P.isSpriteEnable() {
				sprites = P.sprites
			}
			return NewRenderingData(
				P.getPalette(),
				bg,
				sprites,
			)
		}
	}
	return nil
}
