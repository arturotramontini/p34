// #version 410 core

// out vec4 fragColor;

// uniform vec4 uColor;

// void main() {
//     fragColor = uColor;
// }


// #version 150 core

// out vec4 fragColor;

// uniform vec4 uColor;

// void main() {
//     fragColor = uColor;
// }

#version 410 core

out vec4 fragColor;

uniform vec4 uColor;

void main() {
    fragColor = uColor;
		// fragColor = vec4(1,0,0,1);   // TUTTO rosso
}
