package main

import (
	"math"
	"net"

	dati "p33/DATI"
	"p33/gl41"
	"p33/json"
	"p33/keys"
	"p33/mandel"
	xxx1 "p33/mandel/datiX"

	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math/big"
	"math/cmplx"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	Kw = 1000       //3504
	Kh = Kw * 3 / 4 //* 9 / 16 //
	//2
	winSize  = Kw * 1
	winSizeh = Kh * 1
	//
	blockSize = 50
	maxIter   = 250
	wsz1      = Kw * 1
	hsz1      = Kh * 1
)

var (
	Conn       *net.UDPConn
	LocalAddr  *net.UDPAddr
	RemoteAddr *net.UDPAddr
	Errore     error
	Flag       bool = false
	Flag2      bool = false
	F1         bool = false
	F2         bool = false
	FlagRun         = true
	Mode       int  = 0
	ch              = make(chan string) // canale globale
	conn       *net.UDPConn
	remoteAddr *net.UDPAddr
	mouseChan  = make(chan [2]float64, 10) // buffer di eventi mouse
	prec       = uint(256)
	rd         = new(big.Float).SetPrec(prec)
	cx         = new(big.Float).SetPrec(prec)
	cy         = new(big.Float).SetPrec(prec)
	tex        uint32
)

type Message struct {
	Data []byte
	Addr *net.UDPAddr
}

func printMem() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v\n",
		m.Alloc/1024/1024,
		m.TotalAlloc/1024/1024,
		m.Sys/1024/1024,
		m.NumGC)
}

// import mUdp "p3/mUdp"

// struttura per blocchi random
type TspiralPosition struct {
	count    int
	xDir     int
	yDir     int
	xInc     int
	yInc     int
	xPos     int
	yPos     int
	xa       int
	ya       int
	blkcnt   int
	flag     bool
	NotFirst bool
}

var sp TspiralPosition

var (
	bx      int
	by      int
	Counter int32
)

var pixels []uint32

// genera colori semplici in base al numero di iterazioni
func MandelbrotColor(cx, cy float64) uint32 {
	z := complex(0, 0)
	c := complex(cx, cy)
	var n int
	for n = 0; n < maxIter; n++ {
		if cmplx.Abs(z) > 2 {
			break
		}
		z = z*z + c //*k
	}
	// mapping colore (gradiente blu)
	if n == maxIter {
		return 0xFF008000 // verdino
		// return 0xFF000000 // nero
	}
	shade := uint32(10 + (n*255)/maxIter)
	// shade := uint32(255 - (n*255)/maxIter)
	return (0xFF << 24) | (shade << 16) | (shade << 8) | 0x1f
}

// calcola un blocco 8x8 di Mandelbrot
const KeachPixel = 5

func drawBlock(px, py int) {

	Vblkm := mandel.Tblkm{}

	Vblkm.Ny = blockSize
	Vblkm.Hy = make([]float64, blockSize)

	for y := 0; y < blockSize-0; y += KeachPixel {

		iy := py*blockSize + y
		if iy >= winSizeh || iy < 0 {
			continue
		}

		Vblkm.Iy = y
		// Vblkm.Fprimox = true

		for x := 0; x < blockSize-0; x += KeachPixel {

			Vblkm.Ix = x

			ix := px*blockSize + x
			if ix >= winSize || ix < 0 {
				continue
			}

			color := uint32(0)
			switch Mode {
			case 0:
				//float32
				color = mandel.MandelColor1(ix, iy, winSize*1.0, winSizeh)
			case 1:
				//bigfloat
				color = mandel.MandelColor2(ix, iy, winSize*1.0, winSizeh, &Vblkm)
			case 2:
				//float32
				color = mandel.PaintColor1(ix, iy, winSize*1.0, winSizeh)
			}

			idx := iy*winSize + ix
			idxm := len(pixels)
			if idx >= idxm || idx < 0 {
				fmt.Println("supero limiti: ", idx, idxm)
			} else {
				pixels[(winSizeh-iy-1)*winSize+ix] = color
			}

		}
	}
}

// // struttura per blocchi random
// type BlockPicker struct {
// 	blocks [][2]int
// }

// func newBlockPicker(width, height, blockSize int) *BlockPicker {
// 	blocksX := width / blockSize
// 	blocksY := height / blockSize
// 	total := blocksX * blocksY
// 	blocks1xx := make([][2]int, total)
// 	k := 0
// 	for by := 0; by < blocksY; by += 1 {
// 		for bx := 0; bx < blocksX; bx += 1 {
// 			blocks1xx[k] = [2]int{bx, by}
// 			k++
// 		}
// 	}
// 	rand.Seed(time.Now().UnixNano())
// 	rand.Shuffle(len(blocks1xx), func(i, j int) {
// 		blocks1xx[i], blocks1xx[j] = blocks1xx[j], blocks1xx[i]
// 	})
// 	return &BlockPicker{blocks: blocks1xx}
// }

// func (bp *BlockPicker) Next() (bx, by int, ok bool) {
// 	if len(bp.blocks) == 0 {
// 		return 0, 0, false
// 	}
// 	b := bp.blocks[len(bp.blocks)-1]
// 	bp.blocks = bp.blocks[:len(bp.blocks)-1]
// 	return b[0], b[1], true
// }

var mu sync.Mutex

func (sp *TspiralPosition) Next1() (x, y int, ok bool, kn int) {
	mu.Lock()
	defer mu.Unlock()

	if sp.NotFirst {
		if sp.flag { //first time must be true
			switch sp.count & 3 {
			case 0:
				sp.xInc += 1
				sp.xPos += sp.xInc
			case 1:
				sp.yInc += 1
				sp.yPos -= sp.yInc
			case 2:
				sp.xInc += 1
				sp.xPos -= sp.xInc
			case 3:
				sp.yInc += 1
				sp.yPos += sp.yInc
			}
			sp.flag = false
		}

		switch sp.count & 3 {
		case 0:
			sp.xa += 1
			if sp.xa >= sp.xPos {
				sp.flag = true
				sp.count += 1
			}
		case 1:
			sp.ya -= 1
			if sp.ya <= sp.yPos {
				sp.flag = true
				sp.count += 1
			}
		case 2:
			sp.xa -= 1
			if sp.xa <= sp.xPos {
				sp.flag = true
				sp.count += 1
			}
		case 3:
			sp.ya += 1
			if sp.ya >= sp.yPos {
				sp.flag = true
				sp.count += 1
			}
		}
	} else {
		sp.NotFirst = true
	}
	x1 := winSize/blockSize/2 + sp.xa
	y1 := winSizeh/blockSize/2 + sp.ya

	kw := int(math.Max(winSize, winSizeh))
	// kn = (winSize / blockSize) * (winSizeh / blockSize) // (conta 4 ogni blocchetto)
	kn = (kw/blockSize)*(kw/blockSize) + (kw / blockSize) // (conta 4 ogni blocchetto)
	return x1, y1, sp.blkcnt < kn, kn
}

func init() {
	runtime.LockOSThread()
}

// var mu1 sync.Mutex

func externWork() {

	// goroutine che calcola blocchi progressivamente
	go func() {

		drawBlock(bx, by)

		for {

			if !FlagRun {
				break
			}

			bx, by, ok, kn := sp.Next1()
			_ = kn
			if !ok {
				if !ok {
					break
				}
			}
			drawBlock(bx, by)
			sp.blkcnt += 1

			time.Sleep(10 * time.Microsecond) // regola la velocità
		}
		atomic.AddInt32(&Counter, -1)
	}()
}

func Uint32ToImage(pixels []uint32, width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		// for y := height - 1; y > 0; y-- {
		for x := 0; x < width; x++ {
			idx := y*width + x
			p := pixels[idx]

			a := uint8(p >> 24)
			b := uint8(p >> 16)
			g := uint8(p >> 8)
			r := uint8(p)

			img.Set(x, height-y-1, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}
	return img
}

func repaint() {

	if !F1 {
		return
	}
	F1 = false
	fmt.Println("REPAINT")

	mu.Lock()
	sp.blkcnt = 1
	mu.Unlock()

	atomic.StoreInt32(&Counter, 10)
	// atomic.StoreInt32(&Counter, 1)
	externWork()
	externWork()
	externWork()
	externWork()
	externWork()
	externWork()
	externWork()
	externWork()
	externWork()
	externWork()

	go func() {
		for {
			v := atomic.LoadInt32(&Counter)
			if v <= 0 {
				fmt.Println("FATTO.")
				if FlagRun {
					image := Uint32ToImage(pixels, winSize/1, winSizeh/1)
					f1, _ := os.Create("mandelbrot_shaded2.jpg")
					defer f1.Close()
					jpeg.Encode(f1, image, &jpeg.Options{Quality: 90})
					fmt.Println("salvato.")

					f2, _ := os.Create("mandelbrot_shaded2.png")
					defer f2.Close()
					png.Encode(f2, image)
					fmt.Println("salvato.")
				} else {
					fmt.Println("INTERROTTO")
					atomic.StoreInt32(&Counter, 0)
				}
				break
			}
		}
	}()

	if !FlagRun {
		FlagRun = true
	}
	fmt.Println("\n\n\nlen(pixels): ", len(pixels))

}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {

		case glfw.Key5:
			gl41.AddLine(float32(mouseX), float32(mouseY), 50, 90, 1)
		case glfw.Key6:
			gl41.RemoveLine(0)

		case glfw.Key7:
			// nascondi tutte le guide
			gl41.HideGroup(1)
		case glfw.Key8:
			// mostra gruppo 1
			gl41.ShowGroup(1)

		case glfw.Key9:
			// elimina gruppo 1 definitivamente
			gl41.RemoveGroup(1)

		case glfw.Key1:
			cfg, err := json.LoadConfig("config.json")
			_ = err
			json.Cfg = *cfg
			// fmt.Print("\033[20;1H", json.Cfg)
			fmt.Println(json.Cfg)

		case glfw.Key2:
			err := json.SaveConfig("config.json", &json.Cfg)
			_ = err

		case glfw.KeyEscape:
			w.SetShouldClose(true) // chiude finestra

		case glfw.KeySpace:
			log.Println("Hai premuto SPAZIO")
			printMem()
			FlagRun = !FlagRun

		case glfw.KeyUp:
			log.Println("Freccia SU")

		case glfw.KeyDown:
			log.Println("Freccia GIÙ")

		case glfw.KeyA:
			log.Println("tasto A on")
			Mode = 0
			F1 = true

			sp = TspiralPosition{}
			bx = winSize / blockSize / 2
			by = winSizeh / blockSize / 2
			_, _ = bx, by

			repaint()

		case glfw.KeyB:
			log.Println("tasto B on")
			Mode = 1
			F1 = true

			sp = TspiralPosition{}

			bx = winSize / blockSize / 2
			by = winSizeh / blockSize / 2

			repaint()

		case glfw.KeyC:
			log.Println("tasto C on")
			Flag2 = true
			clear(pixels)

		case glfw.KeyD: // coloarione painting sinusoidale
			log.Println("tasto B on")
			Mode = 2
			F1 = true

			sp = TspiralPosition{}

			bx = winSize / blockSize / 2
			by = winSizeh / blockSize / 2

			repaint()

		case glfw.KeyX:
			// salva cx,cy,rd e calcola nuovi in base ai click mouse salvati
			log.Println("tasto X: nuovo cx,cy,rd")
			fmt.Println(Msx, Msy, Mdx, Mdy)

		}
	}
	if action == glfw.Release {
		switch key {
		case glfw.KeyEscape:
			w.SetShouldClose(true) // chiude finestra
		case glfw.KeySpace:
			log.Println("SPAZIO off")
		case glfw.KeyUp:
			log.Println("Freccia SU  off")
		case glfw.KeyDown:
			log.Println("Freccia GIÙ  off")
		case glfw.KeyA:
			log.Println("tasto A off")
		}
	}
}

var (
	//mouse destro x,y per calcolo raggio dal centro
	Mdx int
	Mdy int
	//mouse sinistro x,y per centro x,y
	Msx int
	Msy int
)

var (
	SNuovoCentroX string
	SNuovoCentroY string
	SNuovoRaggio  string
)

func CalcolaNuovoCentro() {
	var (
		prec = uint(300)
		rd   = new(big.Float).SetPrec(prec)
		cx   = new(big.Float).SetPrec(prec)
		cy   = new(big.Float).SetPrec(prec)
		// step  = new(big.Float).SetPrec(prec)
		// delta = new(big.Float).SetPrec(prec)
		// bdist = new(big.Float).SetPrec(prec)
		// zx    = new(big.Float).SetPrec(prec)
		// zy    = new(big.Float).SetPrec(prec)
		// zx2   = new(big.Float).SetPrec(prec)
		// zy2   = new(big.Float).SetPrec(prec)
		// B2    = new(big.Float).SetPrec(prec).SetInt64(2)
		// // B1    = new(big.Float).SetPrec(prec).SetInt64(1)
		// // B3    = new(big.Float).SetPrec(prec).SetInt64(3)
		// eps = new(big.Float).SetPrec(prec)
	)

	rd.SetString(xxx1.Srd)
	cx.SetString(xxx1.Scx)
	cy.SetString(xxx1.Scy)
}

func CalcolaNuovoRaggio() {

}

func mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch button {
		case glfw.MouseButtonLeft:
			log.Println("Click sinistro on")
			Msx = dati.MouseX
			Msy = dati.MouseY
			CalcolaNuovoCentro()
		case glfw.MouseButtonRight:
			log.Println("Click destro on")
			Mdx = dati.MouseX
			Mdy = dati.MouseY
			CalcolaNuovoRaggio()
		}
	}
	if action == glfw.Release {
		switch button {
		case glfw.MouseButtonLeft:
			log.Println("Click sinistro off")
		case glfw.MouseButtonRight:
			log.Println("Click destro off")
		}
	}
}

func mouseScrollCallback(w *glfw.Window, xoff, yoff float64) {
	// log.Printf("Scroll: (%.1f, %.1f)\n", xoff, yoff)

}

var (
	mouseX int32
	mouseY int32
)

var mouseInside bool = true

func cursorEnterCallback(w *glfw.Window, entered bool) {
	mouseInside = entered // true se il mouse è dentro, false se è uscito
	_ = mouseInside
}

func mousePosCallback(w *glfw.Window, xpos, ypos float64) {

	mouseX = int32(xpos) // - winSize/2)
	mouseY = int32(ypos) //(winSizeh - int32(ypos)) //+winSizeh/2))
	// _, _ = mouseX, mouseY

	dati.MouseX = int(mouseX)
	dati.MouseY = int(mouseY)

	keys.ShowPos(xpos, ypos)

	select {
	case mouseChan <- [2]float64{xpos, ypos}:
	default:
	}
}

func main() {

	// fmt.Printf("\033[2J")
	fmt.Printf("\n\n-- start p3 in floder p33 --\n\n")

	json.IniziaJson()

	//
	Srd := ""
	Scx := ""
	Scy := ""
	SMaxiter := 100.0
	rd.SetString(Srd)
	cx.SetString(Scx)
	cy.SetString(Scy)
	Maxiter := SMaxiter
	_ = Maxiter

	localPort := "9002"
	remotePort := "9001"
	localAddr, err := net.ResolveUDPAddr("udp", ":"+localPort)
	if err != nil {
		panic(err)
	}
	remoteAddr, err = net.ResolveUDPAddr("udp", "localhost:"+remotePort)
	if err != nil {
		panic(err)
	}
	conn, err = net.ListenUDP("udp", localAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	text := "qui0001\n"
	_, err = conn.WriteToUDP([]byte(text), remoteAddr)
	if err != nil {
		fmt.Println("Errore invio:", err)
		return
	}

	// Goroutine per ricevere messaggi
	go func() {
		buf := make([]byte, 1024)
		for {
			n, addr, err := conn.ReadFromUDP(buf)
			_ = addr
			if err != nil {
				fmt.Println("Errore ricezione:", err)
				return
			}
			s := string(buf[:n])
			ch <- s // evento nel canale
		}
	}()

	sp = TspiralPosition{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, true, false}
	bx = winSize / blockSize / 2
	by = winSizeh / blockSize / 2

	// window, err := glfw.CreateWindow(winSize, winSizeh, "Mandelbrot progressivo", nil, nil)
	wws := wsz1
	whs := hsz1
	// gl41.SetWinsize(winSize, winSizeh)
	// init GLFW

	err1 := glfw.Init()
	if err1 != nil {
		log.Fatalln("failed to init glfw:", err1)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// glfw.WindowHint(glfw.Samples, 4) // <-- MSAA attivo 4x
	glfw.WindowHint(glfw.Samples, 8) // <-- MSAA attivo 8x

	window, err := glfw.CreateWindow(wws, whs, "Mandelbrot progressivo", nil, nil)

	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL versione:", version)

	shaderv := gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	fmt.Println("GLSL:", shaderv)

	gl.Enable(gl.MULTISAMPLE)

	gl41.InitGL41() // <--- INIZIALIZZA pipeline moderna

	// setup viewport
	// gl.Viewport(0, 0, winSize/1, winSizeh/1)
	gl.Viewport(0, 0, wsz1, hsz1)

	gl41.SetWindowSize(winSize, winSizeh, wsz1, hsz1)
	// gl41.SetWindowSize(wsz1, hsz1, wsz1, hsz1)

	// pixel buffer
	pixels = make([]uint32, winSize*winSizeh)

	tex = gl41.CreateTexture()

	repaint()

	window.SetKeyCallback(keyCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetScrollCallback(mouseScrollCallback)
	window.SetCursorPosCallback(mousePosCallback)
	window.SetCursorEnterCallback(cursorEnterCallback)

	//-----------
	overlayImg := keys.CreateOverlay(winSize, winSizeh) // con parametri dinamici

	// // fmt.Printf("%X\n", overlayImg.(*image.RGBA).Pix)
	data := keys.RGBAtoUint32Fast(overlayImg.(*image.RGBA))
	_ = data

	var texOverlay uint32
	gl.GenTextures(1, &texOverlay)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	//-----------
	oneTime1 := 1
	_ = oneTime1

	// main loop
	var value string = ""
	// Controlla il massimo supportato
	var maxTexSize int32
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &maxTexSize)
	fmt.Println("Max texture size =", maxTexSize)

	// --- Loop principale ---
	go func() {
		var count int = 1000
		for {
			count += 1
			pos := <-mouseChan
			msg := fmt.Sprintf("x=%.2f y=%.2f cnt:%d---\n", pos[0], pos[1], count)
			// conn.Write([]byte(msg))
			_, _ = conn.WriteToUDP([]byte(msg), remoteAddr)
			// fmt.Println(msg)
		}
	}()

	gl41.SetGroupColor(1, 1, 0, 0, 1) // rosso
	gl41.SetGroupColor(2, 0, 1, 0, 1) // verde
	gl41.SetGroupColor(3, 0, 0, 1, 1) // blu

	gl41.AddLine(100, 100, 500, 400, 2)
	gl41.AddLine(200, 800, 900, 200, 2)
	gl41.AddLine(float32(mouseX), float32(mouseY), 1000, 1000, 2)

	for !window.ShouldClose() {

		// --- EVENTI UDP (non blocca)
		select {
		case val := <-ch:
			value = val
			fmt.Println(val, "Ho letto:", value)
		default:
		}

		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl41.RenderFullscreen(tex)
		// gl41.RenderFullscreen(tex, wsz1, hsz1)

		// Disegna la texture fullscreen
		gl41.UpdateTexture(tex, pixels)
		// r := math.Min(float64(winSize), float64(winSizeh)) / 20
		r := math.Min(float64(wsz1), float64(hsz1)) / 20
		mx := float64(mouseX) // / float64(wsz1) //* winSize
		my := float64(mouseY) // / float64(hsz1) //* winSizeh
		gl41.DrawOverlay(float32(mx), float32(my), float32(r))

		gl41.RenderLines()

		// ----------------------------------------------------

		window.SwapBuffers()
		glfw.PollEvents()
	}

}
