package apu

import (
	"github.com/faiface/beep"
)

type Noise struct {
}

func NewNoise() *Noise {
	return &Noise{}
}

func (n *Noise) GetStreamer() beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			samples[i][0] = 0 //rand.Float64()*2 - 1
			samples[i][1] = 0 //rand.Float64()*2 - 1
		}
		return len(samples), true
	})
}

func (n *Noise) Write(addr byte, data byte) {
	//fmt.Printf("NO: 0x%02x - 0x%02x\n", addr, data)
	//n.ram.Write(addr, data)
}
