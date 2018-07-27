package main

import (
	"github.com/go-gl/mathgl/mgl32"
	g "github.com/markov/gojira2d/pkg/graphics"
	"github.com/go-gl/glfw/v3.2/glfw"
	"container/list"
	"math/rand"
	"github.com/markov/gojira2d/pkg/ui"
)

const windowOfOpportunity = 0.2

type window struct {
	w, h int
}
type bar struct {
	creationTime float32
	endTime      float32
	size         float32
	quad         *g.Primitive2D
}

var (
	win              = window{w: 1920, h: 1080}
	buttonPressed    *g.Primitive2D
	buttonReleased   *g.Primitive2D
	bars             *list.List
	sizeInterpolator = float32(win.w-80) / 3
	barStart         = float32(20)
	barEnd           = 3 * sizeInterpolator
	gogoTc           *ui.Text
	gogoTimer        float32
)

func createHud() {
	buttonPressed = g.NewQuadPrimitive(
		mgl32.Vec3{float32(win.w - 60), 30, -1},
		mgl32.Vec2{48, 40},
	)
	buttonPressed.SetTexture(g.NewTextureFromFile("bojack/sprites/button/button_pressed.png"))

	buttonReleased = g.NewQuadPrimitive(
		mgl32.Vec3{float32(win.w - 60), 30, -1},
		mgl32.Vec2{48, 40},
	)
	buttonReleased.SetTexture(g.NewTextureFromFile("bojack/sprites/button/button_unpressed.png"))
	bars = list.New()
}

func createGoGoGo() {
	gogoFont := ui.NewFontFromFiles(
		"regular",
		"examples/assets/fonts/roboto-mono-regular.fnt",
		"examples/assets/fonts/roboto-mono-regular.png",
	)
	gogoColor := g.Color{
		255,
		0,
		0,
		0.9,
	}
	gogoTc = ui.NewText(
		"GO GO GO!!!!",
		gogoFont,
		mgl32.Vec3{300, 200, 0.1},
		mgl32.Vec2{200, 150},
		gogoColor,
		mgl32.Vec4{0, 0, 0, -.17},
	)
}

func drawGoGoGo(ctx *g.Context, player *Player) {
	gogoTimer = float32(glfw.GetTime())
	//if gogoTimer > 8.5 && gogoTimer < 9.5 {
	//	gogoTc.EnqueueForDrawing(ctx)
	//}
	if player.canStart && gogoTimer < 9.8 {
		gogoTc.EnqueueForDrawing(ctx)
	}
}

func updateHud() {
	time := float32(glfw.GetTime())
	if bars.Len() == 0 || bars.Front().Value.(bar).endTime < time && rand.Int31n(100) > 95 {
		duration := rand.Float32()
		size := duration * sizeInterpolator
		newBar := bar{
			time,
			time + duration,
			size,
			g.NewQuadPrimitive(
				mgl32.Vec3{0, 10, -1},
				mgl32.Vec2{size, 60},
			),
		}
		newBar.quad.SetTexture(g.NewTextureFromFile("bojack/sprites/colors/blue.png"))
		bars.PushFront(newBar)
	}

	for e := bars.Front(); e != nil; e = e.Next() {
		bar := e.Value.(bar)
		if bar.endTime+3 < time {
			bars.Remove(e)
			continue
		}

		barX := (time-bar.creationTime)*sizeInterpolator - bar.size
		barCutOff := float32(0)
		if barX < barStart {
			barCutOff = barStart - barX
			barX = barStart
		}
		barWidth := bar.size - barCutOff
		if barX+barWidth > barEnd {
			barCutOff = barX + barWidth - barEnd
		}
		barWidth = bar.size - barCutOff
		if barCutOff > bar.size {
			barWidth = 0
		}
		bar.quad.SetSize(mgl32.Vec2{
			barWidth,
			60,
		})
		bar.quad.SetPosition(mgl32.Vec3{
			barX,
			10,
			-1,
		})
	}
}

func shouldPress() bool {
	if bars.Len() == 0 {
		return false
	}
	lastBar := bars.Back().Value.(bar)
	endTime := float32(glfw.GetTime()) - 3
	return lastBar.creationTime < endTime && lastBar.endTime > endTime
}

func pressOpportunity() bool {
	if bars.Len() == 0 {
		return false
	}
	lastBar := bars.Back().Value.(bar)
	endTime := float32(glfw.GetTime()) - 3
	return mgl32.Abs(lastBar.creationTime-endTime) < windowOfOpportunity
}

func releaseOpportunity() bool {
	if bars.Len() == 0 {
		return false
	}
	lastBar := bars.Back().Value.(bar)
	endTime := float32(glfw.GetTime()) - 3
	return mgl32.Abs(lastBar.endTime-endTime) < windowOfOpportunity
}

func drawHud(ctx *g.Context) {
	if bars.Back() == nil {
		buttonReleased.EnqueueForDrawing(ctx)
		return
	}
	for e := bars.Front(); e != nil; e = e.Next() {
		e.Value.(bar).quad.EnqueueForDrawing(ctx)
	}

	if shouldPress() {
		buttonPressed.EnqueueForDrawing(ctx)
	} else {
		buttonReleased.EnqueueForDrawing(ctx)
	}
}
