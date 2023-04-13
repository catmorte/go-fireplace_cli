package main

import (
	"flag"
	"math"
	"math/rand"
	"sync"
	"time"

	tm "github.com/buger/goterm"
	"github.com/gookit/color"
)

var (
	speed  *time.Duration = flag.Duration("speed", 100*time.Millisecond, "speed of fire")
	symbol                = flag.String("symbol", "$", "symbol of fire")
	isBG                  = flag.Bool("bg", true, "use background color")

	xl, yl, sourceX, sparkY int
	fireMatrix              [][]float32
	fireLock                sync.RWMutex
	sparkColorizer          func(float32) string
)

func init() {
	flag.Parse()
	if *isBG {
		sparkColorizer = bgFire
	} else {
		sparkColorizer = fgFire
	}
	xl, yl = getFireSize()
	resetFire()
}

func getFireSize() (int, int) {
	return tm.Width(), tm.Height()
}

func resetFire() {
	tm.Clear()
	tm.MoveCursor(0, 0)
	sourceX, sparkY = int(float32(xl)*0.1), int(float32(yl)*0.5)
	fireMatrix = make([][]float32, yl)
	for y := 0; y < yl; y++ {
		fireMatrix[y] = make([]float32, xl)
		for x := 0; x < xl; x++ {
			fireMatrix[y][x] = 0
		}
	}
}

func bgFire(r float32) string {
	return color.RGB(uint8(math.Floor(float64(r))), 0, 0, true).Sprint(tm.Color(*symbol, tm.BLACK))
}

func fgFire(r float32) string {
	return color.RGB(uint8(math.Floor(float64(r))), 0, 0, false).Sprint(*symbol)
}

func main() {
	go func() {
		for {
			newXl, newYl := getFireSize()
			if xl != newXl || yl != newYl {
				fireLock.Lock()
				xl, yl = newXl, newYl
				resetFire()
				fireLock.Unlock()
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		fireLock.RLock()
		topFire()
		for y := yl - 1; y >= 0; y-- {
			for x := xl - 1; x >= 0; x-- {
				val := fireMatrix[y][x]
				if val < 0.2 {
					val = 0.1
				}
				if val > 0.7 {
					val = 1
				}
				r := (255.0 * float32(float32(yl-y)/float32(yl))) * float32(val)
				tm.MoveCursor(xl-x, yl-y)
				tm.Print(sparkColorizer(r))
			}
		}
		tm.Flush()
		fireLock.RUnlock()
		time.Sleep(*speed)
	}
}

func topFire() {
	for x := 0; x < xl; x++ {
		newVal := rand.Float32()
		fireMatrix[0][x] = newVal
	}

	for i := 0; i < sourceX; i++ {
		newX := rand.Intn(xl)
		newY := rand.Intn(sparkY)
		fireMatrix[newY][newX] = 0.
	}

	for y := yl - 1; y > 0; y-- {
		for x := 0; x < xl; x++ {
			order := rand.Float32()
			val := float32(0.0)
			if order > 0.7 {
				if (x - 1) > 0 {
					val = fireMatrix[y-1][x-1]
				}
			} else if order > 0.3 {
				val = fireMatrix[y-1][x] * 0.9
			} else {
				if (x + 1) < xl {
					val = fireMatrix[y-1][x+1] * 0.97
				}
			}
			fireMatrix[y][x] = val
		}
	}
}
