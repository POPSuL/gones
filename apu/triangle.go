package apu

import (
	"github.com/faiface/beep"
	"github.com/popsul/gones/common"
)

type Triangle struct {
	linearCounter         uint
	lengthCounter         uint
	dividerForFrequency   uint
	frequency             uint
	lastLevel             float32
	direction             int16
	isLengthCounterEnable bool
	stopped               bool
}

func NewTriangle() *Triangle {
	return &Triangle{
		direction: 1,
	}
}

func (t *Triangle) getNextSample() int8 {
	if t.frequency == 0 {
		return 0
	}

	v := t.lastLevel
	x := v + float32(t.frequency)/0.1/256.*float32(t.direction)
	//x := (v + float32(common.AudioFreq) / float32(t.frequency) / 4) * float32(t.direction)
	//x := (v + (float32(t.frequency) / float32(common.AudioFreq) * 256)) * float32(t.direction)
	if x <= -127 {
		t.direction = 1
		x = -127
	} else if x >= 128 {
		t.direction = -1
		x = 128
	}
	t.lastLevel = x
	//println(int64(float32(common.AudioFreq) / float32(t.frequency)))
	//println(int64(x))
	//println(int64(float32(common.AudioFreq) / float32(t.frequency) / 2))
	return int8(x)
}

func (t *Triangle) GetStreamer() beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		//println(len(samples))
		//println(t.frequency)
		for i := range samples {
			sample := float64(t.getNextSample()+127) / 256.0
			//println(t.getNextSample()+127)
			//println(int64(float64(t.getNextSample()+127)))
			if t.stopped {
				samples[i][0] = 0 //rand.Float64()*2 - 1
				samples[i][1] = 0 //rand.Float64()*2 - 1
			} else {
				samples[i][0] = sample //rand.Float64()*2 - 1
				samples[i][1] = sample //rand.Float64()*2 - 1
			}
		}
		return len(samples), true
	})
}

func (t *Triangle) Write(addr byte, data byte) {
	//fmt.Printf("TR: 0x%02x - 0x%02x\n", addr, data)
	if addr == 0x00 {
		t.isLengthCounterEnable = !common.I2b(uint(data & 0x80))
		t.linearCounter = uint(data & 0x7F)
		//this.oscillator.setVolume(this.volume);
		return
	}
	if addr == 0x02 {
		t.dividerForFrequency &= 0x700
		t.dividerForFrequency |= uint(data)
		return
	}
	if addr == 0x03 {
		// Programmable timer, length counter
		t.dividerForFrequency &= 0xFF
		t.dividerForFrequency |= uint((data & 0x7) << 8)
		if t.isLengthCounterEnable {
			t.lengthCounter = CounterTable[data>>3]
		}
		t.frequency = common.CpuClock / ((t.dividerForFrequency + 1) * 32)
		t.stopped = false
	}
}

func (t *Triangle) UpdateCounter() {
	if t.isLengthCounterEnable && t.lengthCounter > 0 {
		t.lengthCounter--
	}
	if t.linearCounter > 0 {
		t.linearCounter--
	}
	if t.lengthCounter == 0 && t.linearCounter == 0 {
		t.stopped = true
	}
}
