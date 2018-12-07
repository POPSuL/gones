package ppu

type Tile struct {
	paletteId, scrollX, scrollY uint
	pattern                     []interface{}
}

func NewTile(pattern []interface{}, paletteId uint, scrollX uint, scrollY uint) *Tile {
	return &Tile{
		paletteId,
		scrollX,
		scrollY,
		pattern,
	}
}
