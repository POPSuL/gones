package ppu

type Tile struct {
	paletteId, scrollX, scrollY uint
	Pattern                     [][]uint
}

func NewTile(pattern [][]uint, paletteId uint, scrollX uint, scrollY uint) *Tile {
	return &Tile{
		paletteId,
		scrollX,
		scrollY,
		pattern,
	}
}
