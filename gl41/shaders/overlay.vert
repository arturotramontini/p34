
#version 410 core

layout (location = 0) in vec2 position2;

uniform vec2 uCenterPx;     // centro (pixel)
uniform float uRadiusPx;    // raggio (pixel)
uniform vec2 uWinSize;      // dimensione finestra (pixel)
uniform float uThickness;   // spessore linee (pixel)
uniform int uMode;          // 0=circle, 1=hor, 2=vert

out vec2 vLocal;            // per shader fragment

void main()
{
    vec2 p;

    if (uMode == 0) {
        // Cerchio (coordinate già unit circle)
        p = position2 * uRadiusPx + uCenterPx;
        vLocal = position2;
    }
    else if (uMode == 1) {
        // ORIZZONTALE → rettangolo ([-1,+1] x [-0.5,+0.5])
        vec2 scale = vec2(uRadiusPx, uThickness * 0.5);
        p = position2 * scale;
        p += uCenterPx;
        vLocal = position2;
    }
    else {
        // VERTICALE → rettangolo ([-1,+1] x [-0.5,+0.5])
        vec2 scale = vec2(uThickness * 0.5, uRadiusPx);
        p = position2 * scale;
        p += uCenterPx;
        vLocal = position2;
    }

    // Pixel → NDC
    float x = (p.x / uWinSize.x) * 2.0 - 1.0;
    float y = 1.0 - (p.y / uWinSize.y) * 2.0;

    gl_Position = vec4(x, y, 0.0, 1.0);
}

