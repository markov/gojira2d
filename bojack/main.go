package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/markov/gojira2d/pkg/app"
	"log"
	"github.com/go-gl/mathgl/mgl32"
	g "github.com/markov/gojira2d/pkg/graphics"
	)

var (
	keyCallbackFunc glfw.KeyCallback
)

func main() {
	app.Init(640, 480, false, "Run For Your Life!")
	defer app.Terminate()
	defer UnregisterKeyCallback()

	player := NewPlayer(mgl32.Vec3{15,15,0})
	RegisterKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Release {
			log.Printf("#%d key:", key)
			player.Update()
		}
	})
	app.MainLoop(func(speed float64) {
		//NOOP
	}, func() {
		player.Draw(app.Context)
	})
}
type Player struct {
	quad      *g.Primitive2D
	position   mgl32.Vec3
}

func NewPlayer(position mgl32.Vec3) *Player {
	p := &Player{}
	p.position = position
	p.quad = g.NewQuadPrimitive(mgl32.Vec3{position.X(), position.Y(), 0}, mgl32.Vec2{30, 30})
	p.quad.SetAnchorToCenter()
	p.quad.SetTexture(g.NewTextureFromFile("examples/assets/texture.png"))
	return p
}

func (p *Player) Update() {
	absPos := p.position
	absPos = absPos.Add(mgl32.Vec3{10, 0, 0})
	p.position = absPos
	p.quad.SetPosition(p.position)
}

func (p *Player) Draw(ctx *g.Context) {
	p.quad.EnqueueForDrawing(ctx)
}

func RegisterKeyCallback(callback glfw.KeyCallback) {
	if keyCallbackFunc != nil {
		log.Panic("A keyboard key-callback is already registered!")
	}
	keyCallbackFunc = callback
	app.GetWindow().SetKeyCallback(callback)
}

func UnregisterKeyCallback() {
	keyCallbackFunc = nil
	app.GetWindow().SetKeyCallback(nil)
}