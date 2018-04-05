package graphics

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v4.1-core/gl"
	"gojira2d/pkg/utils"
)

type Matrix struct {
	Matrix mgl32.Mat4
	Dirty  bool
}

type Primitive2D struct {
	Primitive
	position          mgl32.Vec3
	scale             mgl32.Vec2
	size              mgl32.Vec2
	anchor            mgl32.Vec2
	angle             float32
	flipX             bool
	flipY             bool
	matrixSize        Matrix
	matrixTranslation Matrix
	matrixRotation    Matrix
	matrixScale       Matrix
	matrixAnchor      Matrix
}

func (p *Primitive2D) SetPosition(position mgl32.Vec3) {
	p.position = position
	p.matrixTranslation.Dirty = true
}

func (p *Primitive2D) SetAnchor(anchor mgl32.Vec2) {
	p.anchor = anchor
	p.matrixAnchor.Dirty = true
}

func (p *Primitive2D) SetAngle(radians float32) {
	p.angle = radians
	p.matrixRotation.Dirty = true
}

func (p *Primitive2D) SetScale(scale mgl32.Vec2) {
	p.scale = scale
	p.matrixScale.Dirty = true
}

func (p *Primitive2D) SetSize(size mgl32.Vec2) {
	p.size = size
	p.matrixSize.Dirty = true
}

func (p *Primitive2D) SetFlipX(flipX bool) {
	p.flipX = flipX
	p.matrixScale.Dirty = true
}

func (p *Primitive2D) SetFlipY(flipY bool) {
	p.flipY = flipY
	p.matrixScale.Dirty = true
}

func (p *Primitive2D) rebuildMatrices() {
	p.matrixTranslation.Matrix = mgl32.Translate3D(p.position.X(), p.position.Y(), p.position.Z())
	p.matrixScale.Matrix = mgl32.Scale3D(p.scale.X(), p.scale.Y(), 1)
	p.matrixRotation.Matrix = mgl32.HomogRotate3DZ(p.angle)
	p.matrixSize.Matrix = mgl32.Scale3D(p.size.X(), p.size.Y(), 1)
	p.matrixAnchor.Matrix = mgl32.Translate3D(-p.anchor.X(), -p.anchor.Y(), 0)
}

func (p *Primitive2D) Draw(context *Context) {
	shaderId := p.shaderProgram.Id()
	gl.BindTexture(gl.TEXTURE_2D, p.texture.Id())
	gl.UseProgram(shaderId)
	p.shaderProgram.SetUniform4fv("mProjection", &context.projectionMatrix)
	p.SetMatrices()
	gl.BindVertexArray(p.vaoId)
	gl.DrawArrays(p.arrayMode, 0, p.arraySize)
}

func (p *Primitive2D) EnqueueForDrawing(context *Context) {
	context.enqueueForDrawing(p)
}

func (p *Primitive2D) drawInBatch(context *Context) {
	p.SetMatrices()
	gl.BindVertexArray(p.vaoId)
	gl.DrawArrays(p.arrayMode, 0, p.arraySize)
}

func (p *Primitive2D) SetSizeFromTexture() {
	p.SetSize(mgl32.Vec2{float32(p.texture.width), float32(p.texture.height)})
}

func (p *Primitive2D) SetAnchorToCenter() {
	p.SetAnchor(mgl32.Vec2{p.size[0] / 2.0, p.size[1] / 2.0})
}

func (p *Primitive2D) invalidateMatrices() {
	p.matrixTranslation.Dirty = true
	p.matrixSize.Dirty = true
	p.matrixRotation.Dirty = true
	p.matrixScale.Dirty = true
	p.matrixAnchor.Dirty = true
}

func (p *Primitive2D) SetMatrices() {
	if p.matrixTranslation.Dirty {
		p.matrixTranslation.Matrix = mgl32.Translate3D(p.position.X(), p.position.Y(), p.position.Z())
		p.matrixTranslation.Dirty = false
		p.shaderProgram.SetUniform4fv("mTranslate", &p.matrixTranslation.Matrix)
	}
	if p.matrixScale.Dirty {
		scaleX := p.scale.X()
		if p.flipX {
			scaleX *= -1
		}
		scaleY := p.scale.Y()
		if p.flipY {
			scaleY *= -1
		}
		p.matrixScale.Matrix = mgl32.Scale3D(scaleX, scaleY, 1)
		p.shaderProgram.SetUniform4fv("mScale", &p.matrixScale.Matrix)
	}
	if p.matrixSize.Dirty {
		p.matrixSize.Matrix = mgl32.Scale3D(p.size.X(), p.size.Y(), 1)
		p.shaderProgram.SetUniform4fv("mSize", &p.matrixSize.Matrix)
	}
	if p.matrixRotation.Dirty {
		p.matrixRotation.Matrix = mgl32.HomogRotate3DZ(p.angle)
		p.shaderProgram.SetUniform4fv("mRotation", &p.matrixRotation.Matrix)
	}
	if p.matrixRotation.Dirty {
		p.matrixAnchor.Matrix = mgl32.Translate3D(-p.anchor.X(), -p.anchor.Y(), 0)
		p.shaderProgram.SetUniform4fv("mAnchor", &p.matrixAnchor.Matrix)
	}
}

func NewQuadPrimitive(position mgl32.Vec3, size mgl32.Vec2) (*Primitive2D) {
	q := Primitive2D{}
	q.position = position
	q.size = size
	q.scale = mgl32.Vec2{1, 1}
	q.shaderProgram = NewShaderProgram(vertexShaderPrimitive2D, "", FragmentShaderTexture)
	q.invalidateMatrices()

	q.arrayMode = gl.TRIANGLE_FAN
	q.arraySize = 4

	// Build the VAO
	gl.GenVertexArrays(1, &q.vaoId)
	gl.BindVertexArray(q.vaoId)

	// Vertices
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	vertices := []float32{0, 0, 0, 1, 1, 1, 1, 0}
	// 4 -> magic constant for float32 size
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	// Texture coordinates
	var vboUV uint32
	gl.GenBuffers(1, &vboUV)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboUV)
	uvCoordinates := []float32{0, 0, 0, 1, 1, 1, 1, 0}
	gl.BufferData(gl.ARRAY_BUFFER, len(uvCoordinates)*4, gl.Ptr(uvCoordinates), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	gl.BindVertexArray(0)
	return &q
}

func NewRegularPolygonPrimitive(position mgl32.Vec3, radius float32, numSegments int, showRotation bool) (*Primitive2D) {
	circlePoints := utils.CircleToPolygon(mgl32.Vec2{0, 0}, 0.5, numSegments, 0)

	q := Primitive2D{}
	q.position = position
	q.size = mgl32.Vec2{radius * 2, radius * 2}
	q.scale = mgl32.Vec2{1, 1}
	q.shaderProgram = NewShaderProgram(vertexShaderPrimitive2D, "", FragmentShaderSolidColor)
	q.invalidateMatrices()

	// Vertices
	vertices := make([]float32, 0, numSegments*2)
	for _, v := range circlePoints {
		vertices = append(vertices, v[0])
		vertices = append(vertices, v[1])
	}
	// Add one vertex for the last line
	vertices = append(vertices, circlePoints[0][0])
	vertices = append(vertices, circlePoints[0][1])
	// To show the rotation adds a line to the center
	if showRotation {
		vertices = append(vertices, 0)
		vertices = append(vertices, 0)
	}

	q.arrayMode = gl.LINE_STRIP
	q.arraySize = int32(len(vertices) / 2)

	// Build the VAO
	gl.GenVertexArrays(1, &q.vaoId)
	gl.BindVertexArray(q.vaoId)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	// 4 -> magic constant for float32 size
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	gl.BindVertexArray(0)
	return &q
}

const (
	vertexShaderPrimitive2D = `
        #version 410 core

        uniform mat4 mProjection;
        uniform mat4 mScale;
        uniform mat4 mTranslate;
        uniform mat4 mSize;
        uniform mat4 mRotation;
        uniform mat4 mAnchor;

        layout(location=0) in vec2 vertex;
        layout(location=1) in vec2 uv;

        out vec2 uv_out;

        void main() {
			mat4 mModel = mTranslate * mRotation * mScale * mAnchor * mSize;
            vec4 vertex_world = mModel * vec4(vertex, 0, 1);
            gl_Position = mProjection * vertex_world;
            uv_out = uv;
        }
        ` + "\x00"
)