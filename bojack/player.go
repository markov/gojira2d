package main

import (
	"github.com/go-gl/mathgl/mgl32"
	g "github.com/markov/gojira2d/pkg/graphics"
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Player struct {
	quad                  *g.Primitive2D
	speed                 float32
	key                   glfw.Key
	keyPressed            bool
	position              mgl32.Vec3
	runningSprites        []*g.Texture
	numberOfFrames        int
	currentSpritePosition int
	animationSpeed        float32
}

func NewPlayer(position mgl32.Vec3, scale mgl32.Vec2, playerName string, numberOfFrames int, key glfw.Key) *Player {
	p := &Player{}
	p.runningSprites = make([]*g.Texture, 0, numberOfFrames + 1)
	for i := 0; i < numberOfFrames; i++ {
		var spriteNumber string
		if i < 10 {
			spriteNumber = fmt.Sprintf("0%d", i)
		} else {
			spriteNumber = fmt.Sprintf("%d", i)
		}
		p.runningSprites = append(
			p.runningSprites,
			g.NewTextureFromFile(fmt.Sprintf("bojack/sprites/%s/%s_%s.png", playerName, playerName, spriteNumber)))

	}
	p.key = key
	p.position = position
	p.numberOfFrames = numberOfFrames
	p.currentSpritePosition = 0
	p.quad = g.NewQuadPrimitive(position, mgl32.Vec2{0, 0})
	p.quad.SetTexture(p.runningSprites[p.currentSpritePosition])
	p.quad.SetSizeFromTexture()
	p.quad.SetScale(scale)
	p.quad.SetAnchorToBottomCenter()
	return p
}

func (p *Player) Update(scene *Scene) {
	if p.keyPressed { //&& shouldPress() {
		p.speed += 0.1
	} else {
		p.speed /= 2
	}
	p.animationSpeed += p.speed
	//if p.animationSpeed > 12 {
	//	p.animationSpeed = 12
	//}
	p.currentSpritePosition = int(p.animationSpeed/10) % p.numberOfFrames
	absPos := p.position
	absPos = absPos.Add(mgl32.Vec3{p.speed, 0, 0})
	p.position = absPos
	p.quad.SetPosition(p.position.Sub(mgl32.Vec3{scene.X(), 0, 0}))
	p.quad.SetTexture(p.runningSprites[p.currentSpritePosition])
	//log.Printf("CURRENT POSITION #%d:", p.currentSpritePosition)
	scene.UpdatePlayerPos(p.position.X())
}

func (p *Player) Draw(ctx *g.Context) {
	p.quad.EnqueueForDrawing(ctx)
}
