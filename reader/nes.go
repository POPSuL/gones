package reader

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const NES_HEADER_SIZE = 0x0010
const PROGRAM_ROM_SIZE = 0x4000
const CHARACTER_ROM_SIZE = 0x2000

type NesRom struct {
	HorizontalMirror  bool
	Program           []byte
	Character         []byte
	ProgramRomPages   uint
	CharacterRomPages uint
	Mapper            uint
}

func ReadRom(file string) *NesRom {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	if buffer[0] != 'N' || buffer[1] != 'E' || buffer[2] != 'S' {
		panic(errors.New("Invalid NES file"))
	}

	fmt.Printf("Rom size: %d (0x%x)\n", len(buffer), len(buffer))

	rom := new(NesRom)
	rom.ProgramRomPages = uint(buffer[4])
	rom.CharacterRomPages = uint(buffer[5])
	rom.HorizontalMirror = uint(buffer[6])&0x01 != 1
	rom.Mapper = (uint(buffer[6])&0xf0)>>4 | (uint(buffer[7]) & 0xf0)

	fmt.Printf("Program ROM pages: %d\n", rom.ProgramRomPages)
	fmt.Printf("Character ROM pages: %d\n", rom.CharacterRomPages)
	fmt.Printf("Mapper: %d\n", rom.Mapper)

	characterRomStart := NES_HEADER_SIZE + rom.ProgramRomPages*PROGRAM_ROM_SIZE
	characterRomEnd := characterRomStart + rom.CharacterRomPages*CHARACTER_ROM_SIZE
	fmt.Printf("Character ROM start: 0x%x (%d)\n", characterRomStart, characterRomStart)
	fmt.Printf("Character ROM end: 0x%x (%d)\n", characterRomEnd, characterRomEnd)

	rom.Program = buffer[NES_HEADER_SIZE:characterRomStart]
	rom.Character = buffer[characterRomStart : characterRomStart+(characterRomEnd-characterRomStart)]

	fmt.Printf(
		"Program   ROM: 0x0000 - 0x%x (%d bytes)\n",
		len(rom.Program)-1,
		len(rom.Program))
	fmt.Printf(
		"Character   ROM: 0x0000 - 0x%x (%d bytes)\n",
		len(rom.Character)-1,
		len(rom.Character))

	return rom
}
