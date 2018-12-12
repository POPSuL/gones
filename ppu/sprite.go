package ppu

type SpriteWithAttribute struct {
	x, y, attribute, id uint
	sprite              [][]uint
}

func NewSpriteWithAttribute(sprite [][]uint, x uint, y uint, attribute uint, id uint) *SpriteWithAttribute {
	return &SpriteWithAttribute{
		x,
		y,
		attribute,
		id,
		sprite,
	}
}
