// Bubbles ScreenSaver
package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math"
	"math/rand"
	"time"
)

type inst struct {
	x, y, r, xS, yS float64
	col             color.RGBA
	xFlag, yFlag    bool
}

var w, h float64 = 1920, 1080
var numBubbles = 50
var fps time.Duration = 120
var instSet []*inst
var boolSet = []bool{false, true}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Bubbles ScreenSaver",
		Bounds: pixel.R(0, 0, w, h),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	tick := time.Tick(time.Second / fps)
	start := time.Now()
	i := 0.0
	for !win.Closed() {
		// Calculate FPS
		if i >= math.MaxFloat64 {
			i = 0.0
			start = time.Now()
		}
		i++
		timer := time.Since(start).Seconds()

		// By click generate bubble configuration
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			config(win)
		}

		// Trick for clear last frame
		imd.Color = pixel.RGB(0, 0, 0)
		imd.Push(pixel.V(0, 0))
		imd.Push(pixel.V(w, h))
		imd.Rectangle(0)
		imd.Draw(win)

		imd.Clear()
		for _, v := range instSet {
			draw(v, imd, win)
		}

		atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		txt := text.New(pixel.V(40, 550), atlas)
		_, _ = fmt.Fprint(txt, "Bubbles ScreenSaver\n")
		_, _ = fmt.Fprint(txt, "Date: "+time.Now().Format(time.RFC1123)+"\n")
		_, _ = fmt.Fprint(txt, "FPS: ", math.Round(i/timer))
		txt.Draw(win, pixel.IM)

		win.Update()
		<-tick
	}
}

func config(win *pixelgl.Window) {
	if len(instSet) <= numBubbles {
		r := float64(rand.Intn(60) + 30)
		s1 := (100 - r) * 0.025
		if r == 100 {
			s1 = 0.025
		}
		s2 := s1 / float64(rand.Intn(60)+30) * 12
		if s2 > s1 {
			s2 = s1
		}
		spdSet := []float64{s1, s2}
		spdSel := uint8(rand.Intn(2))

		instSet = append(instSet, &inst{
			x:  win.MousePosition().X,
			y:  win.MousePosition().Y,
			r:  r,
			xS: spdSet[spdSel],
			yS: spdSet[1&^uint8(spdSel)],
			col: color.RGBA{
				R: uint8(rand.Intn(255)),
				G: uint8(rand.Intn(255)),
				B: uint8(rand.Intn(255)),
				A: uint8(200),
			},
			xFlag: boolSet[rand.Intn(2)],
			yFlag: boolSet[rand.Intn(2)],
		})
	}
}

func draw(v *inst, imd *imdraw.IMDraw, win *pixelgl.Window) {
	x, y, r, xS, yS, xFlag, yFlag := &v.x, &v.y, v.r, v.xS, v.yS, &v.xFlag, &v.yFlag
	// State for x-axis
	switch {
	case *x+r > w:
		*x -= xS
		*xFlag = true
	case (w+*x)-w < r:
		*x += xS
		*xFlag = false
	case *xFlag == true:
		*x -= xS
	default:
		*x += xS
	}

	// State for y-axis
	switch {
	case *y+r > h:
		*y -= yS
		*yFlag = true
	case (h+*y)-h < r:
		*y += yS
		*yFlag = false
	case *yFlag == true:
		*y -= yS
	default:
		*y += yS
	}

	imd.Color = v.col
	imd.Push(pixel.V(*x, *y))
	imd.Circle(r, 0)
	imd.Draw(win)
}
