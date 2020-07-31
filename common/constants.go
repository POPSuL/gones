package common

import "time"

const CpuClock = uint(1789772)
const FrameCounterRate = CpuClock / 240
const AudioFreq = 44100
const AudioBuffer = time.Second / 10
