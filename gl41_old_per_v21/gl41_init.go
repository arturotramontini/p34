package gl41

import "github.com/go-gl/gl/v4.1-core/gl"

var GlobalVAO uint32

func InitGL41() {
	gl.GenVertexArrays(1, &GlobalVAO)
	gl.BindVertexArray(GlobalVAO)
}
