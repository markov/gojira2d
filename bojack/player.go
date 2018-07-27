package main

import (
	"github.com/go-gl/mathgl/mgl32"
	g "github.com/markov/gojira2d/pkg/graphics"
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math"
)

var (
	playersStopAtX = float32(550)
)

type Player struct {
	quad              *g.Primitive2D
	shadowQuad        *g.Primitive2D
	speed             float32
	key               glfw.Key
	keyPressed        bool
	position          mgl32.Vec3
	runningSprites    []*g.Texture
	numberOfFrames    int
	currentFrameIndex int
	animationSpeed    float32
	canStart          bool
	offsetXStartLine  float32
	playerName        string
}

func NewPlayer(
	position mgl32.Vec3,
	scale mgl32.Vec2,
	playerName string,
	numberOfFrames int,
	key glfw.Key,
	offsetXStartLine float32) *Player {
	p := &Player{}
	p.canStart = false
	p.offsetXStartLine = offsetXStartLine
	p.runningSprites = make([]*g.Texture, 0, numberOfFrames+1)
	for i := 0; i < numberOfFrames; i++ {
		p.runningSprites = append(
			p.runningSprites,
			g.NewTextureFromFile(fmt.Sprintf("bojack/sprites/%s/%s_%02d.png", playerName, playerName, i)))
	}
	p.playerName = playerName
	p.key = key
	p.position = position
	p.numberOfFrames = numberOfFrames
	p.currentFrameIndex = 0
	p.quad = g.NewQuadPrimitive(position, mgl32.Vec2{0, 0})
	p.quad.SetTexture(p.runningSprites[p.currentFrameIndex])
	p.quad.SetSizeFromTexture()
	p.quad.SetScale(scale)
	p.quad.SetAnchorToBottomCenter()

	p.shadowQuad = g.NewQuadPrimitive(position, mgl32.Vec2{0, 0})
	p.shadowQuad.SetTexture(g.NewTextureFromFile("bojack/sprites/shadow.png"))
	p.shadowQuad.SetSizeFromTexture()
	p.shadowQuad.SetScale(mgl32.Vec2{0.8, 0.6})
	p.shadowQuad.SetAnchorToCenter()
	return p
}

func (p *Player) UpdateIntro(scene *Scene) {
	p.updateSprite(scene)
	//log.Printf("%s calc: %f:", p.playerName, p.position.X() + p.offsetXStartLine)
	if p.position.X() >= playersStopAtX {
		p.canStart = true
	}
}

func (p *Player) Update(scene *Scene) {
	if !p.canStart {
		p.speed = 1.9
		p.UpdateIntro(scene)
	} else {
		p.RunRunRun(scene)
	}
}

func (p *Player) RunRunRun(scene *Scene) {
	if p.keyPressed {
		if p.speed < 9 {
			p.speed += 0.1
		}
	} else {
		p.speed /= 2
	}

	p.updateSprite(scene)
	//log.Printf("CURRENT POSITION #%d:", p.currentFrameIndex)
}

func (p *Player) updateSprite(scene *Scene) {
	if p.position.X() < scene.X()+100 && p.canStart == true {
		p.position = mgl32.Vec3{scene.X() + 100, p.position.Y(), p.position.Z()}
		p.speed = 0.8
	}

	p.animationSpeed += float32(math.Min(float64(p.speed), 3))
	p.currentFrameIndex = int(p.animationSpeed/10) % p.numberOfFrames
	absPos := p.position
	absPos = absPos.Add(mgl32.Vec3{p.speed, 0, 0})
	p.position = absPos
	p.quad.SetPosition(p.position.Sub(mgl32.Vec3{scene.X(), 0, 0}))
	p.shadowQuad.SetPosition(p.position.Sub(mgl32.Vec3{scene.X(), 0, -0.05}))
	p.quad.SetTexture(p.runningSprites[p.currentFrameIndex])
	scene.UpdatePlayerPos(p.position.X())
}

func (p *Player) Draw(ctx *g.Context) {
	p.quad.EnqueueForDrawing(ctx)
	p.shadowQuad.EnqueueForDrawing(ctx)
}
