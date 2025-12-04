package mandel

import (
	"math"
	"math/big"
	xxx1 "p33/mandel/datiX"
	"strconv"
)

var Test1 float64 = 123.456

func NewTMandel(prec uint) *TMandel {
	return &TMandel{
		cx:   new(big.Float).SetPrec(prec),
		cy:   new(big.Float).SetPrec(prec),
		zx:   new(big.Float).SetPrec(prec),
		zy:   new(big.Float).SetPrec(prec),
		zx2:  new(big.Float).SetPrec(prec),
		zy2:  new(big.Float).SetPrec(prec),
		dist: new(big.Float).SetPrec(prec),
		step: new(big.Float).SetPrec(prec),
		Two:  new(big.Float).SetPrec(prec).SetInt64(2),
	}
}

type TMandel struct {
	cx   *big.Float
	cy   *big.Float
	zx   *big.Float
	zy   *big.Float
	zx2  *big.Float
	zy2  *big.Float
	dist *big.Float
	Two  *big.Float
	step *big.Float
}
type Trxy struct {
	rd *big.Float
	cx *big.Float
	cy *big.Float
}

func NewTrxy(prec uint) *Trxy {
	return &Trxy{
		rd: new(big.Float).SetPrec(prec),
		cx: new(big.Float).SetPrec(prec),
		cy: new(big.Float).SetPrec(prec),
	}
}

var (
	Brxy         = NewTrxy(128)
	XBrxy        = NewTrxy(128)
	Perturb  int = 0
	Bprec        = uint(256)
	Brd          = new(big.Float).SetPrec(Bprec)
	Bcx          = new(big.Float).SetPrec(Bprec)
	Bcy          = new(big.Float).SetPrec(Bprec)
	BMaxiter int = 100
)

func GetRXY(ix, iy, ws, hs int) (string, string, string) {

	prec := uint(128)

	rd := new(big.Float).SetPrec(prec)
	cx := new(big.Float).SetPrec(prec)
	cy := new(big.Float).SetPrec(prec)
	step := new(big.Float).SetPrec(prec)
	delta := new(big.Float).SetPrec(prec)

	rd.SetString(xxx1.Srd)
	cx.SetString(xxx1.Scx)
	cy.SetString(xxx1.Scy)

	// // prova
	// kv := new(big.Float).SetPrec(prec)
	// kv.SetFloat64(10)
	// rd.Mul(rd, kv)
	// //-------

	min := math.Min(float64(ws), float64(hs)) / 2

	step.SetFloat64(min)
	step.Quo(rd, step)

	// eps = step.Quo(step, BTwo)
	// eps.Quo(step, BTwo)

	// deltaX
	d := ix - ws/2
	delta.SetFloat64(float64(d))
	delta.Mul(delta, step)
	cx.Add(cx, delta)

	// deltaY
	d = iy - hs/2
	delta.SetFloat64(float64(d))
	delta.Mul(delta, step)
	cy.Add(cy, delta)

	srd := rd.Text('g', int(prec))
	scx := cx.Text('g', int(prec))
	scy := cy.Text('g', int(prec))

	return srd, scx, scy
}

func MandelColor2_xxx2(ix, iy, ws, hs, mode int) float64 {

	var (
		prec  = uint(300)
		rd    = new(big.Float).SetPrec(prec)
		cx    = new(big.Float).SetPrec(prec)
		cy    = new(big.Float).SetPrec(prec)
		step  = new(big.Float).SetPrec(prec)
		delta = new(big.Float).SetPrec(prec)
		bdist = new(big.Float).SetPrec(prec)
		zx    = new(big.Float).SetPrec(prec)
		zy    = new(big.Float).SetPrec(prec)
		zx2   = new(big.Float).SetPrec(prec)
		zy2   = new(big.Float).SetPrec(prec)
		B2    = new(big.Float).SetPrec(prec).SetInt64(2)
		// B1    = new(big.Float).SetPrec(prec).SetInt64(1)
		// B3    = new(big.Float).SetPrec(prec).SetInt64(3)
		eps = new(big.Float).SetPrec(prec)
	)

	rd.SetString(xxx1.Srd)
	cx.SetString(xxx1.Scx)
	cy.SetString(xxx1.Scy)

	// // provato per immagine quadrata
	// kv := new(big.Float).SetPrec(prec)
	// kv.SetFloat64(20) //20 è per numero di  piastrelle per 1 lato
	// rd.Quo(rd, kv)    // ovvero raggio della piastrella
	// // rd.Mul(rd, B2)    // ovvero raggio per numero piastrelle
	// //-------

	min := math.Min(float64(ws), float64(hs)) / 2

	step.SetFloat64(min)
	step.Quo(rd, step)

	// eps = step.Quo(step, BTwo)
	// eps.Quo(step, B1)
	eps = step

	// deltaX
	d := ix - ws/2
	delta.SetFloat64(float64(d))
	delta.Mul(delta, step)
	cx.Add(cx, delta)
	if mode == 1 {
		cx.Add(cx, eps)
	}

	// deltaY
	d = iy - hs/2
	delta.SetFloat64(float64(d))
	delta.Mul(delta, step)
	cy.Add(cy, delta)
	if mode == 2 {
		cy.Sub(cy, eps)
	}

	zx.SetInt64(0)
	zy.SetInt64(0)

	mxi := int(xxx1.SMaxiter)
	n := 0
	_ = n
	dist := 0.0
	_ = dist

	escapeRadius := 1000000.0

	for n = range mxi {
		zx2.Mul(zx, zx)
		zy2.Mul(zy, zy)
		bdist.Add(zx2, zy2)

		dist, _ = bdist.Float64()
		_ = dist
		if dist >= escapeRadius {
			break
		}

		// zy = 2*zx*zy + cy
		zy.Mul(zx, zy)
		zy.Mul(zy, B2)
		zy.Add(zy, cy)

		// zx = zx2 - zy2 + cx
		zx.Sub(zx2, zy2)
		zx.Add(zx, cx)
	}

	mxi = mxi - 2
	iteration := n
	_ = iteration

	zxf, _ := zx.Float64()
	zyf, _ := zy.Float64()

	if iteration >= mxi { // per escludere momentaneamente
		zxf = 100000.0 / zxf
		zyf = 100000.0 / zyf
	} else {
		zxf = 1 * zxf
		zyf = 1 * zyf
	}

	// https://www.youtube.com/watch?v=zmWkhlocBRY&t=636s
	// minuto 24:23
	distf := math.Sqrt(zxf*zxf + zyf*zyf)

	r := float64(escapeRadius)
	fractIter := math.Log(distf) / math.Log(r)
	fractIter = math.Log(fractIter) / math.Log(2)
	iter := float64(iteration) - fractIter
	mu1 := math.Sqrt(iter) / float64(mxi)
	_ = mu1

	return mu1
}

// *
// *
// *
// *
// *
// var Hp, h, hx, hy float64 = 0.0, 0.0, 0.0, 0.0

type Tblkm struct {
	Prova   float64
	Ny      int
	Iy      int
	Ix      int
	Hy      []float64
	Hx      float64
	Fprimox bool
}

// var Vblkm Tblkm

func MandelColor2(ix, iy, ws, hs int, Vblkm *Tblkm) uint32 {

	h := 0.0
	hx := 0.0
	hy := 0.0

	//--- per calcolare solo una volta invece di tre
	//--- vale solo quando eps è multiplo di step ---------
	if Vblkm.Ix == 0 {
		// Vblkm.Fprimox = false
		h = MandelColor2_xxx2(ix, iy, ws, hs, 0)
	} else {
		h = Vblkm.Hx
	}

	// hx è il pixel X successivo
	hx = MandelColor2_xxx2(ix, iy, ws, hs, 1) // modo 1 : ix + 1
	Vblkm.Hx = hx

	// la rima riga la calcola tutta
	if Vblkm.Iy == 0 {
		hy = MandelColor2_xxx2(ix, iy, ws, hs, 2) // modo 2 : iy - 1
		Vblkm.Hy[Vblkm.Ix] = h
	} else {
		hy = Vblkm.Hy[Vblkm.Ix]
		Vblkm.Hy[Vblkm.Ix] = h
	}
	//---------------------

	// h = MandelColor2_xxx2(ix, iy, ws, hs, 0)
	// hx = MandelColor2_xxx2(ix, iy, ws, hs, 1)
	// hy = MandelColor2_xxx2(ix, iy, ws, hs, 2)

	light := [3]float64{.5, +0.04, .5}

	// normalizza
	lLen := math.Sqrt(light[0]*light[0] + light[1]*light[1] + light[2]*light[2])
	for i := 0; i < 3; i++ {
		light[i] /= lLen
	}

	// gradiente
	eps := 1.0e-18
	dx := 1.0 * (hx - h) / eps / 1
	dy := 1.0 * (hy - h) / eps / 1

	// // gradiente
	// eps := 1.0e-9
	// dx := 1.0 * (hx - h) / eps / 1
	// dy := 1.0 * (hy - h) / eps / 1

	// dx := h - h/100
	// dy := h + h/100
	// dx := (h - Hp) / 1e-10
	// dy := (h - Hp) / 1e-10

	// normale
	// n := [3]float64{1.0, -dx, -dy}
	n := [3]float64{-dx, 1.0, -dy}
	nLen := math.Sqrt(n[0]*n[0] + n[1]*n[1] + n[2]*n[2])
	for i := 0; i < 3; i++ {
		n[i] /= nLen
	}

	// intensità luce
	intensity := n[0]*light[0] + n[1]*light[1] + n[2]*light[2]
	// if intensity < 0 {
	// 	// intensity = 0 - intensity
	// 	intensity = math.Abs(math.Log(intensity))
	// }
	intensity = math.Abs(math.Log(math.Abs(intensity)))
	// intensity = float64(h) / float64(maxIter)

	alpha := uint8(255)
	_ = alpha

	// intensity = 1
	h *= 201
	red := uint32((math.Sin(0.30*h*220)*0.5 + 0.5) * 255.0 * intensity * 1.0)
	gre := uint32((math.Sin(0.45*h*220)*0.5 + 0.5) * 255.0 * intensity * 1.0)
	blu := uint32((math.Sin(0.65*h*220)*0.5 + 0.5) * 255.0 * intensity * 1.0)

	color := (uint32(0xFF) << 24) | ((red & 0xFF) << 16) | ((gre & 0xFF) << 8) | (blu & 0xff)
	return color
}

// *
// *
// *
// *
// *

func MandelColor1_xxx1(ix, iy, ws, hs, mode int) (float64, float64) {

	rd, _ := strconv.ParseFloat(xxx1.Srd, 64)
	cx, _ := strconv.ParseFloat(xxx1.Scx, 64)
	cy, _ := strconv.ParseFloat(xxx1.Scy, 64)

	min := math.Min(float64(ws), float64(hs)) / 2

	step := min
	step = rd / step

	eps := step / 1
	// eps := step / 2

	// deltaX
	d := ix - ws/2
	delta := float64(d) * step
	cx = cx + delta
	if mode == 1 {
		cx = cx + eps
	}

	// deltaY
	d = iy - hs/2
	delta = float64(d) * step
	cy = cy + delta
	if mode == 2 {
		cy = cy - eps
	}

	//----------

	zx, zy := 0.0, 0.0

	maxIter := int(xxx1.SMaxiter)
	mxi := maxIter
	n := 0
	zx2 := 0.0
	zy2 := 0.0
	dist := 0.0
	escapeRadius := 1000000.0

	// cx = 1.0 / cx
	// cy = 1.0 / cy

	for n = range mxi {

		zx2 = zx * zx
		zy2 = zy * zy
		dist = zx2 + zy2
		if dist > escapeRadius {
			break
		}

		zy = 2*zx*zy + cy

		zx = zx2 - zy2 + cx
	}

	iteration := n
	mxi -= 2

	if iteration >= mxi { // per escludere momentaneamente

		zx = 1000.0 / zx
		zy = 1000.0 / zy
	} else {
		zx = 1 * zx
		zy = 1 * zy
		// zx = 1e6 / zx
		// zy = 1e6 / zy
	}

	// https://www.youtube.com/watch?v=zmWkhlocBRY&t=636s
	// minuto 24:23
	dist = math.Sqrt(zx*zx + zy*zy)

	r := float64(escapeRadius)
	fractIter := math.Log(dist) / math.Log(r)
	fractIter = math.Log(fractIter) / math.Log(2)
	iter := float64(iteration) - fractIter
	mu1 := math.Sqrt(iter) / float64(maxIter)
	_ = mu1

	return mu1, (float64(iteration) / float64(mxi))

}

func MandelColor1(ix, iy, ws, hs int) uint32 {

	h, iv := MandelColor1_xxx1(ix, iy, ws, hs, 0)
	hx, _ := MandelColor1_xxx1(ix, iy, ws, hs, 1)
	hy, _ := MandelColor1_xxx1(ix, iy, ws, hs, 2)

	_ = iv

	// luce da sinistra/alto
	light := [3]float64{.5, +0.04, .5}

	// normalizza
	lLen := math.Sqrt(light[0]*light[0] + light[1]*light[1] + light[2]*light[2])
	for i := 0; i < 3; i++ {
		light[i] /= lLen
	}

	// gradiente
	// eps := 1.0e-9
	eps := 1.0e-18
	dx := 1.0 * (hx - h) / eps / 1
	dy := 1.0 * (hy - h) / eps / 1

	// normale
	n := [3]float64{-dx, 1.0, -dy}
	nLen := math.Sqrt(n[0]*n[0] + n[1]*n[1] + n[2]*n[2])
	for i := 0; i < 3; i++ {
		n[i] /= nLen
	}

	// intensità luce
	intensity := n[0]*light[0] + n[1]*light[1] + n[2]*light[2]
	// if intensity < 0 {
	// 	// intensity = 0 - intensity
	// 	intensity = math.Abs(math.Log(intensity))
	// }
	intensity = math.Abs(math.Log(math.Abs(intensity)))
	// intensity = float64(h) / float64(maxIter)

	// // intensity *= iv
	// if iv >= 1.0 {
	// 	intensity = 1.0
	// 	h = 1.0
	// } else {
	// 	intensity = 0.0
	// }

	// // colore base (scala di grigi)
	alpha := uint8(255)
	_ = alpha
	h *= 201
	red := uint32((math.Sin(0.30*h*120)*0.5 + 0.5) * 255.0 * intensity * 1.0)
	gre := uint32((math.Sin(0.45*h*120)*0.5 + 0.5) * 255.0 * intensity * 1.0)
	blu := uint32((math.Sin(0.65*h*120)*0.5 + 0.5) * 255.0 * intensity * 1.0)

	// var color uint32
	color := (uint32(0xFF) << 24) | ((red & 0xFF) << 16) | ((gre & 0xFF) << 8) | (blu & 0xff)
	return color
}

func PaintColor1(ix, iy, ws, hs int) uint32 {

	rhw := float64(hs) / float64(ws)
	// nx e ny vanno da 0 a 1
	nx := float64(ix) / float64(ws)
	ny := float64(iy) / float64(hs)
	nxy := nx * ny
	_ = nxy

	intensity1 := 0.7
	intensity2 := 0.7

	// hx := (math.Sin(2.0*math.Pi*8.0*nxy)*0.45 + 0.5) * intensity1
	// hy := (math.Cos(2.0*math.Pi*8.0*nxy)*0.45 + 0.5) * intensity1

	freq := 9.0 * 8

	hx := (math.Cos(freq*2.0*math.Pi*nx)*-0.5 + 0.5 + 0.0) * intensity1
	hy := (math.Cos(freq*2.0*math.Pi*ny*rhw)*-0.5 + 0.5 + 0.0) * intensity1

	hxy := hx * hy

	red := uint32(hxy * 255.0 * intensity2)
	gre := uint32(hxy * 255.0 * intensity2)
	blu := uint32(hxy * 255.0 * intensity2)

	// var color uint32
	color := (uint32(0xFF) << 24) | ((red & 0xFF) << 16) | ((gre & 0xFF) << 8) | (blu & 0xff)
	return color
}
