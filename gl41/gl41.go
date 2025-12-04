package gl41

import (
	"fmt"
	"math"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

//
// ───────────────────────────────────────────────────────────────
//   VARIABILI GLOBALI
// ───────────────────────────────────────────────────────────────
//

// shader programs
var progQuad uint32
var progOverlay uint32
var progLine uint32

// fullscreen quad
var vaoQuad uint32
var vboQuad uint32

// overlay (cerchio + due rettangoli)
var vaoOverlay uint32
var vboOverlay uint32

// linee custom
var vaoLine uint32
var vboLine uint32

var winW int
var winH int

var winW1 int
var winH1 int

// ───────── linee custom ─────────

type Line struct {
	X1, Y1 float32
	X2, Y2 float32
	Alive  bool
	Group  int
	Color  [4]float32
}

var lines []Line
var lineVerts []float32

var GroupColors = map[int][4]float32{}
var defaultColor = [4]float32{0, 1, 0, 1} // verde

//
// ───────────────────────────────────────────────────────────────
//   UTILS
// ───────────────────────────────────────────────────────────────
//

func loadFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(b) + "\x00"
}

func compile(src string, shaderType uint32) uint32 {
	sh := gl.CreateShader(shaderType)
	csrc, free := gl.Strs(src)
	gl.ShaderSource(sh, 1, csrc, nil)
	free()
	gl.CompileShader(sh)

	var ok int32
	gl.GetShaderiv(sh, gl.COMPILE_STATUS, &ok)
	if ok == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(sh, gl.INFO_LOG_LENGTH, &logLen)
		log := make([]byte, logLen)
		gl.GetShaderInfoLog(sh, logLen, nil, &log[0])
		panic("GL SHADER ERROR:\n" + string(log))
	}
	return sh
}

func buildProgram(vs, fs string) uint32 {
	v := compile(vs, gl.VERTEX_SHADER)
	f := compile(fs, gl.FRAGMENT_SHADER)
	p := gl.CreateProgram()
	gl.AttachShader(p, v)
	gl.AttachShader(p, f)
	gl.LinkProgram(p)
	return p
}

// ───────────────────────────────────────────────────────────────
//
//	INIT 4.1
//
// ───────────────────────────────────────────────────────────────
func SetWindowSize(winSize, winSizeh, wsz1, hsz1 int) {
	winW = winSize
	winH = winSizeh
	winW1 = wsz1
	winH1 = hsz1
}

func InitGL41() {

	fmt.Println("Carico quad shaders")
	vsq := loadFile("gl41/shaders/quad.vert")
	fsq := loadFile("gl41/shaders/quad.frag")
	progQuad = buildProgram(vsq, fsq)
	fmt.Println("OK quad")

	// fullscreen quad
	quad := []float32{
		-1, -1,
		1, -1,
		1, 1,
		-1, -1,
		1, 1,
		-1, 1,
	}

	gl.GenVertexArrays(1, &vaoQuad)
	gl.BindVertexArray(vaoQuad)

	gl.GenBuffers(1, &vboQuad)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboQuad)
	gl.BufferData(gl.ARRAY_BUFFER, len(quad)*4, gl.Ptr(quad), gl.STATIC_DRAW)

	loc := uint32(gl.GetAttribLocation(progQuad, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(loc)
	gl.VertexAttribPointer(loc, 2, gl.FLOAT, false, 0, nil)

	gl.BindVertexArray(0)

	//
	// ─────────── overlay (cerchio + croci) ────────────
	//

	vso := loadFile("gl41/shaders/overlay.vert")
	fso := loadFile("gl41/shaders/overlay.frag")
	progOverlay = buildProgram(vso, fso)
	fmt.Println("OK overlay")

	const N = 64
	verts := make([]float32, 0, (N+8)*2)

	// circle
	for i := 0; i < N; i++ {
		ang := float64(i) * 2 * math.Pi / N
		verts = append(verts, float32(math.Cos(ang)), float32(math.Sin(ang)))
	}

	// horizontal rectangle (4 vertices)
	verts = append(verts,
		-1, -0.25,
		1, -0.25,
		-1, 0.25,
		1, 0.25,
	)

	// vertical rectangle
	verts = append(verts,
		-0.25, -1,
		0.25, -1,
		-0.25, 1,
		0.25, 1,
	)

	gl.GenVertexArrays(1, &vaoOverlay)
	gl.BindVertexArray(vaoOverlay)

	gl.GenBuffers(1, &vboOverlay)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboOverlay)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)

	//
	// ─────────── LINE SHADER ────────────
	//

	vsl := loadFile("gl41/shaders/line.vert")
	fsl := loadFile("gl41/shaders/line.frag")
	progLine = buildProgram(vsl, fsl)
	fmt.Println("OK line shader")

	gl.GenVertexArrays(1, &vaoLine)
	gl.GenBuffers(1, &vboLine)

	gl.BindVertexArray(vaoLine)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboLine)
	gl.BufferData(gl.ARRAY_BUFFER, 0, nil, gl.DYNAMIC_DRAW)

	// attribs: pos(2) + color(4) = 6 float
	stride := int32(6 * 4)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, stride, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

//
// ───────────────────────────────────────────────────────────────
//   LINEE CUSTOM
// ───────────────────────────────────────────────────────────────
//

func SetGroupColor(g int, r, b, g2, a float32) {
	GroupColors[g] = [4]float32{r, b, g2, a}
}

func AddLine(x1, y1, x2, y2 float32, group int) int {
	col, ok := GroupColors[group]
	if !ok {
		col = defaultColor
	}

	l := Line{x1, y1, x2, y2, true, group, col}
	lines = append(lines, l)
	rebuildLineVBO()
	return len(lines) - 1
}

func HideGroup(g int) {
	for i := range lines {
		if lines[i].Group == g {
			lines[i].Alive = false
		}
	}
	rebuildLineVBO()
}

func ShowGroup(g int) {
	for i := range lines {
		if lines[i].Group == g {
			lines[i].Alive = true
		}
	}
	rebuildLineVBO()
}

func RemoveGroup(g int) {
	for i := range lines {
		if lines[i].Group == g {
			lines[i].Alive = false
		}
	}
	CompactLines()
}

func RemoveLine(id int) {
	if id < 0 || id >= len(lines) {
		return
	}
	lines[id].Alive = false
	rebuildLineVBO()
}

func CompactLines() {
	tmp := lines[:0]
	for _, l := range lines {
		if l.Alive {
			tmp = append(tmp, l)
		}
	}
	lines = tmp
	rebuildLineVBO()
}

func rebuildLineVBO() {
	lineVerts = lineVerts[:0]

	for _, l := range lines {
		if !l.Alive {
			continue
		}
		r, g, b, a := l.Color[0], l.Color[1], l.Color[2], l.Color[3]

		lineVerts = append(lineVerts,
			l.X1, l.Y1, r, g, b, a,
			l.X2, l.Y2, r, g, b, a,
		)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, vboLine)
	if len(lineVerts) > 0 {
		gl.BufferData(gl.ARRAY_BUFFER, len(lineVerts)*4, gl.Ptr(lineVerts), gl.DYNAMIC_DRAW)
	} else {
		gl.BufferData(gl.ARRAY_BUFFER, 0, nil, gl.DYNAMIC_DRAW)
	}
}

func RenderLines() {
	// winW1 = w
	// winH1 = h
	gl.UseProgram(progLine)

	gl.Uniform2f(gl.GetUniformLocation(progLine, gl.Str("uWinSize\x00")),
		float32(winW1), float32(winH1))

	gl.BindVertexArray(vaoLine)
	gl.DrawArrays(gl.LINES, 0, int32(len(lineVerts)/6))
	gl.BindVertexArray(0)
}

//
// ───────────────────────────────────────────────────────────────
//   TEXTURE
// ───────────────────────────────────────────────────────────────
//

func CreateTexture() uint32 {
	// winW = w
	// winH = h

	var t uint32
	gl.GenTextures(1, &t)
	gl.BindTexture(gl.TEXTURE_2D, t)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8,
		int32(winW), int32(winH), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, nil)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	return t
}

func UpdateTexture(tex uint32, pixels []uint32) {
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0,
		0, 0,
		int32(winW), int32(winH),
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
}

//
// ───────────────────────────────────────────────────────────────
//   RENDER
// ───────────────────────────────────────────────────────────────
//

func RenderFullscreen(tex uint32) {
	gl.UseProgram(progQuad)

	loc := gl.GetUniformLocation(progQuad, gl.Str("tex\x00"))
	gl.Uniform1i(loc, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	gl.BindVertexArray(vaoQuad)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)
}

func DrawOverlay(cx, cy, radius float32) {
	gl.UseProgram(progOverlay)
	gl.BindVertexArray(vaoOverlay)

	gl.Uniform2f(gl.GetUniformLocation(progOverlay, gl.Str("uCenterPx\x00")), cx, cy)
	gl.Uniform1f(gl.GetUniformLocation(progOverlay, gl.Str("uRadiusPx\x00")), radius)
	gl.Uniform2f(gl.GetUniformLocation(progOverlay, gl.Str("uWinSize\x00")), float32(winW1), float32(winH1))

	// cerchio
	gl.Uniform1i(gl.GetUniformLocation(progOverlay, gl.Str("uMode\x00")), 0)
	gl.Uniform1f(gl.GetUniformLocation(progOverlay, gl.Str("uThickness\x00")), 4)
	gl.Uniform4f(gl.GetUniformLocation(progOverlay, gl.Str("uColor\x00")), 1, 1, 0, 1)
	gl.DrawArrays(gl.LINE_LOOP, 0, 64)

	// rettangolo orizzontale
	gl.Uniform1i(gl.GetUniformLocation(progOverlay, gl.Str("uMode\x00")), 1)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 64, 4)

	// rettangolo verticale
	gl.Uniform1i(gl.GetUniformLocation(progOverlay, gl.Str("uMode\x00")), 2)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 68, 4)

	gl.BindVertexArray(0)
}
