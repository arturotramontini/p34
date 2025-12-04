package gl41

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	ovVAO       uint32
	ovVBO       uint32
	ovProg      uint32
	ovCircleCnt = 36
)

//------------------------------------------------------------
// SHADERS OPENGL 4.1 CORE
//------------------------------------------------------------

var overlayVert410 = `
#version 410 core

layout(location = 0) in vec2 inPos;

uniform vec2  uCenterPx;    // centro in pixel
uniform float uRadiusPx;    // raggio in pixel
uniform vec2  uWinSize;     // dimensioni finestra (pixel)
uniform int   uMode;        // 0 = vert-line, 1 = horiz-line, 2 = circle
uniform float uThickness;   // spessore linea in pixel

out float vDist;            // distanza dal centro (per thickness)
out vec2  vPosPx;           // posizione in pixel

void main() {

    vec2 p = inPos;

    if (uMode == 0) {
        // linea verticale: x=0
        p = vec2(0.0, inPos.y) * uRadiusPx + uCenterPx;
    }
    else if (uMode == 1) {
        // linea orizzontale: y=0
        p = vec2(inPos.x, 0.0) * uRadiusPx + uCenterPx;
    }
    else if (uMode == 2) {
        // cerchio: posizione già unitaria
        p = inPos * uRadiusPx + uCenterPx;
    }

    vPosPx = p;

    // distanza dal centro (per linewidth simulato)
    vDist = length(p - uCenterPx);

    // convert pixel → clip space
    float cx = (p.x / uWinSize.x) * 2.0 - 1.0;
    float cy = 1.0 - (p.y / uWinSize.y) * 2.0;
    gl_Position = vec4(cx, cy, 0.0, 1.0);
}
` + "\x00"

var overlayFrag410 = `
#version 410 core

in float vDist;
in vec2  vPosPx;

out vec4 outColor;

uniform vec4  uColor;
uniform float uRadiusPx;
uniform float uThickness;
uniform int   uMode;

void main() {

    float alpha = 1.0;

    if (uMode == 2) {
        // ---- CERCHIO ----
        float d = abs(vDist - uRadiusPx);
        if (d > uThickness) discard;   // bordo "sottile"
        alpha = 1.0 - d / uThickness;  // anti-alias leggero
    }
    else {
        // ---- LINEE ----
        // distance from line center = distance from uCenter
        float dx = abs(vPosPx.x - vPosPx.y); // trick per anti-alias
        alpha = 1.0;
    }

    outColor = vec4(uColor.rgb, uColor.a * alpha);
}
` + "\x00"

//------------------------------------------------------------
// CREAZIONE VBO: cerchio + 2 linee
//------------------------------------------------------------

func buildOverlayVBO() (uint32, uint32) {

	// cerchio unitario
	N := ovCircleCnt
	verts := make([]float32, 0, (N+4)*2)

	for i := 0; i < N; i++ {
		ang := float64(i) * 2 * math.Pi / float64(N)
		verts = append(verts, float32(math.Cos(ang)))
		verts = append(verts, float32(math.Sin(ang)))
	}

	// linea orizzontale (-1..1)
	verts = append(verts, -1, 0)
	verts = append(verts, 1, 0)

	// linea verticale
	verts = append(verts, 0, -1)
	verts = append(verts, 0, 1)

	var vao, vbo uint32

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return vao, vbo
}

//------------------------------------------------------------
// INIT OVERLAY (da chiamare dopo gl.Init())
//------------------------------------------------------------

func InitOverlayGL41(winWpx, winHpx int) error {

	var err error
	ovProg, err = CreateProgram(overlayVert410, overlayFrag410)
	if err != nil {
		return fmt.Errorf("overlay shader error: %v", err)
	}

	ovVAO, ovVBO = buildOverlayVBO()

	return nil
}

//------------------------------------------------------------
// DISEGNO SEMPLIFICATO (cross + circle)
//------------------------------------------------------------

func drawOne(mode int32, cx, cy, radius float32, rgba [4]float32, thickness float32, winW, winH float32) {

	gl.UseProgram(ovProg)

	gl.BindVertexArray(ovVAO)

	// uniform pixel-space
	gl.Uniform2f(gl.GetUniformLocation(ovProg, gl.Str("uCenterPx\x00")), cx, cy)
	gl.Uniform1f(gl.GetUniformLocation(ovProg, gl.Str("uRadiusPx\x00")), radius)
	gl.Uniform2f(gl.GetUniformLocation(ovProg, gl.Str("uWinSize\x00")), winW, winH)
	gl.Uniform1f(gl.GetUniformLocation(ovProg, gl.Str("uThickness\x00")), thickness)
	gl.Uniform1i(gl.GetUniformLocation(ovProg, gl.Str("uMode\x00")), mode)
	gl.Uniform4f(gl.GetUniformLocation(ovProg, gl.Str("uColor\x00")),
		rgba[0], rgba[1], rgba[2], rgba[3])

	N := int32(ovCircleCnt)

	switch mode {
	case 0: // vertical line
		gl.DrawArrays(gl.LINES, N+2, 2)
	case 1: // horizontal line
		gl.DrawArrays(gl.LINES, N, 2)
	case 2: // circle
		gl.DrawArrays(gl.LINE_LOOP, 0, N)
	}

	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

//------------------------------------------------------------
// FUNZIONE PUBBLICA (da chiamare nel main loop)
//------------------------------------------------------------

func DrawOverlayGL41(cx, cy, radius float32, winW, winH float32) {

	thickness := float32(2.0) // pixel

	// cerchio blu
	drawOne(2, cx, cy, radius, [4]float32{0, 0, 1, 1}, thickness, winW, winH)

	// linea orizzontale verde
	drawOne(1, cx, cy, radius, [4]float32{0, 1, 0, 1}, thickness, winW, winH)

	// linea verticale rossa
	drawOne(0, cx, cy, radius, [4]float32{1, 0, 0, 1}, thickness, winW, winH)
}
