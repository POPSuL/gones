package apu

import (
	"github.com/faiface/beep"
)

type Square struct {
	stopped bool
}

func NewSquare() *Square {
	return &Square{}
}

func (s *Square) GetStreamer() beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			if s.stopped {
				samples[i][0] = 0
				samples[i][1] = 0
			} else {

			}
		}
		return len(samples), true
	})
}

func (s *Square) Write(addr byte, data byte) {
	//fmt.Printf("SQ: 0x%02x - 0x%02x\n", addr, data)
}
